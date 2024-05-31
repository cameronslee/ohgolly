package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gocolly/colly"
)

const DDG string = "https://html.duckduckgo.com/html/?q="

func scrape(search_query string) {
	c := colly.NewCollector(
		colly.AllowedDomains("html.duckduckgo.com"),
		colly.Async(true),
	)

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
	c.Visit(DDG + search_query)

	c.Wait()
}

func NewResultsHandler() ResultsHandler {
	// Replace this in-memory function with a call to a database.
	resultsGetter := func() (result []Result, err error) {

		return []Result{{Title: "templ", Link: "author"}}, nil
	}
	return ResultsHandler{
		GetResults: resultsGetter,
		Log:        log.Default(),
	}
}

type ResultsHandler struct {
	Log        *log.Logger
	GetResults func() ([]Result, error)
}

func (rh ResultsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps, err := rh.GetResults()
	if err != nil {
		rh.Log.Printf("failed to get results: %v", err)
		http.Error(w, "failed to retrieve results", http.StatusInternalServerError)
		return
	}
	templ.Handler(results(ps)).ServeHTTP(w, r)
}

type Result struct {
	Title string
	Link  string
}

func main() {
	// Use a template that doesn't take parameters.
	http.Handle("/", templ.Handler(home()))

	// Use a template that accesses data or handles form posts.
	http.Handle("/results", NewResultsHandler())

	// Start the server.
	fmt.Println("listening on http://localhost:8090")
	if err := http.ListenAndServe("localhost:8090", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}
