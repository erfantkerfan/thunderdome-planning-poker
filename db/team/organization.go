package team

import (
	"context"
	"database/sql"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// OrganizationService represents a PostgreSQL implementation of thunderdome.OrganizationDataSvc.
type OrganizationService struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// OrganizationGet gets an organization
func (d *OrganizationService) OrganizationGet(ctx context.Context, OrgID string) (*thunderdome.Organization, error) {
	var org = &thunderdome.Organization{}

	e := d.DB.QueryRowContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date
        FROM thunderdome.organization o
        WHERE o.id = $1;`,
		OrgID,
	).Scan(
		&org.Id,
		&org.Name,
		&org.CreatedDate,
		&org.UpdatedDate,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("organization_get_by_id query error", zap.Error(e))
		return nil, errors.New("error getting organization")
	}

	return org, nil
}

// OrganizationUserRole gets a user's role in organization
func (d *OrganizationService) OrganizationUserRole(ctx context.Context, UserID string, OrgID string) (string, error) {
	var role string

	e := d.DB.QueryRowContext(ctx,
		`SELECT ou.role
    FROM thunderdome.organization_user ou
    WHERE ou.organization_id = $2 AND ou.user_id = $1;`,
		UserID,
		OrgID,
	).Scan(
		&role,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("organization_get_user_role query error", zap.Error(e))
		return "", errors.New("error getting organization users role")
	}

	return role, nil
}

// OrganizationListByUser gets a list of organizations the user is apart of
func (d *OrganizationService) OrganizationListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.Organization {
	var organizations = make([]*thunderdome.Organization, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.organization o ON ou.organization_id = o.id
        WHERE ou.user_id = $1
        ORDER BY created_date
		LIMIT $2
		OFFSET $3;`,
		UserID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org thunderdome.Organization

			if err := rows.Scan(
				&org.Id,
				&org.Name,
				&org.CreatedDate,
				&org.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("organization_list_by_user query scan error", zap.Error(err))
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("organization_list_by_user query error", zap.Error(err))
	}

	return organizations
}

// OrganizationCreate creates an organization
func (d *OrganizationService) OrganizationCreate(ctx context.Context, UserID string, OrgName string) (*thunderdome.Organization, error) {
	o := &thunderdome.Organization{}

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM thunderdome.organization_create($1, $2);`,
		UserID,
		OrgName,
	).Scan(&o.Id, &o.Name, &o.CreatedDate, &o.UpdatedDate)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to create organization", zap.Error(err))
		return nil, err
	}

	return o, nil
}

// OrganizationUserList gets a list of organization users
func (d *OrganizationService) OrganizationUserList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.OrganizationUser {
	var users = make([]*thunderdome.OrganizationUser, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT u.id, u.name, COALESCE(u.email, ''), ou.role, u.avatar
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.users u ON ou.user_id = u.id
        WHERE ou.organization_id = $1
        ORDER BY ou.created_date
		LIMIT $2
		OFFSET $3;`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr thunderdome.OrganizationUser

			if err := rows.Scan(
				&usr.Id,
				&usr.Name,
				&usr.Email,
				&usr.Role,
				&usr.Avatar,
			); err != nil {
				d.Logger.Ctx(ctx).Error("organization_user_list query scan error", zap.Error(err))
			} else {
				usr.GravatarHash = db.CreateGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("organization_user_list query error", zap.Error(err))
	}

	return users
}

// OrganizationAddUser adds a user to an organization
func (d *OrganizationService) OrganizationAddUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.organization_user (organization_id, user_id, role) VALUES ($1, $2, $3);`,
		OrgID,
		UserID,
		Role,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to add user to organization", zap.Error(err))
		return "", err
	}

	return OrgID, nil
}

// OrganizationRemoveUser removes a user from a organization
func (d *OrganizationService) OrganizationRemoveUser(ctx context.Context, OrganizationID string, UserID string) error {
	_, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.organization_user_remove($1, $2);`,
		OrganizationID,
		UserID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to remove user from organization", zap.Error(err))
		return err
	}

	return nil
}

// OrganizationTeamList gets a list of organization teams
func (d *OrganizationService) OrganizationTeamList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.Team {
	var teams = make([]*thunderdome.Team, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT t.id, t.name, t.created_date, t.updated_date
        FROM thunderdome.organization_team ot
        LEFT JOIN thunderdome.team t ON ot.team_id = t.id
        WHERE ot.organization_id = $1
        ORDER BY t.created_date
		LIMIT $2
		OFFSET $3;`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.Team

			if err := rows.Scan(
				&team.Id,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("organization_team_list query scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("organization_team_list query error", zap.Error(err))
	}

	return teams
}

// OrganizationTeamCreate creates an organization team
func (d *OrganizationService) OrganizationTeamCreate(ctx context.Context, OrgID string, TeamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM thunderdome.organization_team_create($1, $2);`,
		OrgID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to create organization team", zap.Error(err))
		return nil, err
	}

	return t, nil
}

// OrganizationTeamUserRole gets a user's role in organization team
func (d *OrganizationService) OrganizationTeamUserRole(ctx context.Context, UserID string, OrgID string, TeamID string) (string, string, error) {
	var orgRole string
	var teamRole string

	e := d.DB.QueryRowContext(ctx,
		`SELECT ou.role AS orgRole, COALESCE(tu.role, '') AS teamRole
        FROM thunderdome.organization_user ou
        LEFT JOIN thunderdome.team_user tu ON tu.user_id = $1 AND tu.team_id = $3
        WHERE ou.organization_id = $2 AND ou.user_id = $1;`,
		UserID,
		OrgID,
		TeamID,
	).Scan(
		&orgRole,
		&teamRole,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("organization_team_user_role query error", zap.Error(e))
		return "", "", errors.New("error getting organization team users role")
	}

	return orgRole, teamRole, nil
}

// OrganizationDelete deletes an organization
func (d *OrganizationService) OrganizationDelete(ctx context.Context, OrgID string) error {
	_, err := d.DB.ExecContext(ctx,
		`CALL thunderdome.organization_delete($1);`,
		OrgID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("organization_delete query error", zap.Error(err))
		return err
	}

	return nil
}

// OrganizationList gets a list of organizations
func (d *OrganizationService) OrganizationList(ctx context.Context, Limit int, Offset int) []*thunderdome.Organization {
	var organizations = make([]*thunderdome.Organization, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date
        FROM thunderdome.organization o
        ORDER BY o.created_date
		LIMIT $1
		OFFSET $2;`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org thunderdome.Organization

			if err := rows.Scan(
				&org.Id,
				&org.Name,
				&org.CreatedDate,
				&org.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("organization_list scan error", zap.Error(err))
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("organization_list query error", zap.Error(err))
	}

	return organizations
}
