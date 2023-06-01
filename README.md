
fork from https://github.com/projectdiscovery/subfinder

# intro

<h4 align="center">Fast passive url enumeration tool.</h4>

![](https://socialify.git.ci/chainreactors/urlfounder/image?description=1&font=Inter&forks=1&issues=1&language=1&name=1&owner=1&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Light)

<p align="center">
  <a href="#features">Features</a> •
  <a href="#running-urlfounder">Usage</a> •
  <a href="#post-installation-instructions">API Setup</a> •
  <a href="#urlfounder-go-library">Library</a> •
</p>

`urlfounder` is a url discovery tool that returns valid urls for domain, using passive online sources. 

# Features

<h1 align="left">
  <img src="static/urlfounder-run.png" alt="urlfounder" width="700px" >
</h1>

- Probe and extract key information, such as title, status code
- Support for multiple output formats (JSON, file, standard output)
- Multiple data sources to maximise coverage and accuracy of search results
- Fast and lightweight
- Bulk input support

# Usage

```
urlfounder -h
```

This will display help for the tool. Here are all the switches it supports.

```
Usage:
./urlfounder [flags]

Flags:
INPUT:
    -d, -domain string[]  domains to find urls for
    -dL, -list string     file containing list of domains for url discovery

    SOURCE:
    -s, -sources string[]           specific sources to use for discovery. Use -ls to display all available sources.
    -all                            use all sources for enumeration (slow)
    -es, -exclude-sources string[]  sources to exclude from enumeration (-es alienvault,zoomeye)

FILTER:
    -m, -match string[]   url or list of url to match (file or comma separated)
    -f, -filter string[]   url or list of url to filter (file or comma separated)

RATE-LIMIT:
    -rl, -rate-limit int  maximum number of http requests to send per second
    -t int                number of concurrent goroutines for resolving (-active only) (default 10)

OUTPUT:
    -o, -output string       file to write output to
    -oJ, -json               write output in JSONL(ines) format
    -oD, -output-dir string  directory to write output (-dL only)
    -cs, -collect-sources    include all sources in the output (-json only)
    -sc, -status             include StatusCode in output
    -tI, -title              include url titles in output

CONFIGURATION:
    -config string                flag config file (default "/root/.config/urlfounder/config.yaml")
    -pc, -provider-config string  provider config file (default "/root/.config/urlfounder/provider-config.yaml")
    -r string[]                   comma separated list of resolvers to use
    -aC, -active                  display active urls only
    -proxy string                 http proxy to use with urlfounder

DEBUG:
    -silent             show only urls in output
    -version            show version of urlfounder
    -v                  show verbose output
    -nc, -no-color      disable color in output
    -ls, -list-sources  list all available sources
    -stats              report source statistics

OPTIMIZATION:
    -timeout int   seconds to wait before timing out (default 30)
    -max-time int  minutes to wait for enumeration results (default 10)
```

## Post Installation Instructions

You can use the `urlfounder -ls` command to display all the available sources.

The following services require set http proxy to use：

- webarchive
- alienvault

The following services require configuring API keys to work:

- [BeVigil](https://bevigil.com/osint-api)

These values are stored in the `$HOME/.config/urlfounder/provider-config.yaml` file which will be created when you run the tool for the first time. 

An example provider config file:

```
alienvault: []
bevigil:
  - Tu8DSd6GqM1jDDDD
webarchive: []
```

# Running Urlfounder

To run the tool on a target, just use the following command.

```console
./urlfounder -d projectdiscovery.io

_____  _______________________                   _________
__/ / / /__/ __ \__/ /___/ __/_________  ______________  /____________
_/ / / /__/ /_/ /_/ / __/ /_ _/ __ \  / / /_  __ \  __  /_  _ \_  ___/
/ /_/ / _  _, _/_/ /___/ __/ / /_/ / /_/ /_  / / / /_/ / /  __/  /
\____/  /_/ |_| /_____/_/    \____/\__,_/ /_/ /_/\__,_/  \___//_/   V0.0.1

[INF] Loading provider config from /root/.config/urlfounder/provider-config.yaml
[INF] Enumerating urls for projectdiscovery.io
https://projectdiscovery.io/
https://chaos.projectdiscovery.io/
https://dns.projectdiscovery.io/dns/%s/subdomainsinvalid
https://dns.projectdiscovery.io/dns/%s/public
https://dns.projectdiscovery.io/dns/projectdiscovery.io
https://projectdiscovery.io
https://dns.projectdiscovery.io/dns/%s/subdomains
https://blog.projectdiscovery.io/abusing
https://dns.projectdiscovery.io/dns/
https://dns.projectdiscovery.io/dns/addinvalid
http://projectdiscovery.io/
https://projectdiscovery.io/open
https://dns.projectdiscovery.io/dns/%s/subdomainsinternal
https://nuclei.projectdiscovery.io
https://nuclei.projectdiscovery.io/faq/templates/
https://nuclei.projectdiscovery.io/templating
https://cdn.projectdiscovery.io/cdn/cdn
https://nuclei.projectdiscovery.io/faq/nuclei/
[INF] Found 18 urls for projectdiscovery.io in 564 milliseconds 619 microseconds
```

## Urlfounder Go library

Usage example：

```
package main

import (
	"bytes"
	"fmt"
	"github.com/chainreactors/urlfounder/v2/pkg/resolve"
	"github.com/chainreactors/urlfounder/v2/pkg/runner"
	"io"
	"log"
)

func main() {
	runnerInstance, err := runner.NewRunner(&runner.Options{
		Threads:            10,                       // Thread controls the number of threads to use for active enumerations
		Timeout:            30,                       // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10,                       // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		Resolvers:          resolve.DefaultResolvers, // Use the default list of resolvers by marshaling it to the config
		ResultCallback: func(s *resolve.HostEntry) { // Callback function to execute for available host
			log.Println(s.Host, s.Source)
		},
	})

	buf := bytes.Buffer{}
	err = runnerInstance.EnumerateSingleDomain("projectdiscovery.io", []io.Writer{&buf})
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(&buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", data)
}
```

### Wiki

todo

# THANKS

- [urlfounder](https://github.com/chainreactors/urlfounder)
