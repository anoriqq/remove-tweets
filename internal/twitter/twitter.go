package twitter

import (
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Tweet = twitter.Tweet

type twitterService struct {
	c *twitter.Client
}

func (s twitterService) GetUser(screenName string) (*twitter.User, error) {
	u, _, err := s.c.Users.Show(&twitter.UserShowParams{ScreenName: screenName})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s twitterService) GetTimeline(userID, maxID int64) ([]twitter.Tweet, error) {
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
		return nil, err
	}

	if ts[0].ID == maxID {
		ts = ts[1:]
	}

	return ts, nil
}

func (s twitterService) Delete(id int64) error {
	t, _, err := s.c.Statuses.Destroy(id, &twitter.StatusDestroyParams{})
	if err != nil {
		return err
	}

	fmt.Printf("Deleted: %d\n%s\n========\n", t.ID, t.Text)

	return nil
}

func (s twitterService) Unretweet(id int64) error {
	t, _, err := s.c.Statuses.Unretweet(id, &twitter.StatusUnretweetParams{})
	if err != nil {
		return err
	}

	fmt.Printf("Unretweeted: %d\n%s\n========\n", t.ID, t.Text)

	return nil
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
