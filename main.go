package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const duckDuckGoURL = "https://duckduckgo.com/html/"

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

type SearchResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

func randomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(len(userAgents))
	return userAgents[randNum]
}

func buildDuckDuckGoURL(searchTerm string) string {
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = url.QueryEscape(searchTerm)
	return fmt.Sprintf("%s?q=%s", duckDuckGoURL, searchTerm)
}

func duckDuckGoResultParsing(response *http.Response, rank int) ([]SearchResult, error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return nil, err
	}

	results := []SearchResult{}
	sel := doc.Find("div.result")
	rank++
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a.result__url")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("a.result__a")
		descTag := item.Find("div.result__snippet")
		desc := descTag.Text()
		title := titleTag.Text()
		link = strings.Trim(link, " ")

		if link != "" {
			result := SearchResult{
				rank,
				link,
				title,
				desc,
			}
			results = append(results, result)
			rank++
		}
	}
	return results, nil
}

func scrapeDuckDuckGo(searchTerm string, resultCount int) ([]SearchResult, error) {
	results := []SearchResult{}
	resultCounter := 0
	searchURL := buildDuckDuckGoURL(searchTerm)

	res, err := scrapeClientRequest(searchURL)
	if err != nil {
		return nil, err
	}

	data, err := duckDuckGoResultParsing(res, resultCounter)
	if err != nil {
		return nil, err
	}

	if resultCount < len(data) {
		data = data[:resultCount]
	}

	results = append(results, data...)
	return results, nil
}

func scrapeClientRequest(searchURL string) (*http.Response, error) {
	baseClient := &http.Client{}
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("scraper received a non-200 status code suggesting a ban")
		return nil, err
	}

	return res, nil
}

func printResultsToFile(results []SearchResult, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, result := range results {
		fmt.Fprintf(file, "Rank: %d\nTitle: %s\nURL: %s\nDescription: %s\n\n", result.ResultRank, result.ResultTitle, result.ResultURL, result.ResultDesc)
	}
}

func main() {
	var searchTerm string
	var userAgent string
	var outputFileName string
	var resultCount int

	flag.StringVar(&searchTerm, "search", "", "Search term")
	flag.StringVar(&searchTerm, "s", "", "Search term (shorthand)")
	flag.StringVar(&userAgent, "user-agent", "", "User agent string (optional)")
	flag.StringVar(&outputFileName, "output", "", "Output file name")
	flag.IntVar(&resultCount, "count", 10, "Number of search results")

	flag.Parse()

	if searchTerm == "" {
		fmt.Println("Please provide a search term using the -search or -s flag.")
		os.Exit(1)
	}

	if userAgent == "" {
		userAgent = randomUserAgent()
	}

	fmt.Printf("Using User-Agent: %s\n", userAgent)

	results, err := scrapeDuckDuckGo(searchTerm, resultCount)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if outputFileName == "" {
		outputFileName = searchTerm + ".txt"
	}

	printResultsToFile(results, outputFileName)
	fmt.Printf("Results saved to %s\n", outputFileName)
}
