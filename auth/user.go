/**
 * @Author: steven
 * @Description:
 * @File: auth
 * @Date: 24/12/23 21.50
 */

package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/inmemory"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/utils"
	"github.com/huandu/go-sqlbuilder"
)

type UserManager interface {
	Add(ctx context.Context, records ...UserAuthRecord) error
	Save(ctx context.Context, record UserAuthRecord) (UserAuthRecord, error)
	RemoveByIds(ctx context.Context, ids ...int) error
	RemoveByUserIds(ctx context.Context, userIds ...int) error

	GetByUserIds(ctx context.Context, userIds ...int) (UserAuthRecords, error)

	Init() error
	MustInit() UserManager
}

type userManager struct {
	dbw    db.Manager
	dbr    db.Manager
	im     inmemory.Manager
	log    logger.Manager
	driver db.SupportedDriver
}

const (
	tableUserAuth               = "user_auth"
	inMemoryUserCredsHashKey    = "user_creds"    // user_id -> creds
	inMemoryUserDisabledHashKey = "user_disabled" // user_id -> disabled state
)

//goland:noinspection SqlResolve
var (
	userAuthTableExistenceCheckQuery = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf(`select count(table_name) as tableCount from information_schema.tables ist
				       where ist.table_name in ('%s')`, tableUserAuth),
		},
	}
	userAuthColumnDefinitions = map[db.SupportedDriver][][]string{
		db.DriverPostgreSQL: {
			{"id", "serial", "primary key"},
			{"user_id", "int"},
			{"creds", "varchar(128)"},
			{"disabled", "boolean", "default false"},
			{"created_at", "timestamp with time zone", "default current_timestamp"},
			{"updated_at", "timestamp with time zone"},
			{"disabled_at", "timestamp with time zone"},
			{"expired_at", "timestamp with time zone"},
		},
	}
	userAuthIndexDefinitions = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf("create unique index %s_user_id_uidx on public.%s(user_id)", tableUserAuth, tableUserAuth),
			fmt.Sprintf("create index %s_disabled_idx on public.%s(disabled)", tableUserAuth, tableUserAuth),
			fmt.Sprintf("create index %s_created_at_idx on public.%s(created_at)", tableUserAuth, tableUserAuth),
		},
	}
	userAuthScopesColumnDefinitions = map[db.SupportedDriver][][]string{}
	userAuthScopesIndexDefinitions  = map[db.SupportedDriver][]string{}
	usarAuthSaveQuery               = map[db.SupportedDriver]string{
		db.DriverPostgreSQL: fmt.Sprintf(
			`insert into %s (id, user_id, creds, disabled)
					values(:id, :user_id, :creds, :disabled)
				on conflict (user_id) do 
					update set
						case when disabled <> excluded.disabled
							then disabled = excluded.disabled
						end,
						creds = coalesce(nullif(excluded.creds,''),creds),
					where user_id = :user_id 
				returning id, user_id, disabled
		`, tableUserAuth),
	}
)

func (m *userManager) Add(ctx context.Context, records ...UserAuthRecord) error {
	builder := sqlbuilder.NewInsertBuilder().
		InsertInto(tableUserAuth).
		Cols("user_id", "email", "phone", "creds", "disabled")
	for _, record := range records {
		builder.Values(
			record.UserID, record.Email, record.Phone, record.Creds, record.Disabled,
		)
	}
	q, args := builder.BuildWithFlavor(m.driver.ToSqlBuilderFlavor())
	_, err := m.dbw.Exec(ctx, q, args...)
	return err
}

func (m *userManager) Save(ctx context.Context, record UserAuthRecord) (UserAuthRecord, error) {
	q, ok := usarAuthSaveQuery[m.driver]
	if !ok {
		return record, errors.New("not supported yet")
	}
	rows, err := m.dbw.NamedQuery(ctx, q, record)
	if err != nil {
		return record, err
	}
	err = rows.Scan(&record.ID, &record.Disabled)
	if err != nil {
		return record, err
	}
	return record, nil
}

