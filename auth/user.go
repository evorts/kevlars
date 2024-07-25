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
	"github.com/evorts/kevlars/audit"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/inmemory"
	"github.com/evorts/kevlars/jwe"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/utils"
	"github.com/huandu/go-sqlbuilder"
	"strconv"
	"strings"
)

type UserManager interface {
	Add(ctx context.Context, records ...UserAuthRecord) error
	Save(ctx context.Context, record UserAuthRecord) (UserAuthRecord, error)
	RemoveByIds(ctx context.Context, ids ...int) error
	RemoveByUserIds(ctx context.Context, userIds ...int) error

	GetByUserIds(ctx context.Context, userIds ...int) (UserAuthRecords, error)

	AddAccess(ctx context.Context, records ...UserAccessRecord) error
	DisabledAccessByIds(ctx context.Context, ids ...int) error

	// IsAllowed user id to access resources
	IsAllowed(ctx context.Context, id int64, resource string, scope Scope) (bool, error)

	Introspect(ctx context.Context, token string, bindTo interface{}) error
	Authenticate(ctx context.Context, id int64, creds string) (token string, err error)

	common.Init[UserManager]
}

type userManager struct {
	dbw    db.Manager
	dbr    db.Manager
	im     inmemory.Manager
	log    logger.Manager
	driver db.SupportedDriver
	audit  audit.Manager
	jwe    jwe.Manager
}

const (
	tableUserAuth   = "user_auth"
	tableUserAccess = "user_access"

	inMemoryUserCredsHashKey    = "user_creds"    // user_id -> creds
	inMemoryUserTokenHashKey    = "user_token"    // user_token -> detail/claim
	inMemoryUserDisabledHashKey = "user_disabled" // user_id -> disabled state
)

