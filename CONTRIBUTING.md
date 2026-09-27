# Contributing

## Commit messages

This repo uses [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

- **type** — one of:
  - `feat` — a new feature or capability
  - `fix` — a bug fix
  - `docs` — documentation only (requirements, architecture, README, comments)
  - `refactor` — code change that neither fixes a bug nor adds a feature
  - `perf` — a change that improves performance
  - `test` — adding or correcting tests
  - `build` — build system, dependencies, migrations
  - `ci` — CI/CD configuration
  - `chore` — everything else (tooling, formatting, repo housekeeping)
- **scope** *(optional)* — the area touched, e.g. `api`, `data-model`,
  `architecture`, `migrations`, `web`, `auth`, `notifications`.
- **subject** — imperative mood ("add", not "added"/"adds"), no trailing
  period, ideally ≤72 characters.
- **body** *(optional)* — explain *why*, not what (the diff already shows
  what); wrap around 72–100 characters per line.
- **footer** *(optional)* — `BREAKING CHANGE: ...` for breaking changes,
  `Refs: #123` / `Closes: #123` for issue references, and a
  `Co-Authored-By:` trailer for pair-authored or AI-assisted commits.

### Examples

```
feat(api): add POST /time-logs/start endpoint

Creates a time_logs row with start_at=now() and end_at=NULL. No check
against other running logs — overlapping timers are allowed by FR-4.4.
```

```
fix(recurrence): stop weekly rule from skipping DST-boundary occurrences

Materialization compared naive local times across the spring-forward
transition, dropping one occurrence per year. Convert via the user's
IANA timezone before comparing.

Refs: #42
```

```
docs(architecture): document notification delivery reliability trade-off
```

### Rules

- One logical change per commit; don't mix an unrelated refactor into a
  feature commit.
- Every commit should build/pass lint on its own once the codebase has a
  build (avoid "WIP" or "fix typo" follow-up commits — squash locally
  before pushing, or use `git commit --amend` on unpushed work).
- Reference the relevant `architecture/*.md` section in the body when a
  commit implements or deviates from a documented design decision.
