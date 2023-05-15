package runner

import (
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/hako/durafmt"

	"github.com/projectdiscovery/gologger"

	"github.com/chainreactors/urlfounder/v2/pkg/resolve"
	"github.com/chainreactors/urlfounder/v2/pkg/subscraping"
)

const maxNumCount = 2

// EnumerateSingleURL wraps EnumerateSingleURLWithCtx with an empty context
func (r *Runner) EnumerateSingleURL(domain string, writers []io.Writer) error {
	return r.EnumerateSingleURLWithCtx(context.Background(), domain, writers)
}

// EnumerateSingleURLWithCtx performs url enumeration against a single domain
func (r *Runner) EnumerateSingleURLWithCtx(ctx context.Context, domain string, writers []io.Writer) error {
	gologger.Info().Msgf("Enumerating urls for %s\n", domain)

	//Check if the user has asked to remove wildcards explicitly.
	//If yes, create the resolution pool and get the wildcards for the current domain
	var resolutionPool *resolve.ResolutionPool
	if r.options.RemoveWildcard {
		resolutionPool = r.resolverClient.NewResolutionPool(r.options.Threads, r.options.RemoveWildcard)
		err := resolutionPool.InitWildcards(domain)
		if err != nil {
			// Log the error but don't quit.
			gologger.Warning().Msgf("Could not get wildcards for domain %s: %s\n", domain, err)
		}
	}

	// Run the passive url enumeration
	now := time.Now()
	passiveResults := r.passiveAgent.EnumerateURLsWithCtx(ctx, domain, r.options.Proxy, r.options.RateLimit, r.options.Timeout, time.Duration(r.options.MaxEnumerationTime)*time.Minute)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// Create a unique map for filtering duplicate urls out 过滤重复子域
	uniqueMap := make(map[string]resolve.HostEntry)
	// Create a map to track sources for each host 跟踪host源
	sourceMap := make(map[string]map[string]struct{})
	// Process the results in a separate goroutine
	go func() {
		for result := range passiveResults {
			switch result.Type {
			case subscraping.Error:
				gologger.Warning().Msgf("Could not run source %s: %s\n", result.Source, result.Error)
			case subscraping.URL:
				// 验证找到的子域并删除通配符
				url := strings.ReplaceAll(strings.ToLower(result.Value), "*.", "")

				if matchURL := r.filterAndMatchURL(url); matchURL {
					if _, ok := uniqueMap[url]; !ok {
						sourceMap[url] = make(map[string]struct{})
					}

					// Log the verbose message about the found url per source
					if _, ok := sourceMap[url][result.Source]; !ok {
						gologger.Verbose().Label(result.Source).Msg(url)
					}

					sourceMap[url][result.Source] = struct{}{}

					// Check if the url is a duplicate. If not,
					// send the url for resolution.
					if _, ok := uniqueMap[url]; ok {
						continue
					}

					hostEntry := resolve.HostEntry{Host: url, Source: result.Source}

					uniqueMap[url] = hostEntry
					// If the user asked to remove wildcard then send on the resolve
					// queue. Otherwise, if mode is not verbose print the results on
					// the screen as they are discovered.
					if r.options.RemoveWildcard {
						resolutionPool.Tasks <- hostEntry
					}
				}
			}
		}
		// Close the task channel only if wildcards are asked to be removed
		if r.options.RemoveWildcard {
			close(resolutionPool.Tasks)
		}
		wg.Done()
	}()

	// If the user asked to remove wildcards, listen from the results
	// queue and write to the map. At the end, print the found results to the screen
	foundResults := make(map[string]resolve.Result)
	if r.options.RemoveWildcard {
		// Process the results coming from the resolutions pool
		for result := range resolutionPool.Results {
			switch result.Type {
			case resolve.Error:
				gologger.Warning().Msgf("Could not resolve host: %s\n", result.Error)
			case resolve.URL:
				// Add the found url to a map.
				if _, ok := foundResults[result.Host]; !ok {
					foundResults[result.Host] = result
				}
			}
		}
	}

	wg.Wait()
	outputWriter := NewOutputWriter(r.options.JSON)
	// Now output all results in output writers
	var err error
	for _, writer := range writers {
		if r.options.StatusCode && r.options.Title {
			err = outputWriter.WriteStatusCodeAndTitle(domain, foundResults, writer)
		} else if r.options.StatusCode {
			err = outputWriter.WriteStatusCode(domain, foundResults, writer)
		} else {
			if r.options.RemoveWildcard {
				err = outputWriter.WriteHostNoWildcard(domain, foundResults, writer)
			} else {
				if r.options.CaptureSources {
					err = outputWriter.WriteSourceHost(domain, sourceMap, writer)
				} else {
					err = outputWriter.WriteHost(domain, uniqueMap, writer)
				}
			}
		}
		if err != nil {
			gologger.Error().Msgf("Could not write results for %s: %s\n", domain, err)
			return err
		}
	}

	// Show found url count in any case.
	duration := durafmt.Parse(time.Since(now)).LimitFirstN(maxNumCount).String()
	var numberOfURLs int
	if r.options.RemoveWildcard {
		numberOfURLs = len(foundResults)
	} else {
		numberOfURLs = len(uniqueMap)
	}

	if r.options.ResultCallback != nil {
		for _, v := range uniqueMap {
			r.options.ResultCallback(&v)
		}
	}
	gologger.Info().Msgf("Found %d urls for %s in %s\n", numberOfURLs, domain, duration)

	if r.options.Statistics {
		gologger.Info().Msgf("Printing source statistics for %s", domain)
		printStatistics(r.passiveAgent.GetStatistics())
	}

	return nil
}

func (r *Runner) filterAndMatchURL(url string) bool {
	if r.options.filterRegexes != nil {
		for _, filter := range r.options.filterRegexes {
			if m := filter.MatchString(url); m {
				return false
			}
		}
	}
	if r.options.matchRegexes != nil {
		for _, match := range r.options.matchRegexes {
			if m := match.MatchString(url); m {
				return true
			}
		}
		return false
	}
	return true
}
