package runner

import (
	"net"
	"strings"

	"github.com/chainreactors/urlfounder/v2/pkg/passive"
	"github.com/chainreactors/urlfounder/v2/pkg/resolve"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
)

// initializePassiveEngine creates the passive engine and loads sources etc
func (r *Runner) initializePassiveEngine() {
	r.passiveAgent = passive.New(r.options.Sources, r.options.ExcludeSources, r.options.All, r.options.OnlyRecursive)
}

// initializeResolver creates the resolver used to resolve the found urls
func (r *Runner) initializeResolver() error {
	var resolvers []string

	if len(r.options.Resolvers) > 0 {
		resolvers = append(resolvers, r.options.Resolvers...)
	} else {
		resolvers = append(resolvers, resolve.DefaultResolvers...)
	}

	// Add default 53 UDP port if missing
	for i, resolver := range resolvers {
		if !strings.Contains(resolver, ":") {
			resolvers[i] = net.JoinHostPort(resolver, "53")
		}
	}

	r.resolverClient = resolve.New()
	var err error
	r.resolverClient.DNSClient, err = dnsx.New(dnsx.Options{BaseResolvers: resolvers, MaxRetries: 5})
	if err != nil {
		return nil
	}

	return nil
}
