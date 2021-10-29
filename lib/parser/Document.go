package parser

import (
	"bufio"
	"bytes"
	"download-delegator/lib/parser/model"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
	"strings"
)

type HtmlDocNode struct {
	node *html.Node
}

func (docNode *HtmlDocNode) IsTextNode() bool {
	return docNode.node.Type == html.TextNode
}

func (docNode *HtmlDocNode) RemoveFromParent() {
	docNode.node.Parent.RemoveChild(docNode.node)
}

func (docNode *HtmlDocNode) PrependNewChild(tagName string) model.DocNode {
	node := new(html.Node)
	node.Type = html.ElementNode
	node.DataAtom = atom.Lookup([]byte(tagName))
	node.Data = tagName

	child := FromNode(node)

	docNode.PrependChild(child)

	return child
}

func (docNode *HtmlDocNode) FindSingle(selector string) model.DocNode {
	list := docNode.Find(selector)

	if len(list) == 0 {
		return nil
	}

	return list[0]
}

func (docNode *HtmlDocNode) Find(selector string) []model.DocNode {
	if strings.Contains(selector, ",") {
		var result []model.DocNode

		parts := strings.Split(selector, ",")

		for _, s := range parts {
			for _, item := range docNode.find(s) {
				result = append(result, item)
			}
		}

		return result
	} else {
		return docNode.find(selector)
	}
}

func (docNode *HtmlDocNode) find(selector string) []model.DocNode {
	selector = strings.Replace(selector, " ", "", -1)

	var result []model.DocNode

	queryNode := goquery.NewDocumentFromNode(docNode.node)

	for _, ele := range queryNode.Find(selector).Nodes {
		result = append(result, FromNode(ele))
	}

	return result
}

func (docNode *HtmlDocNode) Attr(key string) string {
	for _, attr := range docNode.node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func (docNode *HtmlDocNode) HasAttr(key string) bool {
	for _, attr := range docNode.node.Attr {
		if attr.Key == key {
			return true
		}
	}
	return false
}

func (docNode *HtmlDocNode) RemoveAttr(key string) {
	var filteredAttributes []html.Attribute

	for _, attr := range docNode.node.Attr {
		if attr.Key != key {
			filteredAttributes = append(filteredAttributes, attr)
		}
	}

	docNode.node.Attr = filteredAttributes
}

func (docNode *HtmlDocNode) SetAttr(key string, value string) {
	for _, attr := range docNode.node.Attr {
		if attr.Key == key {
			attr.Val = value
			return
		}
	}

	docNode.node.Attr = append(docNode.node.Attr, html.Attribute{
		Key: key,
		Val: value,
	})
}

func (docNode *HtmlDocNode) Parent() model.DocNode {
	return FromNode(docNode.node.Parent)
}

func (docNode *HtmlDocNode) NextElementSibling() model.DocNode {
	// find next element
	nextElem := docNode.node.NextSibling

	for nextElem != nil {
		if nextElem.Type == html.ElementNode {
			return FromNode(docNode.node.NextSibling)
		}

		nextElem = nextElem.NextSibling
	}

	return nil
}

func (docNode *HtmlDocNode) Clone() model.DocNode {
	return FromNode(docNode.node)
}

func (docNode *HtmlDocNode) PrependChild(element model.DocNode) {
	if docNode.node.FirstChild != nil {
		docNode.node.InsertBefore(element.Node(), docNode.node.FirstChild)
	} else {
		docNode.node.AppendChild(element.Node())
	}
}

func (docNode *HtmlDocNode) ChildNodes() []model.DocNode {
	var nodes []model.DocNode

	curChild := docNode.node.FirstChild

	for curChild != nil {
		nodes = append(nodes, FromNode(curChild))

		curChild = curChild.NextSibling
	}

	return nodes
}

func (docNode *HtmlDocNode) Text() string {
	var buf bytes.Buffer

	isFirst := true

	// Slightly optimized vs calling Each: no single selection object created
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			val := strings.TrimSpace(n.Data)

			if val == "" {
				return
			}

			if !isFirst {
				buf.WriteString(" ")
			}
			isFirst = false
			// Keep newlines and spaces, like jQuery
			buf.WriteString(val)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}

	f(docNode.node)

	return buf.String()
}

func (docNode *HtmlDocNode) Html(outer bool) string {
	queryNode := goquery.NewDocumentFromNode(docNode.node)

	var content string
	var e error

	if outer {
		content, e = goquery.OuterHtml(queryNode.Selection)

		check(e)
	} else {
		content, e = queryNode.Html()

		check(e)
	}

	return content
}

func (docNode *HtmlDocNode) IsLeafElement() bool {
	curChild := docNode.node.FirstChild

	for curChild != nil {
		if curChild.Type == html.ElementNode {
			return false
		}

		curChild = curChild.NextSibling
	}

	return true
}

func (docNode *HtmlDocNode) Node() *html.Node {
	return docNode.node
}

type TextDocNode struct {
}

func (docNode *HtmlDocNode) String() string {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := html.Render(w, docNode.node)

	if err != nil {
		log.Panic(err)
	}

	w.Flush()

	return b.String()
}

func ParseHtml(htmlContent string) (model.DocNode, error) {
	r := strings.NewReader(htmlContent)

	doc, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	htmlDocNode := FromNode(doc)

	return htmlDocNode, nil
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func FromNode(node *html.Node) model.DocNode {
	if node == nil {
		return nil
	}

	htmlDocNode := new(HtmlDocNode)
	htmlDocNode.node = node

	return htmlDocNode
}
