package bevigil

import (
	"context"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/chainreactors/urlfounder/v2/pkg/subscraping"
)

type Response struct {
	Domain string   `json:"domain"`
	Urls   []string `json:"urls"`
}

type Source struct {
	apiKeys   []string
	timeTaken time.Duration
	errors    int
	results   int
	skipped   bool
}

func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)
	s.errors = 0
	s.results = 0

	go func() {
		defer func(startTime time.Time) {
			s.timeTaken = time.Since(startTime)
			close(results)
		}(time.Now())

		randomApiKey := subscraping.PickRandom(s.apiKeys, s.Name())
		if randomApiKey == "" {
			s.skipped = true
			return
		}

		getUrl := fmt.Sprintf("https://osint.bevigil.com/api/%s/urls/", domain)

		resp, err := session.Get(ctx, getUrl, "", map[string]string{
			"X-Access-Token": randomApiKey, "User-Agent": "urlfounder",
		})
		isForbidden := resp != nil && resp.StatusCode == http.StatusForbidden
		if err != nil {
			println(err.Error())
			if !isForbidden {
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
				session.DiscardHTTPResponse(resp)
				return
			}
		}

		var urls []string
		var response Response
		err = jsoniter.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			resp.Body.Close()
			return
		}

		resp.Body.Close()

		if len(response.Urls) > 0 {
			urls = response.Urls
		}

		for _, url := range urls {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: url}
		}

	}()
	return results
}

func (s *Source) Name() string {
	return "bevigil"
}

func (s *Source) IsDefault() bool {
	return false
}

func (s *Source) HasRecursiveSupport() bool {
	return false
}

func (s *Source) NeedsKey() bool {
	return true
}

func (s *Source) AddApiKeys(keys []string) {
	s.apiKeys = keys
}

func (s *Source) Statistics() subscraping.Statistics {
	return subscraping.Statistics{
		Errors:    s.errors,
		Results:   s.results,
		TimeTaken: s.timeTaken,
		Skipped:   s.skipped,
	}
}
