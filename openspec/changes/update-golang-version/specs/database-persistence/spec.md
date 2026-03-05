## MODIFIED Requirements

### Requirement: Database Interface Abstraction
The system SHALL provide a Database interface for testability and flexibility.

#### Scenario: Interface implementation
- **WHEN** database operations are performed
- **THEN** they use the Database interface (QueryRow, Query, and Exec methods)

#### Scenario: PostgreSQL implementation
- **WHEN** using PostgreSQL
- **THEN** PgxDatabase struct implements the Database interface

#### Scenario: Multi-row query support
- **WHEN** a query returns multiple rows
- **THEN** the Query method returns pgx.Rows for iteration

