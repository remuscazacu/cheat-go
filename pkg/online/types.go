package online

import (
	"cheat-go/pkg/apps"
	"time"
)

type Repository struct {
	URL         string    `json:"url" yaml:"url"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	LastUpdated time.Time `json:"last_updated" yaml:"last_updated"`
	Stars       int       `json:"stars" yaml:"stars"`
	Author      string    `json:"author" yaml:"author"`
}

type CheatSheet struct {
	ID          string    `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Description string    `json:"description" yaml:"description"`
	App         apps.App  `json:"app" yaml:"app"`
	Repository  string    `json:"repository" yaml:"repository"`
	Downloads   int       `json:"downloads" yaml:"downloads"`
	Rating      float64   `json:"rating" yaml:"rating"`
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" yaml:"updated_at"`
	Tags        []string  `json:"tags" yaml:"tags"`
}

type SearchOptions struct {
	Query      string   `json:"query" yaml:"query"`
	Tags       []string `json:"tags" yaml:"tags"`
	Repository string   `json:"repository" yaml:"repository"`
	MinRating  float64  `json:"min_rating" yaml:"min_rating"`
	SortBy     string   `json:"sort_by" yaml:"sort_by"`
	Limit      int      `json:"limit" yaml:"limit"`
	Offset     int      `json:"offset" yaml:"offset"`
}

type SyncConfig struct {
	AutoSync     bool          `json:"auto_sync" yaml:"auto_sync"`
	SyncInterval time.Duration `json:"sync_interval" yaml:"sync_interval"`
	LastSync     time.Time     `json:"last_sync" yaml:"last_sync"`
	Repositories []string      `json:"repositories" yaml:"repositories"`
}

type Client interface {
	GetRepositories() ([]Repository, error)
	SearchCheatSheets(opts SearchOptions) ([]CheatSheet, error)
	GetCheatSheet(id string) (*CheatSheet, error)
	DownloadCheatSheet(id string) (*apps.App, error)
	SubmitCheatSheet(sheet CheatSheet) error
	RateCheatSheet(id string, rating float64) error
}
