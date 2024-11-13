package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/malinatrash/tabhub/internal/storage/myErrors"
)

func (s *Storage) CreatePermission(ctx context.Context, projectID int, userID int, ownerID int) error {
	var existingPermissonID int
	query := `SELECT id FROM project_permissions WHERE project_id = $1`
	err := s.db.GetContext(ctx, &existingPermissonID, query, projectID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", myErrors.ErrDBInsert, err)
	}
	if existingPermissonID != 0 {
		return fmt.Errorf("%w: %s", myErrors.ErrPermissonAlreadyExists, existingPermissonID)
	}

	project, err := s.Project(ctx, projectID)
	if err != nil {
		return err
	}

	if ownerID != project.OwnerID {
		return errors.New("project ownership does not match project ownership")
	}

	_, err = s.db.ExecContext(ctx, `INSERT INTO project_permissions (user_id, project_id, created_at, updated_at) VALUES ($1, $2, now(), now())`, userID, projectID)
	if err != nil {
		return fmt.Errorf("%w: %v", myErrors.ErrDBInsert, err)
	}

	return nil
}

func (s *Storage) DeletePermission(ctx context.Context, projectID int, userID int, ownerID int) error {
	var existingPermissionID int
	query := `SELECT id FROM project_permissions WHERE project_id = $1 AND user_id = $2`
	err := s.db.GetContext(ctx, &existingPermissionID, query, projectID, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %v", errors.New("error while select"), err)
	}
	if existingPermissionID == 0 {
		return fmt.Errorf("%w: no permission found", errors.New("permission not found"))
	}

	project, err := s.Project(ctx, projectID)
	if err != nil {
		return err
	}

	if ownerID != project.OwnerID {
		return errors.New("project ownership does not match project ownership")
	}

	_, err = s.db.ExecContext(ctx, `DELETE FROM project_permissions WHERE project_id = $1 AND user_id = $2`, projectID, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errors.New("error while deleting permission"), err)
	}

	return nil
}
