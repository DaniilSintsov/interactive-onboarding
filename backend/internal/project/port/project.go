package port

import (
	"strings"

	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/entity"
	"github.com/DaniilSintsov/interactive-onboarding/backend/internal/project/errs"
)

type CreateProjectParams struct {
	Name     string
	Elements []CreateProjectElementParams
}

type CreateProjectElementParams struct {
	Key         string
	Label       string
	Description string
	Page        string
}

func (params *CreateProjectParams) Normalize() {
	params.Name = strings.TrimSpace(params.Name)
	for i := range params.Elements {
		params.Elements[i].Normalize()
	}
}

func (params *CreateProjectElementParams) Normalize() {
	params.Key = strings.TrimSpace(params.Key)
	params.Label = strings.TrimSpace(params.Label)
	params.Description = strings.TrimSpace(params.Description)
	params.Page = strings.TrimSpace(params.Page)
}

func (params *CreateProjectParams) Validate() error {
	params.Normalize()
	if params.Name == "" {
		return errs.ErrProjectNameRequired
	}
	for _, element := range params.Elements {
		err := element.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (params *CreateProjectElementParams) Validate() error {
	return entity.CheckElementFields(params.Key, params.Label, params.Description, params.Page)
}

type ListProjectsResult struct {
	Projects []entity.Project
	Total    int64
	Limit    int
	Offset   int
}

type ProjectWithElements struct {
	Project  entity.Project
	Elements []entity.Element
}
