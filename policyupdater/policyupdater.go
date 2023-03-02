package policyupdater

import (
	"fmt"

	"github.com/docker/dhe-deploy/gocode/pkg/api-client/models"
	"github.com/sirupsen/logrus"

	"github.com/squizzi/msr-policy-updater/msrclient"
)

type PolicyUpdater struct {
	Username    string
	Password    string
	Host        string
	PollMirrors bool
	PushMirrors bool
	BatchSize   int64
}

func New(username, password, host string, pollMirror, pushMirror bool, batchSize int64) *PolicyUpdater {
	return &PolicyUpdater{
		Username:    username,
		Password:    password,
		Host:        host,
		PollMirrors: pollMirror,
		PushMirrors: pushMirror,
		BatchSize:   batchSize,
	}
}

func (u *PolicyUpdater) Update() error {
	client, err := msrclient.NewMSRAPIClient(u.Username, u.Password, u.Host, true)
	if err != nil {
		return fmt.Errorf("failed to get an MSR API client: %w", err)
	}

	logrus.Infof("Performing policy updates on all repositories in batches of %d", u.BatchSize)

	for {
		// Iterate repositories in batches until nextPage is empty.
		reposBatch, nextPage, err := client.ListRepositories(u.BatchSize, "")
		if err != nil {
			return fmt.Errorf("failed to list batch of repositories (page start: %s, page size: %d): %w", nextPage, u.BatchSize, err)
		}

		if err := u.updatePoliciesOnRepos(client, reposBatch); err != nil {
			return err
		}

		if nextPage == "" {
			logrus.Info("No additional batches of repositories to process")
			break
		}
	}

	return nil
}

func (u *PolicyUpdater) updatePoliciesOnRepos(client *msrclient.MsrAPIClient, repos []*models.ResponsesRepository) error {
	for _, r := range repos {
		if u.PollMirrors {
			logrus.Debugf("Getting list of poll mirror policies for repository: %q", *r.Name)

			pollPolicies, err := client.ListPollMirrorPolicies(*r.Name)
			if err != nil {
				return fmt.Errorf("failed to list poll mirror policies for %q repository: %w", *r.Name, err)
			}

			for _, p := range pollPolicies {
				logrus.Debugf("Updating poll mirror policy: %q, repository: %q", *p.ID, *r.Name)

				err := client.UpdatePollMirrorPolicyUsernamePassword(*p.ID, u.Username, u.Password)
				if err != nil {
					return fmt.Errorf("failed to update poll mirroring policy: %q from repository: %q: %w", *p.ID, *r.Name, err)
				}
			}
		}

		if u.PushMirrors {
			logrus.Debugf("Getting list of push mirror policies for repository: %q", *r.Name)

			pushPolicies, err := client.ListPushMirrorPolicies(*r.Name)
			if err != nil {
				return fmt.Errorf("failed to list push mirror policies for %q repository: %w", *r.Name, err)
			}

			for _, p := range pushPolicies {
				logrus.Debugf("Updating push mirror policy: %q, repository: %q", *p.ID, *r.Name)

				err := client.UpdatePushMirrorPolicyUsernamePassword(*p.ID, u.Username, u.Password)
				if err != nil {
					return fmt.Errorf("failed to update push mirroring policy: %q from repository: %q: %w", *p.ID, *r.Name, err)
				}
			}
		}
	}

	return nil
}
