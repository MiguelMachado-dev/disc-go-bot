package scraper

import (
	"regexp"
	"time"

	"github.com/MiguelMachado-dev/disc-go-bot/config"
	"github.com/gocolly/colly"
)

type Stats struct {
	Players string
}

var log = config.NewLogger("huntScraper")

func hunt(playersCh chan string) {
	// Prints time that the scraper started
	log.Infoln("Start scraping at", time.Now())

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
		log.Infoln("Status code:", r.StatusCode)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Infoln("Visiting:", r.URL)
	})

	c.Visit("https://steamcommunity.com/app/594650")
}
