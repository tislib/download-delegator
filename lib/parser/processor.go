package parser

import (
	"download-delegator/lib/parser/markers"
	"download-delegator/lib/parser/model"
	"log"
)

type Processor struct {
}

func (p Processor) ProcessData(data model.ProcessData) string {
	// parse html
	doc, err := ParseHtml(*data.Html)

	if err != nil {
		log.Panic(err)
	}

	pageData := model.PageData{
		Document: doc,
		Url:      *data.Url,
	}

	docNode := p.ProcessWithModel(data.Model, &pageData)

	return docNode.String()
}

func (p Processor) ProcessWithModel(model *model.Model, data *model.PageData) model.DocNode {
	// apply markers

	for _, markerData := range model.Markers {
		// locate marker
		marker := locateMarkerByType(markerData.Type)

		if marker == nil {
			//log.Printf("markerType %s not found", markerData.Type)
			//continue
			log.Panicf("markerType %s not found", markerData.Type)
		}

		appliedPageData := marker.Apply(data, markerData)

		if appliedPageData != nil {
			data = appliedPageData
		}
	}

	return data.Document
}

func locateMarkerByType(name string) model.Marker {
	for _, markerType := range markers.GetMarkerTypes() {
		if markerType.GetName() == name {
			return markerType
		}
	}

	return nil
}
