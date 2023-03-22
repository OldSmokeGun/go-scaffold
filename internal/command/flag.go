package command

import (
	"go-scaffold/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	flagAppName        = flag{"name", "n", "go-scaffold", "set the application name"}
	flagAppEnvironment = flag{"env", "e", "dev", "set the application environment (dev, test, prod)"}

	flagConfig = flag{"config", "f", "./etc/config.yaml", "configuration file path"}

	flagApolloConfigEnable    = flag{"config.apollo.enable", "", false, "enable apollo"}
	flagApolloConfigEndpoint  = flag{"config.apollo.endpoint", "", "http://localhost:8080", "apollo endpoint"}
	flagApolloConfigAppID     = flag{"config.apollo.appid", "", "", "apollo appID"}
	flagApolloConfigCluster   = flag{"config.apollo.cluster", "", "default", "apollo cluster"}
	flagApolloConfigNamespace = flag{"config.apollo.namespace", "", "application", "apollo namespace"}
	flagApolloConfigSecret    = flag{"config.apollo.secret", "", "", "apollo secret"}

	flagLoggerPath   = flag{"log.path", "", "logs/%Y%m%d.log", "log output path"}
	flagLoggerLevel  = flag{"log.level", "", "info", "log level (debug, info, warn, error)"}
	flagLoggerFormat = flag{"log.format", "", "json", "log output format (text, json)"}

	migrationDirConfig = flag{"migration", "m", "file://migrations", "migration directory"}
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

func addApolloConfigFlag(cmd *cobra.Command, persistent bool) {
	flags := getFlags(cmd, persistent)
	flags.BoolP(flagApolloConfigEnable.name, flagApolloConfigEnable.shortName, flagApolloConfigEnable.defaultValue.(bool), flagApolloConfigEnable.usage)
	flags.StringP(flagApolloConfigEndpoint.name, flagApolloConfigEndpoint.shortName, flagApolloConfigEndpoint.defaultValue.(string), flagApolloConfigEndpoint.usage)
	flags.StringP(flagApolloConfigAppID.name, flagApolloConfigAppID.shortName, flagApolloConfigAppID.defaultValue.(string), flagApolloConfigAppID.usage)
	flags.StringP(flagApolloConfigCluster.name, flagApolloConfigCluster.shortName, flagApolloConfigCluster.defaultValue.(string), flagApolloConfigCluster.usage)
	flags.StringP(flagApolloConfigNamespace.name, flagApolloConfigNamespace.shortName, flagApolloConfigNamespace.defaultValue.(string), flagApolloConfigNamespace.usage)
	flags.StringP(flagApolloConfigSecret.name, flagApolloConfigSecret.shortName, flagApolloConfigSecret.defaultValue.(string), flagApolloConfigSecret.usage)
}

func addLoggerFlag(cmd *cobra.Command, persistent bool) {
	flags := getFlags(cmd, persistent)
	flags.StringP(flagLoggerPath.name, flagLoggerPath.shortName, flagLoggerPath.defaultValue.(string), flagLoggerPath.usage)
	flags.StringP(flagLoggerLevel.name, flagLoggerLevel.shortName, flagLoggerLevel.defaultValue.(string), flagLoggerLevel.usage)
	flags.StringP(flagLoggerFormat.name, flagLoggerFormat.shortName, flagLoggerFormat.defaultValue.(string), flagLoggerFormat.usage)
}

func addMigrationFlag(cmd *cobra.Command, persistent bool) {
	getFlags(cmd, persistent).StringP(migrationDirConfig.name, migrationDirConfig.shortName, migrationDirConfig.defaultValue.(string), migrationDirConfig.usage)
}

func getAppName(cmd *cobra.Command) config.AppName {
	return config.AppName(cmd.Flag(flagAppName.name).Value.String())
}

func getAppEnvironment(cmd *cobra.Command) config.Env {
	return config.Env(cmd.Flag(flagAppEnvironment.name).Value.String())
}
