package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/anoriqq/remove-tweets/internal/twitter"
)

var (
	maxID      string
	screenName string
	until      string
)

func init() {
	flag.StringVar(&maxID, "maxid", "", "maxID")
	flag.StringVar(&screenName, "screenname", "", "screenName")
	flag.StringVar(&until, "until", "", "until")
}

func main() {
	flag.Parse()
	if len(maxID) < 1 {
		return
	}
	if len(screenName) < 1 {
		return
	}
	if len(until) < 1 {
		return
	}

	s := twitter.NewTwitterService()

	u, err := s.GetUser(screenName)
	if err != nil {
		panic(err)
	}

	i, err := strconv.Atoi(maxID)
	if err != nil {
		panic(err)
	}
	maxID := int64(i)

	until, err := time.Parse(time.RFC3339, until)
	if err != nil {
		panic(err)
	}

	for {
		ts, err := s.GetTimeline(u.ID, maxID)
		if err != nil {
			panic(err)
		}
		if len(ts) < 1 {
			break
		}

		for _, t := range ts {
			if !IsCreatedBeforeThresholdDateTime(t, until) {
				continue
			}

			if t.Retweeted {
				err := s.Unretweet(t.ID)
				if err != nil {
					panic(err)
				}
			} else {
				err := s.Delete(t.ID)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

// IsCreatedBeforeThresholdDateTime tがuntilよりも以前に作成されていたらtrue
func IsCreatedBeforeThresholdDateTime(t twitter.Tweet, until time.Time) bool {
	createdAt, err := t.CreatedAtTime()
	if err != nil {
		panic(err)
	}

	return createdAt.Before(until)
}
