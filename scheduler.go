package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func NewScheduler() *gocron.Scheduler {
	// TimeLocation requires tzdata (or similar on your linux disto)
	TimeLocation, err := time.LoadLocation("Europe/Berlin")
	s, err := gocron.NewScheduler(
		gocron.WithLocation(TimeLocation),
	)
	if err != nil {
		fmt.Println(err)
	}

	return &s
}
