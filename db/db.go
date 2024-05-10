/**
 * @Author: steven
 * @Description:
 * @File: db
 * @Date: 17/12/23 23.53
 */

package db

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/telemetry"
	_ "github.com/go-sql-driver/mysql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	otelAttr "go.opentelemetry.io/otel/attribute"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	otelTrace "go.opentelemetry.io/otel/trace"
)

type Manager interface {
	MustConnect(ctx context.Context) Manager
	Connect(ctx context.Context) error
	SqlMock() sqlmock.Sqlmock

	Rebind(query string) string

	Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	NamedQuery(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	MustExec(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	MustBegin(ctx context.Context, opts *sql.TxOptions) *sqlx.Tx
	Begin(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)

	Prepare(ctx context.Context, query string) (*sqlx.Stmt, error)
	PrepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error)

	Driver() SupportedDriver
	Ping() error
	SetTelemetry(tm telemetry.Manager) Manager
}

type SupportedDriver string

func (d SupportedDriver) String() string {
	return string(d)
}

func (d SupportedDriver) ToSqlBuilderFlavor() sqlbuilder.Flavor {
	switch d {
	case DriverPostgreSQL:
		return sqlbuilder.PostgreSQL
	case DriverSqlServer:
		return sqlbuilder.SQLServer
	default:
		return sqlbuilder.MySQL
	}
}

const (
	DriverMySQL      SupportedDriver = "mysql"
	DriverPostgreSQL SupportedDriver = "postgres"
	// DriverSqlServer ref: https://gist.github.com/rossnelson/cbb192d314b6c89b7148e919ab25986c
	DriverSqlServer SupportedDriver = "mssql"
	DriverMock      SupportedDriver = "sqlmock"
)

type manager struct {
	db       *sqlx.DB
	driver   SupportedDriver
	dialect  string
	dsn      string
	mockMode bool
	sqlMock  sqlmock.Sqlmock
	scope    string

	maxOpenConnection int
	maxIdleConnection int

	telemetryEnabled bool
	oTelOpenConnect  bool
	tm               telemetry.Manager
}

func spanName(scope, name string) string {
	if len(scope) > 0 {
		return scope + "." + name
	}
	return name
}

func wrapE[T any](m *manager, ctx context.Context, name, q string, f func(newCtx context.Context) (T, error)) (T, error) {
	if !m.telemetryEnabled {
		return f(ctx)
	}
	opts := make([]otelTrace.SpanStartOption, 0)
	opts = append(opts, otelTrace.WithSpanKind(otelTrace.SpanKindClient))
	if len(q) > 0 {
		opts = append(opts, otelTrace.WithAttributes(
			otelAttr.String("sql", q),
		))
	}
	newCtx, span := m.tm.Tracer().Start(ctx, name, opts...)
	defer span.End()
	return f(newCtx)
}

func wrap[T any](m *manager, ctx context.Context, name, q string, f func(newCtx context.Context) T) T {
	if !m.telemetryEnabled {
		return f(ctx)
	}
	opts := make([]otelTrace.SpanStartOption, 0)
	opts = append(opts, otelTrace.WithSpanKind(otelTrace.SpanKindClient))
	if len(q) > 0 {
		opts = append(opts, otelTrace.WithAttributes(
			otelAttr.String("sql", q),
		))
	}
	newCtx, span := m.tm.Tracer().Start(ctx, name, opts...)
	defer span.End()
	return f(newCtx)
}

func (m *manager) spanName(v string) string {
	return rules.WhenTrueR1(len(m.scope) > 0, func() string {
		return m.scope + ".db." + v
	}, func() string {
		return "db." + v
	})
}

func (m *manager) Rebind(query string) string {
	return m.db.Rebind(query)
}

func (m *manager) Ping() error {
	return m.db.Ping()
}

func (m *manager) Driver() SupportedDriver {
	return m.driver
}

func (m *manager) SetTelemetry(tm telemetry.Manager) Manager {
	m.tm = tm
	return m
}

func (m *manager) SqlMock() sqlmock.Sqlmock {
	return m.sqlMock
}

func (m *manager) Query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return wrapE(m, ctx, m.spanName("query"), query, func(newCtx context.Context) (*sqlx.Rows, error) {
		return m.db.QueryxContext(newCtx, query, args...)
	})
}

