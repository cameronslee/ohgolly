package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gocolly/colly"
)

const DDG string = "https://html.duckduckgo.com/html/?q="

func scrape(search_query string) []Result {
	var res []Result
	c := colly.NewCollector(
		colly.AllowedDomains("html.duckduckgo.com"),
		colly.Async(true),
	)

	// callback
	c.OnHTML("h2.result__title a.result__a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		temp := Result{Title: e.Text, Link: link}
		res = append(res, temp)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// scrape
	c.Visit(DDG + search_query)

	c.Wait()

	return res
}

func NewResultsHandler(search_query string) ResultsHandler {
	resultsGetter := func() (result []Result, err error) {

		return scrape(search_query), nil
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
	http.Handle("/", templ.Handler(home()))
	http.Handle("/results", NewResultsHandler("golang"))

	fmt.Println("listening on http://localhost:8090")
	if err := http.ListenAndServe("localhost:8090", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}
