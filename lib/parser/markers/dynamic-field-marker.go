package markers

import (
	"download-delegator/lib/parser/model"
)

type DynamicFieldMarker struct {
}

func (marker *DynamicFieldMarker) GetName() string {
	return "dynamic-field"
}

func (marker *DynamicFieldMarker) GetParameters() []model.MarkerParameter {
	return makeParameters(
		model.MarkerParameter{
			Name:          ParamTransitive,
			Caption:       "Transitive parameter",
			DefaultValue:  false,
			ParameterType: model.CHECKBOX,
			Required:      true,
		},
		model.MarkerParameter{
			Name:          ParamDynamicAttribute,
			Caption:       "Dynamic Attribute",
			DefaultValue:  "",
			ParameterType: model.TEXT,
			Required:      true,
		},
	)
}

func (marker *DynamicFieldMarker) Apply(data *model.PageData, markerData model.MarkerData) *model.PageData {
	dynamicAttribute := model.GetMarkerStringParameter(markerData, ParamDynamicAttribute)
	isTransitive := model.GetMarkerBooleanParameter(markerData, ParamTransitive)

	for _, element := range data.Document.Find("[" + dynamicAttribute + "]") {
		if !isTransitive {
			element.SetAttr("ug-field", element.Attr(dynamicAttribute))
		} else {
			element.SetAttr("ug-field-t", element.Attr(dynamicAttribute))
		}

		element.SetAttr("ug-marker", marker.GetName())
	}

	return data
}
