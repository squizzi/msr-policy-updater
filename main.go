package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/squizzi/msr-policy-updater/updater"
)

func main() {
	c := Command()

	if err := c.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   os.Args[0],
		Short: "MSR Mirroring Policy Password Update Tool",
		Long: `
This tool can be used to update all push and poll mirroring policies affliated
with a target MSR domain name with a new username and password combo.`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := markFlagsRequired(c, []string{"username", "password", "msr-url", "msr-username", "msr-password"}); err != nil {
				logrus.Fatalf("failed to mark flags required: %s", err)
			}

			if !c.Flag("poll-mirroring").Changed || !c.Flag("push-mirroring").Changed {
				err := c.Usage()
				if err != nil {
					logrus.Fatalf("failed to print command usage: %s", err)
				}
				logrus.Fatalf("either 'poll-mirroring' or 'push-mirroring' flags must be specified")
			}

			username := c.Flag("username").Value
			password := c.Flag("password").Value
			host := c.Flag("msr-url").Value
			pushMirror := c.Flag("push-mirroring").Changed
			pollMirror := c.Flag("poll-mirroring").Changed

			u := updater.New(
				username.String(),
				password.String(),
				host.String(),
				pollMirror,
				pushMirror)

			if err := u.Update(); err != nil {
				return err
			}

			return nil
		},
	}

	c.Flags().SortFlags = false

	c.PersistentFlags().String("log-level", "info", "Log level to use (default: info)")

	c.PersistentFlags().String("msr-url", "", "URL of MSR to perform update against (required)")
	c.PersistentFlags().String("msr-username", "", "Username to use when authenticating against the target MSR to update policies (required)")
	c.PersistentFlags().String("msr-password", "", "Password to use when authenticating against the target MSR to update policies (required)")

	c.PersistentFlags().Bool("poll-mirroring", false, "Update poll mirroring policies")
	c.PersistentFlags().Bool("push-mirroring", false, "Update push mirroring policies")
	c.PersistentFlags().StringP("username", "u", "", "Username to update across policies (required)")
	c.PersistentFlags().StringP("password", "p", "", "Password to update across policies (required)")

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
