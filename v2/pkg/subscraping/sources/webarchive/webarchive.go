package webarchive

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chainreactors/urlfounder/v2/pkg/subscraping"
	"net/http"
	"time"
)

// Source is the passive scraping agent
type Source struct {
	timeTaken time.Duration
	errors    int
	results   int
	skipped   bool
}

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)
	s.errors = 0
	s.results = 0

	go func() {
		defer func(startTime time.Time) {
			s.timeTaken = time.Since(startTime)
			close(results)
		}(time.Now())

		headers := map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		}

		api := fmt.Sprintf("https://web.archive.org/web/timemap/json?url=%s&matchType=prefix&collapse=urlkey&output=json&fl=original%%2Cmimetype%%2Ctimestamp%%2Cendtimestamp%%2Cgroupcount%%2Cuniqcount&limit=1000", domain)
		resp, err := session.Get(ctx, api, "", headers)
		isForbidden := resp != nil && resp.StatusCode == http.StatusForbidden
		if err != nil {
			//println(err.Error())
			if !isForbidden {
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
				s.errors++
				session.DiscardHTTPResponse(resp)
			}
			return
		}

		var res [][]string
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			s.errors++
			resp.Body.Close()
			return
		}
		resp.Body.Close()

		for _, r := range res[1:] {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: r[0]}
			s.results++
		}

	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "webarchive"
}

func (s *Source) IsDefault() bool {
	return false
}

func (s *Source) HasRecursiveSupport() bool {
	return false
}

func (s *Source) AddApiKeys(keys []string) {

}

func (s *Source) NeedsKey() bool {
	return true
}

func (s *Source) Statistics() subscraping.Statistics {
	return subscraping.Statistics{
		Errors:    s.errors,
		Results:   s.results,
		TimeTaken: s.timeTaken,
		Skipped:   s.skipped,
	}
}
