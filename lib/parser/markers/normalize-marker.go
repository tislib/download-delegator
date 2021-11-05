package markers

import (
	"download-delegator/lib/parser/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"regexp"
	"strconv"
	"strings"
)

type NormalizeMarker struct {
}

func (marker *NormalizeMarker) GetName() string {
	return "normalize"
}

func (marker *NormalizeMarker) GetParameters() []model.MarkerParameter {
	return makeParameters(
		model.MarkerParameter{
			Name:          ParamElement,
			Caption:       "Element",
			ParameterType: model.INSPECTOR,
			Required:      true,
		},
	)
}

func (marker *NormalizeMarker) Apply(data *model.PageData, markerData model.MarkerData) *model.PageData {
	elementSelector := model.GetMarkerStringParameter(markerData, ParamElement)

	if elementSelector != "" {
		for _, selection := range data.Document.Find(elementSelector) {
			marker.wrapText(selection)
			marker.normalizeTableSpan(selection)
			marker.normalizeId(selection)
		}
	}

	return data
}

func (marker *NormalizeMarker) wrapText(selection model.DocNode) {
	for _, node := range selection.ChildNodes() {
		if node.IsTextNode() {
			text := strings.TrimSpace(selection.Text())

			if text == "" {
				continue
			}

			if node.Parent().IsLeafElement() {
				continue
			}

			newElem := &html.Node{}
			newElem.Type = html.ElementNode
			newElem.Data = "text"
			parent := node.Node().Parent

			parent.InsertBefore(newElem, node.Node())
			node.RemoveFromParent()
			newElem.AppendChild(node.Node())
		}
	}
}

func (marker *NormalizeMarker) normalizeTableSpan(element model.DocNode) {
	elementsWithRowspan := element.Find("td[rowspan],th[rowspan]")

	for _, item := range elementsWithRowspan {
		marker.fixRowSpan(item)
	}
}

func (marker *NormalizeMarker) fixRowSpan(element model.DocNode) {
	rowSpan, err := strconv.Atoi(element.Attr("rowspan"))

	check(err)

	element.RemoveAttr("rowspan")
	tr := element.Parent()
	nextSibling := tr.NextElementSibling()

	for i := 0; i < rowSpan-1; i++ {
		if nextSibling != nil {
			nextSibling.PrependChild(element.Clone())
			nextSibling = nextSibling.NextElementSibling()
		}
	}
}

func (marker *NormalizeMarker) normalizeId(selection model.DocNode) {
	var invalidIdPattern, err = regexp.Compile("^\\d[\\w\\W]+")

	check(err)

	for _, selection := range selection.Find("[id]") {
		attr := selection.Attr("id")

		if invalidIdPattern.Match([]byte(attr)) {
			selection.RemoveAttr(attr)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