func (m *userManager) RemoveByIds(ctx context.Context, ids ...int) error {
	builder := sqlbuilder.NewDeleteBuilder().
		DeleteFrom(tableUserAuth)
	builder.Where(builder.In("id", ids))
	q, args := builder.BuildWithFlavor(m.driver.ToSqlBuilderFlavor())
	_, err := m.dbw.Exec(ctx, q, args...)
	return err
}

func (m *userManager) RemoveByUserIds(ctx context.Context, userIds ...int) error {
	builder := sqlbuilder.NewDeleteBuilder().
		DeleteFrom(tableUserAuth)
	builder.Where(builder.In("user_id", userIds))
	q, args := builder.BuildWithFlavor(m.driver.ToSqlBuilderFlavor())
	_, err := m.dbw.Exec(ctx, q, args...)
	return err
}

func (m *userManager) GetByUserIds(ctx context.Context, userIds ...int) (UserAuthRecords, error) {
	rs := make(UserAuthRecords, 0)
	builder := sqlbuilder.NewSelectBuilder().From(tableUserAuth)
	builder.Where(builder.In("user_id", userIds))
	q, args := builder.BuildWithFlavor(m.driver.ToSqlBuilderFlavor())
	rows, err := m.dbw.Query(ctx, q, args...)
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
		var record UserAuthRecord
		if err = rows.StructScan(&record); err != nil {
			return rs, err
		}
		rs = append(rs, &record)
	}
	return rs, nil
}

func (m *userManager) Init() error {
	ctx := context.Background()
	if err := m.initSchema(ctx); err != nil {
		return err
	}
	return nil
}

func (m *userManager) MustInit() UserManager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *userManager) tableCheck(ctx context.Context, driver db.SupportedDriver) (int, error) {
	total := 0
	if !utils.KeyExistsInMap(userAuthTableExistenceCheckQuery, driver) {
		return total, errors.New("driver not supported by table existence check")
	}
	tableChecks := utils.GetValueOnMap(userAuthTableExistenceCheckQuery, driver, []string{})
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

func (m *userManager) initSchema(ctx context.Context) error {
	driver := m.dbw.Driver()
	flavor := getFlavorByDriver(driver)
	// to avoid unnecessary execution of schema scaffolding,
	// check the existence of tables -- should return total table of 2
	total, err := m.tableCheck(ctx, driver)
	if err != nil {
		return err
	}
	if total == 1 {
		return nil
	}
	if !utils.KeyExistsInMap(userAuthColumnDefinitions, driver) {
		return errors.New("driver not supported by column definition")
	}
	columnDefinitions := utils.GetValueOnMap(userAuthColumnDefinitions, driver, [][]string{})
	if len(columnDefinitions) < 1 {
		return errors.New("user auth column definitions is empty")
	}
	indexDefinitions := utils.GetValueOnMap(userAuthIndexDefinitions, driver, []string{})
	if len(indexDefinitions) < 1 {
		return errors.New("user auth index definitions is empty")
	}
	builder := sqlbuilder.NewCreateTableBuilder().CreateTable(tableUserAuth).IfNotExists()
	for _, definition := range columnDefinitions {
		builder = builder.Define(definition...)
	}
	tx := m.dbw.MustBegin(ctx, &sql.TxOptions{})
	// create user auth table
	q, _ := builder.BuildWithFlavor(flavor)
	_, err = tx.Exec(q)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// execute index creation on user auth table if not exists
	for _, definition := range indexDefinitions {
		_, err = tx.Exec(definition)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func NewUserAuthManager(dbm db.Manager, opts ...common.Option[userManager]) UserManager {
	m := &userManager{
		dbw: dbm, dbr: dbm, driver: rules.WhenTrueR1(dbm == nil, func() db.SupportedDriver {
			return db.DriverPostgreSQL
		}, func() db.SupportedDriver {
			return dbm.Driver()
		}),
	}
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}
