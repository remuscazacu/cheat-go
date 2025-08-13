package apps

// App represents a single application with its shortcuts
type App struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description"`
	Categories  []string          `yaml:"categories" json:"categories"`
	Shortcuts   []Shortcut        `yaml:"shortcuts" json:"shortcuts"`
	Metadata    map[string]string `yaml:"metadata" json:"metadata"`
	Version     string            `yaml:"version" json:"version"`
}

// Shortcut represents a single keyboard shortcut
type Shortcut struct {
	Keys        string   `yaml:"keys" json:"keys"`
	Description string   `yaml:"description" json:"description"`
	Category    string   `yaml:"category" json:"category"`
	Tags        []string `yaml:"tags" json:"tags"`
	Platform    string   `yaml:"platform,omitempty" json:"platform,omitempty"`
}

// AppRegistry holds all registered applications
type AppRegistry struct {
	apps map[string]*App
}

// NewAppRegistry creates a new app registry
func NewAppRegistry() *AppRegistry {
	return &AppRegistry{
		apps: make(map[string]*App),
	}
}

// Register adds an app to the registry
func (r *AppRegistry) Register(app *App) {
	r.apps[app.Name] = app
}

// Get retrieves an app by name
func (r *AppRegistry) Get(name string) (*App, bool) {
	app, exists := r.apps[name]
	return app, exists
}

// GetAll returns all registered apps
func (r *AppRegistry) GetAll() map[string]*App {
	return r.apps
}

// List returns the names of all registered apps
func (r *AppRegistry) List() []string {
	names := make([]string, 0, len(r.apps))
	for name := range r.apps {
		names = append(names, name)
	}
	return names
}
