/**
 * @Author: steven
 * @Description:
 * @File: log
 * @Date: 18/12/23 06.38
 */

package audit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/utils"
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
	Add(ctx context.Context, records ...Record) error
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
	tableExistenceCheckQuery = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf(`select count(table_name) as tableCount from information_schema.tables ist
				       where ist.table_name in ('%s')`, table),
		},
	}
	tableColumnDefinitions = map[db.SupportedDriver][][]string{
		db.DriverPostgreSQL: {
			{"id", "serial", "primary key"},
			{"action", "varchar(150)", "not null"},
			{"created_by_id", "varchar(25)", "not null"},
			{"created_by_name", "varchar(50)", "not null"},
			{"before_changed", "jsonb"},
			{"after_changed", "jsonb"},
			{"additional_props", "jsonb"},
			{"notes", "text"},
			{"created_at", "timestamp with time zone", "default current_timestamp"},
			{"updated_at", "timestamp with time zone"},
		},
	}
	tableIndexDefinitions = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf("create  index if not exists %s_action_idx on public.%s(action)", table, table),
			fmt.Sprintf("create index if not exists %s_created_by_id_idx on public.%s(created_by_id)", table, table),
			fmt.Sprintf("create index if not exists %s_created_at_idx on public.%s(created_at)", table, table),
		},
	}
)

func (m *manager) Add(ctx context.Context, records ...Record) error {
	builder := sqlbuilder.NewInsertBuilder().
		InsertInto(table).
		Cols(columns...)
	for _, record := range records {
		builder.Values(record.Action, record.CreatedById, record.CreatedByName, record.Role,
			db.ToJsonObjectFromMap(record.BeforeChanged), db.ToJsonObjectFromMap(record.AfterChanged),
			db.ToJsonObjectFromMap(record.AdditionalProps), record.Notes, sqlbuilder.Raw("now()"))
	}
	q, args := builder.BuildWithFlavor(m.dbw.Driver().ToSqlBuilderFlavor())
	_, err := m.dbw.Exec(ctx, q, args...)
	return err
}

func (m *manager) tableCheck(ctx context.Context, driver db.SupportedDriver) (int, error) {
	total := 0
	if !utils.KeyExistsInMap(tableExistenceCheckQuery, driver) {
		return total, errors.New("driver not supported by table existence check")
	}
	tableChecks := utils.GetValueOnMap(tableExistenceCheckQuery, driver, []string{})
	if len(tableChecks) < 1 {
		return total, errors.New("no table existence check query exists")
	}
	for _, checkQuery := range tableChecks {
		var count sql.NullInt16
		err := m.dbw.QueryRow(ctx, checkQuery).Scan(&count)
		if err != nil {
			return total, err
		}
		total += int(count.Int16)
	}
	return total, nil
}

func (m *manager) getFlavorByDriver(driver db.SupportedDriver) sqlbuilder.Flavor {
	if driver == db.DriverPostgreSQL {
		return sqlbuilder.PostgreSQL
	}
	panic("unsupported flavor")
}

func (m *manager) Init() error {
	ctx := context.Background()
	driver := m.dbw.Driver()
	flavor := m.getFlavorByDriver(driver)
	// to avoid unnecessary execution of schema scaffolding,
	// check the existence of tables -- should return total table of 2
	total, err := m.tableCheck(ctx, driver)
	if err != nil {
		return err
	}
	if total >= 1 {
		return nil
	}
	columnDefinitions := utils.GetValueOnMap(tableColumnDefinitions, driver, [][]string{})
	if len(columnDefinitions) < 1 {
		return errors.New("clients column definitions is empty")
	}
	indexDefinitions := utils.GetValueOnMap(tableIndexDefinitions, driver, []string{})
	if len(indexDefinitions) < 1 {
		return errors.New("clients index definitions is empty")
	}
	builderTable := sqlbuilder.NewCreateTableBuilder().CreateTable(table).IfNotExists()
	for _, definition := range columnDefinitions {
		builderTable = builderTable.Define(definition...)
	}
	tx := m.dbw.MustBegin(ctx, &sql.TxOptions{})
	q, _ := builderTable.BuildWithFlavor(flavor)
	_, err = tx.Exec(q)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// execute index creation on clients table if not exists
	for _, definition := range indexDefinitions {
		_, err = tx.Exec(definition)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
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
