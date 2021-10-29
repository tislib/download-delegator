package model

import (
	"fmt"
	"golang.org/x/net/html"
)

type DocNode interface {
	fmt.Stringer

	Find(selector string) []DocNode

	Attr(s string) string

	RemoveAttr(attr string)

	Parent() DocNode

	NextElementSibling() DocNode

	Clone() DocNode

	PrependChild(element DocNode)

	PrependNewChild(tagName string) DocNode

	ChildNodes() []DocNode

	Text() string

	Html(outer bool) string

	IsLeafElement() bool

	Node() *html.Node

	SetAttr(key string, value string)

	FindSingle(s string) DocNode

	IsTextNode() bool

	RemoveFromParent()

	HasAttr(s string) bool
}
