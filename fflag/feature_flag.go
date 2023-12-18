/**
 * @Author: steven
 * @Description:
 * @File: feature_flag
 * @Date: 18/12/23 08.07
 */

package fflag

import (
	"context"
	"github.com/evorts/kevlars/db"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

type Record struct {
	Id            int64      `db:"id"`
	Feature       string     `db:"feature"`
	Enabled       bool       `db:"enabled"`
	LastChangedBy string     `db:"last_changed_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

type Manager interface {
	Add(ctx context.Context, record Record) error
	AddMultiple(ctx context.Context, records []Record) error
	ExecWhenEnabled(ctx context.Context, feature string, f func())
	IsEnabled(ctx context.Context, feature string) bool

	Init() error
	MustInit() Manager
}

type manager struct {
	db db.Manager
}

const (
	table = "feature_flag"
)

var (
	columns = []string{"feature", "enabled", "last_changed_by", "created_at"}
)

func (m *manager) ExecWhenEnabled(ctx context.Context, feature string, f func()) {
	if m.IsEnabled(ctx, feature) {
		f()
	}
}

func (m *manager) IsEnabled(ctx context.Context, feature string) bool {
	q := m.db.Rebind(`select enabled from ` + table + ` where feature = ?`)
	var value int
	if err := m.db.QueryRow(ctx, q, feature).Scan(&value); err != nil {
		return false
	}
	return value == 1
}

//goland:noinspection SqlResolve
func (m *manager) Add(ctx context.Context, record Record) error {
	sql, args := sqlbuilder.NewInsertBuilder().
		InsertInto(table).
		Cols(columns...).
		Values(record.Feature, record.Enabled, record.LastChangedBy, sqlbuilder.Raw("now()")).
		BuildWithFlavor(m.db.Driver().ToSqlBuilderFlavor())
	sql = m.db.Rebind(sql)
	err := m.db.QueryRow(ctx, sql, args...).Err()
	return err
}

func (m *manager) AddMultiple(ctx context.Context, records []Record) error {
	builder := sqlbuilder.NewInsertBuilder().
		InsertInto(table).
		Cols(columns...)
	for _, record := range records {
		builder.Values(record.Feature, record.Enabled, record.LastChangedBy, sqlbuilder.Raw("now()"))
	}
	sql, args := builder.BuildWithFlavor(m.db.Driver().ToSqlBuilderFlavor())
	_, err := m.db.Exec(ctx, sql, args...)
	return err
}

func (m *manager) Init() error {
	return m.ensureTableExist()
}

func (m *manager) MustInit() Manager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *manager) ensureTableExist() error {
	var (
		ctx          = context.Background()
		colIdDef     = []string{"id", "bigint", "primary key"}
		colAt        = []string{"created_at"}
		colUpdatedAt = []string{"updated_at"}
	)
	if m.db.Driver() == db.DriverPostgreSQL {
		colIdDef = append(colIdDef, "generated always as identity")
		colAt = append(colAt, "timestamp")
		colUpdatedAt = append(colAt, "timestamp")
	} else {
		colIdDef = append(colIdDef, "auto_increment")
		colAt = append(colAt, "datetime")
		colUpdatedAt = append(colAt, "datetime")
	}
	builder := sqlbuilder.NewCreateTableBuilder().CreateTable(table).IfNotExists().
		Define(colIdDef...).
		Define("feature", "varchar(50)", "not null").
		Define("enabled", "tinyint").
		Define("last_changed_by", "varchar(50)").
		Define(colAt...).
		Define(colUpdatedAt...).Option("DEFAULT CHARACTER SET", "utf8mb4")

	_, err := m.db.Exec(ctx, builder.String())
	if err == nil {
		_, _ = m.db.Exec(ctx, `ALTER TABLE `+table+` ADD INDEX(action), ADD INDEX(feature), ADD INDEX(enabled)`)
	}
	return err
}

func New(dbm db.Manager) Manager {
	return &manager{db: dbm}
}
