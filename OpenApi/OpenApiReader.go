package OpenApi

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadTitle(path string) (string, error) {
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(path)
	if err != nil {
		return "", fmt.Errorf("Failed to load spec: %v", err)
	}
	if doc.Info == nil || doc.Info.Title == "" {
		return "", fmt.Errorf("missing info.title in spec")
	}
	return doc.Info.Title, nil
}

func PrintTitle(path string) error {
	title, err := LoadTitle(path)
	if err != nil {
		return err
	}
	fmt.Println(title)
	return nil
}
