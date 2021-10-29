package markers

import (
	"download-delegator/lib/parser/model"
	"strings"
)

type MetaDataMarker struct {
}

func (marker *MetaDataMarker) GetName() string {
	return "meta-data"
}

func (marker *MetaDataMarker) GetParameters() []model.MarkerParameter {
	return makeParameters(
		textParameter(ParamName, "Name", false),
		inspectorParameter(ParamSelector, "Selector", false),
		model.MarkerParameter{
			Name:          ParamMetaTags,
			Caption:       "Extract Meta Tags",
			DefaultValue:  true,
			ParameterType: model.CHECKBOX,
		},
	)
}

func (marker *MetaDataMarker) Apply(data *model.PageData, markerData model.MarkerData) *model.PageData {
	metaTags := model.GetMarkerBooleanParameter(markerData, "meta-tags")

	if metaTags {
		for _, element := range data.Document.Find("meta") {
			key := strings.TrimSpace(element.Attr("name"))

			if key == "" {
				key = strings.TrimSpace(element.Attr("property"))
			}

			content := strings.TrimSpace(element.Attr("content"))

			if key != "" && content != "" {
				element.SetAttr("ug-field", "meta."+key)
				element.SetAttr("ug-value", content)
			}
		}
	}

	return data
}
