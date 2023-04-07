package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Stats struct {
	Players string
}

func GuiltyGear(playersCh chan string) {
	fmt.Println("Start scraping")

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