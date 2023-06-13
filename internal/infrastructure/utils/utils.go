package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"

	"go-parser/internal/dto"
)

// traverse is a template function for traversing tree of Nodes
func traverse(n *html.Node, f func(n *html.Node)) {
	f(n)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, f)
	}
}

// extractText is a utility function
// for extracting text from child Text Nodes
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var text string
	for curr := n.FirstChild; curr != nil; curr = curr.NextSibling {
		text += extractText(curr)
	}
	return text
}

// checkClass is a utility function for checking
// if HTML tag has class attribute with 'value' value
func checkClass(n *html.Node, value string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && attr.Val == value {
			return true
		}
	}
	return false
}

// getPagesNum is a utility function that returns
// number of pages to parse
func getPagesNum(url string) (int, error) {
	page, err := GetHtmlPage(url)
	if err != nil {
		return 0, nil
	}

	doc, err := html.Parse(strings.NewReader(string(page)))
	if err != nil {
		return 0, nil
	}

	pagination := FindNode(doc, "tg-pagination")

	maxPage := 0

	for link := FindNode(pagination, "active"); link != nil; link = link.NextSibling {
		text := strings.Trim(extractText(link), "\t \n")

		num, err := strconv.Atoi(text)
		if err != nil {
			continue
		}

		if num > maxPage {
			maxPage = num
		}
	}

	return maxPage, nil
}

func GetHtmlPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("page is not found")
	}

	page, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// FindNode finds node by its class attribute value
func FindNode(n *html.Node, value string) *html.Node {
	var node *html.Node

	traverse(n, func(n *html.Node) {
		if n.Type == html.ElementNode && checkClass(n, value) {
			node = n
		}
	})
	return node
}

// FindNextPageLink is a function that returns
// next page link for current page
func FindNextPageLink(n *html.Node) string {
	var link string

	traverse(n, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" && checkClass(n, "tg-tg-nextpage") {
			for _, attr := range n.FirstChild.Attr {
				if attr.Key == "href" {
					link = attr.Val
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if next := FindNextPageLink(c); next != "" {
				link = next
			}
		}
	})
	return link
}

// GetPagesLinks is a function for getting all pages links
func GetPagesLinks(parseUrl string) []string {
	pages := []string{parseUrl}

	num, err := getPagesNum(parseUrl)
	if err != nil {
		log.Fatal(err)
	}

	queue := make(chan string, 1)

	var wg sync.WaitGroup
	wg.Add(num - 1)
	for i := 2; i <= num; i++ {
		go func(i int) {
			url := parseUrl + "&pid=" + strconv.Itoa(i)
			queue <- url
		}(i)
	}

	go func() {
		for url := range queue {
			pages = append(pages, url)
			wg.Done()
		}
	}()

	wg.Wait()

	return pages
}

// GetN	odesByClass gets all nodes with class value
// equal to 'classValue' on HTML page and returns them
func GetNodesByClass(n *html.Node, classValue string) []*html.Node {
	nodes := make([]*html.Node, 0)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found := FindNode(c, classValue)
		if found != nil {
			nodes = append(nodes, found)
		}
	}
	return nodes
}

// ParseBookCards function parses all book cards nodes to
// retrieve necessary information, puts retrieved data into
// the Book struct and returns slice of parsed Books
func ParseBookCards(bookCards []*html.Node) []*dto.Book {
	var books []*dto.Book

	for _, card := range bookCards {
		titleNode := FindNode(card, "tg-booktitle")
		authorNode := FindNode(card, "tg-bookwriter")
		priceNode := FindNode(card, "tg-bookprice")

		// Delete unnecessary symbols from string and convert book price to int
		title := strings.Trim(extractText(titleNode), "\t \n")
		author := strings.Trim(extractText(authorNode), "\t \n")

		// If last symbol of author string is non-break space
		// change value of author string to "Нет"
		if author[1] == 160 {
			author = "Не указан"
		}

		// If book is not in stock, its price is 0
		var price int
		if priceNode != nil {
			price, _ = strconv.Atoi(strings.Trim(extractText(priceNode), "\t ₽\n"))
		} else {
			title = title + " (НЕТ В НАЛИЧИИ)"
		}

		book := &dto.Book{
			Title:  title,
			Author: author,
			Price:  price,
		}

		books = append(books, book)
	}

	return books
}
