package store

import (
	"context"
	"database/sql"

	"github.com/KORLA2/SocialMedia/models"
)

type RoleStore struct {
	db *sql.DB
}

func (r *RoleStore) GetRoleByName(ctx context.Context, minimumRequiredRole string) (*models.Role, error) {

	query := ` select role_id,level , description from roles where name=$1
`

	role := &models.Role{}

	if err := r.db.QueryRowContext(ctx, query, minimumRequiredRole).
		Scan(&role.ID,
			&role.Level,
			&role.Description,
		); err != nil {
		return nil, err
	}
	return role, nil

}
