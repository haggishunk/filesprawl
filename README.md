# filesprawl

## indexing

filesprawl indexes files from remote storage locations and local filesystems to identify duplicates and optimize storage usage.

### Remote Storage Indexing

filesprawl assumes it is talking to a **local** rclone `rcd` session on `http://localhost:5572`, and that the remote names you pass on the CLI are the locally configured rclone remote names visible to that `rcd` session.

example workflow for remote indexing:

- `filesprawl list-remotes` — List available rclone remotes
- `filesprawl index --remote media` — Index entire remote
- `filesprawl index --remote media --path code/flux` — Index specific path within remote

`list-remotes` prints the local rclone remote names as configured. `index --remote` accepts either `media` or `media:` and normalizes that name before persisting it.

All indexing commands require `DATABASE_URL` to be set so scan results can be persisted with the local host name and the selected remote name.

### Local File Indexing

You can also index local filesystem paths:

```bash
filesprawl scan-local /path/to/local/directory
```

Local files are indexed using direct filesystem operations (not through rclone) and participate in the same duplicate detection as remote files.

## finding duplicates

once indexed, you can find duplicate files:

### within a specific remote

```bash
filesprawl duplicates --remote media
```

### across all remotes

```bash
filesprawl duplicates --across-remotes
```

### with filters

```bash
# find duplicates larger than 10MB
filesprawl duplicates --remote media --min-size 10485760

# find duplicates using specific hash type
filesprawl duplicates --across-remotes --hash-type sha256

# limit results and use pagination
filesprawl duplicates --across-remotes --limit 20 --offset 0

# combine filters
filesprawl duplicates --remote media --hash-type dropbox --min-size 1000000 --limit 10
```

available filters:
- `--hash-type`: filter by hash type (md5, sha1, sha256, dropbox)
- `--min-size`: minimum file size in bytes
- `--limit`: maximum number of duplicate groups to return
- `--offset`: skip this many duplicate groups (for pagination)

with this index we can determine the following:

- duplicates between remote storage locations
- duplicates in same remote storage locations


## locality interface

The locality interface maps remote storage paths to local filesystem paths, enabling you to understand where remote files correspond to on your local system.

### Path Resolution

Resolve remote paths to local paths:

```bash
filesprawl locality resolve-remote /documents/report.pdf
# Output: Remote: /documents/report.pdf
#         Local:  /home/user/Documents/report.pdf
```

Resolve local paths to remote paths:

```bash
filesprawl locality resolve-local /home/user/Documents/report.pdf
# Output: Local: /home/user/Documents/report.pdf
#         Remote paths:
#           - /documents/report.pdf
```

### Configuration

Locality configuration uses JSON and supports pattern-based mappings:

```json
{
  "rules": [
    {
      "remotePattern": "/documents",
      "localPattern": "/home/user/Documents",
      "priority": 1
    },
    {
      "remotePattern": "/data/{category}/*",
      "localPattern": "/home/user/{category}/*",
      "priority": 2
    }
  ]
}
```

View your current configuration:

```bash
filesprawl locality show-config
```

Configuration is loaded from (in priority order):
1. `FILESPRAWL_LOCALITY_CONFIG` environment variable
2. `$HOME/.config/filesprawl/locality.json` (XDG Base Directory)
3. `$HOME/.filesprawl/locality.json` (legacy location)

See [LOCALITY_CONFIG.md](LOCALITY_CONFIG.md) for complete configuration documentation.

## cost savings

how can we optimize for cost savings?

not all remote storages cost the same (eg. storage size, transfer).

## feature and integration support

how can we optimize for features and integration support?

some content is meant to be consumed by tools (eg. terraform s3 backend) and integrations (eg. dropbox user sharing) that are better supported by certain remote backends.

this can be refined by cost savings.
