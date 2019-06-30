package cli

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// localExporterCmd is the top-level cobra command for `localExporter`
	localExporterCmd = &cobra.Command{
		Use:               "prombus",
		PersistentPreRunE: commonSetup,
	}

	verbosity  int
	configPath string
)

// Init initializes the CLI environment for localExporter
func Init() (*cobra.Command, error) {
	localExporterCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "increase verbosity level")

	localExporterCmd.AddCommand(cmdServe)

	return localExporterCmd, nil
}

// commonSetup perform actions commons to all CLI subcommands
func commonSetup(cmd *cobra.Command, cmdArgs []string) error {
	logrus.SetLevel(verbosityLevel(verbosity))

	return nil
}

// verbosityLevel parses `-v` count into logrus log-level
func verbosityLevel(verbCount int) logrus.Level {
	switch verbCount {
	case 0:
		return logrus.WarnLevel
	case 1:
		return logrus.InfoLevel
	case 2:
		return logrus.DebugLevel
	default:
		return logrus.TraceLevel
	}
}
