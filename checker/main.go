// checker is a simpler command line tool for checking entries to a raffle
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	envTwitterConsumerKey    = "TWITTER_CONSUMER_KEY"
	envTwitterConsumerSecret = "TWITTER_CONSUMER_SECRET"
	envTwitterAccessToken    = "TWITTER_ACCESS_TOKEN"
	envTwitterAccessSecret   = "TWITTER_ACCESS_SECRET"

	maxTwitterMentions = 200

	timeFormat = "_2 Jan 2006 15:04"
)

var (
	flagSet         = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fKey            = flagSet.String("key", "", "the key announced to participants")
	fTwitter        = flagSet.Bool("twitter", false, "find entries from Twitter mention timeline")
	fConsumerKey    = flagSet.String("consumerkey", "", "Twitter Consumer Key - (default env "+envTwitterConsumerKey+")")
	fConsumerSecret = flagSet.String("consumersecret", "", "Twitter Consumer Secret - (default env "+envTwitterConsumerSecret+")")
	fAccessToken    = flagSet.String("accesstoken", "", "Twitter Access Token - (default env "+envTwitterAccessToken+")")
	fAccessSecret   = flagSet.String("accesssecret", "", "Twitter Access Secret - (default env "+envTwitterAccessSecret+")")
	fStart          = flagSet.String("start", "", "the start time of the Twitter raffle. Defaults to the start of time :)")
	fEnd            = flagSet.String("end", "", "the end time of the Twitter raffle - (default time.Now())")
	fLocation       = flagSet.String("location", "", "the location to use for parsing -start and -end - defaults to system's local time zone")
	fEntryRegex     = flagSet.String("entryregex", "please enter me into the @jetbrains raffle", "the regex used to identify Twitter-based entries to the raffle")
	fWinners        = flagSet.Int("winners", 3, "the max number of winners to draw")

	hashRegex = regexp.MustCompile("[a-f0-9]{64}")
)

// Takes either handle:entry pairs args command line args, or if there are none
// read lines from stdin

func main() {
	os.Exit(main1())
}

func main1() int {
	switch err := mainerr(); err {
	case nil:
		return 0
	case flag.ErrHelp:
		return 2
	default:
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
}

func mainerr() (err error) {
	defer func() {
		switch r := recover().(type) {
		case nil:
		case error:
			err = r
		default:
			panic(r)
		}
	}()

	*fConsumerKey = os.Getenv(envTwitterConsumerKey)
	*fConsumerSecret = os.Getenv(envTwitterConsumerSecret)
	*fAccessToken = os.Getenv(envTwitterAccessToken)
	*fAccessSecret = os.Getenv(envTwitterAccessSecret)

	flagSet.Usage = func() {
		mainUsage(os.Stderr)
	}
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	if *fKey == "" {
		errorf("must provide -key")
	}
	if *fWinners < 1 {
		errorf("invalid max number of winners %v", *fWinners)
	}

	var entries []entry
	if *fTwitter {
		entries = twitterEntries()
	} else {
		entries = nonTwitterEntries()
	}

	var winners []entry
	for _, e := range entries {
		if check(*fKey, e) {
			winners = append(winners, e)
		}
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(winners), func(i, j int) {
		winners[i], winners[j] = winners[j], winners[i]
	})
	fmt.Printf("Number of entries: %v\n", len(entries))
	fmt.Printf("Number of valid entries: %v\n", len(winners))
	fmt.Println("")
	if len(winners) == 0 {
		fmt.Println("No winners")
	} else {
		for i := 0; i < *fWinners && i < len(winners); i++ {
			fmt.Printf("%v is a winner!\n", winners[i].handle)
		}
	}

	return nil
}

func twitterEntries() (res []entry) {
	if *fConsumerKey == "" || *fConsumerSecret == "" || *fAccessToken == "" || *fAccessSecret == "" {
		errorf("Consumer key/secret and Access token/secret required")
	}
	config := oauth1.NewConfig(*fConsumerKey, *fConsumerSecret)
	token := oauth1.NewToken(*fAccessToken, *fAccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	mentionTimelineParams := &twitter.MentionTimelineParams{
		Count:     maxTwitterMentions,
		TweetMode: "extended",
	}
	tweets, _, err := client.Timelines.MentionTimeline(mentionTimelineParams)
	if err != nil {
		errorf("failed to retrieve Twitter mentions: %v", err)
	}
	var start time.Time
	end := time.Now()
	loc := time.Local
	if *fLocation != "" {
		ploc, err := time.LoadLocation(*fLocation)
		if err != nil {
			errorf("failed to load time location %v: %v", *fLocation, err)
		}
		loc = ploc
	}
	if *fStart != "" {
		stime, err := time.ParseInLocation(timeFormat, *fStart, loc)
		if err != nil {
			errorf("failed to parse start time %q: %v", *fStart, err)
		}
		start = stime
	}
	if *fEnd != "" {
		etime, err := time.ParseInLocation(timeFormat, *fEnd, loc)
		if err != nil {
			errorf("failed to parse end time %q: %v", *fEnd, err)
		}
		end = etime
	}
	match, err := regexp.Compile(*fEntryRegex)
	if err != nil {
		errorf("failed to compile regex %q for matching entries: %v", *fEntryRegex, err)
	}
	var minTime time.Time
	for _, t := range tweets {
		created, err := t.CreatedAtTime()
		if err != nil {
			errorf("failed to get created time for tweet %v: %v", t.ID, err)
		}
		if minTime.IsZero() || created.Before(minTime) {
			minTime = created
		}
		if created.Before(start) || end.Before(created) {
			continue
		}
		if !match.MatchString(t.FullText) {
			continue
		}
		hashes := hashRegex.FindStringSubmatch(t.FullText)
		if len(hashes) != 1 {
			errorf("failed to find single hash in tweet %v", t.ID)
		}
		res = append(res, entry{
			handle: t.User.ScreenName,
			hash:   hashes[0],
		})
	}
	if start.Before(minTime) {
		fmt.Fprintf(os.Stderr, "******** WARNING: start time of raffle before first entry in mentiones timeline (limited to %v)", maxTwitterMentions)
	}
	return
}

func nonTwitterEntries() (res []entry) {
	var pairs []string
	if len(flagSet.Args()) > 0 {
		pairs = flagSet.Args()
	} else {
		fmt.Fprintf(os.Stderr, "reading entries from stdin...\n")
		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			errorf("failed to read stdin: %v\n", err)
		}
		pairs = strings.Split(string(in), "\n")
	}
	for _, p := range pairs {
		parts := strings.Split(strings.TrimSpace(p), ":")
		if len(parts) != 2 {
			errorf("pair %q invalid format; should be handle:hash", p)
		}
		res = append(res, entry{
			handle: parts[0],
			hash:   parts[1],
		})
	}
	return
}

func errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}

type entry struct {
	handle string
	hash   string
}

func check(secret string, e entry) bool {
	hash := sha256.New()
	fmt.Fprintf(hash, "Handle: %v\n", e.handle)
	fmt.Fprintf(hash, "Key: %v\n", *fKey)

	valid := fmt.Sprintf("%x", hash.Sum(nil))
	return valid == e.hash
}
