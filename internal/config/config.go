package config

import (
	"errors"
	"strconv"
	"time"
)

type Config struct {
	ScreenName string
	MaxID      int64
	Until      time.Time
}

func NewConfig(screenName, maxIDStr, untilStr string) (Config, error) {
	if len(screenName) < 1 {
		return Config{}, errors.New("screenName is required")
	}
	if len(maxIDStr) < 1 {
		return Config{}, errors.New("maxID is required")
	}
	if len(untilStr) < 1 {
		return Config{}, errors.New("until is required")
	}

	maxID, err := strconv.Atoi(maxIDStr)
	if err != nil {
		return Config{}, err
	}

	until, err := time.Parse(time.RFC3339, untilStr)
	if err != nil {
		return Config{}, err
	}

	return Config{
		ScreenName: screenName,
		MaxID:      int64(maxID),
		Until:      until,
	}, nil
}
