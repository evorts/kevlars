/**
 * @Author: steven
 * @Description:
 * @File: log
 * @Date: 18/12/23 06.38
 */

package audit

import (
	"context"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

type Record struct {
	Id              int64                  `db:"id"`
	Action          string                 `db:"action"`
	CreatedById     string                 `db:"created_by_id"`
	CreatedByName   string                 `db:"created_by_name"`
	Role            string                 `db:"role"`
	BeforeChanged   map[string]interface{} `db:"before_changed"`
	AfterChanged    map[string]interface{} `db:"after_changed"`
	AdditionalProps map[string]interface{} `db:"additional_props"`
	Notes           string                 `db:"notes"`
	CreatedAt       *time.Time             `db:"created_at"`
}

type Manager interface {
	Add(ctx context.Context, record Record) error
	AddMultiple(ctx context.Context, records []Record) error
	Init() error
	MustInit() Manager
}

type manager struct {
	dbw db.Manager
}

const (
	table = "audit_log"
)

var (
	columns = []string{"action", "created_by_id", "created_by_name", "role", "before_changed", "after_changed",
		"additional_props", "notes", "created_at"}
)

//goland:noinspection SqlResolve
func (m *manager) Add(ctx context.Context, record Record) error {
	sql, args := sqlbuilder.NewInsertBuilder().
		InsertInto(table).
		Cols(columns...).
		Values(record.Action, record.CreatedById, record.CreatedByName, record.Role,
			db.ToJsonObjectFromMap(record.BeforeChanged), db.ToJsonObjectFromMap(record.AfterChanged),
			db.ToJsonObjectFromMap(record.AdditionalProps), record.Notes, sqlbuilder.Raw("now()")).
		BuildWithFlavor(m.dbw.Driver().ToSqlBuilderFlavor())
	sql = m.dbw.Rebind(sql)
	err := m.dbw.QueryRow(ctx, sql, args...).Err()
	return err
}

func (m *manager) AddMultiple(ctx context.Context, records []Record) error {
	builder := sqlbuilder.NewInsertBuilder().
		InsertInto(table).
		Cols(columns...)
	for _, record := range records {
		builder.Values(record.Action, record.CreatedById, record.CreatedByName, record.Role,
			db.ToJsonObjectFromMap(record.BeforeChanged), db.ToJsonObjectFromMap(record.AfterChanged),
			db.ToJsonObjectFromMap(record.AdditionalProps), record.Notes, sqlbuilder.Raw("now()"))
	}
	sql, args := builder.BuildWithFlavor(m.dbw.Driver().ToSqlBuilderFlavor())
	_, err := m.dbw.Exec(ctx, sql, args...)
	return err
}

func (m *manager) Init() error {
	var (
		ctx      = context.Background()
		colIdDef = []string{"id", "bigint", "primary key"}
		colJson  []string
		colAt    = []string{"created_at"}
	)
	if m.dbw.Driver() == db.DriverPostgreSQL {
		colIdDef = append(colIdDef, "generated always as identity")
		colJson = append(colJson, "jsonb")
		colAt = append(colAt, "timestamp")
	} else {
		colIdDef = append(colIdDef, "auto_increment")
		colJson = append(colJson, "json")
		colAt = append(colAt, "datetime")
	}
	builder := sqlbuilder.NewCreateTableBuilder().CreateTable(table).IfNotExists().
		Define(colIdDef...).
		Define("action", "varchar(150)", "not null").
		Define("created_by_id", "varchar(50)").
		Define("created_by_name", "varchar(100)").
		Define("role", "varchar(50)").
		Define(append([]string{"before_changed"}, colJson...)...).
		Define(append([]string{"after_changed"}, colJson...)...).
		Define(append([]string{"additional_props"}, colJson...)...).
		Define("notes", "text").
		Define(colAt...).Option("DEFAULT CHARACTER SET", "utf8mb4")

	_, err := m.dbw.Exec(ctx, builder.String())
	if err == nil {
		_, _ = m.dbw.Exec(ctx, `ALTER TABLE `+table+` ADD INDEX(action), ADD INDEX(created_by_id), ADD INDEX(created_at)`)
	}
	return err
}

func (m *manager) MustInit() Manager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func New(db db.Manager, opts ...common.Option[manager]) Manager {
	m := &manager{dbw: db}
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}
