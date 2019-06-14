## `checker`

The `checker` app requires Twitter API credentials. See [the
wiki](https://github.com/go-london-user-group/raffle/wiki/checker-Twitter-credentials).

### Usage

```
The checker command checks entries to a raffle.

checker defines the following flags:

  -accesssecret string
    	Twitter Access Secret - (default env TWITTER_ACCESS_SECRET)
  -accesstoken string
    	Twitter Access Token - (default env TWITTER_ACCESS_TOKEN)
  -consumerkey string
    	Twitter Consumer Key - (default env TWITTER_CONSUMER_KEY)
  -consumersecret string
    	Twitter Consumer Secret - (default env TWITTER_CONSUMER_SECRET)
  -end string
    	the end time of the Twitter raffle - (default time.Now())
  -entryregex string
    	the regex used to identify Twitter-based entries to the raffle (default "please enter me into the @jetbrains raffle")
  -key string
    	the key announced to participants
  -location string
    	the location to use for parsing -start and -end - defaults to system's local time zone
  -start string
    	the start time of the Twitter raffle. Defaults to the start of time :)
  -twitter
    	find entries from Twitter mention timeline
  -winners int
    	the max number of winners to draw (default 3)

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
```
