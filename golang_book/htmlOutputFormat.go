package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		outline(url)
	}
}

func outline(url string) error {
	if strings.Contains(url, "www.") != true {
		url = "www." + url
	}
	if strings.Contains(url, "http://") != true {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	fmt.Println(doc)
	if err != nil {
		return err
	}

	//!+call
	forEachNode(doc, startElement, endElement)
	//!-call

	return nil
}

//!+forEachNode
// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
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

//!-forEachNode

//!+startend
var depth int

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		if n.FirstChild == nil {
			fmt.Printf("%*s<%s/>\n", depth*2, "", n.Data)
			depth++
		} else {
			var attributes []string
			for _, attr := range n.Attr {
				attributes = append(attributes, attr.Key)
			}
			if len(attributes) == 1 {
				fmt.Printf("%*s<%s %s='...'>\n", depth*2, "", n.Data, attributes[0])
			} else if len(attributes) > 1 {
				var attr []string
				for _, a := range attributes {
					attr = append(attr, a+"='...'")
				}
				trimmedAttr := strings.Trim(strings.Join(attr, " "), "[]")
				fmt.Printf("%*s<%s %s>\n", depth*2, "", n.Data, trimmedAttr)
			}
			depth++
		}
	} else if n.Type == html.CommentNode {
		fmt.Printf("%s\n", n.Data)
	} else if n.Type == html.TextNode {
		fmt.Printf("%v\n", n.Data)
	}

}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode && n.LastChild != nil {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}
