package auth_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/nedpals/supabase-go"
	"github.com/stretchr/testify/assert"

	. "nory/common/auth"
	"nory/internal/user"
)

func TestAuthMiddleware(t *testing.T) {
	supabaseKey := os.Getenv("SUAPABASE_KEY")
	supabaseUrl := os.Getenv("SUAPABASE_URL")
	email := os.Getenv("AUTH_USER_EMAIL")
	password := os.Getenv("AUTH_USER_PASSWORD")

	if supabaseKey == "" || supabaseUrl == "" || email == "" || password == "" {
		t.Skip()
	}

	userRepository := user.NewUserRepositoryMem()
	supa := supabase.CreateClient(supabaseUrl, supabaseKey)
	a := &Auth{
		SupabaseAuth:   supa.Auth,
		UserRepository: userRepository,
	}

	supaUser, err := a.SupabaseAuth.SignIn(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	assert.Nil(t, err)
	bearer := fmt.Sprintf("Bearer %s", supaUser.AccessToken)

	user, err := a.UserFromBearer(context.Background(), bearer)
	assert.Nil(t, err)
	assert.Equal(t, supaUser.User.ID, user.UserId, "unknown login data received")
	assert.Equal(t, supaUser.User.Email, user.Email, "unknown login data received")
}
