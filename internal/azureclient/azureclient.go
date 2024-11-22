package azureclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"
)

type Group struct {
	Name       string
	NumMembers int
}

type Client interface {
	GetGroup(ctx context.Context, groupID uuid.UUID) (*Group, error)
}

type client struct {
	client *http.Client
}

func NewFromConfig(ctx context.Context, tenantID, clientID, clientSecret string) Client {
	endpoint := microsoft.AzureADEndpoint(tenantID)
	conf := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     endpoint.TokenURL,
		AuthStyle:    endpoint.AuthStyle,
		Scopes:       []string{"https://graph.microsoft.com/.default"},
	}

	return New(conf.Client(ctx))
}

func New(c *http.Client) Client {
	return &client{client: c}
}

func (s *client) GetGroup(ctx context.Context, groupID uuid.UUID) (*Group, error) {
	group, err := s.getGroup(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("get group: %w", err)
	}

	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s/members/$count", groupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("ConsistencyLevel", "eventual")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	num, err := strconv.Atoi(string(body))
	if err != nil {
		return nil, fmt.Errorf("unable to convert response body to an integer: %w", err)
	}

	group.NumMembers = num
	return group, nil
}

func (s *client) getGroup(ctx context.Context, groupID uuid.UUID) (*Group, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s?$select=displayName", groupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected HTTP status: %s: %s", resp.Status, string(body))
	}

	info := &struct {
		Name string `json:"displayName"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return &Group{Name: info.Name}, nil
}
