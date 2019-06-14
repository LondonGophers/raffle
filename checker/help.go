package main

import (
	"bytes"
	"fmt"
	"io"
)

func mainUsage(f io.Writer) {
	var flagUsage bytes.Buffer
	flagSet.SetOutput(&flagUsage)
	flagSet.PrintDefaults()
	fmt.Fprintf(f, mainHelp, flagUsage.Bytes())
}

var mainHelp = `
The checker command checks entries to a raffle.

checker defines the following flags:

%s
The -key flag must be provided; this is the secret announced to participants.

The -twitter flag instruct checker to take entries from a user's mentions
timeline. This requires the associated consumer key/secret and access
token/secret flags (or their associated environment variables) to have been
set. For more details on Twitter credentials, see the wiki.

The -start and -end flags can be used to control the start/end of of the
Twitter raffle. Times should be specified according to the reference format:

  _2 Jan 2006 15:04

Times will be parsed according to the system's local time zone, unless
overridden with -location.

Entries to the raffle will be identified from the mentions timeline according
to the -entrymatch flag which is parse as a regular expression. Within an
entry, the hash is further identified as hex string of length 64.

If the -twitter flag is not provided, handle:hash entry pairs will be read from
args command line args, or if there are none read lines from stdin.

The -winners flag determines the maximum number of valid entries to draw as
winners, after a pseudo-random shuffle of valid entries.

Examples:

  checker -twitter -key gophers-with-flare -start "15 May 2019 19:00" -end "15 may 2019 21:00"
`[1:]
