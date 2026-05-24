package main

import (
	"fmt"
	"os"
	"strconv"

	"social-network/internal/config"
	"social-network/internal/database"
)

const usage = `Usage:
  migrate up                 Apply all pending up migrations
  migrate down               Roll back ALL migrations (destructive — empties schema)
  migrate down <N>           Roll back N migrations
  migrate up <N>             Apply N pending migrations
  migrate version            Print current schema version

Database path is taken from DATABASE_PATH (default ./data/database/social.db).`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	db, err := database.NewDatabase(cfg.Database.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	cmd := os.Args[1]
	switch cmd {
	case "up":
		if len(os.Args) == 2 {
			if err := database.RunMigrations(db); err != nil {
				fail(err)
			}
			fmt.Println("up: applied all pending migrations")
			return
		}
		n, err := parseSteps(os.Args[2])
		if err != nil {
			fail(err)
		}
		if err := database.MigrateSteps(db, n); err != nil {
			fail(err)
		}
		fmt.Printf("up: applied %d migration(s)\n", n)

	case "down":
		if len(os.Args) == 2 {
			if err := database.RollbackMigrations(db); err != nil {
				fail(err)
			}
			fmt.Println("down: rolled back all migrations")
			return
		}
		n, err := parseSteps(os.Args[2])
		if err != nil {
			fail(err)
		}
		if err := database.MigrateSteps(db, -n); err != nil {
			fail(err)
		}
		fmt.Printf("down: rolled back %d migration(s)\n", n)

	case "version":
		v, dirty, err := database.MigrationVersion(db)
		if err != nil {
			fail(err)
		}
		if v == 0 && !dirty {
			fmt.Println("version: no migrations applied")
			return
		}
		fmt.Printf("version: %d (dirty=%t)\n", v, dirty)

	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
}

func parseSteps(arg string) (int, error) {
	n, err := strconv.Atoi(arg)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("expected a positive integer step count, got %q", arg)
	}
	return n, nil
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
