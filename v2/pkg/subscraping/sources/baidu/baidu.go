package baidu

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chainreactors/urlfounder/v2/pkg/subscraping"
	"strings"
	"time"
)

// Source is the baidu scraping agent
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

		baseURL := "https://www.baidu.com/s?wd=site:%s&pn=%d"
		page := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				url := fmt.Sprintf(baseURL, domain, page)
				resp, err := session.Get(ctx, url, "", nil)
				if err != nil {
					results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
					s.errors++
					return
				}
				defer resp.Body.Close()

				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
					s.errors++
					return
				}

				doc.Find("a").Each(func(_ int, sel *goquery.Selection) {
					href, exists := sel.Attr("href")
					if exists {
						if strings.HasPrefix(href, "http") {
							results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: href}
							s.results++
						}
					}
				})

				page += 10
			}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "baidu"
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
	return false
}

func (s *Source) Statistics() subscraping.Statistics {
	return subscraping.Statistics{
		Errors:    s.errors,
		Results:   s.results,
		TimeTaken: s.timeTaken,
		Skipped:   s.skipped,
	}
}
