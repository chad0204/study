package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

//限制Extract同时最多只有20的并发量
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Printf("get url: %v \n", url)
	//acquire
	tokens <- struct{}{}
	list, err := Extract(url)
	//release
	<-tokens
	if err != nil {
		log.Print(err)
	}
	return list
}

//串行广度优先
func breadthFirst(crawl func(item string) []string, workList []string) {
	seen := make(map[string]bool)
	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				workList = append(workList, crawl(item)...)
			}
		}
	}
}

func main() {
	//breadthFirst(crawl, os.Args[1:])
	//下面的方式是取代breadthFirst函数
	workList := make(chan []string)

	n := 1
	go func() {
		workList <- os.Args[1:]
	}()

	//保存已经遍历过的link
	seen := make(map[string]bool)

	//for list := range workList {//会导致无法退出, 因为不知道何时close, 使用n计数退出
	for ; n > 0; n-- {
		list := <-workList
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(url string) {
					workList <- crawl(url)
				}(link)
			}
		}
	}

}
