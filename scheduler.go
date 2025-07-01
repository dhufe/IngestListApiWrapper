package main

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"time"
)

func NewScheduler() *gocron.Scheduler {
	TimeLocation, err := time.LoadLocation("Europe/Berlin")
	s, err := gocron.NewScheduler(
		gocron.WithLocation(TimeLocation),
	)

	if err != nil {
		fmt.Println(err)
	}

	return &s
}
