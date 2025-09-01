package ui

func (m Model) ViewHelp() string {
	return `
╭─ Help ───────────────────────────────────────────────────╮
│                                                          │
│  NAVIGATION                                              │
│    ↑/k, ↓/j, ←/h, →/l  Navigate                        │
│    Home/Ctrl+A          Go to first row                 │
│    End/Ctrl+E           Go to last row                  │
│                                                          │
│  FEATURES                                                │
│    /                    Search mode                     │
│    f                    Filter apps                     │
│    n                    Notes manager                   │
│    p                    Plugin manager                  │
│    o                    Browse online                   │
│    s                    Sync status                     │
│    Ctrl+S               Force sync                      │
│    Ctrl+R               Refresh data                    │
│    ?                    This help screen                │
│    q/Ctrl+C             Quit                           │
│                                                          │
│  SEARCH MODE                                            │
│    Type to search                                       │
│    Enter                Confirm search                  │
│    Esc                  Cancel search                   │
│    Ctrl+U               Clear search                    │
│                                                          │
╰──────────────────────────────────────────────────────────╯

Press ? or Esc to close help
`
}
