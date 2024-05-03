/**
 * @Author: steven
 * @Description:
 * @File: db_noop
 * @Date: 27/03/24 10.35
 */

package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/evorts/kevlars/telemetry"

	"github.com/jmoiron/sqlx"
)

type managerNoop struct {
	driver SupportedDriver
}

func (m *managerNoop) MustConnect(ctx context.Context) Manager {
	return m
}

func (m *managerNoop) Connect(ctx context.Context) error {
	return errors.New("noop doesnt support this")
}

func (m *managerNoop) SqlMock() sqlmock.Sqlmock {
	return nil
}

func (m *managerNoop) Rebind(query string) string {
	return query
}

func (m *managerNoop) Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return nil, errors.New("noop doesnt support this")
}

func (m *managerNoop) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return nil
}

func (m *managerNoop) NamedQuery(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return nil, errors.New("noop doesnt support this")
}

func (m *managerNoop) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, errors.New("noop doesnt support this")
}

func (m *managerNoop) MustExec(ctx context.Context, query string, args ...interface{}) sql.Result {
	return nil
}

func (m *managerNoop) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return nil, errors.New("noop doesnt support this")
}

func (m *managerNoop) MustBegin(ctx context.Context, opts *sql.TxOptions) *sqlx.Tx {
	return nil
}

func (m *managerNoop) Begin(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return nil, errors.New("noop doesnt support this")
}

func (m *managerNoop) Prepare(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return nil, errors.New("noop doesnt support this")
}

func (m *managerNoop) PrepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return nil, errors.New("noop doesnt support this")
}

func (m *managerNoop) Driver() SupportedDriver {
	return m.driver
}

func (m *managerNoop) Ping() error {
	return nil
}

func (m *managerNoop) SetTelemetry(tm telemetry.Manager) Manager {
	return m
}

func NewNoop() Manager {
	return &managerNoop{driver: DriverMock}
}
