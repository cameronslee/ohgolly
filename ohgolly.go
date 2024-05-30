package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

const DDG string = "https://html.duckduckgo.com/html/?q=golang"

func main() {
	//search_query := ""

	c := colly.NewCollector(
		colly.AllowedDomains("html.duckduckgo.com"),
		colly.Async(true),
	)

	fmt.Print("FOO\n")
	// On every a element which has href attribute call callback
	c.OnHTML("h2.result__title a.result__a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	c.Visit(DDG)

	c.Wait()
}
