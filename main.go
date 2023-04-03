package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

var links []string
var visited map[string]bool = map[string]bool{}

func main() {
	visitLinks("https://github.com/michelbernardods?tab=repositories")

	fmt.Println(len(links), "links")
}

func extractLinks(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}

			link, err := url.Parse(attr.Val)
			if err != nil || link.Scheme == "" {
				continue
			}

			links = append(links, link.String())
			visitLinks(link.String())

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
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("status not is 200: %d", res.StatusCode))
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		panic(err)
	}

	extractLinks(doc)
}
