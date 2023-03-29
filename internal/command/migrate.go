package command

import (
	"errors"
	"math"
	"strconv"

	"go-scaffold/internal/config"
	"go-scaffold/pkg/migrate"

	printer "github.com/fatih/color"
	sdkmg "github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

type migrateCmd struct {
	*baseCmd
	migrate *sdkmg.Migrate
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
		newMigrateStepCmd(),
		newMigrateToCmd(),
		newMigrateForceCmd(),
		newMigrateVersionCmd(),
	)

	return c
}

func (c *migrateCmd) initMigrate(cmd *cobra.Command) {
	c.mustConfig()

	dir := cmd.Flag(migrationDirConfig.name).Value.String()
	if dir == "" {
		panic("migration directory must be specified")
	}

	dbConfig, err := config.GetDBConn()
	if err != nil {
		panic(err)
	}

	// supported for multi sql statement
	if err := dbConfig.EnableMultiStatement(); err != nil {
		panic(err)
	}

	db, _, err := initDB(cmd.Context(), dbConfig, nil)
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDB(dir, migrate.Driver(dbConfig.Driver), db)
	if err != nil {
		panic(err)
	}

	c.migrate = m
}

func (c *migrateCmd) closeMigrate() {
	se, de := c.migrate.Close()
	if se != nil {
		panic(se)
	}
	if de != nil {
		panic(de)
	}
}

func (c *migrateCmd) getVersionUint(args []string) (uint, error) {
	arg := args[0]
	n, err := strconv.ParseUint(arg, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}

func (c *migrateCmd) getVersionInt(args []string) (int, error) {
	arg := args[0]
	return strconv.Atoi(arg)
}

func (c *migrateCmd) printError(err error) {
	if errors.Is(err, sdkmg.ErrNilVersion) {
		printer.Red("no migration record")
	} else if errors.Is(err, sdkmg.ErrNoChange) {
		printer.Red("no migration changed")
	} else if err != nil {
		printer.Red("migration error: %s", err)
	}
}

func (c *migrateCmd) printCurrentVersion() {
	version, dirty, err := c.migrate.Version()
	if err != nil {
		c.printError(err)
		return
	}

	printer.Green("current version: %d, dirty status: %t\n", version, dirty)
}

type migrateUpCmd struct {
	*migrateCmd
}

func newMigrateUpCmd() *migrateUpCmd {
	c := &migrateUpCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:   "up",
		Short: "migrate up (migrate all the way up from the currently active migration version)",
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

func (c *migrateUpCmd) run() {
	if err := c.migrate.Up(); err != nil {
		c.printError(err)
		return
	}
	c.printCurrentVersion()
	printer.Green("all migrations are complete")
}

type migrateDownCmd struct {
	*migrateCmd
}

func newMigrateDownCmd() *migrateDownCmd {
	c := &migrateDownCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:   "down",
		Short: "migrate down (migrate all the way down from the currently active migration version)",
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

func (c *migrateDownCmd) run() {
	if err := c.migrate.Down(); err != nil {
		c.printError(err)
		return
	}
	printer.Green("all migrations are rolled back")
}

type migrateStepCmd struct {
	*migrateCmd
}

func newMigrateStepCmd() *migrateStepCmd {
	c := &migrateStepCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:     "step",
		Short:   "migrate from the specified version n (up if n > 0, down if n < 0)",
		Example: "  app migrate step n",
		Args:    cobra.MinimumNArgs(1),
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

func (c *migrateStepCmd) run(args []string) {
	n, err := c.getVersionInt(args)
	if err != nil {
		printer.Red(err.Error())
		return
	}

	if n > 0 {
		printer.Green("migrate %d versions up", n)
	} else if n < 0 {
		printer.Green("migrate %d versions down", int(math.Abs(float64(n))))
	}

	if err := c.migrate.Steps(n); err != nil {
		c.printError(err)
		return
	}
	c.printCurrentVersion()
}

type migrateToCmd struct {
	*migrateCmd
}

func newMigrateToCmd() *migrateToCmd {
	c := &migrateToCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:     "to",
		Short:   "migrate from the currently active version to the specified version n",
		Example: "  app migrate to n",
		Args:    cobra.MinimumNArgs(1),
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

func (c *migrateToCmd) run(args []string) {
	v, err := c.getVersionUint(args)
	if err != nil {
		printer.Red(err.Error())
		return
	}

	printer.Green("attempt to migrate to version %d", v)

	if err := c.migrate.Migrate(v); err != nil {
		c.printError(err)
		return
	}

	c.printCurrentVersion()
}

type migrateForceCmd struct {
	*migrateCmd
}

func newMigrateForceCmd() *migrateForceCmd {
	c := &migrateForceCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:     "force",
		Short:   "force set the currently active migration version (it doesn't check any currently active version in the database, it resets the dirty state of the migration version to false)",
		Example: "  app migrate force n",
		Args:    cobra.MinimumNArgs(1),
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

func (c *migrateForceCmd) run(args []string) {
	v, err := c.getVersionInt(args)
	if err != nil {
		printer.Red(err.Error())
		return
	}

	if err := c.migrate.Force(v); err != nil {
		c.printError(err)
		return
	}

	printer.Green("reset the dirty state to false for version %d, please run migrate again", v)
	c.printCurrentVersion()
}

type migrateVersionCmd struct {
	*migrateCmd
}

func newMigrateVersionCmd() *migrateVersionCmd {
	c := &migrateVersionCmd{&migrateCmd{baseCmd: new(baseCmd)}}

	c.cmd = &cobra.Command{
		Use:   "version",
		Short: "view currently active migration versions",
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

func (c *migrateVersionCmd) run() {
	c.printCurrentVersion()
}
