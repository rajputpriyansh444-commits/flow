package api

import (
	"context"
	"testing"

	generatedapi "github.com/rajputpriyansh444-commits/flow/api/v1"
	"github.com/rajputpriyansh444-commits/flow/internal/auth"
	authservice "github.com/rajputpriyansh444-commits/flow/internal/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type changePasswordAuthService struct{ AuthService }

func (changePasswordAuthService) ChangePassword(context.Context, string, string, string) error {
	return authservice.ErrExternalAuthPasswordManagement
}

func TestChangePasswordRejectsExternalUser(t *testing.T) {
	t.Parallel()

	a := &API{authService: changePasswordAuthService{}}
	ctx := auth.WithUser(context.Background(), &auth.User{ID: "external-user-id", Role: auth.RoleViewer})

	result, err := a.ChangePassword(ctx, generatedapi.ChangePasswordRequestObject{
		Body: &generatedapi.ChangePasswordRequest{
			CurrentPassword: "oldpassword",
			NewPassword:     "newpassword1",
		},
	})

	require.NoError(t, err)
	response, ok := result.(generatedapi.ChangePassword403JSONResponse)
	require.True(t, ok)
	assert.Equal(t, generatedapi.ErrorCodeForbidden, response.Code)
	assert.Equal(t, "Password is managed by the authentication provider for this user", response.Message)
}
