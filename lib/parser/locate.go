package parser

import (
	"download-delegator/lib/parser/model"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Locate struct {
	models  []model.Model
	schemas []model.Schema
}

func (l Locate) LocateModel(url *string, models []model.Model) (*model.Model, error) {
	for _, m := range models {
		if m.UrlCheck == "" {
			continue
		}
		pattern := regexp.MustCompile(m.UrlCheck)

		if pattern.Match([]byte(*url)) {
			return &m, nil
		}
	}

	return nil, model.Error{
		Message: "model not found for url " + *url,
	}
}

func (l Locate) LocateSchema(schemaName string) (*model.Schema, error) {
	for _, schema := range l.schemas {
		if schema.Name == schemaName {
			return &schema, nil
		}
	}

	return nil, model.Error{
		Message: "schema not found for name " + schemaName,
	}
}

func (l *Locate) Init() {
	go func() {
		for {
			l.reload()
			time.Sleep(10 * time.Second)
		}
	}()
}

func (l *Locate) reload() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panicing while reloading locate: %s ", r)
		}
	}()
	log.Print("starting reload")

	var err error

	l.models, err = readModels()
	check(err)

	l.schemas, err = readSchemas()
	check(err)

	log.Printf("loaded %d / %d", len(l.models), len(l.schemas))
}
func (l *Locate) PrepareProcessData(p model.ProcessDataLight) (*model.ProcessData, error) {
	modelItem, err := l.LocateModel(p.Url, l.models)

	if err != nil {
		return nil, err
	}

	schema, err := l.LocateSchema(modelItem.Schema)

	if err != nil {
		return nil, err
	}

	processData := new(model.ProcessData)
	processData.Url = p.Url
	processData.Html = p.Html
	processData.Schema = schema
	processData.Model = modelItem
	processData.AdditionalModels = l.models

	return processData, nil
}

func readModels() ([]model.Model, error) {
	modelsUrl := "http://10.0.1.77:30006/api/1.0/models?token=eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ0YWxlaCIsImV4cCI6MTkxODQ2MTYwN30.GyMJCuJtB9nHEo7aa96dHmqrwVKQbbNJRC7FQkIKv8uCmlGTR7LvlBBE8ATKAwbLmGbKuVpD7GoOwgYu_wKM9w"

	resp, err := http.Get(modelsUrl)

	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var models []model.Model

	err = json.Unmarshal(respBytes, &models)

	if err != nil {
		return nil, err
	}
	return models, nil
}

func readSchemas() ([]model.Schema, error) {
	modelsUrl := "http://10.0.1.77:30006/api/1.0/schemas?token=eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ0YWxlaCIsImV4cCI6MTkxODQ2MTYwN30.GyMJCuJtB9nHEo7aa96dHmqrwVKQbbNJRC7FQkIKv8uCmlGTR7LvlBBE8ATKAwbLmGbKuVpD7GoOwgYu_wKM9w"

	resp, err := http.Get(modelsUrl)

	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var models []model.Schema

	err = json.Unmarshal(respBytes, &models)

	if err != nil {
		return nil, err
	}
	return models, nil
}
