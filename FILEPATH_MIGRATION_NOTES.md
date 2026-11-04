# Filepath Storage & Display Schema

## New Storage Design (Fixed)

Local and rclone scans now use consistent, clean storage:

**LOCAL SCANS:**

- `location.name`: hostname (identifies machine)
- `remote.remote_name`: hostname (identifies source)
- `object_meta.object_path`: relative path (e.g., `photos/IMG_0587.MOV`)
- `object_meta.scan_root`: scan directory (filter for incremental scans)

**RCLONE SCANS:**

- `location.name`: remote name (identifies storage)
- `remote.remote_name`: remote name (identifies source)
- `object_meta.object_path`: path within remote (e.g., `photos/IMG_0587.MOV`)
- `object_meta.scan_root`: NULL (not applicable for remotes)

## UI Display & File Access

**LOCAL FILES:**

```
Remote display:   archlinux
Filepath display: /goblet/light/photos/IMG_0587.MOV
File access:      Direct filesystem path
```

**RCLONE FILES:**

```
Remote display:   goblet
Filepath display: photos/IMG_0587.MOV
File access:      goblet:photos/IMG_0587.MOV
```

## Data Migration

### For New Scans

Automatically uses correct schema:

```bash
filesprawl index --remote archlinux --path /goblet/light
```

### For Existing Data

Optional cleanup and re-scan:

```bash
# Clear old data
DELETE FROM object_meta WHERE scan_root = '/goblet/light';
DELETE FROM remote WHERE remote_name = '/goblet/light';

# Re-index
filesprawl index --remote archlinux --path /goblet/light
```

## Benefits

- **Direct mapping**: No adapter layer needed
- **Clean semantics**: Hostname = source identifier for local files
- **Type detection**: `scan_root != NULL` indicates local file
- **Flexible**: Can scan multiple roots from same machine
