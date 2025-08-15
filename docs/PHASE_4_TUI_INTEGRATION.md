# Phase 4 TUI Integration Plan

## Overview
Integrate all Phase 4 features (plugins, notes, online, sync, cache) into the main TUI interface with keyboard shortcuts and interactive views.

## Implementation Plan

### 1. Update Main Model Structure
- Add new fields to the main model for Phase 4 features
- Add view states for different modes (notes, plugins, online, sync)
- Integrate cache for performance

### 2. New Keyboard Shortcuts
- `n` - Open notes manager view
- `p` - Open plugin manager view  
- `o` - Browse online repositories
- `s` - Show sync status
- `Ctrl+S` - Force sync

### 3. View Components

#### Notes View (`n` key)
- List personal notes
- **Create/Edit/Delete notes** (✅ Editor integration fixed)
  - `e` key opens external editor ($EDITOR or nano)
  - Structured editing with title, category, tags, content
  - Support for vim, nano, emacs, code, and other editors
- Search notes
- Link notes to apps
- Export/Import notes

#### Plugin Manager View (`p` key)
- List installed plugins
- Load/Unload plugins
- Show plugin info
- Configure plugins
- Install from directory

#### Online Repository View (`o` key)
- Browse repositories
- Search cheat sheets
- Download sheets
- Rate sheets
- Submit sheets

#### Sync Status View (`s` key)
- Show sync status
- Display conflicts
- Resolve conflicts
- Configure sync settings
- Manual sync trigger

### 4. Integration Points

#### Cache Integration
- Cache app data for faster loading
- Cache search results
- Cache online repository data

#### Background Services
- Auto-sync timer
- Plugin auto-loading
- Cache cleanup

### 5. Configuration Updates
Add new config sections:
```yaml
plugins:
  enabled: true
  auto_load: true
  directories:
    - ~/.config/cheat-go/plugins

notes:
  enabled: true
  directory: ~/.config/cheat-go/notes
  
online:
  enabled: true
  repositories:
    - https://github.com/cheat-go/community
  cache_ttl: 1h
  
sync:
  enabled: false
  auto_sync: true
  interval: 15m
  endpoint: https://sync.cheatsheets.com
  
cache:
  enabled: true
  memory_size: 10485760
  disk_cache: ~/.cache/cheat-go
```

## Implementation Steps

1. **Update Model** - Add Phase 4 feature managers to main model
2. **Add View States** - Create new view states enum
3. **Implement Views** - Create view components for each feature
4. **Add Key Handlers** - Implement keyboard shortcuts
5. **Update Config** - Add configuration loading for Phase 4
6. **Add Tests** - Test all new functionality
7. **Update Documentation** - Update README with new features