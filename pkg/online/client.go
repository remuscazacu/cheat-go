package online

import (
	"bytes"
	"cheat-go/pkg/apps"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	cache      *cache
	mu         sync.RWMutex
}

type cache struct {
	repositories []Repository
	cheatSheets  map[string]*CheatSheet
	lastUpdated  time.Time
	ttl          time.Duration
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: &cache{
			cheatSheets: make(map[string]*CheatSheet),
			ttl:         15 * time.Minute,
		},
	}
}

func (c *HTTPClient) GetRepositories() ([]Repository, error) {
	c.mu.RLock()
	if time.Since(c.cache.lastUpdated) < c.cache.ttl && len(c.cache.repositories) > 0 {
		repos := c.cache.repositories
		c.mu.RUnlock()
		return repos, nil
	}
	c.mu.RUnlock()

	resp, err := c.httpClient.Get(c.baseURL + "/api/repositories")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %w", err)
	}

	c.mu.Lock()
	c.cache.repositories = repos
	c.cache.lastUpdated = time.Now()
	c.mu.Unlock()

	return repos, nil
}

func (c *HTTPClient) SearchCheatSheets(opts SearchOptions) ([]CheatSheet, error) {
	params := url.Values{}
	if opts.Query != "" {
		params.Add("q", opts.Query)
	}
	if len(opts.Tags) > 0 {
		params.Add("tags", strings.Join(opts.Tags, ","))
	}
	if opts.Repository != "" {
		params.Add("repository", opts.Repository)
	}
	if opts.MinRating > 0 {
		params.Add("min_rating", fmt.Sprintf("%.1f", opts.MinRating))
	}
	if opts.SortBy != "" {
		params.Add("sort", opts.SortBy)
	}
	if opts.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset > 0 {
		params.Add("offset", fmt.Sprintf("%d", opts.Offset))
	}

	url := fmt.Sprintf("%s/api/cheatsheets?%s", c.baseURL, params.Encode())
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search cheat sheets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var sheets []CheatSheet
	if err := json.NewDecoder(resp.Body).Decode(&sheets); err != nil {
		return nil, fmt.Errorf("failed to decode cheat sheets: %w", err)
	}

	return sheets, nil
}

func (c *HTTPClient) GetCheatSheet(id string) (*CheatSheet, error) {
	c.mu.RLock()
	if cached, exists := c.cache.cheatSheets[id]; exists {
		c.mu.RUnlock()
		return cached, nil
	}
	c.mu.RUnlock()

	resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/cheatsheets/%s", c.baseURL, id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cheat sheet: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("cheat sheet not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var sheet CheatSheet
	if err := json.NewDecoder(resp.Body).Decode(&sheet); err != nil {
		return nil, fmt.Errorf("failed to decode cheat sheet: %w", err)
	}

	c.mu.Lock()
	c.cache.cheatSheets[id] = &sheet
	c.mu.Unlock()

	return &sheet, nil
}

func (c *HTTPClient) DownloadCheatSheet(id string) (*apps.App, error) {
	sheet, err := c.GetCheatSheet(id)
	if err != nil {
		return nil, err
	}

	return &sheet.App, nil
}

func (c *HTTPClient) SubmitCheatSheet(sheet CheatSheet) error {
	data, err := json.Marshal(sheet)
	if err != nil {
		return fmt.Errorf("failed to marshal cheat sheet: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/api/cheatsheets",
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("failed to submit cheat sheet: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to submit cheat sheet: %s", body)
	}

	return nil
}

func (c *HTTPClient) RateCheatSheet(id string, rating float64) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	data, _ := json.Marshal(map[string]float64{"rating": rating})

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/cheatsheets/%s/rate", c.baseURL, id),
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to rate cheat sheet: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to rate cheat sheet: %s", body)
	}

	delete(c.cache.cheatSheets, id)

	return nil
}

type MockClient struct {
	repositories []Repository
	cheatSheets  []CheatSheet
	mu           sync.RWMutex
}

func NewMockClient() *MockClient {
	return &MockClient{
		repositories: defaultRepositories(),
		cheatSheets:  defaultCheatSheets(),
	}
}

func defaultRepositories() []Repository {
	return []Repository{
		{
			URL:         "https://github.com/cheat-go/community",
			Name:        "Official Community Repository",
			Description: "Official cheat sheets maintained by the community",
			LastUpdated: time.Now().Add(-24 * time.Hour),
			Stars:       1250,
			Author:      "cheat-go",
		},
		{
			URL:         "https://github.com/awesome/cheatsheets",
			Name:        "Awesome Cheat Sheets",
			Description: "A curated collection of awesome cheat sheets",
			LastUpdated: time.Now().Add(-48 * time.Hour),
			Stars:       890,
			Author:      "awesome",
		},
	}
}

func defaultCheatSheets() []CheatSheet {
	return []CheatSheet{
		{
			ID:          "vim-advanced",
			Name:        "Vim Advanced",
			Description: "Advanced Vim shortcuts and commands",
			Repository:  "https://github.com/cheat-go/community",
			Downloads:   5420,
			Rating:      4.8,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-2 * 24 * time.Hour),
			Tags:        []string{"vim", "editor", "advanced"},
		},
		{
			ID:          "git-workflow",
			Name:        "Git Workflow",
			Description: "Complete Git workflow commands",
			Repository:  "https://github.com/cheat-go/community",
			Downloads:   3210,
			Rating:      4.6,
			CreatedAt:   time.Now().Add(-45 * 24 * time.Hour),
			UpdatedAt:   time.Now().Add(-5 * 24 * time.Hour),
			Tags:        []string{"git", "vcs", "workflow"},
		},
	}
}

func (m *MockClient) GetRepositories() ([]Repository, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.repositories, nil
}

func (m *MockClient) SearchCheatSheets(opts SearchOptions) ([]CheatSheet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := []CheatSheet{}
	for _, sheet := range m.cheatSheets {
		if opts.Query != "" && !strings.Contains(strings.ToLower(sheet.Name), strings.ToLower(opts.Query)) {
			continue
		}
		if opts.MinRating > 0 && sheet.Rating < opts.MinRating {
			continue
		}
		if opts.Repository != "" && sheet.Repository != opts.Repository {
			continue
		}
		results = append(results, sheet)
	}

	if opts.Limit > 0 && len(results) > opts.Limit {
		results = results[:opts.Limit]
	}

	return results, nil
}

func (m *MockClient) GetCheatSheet(id string) (*CheatSheet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, sheet := range m.cheatSheets {
		if sheet.ID == id {
			return &sheet, nil
		}
	}
	return nil, fmt.Errorf("cheat sheet not found")
}

func (m *MockClient) DownloadCheatSheet(id string) (*apps.App, error) {
	sheet, err := m.GetCheatSheet(id)
	if err != nil {
		return nil, err
	}
	return &sheet.App, nil
}

func (m *MockClient) SubmitCheatSheet(sheet CheatSheet) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sheet.ID = fmt.Sprintf("custom-%d", time.Now().Unix())
	sheet.CreatedAt = time.Now()
	sheet.UpdatedAt = time.Now()
	m.cheatSheets = append(m.cheatSheets, sheet)
	return nil
}

func (m *MockClient) RateCheatSheet(id string, rating float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, sheet := range m.cheatSheets {
		if sheet.ID == id {
			m.cheatSheets[i].Rating = (sheet.Rating + rating) / 2
			return nil
		}
	}
	return fmt.Errorf("cheat sheet not found")
}
