package markers

import (
	"download-delegator/lib/parser/model"
)

type BaseHrefMarker struct {
	a uintptr
}

func (marker *BaseHrefMarker) GetName() string {
	return "base-marker"
}

func (marker *BaseHrefMarker) GetParameters() []model.MarkerParameter {
	return makeParameters(
		textParameter(ParamBaseHref, "Base Href", true),
	)
}

func (marker *BaseHrefMarker) Apply(data *model.PageData, markerData model.MarkerData) *model.PageData {
	baseHref := model.GetMarkerStringParameter(markerData, ParamBaseHref)

	head := data.Document.FindSingle("head")

	baseElem := head.PrependNewChild("base")
	baseElem.SetAttr("href", baseHref)

	return data
}
