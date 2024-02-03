package adapters

import (
	"context"
	"test_ina_bank/internal/common/errs"
	"test_ina_bank/internal/domain/tasks"
)

func (m *mysqlRepository) InsertTask(ctx context.Context, in *tasks.Task) (err error) {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.insert_task")
	defer span.End()

	_, err = m.Exec(sql_InsertTask, in.UserId, in.Title, in.Description, tasks.Pending.String())
	if err != nil {
		return errs.Failed.Wrap(err, "error query task")
	}
	return nil
}

func (m *mysqlRepository) ListTask(ctx context.Context) ([]*tasks.TaskModel, error) {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.list_task")
	defer span.End()

	rows, err := m.Query(sql_ListTask)
	if err != nil {
		return nil, errs.Failed.Wrap(err, "error query task")
	}
	defer rows.Close()

	out := make([]*tasks.TaskModel, 0)
	for rows.Next() {
		each := &tasks.TaskModel{}
		dataRows := []interface{}{
			&each.Id,
			&each.UserId,
			&each.Title,
			&each.Description,
			&each.Status,
		}
		err = rows.Scan(dataRows...)
		if err != nil {
			return nil, errs.Failed.Wrap(err, "error scan value")
		}
		out = append(out, each)
	}
	return out, nil
}

func (m *mysqlRepository) GetTaskById(ctx context.Context, id int) (*tasks.TaskModel, error) {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.task_by_id")
	defer span.End()

	row, err := m.Query(sql_GetTaskById, id)
	if err != nil {
		return nil, errs.Failed.Wrap(err, "error query task")
	}
	defer row.Close()

	out := &tasks.TaskModel{}
	for row.Next() {
		dataRows := []interface{}{
			&out.Id,
			&out.UserId,
			&out.Title,
			&out.Description,
			&out.Status,
		}
		err = row.Scan(dataRows...)
		if err != nil {
			return nil, errs.Failed.Wrap(err, "error scan value")
		}
	}
	return out, nil
}

func (m *mysqlRepository) UpdateTask(ctx context.Context, id int, updateFn func(t *tasks.Task) (*tasks.Task, error)) error {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.update_task")
	defer span.End()

	_, err := m.Transaction(
		func() (interface{}, error) {
			existing, err := m.GetTaskById(ctx, id)
			if err != nil {
				return nil, err
			}
			if !existing.IsExist() {
				return nil, errs.Invalidated.New(errs.ErrTaskNotFound.Error())
			}
			taskExisting, err := existing.UnmarshalTask()
			if err != nil {
				return nil, err
			}

			updated, err := updateFn(taskExisting)
			if err != nil {
				return nil, err
			}

			_, err = m.Exec(sql_UpdateTaskById, updated.UserId, updated.Title, updated.Description, updated.Status.String(), updated.Id)
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

func (m *mysqlRepository) DeleteTask(ctx context.Context, id int) error {
	ctx, span := m.tracer.Tracer.Start(ctx, "db.delete_task")
	defer span.End()

	existing, err := m.GetTaskById(ctx, id)
	if err != nil {
		return err
	}
	if !existing.IsExist() {
		return errs.Invalidated.New(errs.ErrTaskNotFound.Error())
	}

	_, err = m.Exec(sql_DeleteTaskById, id)
	if err != nil {
		return errs.Failed.Wrap(err, "error query task")
	}
	return nil
}
