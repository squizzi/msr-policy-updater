package updater

import (
	"fmt"

	"github.com/squizzi/msr-policy-updater/msrclient"
)

type UpdatePolicies struct {
	Username    string
	Password    string
	Host        string
	PollMirrors bool
	PushMirrors bool
}

func New(username, password, host string, pollMirror, pushMirror bool) *UpdatePolicies {
	return &UpdatePolicies{
		Username:    username,
		Password:    password,
		Host:        host,
		PollMirrors: pollMirror,
		PushMirrors: pushMirror,
	}
}

func (u *UpdatePolicies) Update() error {
	client, err := msrclient.NewMSRAPIClient(u.Username, u.Password, u.Host, true)
	if err != nil {
		return fmt.Errorf("failed to get an MSR API client: %w", err)
	}

	repos, err := client.ListAllRepositories()
	if err != nil {
		return fmt.Errorf("failed to list all repositories: %w", err)
	}

	for _, r := range repos {
		if u.PollMirrors {
			pollPolicies, err := client.ListPollMirrorPolicies(*r.Name)
			if err != nil {
				return fmt.Errorf("failed to list poll mirror policies for %q repository: %w", *r.Name, err)
			}

			for _, p := range pollPolicies {
				err := client.UpdatePollMirrorPolicyUsernamePassword(*p.ID, u.Username, u.Password)
				if err != nil {
					return fmt.Errorf("failed to update poll mirroring policy: %q from repository: %q: %w", *p.ID, *r.Name, err)
				}
			}
		}

		if u.PushMirrors {
			pushPolicies, err := client.ListPushMirrorPolicies(*r.Name)
			if err != nil {
				return fmt.Errorf("failed to list push mirror policies for %q repository: %w", *r.Name, err)
			}

			for _, p := range pushPolicies {
				err := client.UpdatePushMirrorPolicyUsernamePassword(*p.ID, u.Username, u.Password)
				if err != nil {
					return fmt.Errorf("failed to update push mirroring policy: %q from repository: %q: %w", *p.ID, *r.Name, err)
				}
			}
		}
	}

	return nil
}
