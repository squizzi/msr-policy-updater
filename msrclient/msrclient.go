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

var ErrClientNotAuthenticated = errors.New("problem during authentication, check provided username and password")
var ErrClientForbidden = errors.New("provided username and password does not have access to this resource")

// MSRAPIClient is an interface that encapsulates API requests to a MSR.
type MSRAPIClient interface {
	ListAllRepositories()
	ListPollMirrorPolicies()
	ListPushMirrorPolicies()
	UpdatePollMirrorPolicyUsernamePassword()
	UpdatePushMirrorPolicyUsernamePassword()
}

// MsrAPIClient is a client that invokes API requests to a MSR.
// Implements interface MSRAPIClient
type MsrAPIClient struct {
	client        *apiclient.MirantisSecureRegistry
	username      string
	password      string
	host          string
	httpTransport http.RoundTripper
}

// NewMSRAPIClient creates an instance of MSRAPIClient corresponding to the
// given host.
func NewMSRAPIClient(username, password, host string, insecure bool) (*MsrAPIClient, error) {
	transport := httptransport.New(host, "/", []string{"https"})

	if insecure {
		transport.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if username == "" || password == "" {
		return nil, ErrClientNotAuthenticated
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

func (c *MsrAPIClient) ListAllRepositories() ([]*models.ResponsesRepository, error) {
	r, err := c.client.Repositories.ListRepositories(
		repositories.NewListRepositoriesParams().
			WithPageSize(toInt64Ptr(10000)))
	if err != nil {
		return nil, err
	}

	return r.Payload.Repositories, nil
}

func (c *MsrAPIClient) ListPollMirrorPolicies(reponame string) ([]*models.ResponsesPollMirroringPolicy, error) {
	pm, err := c.client.Repositories.ListRepoPollMirroringPolicies(
		repositories.NewListRepoPollMirroringPoliciesParams().
			WithReponame(reponame))

	if err != nil {
		return nil, err
	}

	return pm.Payload, nil
}

func (c *MsrAPIClient) ListPushMirrorPolicies(reponame string) ([]*models.ResponsesMirroringPolicy, error) {
	pm, err := c.client.Repositories.ListRepoMirroringPolicies(
		repositories.NewListRepoMirroringPoliciesParams().
			WithReponame(reponame))

	if err != nil {
		return nil, err
	}

	return pm.Payload, nil
}

func (c *MsrAPIClient) UpdatePushMirrorPolicyUsernamePassword(policyID, username, password string) error {
	_, err := c.client.Repositories.UpdateRepoPushMirroringPolicy(
		repositories.NewUpdateRepoPushMirroringPolicyParams().
			WithPushmirroringpolicyid(policyID).
			WithInitialEvaluation(toBoolPtr(false)).
			WithBody(&models.FormsUpdatePushMirroringPolicy{
				Username: username,
				Password: password,
			}))
	if err != nil {
		return err
	}

	return nil
}

func (c *MsrAPIClient) UpdatePollMirrorPolicyUsernamePassword(policyID, username, password string) error {
	_, err := c.client.Repositories.UpdateRepoPollMirroringPolicy(
		repositories.NewUpdateRepoPollMirroringPolicyParams().
			WithPollmirroringpolicyid(policyID).
			WithBody(&models.FormsUpdatePollMirroringPolicy{
				Username: username,
				Password: password,
			}))
	if err != nil {
		return err
	}

	return nil
}

func toInt64Ptr(i int64) *int64 {
	return &i
}

func toBoolPtr(v bool) *bool {
	return &v
}
