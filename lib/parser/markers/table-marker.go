package markers

import (
	"download-delegator/lib/parser/model"
)

// on development
type TableMarker struct {
}

func (marker *TableMarker) GetName() string {
	return "table-marker"
}

func (marker *TableMarker) GetParameters() []model.MarkerParameter {
	return makeParameters(
		model.MarkerParameter{
			Name:          ParamTableLoc,
			Caption:       "Table location",
			DefaultValue:  true,
			ParameterType: model.TRANSITIVE_ELEMENT,
		})
}

func (marker *TableMarker) Apply(data *model.PageData, markerData model.MarkerData) *model.PageData {
	tableLocation := model.GetMarkerStringParameter(markerData, ParamTableLoc)

	for _, table := range data.Document.Find("[ug-field-t=" + tableLocation + "]") {
		marker.applyTable(table)
	}

	return data
}

func (marker *TableMarker) applyTable(table model.DocNode) {
	table.Find("thead > ")
}
