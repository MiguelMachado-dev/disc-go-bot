package scraper

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly"
)

func Caliber(playersCh chan string) {
	// Prints time that the scraper started
	fmt.Println("Start scraping at", time.Now())

	c := colly.NewCollector(
		colly.AllowedDomains("steamcommunity.com"),
	)

	c.OnHTML(".apphub_HeaderBottom div.apphub_Stats", func(e *colly.HTMLElement) {
		stats := Stats{}
		stats.Players = e.ChildText(".apphub_NumInApp")

		// Use a regular expression to match numbers and commas
		re := regexp.MustCompile(`[\d,]+`)
		matches := re.FindString(stats.Players)

		// Send the players count through the channel
		playersCh <- matches
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Status code:", r.StatusCode)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://steamcommunity.com/app/307950")
}
