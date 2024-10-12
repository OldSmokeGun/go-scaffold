package command

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"go-scaffold/internal/config"
)

var (
	flagAppName        = flag{"name", "n", "go-scaffold", "set the application name"}
	flagAppEnvironment = flag{"env", "e", "dev", "set the application environment (dev, test, prod)"}

	flagConfig = flag{"config", "f", "./etc/config.yaml", "configuration file path"}

	flagRemoteConfigEnable     = flag{"config.remote.enable", "", false, "enable remote config"}
	flagRemoteConfigEndpoints  = flag{"config.remote.endpoints", "", []string{"http://localhost:2379"}, "remote config endpoint"}
	flagRemoteConfigTimeout    = flag{"config.remote.timeout", "", time.Second * 3, "remote config timeout"}
	flagRemoteConfigPathPrefix = flag{"config.remote.path-prefix", "", "/go-scaffold/etc", "remote config path prefix"}

	flagLoggerPath   = flag{"log.path", "", "logs/%Y%m%d.log", "log output path"}
	flagLoggerLevel  = flag{"log.level", "", "info", "log level (debug, info, warn, error)"}
	flagLoggerFormat = flag{"log.format", "", "json", "log output format (text, json)"}

	flagMigrationDir           = flag{"migration", "m", "migrations", "migration directory"}
	flagMigrationDBGroup       = flag{"db-group", "", "default", "migration database group"}
	flagMigrationIgnoreUnknown = flag{"ignore-unknown", "", false, "whether to skip checking the database for migrations that are not in the migration source"}
)

type flag struct {
	name         string
	shortName    string
	defaultValue interface{}
	usage        string
}

func getFlags(cmd *cobra.Command, persistent bool) *pflag.FlagSet {
	flags := cmd.Flags()
	if persistent {
		flags = cmd.PersistentFlags()
	}
	return flags
}

func addAppRuntimeFlag(cmd *cobra.Command, persistent bool) {
	getFlags(cmd, persistent).StringP(flagAppName.name, flagAppName.shortName, flagAppName.defaultValue.(string), flagAppName.usage)
	getFlags(cmd, persistent).StringP(flagAppEnvironment.name, flagAppEnvironment.shortName, flagAppEnvironment.defaultValue.(string), flagAppEnvironment.usage)
}

func addConfigFlag(cmd *cobra.Command, persistent bool) {
	getFlags(cmd, persistent).StringP(flagConfig.name, flagConfig.shortName, flagConfig.defaultValue.(string), flagConfig.usage)
}

func addRemoteConfigFlag(cmd *cobra.Command, persistent bool) {
	flags := getFlags(cmd, persistent)
	flags.BoolP(flagRemoteConfigEnable.name, flagRemoteConfigEnable.shortName, flagRemoteConfigEnable.defaultValue.(bool), flagRemoteConfigEnable.usage)
	flags.StringSliceP(flagRemoteConfigEndpoints.name, flagRemoteConfigEndpoints.shortName, flagRemoteConfigEndpoints.defaultValue.([]string), flagRemoteConfigEndpoints.usage)
	flags.DurationP(flagRemoteConfigTimeout.name, flagRemoteConfigTimeout.shortName, flagRemoteConfigTimeout.defaultValue.(time.Duration), flagRemoteConfigTimeout.usage)
	flags.StringP(flagRemoteConfigPathPrefix.name, flagRemoteConfigPathPrefix.shortName, flagRemoteConfigPathPrefix.defaultValue.(string), flagRemoteConfigPathPrefix.usage)
}

func addLoggerFlag(cmd *cobra.Command, persistent bool) {
	flags := getFlags(cmd, persistent)
	flags.StringP(flagLoggerPath.name, flagLoggerPath.shortName, flagLoggerPath.defaultValue.(string), flagLoggerPath.usage)
	flags.StringP(flagLoggerLevel.name, flagLoggerLevel.shortName, flagLoggerLevel.defaultValue.(string), flagLoggerLevel.usage)
	flags.StringP(flagLoggerFormat.name, flagLoggerFormat.shortName, flagLoggerFormat.defaultValue.(string), flagLoggerFormat.usage)
}

func addMigrationFlag(cmd *cobra.Command, persistent bool) {
	getFlags(cmd, persistent).StringP(flagMigrationDir.name, flagMigrationDir.shortName, flagMigrationDir.defaultValue.(string), flagMigrationDir.usage)
	getFlags(cmd, persistent).StringP(flagMigrationDBGroup.name, flagMigrationDBGroup.shortName, flagMigrationDBGroup.defaultValue.(string), flagMigrationDBGroup.usage)
	getFlags(cmd, persistent).BoolP(flagMigrationIgnoreUnknown.name, flagMigrationIgnoreUnknown.shortName, flagMigrationIgnoreUnknown.defaultValue.(bool), flagMigrationIgnoreUnknown.usage)
}

func getAppName(cmd *cobra.Command) config.AppName {
	return config.AppName(cmd.Flag(flagAppName.name).Value.String())
}

func getAppEnvironment(cmd *cobra.Command) config.Env {
	return config.Env(cmd.Flag(flagAppEnvironment.name).Value.String())
}