//goland:noinspection SqlResolve
var (
	userCustomDefinitions = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf(`CREATE TYPE access_scope AS ENUM('read', 'write', 'delete', 'undefined')`),
		},
	}
	userAuthTableExistenceCheckQuery = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf(`select count(table_name) as tableCount from information_schema.tables ist
				       where ist.table_name in ('%s,%s')`, tableUserAuth, tableUserAccess),
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
	userAccessColumnsDefinition = map[db.SupportedDriver][][]string{
		db.DriverPostgreSQL: {
			{"id", "serial", "primary key"},
			{"user_id", "bigint", "not null"},
			{"resource", "varchar(255)", "not null"},
			{"scopes", "access_scope[]", "default '[]'::access_scope[]"},
			{"disabled", "boolean", "default false"},
			{"created_at", "timestamp with time zone", "default current_timestamp"},
			{"updated_at", "timestamp with time zone"},
			{"disabled_at", "timestamp with time zone"},
		},
	}
	userAccessIndexDefinition = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf("create index if not exists %s_user_id_resource_uidx on public.%s(user_id, resource)", tableUserAccess, tableUserAccess),
			fmt.Sprintf("create index if not exists %s_disabled_idx on public.%s(disabled)", tableUserAccess, tableUserAccess),
			fmt.Sprintf("create index if not exists %s_created_at_idx on public.%s(created_at)", tableUserAccess, tableUserAccess),
		},
	}
	userAuthSaveQuery = map[db.SupportedDriver]string{
		db.DriverPostgreSQL: fmt.Sprintf(
			`INSERT INTO %s (id, user_id, creds, disabled)
					VALUES(:id, :user_id, :creds, :disabled)
				ON CONFLICT (user_id) DO 
					UPDATE SET
						CASE WHEN disabled <> excluded.disabled
							THEN disabled = excluded.disabled
						END,
						creds = COALESCE(NULLIF(excluded.creds,''),creds),
						CASE WHEN disabled 
							THEN disabled_at = current_timestamp
						END
				RETURNING id, user_id, disabled
		`, tableUserAuth),
	}
	userAccessSaveQuery = map[db.SupportedDriver]string{
		db.DriverPostgreSQL: fmt.Sprintf(
			`INSERT INTO %s (id, user_id, resource, scopes, disabled, disabled_at)
					VALUES(:id, :user_id, :creds, :disabled, CASE WHEN :disabled THEN current_timestamp ELSE NULL END)
				ON CONFLICT (user_id) DO 
					UPDATE SET
						CASE WHEN disabled <> excluded.disabled
							THEN disabled = excluded.disabled
						END,
						creds = COALESCE(NULLIF(excluded.creds,''),creds),
						CASE WHEN disabled 
							THEN disabled_at = current_timestamp
						END
				RETURNING id, user_id, disabled
		`, tableUserAccess),
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
	q, ok := userAuthSaveQuery[m.driver]
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
	if len(tableChecks) < 2 {
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
	// create custom type definitions
	if !utils.KeyExistsInMap(userCustomDefinitions, driver) {
		return errors.New("driver not supported by custom definition")
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
	accessColumnsDefinition := utils.GetValueOnMap(userAccessColumnsDefinition, driver, [][]string{})
	if len(accessColumnsDefinition) < 1 {
		return errors.New("user access column definitions is empty")
	}
	accessIndexDefinition := utils.GetValueOnMap(userAccessIndexDefinition, driver, []string{})
	if len(accessIndexDefinition) < 1 {
		return errors.New("user access index definitions is empty")
	}
	builder := sqlbuilder.NewCreateTableBuilder().CreateTable(tableUserAuth).IfNotExists()
	for _, definition := range columnDefinitions {
		builder = builder.Define(definition...)
	}
	tx := m.dbw.MustBegin(ctx, &sql.TxOptions{})
	// create custom type
	for _, definition := range userCustomDefinitions[driver] {
		_, err = tx.Exec(definition)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
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
	builderAccess := sqlbuilder.NewCreateTableBuilder().CreateTable(tableUserAccess).IfNotExists()
	for _, definition := range accessColumnsDefinition {
		builderAccess = builderAccess.Define(definition...)
	}
	// create user access table
	q, _ = builderAccess.BuildWithFlavor(flavor)
	_, err = tx.Exec(q)
	if err != nil {
		_ = tx.Rollback()
	}
	// execute index creation on user access table if not exists
	for _, definition := range accessIndexDefinition {
		_, err = tx.Exec(definition)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (m *userManager) AddAccess(ctx context.Context, records ...UserAccessRecord) error {
	// build select values
	svp := make([]string, len(records))
	svArgs := make([]interface{}, 0)
	for _, record := range records {
		svp = append(svp, fmt.Sprintf(
			`(%s, CASE WHEN %s THEN current_timestamp ELSE NULL END)`,
			strings.Join(utils.RepeatInSlice("?", 5), ","),
			strconv.FormatBool(record.Disabled),
		))
		svArgs = append(svArgs, record.ID, record.UserID, record.Resource, record.Scopes, record.Disabled)
	}
	q := `INSERT INTO ` + tableUserAccess + `(id, user_id, resource, scopes, disabled)
		VALUES ` + strings.Join(svp, ",") + `
		ON CONFLICT(user_id, resource) DO
			UPDATE SET
				scopes = 
				disabled = (CASE WHEN disabled <> excluded.disabled THEN excluded.disabled ELSE disabled END),
				updated_at = current_timestamp,
				disabled_at = (CASE WHEN excluded.disabled THEN current_timestamp ELSE disabled_at END)
		RETURNING id, user_id, resource
	`
	_, err := m.dbw.Exec(ctx, m.dbw.Rebind(q), svArgs...)
	return err
}

func (m *userManager) IsAllowed(ctx context.Context, id int64, resource string, scope Scope) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *userManager) Introspect(ctx context.Context, token string, bindTo interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *userManager) Authenticate(ctx context.Context, id int64, creds string) (token string, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *userManager) DisabledAccessByIds(ctx context.Context, ids ...int) error {
	//TODO implement me
	panic("implement me")
}

func NewUserAuthManager(dbm db.Manager, opts ...common.Option[userManager]) UserManager {
	m := &userManager{
		dbw: dbm, dbr: dbm, driver: rules.WhenTrueRE1(dbm == nil, func() db.SupportedDriver {
			return db.DriverPostgreSQL
		}, func() db.SupportedDriver {
			return dbm.Driver()
		}),
		audit: audit.NewNoop(),
		im:    inmemory.NewNoop(),
	}
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}
