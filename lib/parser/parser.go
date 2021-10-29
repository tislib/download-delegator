package parser

import (
	"download-delegator/lib/parser/model"
	"github.com/google/uuid"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct {
	data     model.ProcessData
	document model.DocNode
}

func (p *Parser) Parse(data model.ProcessData) *model.Record {
	p.data = data

	processor := new(Processor)

	// parse html
	doc, err := ParseHtml(*data.Html)

	if err != nil {
		log.Panic(err)
	}

	pageData := model.PageData{
		Document: doc,
		Url:      *data.Url,
	}

	docNode := processor.ProcessWithModel(data.Model, &pageData)
	p.document = docNode

	return p.ParseWithDocNode(docNode)
}

func (p Parser) ParseWithDocNode(processedDocument model.DocNode) *model.Record {
	recordData := p.Extract(p.data.Schema, processedDocument)
	recordMeta := p.ExtractMeta()

	record := new(model.Record)
	record.Data = recordData
	record.Meta = recordMeta
	record.Tags = append([]string{})

	record.Schema = p.data.Model.Schema
	record.Source = p.data.Model.Source
	record.SourceUrl = *p.data.Url
	record.Ref = p.extractRef(p.data.Model, *p.data.Url)
	record.ObjectType = p.data.Model.ObjectType

	id, err := uuid.Parse(NamedUUID([]byte(*p.data.Url)))

	check(err)

	record.Id = id

	record.Name = Coalesce(recordData["name"], recordMeta["title"], recordMeta["og:title"])
	record.Description = Coalesce(recordData["description"], recordMeta["description"], recordMeta["og:description"])

	return record
}

func (p Parser) Extract(objectProperties model.ObjectProperty, parent model.DocNode) model.RecordData {
	var data = make(model.RecordData)

	for key, property := range objectProperties.GetProperties() {
		value := p.locatePropertyValue(parent, key, &property)

		if value != nil {
			data[key] = value
		}
	}

	return data
}

func (p Parser) locatePropertyValue(parent model.DocNode, key string, property *model.SchemaProperty) model.Value {
	fields := parent.Find("[ug-field=\"" + key + "\"]")

	return p.locatePropertyValueForField(property, fields)
}

func (p Parser) ExtractMeta() model.RecordMeta {
	metaData := make(model.RecordMeta)

	metaFields := p.document.Find("meta[ug-field]")

	for _, metaField := range metaFields {
		key := metaField.Attr("ug-field")
		value := metaField.Attr("ug-value")

		if strings.HasPrefix(key, "meta.") {
			key = key[len("meta."):]
		}

		metaData[key] = value
	}

	return metaData
}

func (p Parser) locatePropertyValueForField(property *model.SchemaProperty, fields []model.DocNode) model.Value {
	if property.IsArrayProperty() {
		itemsProperty := property.Items

		var result []model.Value

		for _, field := range fields {
			val := p.locatePropertyValueField(itemsProperty, field)
			result = append(result, val)
		}

		return result
	} else {
		if len(fields) > 1 && property.IsStringProperty() {
			result := new(strings.Builder)

			isFirst := true
			for _, field := range fields {
				val := strings.TrimSpace(p.locatePropertyValueField(property, field).(string))

				if val == "" {
					continue
				}
				if !isFirst {
					result.WriteString(" ")
				}

				isFirst = false
				result.WriteString(val)
			}

			return result.String()
		} else if len(fields) > 0 {
			return p.locatePropertyValueField(property, fields[0])
		} else {
			return nil
		}
	}
}

func (p Parser) locatePropertyValueField(property *model.SchemaProperty, field model.DocNode) model.Value {
	if property.IsStringProperty() {
		return p.getValue(field)
	} else if property.IsNumberProperty() {
		val := p.getValue(field)

		valStrNum := string(regexp.MustCompile("[^\\d.]+").ReplaceAll([]byte(val), []byte("")))

		if len(valStrNum) == 0 {
			return nil
		} else {
			number, err := strconv.ParseFloat(valStrNum, 64)

			check(err)

			return number
		}
	} else if property.IsArrayProperty() {
		return p.locatePropertyValueForField(property, append([]model.DocNode{}, field))
	} else if property.IsObjectProperty() {
		return p.Extract(property, field)
	} else if property.IsReferenceProperty() {
		return p.extractReference(property, field)
	} else {
		log.Print("invalid property type: ", property.Type)
		return nil
	}

}

func (p Parser) extractReference(referenceProperty *model.SchemaProperty, field model.DocNode) *model.Reference {
	reference := new(model.Reference)

	reference.Name = p.getValue(field)
	if !field.HasAttr("href") {
		return reference
	}

	href := field.Attr("href")

	if !strings.HasPrefix(href, "http") {
		if strings.HasPrefix(href, "/") {
			pageUrlObj, err := url.Parse(*p.data.Url)

			check(err)

			href = pageUrlObj.Scheme + "://" + pageUrlObj.Host + href
		}
	}

	schemaName := referenceProperty.Schema

	m := p.locateModel(schemaName)

	if m == nil {
		return reference
	}

	reference.Source = m.Source
	reference.SourceUrl = href
	reference.Ref = p.extractRef(m, href)
	reference.ObjectType = m.ObjectType

	return reference
}

func (p Parser) getValue(field model.DocNode) string {
	if field.HasAttr("ug-value") {
		return strings.TrimSpace(field.Attr("ug-value"))
	}

	return strings.TrimSpace(field.Text())
}

func (p *Parser) locateModel(schemaName string) *model.Model {
	for _, item := range p.data.AdditionalModels {
		if item.Schema == schemaName {
			return &item
		}
	}

	return nil
}

func (p *Parser) extractRef(m *model.Model, href string) string {
	ref := m.Ref

	if ref != "" {
		refExp, err := regexp.Compile(ref)

		check(err)

		h := refExp.FindAllSubmatchIndex([]byte(href), -1)

		if len(h) > 0 {
			return href[h[0][2]:h[0][3]]
		}
	}

	return ""
}
