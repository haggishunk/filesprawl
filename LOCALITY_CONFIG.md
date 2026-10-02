# Locality Configuration

The filesprawl locality interface uses a JSON configuration file to map between remote storage paths and local filesystem paths.

## Configuration File Location

The system searches for configuration in this priority order:

1. **`FILESPRAWL_LOCALITY_CONFIG` environment variable** (if set)
   ```bash
   export FILESPRAWL_LOCALITY_CONFIG=/etc/filesprawl/locality.json
   filesprawl locality resolve-remote /documents/file.pdf
   ```

2. **XDG Base Directory Specification**
   - `$XDG_CONFIG_HOME/filesprawl/locality.json` (if `XDG_CONFIG_HOME` is set)
   - `$HOME/.config/filesprawl/locality.json` (standard XDG default)
   ```bash
   # Automatic: ~/.config/filesprawl/locality.json
   filesprawl locality resolve-remote /documents/file.pdf
   ```

3. **Legacy Location** (backward compatibility)
   - `$HOME/.filesprawl/locality.json`
   ```bash
   # Automatic: ~/.filesprawl/locality.json
   filesprawl locality resolve-remote /documents/file.pdf
   ```

## Configuration Format

```json
{
  "rules": [
    {
      "remotePattern": "/documents",
      "localPattern": "/home/user/Documents",
      "priority": 10
    },
    {
      "remotePattern": "/data/{category}/*",
      "localPattern": "/home/user/{category}/*",
      "priority": 20
    },
    {
      "remotePattern": "/*",
      "localPattern": "/mnt/backup/*",
      "priority": 1
    }
  ]
}
```

## Creating a Configuration File

```bash
# Create config directory (XDG standard location)
mkdir -p ~/.config/filesprawl

# Create configuration file
cat > ~/.config/filesprawl/locality.json << 'EOF'
{
  "rules": [
    {
      "remotePattern": "/documents",
      "localPattern": "/home/user/Documents",
      "priority": 1
    }
  ]
}
EOF

# Test it
filesprawl locality show-config
filesprawl locality resolve-remote /documents/report.pdf
```

## Rule Priority

Rules are evaluated in priority order (highest first):
- Higher numbers = higher priority
- More specific patterns should have higher priority
- Wildcard patterns typically have lower priority (catch-all)

## Pattern Syntax

**Wildcards:** `*` matches the rest of the path
- `remotePattern: "/data/*"` matches `/data/file1.txt`, `/data/subdir/file2.txt`, etc.
- `localPattern: "/home/user/data/*"` maps matched files to corresponding local paths

**Variables:** `{variable}` captures path segments
- `remotePattern: "/{category}/files"` captures `documents` in `/documents/files`
- `localPattern: "/home/user/{category}"` expands to `/home/user/documents`

## Environment Variables

### FILESPRAWL_LOCALITY_CONFIG
Override the default configuration file location:
```bash
export FILESPRAWL_LOCALITY_CONFIG=/custom/path/locality.json
filesprawl locality resolve-remote /documents/file.pdf
```

### XDG_CONFIG_HOME
Set custom XDG config directory:
```bash
export XDG_CONFIG_HOME=/etc/xdg
# Looks for: /etc/xdg/filesprawl/locality.json
```

## File Permissions

For security, restrict access to the configuration file:
```bash
chmod 600 ~/.config/filesprawl/locality.json
```

## Troubleshooting

### Configuration file not found
```bash
filesprawl locality show-config
# Shows current configuration path and rules
# If empty, configuration file doesn't exist or is invalid JSON
```

### Path not mapping
- Check rule priority (higher number = checked first)
- Verify pattern syntax (wildcards and variables)
- Use `show-config` to inspect loaded rules

### Invalid JSON
Error message will indicate parsing failure. Validate JSON:
```bash
python3 -m json.tool ~/.config/filesprawl/locality.json
```
