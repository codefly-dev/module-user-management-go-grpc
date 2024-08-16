package business_test

import (
	"backend/pkg/business"
	"backend/pkg/infra"
	"context"
	"testing"
	"time"

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

	// Register a new user
	rootAuthId := "root-auth-id"
	rootEmail := "root@email.com"
	exists, err := store.GetUserByAuthId(ctx, rootAuthId)
	require.NoError(t, err)
	if exists != nil {
		err = store.DeleteUser(ctx, exists.Id)
		require.NoError(t, err)
	}

	registerInput := &business.RegisterUserInput{Email: rootEmail, SignupAuthId: rootAuthId}

	respRegisterUser, err := service.RegisterUser(ctx, registerInput)
	require.NoError(t, err)
	root := respRegisterUser.User
	require.NotNil(t, root)
	require.Equal(t, rootEmail, root.Email)
	org := respRegisterUser.Organization
	require.NotNil(t, org)
	require.Equal(t, business.DefaultOrganizationName, respRegisterUser.Organization.Name)

	userBack, err := service.GetUserByAuthId(ctx, rootAuthId)
	require.NoError(t, err)
	require.NotNil(t, userBack)
	require.Equal(t, root.Email, userBack.Email)

	//adapters.WithService(service)
	//
	//restPort := uint16(50002)
	//config := &adapters.Configuration{
	//	EndpointGrpcPort: 50001,
	//	EndpointHttpPort: shared.Pointer(restPort),
	//}
	//server, err := adapters.NewServer(config)
	//require.NoError(t, err)
	//
	//// At the gRPC level to test out authz
	//ctx = context.Background()
	//w := wool.Get(ctx)
	//w.WithUserAuthID(admin.SignupAuthId)
	//
	//// Login
	//
	//login, err := server.Grpc.Login(w.GRPC().Out(), &gen.LoginRequest{})
	//require.NoError(t, err)
	//require.NotNil(t, login.User)
	//require.Equal(t, admin.Email, login.User.Email)
	//
	//// Create and delete a user
	//createUser, err := server.Grpc.CreateUser(w.GRPC().Out(), &gen.CreateUserRequest{Email: "test@email.com"})
	//require.NoError(t, err)
	//deleteUser, err := server.Grpc.DeleteUser(w.GRPC().Out(), &gen.DeleteUserRequest{UserId: createUser.User.Id})
	//require.NoError(t, err)
	//require.NotNil(t, deleteUser.User)
	//// Delete again should get not found err
	//_, err = server.Grpc.DeleteUser(w.GRPC().Out(), &gen.DeleteUserRequest{UserId: createUser.User.Id})
	//require.Error(t, err)
	//require.True(t, wool.IsNotFound(err))
	//
	//// Assume the guest user -- ensure we can't add a user
	//w.WithUserAuthID(guest.SignupAuthId)
	//_, err = server.Grpc.CreateUser(w.GRPC().Out(), &gen.CreateUserRequest{Email: "willnotwork@email.com"})
	//require.Error(t, err)
	//require.True(t, wool.IsUnauthorized(err))
	//
	//// Try the REST APIs for completeness
	//go func() {
	//	err = server.Start(context.Background())
	//	if err != nil {
	//		panic(err)
	//	}
	//}()
	//
	//timeout := 500 * time.Second
	//for tries := 0; tries < 5; tries++ {
	//	// HTTP Requests
	//
	//	req := &http.Request{
	//		URL: shared.Must(url.Parse(fmt.Sprintf("http://localhost:%d/version", restPort)))}
	//	client := http.Client{Timeout: timeout}
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		// "Ready" check
	//		continue
	//	}
	//	require.Equal(t, resp.StatusCode, http.StatusOK)
	//	body, err := io.ReadAll(resp.Body)
	//	require.NoError(t, err)
	//	var ver gen.VersionResponse
	//	err = protojson.Unmarshal(body, &ver)
	//	require.NoError(t, err)
	//	require.Equal(t, codefly.Version(), ver.Version)
	//
	//	// Login
	//	callContext := context.Background()
	//	w := wool.Get(callContext).In("TestUser")
	//	w.WithUserAuthID(admin.SignupAuthId)
	//
	//	req = &http.Request{
	//		Method: http.MethodPost,
	//		Header: w.HTTP().Headers(),
	//		URL:    shared.Must(url.Parse(fmt.Sprintf("http://localhost:%d/login", restPort)))}
	//	client = http.Client{Timeout: timeout}
	//
	//	resp, err = client.Do(req)
	//	require.NoError(t, err)
	//	require.Equal(t, resp.StatusCode, http.StatusOK)
	//
	//	var u gen.LoginResponse
	//	body, err = io.ReadAll(resp.Body)
	//	require.NoError(t, err)
	//
	//	err = protojson.Unmarshal(body, &u)
	//	require.NoError(t, err)
	//	require.Equal(t, admin.SignupAuthId, u.User.SignupAuthId)
	//	require.Equal(t, admin.Email, u.User.Email)
	//
	//	return
	//}
	//t.Fatal("REST call failed")

}
