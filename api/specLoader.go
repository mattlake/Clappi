package api

import (
	"Clappi/constants"
	"errors"
	"fmt"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/utils"
	"os"
	"path/filepath"
	"strings"
)

func (am *APIManager) LoadSpecs() error {
	am.apis = make([]*API, 0)

	return filepath.Walk(am.specPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf(constants.FileAccessError, path, err)
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			return nil
		}

		api, err := am.loadSpec(path)
		if err != nil {
			am.apis = append(am.apis, &API{
				Name:      filepath.Base(path),
				FilePath:  path,
				LoadError: err,
			})
			return nil
		}

		am.apis = append(am.apis, api)
		return nil
	})
}

func (am *APIManager) loadSpec(path string) (*API, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(constants.FileReadingError, path, err)
	}

	doc, err := libopenapi.NewDocument(data)
	if err != nil {
		return nil, fmt.Errorf(constants.FileParsingError, path, err)
	}

	if doc.GetSpecInfo().SpecType == utils.OpenApi2 {
		return nil, errors.New(constants.UnsupportedVersionError)
	}

	model, errs := doc.BuildV3Model()
	if len(errs) > 0 {
		return nil, fmt.Errorf(constants.ModelBuildingError, path)
	}

	info := model.Model.Info
	title := filepath.Base(path)
	if info != nil && info.Title != "" {
		title = info.Title
	}

	return &API{
		Name:     title,
		Document: &doc,
		Model:    model,
		FilePath: path,
	}, nil
}
