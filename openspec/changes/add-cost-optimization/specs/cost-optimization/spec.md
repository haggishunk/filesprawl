## ADDED Requirements

### Requirement: Cost Model Configuration
The system SHALL support configurable cost parameters for each remote storage backend.

#### Scenario: Storage cost configuration
- **WHEN** configuring a remote
- **THEN** storage cost per GB per month can be specified

#### Scenario: Transfer cost configuration
- **WHEN** configuring a remote
- **THEN** transfer cost per GB (ingress and egress) can be specified

#### Scenario: Minimum storage duration
- **WHEN** configuring a remote
- **THEN** minimum storage duration (e.g., 30 days for glacier) can be specified

### Requirement: Current Cost Analysis
The system SHALL calculate current storage costs across all remotes.

#### Scenario: Per-remote cost calculation
- **WHEN** analyzing costs
- **THEN** calculate total storage cost for each remote based on file sizes and cost parameters

#### Scenario: Total cost calculation
- **WHEN** analyzing costs
- **THEN** calculate total storage cost across all remotes

#### Scenario: Cost breakdown by file type
- **WHEN** analyzing costs
- **THEN** provide cost breakdown by file type or directory

### Requirement: Duplicate Cost Analysis
The system SHALL identify cost savings from duplicate file elimination.

#### Scenario: Duplicate storage cost
- **WHEN** analyzing duplicates
- **THEN** calculate total cost of storing duplicate files

#### Scenario: Potential savings calculation
- **WHEN** duplicates are identified
- **THEN** calculate potential savings from removing duplicates from higher-cost remotes

#### Scenario: Redundancy preservation
- **WHEN** calculating savings
- **THEN** ensure at least one copy is preserved on a remote

### Requirement: Cost Optimization Recommendations
The system SHALL generate recommendations for cost-optimal file placement.

#### Scenario: File migration recommendations
- **WHEN** generating optimization recommendations
- **THEN** suggest moving files from high-cost to low-cost remotes

#### Scenario: Duplicate elimination recommendations
- **WHEN** generating optimization recommendations
- **THEN** suggest which duplicate copies to remove based on cost

#### Scenario: Projected savings
- **WHEN** recommendations are generated
- **THEN** calculate projected cost savings from implementing recommendations

### Requirement: Cost-Aware Migration Planning
The system SHALL support planning file migrations with cost considerations.

#### Scenario: Migration cost calculation
- **WHEN** planning a file migration
- **THEN** calculate transfer costs for moving files between remotes

#### Scenario: Break-even analysis
- **WHEN** planning a migration
- **THEN** calculate time to break even on transfer costs vs storage savings

#### Scenario: Batch migration planning
- **WHEN** planning multiple migrations
- **THEN** optimize batch order to minimize total cost