func (m *manager) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return wrap(m, ctx, m.spanName("query_row"), query, func(newCtx context.Context) *sqlx.Row {
		return m.db.QueryRowxContext(newCtx, query, args...)
	})
}

func (m *manager) NamedQuery(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return wrapE(m, ctx, m.spanName("query_named"), query, func(newCtx context.Context) (*sqlx.Rows, error) {
		return m.db.NamedQueryContext(newCtx, query, arg)
	})
}

func (m *manager) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return wrapE(m, ctx, m.spanName("exec"), query, func(newCtx context.Context) (sql.Result, error) {
		return m.db.ExecContext(newCtx, query, args...)
	})
}

func (m *manager) MustExec(ctx context.Context, query string, args ...interface{}) sql.Result {
	return wrap(m, ctx, m.spanName("exec_must"), query, func(newCtx context.Context) sql.Result {
		return m.db.MustExecContext(newCtx, query, args...)
	})
}

func (m *manager) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return wrapE(m, ctx, m.spanName("exec_named"), query, func(newCtx context.Context) (sql.Result, error) {
		return m.db.NamedExecContext(ctx, query, arg)
	})
}

func (m *manager) MustBegin(ctx context.Context, opts *sql.TxOptions) *sqlx.Tx {
	return wrap(m, ctx, m.spanName("tx_begin_must"), "", func(newCtx context.Context) *sqlx.Tx {
		return m.db.MustBeginTx(newCtx, opts)
	})
}

func (m *manager) Begin(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return wrapE(m, ctx, m.spanName("tx_begin"), "", func(newCtx context.Context) (*sqlx.Tx, error) {
		return m.db.BeginTxx(newCtx, opts)
	})
}

func (m *manager) Prepare(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return wrapE(m, ctx, m.spanName("query_prepare"), "", func(newCtx context.Context) (*sqlx.Stmt, error) {
		return m.db.PreparexContext(newCtx, query)
	})
}

func (m *manager) PrepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return wrapE(m, ctx, m.spanName("query_prepare_named"), "", func(newCtx context.Context) (*sqlx.NamedStmt, error) {
		return m.db.PrepareNamedContext(newCtx, query)
	})
}

func (m *manager) MustConnect(ctx context.Context) Manager {
	if err := m.Connect(ctx); err != nil {
		panic(err)
	}
	return m
}

func (m *manager) Connect(ctx context.Context) (err error) {
	if m.oTelOpenConnect {
		attrs := make([]otelAttr.KeyValue, 0)
		if m.driver == DriverMySQL {
			attrs = append(attrs, semConv.DBSystemMySQL)
		} else if m.driver == DriverSqlServer {
			attrs = append(attrs, semConv.DBSystemMSSQL)
		} else {
			attrs = append(attrs, semConv.DBSystemPostgreSQL)
		}
		m.db, err = otelsqlx.Open(string(m.driver), m.dsn, otelsql.WithAttributes(attrs...))
		//m.db, err = otelsqlx.ConnectContext(ctx, string(m.driver), m.dsn)
	} else if m.mockMode {
		var mockDB *sql.DB
		mockDB, m.sqlMock, err = sqlmock.New()
		m.db = sqlx.NewDb(mockDB, DriverMock.String())
	} else {
		m.db, err = sqlx.ConnectContext(ctx, string(m.driver), m.dsn)
	}
	if err == nil {
		err = m.db.Ping()
	}
	if m.maxOpenConnection > 0 {
		m.db.SetMaxOpenConns(m.maxOpenConnection)
	}
	return
}

func New(driver SupportedDriver, dsn string, opts ...Option) Manager {
	m := &manager{
		driver: driver,
		dsn:    dsn,
		tm:     telemetry.NewNoop(),
	}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}

func NewWithMock() Manager {
	return &manager{driver: DriverMock, dsn: "", telemetryEnabled: false, mockMode: true}
}
