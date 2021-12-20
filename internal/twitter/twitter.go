package twitter

import (
	"os"

	"github.com/anoriqq/remove-tweets/internal/logger"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Tweet = twitter.Tweet

type twitterService struct {
	c      *twitter.Client
	logger logger.Logger
}

func (s twitterService) GetUser(screenName string) (*twitter.User, error) {
	u, _, err := s.c.Users.Show(&twitter.UserShowParams{ScreenName: screenName})
	if err != nil {
		return nil, err
	}

	s.logger.Infof("get user: %v", screenName)

	return u, nil
}

// GetTimeline userIDのmaxID以前のtweetsを取得する
func (s twitterService) GetTimeline(userID, maxID int64) ([]twitter.Tweet, error) {
	includeRetweets := true
	excludeReplies := false

	tt, _, err := s.c.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          userID,
		MaxID:           maxID,
		Count:           200,
		IncludeRetweets: &includeRetweets,
		ExcludeReplies:  &excludeReplies,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Infof("get timeline: %v", len(tt))

	if tt[0].ID == maxID {
		tt = tt[1:]
	}

	return tt, nil
}

// Delete IDのtweetを削除する
func (s twitterService) Delete(id int64) error {
	t, _, err := s.c.Statuses.Destroy(id, &twitter.StatusDestroyParams{})
	if err != nil {
		return err
	}

	s.logger.Infof("deleted: %v: %v", t.CreatedAtTime, t.Text)

	return nil
}

// Unretweet idのretweetを取り消す
func (s twitterService) Unretweet(id int64) error {
	t, _, err := s.c.Statuses.Unretweet(id, &twitter.StatusUnretweetParams{})
	if err != nil {
		return err
	}

	s.logger.Infof("unretweeted: %v: %v", t.CreatedAtTime, t.Text)

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
	logger := logger.NewLogger()

	return &twitterService{
		c:      c,
		logger: logger,
	}
}
