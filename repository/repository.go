package repository

import (
	"errors"
)

type ProjectSettings struct {
	ProjectsDirectory string
	FlaxLocation      string
}

type Project struct {
	Name        string
	Preview     string
	ProjectFile string
}

var (
	errUpdateFailed = errors.New("update failed")
)

type Repository interface {
	Migrate() error
	InsertFlaxerSettings(h FlaxerSettings) (*FlaxerSettings, error)
	GetFlaxerSettings() (*FlaxerSettings, error)
	UpdateFlaxerSettings(id int64, update FlaxerSettings) error
}

type FlaxerSettings struct {
	ID                int64  `json:"id"`
	ProjectsDirectory string `json:"projectsdirectory"`
	FlaxLocation      string `json:"flaxlocation"`
}
