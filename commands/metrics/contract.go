package metrics

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

type MetricType uint8

const (
	CPU MetricType = iota
	Memory
	NetworkIn
	NetworkOut
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "metrics",
	ShortHelp: "Print service and environment metrics in your local time zone",
	LongHelp: "The `metrics` command gives access to environment metrics or individual service metrics through a variety of formats. " +
		"This is useful for checking on the status and performance of your application or environment as a whole. " +
		"The metrics command cannot be run directly but has sub commands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(CPUSubCmd.Name, CPUSubCmd.ShortHelp, CPUSubCmd.CmdFunc(settings))
			cmd.Command(MemorySubCmd.Name, MemorySubCmd.ShortHelp, MemorySubCmd.CmdFunc(settings))
			cmd.Command(NetworkInSubCmd.Name, NetworkInSubCmd.ShortHelp, NetworkInSubCmd.CmdFunc(settings))
			cmd.Command(NetworkOutSubCmd.Name, NetworkOutSubCmd.ShortHelp, NetworkOutSubCmd.CmdFunc(settings))
		}
	},
}

var CPUSubCmd = models.Command{
	Name:      "cpu",
	ShortHelp: "Print service and environment CPU metrics in your local time zone",
	LongHelp: "`metrics cpu` prints out CPU metrics for your environment or individual services. " +
		"You can print out metrics in csv, json, plain text, or spark lines format. " +
		"If you want plain text format, simply omit the `--json`, `--csv`, and `--spark` flags. " +
		"You can only stream metrics using plain text or spark lines formats. " +
		"To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. " +
		"Otherwise you may choose a service, such as an app service, to retrieve metrics for. " +
		"Here are some sample commands\n\n" +
		"```catalyze metrics cpu\n" +
		"catalyze metrics cpu app01 --stream\n" +
		"catalyze metrics cpu --json\n" +
		"catalyze metrics cpu db01 --csv -m 60```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := subCmd.BoolOpt("json", false, "Output the data as json")
			csv := subCmd.BoolOpt("csv", false, "Output the data as csv")
			spark := subCmd.BoolOpt("spark", false, "Output the data using spark lines")
			stream := subCmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := subCmd.IntOpt("m mins", 1, "How many minutes worth of metrics to retrieve.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdMetrics(*serviceName, CPU, *json, *csv, *spark, *stream, *mins, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]"
		}
	},
}

var MemorySubCmd = models.Command{
	Name:      "memory",
	ShortHelp: "Print service and environment memory metrics in your local time zone",
	LongHelp: "`metrics memory` prints out memory metrics for your environment or individual services. " +
		"You can print out metrics in csv, json, plain text, or spark lines format. " +
		"If you want plain text format, simply omit the `--json`, `--csv`, and `--spark` flags. " +
		"You can only stream metrics using plain text or spark lines formats. " +
		"To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. " +
		"Otherwise you may choose a service, such as an app service, to retrieve metrics for. " +
		"Here are some sample commands\n\n" +
		"```catalyze metrics memory\n" +
		"catalyze metrics memory app01 --stream\n" +
		"catalyze metrics memory --json\n" +
		"catalyze metrics memory db01 --csv -m 60```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := subCmd.BoolOpt("json", false, "Output the data as json")
			csv := subCmd.BoolOpt("csv", false, "Output the data as csv")
			spark := subCmd.BoolOpt("spark", false, "Output the data using spark lines")
			stream := subCmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := subCmd.IntOpt("m mins", 1, "How many minutes worth of metrics to retrieve.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdMetrics(*serviceName, Memory, *json, *csv, *spark, *stream, *mins, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]"
		}
	},
}

var NetworkInSubCmd = models.Command{
	Name:      "network-in",
	ShortHelp: "Print service and environment received network data metrics in your local time zone",
	LongHelp: "`metrics network-in` prints out received network metrics for your environment or individual services. " +
		"You can print out metrics in csv, json, plain text, or spark lines format. " +
		"If you want plain text format, simply omit the `--json`, `--csv`, and `--spark` flags. " +
		"You can only stream metrics using plain text or spark lines formats. " +
		"To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. " +
		"Otherwise you may choose a service, such as an app service, to retrieve metrics for. Here are some sample commands\n\n" +
		"```catalyze metrics network-in\n" +
		"catalyze metrics network-in app01 --stream\n" +
		"catalyze metrics network-in --json\n" +
		"catalyze metrics network-in db01 --csv -m 60```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := subCmd.BoolOpt("json", false, "Output the data as json")
			csv := subCmd.BoolOpt("csv", false, "Output the data as csv")
			spark := subCmd.BoolOpt("spark", false, "Output the data using spark lines")
			stream := subCmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := subCmd.IntOpt("m mins", 1, "How many minutes worth of metrics to retrieve.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdMetrics(*serviceName, NetworkIn, *json, *csv, *spark, *stream, *mins, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]"
		}
	},
}

var NetworkOutSubCmd = models.Command{
	Name:      "network-out",
	ShortHelp: "Print service and environment transmitted network data metrics in your local time zone",
	LongHelp: "`metrics network-out` prints out transmitted network metrics for your environment or individual services. " +
		"You can print out metrics in csv, json, plain text, or spark lines format. " +
		"If you want plain text format, simply omit the `--json`, `--csv`, and `--spark` flags. " +
		"You can only stream metrics using plain text or spark lines formats. " +
		"To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. " +
		"Otherwise you may choose a service, such as an app service, to retrieve metrics for. " +
		"Here are some sample commands\n\n" +
		"```catalyze metrics network-out\n" +
		"catalyze metrics network-out app01 --stream\n" +
		"catalyze metrics network-out --json\n" +
		"catalyze metrics network-out db01 --csv -m 60```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := subCmd.BoolOpt("json", false, "Output the data as json")
			csv := subCmd.BoolOpt("csv", false, "Output the data as csv")
			spark := subCmd.BoolOpt("spark", false, "Output the data using spark lines")
			stream := subCmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := subCmd.IntOpt("m mins", 1, "How many minutes worth of metrics to retrieve.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdMetrics(*serviceName, NetworkOut, *json, *csv, *spark, *stream, *mins, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]"
		}
	},
}

// IMetrics
type IMetrics interface {
	RetrieveEnvironmentMetrics(mins int) (*[]models.Metrics, error)
	RetrieveServiceMetrics(mins int, svcID string) (*models.Metrics, error)
}

// SMetrics is a concrete implementation of IMetrics
type SMetrics struct {
	Settings *models.Settings
}

// New returns an instance of IMetrics
func New(settings *models.Settings) IMetrics {
	return &SMetrics{
		Settings: settings,
	}
}
