package api

import (
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type API struct {
	Name      string
	Document  *libopenapi.Document
	Model     *libopenapi.DocumentModel[v3.Document]
	FilePath  string
	LoadError error
}

type APIManager struct {
	apis     []*API
	specPath string
}

func NewAPIManager(specPath string) *APIManager {
	return &APIManager{
		specPath: specPath,
		apis:     make([]*API, 0),
	}
}

func (am *APIManager) GetAPIs() []*API {
	return am.apis
}
