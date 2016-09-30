package rollback

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/releases"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "rollback",
	ShortHelp: "Rollback a code service to a specific release",
	LongHelp: "`rollback` is a way to redeploy older versions of your code service. " +
		"You must specify the name of the service to rollback and the name of an existing release to rollback to. " +
		"Releases can be found with the [releases list](#releases-list) command. Here are some sample commands\n\n" +
		"```catalyze rollback code-1 f93ced037f828dcaabccfc825e6d8d32cc5a1883```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to rollback")
			releaseName := cmd.StringArg("RELEASE_NAME", "", "The name of the release to rollback to")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRollback(*serviceName, *releaseName, jobs.New(settings), releases.New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "SERVICE_NAME RELEASE_NAME"
		}
	},
}
