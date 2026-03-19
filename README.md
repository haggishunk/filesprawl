# filesprawl

## indexing

first job is to index all object and remote drive storage files

filesprawl assumes it is talking to a **local** rclone `rcd` session on `http://localhost:5572`, and that the remote names you pass on the CLI are the locally configured rclone remote names visible to that `rcd` session.

example workflow:

- `filesprawl list-remotes`
- `filesprawl index --remote media`
- `filesprawl index --remote media --path code/flux`

`list-remotes` prints the local rclone remote names as configured. `index --remote` accepts either `media` or `media:` and normalizes that name before persisting it.

`index` also requires `DATABASE_URL` to be set so scan results can be persisted with the local host name and the selected local rclone remote name.

with this index we can determine the following:

- duplicates between remote storage locations
- duplicates in same remote storage locations


## locality interface

rclone is great for configuring remote storage locations but we need a way to configure a mapping to local storage locations.  this interface configuration could be stored in a known location in the remote with optional local overrides for user (system is out of scope).

## cost savings

how can we optimize for cost savings?

not all remote storages cost the same (eg. storage size, transfer).

## feature and integration support

how can we optimize for features and integration support?

some content is meant to be consumed by tools (eg. terraform s3 backend) and integrations (eg. dropbox user sharing) that are better supported by certain remote backends.

this can be refined by cost savings.
