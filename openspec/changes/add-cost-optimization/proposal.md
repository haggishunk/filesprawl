# Change: Add Cost Optimization

## Why
Different remote storage backends have different cost structures for storage size and data transfer. Users need to analyze current storage costs and receive recommendations for optimizing file placement across remotes to minimize total cost while maintaining required redundancy and access patterns.

## What Changes
- Add cost modeling for different remote storage backends
- Support configurable cost parameters (storage cost per GB, transfer cost per GB)
- Analyze current storage costs across all remotes
- Provide recommendations for cost-optimal file placement
- Consider duplicate files when calculating potential savings
- Support cost-aware migration planning

## Impact
- Affected specs: New capability `cost-optimization`
- Affected code:
  - New package `internal/cost` for cost modeling and optimization
  - Extend remote configuration to include cost parameters
  - Database schema updates to store cost information
  - Analysis tools to calculate current and projected costs

