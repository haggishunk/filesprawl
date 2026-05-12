# Project Context

## Purpose
Filesprawl is a file indexing and deduplication system for managing files across multiple remote storage locations. The project aims to:
- Index all objects across remote drive storage locations (Dropbox, S3, etc.)
- Identify duplicates between and within remote storage locations
- Optimize storage costs by analyzing file placement across different storage backends
- Support feature-based optimization (e.g., Terraform S3 backends, Dropbox sharing)
- Provide a locality interface mapping remote storage to local paths

## Tech Stack
- **Language**: Go 1.21.3
- **Database**: PostgreSQL (via Docker Compose)
- **Remote Storage Integration**: rclone (v1.65.0) with RC (remote control) API
- **Database Driver**: pgx/v5 (PostgreSQL driver and toolkit)
- **Container Orchestration**: Docker Compose
- **Database Admin**: pgAdmin4 (development)
- **Additional Libraries**:
  - mitchellh/mapstructure (v1.5.0) - struct mapping
  - rclone/rclone (v1.65.0) - remote storage operations

## Project Conventions

### Code Style
- **Package Organization**: Internal packages under `internal/` following Go best practices
- **Naming Conventions**:
  - Packages: lowercase, single-word names (e.g., `object`, `remote`, `operation`)
  - Interfaces: Descriptive names (e.g., `Database`)
  - Structs: PascalCase (e.g., `ObjectRepository`, `PgxDatabase`)
  - Factory functions: `New` prefix (e.g., `NewObjectRepository`, `NewHash`)
  - Functional options: `With` prefix (e.g., `WithRemote`, `WithRepo`)
- **Error Handling**: Wrap errors with context using `fmt.Errorf` with `%w` verb
- **Logging**: Standard library `log` package with descriptive messages
- **Comments**: Document exported types, functions, and complex logic

### Architecture Patterns
- **Repository Pattern**: `repository.ObjectRepository` abstracts database operations
- **Database Abstraction**: `database.Database` interface allows for testability and flexibility
- **Functional Options Pattern**: Used for constructors (e.g., `Scanner`, `Hash`, `Meta`)
- **Separation of Concerns**:
  - `internal/database`: Database connection and interface
  - `internal/object`: Domain models (Hash, Meta, MetaHashJunction)
  - `internal/operation`: Business logic (scanning operations)
  - `internal/rclone`: External API integration (rclone RC client)
  - `internal/remote`: Remote storage configuration
  - `internal/repository`: Data persistence layer
  - `internal/analysis`: Analysis logic (duplicate detection, report formatting)
- **Context Propagation**: All database and external operations accept `context.Context`
- **Connection Pooling**: pgxpool for efficient database connection management

### Testing Strategy
- **Test Files**: Co-located with source files using `_test.go` suffix
- **Test Framework**: Go standard library `testing` package
- **Test Coverage**: Unit tests for individual packages (currently minimal/scaffolded)
- **Database Testing**: Uses `database.Database` interface for mocking in unit tests
- **Integration Testing**: Docker Compose provides PostgreSQL for integration tests;
  integration tests are skipped when `DATABASE_URL` is not set

### Git Workflow
- **Repository**: github.com/haggishunk/filesprawl
- **Branching**: (To be defined based on team preferences)
- **Commits**: (To be defined based on team preferences)

## Domain Context

### File Indexing Concepts
- **Object Hash**: Cryptographic hash of file content (MD5, SHA1, SHA256, Dropbox-specific)
- **Object Meta**: File metadata (name, path, MIME type, size in bytes)
- **Junction Table**: Links metadata to hashes with scan timestamps for temporal tracking
- **Remote**: Storage backend configuration (hostname, remote name)
- **Scanning**: Recursive directory traversal that persists file information
- **Duplicate Detection**: Identifies files with identical content hashes within or across
  remotes, grouped by hash value with optional filtering by hash type and minimum file size

### Storage Optimization Goals
1. **Deduplication**: Identify identical files across remotes using content hashes
2. **Cost Optimization**: Different remotes have different storage/transfer costs
3. **Feature Support**: Some backends better support specific use cases (S3 for Terraform, Dropbox for sharing)
4. **Locality Mapping**: Map remote paths to local filesystem locations

### rclone Integration
- Uses rclone's RC (remote control) server running on localhost:5572
- Communicates via HTTP POST with JSON payloads
- Supports operations like `operations/list` for directory listing
- Handles various remote types (Dropbox, S3, etc.) through rclone's unified interface

## Important Constraints

### Technical Constraints
- **rclone Dependency**: Requires rclone RC server running (`rclone rcd --rc-serve --rc-no-auth`)
- **Database**: PostgreSQL required (currently no support for other databases)
- **Go Version**: Requires Go 1.21.3 or compatible
- **Environment Variables**: `DATABASE_URL` must be set for database connection

### Development Constraints
- **Docker Compose**: Development database runs in Docker
- **Database Credentials**: Test credentials in docker-compose.yaml (not for production)
- **Port Conflicts**: PostgreSQL (5432), pgAdmin (8080), rclone RC (5572)

### Data Model Constraints
- **Hash Types**: Limited to enum values (md5, dropbox, sha1, sha256)
- **Temporal Tracking**: Scan timestamps track when file-hash relationships were observed
- **No Deletion Tracking**: Current schema doesn't track file deletions

## External Dependencies

### rclone RC Server
- **Purpose**: Unified interface for remote storage operations
- **Endpoint**: http://localhost:5572
- **Authentication**: Currently disabled (--rc-no-auth)
- **Operations Used**: `operations/list` for directory/file listing
- **Configuration**: Uses rclone's config file for remote definitions

### PostgreSQL Database
- **Version**: Latest (via Docker)
- **Database Name**: filesprawl
- **User**: sprawler
- **Schema**: Defined in `db/init/schema.sql`
- **Tables**:
  - `object_hash`: Content hashes with type; indexed on `hash_value` for duplicate queries
  - `object_meta`: File metadata (name, path, MIME type, size in bytes)
  - `object_hash_junction`: Many-to-many relationship with scan timestamps
  - `object_remote_junction`: Links files to remotes
  - `remote`: Remote storage configuration

### pgAdmin
- **Purpose**: Database administration during development
- **Access**: http://localhost:8080
- **Credentials**: admin@example.com / adminpassword (development only)

### Remote Storage Providers
- Configured through rclone (Dropbox, S3, Google Drive, etc.)
- Each provider may have different capabilities and cost structures
- Accessed uniformly through rclone's abstraction layer
