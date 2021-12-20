package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/anoriqq/remove-tweets/internal/logger"
	"github.com/anoriqq/remove-tweets/internal/twitter"
)

var (
	maxID      string
	screenName string
	until      string
)

type config struct {
	maxID      string
	screenName string
	until      string
}

func (c config) Valid() error {
	if len(maxID) < 1 {
		return errors.New("maxID is required")
	}
	if len(screenName) < 1 {
		return errors.New("screenName is required")
	}
	if len(until) < 1 {
		return errors.New("until is required")
	}

	return nil
}

func init() {
	flag.StringVar(&maxID, "maxid", "", "maxID")
	flag.StringVar(&screenName, "screenname", "", "screenName")
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

	c := config{
		maxID:      maxID,
		screenName: screenName,
		until:      until,
	}
	err := c.Valid()
	if err != nil {
		return err
	}

	s := twitter.NewTwitterService()

	u, err := s.GetUser(screenName)
	if err != nil {
		return err
	}

	i, err := strconv.Atoi(maxID)
	if err != nil {
		return err
	}
	maxID := int64(i)

	until, err := time.Parse(time.RFC3339, until)
	if err != nil {
		return err
	}

	for {
		ts, err := s.GetTimeline(u.ID, maxID)
		if err != nil {
			return err
		}
		if len(ts) < 1 {
			break
		}

		for _, t := range ts {
			// tがuntilよりも後に作成されていたらskip
			createdAt, err := t.CreatedAtTime()
			if err != nil {
				return err
			}
			if createdAt.After(until) {
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
