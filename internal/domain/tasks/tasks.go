package tasks

import (
	"test_ina_bank/config"
	"test_ina_bank/internal/common/errs"
)

// StatusType is enum-like type.
type StatusType struct {
	s string
}

func (s StatusType) IsZero() bool {
	return s == StatusType{}
}

func (s StatusType) String() string {
	return s.s
}

var (
	Pending    = StatusType{"pending"}
	OnProgress = StatusType{"on progress"}
	Done       = StatusType{"done"}
)

func NewStatusTypeFromString(statusType string) (StatusType, error) {
	switch statusType {
	case "pending":
		return Pending, nil
	case "on progress":
		return OnProgress, nil
	case "done":
		return Done, nil
	}

	return StatusType{}, errs.Invalidated.New(config.Config.ApplicationMessage[config.Config.General.CurrentLanguage].InvalidStatus.Key)
}

type TaskModel struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

func (t *TaskModel) IsExist() bool {
	return t.Id != 0
}

func (t *TaskModel) UnmarshalTask() (*Task, error) {
	statusType, err := NewStatusTypeFromString(t.Status)
	if err != nil {
		return nil, errs.Invalidated.New(err.Error())
	}

	out := &Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      statusType,
	}
	return out, nil
}

type Task struct {
	Id          int        `json:"id"`
	UserId      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      StatusType `json:"status"`
}

func (t *Task) Assign(userId int) {
	t.UserId = userId
	t.Status = Pending
}

func (t *Task) Replace(task *Task) {
	*t = *task
}

type TaskId struct {
	Id int `json:"id"`
}

type ListTask struct {
	Data []*TaskModel `json:"data"`
}
