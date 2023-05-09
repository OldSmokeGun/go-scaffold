package command

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	"go-scaffold/internal/config"

	printer "github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
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

	ignoreUnknown, err := cmd.Flags().GetBool(flagMigrationIgnoreUnknown.name)
	if err != nil {
		panic(err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: dir,
	}

	dbConfig, err := config.GetDBConn()
	if err != nil {
		panic(err)
	}

	// supported for multi sql statement
	if err := dbConfig.EnableMultiStatement(); err != nil {
		panic(err)
	}

	db, cleanup, err := initDB(cmd.Context(), dbConfig, nil)
	if err != nil {
		panic(err)
	}

	migrate.SetTable("migrations")
	migrate.SetIgnoreUnknown(ignoreUnknown)

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

func (c *migrateCmd) printRecords() {
	records, err := migrate.GetMigrationRecords(c.db, c.driver)
	if err != nil {
		c.printError(err)
		return
	}

	if len(records) == 0 {
		printer.Green("no migration records!")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"MIGRATION ID", "APPLIED TIME"})
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
			c.initConfig(cmd, false)
			defer c.closeConfig()
			c.initMigrate(cmd)
			defer c.closeMigrate()
			c.run(args)
		},
	}

	return c
}

func (c *migrateUpCmd) run(args []string) {
	v, err := c.getVersionInt(args)
	if err != nil {
		printer.Red(err.Error())
		return
	}

	n := 0
	if v > 0 {
		n, err = migrate.ExecVersion(c.db, c.driver, c.migrations, migrate.Up, v)
		if err != nil {
			c.printError(err)
			return
		}
	} else {
		n, err = migrate.Exec(c.db, c.driver, c.migrations, migrate.Up)
		if err != nil {
			c.printError(err)
			return
		}
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
			c.initConfig(cmd, false)
			defer c.closeConfig()
			c.initMigrate(cmd)
			defer c.closeMigrate()
			c.run(args)
		},
	}

	return c
}

func (c *migrateDownCmd) run(args []string) {
	v, err := c.getVersionInt(args)
	if err != nil {
		printer.Red(err.Error())
		return
	}

	n := 0
	if v > 0 {
		n, err = migrate.ExecVersion(c.db, c.driver, c.migrations, migrate.Down, v)
		if err != nil {
			c.printError(err)
			return
		}
	} else {
		n, err = migrate.Exec(c.db, c.driver, c.migrations, migrate.Down)
		if err != nil {
			c.printError(err)
			return
		}
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
			c.initConfig(cmd, false)
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
