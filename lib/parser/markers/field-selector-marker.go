package markers

import (
	"download-delegator/lib/parser/model"
	"strings"
)

const PARAM_NAME = "name"
const PARAM_SELECTOR = "selector"
const PARAM_OUTPUT_TYPE = "outputType"
const PARAM_TRANSITIVE = "transitive"

type FieldSelectorMarker struct {
}

func (marker *FieldSelectorMarker) GetName() string {
	return "field-selector"
}

func (marker *FieldSelectorMarker) GetParameters() []model.MarkerParameter {
	return makeParameters(
		textParameter(PARAM_NAME, "Name", false),
		inspectorParameter(PARAM_SELECTOR, "Selector", false),
		model.MarkerParameter{
			Name:          PARAM_OUTPUT_TYPE,
			Caption:       "Output type",
			DefaultValue:  "text",
			ParameterType: model.TEXT,
			Required:      true,
		},
		model.MarkerParameter{
			Name:          PARAM_TRANSITIVE,
			Caption:       "Transitive parameter",
			DefaultValue:  false,
			ParameterType: model.CHECKBOX,
			Required:      true,
		},
	)
}

func (marker *FieldSelectorMarker) Apply(data *model.PageData, markerData model.MarkerData) *model.PageData {
	selector := model.GetMarkerStringParameter(markerData, "selector")

	if selector == "" {
		return nil
	}

	for _, element := range data.Document.Find(selector) {
		marker.applyParameter(element, markerData)
	}

	return data
}

func (marker *FieldSelectorMarker) applyParameter(element model.DocNode, markerData model.MarkerData) {
	fieldName := model.GetMarkerStringParameter(markerData, "name")
	outputType := model.GetMarkerStringParameter(markerData, "outputType")
	isTransitive := model.GetMarkerBooleanParameter(markerData, PARAM_TRANSITIVE)

	if !isTransitive {
		element.SetAttr("ug-field", fieldName)
	} else {
		element.SetAttr("ug-field-t", fieldName)
	}

	element.SetAttr("ug-marker", marker.GetName())

	marker.applyValueIf(element, outputType)
}

func (marker *FieldSelectorMarker) applyValueIf(element model.DocNode, outputType string) {
	var value string

	if outputType == "img" {
		value = element.Attr("src")
	} else if outputType == "text" {
		value = element.Text()
	} else if outputType == "html" {
		value = element.Html(false)
	} else if outputType == "outerHtml" {
		value = element.Html(true)
	} else if strings.HasPrefix(outputType, "attr:") {
		attr := outputType[len("attr:"):]

		value = element.Attr(attr)
	}

	if value != "" {
		element.SetAttr("ug-value", value)
	}
}
