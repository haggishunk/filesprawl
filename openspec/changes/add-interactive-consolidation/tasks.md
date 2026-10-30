# Implementation Tasks: Interactive Duplicate Consolidation

## Phase 1: Foundation
- [ ] 1.1 Create spec for `consolidation-workflow` capability
- [ ] 1.2 Design database schema for consolidation history and preferences
- [ ] 1.3 Add database migrations for new tables
- [ ] 1.4 Create `internal/consolidation` package structure

## Phase 2: Data Loss Prevention Foundation
- [ ] 2.1 Create `internal/safety` package for data loss detection
- [ ] 2.2 Implement name collision detector (same filename, different hashes)
- [ ] 2.3 Implement completeness verifier (0-byte vs non-zero detection)
- [ ] 2.4 Implement directory deletion safety checker (verify copies exist elsewhere)
- [ ] 2.5 Add consolidation safety analyzer (combine all three checks)

## Phase 3: Core Consolidation Logic
- [ ] 3.1 Implement archive operations with path preservation (`internal/archive`)
- [ ] 3.2 Implement consolidation preferences storage and retrieval
- [ ] 3.3 Create consolidation decision engine (with memory of past choices)
- [ ] 3.4 Implement file move/copy operations via rclone RC API

## Phase 4: Preview & Discovery
- [ ] 4.1 Implement image thumbnail generation (`internal/previews`)
- [ ] 4.2 Add terminal-friendly image display (ANSI art or external viewer)
- [ ] 4.3 Create interactive duplicate review interface (TUI/CLI)

## Phase 5: Interactive Command with Safety Display
- [ ] 5.1 Implement `consolidate` command (`cmd/filesprawl/consolidate.go`)
- [ ] 5.2 Integrate safety checks into consolidation analysis phase
- [ ] 5.3 Build safety panel display (consolidated view of all three safety checks)
- [ ] 5.4 Build interactive prompt workflow (discovery → safety review → preview → decision → execute)
- [ ] 5.5 Add consolidation progress tracking and rollback capability
- [ ] 5.6 Add flags for batch mode (auto-consolidate using preferences)

## Phase 6: Advanced Features (Optional)
- [ ] 6.1 Implement DigiKam database integration (`internal/digikam`)
- [ ] 6.2 Add smart consolidation (photo detection, rating-based selection)
- [ ] 6.3 Add consolidation report generation
- [ ] 6.4 Integration tests for end-to-end consolidation workflow
