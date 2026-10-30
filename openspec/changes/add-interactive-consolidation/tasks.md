# Implementation Tasks: Interactive Duplicate Consolidation

## Phase 1: Foundation
- [ ] 1.1 Create spec for `consolidation-workflow` capability
- [ ] 1.2 Design database schema for consolidation history and preferences
- [ ] 1.3 Add database migrations for new tables
- [ ] 1.4 Create `internal/consolidation` package structure

## Phase 2: Core Consolidation Logic
- [ ] 2.1 Implement archive operations with path preservation (`internal/archive`)
- [ ] 2.2 Implement consolidation preferences storage and retrieval
- [ ] 2.3 Create consolidation decision engine (with memory of past choices)
- [ ] 2.4 Implement file move/copy operations via rclone RC API

## Phase 3: Preview & Discovery
- [ ] 3.1 Implement image thumbnail generation (`internal/previews`)
- [ ] 3.2 Add terminal-friendly image display (ANSI art or external viewer)
- [ ] 3.3 Create interactive duplicate review interface (TUI/CLI)

## Phase 4: Interactive Command
- [ ] 4.1 Implement `consolidate` command (`cmd/filesprawl/consolidate.go`)
- [ ] 4.2 Build interactive prompt workflow (discovery → preview → decision → execute)
- [ ] 4.3 Add consolidation progress tracking and rollback capability
- [ ] 4.4 Add flags for batch mode (auto-consolidate using preferences)

## Phase 5: Advanced Features (Optional)
- [ ] 5.1 Implement DigiKam database integration (`internal/digikam`)
- [ ] 5.2 Add smart consolidation (photo detection, rating-based selection)
- [ ] 5.3 Add consolidation report generation
- [ ] 5.4 Integration tests for end-to-end consolidation workflow
