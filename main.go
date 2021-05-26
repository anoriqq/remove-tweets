package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	flagMaxID := flag.String("maxid", "", "maxID")
	flagScreenName := flag.String("screenname", "", "screenName")
	flagUntil := flag.String("until", "", "until")
	flag.Parse()
	if len(*flagMaxID) < 1 {
		return
	}
	if len(*flagScreenName) < 1 {
		return
	}
	if len(*flagUntil) < 1 {
		return
	}

	s := NewTwitterService()

	u := s.GetUser(*flagScreenName)

	i, err := strconv.Atoi(*flagMaxID)
	if err != nil {
		panic(err)
	}
	maxID := int64(i)

	until, err := time.Parse(time.RFC3339, *flagUntil)
	if err != nil {
		panic(err)
	}

	for {
		ts := s.GetTimeline(u.ID, maxID)
		if len(ts) < 1 {
			break
		}

		for _, t := range ts {
			if !IsCreatedBeforeThresholdDateTime(t, until) {
				continue
			}

			if t.Retweeted {
				s.Unretweet(t.ID)
			} else {
				s.Delete(t.ID)
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

type twitterService struct {
	c *twitter.Client
}

func (s twitterService) GetUser(screenName string) *twitter.User {
	u, _, err := s.c.Users.Show(&twitter.UserShowParams{ScreenName: screenName})
	if err != nil {
		panic(err)
	}
	return u
}

func (s twitterService) GetTimeline(userID, maxID int64) []twitter.Tweet {
	includeRetweets := true
	excludeReplies := false

	ts, _, err := s.c.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          userID,
		MaxID:           maxID,
		Count:           200,
		IncludeRetweets: &includeRetweets,
		ExcludeReplies:  &excludeReplies,
	})
	if err != nil {
		panic(err)
	}

	if ts[0].ID == maxID {
		ts = ts[1:]
	}

	return ts
}

func (s twitterService) Delete(id int64) {
	t, _, err := s.c.Statuses.Destroy(id, &twitter.StatusDestroyParams{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted: %d\n%s\n========\n", t.ID, t.Text)
}

func (s twitterService) Unretweet(id int64) {
	t, _, err := s.c.Statuses.Unretweet(id, &twitter.StatusUnretweetParams{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Unretweeted: %d\n%s\n========\n", t.ID, t.Text)
}

func NewTwitterService() *twitterService {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerKeySecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerKeySecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	c := twitter.NewClient(httpClient)

	return &twitterService{
		c: c,
	}
}
