package context

import (
	"context"
	"github.com/robbridges/webapp_v2/models"
	"testing"
)

func TestWithUser(t *testing.T) {
	ctx := setup()

	// Retrieve the user from the context
	_, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		t.Fatal("Expected user in context, but not found")
	}
}

func TestUser(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := setup()
		user := User(ctx)

		if user.ID != 1 {
			t.Errorf("Invalid data returned for mock user id")
		}

		if user.Email != "fake@fake.com" {
			t.Errorf("Invalid data returned from mock user email")
		}
	})
	t.Run("Sad path, user never set", func(t *testing.T) {
		ctx := context.Background()
		user := User(ctx)
		if user != nil {
			t.Errorf("This should have been nil")
		}
	})

}

func setup() context.Context {
	// Create a dummy user for testing
	user := &models.User{
		ID:    1,
		Email: "fake@fake.com",
	}

	// Create a context without a user
	ctx := context.Background()

	// Call the WithUser function to add the user to the context
	ctx = WithUser(ctx, user)
	return ctx
}
