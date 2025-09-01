
# `cheat-go` Improvement Plan

This document outlines a detailed plan for improving the `cheat-go` project.

## Phase 1: Refactor the `main` package

*   **Goal:** Break the `main.go` file into smaller, more focused files.
*   **Tasks:**
    *   Create a new `ui` package to house all UI-related code.
    *   Move the `viewMode` enum and the `model` struct to a new `ui/model.go` file.
    *   Create separate files for each view in the `ui` package:
        *   `ui/view_main.go`
        *   `ui/view_notes.go`
        *   `ui/view_plugins.go`
        *   `ui/view_online.go`
        *   `ui/view_sync.go`
        *   `ui/view_help.go`
    *   Move the `Update` and `View` functions to the `ui` package and split them into smaller functions for each view.
    *   Move the `handle...Input` functions to their respective view files.
    *   Refactor `main.go` to be the application's entry point, responsible for initializing the model and starting the Bubble Tea program.

## Phase 2: Increase Test Coverage

*   **Goal:** Increase the test coverage for the `cache` and `sync` packages to over 80%.
*   **Tasks:**
    *   Write unit tests for the `cache` package, focusing on the LRU implementation.
    *   Write unit tests for the `sync` package, focusing on the sync logic and conflict resolution.
    *   Use a mocking framework to mock external dependencies, such as the file system and the network.

## Phase 3: Configuration and Error Handling

*   **Goal:** Move hardcoded values to the configuration file and improve error handling.
*   **Tasks:**
    *   Move the cache size and notes directory to the `config.yaml` file.
    *   Implement a consistent error handling strategy. Instead of printing errors to the console, return them to the caller and handle them appropriately.
    *   Use a logging library to log errors and other important information.

## Phase 4: Implement the `onlineClient`

*   **Goal:** Implement the `onlineClient` to fetch real data from online repositories.
*   **Tasks:**
    *   Use the `net/http` package to make HTTP requests to the online repositories.
    *   Use the `encoding/json` package to parse the JSON responses.
    *   Implement a caching mechanism to avoid making unnecessary HTTP requests.

## Phase 5: Documentation

*   **Goal:** Update the documentation to reflect the changes.
*   **Tasks:**
    *   Update the `README.md` file to reflect the new architecture.
    *   Add comments to the code to explain the changes.
    *   Create a new `ARCHITECTURE.md` file to document the new architecture.
