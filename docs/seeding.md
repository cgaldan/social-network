# Database Seeding

The `cmd/seed` tool populates the SQLite database with realistic data for manual testing and local development. It runs in two layers:

1. A **curated dataset** — a small, hand-authored cast of users, posts, follows, groups, events, direct messages and notifications. Lives in [`backend/cmd/seed/data/seed.json`](../backend/cmd/seed/data/seed.json).
2. An optional **procedural bulk layer** — thousands of synthetic users with realistic activity distribution (lurkers, normals, actives, celebrities) and matching posts, comments, follows and conversations.

Both layers can run together. The procedural layer always builds on top of the curated one so that named users and procedural users coexist.

---

## Quick start

From the `backend/` directory:

```bash
make seed                     # apply the curated dataset on top of existing data
make seed-reset               # wipe the DB first, then apply the curated dataset
make seed-scale               # wipe, apply curated, then add 1000 procedural users
make seed-scale SCALE=5000    # same but with 5000 procedural users
```

Each target invokes `cmd/seed` with the appropriate flags. You can also run it directly:

```bash
go run ./cmd/seed                            # curated only, no reset
go run ./cmd/seed -reset                     # wipe + curated
go run ./cmd/seed -reset -scale 1000         # wipe + curated + 1000 procedural users
go run ./cmd/seed -file path/to/other.json   # use a different curated dataset
```

> Migrations are applied automatically at the start of every seed run, so you don't need to run `make migrate` separately.

---

## What the curated dataset includes

The curated dataset in [`seed.json`](../backend/cmd/seed/data/seed.json) references users by **nickname** rather than by ID, so you can edit it without thinking about auto-increment values. Each top-level key maps to a section of the seed:

| Section | What it produces |
|---|---|
| `users` | Accounts with profile fields, avatars, public/private flag |
| `follows` | Follow edges with `accepted` or `pending` status |
| `posts` | Posts with category, privacy level, audience list, comments |
| `direct_chats` | Two-person conversations with a thread of messages |
| `groups` | Groups with members, invitations, join requests, posts, events, RSVPs |
| `notifications` | Pre-baked notifications targeting curated users |

Edit the JSON freely; the seeder validates that every referenced nickname exists.

---

## Procedural bulk layer

When `-scale N` is provided, the seeder generates `N` extra users distributed across four behavioural tiers:

| Tier | Share | Posts | Comments | Outbound follows | Conversations | Avatar chance |
|---|---|---|---|---|---|---|
| Lurker | ~60% | 0–2 | 0–2 | 0–2 | 0 | 20% |
| Normal | ~30% | 1–6 | 2–8 | 3–12 | 0–3 | 70% |
| Active | ~9% | 5–30 | 10–40 | 20–80 | 3–8 | 95% |
| Celebrity | ~1% (capped at 50) | 30–200 | 20–80 | 20–80 | 5–15 | 100% |

Celebrities also act as **follow magnets**: each one attracts inbound follows from a random `N/20` … `N/5` slice of the population, so their `followers_count` looks realistic relative to everyone else.

All bulk inserts happen inside a single transaction with prepared statements, so:

- `SCALE=5000` finishes in seconds
- `SCALE=50000` finishes in roughly a couple of minutes

---

## Generated images

On every run, `cmd/seed` synthesises a small set of PNG avatars and post media into your `UPLOAD_PATH` (default `backend/data/uploads/`). Curated users get a deterministic avatar; procedural users get one based on their tier's `avatarChance`. Roughly 1 in 3 curated posts gets an image; bulk posts get one based on their tier's `postMediaChance`.

---

## Default credentials

Every seeded user shares the same password:

```
password123
```

A few handy sample logins from the curated dataset:

| Identifier | Password |
|---|---|
| `alice@example.com` | `password123` |
| `bob@example.com` | `password123` |
| `cgaldan` (nickname) | any (dev bypass) |
| `testuser` (nickname) | any (dev bypass) |

> The `dev bypass` accounts only accept any password in development environments — do not rely on this in any deployed config.

---

## Reset behaviour

The `-reset` flag deletes rows from every domain table in child-to-parent order so foreign-key cascades stay happy. It also clears `sqlite_sequence`, so auto-increment IDs restart from 1. Tables wiped:

```
notifications, group_event_rsvps, group_events, group_invitations,
group_join_requests, group_members, groups, messages,
conversation_participants, conversations, post_audiences, comments,
posts, follows, sessions, users
```

Migration history (`schema_migrations`) is left alone, so the database structure stays intact.

After seeding, two cosmetic fixes are applied so direct DB inspection matches what the API returns:

- `users.is_online` is reset to `0` for every user (nobody is actually connected).
- `users.following_count` and `users.followers_count` are recomputed from the `follows` table.

---

## Seeding with Docker

The current Docker image does **not** ship the seed binary — it only contains the compiled server. To seed a Dockerized deployment:

**Seed locally, then start Docker.** Stop the backend container, run the seed against a local DB file, then point the container's volume at the seeded file. Practical only for a fresh setup.

---

## Source

- Entry point: [`backend/cmd/seed/main.go`](../backend/cmd/seed/main.go)
- Image generator: [`backend/cmd/seed/images.go`](../backend/cmd/seed/images.go)
- Curated dataset: [`backend/cmd/seed/data/seed.json`](../backend/cmd/seed/data/seed.json)
