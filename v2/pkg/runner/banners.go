package runner

import (
	"github.com/projectdiscovery/gologger"
)

const banner = `
_____  _______________________                   _________            
__/ / / /__/ __ \__/ /___/ __/_________  ______________  /____________
_/ / / /__/ /_/ /_/ / __/ /_ _/ __ \  / / /_  __ \  __  /_  _ \_  ___/
/ /_/ / _  _, _/_/ /___/ __/ / /_/ / /_/ /_  / / / /_/ / /  __/  /    
\____/  /_/ |_| /_____/_/    \____/\__,_/ /_/ /_/\__,_/  \___//_/   V0.0.1
`

// Name
const ToolName = `urlfounder`

// Version is the current version of urlfounder
const version = `v0.0.1`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
}
