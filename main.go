package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	webpages := map[string]bool{}
	baseURL := os.Args[1]
	fmt.Println("Going to visit: " + baseURL)
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Error with http.Get:", err)
		os.Exit(1)
	}

	webpages[baseURL+"/"] = true

	parsed, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		os.Exit(1)
	}
	resp.Body.Close()

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && (strings.HasPrefix(a.Val, "/") || (strings.HasPrefix(a.Val, baseURL))) {
					if !webpages[a.Val] {
						webpages[a.Val] = true
						url := ""
						if strings.HasPrefix(a.Val, "/") {
							url = baseURL + a.Val
						} else {
							url = a.Val
						}
						fmt.Println("Going to visit: " + url)
						resp, _ := http.Get(url)
						parsed, _ := html.Parse(resp.Body)
						f(parsed)
					} else {
						if _, ok := webpages[a.Val]; !ok {
							webpages[a.Val] = false
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(parsed)

	fmt.Println(webpages)
}
