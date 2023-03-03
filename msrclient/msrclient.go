package msrclient

import (
	"crypto/tls"
	"errors"
	"net/http"

	apiclient "github.com/docker/dhe-deploy/gocode/pkg/api-client/client"
	"github.com/docker/dhe-deploy/gocode/pkg/api-client/client/repositories"
	"github.com/docker/dhe-deploy/gocode/pkg/api-client/models"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

var (
	ErrMirrorCredsIncorrect = errors.New("problem during authentication using the provided mirroring credentials: check provided mirroring username and password")
	ErrUnauthorized         = errors.New("failed to authenticate with target MSR: check provided msr-username and msr-password")
)

// MsrAPIClient is a client that invokes API requests to a MSR.
// Implements interface MSRAPIClient
type MsrAPIClient struct {
	client        *apiclient.MirantisSecureRegistry
	username      string
	password      string
	host          string
	httpTransport http.RoundTripper
}

// New reates an instance of MsrAPIClient corresponding to the
// given host.
func New(username, password, host string, insecure bool) (*MsrAPIClient, error) {
	transport := httptransport.New(host, "/", []string{"https"})

	if insecure {
		transport.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if username == "" || password == "" {
		return nil, ErrUnauthorized
	}

	transport.DefaultMediaType = "application/json"
	transport.DefaultAuthentication = httptransport.BasicAuth(username, password)

	return &MsrAPIClient{
		client:        apiclient.New(transport, strfmt.Default),
		username:      username,
		password:      password,
		host:          host,
		httpTransport: transport.Transport}, nil
}

func (c *MsrAPIClient) ListRepositories(pageSize int64, pageStart string) ([]*models.ResponsesRepository, string, error) {
	r, err := c.client.Repositories.ListRepositories(
		repositories.NewListRepositoriesParams().
			WithPageStart(&pageStart).
			WithPageSize(&pageSize))
	if err != nil {
		var unauthorizedErr *repositories.ListRepositoriesUnauthorized

		if errors.As(err, &unauthorizedErr) {
			return nil, "", ErrUnauthorized
		}

		return nil, "", err
	}

	return r.Payload.Repositories, r.XNextPageStart, nil
}

func (c *MsrAPIClient) ListPollMirrorPolicies(namespace, name string) ([]*models.ResponsesPollMirroringPolicy, error) {
	pm, err := c.client.Repositories.ListRepoPollMirroringPolicies(
		repositories.NewListRepoPollMirroringPoliciesParams().
			WithNamespace(namespace).
			WithReponame(name))

	if err != nil {
		return nil, err
	}

	return pm.Payload, nil
}

func (c *MsrAPIClient) ListPushMirrorPolicies(namespace, name string) ([]*models.ResponsesPushMirroringPolicy, error) {
	pm, err := c.client.Repositories.ListRepoPushMirroringPolicies(
		repositories.NewListRepoPushMirroringPoliciesParams().
			WithNamespace(namespace).
			WithReponame(name))

	if err != nil {
		return nil, err
	}

	return pm.Payload, nil
}

func (c *MsrAPIClient) UpdatePushMirrorPolicyUsernamePassword(policyID, namespace, name, username, password string) error {
	_, err := c.client.Repositories.UpdateRepoPushMirroringPolicy(
		repositories.NewUpdateRepoPushMirroringPolicyParams().
			WithNamespace(namespace).
			WithReponame(name).
			WithPushmirroringpolicyid(policyID).
			WithBody(&models.FormsUpdatePushMirroringPolicy{
				Username: username,
				Password: password,
			}))
	if err != nil {
		var badRequestErr *repositories.UpdateRepoPushMirroringPolicyBadRequest

		if errors.As(err, &badRequestErr) {
			return ErrMirrorCredsIncorrect
		}

		return err
	}

	return nil
}

func (c *MsrAPIClient) UpdatePollMirrorPolicyUsernamePassword(policyID, namespace, name, username, password string) error {
	_, err := c.client.Repositories.UpdateRepoPollMirroringPolicy(
		repositories.NewUpdateRepoPollMirroringPolicyParams().
			WithNamespace(namespace).
			WithReponame(name).
			WithPollmirroringpolicyid(policyID).
			WithBody(&models.FormsUpdatePollMirroringPolicy{
				Username: username,
				Password: password,
			}))
	if err != nil {
		var badRequestErr *repositories.UpdateRepoPollMirroringPolicyBadRequest

		if errors.As(err, &badRequestErr) {
			return ErrMirrorCredsIncorrect
		}
		return err
	}

	return nil
}
