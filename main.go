package main

import (
	"flag"
	"os"

	"github.com/anoriqq/remove-tweets/internal/config"
	"github.com/anoriqq/remove-tweets/internal/logger"
	"github.com/anoriqq/remove-tweets/internal/twitter"
)

var (
	screenName string
	maxID      string
	until      string
)

func init() {
	flag.StringVar(&screenName, "screenname", "", "screenName")
	flag.StringVar(&maxID, "maxid", "", "maxID")
	flag.StringVar(&until, "until", "", "until")
}

func main() {
	logger := logger.NewLogger()
	logger.Info("start")

	err := run()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("completed")
	os.Exit(0)
}

func run() error {
	flag.Parse()

	c, err := config.NewConfig(screenName, maxID, until)
	if err != nil {
		return err
	}

	s := twitter.NewTwitterService()

	u, err := s.GetUser(c.ScreenName)
	if err != nil {
		return err
	}

	// 削除可能なtweetがなくなるまでloop
	for {
		tt, err := s.GetTimeline(u.ID, c.MaxID)
		if err != nil {
			return err
		}
		if len(tt) < 1 {
			break
		}

		for _, t := range tt {
			// tがuntilよりも後に作成されていたらskip
			// TODO: maxID指定してるのでこの判定処理いらないのでは
			createdAt, err := t.CreatedAtTime()
			if err != nil {
				return err
			}
			if createdAt.After(c.Until) {
				continue
			}

			if t.Retweeted {
				err := s.Unretweet(t.ID)
				if err != nil {
					return err
				}
			} else {
				err := s.Delete(t.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
