package store

import (
	"context"
	"time"

	"github.com/KORLA2/SocialMedia/models"
)

func NewTestStorage() *Storage {

	return &Storage{
		Users: &TestUserStore{},
	}
}

type TestUserStore struct{}

func (t *TestUserStore) CreateAndInvite(ctx context.Context, user *models.User, token string, expiry time.Duration) error {
	return nil
}
func (t *TestUserStore) Activate(ctx context.Context, token string) error { return nil }
func (t *TestUserStore) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return nil, nil
}
func (t *TestUserStore) GetUserByUserName(ctx context.Context, username string) (*LoginPayload, error) {
	return nil, nil
}
func (t *TestUserStore) Delete(ctx context.Context, id int) error { return nil }
