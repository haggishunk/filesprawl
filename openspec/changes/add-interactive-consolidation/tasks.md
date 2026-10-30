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

## Phase 2b: Persistence and Multi-User Foundation
- [ ] 2b.1 Create database schema: consolidation_users table (id, username, created_at, last_active)
- [ ] 2b.2 Create database schema: consolidation_flags table (file_id, user_id, flag_type, flagged_at, reason, archive_dest)
- [ ] 2b.3 Create database schema: consolidation_reviews table (file_id, original_flag, original_user, reviewed_by, reviewed_at, new_flag, review_reason)
- [ ] 2b.4 Add database migrations for three new tables
- [ ] 2b.5 Create `internal/consolidation/users` package (user management, authentication)
- [ ] 2b.6 Create `internal/consolidation/flags` package (flag persistence, retrieval, modification)

## Phase 3: Core Consolidation Logic
- [ ] 3.1 Implement archive operations with path preservation (`internal/archive`)
- [ ] 3.2 Implement consolidation preferences storage and retrieval
- [ ] 3.3 Create consolidation decision engine (with memory of past choices)
- [ ] 3.4 Implement file move/copy operations via rclone RC API

## Phase 4: Dual-View Browsing Interface with Persistent Flagging
- [ ] 4.1 Create `internal/consolidation/browsing` package
- [ ] 4.2 Implement hash-centric view renderer (group by content hash)
- [ ] 4.3 Implement filename-centric view renderer (group by filename)
- [ ] 4.4 Implement location presence indicator system (🟢/🟡/🔴 color coding)
- [ ] 4.5 Implement file flagging state machine (KEEP/ARCHIVE/DELETE/SKIP) with immediate DB persistence
- [ ] 4.6 Implement flag persistence layer (write to consolidation_flags table on each flag change)
- [ ] 4.7 Implement view switching (hash ↔ filename)
- [ ] 4.8 Implement decision summary display (grouped by action type)
- [ ] 4.9 Implement quit-and-save functionality (exit without executing consolidation, flags preserved)

## Phase 4b: Preview & Discovery
- [ ] 4b.1 Implement image thumbnail generation (`internal/previews`)
- [ ] 4b.2 Add terminal-friendly image display (ANSI art or external viewer)
- [ ] 4b.3 Integrate previews with browsing interface (overlay on files)

## Phase 5: Interactive Command with Safety, Browsing, and Multi-User
- [ ] 5.1 Implement `consolidate` command (`cmd/filesprawl/consolidate.go`)
- [ ] 5.2 Implement user authentication prompt (username selection/creation at session start)
- [ ] 5.3 Implement "load pending flags" flow (show existing flags from DB, option to accept/modify)
- [ ] 5.4 Integrate safety checks into consolidation analysis phase
- [ ] 5.5 Build safety panel display (consolidated view of all three safety checks)
- [ ] 5.6 Integrate dual-view browsing interface into command workflow
- [ ] 5.7 Build interactive prompt workflow (auth → pending flags → safety → browse by hash/filename → flag decisions → summary → execute)
- [ ] 5.8 Implement keyboard shortcut system (k/a/d for keep/archive/delete, TAB for view switching, Q for quit-and-save)
- [ ] 5.9 Implement multi-user review display (show original flagging user + previous modifications in audit trail)
- [ ] 5.10 Add consolidation progress tracking and rollback capability
- [ ] 5.11 Add flags for batch mode (auto-consolidate using preferences)
- [ ] 5.12 Log all user actions and decisions to consolidation_reviews table

## Phase 6: Advanced Features (Optional)
- [ ] 6.1 Implement DigiKam database integration (`internal/digikam`)
- [ ] 6.2 Add smart consolidation (photo detection, rating-based selection)
- [ ] 6.3 Add consolidation report generation
- [ ] 6.4 Add search/filter in browsing interface (find by hash substring, filename pattern)
- [ ] 6.5 Integration tests for end-to-end consolidation workflow
