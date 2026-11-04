# Implementation Tasks: Consolidation UI Redesign

## Phase 1: Component Architecture (Bubble Tea foundation ready)

- [x] Add Bubble Tea and Lipgloss dependencies
- [x] Create messages.go with typed message system
- [x] Create bubbletea_model.go with ConsolidationModel
- [x] Update app.go to use Bubble Tea instead of gotui
- [ ] Refactor Model to include panel state (selected index, widths)

## Phase 2: Hash List Panel

- [ ] Create ListPane component in components/list_pane.go
- [ ] Implement hash group rendering with indicators + counts
- [ ] Wire j/k keyboard navigation for list scrolling
- [ ] Implement selection highlighting
- [ ] Update Model.View() to render left panel

## Phase 3: File Table Panel

- [ ] Create TablePane component in components/table_pane.go
- [ ] Implement table structure: remote | filename | size | status
- [ ] Add status icon rendering (✓/🟡/🔴)
- [ ] Add flag status display (KEEP/ARCHIVE/DELETE markers)
- [ ] Implement "wide" toggle (filepath column on/off)
- [ ] Add j/k navigation for file rows

## Phase 4: Banner & Status Bar

- [ ] Create Banner component showing version, group count, mode, remotes
- [ ] Create StatusBar component showing available actions
- [ ] Add flag summary to status bar
- [ ] Add deletion safety warning to status

## Phase 5: Keyboard Event Handling

- [ ] Wire j/k scroll events to both panels
- [ ] Implement v (view toggle) for table columns
- [ ] Implement K/a/d/s (action) key handlers
- [ ] Add ? (help) modal trigger
- [ ] Add q (quit) confirmation

## Phase 6: Deletion Safety Logic

- [ ] Create safety checker function in bubbletea_model.go
- [ ] Check (remote, filename) pairs before DELETE flag
- [ ] Show error message if deletion unsafe
- [ ] Block UI action on safety violation

## Phase 7: Integration & Polish

- [ ] Connect panel state to model updates
- [ ] Handle terminal resize (WindowSizeMsg)
- [ ] Test keyboard navigation across panels
- [ ] Add help modal
- [ ] Test with real data from database

## Phase 8: Testing

- [ ] Unit tests for safety logic
- [ ] Integration test: load duplicates, navigate, flag, check safety
- [ ] Manual test with various terminal sizes
