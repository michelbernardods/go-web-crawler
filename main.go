package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

var links []string
var urlFlag string
var visited map[string]bool = map[string]bool{}

func init() {
	flag.StringVar(&urlFlag, "url", "https://github.com/michelbernardods?tab=repositories", "url to start")
}

func main() {
	done := make(chan bool)
	go visitLinks(urlFlag)

	<-done

	fmt.Println(len(links), "links")
}

func saveLinks() {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err = os.Mkdir("data", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	fileLinks, err := os.Create("data/links.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer fileLinks.Close()

	for _, link := range links {
		_, err := fileLinks.WriteString(link + "\n") // adiciona quebra de linha após cada link
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func extractLinks(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}

			link, _ := url.Parse(attr.Val)
			invalidSchemes := []string{"", "mailto", "javascript", "about"}

			if containsString(invalidSchemes, link.Scheme) {
				continue
			}

			links = append(links, link.String())
			saveLinks()
			go visitLinks(link.String())

		}

	}

	for nav := n.FirstChild; nav != nil; nav = nav.NextSibling {
		extractLinks(nav)
	}
}

func visitLinks(url string) {
	if ok := visited[url]; ok {
		return
	}
	visited[url] = true

	fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("status not is 200: %d", res.StatusCode)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	extractLinks(doc)
}
