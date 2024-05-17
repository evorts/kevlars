/**
 * @Author: steven
 * @Description:
 * @File: feature_flag
 * @Date: 18/12/23 08.07
 */

package fflag

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
	Id            int64      `db:"id"`
	Feature       string     `db:"feature"`
	Enabled       bool       `db:"enabled"`
	LastChangedBy string     `db:"last_changed_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

type Manager interface {
	Add(ctx context.Context, records ...Record) error
	ExecWhenEnabled(ctx context.Context, feature string, f func())
	IsEnabled(ctx context.Context, feature string) bool

	Init() error
	MustInit() Manager
}

type manager struct {
	dbw db.Manager
	dbr db.Manager
}

const (
	table = "feature_flag"
)

var (
	columns                  = []string{"feature", "enabled", "last_changed_by"}
	tableExistenceCheckQuery = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf(`select count(table_name) as tableCount from information_schema.tables ist
				       where ist.table_name in ('%s')`, table),
		},
	}
	tableColumnDefinitions = map[db.SupportedDriver][][]string{
		db.DriverPostgreSQL: {
			{"id", "serial", "primary key"},
			{"feature", "varchar(50)", "not null"},
			{"enabled", "tinyint", "default 0"},
			{"last_changed_by", "varchar(30)"},
			{"created_at", "timestamp with time zone", "default current_timestamp"},
			{"updated_at", "timestamp with time zone"},
		},
	}
	tableIndexDefinitions = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf("create unique index if not exists %s_feature_idx on public.%s(feature)", table, table),
			fmt.Sprintf("create index if not exists %s_enabled_idx on public.%s(enabled)", table, table),
		},
	}
)

func (m *manager) ExecWhenEnabled(ctx context.Context, feature string, f func()) {
	if m.IsEnabled(ctx, feature) {
		f()
	}
}

func (m *manager) IsEnabled(ctx context.Context, feature string) bool {
	q := m.dbr.Rebind(`select enabled from ` + table + ` where feature = ?`)
	var value sql.NullInt32
	if err := m.dbr.QueryRow(ctx, q, feature).Scan(&value); err != nil {
		return false
	}
	return value.Int32 == 1
}

func (m *manager) Add(ctx context.Context, records ...Record) error {
	builder := sqlbuilder.NewInsertBuilder().
		InsertInto(table).
		Cols(columns...)
	for _, record := range records {
		builder.Values(record.Feature, record.Enabled, record.LastChangedBy)
	}
	sql, args := builder.BuildWithFlavor(m.dbw.Driver().ToSqlBuilderFlavor())
	_, err := m.dbw.Exec(ctx, sql, args...)
	return err
}

func (m *manager) Init() error {
	ctx := context.Background()
	return m.initSchema(ctx)
}

func (m *manager) MustInit() Manager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *manager) getFlavorByDriver(driver db.SupportedDriver) sqlbuilder.Flavor {
	if driver == db.DriverPostgreSQL {
		return sqlbuilder.PostgreSQL
	}
	panic("unsupported flavor")
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

func (m *manager) initSchema(ctx context.Context) error {
	driver := m.dbw.Driver()
	flavor := m.getFlavorByDriver(driver)
	// to avoid unnecessary execution of schema scaffolding,
	// check the existence of tables -- should return total table of 2
	total, err := m.tableCheck(ctx, driver)
	if err != nil {
		return err
	}
	if total < 1 {
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

func New(db db.Manager, opts ...common.Option[manager]) Manager {
	m := &manager{dbw: db, dbr: db}
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}
