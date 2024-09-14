package command

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	printer "github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"

	"go-scaffold/internal/config"
)

// create a migration file like this:
// sql-migrate new create_<table name>_table
// NOTE: you must create dbconfig.yml first!
// more details: https://github.com/rubenv/sql-migrate

type migrateCmd struct {
	*baseCmd
	driver     string
	db         *sql.DB
	cleanup    func()
	migrations *migrate.FileMigrationSource
}

func newMigrateCmd() *migrateCmd {
	c := &migrateCmd{baseCmd: new(baseCmd)}

	c.cmd = &cobra.Command{
		Use:   "migrate",
		Short: "database migration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				panic(err)
			}
		},
	}

	addMigrationFlag(c.cmd, true)

	c.addCommands(
		newMigrateUpCmd(),
		newMigrateDownCmd(),
		newMigrateStatusCmd(),
	)

	return c
}

func (c *migrateCmd) initMigrate(cmd *cobra.Command) {
	c.mustConfig()

	dir := cmd.Flag(flagMigrationDir.name).Value.String()
	if dir == "" {
		panic("migration directory must be specified")
	}
	dbGroup := cmd.Flag(flagMigrationDBGroup.name).Value.String()
	if dbGroup == "" {
		panic("migration database group must be specified")
	}
	ignoreUnknown, err := cmd.Flags().GetBool(flagMigrationIgnoreUnknown.name)
	if err != nil {
		panic(err)
	}

	var dbConfig config.Database

	switch dbGroup {
	case "default":
		fallthrough
	default:
		dbConfig, err = config.GetDefaultDatabase()
		if err != nil {
			panic(err)
		}
	}

	// supported for multi sql statement
	if err := dbConfig.EnableMultiStatement(); err != nil {
		panic(err)
	}

	migrationDir := path.Join(dir, dbGroup)
	migrations := &migrate.FileMigrationSource{
		Dir: migrationDir,
	}
	migrate.SetTable("migrations")
	migrate.SetIgnoreUnknown(ignoreUnknown)

	db, cleanup, err := initDB(cmd.Context(), dbConfig.DatabaseConn, nil)
	if err != nil {
		panic(err)
	}

	c.driver = dbConfig.Driver.String()
	c.db = db
	c.cleanup = cleanup
	c.migrations = migrations
}

func (c *migrateCmd) closeMigrate() {
	c.cleanup()
}

func (c *migrateCmd) getVersionInt(args []string) (int64, error) {
	if len(args) == 0 {
		return 0, nil
	}
	arg := args[0]
	return strconv.ParseInt(arg, 10, 64)
}

func (c *migrateCmd) printError(err error) {
	printer.Red(err.Error())
}

func (c *migrateCmd) applyMigrations(args []string, direction migrate.MigrationDirection) (n int, err error) {
	version, err := c.getVersionInt(args)
	if err != nil {
		return
	}

	if err = c.printPlanMigrations(direction, version); err != nil {
		return
	}

	if version > 0 {
		n, err = migrate.ExecVersion(c.db, c.driver, c.migrations, direction, version)
	} else {
		n, err = migrate.Exec(c.db, c.driver, c.migrations, direction)
	}
	return
}

func (c *migrateCmd) printPlanMigrations(direction migrate.MigrationDirection, version int64) (err error) {
	var planMigrations []*migrate.PlannedMigration

	if version > 0 {
		planMigrations, _, err = migrate.PlanMigrationToVersion(c.db, c.driver, c.migrations, direction, version)
	} else {
		planMigrations, _, err = migrate.PlanMigration(c.db, c.driver, c.migrations, direction, 0)
	}
	if err != nil {
		return err
	}

	if len(planMigrations) == 0 {
		return errors.New("no migration will be applied")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PLAN MIGRATION ID", "ENABLE TRANSACTION"})
	table.SetRowLine(true)
	for _, planMigration := range planMigrations {
		table.Append([]string{planMigration.Id, fmt.Sprintf("%t", !planMigration.DisableTransaction)})
	}
	table.Render()
	return nil
}

func (c *migrateCmd) printRecords() {
	records, err := migrate.GetMigrationRecords(c.db, c.driver)
	if err != nil {
		c.printError(err)
		return
	}

	if len(records) == 0 {
		c.printError(errors.New("no migration records"))
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"MIGRATION ID", "APPLIED TIME"})
	table.SetRowLine(true)
	for _, record := range records {
		table.Append([]string{record.Id, record.AppliedAt.Format(time.DateTime)})
	}
	table.Render()
}

type migrateUpCmd struct {
	*migrateCmd
}

func newMigrateUpCmd() *migrateUpCmd {
	c := &migrateUpCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:   "up [version]",
		Short: "migrate the database to the most recent version available (you can specify the version number)",
		Run: func(cmd *cobra.Command, args []string) {
			c.initRuntime(cmd)
			c.initConfig(cmd)
			defer c.closeConfig()
			c.initMigrate(cmd)
			defer c.closeMigrate()
			c.run(args)
		},
	}

	return c
}

func (c *migrateUpCmd) run(args []string) {
	n, err := c.applyMigrations(args, migrate.Up)
	if err != nil {
		c.printError(err)
		return
	}

	printer.Green("migrations are completed, applied %d migrations!", n)
}

type migrateDownCmd struct {
	*migrateCmd
}

func newMigrateDownCmd() *migrateDownCmd {
	c := &migrateDownCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:   "down [version]",
		Short: "undo a database migration (you can specify the version number)",
		Run: func(cmd *cobra.Command, args []string) {
			c.initRuntime(cmd)
			c.initConfig(cmd)
			defer c.closeConfig()
			c.initMigrate(cmd)
			defer c.closeMigrate()
			c.run(args)
		},
	}

	return c
}

func (c *migrateDownCmd) run(args []string) {
	if len(args) == 0 {
		printer.Yellow("this will roll back all migrations, yes/no?")

		ret := ""
		_, err := fmt.Scan(&ret)
		if err != nil {
			c.printError(err)
			return
		}

		if strings.ToLower(ret) != "yes" && strings.ToLower(ret) != "y" {
			return
		}
	}

	n, err := c.applyMigrations(args, migrate.Down)
	if err != nil {
		c.printError(err)
		return
	}

	printer.Green("migrations are rolled back, applied %d migrations!", n)
}

type migrateStatusCmd struct {
	*migrateCmd
}

func newMigrateStatusCmd() *migrateStatusCmd {
	c := &migrateStatusCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:   "status",
		Short: "show the migration records",
		Run: func(cmd *cobra.Command, args []string) {
			c.initRuntime(cmd)
			c.initConfig(cmd)
			defer c.closeConfig()
			c.initMigrate(cmd)
			defer c.closeMigrate()
			c.run()
		},
	}

	return c
}

func (c *migrateStatusCmd) run() {
	c.printRecords()
}
