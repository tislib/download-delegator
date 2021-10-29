package markers

import "download-delegator/lib/parser/model"

var markerTypes = append([]model.Marker{},
	new(BaseHrefMarker),
	new(NormalizeMarker),
	new(FieldSelectorMarker),
	new(MetaDataMarker),
	new(ChildToParentTransform),
	new(DynamicFieldMarker),
)

func GetMarkerTypes() []model.Marker {
	return markerTypes
}
