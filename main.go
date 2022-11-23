package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/mattn/go-mastodon"
)

// stats returns current and total number of days in year
func stats() (float64, float64) {
	now := time.Now()
	yearStart := time.Date(now.Year(), time.January, 0, 0, 0, 0, 0, time.Local)
	yearEnd := time.Date(now.Year()+1, time.January, 0, 0, 0, 0, 0, time.Local)
	current := now.Sub(yearStart).Hours() / 24
	total := yearEnd.Sub(now).Hours()/24 + current
	return current, total
}

func createBar(progress float64) string {
	var bars = [][]string{{"⬜", "⬛"}, {"░", "█"}, {"▱", "▰"}, {"▯", "▮"}}
	choice := bars[0]

	pb := ""
	filled := int(math.Floor(progress / 10))
	for i := 0; i < filled; i++ {
		pb += choice[1]
	}
	for i := filled; i < 10; i++ {
		pb += choice[0]
	}

	pb += " " + strconv.Itoa(int(progress)) + "%"
	return pb
}

func main() {
	current, total := stats()
	perc := (current / total) * 100
	pb := createBar(perc)
	fmt.Println(pb)

	MASTODON_SERVER := os.Getenv("MASTODON_SERVER")
	MASTODON_CLIENT_ID := os.Getenv("MASTODON_CLIENT_ID")
	MASTODON_CLIENT_SECRET := os.Getenv("MASTODON_CLIENT_SECRET")
	MASTODON_ACCESS_TOKEN := os.Getenv("MASTODON_ACCESS_TOKEN")
	MASTODON_USERNAME := os.Getenv("MASTODON_USERNAME")
	MASTODON_PASSWORD := os.Getenv("MASTODON_PASSWORD")

	c := mastodon.NewClient(&mastodon.Config{
		Server:       MASTODON_SERVER,
		ClientID:     MASTODON_CLIENT_ID,
		ClientSecret: MASTODON_CLIENT_SECRET,
		AccessToken:  MASTODON_ACCESS_TOKEN,
	})

	err := c.Authenticate(context.Background(), MASTODON_USERNAME, MASTODON_PASSWORD)
	if err != nil {
		log.Fatal(err)
	}

	c.PostStatus(context.Background(), &mastodon.Toot{
		Status:     pb,
		Visibility: mastodon.VisibilityPublic,
	})
}
