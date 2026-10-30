## MODIFIED Requirements

### Requirement: Indexed Remote Identity
The system SHALL persist remote identity using source location (hostname for local scans, remote name for rclone remotes) together with scan root and remote name for disambiguation.

#### Scenario: Local scan identification
- **WHEN** Filesprawl scans local path `/media` on host `archlinux`
- **THEN** persisted identity is: location(archlinux, local) + scan_root(/media)
- **AND** different local roots create separate remote records

#### Scenario: Remote source identification
- **WHEN** Filesprawl indexes rclone remote `dropbox:media`
- **THEN** persisted identity is: location(dropbox, remote) + remote_name(media)
- **AND** hostname is irrelevant (location type disambiguates)

#### Scenario: Distinguish same remote name across locations
- **WHEN** two different hosts or remotes each have a `media` index
- **THEN** they remain distinct by location_id
- **AND** no false duplicates occur between them

## REMOVED Requirements

### Requirement: Hostname Column in Remote
**Reason**: Hostname is only meaningful for local scans; replaced by location table with explicit type.
**Migration**: All hostname references moved to location(name='hostname', type='local').
