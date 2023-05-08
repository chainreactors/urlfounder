package alienvault

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chainreactors/urlfounder/v2/pkg/subscraping"
)

type alienvaultResponse struct {
	Detail  string `json:"detail"`
	Error   string `json:"error"`
	UrlList []struct {
		Url string `json:"url"`
	} `json:"url_list"`
}

// Source is the passive scraping agent
type Source struct {
	timeTaken time.Duration
	errors    int
	results   int
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

		api := fmt.Sprintf("https://otx.alienvault.com/otxapi/indicators/domain/url_list/%s?limit=1000&page=1", domain)
		//resp, err := session.Get(ctx, api, "", headers)
		resp, err := session.SimpleGet(ctx, api)
		if err != nil && resp == nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			s.errors++
			session.DiscardHTTPResponse(resp)
			return
		}

		var response alienvaultResponse
		// Get the response body and decode
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
			s.errors++
			resp.Body.Close()
			return
		}
		resp.Body.Close()

		if response.Error != "" {
			results <- subscraping.Result{
				Source: s.Name(), Type: subscraping.Error, Error: fmt.Errorf("%s, %s", response.Detail, response.Error),
			}
			return
		}

		for _, record := range response.UrlList {
			results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: record.Url}
			s.results++
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "alienvault"
}

func (s *Source) IsDefault() bool {
	return true
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
	}
}
