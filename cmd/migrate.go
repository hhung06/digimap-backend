package cmd

import (
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hhung06/digimap-backend/config"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up [N]",
	Short: "Run all (or N) pending migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down [N]",
	Short: "Revert all (or N) migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMigrateDown,
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current migration version",
	RunE:  runMigrateVersion,
}

var migrateForceCmd = &cobra.Command{
	Use:   "force <version>",
	Short: "Force set migration version (fixes dirty state)",
	Args:  cobra.ExactArgs(1),
	RunE:  runMigrateForce,
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)
}

func newMigrator() (*migrate.Migrate, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	m, err := migrate.New("file://migrations", cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}
	return m, nil
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	m, err := newMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if len(args) == 0 {
		err = m.Up()
	} else {
		n, convErr := strconv.Atoi(args[0])
		if convErr != nil {
			return fmt.Errorf("invalid step count: %w", convErr)
		}
		err = m.Steps(n)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("no migrations to apply")
		return nil
	}
	if err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	fmt.Println("migrations applied successfully")
	return nil
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	m, err := newMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if len(args) == 0 {
		err = m.Down()
	} else {
		n, convErr := strconv.Atoi(args[0])
		if convErr != nil {
			return fmt.Errorf("invalid step count: %w", convErr)
		}
		err = m.Steps(-n)
	}

	if err == migrate.ErrNoChange {
		fmt.Println("no migrations to revert")
		return nil
	}
	if err != nil {
		return fmt.Errorf("migrate down: %w", err)
	}
	fmt.Println("migrations reverted successfully")
	return nil
}

func runMigrateVersion(cmd *cobra.Command, args []string) error {
	m, err := newMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	ver, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("get version: %w", err)
	}
	fmt.Printf("version: %d  dirty: %v\n", ver, dirty)
	return nil
}

func runMigrateForce(cmd *cobra.Command, args []string) error {
	m, err := newMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	ver, convErr := strconv.Atoi(args[0])
	if convErr != nil {
		return fmt.Errorf("invalid version: %w", convErr)
	}
	if err := m.Force(ver); err != nil {
		return fmt.Errorf("force version: %w", err)
	}
	fmt.Printf("forced to version %d\n", ver)
	return nil
}
