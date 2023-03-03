package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/squizzi/msr-policy-updater/policyupdater"
)

func main() {
	c := Command()

	if err := c.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func Command() *cobra.Command {
	var (
		logLevel string

		msrUsername string
		msrPassword string
		username    string
		password    string
		host        string

		pushMirror bool
		pollMirror bool
		batchSize  int64
	)

	c := &cobra.Command{
		Use: os.Args[0],
		Long: `
MSR Mirroring Policy Password Update Tool

This tool can be used to update all push and poll mirroring policies affliated
with a target MSR domain name with a new username and password combo.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(c *cobra.Command, args []string) error {
			lvl, err := logrus.ParseLevel(logLevel)
			if err != nil {
				logrus.Fatalf("failed to parse defined log-level: %q: %s", logLevel, err)
			}

			logrus.SetLevel(lvl)

			if err := markFlagsRequired(c, []string{"username", "password", "msr-url", "msr-username", "msr-password"}); err != nil {
				logrus.Fatalf("failed to mark flags required: %s", err)
			}

			if !c.Flag("poll-mirroring").Changed && !c.Flag("push-mirroring").Changed {
				err := c.Usage()
				if err != nil {
					logrus.Fatalf("failed to print command usage: %s", err)
				}
				logrus.Fatalf("either 'poll-mirroring' or 'push-mirroring' flags must be specified")
			}

			logrus.Infof("Updating mirroring policies (Push: %t, Poll: %t)", pushMirror, pollMirror)

			u, err := policyupdater.New(
				msrUsername,
				msrPassword,
				username,
				password,
				host,
				pollMirror,
				pushMirror,
				batchSize,
			)
			if err != nil {
				return fmt.Errorf("failed to setup new policy updater: %w", err)
			}

			if err := u.Update(); err != nil {
				return err
			}

			return nil
		},
	}

	c.Flags().StringVar(&logLevel, "log-level", "info", "Log level to use")

	c.Flags().StringVar(&host, "msr-url", "", "URL of MSR to perform update against (required)")
	c.Flags().StringVar(&msrUsername, "msr-username", "", "Username to use when authenticating against the target MSR to update policies (required)")
	c.Flags().StringVar(&msrPassword, "msr-password", "", "Password to use when authenticating against the target MSR to update policies (required)")

	c.Flags().BoolVar(&pollMirror, "poll-mirroring", false, "Update poll mirroring policies")
	c.Flags().BoolVar(&pushMirror, "push-mirroring", false, "Update push mirroring policies")
	c.Flags().StringVarP(&username, "username", "u", "", "Username to update across policies (required)")
	c.Flags().StringVarP(&password, "password", "p", "", "Password to update across policies (required)")

	c.Flags().Int64VarP(&batchSize, "batch-size", "b", 100, "Number of repositories to fetch for updating at a time")

	c.Flags().SortFlags = false

	return c
}

func markFlagsRequired(c *cobra.Command, flagNames []string) error {
	for _, f := range flagNames {
		if err := c.MarkFlagRequired(f); err != nil {
			err := c.Usage()
			if err != nil {
				return err
			}
			return err
		}
	}

	return nil
}
