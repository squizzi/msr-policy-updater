package policyupdater

import (
	"fmt"

	"github.com/docker/dhe-deploy/gocode/pkg/api-client/models"
	"github.com/sirupsen/logrus"

	"github.com/squizzi/msr-policy-updater/msrclient"
)

type PolicyUpdater struct {
	MSRClient *msrclient.MsrAPIClient

	Username    string
	Password    string
	PollMirrors bool
	PushMirrors bool
	BatchSize   int64
}

func New(msrUsername, msrPassword, username, password, host string, pollMirror, pushMirror bool, batchSize int64) (*PolicyUpdater, error) {
	client, err := msrclient.New(msrUsername, msrPassword, host, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get an MSR API client: %w", err)
	}

	return &PolicyUpdater{
		MSRClient:   client,
		Username:    username,
		Password:    password,
		PollMirrors: pollMirror,
		PushMirrors: pushMirror,
		BatchSize:   batchSize,
	}, nil
}

func (u *PolicyUpdater) Update() error {

	logrus.Infof("Performing policy updates on all repositories in batches of %d", u.BatchSize)

	var (
		nextPage   string
		reposBatch []*models.ResponsesRepository
		err        error
	)
	for {
		logrus.Debugf("Processing repositories in batches of %d (nextPage: %s)", u.BatchSize, nextPage)

		// Iterate repositories in batches until nextPage is empty.
		reposBatch, nextPage, err = u.MSRClient.ListRepositories(u.BatchSize, nextPage)
		if err != nil {
			return fmt.Errorf("failed to list batch of repositories (nextPage: %s, pageSize: %d): %w", nextPage, u.BatchSize, err)
		}

		if err := u.updatePoliciesOnRepos(u.MSRClient, reposBatch); err != nil {
			return err
		}

		if nextPage == "" {
			logrus.Info("No additional batches of repositories to process, done")
			break
		}
	}

	return nil
}

func (u *PolicyUpdater) updatePoliciesOnRepos(client *msrclient.MsrAPIClient, repos []*models.ResponsesRepository) error {
	for _, r := range repos {
		repoName := *r.Namespace + "/" + *r.Name

		if u.PollMirrors {
			logrus.Debugf("Getting list of poll mirror policies for repository: %q", repoName)

			pollPolicies, err := client.ListPollMirrorPolicies(*r.Namespace, *r.Name)
			if err != nil {
				return fmt.Errorf("failed to list poll mirror policies for %q repository: %w", repoName, err)
			}

			for _, poll := range pollPolicies {
				logrus.Debugf("Updating poll mirror policy (id: %q, repository: %q)", *poll.ID, repoName)

				err := client.UpdatePollMirrorPolicyUsernamePassword(*poll.ID, *r.Namespace, *r.Name, u.Username, u.Password)
				if err != nil {
					return fmt.Errorf("failed to update poll mirroring policy (id: %q, repository: %q): %w", *poll.ID, repoName, err)
				}
			}
		}

		if u.PushMirrors {
			logrus.Debugf("Getting list of push mirror policies for repository: %q", repoName)

			pushPolicies, err := client.ListPushMirrorPolicies(*r.Namespace, *r.Name)
			if err != nil {
				return fmt.Errorf("failed to list push mirror policies for %q repository: %w", repoName, err)
			}

			for _, push := range pushPolicies {
				logrus.Debugf("Updating push mirror policy (id: %q, repository: %q)", *push.ID, repoName)

				err := client.UpdatePushMirrorPolicyUsernamePassword(*push.ID, *r.Namespace, *r.Name, u.Username, u.Password)
				if err != nil {
					return fmt.Errorf("failed to update push mirroring policy (id: %q, repository: %q): %w", *push.ID, repoName, err)
				}
			}
		}
	}

	return nil
}
