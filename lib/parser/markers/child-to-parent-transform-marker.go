package markers

import (
	"download-delegator/lib/parser/model"
	"github.com/gosimple/slug"
)

type ChildToParentTransform struct {
}

func (marker *ChildToParentTransform) GetName() string {
	return "child-to-parent-transform"
}

func (marker *ChildToParentTransform) GetParameters() []model.MarkerParameter {
	return makeParameters(
		inspectorParameter(ParamChildSelector, "Child selector(from parent)", false),
		inspectorParameter(ParamParentSelector, "Parent selector", false),
		textParameter(ParamParentAttr, "Parent attribute", false),
	)
}

func (marker *ChildToParentTransform) Apply(page *model.PageData, markerData model.MarkerData) *model.PageData {
	parentSelector := model.GetMarkerStringParameter(markerData, "parent_selector")
	parentAttr := model.GetMarkerStringParameter(markerData, "parent_attr")
	childSelector := model.GetMarkerStringParameter(markerData, "child_selector")

	if parentSelector != "" && parentAttr != "" && childSelector != "" {
		parentElements := page.Document.Find(parentSelector)

		for _, parentElement := range parentElements {
			child := parentElement.FindSingle(childSelector)

			if child != nil {
				parentElement.SetAttr(parentAttr, slug.Make(child.Text()))
			}
		}
	}

	return page
}
