/**
 * @Author: steven
 * @Description:
 * @File: domain
 * @Date: 08/01/24 06.20
 */

package app

import (
	"context"
	"errors"
	"github.com/evorts/kevlars/db"
	"time"
)

type Todo struct {
	Id          int        `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

type ITodoManager interface {
	Create(ctx context.Context, item Todo) (int, error)
	GetById(ctx context.Context, id int) (*Todo, error)
	UpdateById(ctx context.Context, item Todo) error
	DeleteById(ctx context.Context, id int) error
}

type todoManager struct {
	db db.Manager
}

func (t *todoManager) Create(ctx context.Context, item Todo) (int, error) {
	q := `insert into todo(title,description) values(?,?)`
	rs, err := t.db.Exec(ctx, t.db.Rebind(q))
	if err != nil {
		return 0, err
	}
	id, errId := rs.LastInsertId()
	return int(id), errId
}

func (t *todoManager) GetById(ctx context.Context, id int) (*Todo, error) {
	todo := new(Todo)
	err := t.db.QueryRow(ctx, `select id, title, description, created_at, updated_at from todo where id = ?`, id).StructScan(todo)
	return todo, err
}

func (t *todoManager) UpdateById(ctx context.Context, item Todo) error {
	if item.Id < 1 {
		return errors.New("invalid id")
	}
	q := `update todo set title = :title, description = :description, updated_at = current_timestamp where id = :id`
	_, err := t.db.NamedExec(ctx, t.db.Rebind(q), item)
	return err
}

func (t *todoManager) DeleteById(ctx context.Context, id int) error {
	q := `delete from todo where id = ?`
	_, err := t.db.Exec(ctx, t.db.Rebind(q), id)
	return err
}

func NewTodo(dbm db.Manager) ITodoManager {
	return &todoManager{db: dbm}
}
