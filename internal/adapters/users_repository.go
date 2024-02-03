package adapters

import (
	"context"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/internal/domain/users"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/sqlhandler"
	"test_ina_bank/pkg/telemetry"
)

type mysqlRepository struct {
	log    *baselogger.Logger
	tracer *telemetry.OtelSdk
	sqlhandler.SqlHandler
}

func NewMysqlRepository(log *baselogger.Logger, sqlHandler sqlhandler.SqlHandler, tracer *telemetry.OtelSdk) *mysqlRepository {
	return &mysqlRepository{
		log:        log,
		tracer:     tracer,
		SqlHandler: sqlHandler,
	}
}

func (m *mysqlRepository) InsertUser(ctx context.Context, in *users.User) (err error) {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.insert_user")
	defer span.End()

	_, err = m.Exec(sql_InsertUser, in.Email, in.Name, in.Password)
	if err != nil {
		return errs.Failed.Wrap(err, "error query user")
	}
	return nil
}

func (m *mysqlRepository) ListUser(ctx context.Context) ([]*users.User, error) {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.list_user")
	defer span.End()

	rows, err := m.Query(sql_ListUser)
	if err != nil {
		return nil, errs.Failed.Wrap(err, "error query task")
	}
	defer rows.Close()

	out := make([]*users.User, 0)
	for rows.Next() {
		each := &users.User{}
		dataRows := []interface{}{
			&each.Id,
			&each.Email,
			&each.Name,
		}
		err = rows.Scan(dataRows...)
		if err != nil {
			return nil, errs.Failed.Wrap(err, "error scan value")
		}
		out = append(out, each)
	}
	return out, nil
}

func (m *mysqlRepository) GetUserById(ctx context.Context, id int) (*users.User, error) {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.user_by_id")
	defer span.End()

	row, err := m.Query(sql_UserById, id)
	if err != nil {
		return nil, errs.Failed.Wrap(err, "error query task")
	}
	defer row.Close()

	out := &users.User{}
	for row.Next() {
		dataRows := []interface{}{
			&out.Id,
			&out.Email,
			&out.Name,
		}
		err = row.Scan(dataRows...)
		if err != nil {
			return nil, errs.Failed.Wrap(err, "error scan value")
		}
	}
	return out, nil
}

func (m *mysqlRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.user_by_email")
	defer span.End()

	row, err := m.Query(sql_GetUserByEmail, email)
	if err != nil {
		return nil, errs.Failed.Wrap(err, "error query task")
	}
	defer row.Close()

	out := &users.User{}
	for row.Next() {
		dataRows := []interface{}{
			&out.Id,
			&out.Email,
			&out.Name,
			&out.Password,
		}
		err = row.Scan(dataRows...)
		if err != nil {
			return nil, errs.Failed.Wrap(err, "error scan value")
		}
	}
	return out, nil
}

func (m *mysqlRepository) UpdateUser(ctx context.Context, id int, updateFn func(u *users.User) (*users.User, error)) error {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.update_user")
	defer span.End()

	_, err := m.Transaction(
		func() (interface{}, error) {
			existing, err := m.GetUserById(ctx, id)
			if err != nil {
				return nil, err
			}
			if !existing.IsExist() {
				return nil, errs.Invalidated.Wrap(errs.ErrUserNotFound, "")
			}

			updated, err := updateFn(existing)
			if err != nil {
				return nil, err
			}

			_, err = m.Exec(sql_UpdateUserById, updated.Email, updated.Name, updated.Password, updated.Id)
			if err != nil {
				return nil, errs.Failed.Wrap(err, "error query task")
			}

			return updated, nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlRepository) DeleteUser(ctx context.Context, id int) error {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.delete_user")
	defer span.End()

	existing, err := m.GetUserById(ctx, id)
	if err != nil {
		return err
	}
	if !existing.IsExist() {
		return errs.Invalidated.Wrap(errs.ErrUserNotFound, "")
	}

	_, err = m.Exec(sql_DeleteUserById, id)
	if err != nil {
		return errs.Failed.Wrap(err, "error query task")
	}
	return nil
}
