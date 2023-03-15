package azureclient_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/navikt/aad-developer-groups-monitor/pkg/azureclient"
	"github.com/navikt/aad-developer-groups-monitor/pkg/test"
	"github.com/stretchr/testify/assert"
)

func Test_GetGroup(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()

	t.Run("group does not exist", func(t *testing.T) {
		httpClient := test.NewTestHttpClient(
			func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Contains(t, req.URL.Path, fmt.Sprintf("groups/%s", groupID))
				return test.Response("404 Not Found", `{"error":{"code":"Request_ResourceNotFound"}}`)
			},
		)

		client := azureclient.New(httpClient)
		group, err := client.GetGroup(ctx, groupID)

		assert.ErrorContains(t, err, "404 Not Found")
		assert.Nil(t, group)
	})

	t.Run("successful", func(t *testing.T) {
		httpClient := test.NewTestHttpClient(
			func(req *http.Request) *http.Response {
				return test.Response("200 OK", `{"displayName":"some name"}`)
			},

			func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Contains(t, req.URL.Path, fmt.Sprintf("groups/%s/members/$count", groupID))
				return test.Response("200 OK", "123")
			},
		)

		client := azureclient.New(httpClient)
		group, err := client.GetGroup(ctx, groupID)

		assert.NoError(t, err)
		assert.Equal(t, "some name", group.Name)
		assert.Equal(t, 123, group.NumMembers)
	})
}
