# Implementation Tasks

## 1. Configuration Schema
- [ ] 1.1 Define locality mapping configuration schema (JSON)
- [ ] 1.2 Define mapping rule structure (remote path pattern -> local path pattern)
- [ ] 1.3 Support variable substitution in path patterns
- [ ] 1.4 Document configuration format

## 2. Configuration Storage
- [ ] 2.1 Implement reading configuration from remote (`.filesprawl/locality.json`)
- [ ] 2.2 Implement reading local override configuration
- [ ] 2.3 Implement configuration merging (local overrides remote)
- [ ] 2.4 Add caching for configuration to avoid repeated remote reads

## 3. Locality Package
- [ ] 3.1 Create internal/locality package
- [ ] 3.2 Implement LocalityMapper type
- [ ] 3.3 Implement RemoteToLocal path resolution
- [ ] 3.4 Implement LocalToRemote path resolution
- [ ] 3.5 Support pattern matching and variable substitution
- [ ] 3.6 Write unit tests for path resolution

## 4. Integration
- [ ] 4.1 Add CLI commands to manage locality mappings
- [ ] 4.2 Add CLI commands to resolve paths (remote->local, local->remote)
- [ ] 4.3 Integrate with scanning to store locality information
- [ ] 4.4 Add integration tests
- [ ] 4.5 Update documentation


## 5. Remote Scoping and Regex Patterns
- [x] 5.1 Add `remoteName`/`remoteRegex` remote scoping to mapping rules (mutually exclusive)
- [x] 5.2 Add `isRegex` regex path pattern support with anchored, cached compilation
- [x] 5.3 Evaluate applicable rules in descending priority order
- [x] 5.4 Validate configuration on load with clear per-rule errors
- [x] 5.5 Add `--remote` flag to locality resolve commands and show scope/type in show-config
- [x] 5.6 Add unit and CLI tests for scoping, regex, priority ordering, and validation
- [x] 5.7 Document remote-scoped and regex rules in LOCALITY_CONFIG.md and README.md

