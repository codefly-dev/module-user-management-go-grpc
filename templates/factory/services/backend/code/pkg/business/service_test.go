package business_test

import (
	"backend/pkg/adapters"
	"backend/pkg/business"
	"backend/pkg/gen"
	"backend/pkg/infra"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	codefly "github.com/codefly-dev/sdk-go"

	"github.com/codefly-dev/core/shared"

	"github.com/codefly-dev/core/cli"

	"github.com/codefly-dev/core/wool"

	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	ctx := context.Background()
	wool.SetGlobalLogLevel(wool.DEBUG)

	deps, err := cli.WithDependencies(ctx, cli.WithDebug(),
		cli.WithNamingScope("test"), cli.WithTimeout(30*time.Second),
		cli.WithSilence("store"))
	require.NoError(t, err)

	defer func() {
		err = deps.Destroy(ctx)
		require.NoError(t, err)
	}()

	store, err := infra.NewPostgresStore(ctx)
	require.NoError(t, err)

	service, err := business.NewService(store)
	require.NoError(t, err)

	authID := "test-auth-id"
	email := "email@test.com"

	u := &gen.User{
		SignupAuthId: authID,
		Email:        email,
	}

	_, err = service.DeleteOwner(ctx, authID)
	require.NoError(t, err)

	respRegisterUser, err := service.RegisterUser(ctx, u)
	require.NoError(t, err)
	require.Equal(t, email, respRegisterUser.User.Email)
	require.Equal(t, business.DefaultOrganizationName, respRegisterUser.Organization.Name)

	userBack, err := service.GetUserByAuthID(ctx, authID)
	require.NoError(t, err)
	require.Equal(t, email, userBack.Email)
	require.Equal(t, authID, userBack.SignupAuthId)

	org, err := service.GetOrganizationForOwner(ctx, u)
	require.NoError(t, err)
	require.Equal(t, org.Name, business.DefaultOrganizationName)

	teams, err := service.GetTeams(ctx, org)
	require.NoError(t, err)
	require.Len(t, teams, 1)
	require.Equal(t, teams[0].Name, "Administrators")

	// Try the REST API
	adapters.WithService(service)

	restPort := uint16(50002)
	config := &adapters.Configuration{
		EndpointGrpcPort: 50001,
		EndpointHttpPort: shared.Pointer(restPort),
	}

	server, err := adapters.NewServer(config)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Start(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	timeout := 500 * time.Second
	for tries := 0; tries < 5; tries++ {
		// HTTP Request

		req := &http.Request{
			URL: shared.Must(url.Parse(fmt.Sprintf("http://localhost:%d/version", restPort)))}
		client := http.Client{Timeout: timeout}
		resp, err := client.Do(req)
		if err != nil {
			// "Ready" check
			continue
		}
		require.Equal(t, resp.StatusCode, http.StatusOK)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var ver gen.VersionResponse
		err = protojson.Unmarshal(body, &ver)
		require.NoError(t, err)
		require.Equal(t, codefly.Version(), ver.Version)

		// Login
		callContext := context.Background()
		w := wool.Get(callContext).In("TestUser")
		w.WithUserAuthID(authID)

		req = &http.Request{
			Method: http.MethodPost,
			Header: w.HTTP().Headers(),
			URL:    shared.Must(url.Parse(fmt.Sprintf("http://localhost:%d/login", restPort)))}
		client = http.Client{Timeout: timeout}

		resp, err = client.Do(req)
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, http.StatusOK)

		var u gen.LoginResponse
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		err = protojson.Unmarshal(body, &u)
		require.NoError(t, err)
		require.Equal(t, authID, u.User.SignupAuthId)
		require.Equal(t, email, u.User.Email)

		// Get organization
		req = &http.Request{
			Header: w.HTTP().Headers(),
			URL:    shared.Must(url.Parse(fmt.Sprintf("http://localhost:%d/organization", restPort)))}
		client = http.Client{Timeout: timeout}
		resp, err = client.Do(req)
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, http.StatusOK)
		var rO gen.Organization
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		err = protojson.Unmarshal(body, &rO)
		require.NoError(t, err)
		require.Equal(t, business.DefaultOrganizationName, rO.Name)

		return
	}
	t.Fatal("REST call failed")

}
