# HK Pre-Commit Setup & GitHub Actions Workflow

This document describes the hk pre-commit setup and the GitHub Actions workflow for linting and formatting.

## Local Setup

### Install hk and tools

```bash
# Install hk via mise (or manually)
mise use -g hk@latest

# Install required linting/formatting tools via mise
mise install
```

### Run pre-commit checks locally

```bash
# Check all files
hk run pre-commit --check --all

# Fix all files (auto-format where possible)
hk run pre-commit --fix --all

# Check only changed files vs main
hk run pre-commit --check --from-ref origin/main --to-ref HEAD
```

## Configured Linters

The `hk.pkl` file configures these linters:

| Linter            | Purpose                       | Auto-Fix     |
| ----------------- | ----------------------------- | ------------ |
| **actionlint**    | GitHub Actions validation     | ❌           |
| **dclint**        | Docker Compose validation     | ✅           |
| **go_fmt**        | Go code formatting            | ✅           |
| **golangci_lint** | Go linting                    | ✅ (partial) |
| **prettier**      | YAML/MD/JSON formatting       | ✅           |
| **shellcheck**    | Shell script validation       | ❌           |
| **shellharden**   | Shell script hardening        | ✅           |
| **shfmt**         | Shell script formatting       | ✅           |
| **zizmor**        | GitHub Actions security audit | ✅ (partial) |

## GitHub Actions Workflow

**File**: `.github/workflows/hk-lint.yml`

### Behavior

- **Triggers**: On pull requests (opened, reopened, synchronized)
- **Scope**: Only checks files changed from base ref to PR head
- **Commands**: Uses `--from-ref origin/base --to-ref HEAD --check`
- **Feedback**: Comments on PR if checks fail

### Workflow Steps

1. Checkout code with full history
2. Set up Go environment
3. Install hk
4. Run `hk run pre-commit` on changed files
5. Comment on PR if checks fail

## Common Issues

### Template Injection Security

Use environment variables to avoid template injection:

```yaml
env:
  BASE_REF: ${{ github.base_ref }}
run: git fetch origin "$BASE_REF:refs/remotes/origin/$BASE_REF"
```

### shellharden Not Available

Ensure shfmt version is set:

```bash
mise use -g shfmt@3.14.1
```

### pkl Format Validation

Ensure pkl is installed:

```bash
mise use -g pkl@0.32.2
```

## Development Workflow

1. Make code changes
2. Run `hk run pre-commit --fix --all` locally
3. Commit and push
4. GitHub Actions workflow runs automatically
5. If checks fail, fix locally and push again

## Security Notes

- All GitHub Actions are pinned to commit hashes
- Workflow uses read-only permissions
- Template variables passed via environment variables
- No secrets are exposed in logs
