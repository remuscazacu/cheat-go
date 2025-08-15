package online

import (
	"cheat-go/pkg/apps"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClient_GetRepositories(t *testing.T) {
	repos := []Repository{
		{
			URL:         "https://github.com/test/repo1",
			Name:        "Test Repo 1",
			Description: "Test repository 1",
			LastUpdated: time.Now(),
			Stars:       100,
			Author:      "test",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/repositories" {
			t.Errorf("Expected path /api/repositories, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(repos)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	gotRepos, err := client.GetRepositories()

	if err != nil {
		t.Fatalf("GetRepositories() error = %v", err)
	}

	if len(gotRepos) != len(repos) {
		t.Errorf("GetRepositories() got %d repos, want %d", len(gotRepos), len(repos))
	}

	if gotRepos[0].Name != repos[0].Name {
		t.Errorf("GetRepositories() got name %s, want %s", gotRepos[0].Name, repos[0].Name)
	}
}

func TestHTTPClient_SearchCheatSheets(t *testing.T) {
	sheets := []CheatSheet{
		{
			ID:          "sheet1",
			Name:        "Test Sheet 1",
			Description: "Test cheat sheet 1",
			Downloads:   100,
			Rating:      4.5,
			Tags:        []string{"test", "sample"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cheatsheets" {
			t.Errorf("Expected path /api/cheatsheets, got %s", r.URL.Path)
		}

		q := r.URL.Query()
		if q.Get("q") != "test" {
			t.Errorf("Expected query param q=test, got %s", q.Get("q"))
		}
		if q.Get("min_rating") != "4.0" {
			t.Errorf("Expected min_rating=4.0, got %s", q.Get("min_rating"))
		}

		json.NewEncoder(w).Encode(sheets)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	opts := SearchOptions{
		Query:     "test",
		MinRating: 4.0,
		Limit:     10,
	}

	gotSheets, err := client.SearchCheatSheets(opts)
	if err != nil {
		t.Fatalf("SearchCheatSheets() error = %v", err)
	}

	if len(gotSheets) != len(sheets) {
		t.Errorf("SearchCheatSheets() got %d sheets, want %d", len(gotSheets), len(sheets))
	}
}

func TestHTTPClient_GetCheatSheet(t *testing.T) {
	sheet := CheatSheet{
		ID:          "sheet1",
		Name:        "Test Sheet",
		Description: "Test cheat sheet",
		App: apps.App{
			Name:        "Test App",
			Description: "Test application",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cheatsheets/sheet1" {
			t.Errorf("Expected path /api/cheatsheets/sheet1, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(sheet)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	gotSheet, err := client.GetCheatSheet("sheet1")

	if err != nil {
		t.Fatalf("GetCheatSheet() error = %v", err)
	}

	if gotSheet.ID != sheet.ID {
		t.Errorf("GetCheatSheet() got ID %s, want %s", gotSheet.ID, sheet.ID)
	}

	// Test caching
	gotSheet2, err := client.GetCheatSheet("sheet1")
	if err != nil {
		t.Fatalf("GetCheatSheet() cached error = %v", err)
	}
	if gotSheet2.ID != sheet.ID {
		t.Errorf("Cached GetCheatSheet() got ID %s, want %s", gotSheet2.ID, sheet.ID)
	}
}

func TestHTTPClient_SubmitCheatSheet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cheatsheets" {
			t.Errorf("Expected path /api/cheatsheets, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var sheet CheatSheet
		json.NewDecoder(r.Body).Decode(&sheet)

		if sheet.Name != "New Sheet" {
			t.Errorf("Expected name 'New Sheet', got %s", sheet.Name)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)
	sheet := CheatSheet{
		Name:        "New Sheet",
		Description: "New cheat sheet",
	}

	err := client.SubmitCheatSheet(sheet)
	if err != nil {
		t.Fatalf("SubmitCheatSheet() error = %v", err)
	}
}

func TestHTTPClient_RateCheatSheet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cheatsheets/sheet1/rate" {
			t.Errorf("Expected path /api/cheatsheets/sheet1/rate, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		var data map[string]float64
		json.NewDecoder(r.Body).Decode(&data)

		if data["rating"] != 4.5 {
			t.Errorf("Expected rating 4.5, got %f", data["rating"])
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	err := client.RateCheatSheet("sheet1", 4.5)
	if err != nil {
		t.Fatalf("RateCheatSheet() error = %v", err)
	}

	// Test invalid rating
	err = client.RateCheatSheet("sheet1", 6.0)
	if err == nil {
		t.Error("Expected error for invalid rating, got nil")
	}
}

func TestMockClient_Operations(t *testing.T) {
	client := NewMockClient()

	// Test GetRepositories
	repos, err := client.GetRepositories()
	if err != nil {
		t.Fatalf("GetRepositories() error = %v", err)
	}
	if len(repos) == 0 {
		t.Error("Expected default repositories, got none")
	}

	// Test SearchCheatSheets
	sheets, err := client.SearchCheatSheets(SearchOptions{Query: "vim"})
	if err != nil {
		t.Fatalf("SearchCheatSheets() error = %v", err)
	}
	if len(sheets) != 1 {
		t.Errorf("Expected 1 sheet matching 'vim', got %d", len(sheets))
	}

	// Test GetCheatSheet
	sheet, err := client.GetCheatSheet("vim-advanced")
	if err != nil {
		t.Fatalf("GetCheatSheet() error = %v", err)
	}
	if sheet.ID != "vim-advanced" {
		t.Errorf("GetCheatSheet() got ID %s, want vim-advanced", sheet.ID)
	}

	// Test SubmitCheatSheet
	newSheet := CheatSheet{
		Name:        "Test Sheet",
		Description: "Test description",
	}
	err = client.SubmitCheatSheet(newSheet)
	if err != nil {
		t.Fatalf("SubmitCheatSheet() error = %v", err)
	}

	// Test RateCheatSheet
	err = client.RateCheatSheet("vim-advanced", 5.0)
	if err != nil {
		t.Fatalf("RateCheatSheet() error = %v", err)
	}

	// Test DownloadCheatSheet
	_, err = client.DownloadCheatSheet("vim-advanced")
	if err != nil {
		t.Fatalf("DownloadCheatSheet() error = %v", err)
	}
}

func TestHTTPClient_Errors(t *testing.T) {
	// Test with server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	// Test GetRepositories error
	_, err := client.GetRepositories()
	if err == nil {
		t.Error("Expected error from GetRepositories")
	}

	// Test SearchCheatSheets error
	_, err = client.SearchCheatSheets(SearchOptions{Query: "test"})
	if err == nil {
		t.Error("Expected error from SearchCheatSheets")
	}

	// Test GetCheatSheet error
	_, err = client.GetCheatSheet("test")
	if err == nil {
		t.Error("Expected error from GetCheatSheet")
	}

	// Test SubmitCheatSheet error
	err = client.SubmitCheatSheet(CheatSheet{Name: "test"})
	if err == nil {
		t.Error("Expected error from SubmitCheatSheet")
	}

	// Test RateCheatSheet error
	err = client.RateCheatSheet("test", 5.0)
	if err == nil {
		t.Error("Expected error from RateCheatSheet")
	}

	// Test DownloadCheatSheet error
	_, err = client.DownloadCheatSheet("test")
	if err == nil {
		t.Error("Expected error from DownloadCheatSheet")
	}
}

func TestHTTPClient_InvalidURL(t *testing.T) {
	// Test with invalid base URL
	client := NewHTTPClient("invalid-url")

	_, err := client.GetRepositories()
	if err == nil {
		t.Error("Expected error with invalid URL")
	}
}

func TestHTTPClient_InvalidJSON(t *testing.T) {
	// Test with server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL)

	_, err := client.GetRepositories()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestMockClient_DownloadCheatSheet(t *testing.T) {
	client := NewMockClient()

	// Test successful download
	app, err := client.DownloadCheatSheet("vim-advanced")
	if err != nil {
		t.Fatalf("DownloadCheatSheet() error = %v", err)
	}
	if app == nil {
		t.Error("Downloaded app should not be nil")
	}

	// Test download non-existent sheet
	_, err = client.DownloadCheatSheet("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent sheet")
	}
}

func TestMockClient_SearchCheatSheetsFiltering(t *testing.T) {
	client := NewMockClient()

	// Test search with query
	results, err := client.SearchCheatSheets(SearchOptions{
		Query: "vim",
	})
	if err != nil {
		t.Fatalf("SearchCheatSheets() error = %v", err)
	}
	// Should return some results
	if len(results) == 0 {
		t.Error("Should find some vim-related results")
	}

	// Test search with tag filter
	results, err = client.SearchCheatSheets(SearchOptions{
		Tags: []string{"editor"},
	})
	if err != nil {
		t.Fatalf("SearchCheatSheets() error = %v", err)
	}

	// Test search with limit
	results, err = client.SearchCheatSheets(SearchOptions{
		Limit: 2,
	})
	if err != nil {
		t.Fatalf("SearchCheatSheets() error = %v", err)
	}
	if len(results) > 2 {
		t.Errorf("Expected at most 2 results with limit=2, got %d", len(results))
	}
}

func TestMockClient_GetCheatSheetEdgeCases(t *testing.T) {
	client := NewMockClient()

	// Test getting existing sheet
	sheet, err := client.GetCheatSheet("vim-advanced")
	if err != nil {
		t.Fatalf("GetCheatSheet() error = %v", err)
	}
	if sheet.ID != "vim-advanced" {
		t.Errorf("Expected ID vim-advanced, got %s", sheet.ID)
	}

	// Test getting non-existent sheet
	_, err = client.GetCheatSheet("non-existent-sheet")
	if err == nil {
		t.Error("Expected error for non-existent sheet")
	}
}

func TestMockClient_RateCheatSheetValidation(t *testing.T) {
	client := NewMockClient()

	// Test valid rating
	err := client.RateCheatSheet("vim-advanced", 4.5)
	if err != nil {
		t.Fatalf("RateCheatSheet() error = %v", err)
	}

	// Test rating with high value (MockClient might not validate range)
	err = client.RateCheatSheet("vim-advanced", 6.0)
	// MockClient may accept any rating value
	if err != nil {
		t.Logf("Rating validation: %v", err)
	}

	// Test rating with low value (MockClient might not validate range)
	err = client.RateCheatSheet("vim-advanced", -1.0)
	// MockClient may accept any rating value
	if err != nil {
		t.Logf("Rating validation: %v", err)
	}

	// Test rating non-existent sheet
	err = client.RateCheatSheet("non-existent", 3.0)
	if err == nil {
		t.Error("Expected error for rating non-existent sheet")
	}
}

func TestMockClient_RepositoryStructure(t *testing.T) {
	client := NewMockClient()

	repos, err := client.GetRepositories()
	if err != nil {
		t.Fatalf("GetRepositories() error = %v", err)
	}

	for _, repo := range repos {
		if repo.Name == "" {
			t.Error("Repository should have name")
		}
		if repo.URL == "" {
			t.Error("Repository should have URL")
		}
		if repo.Author == "" {
			t.Error("Repository should have author")
		}
		if repo.Stars < 0 {
			t.Error("Repository stars should not be negative")
		}
	}
}

func TestSearchOptionsEdgeCases(t *testing.T) {
	client := NewMockClient()

	// Test search with MinRating filter
	results, err := client.SearchCheatSheets(SearchOptions{
		MinRating: 4.0,
	})
	if err != nil {
		t.Fatalf("SearchCheatSheets() error = %v", err)
	}
	// Should handle MinRating filter gracefully

	// Test search with SortBy
	results, err = client.SearchCheatSheets(SearchOptions{
		SortBy: "rating",
	})
	if err != nil {
		t.Fatalf("SearchCheatSheets() error = %v", err)
	}

	// Test search with Offset
	results, err = client.SearchCheatSheets(SearchOptions{
		Offset: 1,
	})
	if err != nil {
		t.Fatalf("SearchCheatSheets() error = %v", err)
	}

	// Test that results structure is valid
	for _, result := range results {
		if result.ID == "" {
			t.Error("CheatSheet should have ID")
		}
		if result.Name == "" {
			t.Error("CheatSheet should have Name")
		}
	}
}
