package scraper

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type Stats struct {
	Players string
}

func GuiltyGear(playersCh chan string) {
	// Prints time that the scraper started
	fmt.Println("Start scraping at", time.Now())

	c := colly.NewCollector(
		colly.AllowedDomains("steamcharts.com"),
	)

	c.OnHTML("#app-heading div.app-stat:first-of-type", func(e *colly.HTMLElement) {
		stats := Stats{}
		stats.Players = e.ChildText(".num")

		// Send the players count through the channel
		playersCh <- stats.Players
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status code:", r.StatusCode)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://steamcharts.com/app/1384160")
}
