# Database Migrations

The backend uses [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate) with SQL files embedded into the binary at compile time. Migrations are paired `up`/`down` SQL scripts in [`backend/internal/database/migrations/`](../backend/internal/database/migrations/), numbered sequentially (`001_create_users.up.sql` / `001_create_users.down.sql`, etc.).

The server applies all pending `up` migrations automatically at startup, so day-to-day you don't have to think about them. The CLI documented below is for rolling back, stepping through versions, or inspecting state.

## CLI: `cmd/migrate`

A small Go binary at [`backend/cmd/migrate/main.go`](../backend/cmd/migrate/main.go) drives the migrator. It reuses the same embedded migration source and `schema_migrations` bookkeeping table the server uses, so up and down stay perfectly consistent.

The database path is taken from `DATABASE_PATH` (or the `.env` file), defaulting to `./data/database/social.db`.

### Make targets

Run from the `backend/` directory:

| Command | Description |
|---|---|
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back **all** migrations (destructive — empties the schema) |
| `make migrate-down-one` | Roll back only the most recent migration |
| `make migrate-version` | Print current schema version and dirty flag |

### Direct CLI

The Make targets are thin wrappers around the CLI, which accepts a few extra forms:

```bash
go run ./cmd/migrate up         # all pending
go run ./cmd/migrate up 2       # apply the next 2 pending
go run ./cmd/migrate down       # roll back ALL (destructive)
go run ./cmd/migrate down 1     # roll back 1 migration
go run ./cmd/migrate version    # current version + dirty flag
```

## The `schema_migrations` table

`golang-migrate` maintains a single bookkeeping table:

| Column | Meaning |
|---|---|
| `version` | Number of the last attempted migration |
| `dirty` | `false` = the migration finished cleanly; `true` = it started but didn't complete |

A **dirty** state means a migration was interrupted (process killed mid-flight, an SQL statement errored partway through, etc.). The migrator refuses to apply further changes until the state is cleaned up — the database may be partially migrated and it has no way to know which statements ran.

To recover: inspect the schema by hand to figure out what was applied, then either finish the migration manually and update `schema_migrations` to clear the flag, or roll back the partial work and reset to the previous version. The CLI doesn't expose `force` — direct edits to `schema_migrations` via `sqlite3` are usually the simplest path:

```bash
sqlite3 ./data/database/forum.db "UPDATE schema_migrations SET dirty = 0;"
```

## Authoring a new migration

1. Pick the next sequential number and create both files in [`backend/internal/database/migrations/`](../backend/internal/database/migrations/):
   - `NNN_short_name.up.sql` — the change
   - `NNN_short_name.down.sql` — the exact reverse (drop indexes first, then tables; reverse `ALTER` statements)
2. The `//go:embed migrations/*.sql` directive in [`database.go`](../backend/internal/database/database.go) picks up new files automatically — no code changes required.
3. Test both directions before committing:
   ```bash
   make migrate-up
   make migrate-version
   make migrate-down-one
   make migrate-version
   make migrate-up
   ```
4. Keep migrations small and focused. One file per logical change.
