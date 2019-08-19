package project

import (
	"github.com/h3poteto/yadockeri/app/repositories/project_values"
	"github.com/h3poteto/yadockeri/app/repositories/projects"
	"github.com/h3poteto/yadockeri/app/values"
	"github.com/h3poteto/yadockeri/db"
)

type OverrideValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Project struct {
	ID                int              `json:"id"`
	UserID            int              `json:"user_id"`
	Title             string           `json:"title"`
	BaseURL           string           `json:"base_url"`
	RepositoryOwner   string           `json:"repository_owner"`
	RepositoryName    string           `json:"repository_name"`
	HelmRepositoryUrl string           `json:"helm_repository_url"`
	HelmDirectoryName string           `json:"helm_directory_name"`
	Namespace         string           `json:"namespace"`
	ValueOptions      []*OverrideValue `json:"values"`
}

func GetProjects() ([]*Project, error) {
	p := projects.New(db.SharedInstance().Connection)
	projects, err := p.All()
	if err != nil {
		return nil, err
	}
	valuesRepository := project_values.New(db.SharedInstance().Connection)

	var results []*Project
	for _, proj := range projects {
		overrides, err := valuesRepository.GetByProject(proj.ID)
		if err != nil {
			return nil, err
		}
		proj.ValueOptions = overrides

		var values []*OverrideValue
		for _, value := range proj.ValueOptions {
			v := &OverrideValue{
				Key:   value.Key,
				Value: value.Value,
			}
			values = append(values, v)
		}

		p := &Project{
			ID:                proj.ID,
			UserID:            proj.UserID,
			Title:             proj.Title,
			BaseURL:           proj.BaseURL,
			RepositoryOwner:   proj.RepositoryOwner,
			RepositoryName:    proj.RepositoryName,
			HelmRepositoryUrl: proj.HelmRepositoryUrl,
			HelmDirectoryName: proj.HelmDirectoryName,
			Namespace:         proj.Namespace,
			ValueOptions:      values,
		}
		results = append(results, p)
	}
	return results, nil
}

func CreateProject(userID int, title, baseURL, owner, name, helmRepositoryURL, helmDirectory, namespace string, valueOptions []*values.OverrideValue) (*Project, error) {
	transaction, err := db.SharedInstance().Connection.Begin()
	projectRepository := projects.New(db.SharedInstance().Connection)
	id, err := projectRepository.Create(transaction, userID, title, baseURL, owner, name, helmRepositoryURL, helmDirectory, namespace)
	if err != nil {
		return nil, err
	}

	valuesRepository := project_values.New(db.SharedInstance().Connection)
	for _, value := range valueOptions {
		_, err := valuesRepository.Create(transaction, id, value.Key, value.Value)
		if err != nil {
			return nil, err
		}
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	proj, err := projectRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	overrides, err := valuesRepository.GetByProject(proj.ID)
	if err != nil {
		return nil, err
	}
	proj.ValueOptions = overrides

	var values []*OverrideValue
	for _, value := range proj.ValueOptions {
		v := &OverrideValue{
			Key:   value.Key,
			Value: value.Value,
		}
		values = append(values, v)
	}

	return &Project{
		ID:                proj.ID,
		UserID:            proj.UserID,
		Title:             proj.Title,
		BaseURL:           proj.BaseURL,
		RepositoryOwner:   proj.RepositoryOwner,
		RepositoryName:    proj.RepositoryName,
		HelmRepositoryUrl: proj.HelmRepositoryUrl,
		HelmDirectoryName: proj.HelmDirectoryName,
		Namespace:         proj.Namespace,
		ValueOptions:      values,
	}, nil
}

func GetProjectByID(id int) (*Project, error) {
	repository := projects.New(db.SharedInstance().Connection)
	p, err := repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	valuesRepository := project_values.New(db.SharedInstance().Connection)
	overrides, err := valuesRepository.GetByProject(p.ID)
	if err != nil {
		return nil, err
	}
	p.ValueOptions = overrides

	var values []*OverrideValue
	for _, value := range p.ValueOptions {
		v := &OverrideValue{
			Key:   value.Key,
			Value: value.Value,
		}
		values = append(values, v)
	}
	return &Project{
		ID:                p.ID,
		UserID:            p.UserID,
		Title:             p.Title,
		BaseURL:           p.BaseURL,
		RepositoryOwner:   p.RepositoryOwner,
		RepositoryName:    p.RepositoryName,
		HelmRepositoryUrl: p.HelmRepositoryUrl,
		HelmDirectoryName: p.HelmDirectoryName,
		Namespace:         p.Namespace,
		ValueOptions:      values,
	}, nil
}
