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
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/utils"
	"github.com/huandu/go-sqlbuilder"
	"net/url"
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

type Records []*Record

type Manager interface {
	Add(ctx context.Context, records ...Record) error
	ExecWhenEnabled(ctx context.Context, feature string, f func())
	IsEnabled(ctx context.Context, feature string) bool
	getFeaturesBy(ctx context.Context, by db.IHelper) (Records, error)

	Init() error
	MustInit() Manager
	AddOptions(opts ...common.Option[manager]) Manager

	migrate()
	loadData() error
	dbMigrate() *dbmate.DB
}

type manager struct {
	dbw db.Manager
	dbr db.Manager
	dbm *dbmate.DB
	log logger.Manager

	migrationDir     []string
	migrationEnabled bool

	dataLoaded   bool
	lazyLoadData bool
	mapFeature   map[string]bool
}

const (
	table = "feature_flag"
)

//goland:noinspection SqlResolve
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
			{"enabled", "boolean", "default false"},
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
	rules.WhenTrue(m.lazyLoadData, func() {
		m.log.Info(m.loadData())
	})
	if m.dataLoaded {
		if v, ok := m.mapFeature[feature]; ok {
			return v
		}
	}
	q := m.dbr.Rebind(`select enabled from ` + table + ` where feature = ?`)
	var value sql.NullBool
	if err := m.dbr.QueryRow(ctx, q, feature).Scan(&value); err != nil {
		return false
	}
	return value.Bool
}

func (m *manager) Add(ctx context.Context, records ...Record) error {
	builder := sqlbuilder.NewInsertBuilder().
		InsertInto(table).
		Cols(columns...)
	for _, record := range records {
		builder.Values(record.Feature, record.Enabled, record.LastChangedBy)
	}
	q, args := builder.BuildWithFlavor(m.dbw.Driver().ToSqlBuilderFlavor())
	_, err := m.dbw.Exec(ctx, q, args...)
	return err
}

func (m *manager) AddOptions(opts ...common.Option[manager]) Manager {
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}

func (m *manager) getFeaturesBy(ctx context.Context, by db.IHelper) (Records, error) {
	qf, args := by.BuildSqlAndArgsWithWherePrefix()
	q := fmt.Sprintf(`select id, feature, enabled, last_changed_by, created_at, updated_at from %s where %s`, table, qf)
	rs := make(Records, 0)
	rows, err := m.dbr.Query(ctx, m.dbr.Rebind(q), args...)
	defer func() {
		rules.WhenTrue(rows != nil, func() {
			_ = rows.Close()
		})
	}()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rs, db.ErrorRecordNotFound
		}
		return rs, err
	}
	for rows.Next() {
		var (
			record        Record
			lastChangedBy sql.NullString
			updatedAt     sql.NullTime
		)
		if err = rows.Scan(
			&record.Id, &record.Feature, &record.Enabled,
			&lastChangedBy, &record.CreatedAt, &updatedAt,
		); err != nil {
			return rs, err
		}
		record.LastChangedBy = lastChangedBy.String
		if updatedAt.Valid {
			record.UpdatedAt = &updatedAt.Time
		}
		rs = append(rs, &record)
	}
	return rs, nil
}

func (m *manager) Init() error {
	ctx := context.Background()
	if err := m.initSchema(ctx); err != nil {
		return err
	}
	m.migrate()
	if !m.lazyLoadData {
		if err := m.loadData(); err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) MustInit() Manager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *manager) loadData() error {
	if len(m.mapFeature) > 0 {
		return nil
	}
	var (
		page  = 1
		limit = 20
	)
	for {
		items, err := m.getFeaturesBy(
			context.Background(),
			db.NewHelper(
				db.SeparatorAND,
				db.WithPagination(page, limit),
			),
		)
		if err != nil {
			return err
		}
		// map items into map client authorization
		for _, item := range items {
			m.mapFeature[item.Feature] = item.Enabled
		}
		if len(items) < limit {
			break
		}
		page++
	}
	m.dataLoaded = true
	return nil
}

func (m *manager) migrate() {
	if !m.migrationEnabled || m.migrationDir == nil || len(m.migrationDir) < 1 {
		m.log.Info("migration terms not fulfilled or fs not defined")
		return
	}
	m.dbMigrate().MigrationsDir = m.migrationDir
	m.log.Info("migrations:")
	migrations, err := m.dbMigrate().FindMigrations()
	if err != nil {
		panic(err)
	}
	for _, migration := range migrations {
		m.log.Info(migration.Version, migration.FilePath)
	}
	m.log.Info("applying...")
	err = m.dbMigrate().Migrate()
	if err != nil {
		panic(err)
	}
	return
}

func (m *manager) dbMigrate() *dbmate.DB {
	if m.dbm != nil {
		return m.dbm
	}
	u, err := url.Parse(m.dbw.DSN())
	if err != nil {
		m.log.Error(err.Error())
		return nil
	}
	m.dbm = dbmate.New(u)
	return m.dbm
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

func New(db db.Manager, opts ...common.Option[manager]) Manager {
	m := &manager{
		dbw: db, dbr: db, log: logger.NewNoop(),
		mapFeature:   make(map[string]bool),
		migrationDir: make([]string, 0),
	}
	m.AddOptions(opts...)
	return m
}
