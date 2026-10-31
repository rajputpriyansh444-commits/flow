package api_test

import (
	"context"

	"github.com/rajputpriyansh444-commits/flow/internal/auth"
)

func adminCtx() context.Context {
	return auth.WithUser(context.Background(), &auth.User{
		Username: "admin",
		Role:     auth.RoleAdmin,
	})
}
