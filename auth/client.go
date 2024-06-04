/**
 * @Author: steven
 * @Description:
 * @File: consumers
 * @Date: 13/05/24 16.14
 */

package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/ctime"
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/rules/eval"
	"github.com/evorts/kevlars/utils"
	"github.com/huandu/go-sqlbuilder"
	"github.com/lib/pq"
	"net/url"
	"strings"
)

type ClientManager interface {
	AddClient(ctx context.Context, items Clients) (Clients, error)
	AddScope(ctx context.Context, items ClientScopes) (ClientScopes, error)
	AddClientWithScopes(ctx context.Context, item ClientWithScopes) (*ClientWithScopes, error)

	GetClientsWithScopesBy(ctx context.Context, by db.IHelper) (ClientsWithScopes, error)

	VoidClientsByIds(ctx context.Context, ids ...int) error
	VoidScopeByIds(ctx context.Context, id ...int) error

	ModifyClient(ctx context.Context, items Clients) error

	IsAllowed(secret, resource string, scope Scope) (clientName string, allowed bool)

	Init() error
	MustInit() ClientManager
	AddOptions(opts ...common.Option[clientManager]) ClientManager

	migrate()
	loadData() error
	dbMigrate() *dbmate.DB
}

type clientManager struct {
	dbw              db.Manager
	dbr              db.Manager
	dbm              *dbmate.DB
	driver           db.SupportedDriver
	log              logger.Manager
	migrationDir     []string
	migrationEnabled bool
	lazyLoad         bool
	mapAuthorization mapClientAuthorization
}

const (
	tableClients     = "clients"
	tableClientScope = "client_scopes"
)

//goland:noinspection SqlCurrentSchemaInspection,SqlResolve
var (
	clientCustomDefinitions = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf(`CREATE TYPE client_scope AS ENUM('read', 'write', 'delete', 'undefined')`),
		},
	}
	clientTableExistenceCheckQuery = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf(`select count(table_name) as tableCount from information_schema.tables ist
				       where ist.table_name in ('%s','%s')`, tableClients, tableClientScope),
		},
	}
	clientsColumnsDefinition = map[db.SupportedDriver][][]string{
		db.DriverPostgreSQL: {
			{"id", "serial", "primary key"},
			{"name", "varchar(45)", "not null"},
			{"secret", "varchar(128)", "not null"},
			{"expired_at", "timestamp with time zone"},
			{"disabled", "boolean", "default false"},
			{"created_at", "timestamp with time zone", "default current_timestamp"},
			{"updated_at", "timestamp with time zone"},
			{"disabled_at", "timestamp with time zone"},
		},
	}
	clientsTableIndexDefinition = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf("create unique index if not exists %s_secret_uidx on public.%s(secret)", tableClients, tableClients),
			fmt.Sprintf("create unique index if not exists %s_name_uidx on public.%s(name)", tableClients, tableClients),
		},
	}
	clientScopesColumnsDefinition = map[db.SupportedDriver][][]string{
		db.DriverPostgreSQL: {
			{"id", "serial", "primary key"},
			{"client_id", "int"},
			{"constraint", "fk_" + tableClientScope + "_client_id", "foreign key (client_id)", "references " + tableClients + "(id)"},
			{"resource", "varchar(255)", "not null"},
			{"scopes", "client_scope[]", "default array[]::client_scope[]"},
			{"disabled", "boolean", "default false"},
			{"created_at", "timestamp with time zone", "default current_timestamp"},
			{"updated_at", "timestamp with time zone"},
			{"disabled_at", "timestamp with time zone"},
		},
	}
	clientScopesTableIndexDefinition = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf("create unique index %s_client_id_resource_uidx on %s(client_id, resource)", tableClientScope, tableClientScope),
		},
	}
)

var (
	// addClientQuery
	addClientQuery = map[db.SupportedDriver]struct {
		placeholder func(int) []string
		query       func(v string) string
	}{
		db.DriverPostgreSQL: {
			placeholder: func(repeat int) []string {
				return db.PlaceholderRepeat(
					fmt.Sprintf(
						`(%s,CASE WHEN ? THEN current_timestamp END)`,
						strings.Join(utils.RepeatInSlice("?", 4), ","),
					), repeat,
				)
			},
			query: func(v string) string {
				return `INSERT INTO ` + tableClients + `(name, secret, disabled, expired_at, disabled_at)
					VALUES ` + v + `
					RETURNING id, name, disabled, expired_at, created_at, disabled_at
		`
			},
		},
	}
	// addScopeQuery
	addScopeQuery = map[db.SupportedDriver]struct {
		placeholder func(int) []string
		query       func(v string) string
	}{
		db.DriverPostgreSQL: {
			placeholder: func(repeat int) []string {
				return db.PlaceholderRepeat(
					fmt.Sprintf(
						`(%s,CASE WHEN ? THEN current_timestamp END)`,
						strings.Join(utils.RepeatInSlice("?", 4), ","),
					), repeat,
				)
			},
			query: func(v string) string {
				return `INSERT INTO ` + tableClientScope + `(client_id, resource, scopes, disabled, disabled_at)
						VALUES ` + v + `
						RETURNING id, client_id, resource, scopes, disabled, created_at, disabled_at
				`
			},
		},
	}

	getClientWithScopesByQuery = map[db.SupportedDriver]struct {
		query func(qf string) string
	}{
		db.DriverPostgreSQL: {
			query: func(qf string) string {
				return `SELECT 
		    				c.id, c.name, c.secret, c.expired_at, c.disabled, 
		    				c.created_at, c.updated_at, c.disabled_at,
		    				cs.id as scope_id, cs.resource, cs.scopes, cs.disabled as scope_disabled,
		    				cs.created_at as scope_created_at, cs.updated_at as scope_updated_at,
		    				cs.disabled_at as scope_disabled_at
						FROM` + " " + tableClients + ` c JOIN ` + tableClientScope +
					` cs ON cs.client_id = c.id ` + qf
			},
		},
	}
)

func (m *clientManager) ModifyClient(ctx context.Context, items Clients) error {
	//TODO implement me
	panic("implement me")
}

func (m *clientManager) AddClient(ctx context.Context, items Clients) (Clients, error) {
	rs := make(Clients, 0)
	if eval.IsEmpty(items) {
		return rs, db.ErrorEmptyArguments
	}
	placeholders := addClientQuery[m.driver].placeholder(len(items))
	args := make([]interface{}, 0)
	for _, item := range items {
		args = append(args, item.Name, item.Secret, item.Disabled, item.ExpiredAt, item.Disabled)
	}
	q := addClientQuery[m.driver].query(strings.Join(placeholders, ","))
	rows, err := m.dbw.Query(ctx, m.dbw.Rebind(q), args...)
	defer func() {
		rules.WhenTrue(rows != nil, func() {
			_ = rows.Close()
		})
	}()
	if err != nil {
		return rs, err
	}
	for rows.Next() {
		var item Client
		if err = rows.StructScan(&item); err != nil {
			return rs, err
		}
		rs = append(rs, &item)
	}
	return rs, nil
}

func (m *clientManager) AddScope(ctx context.Context, items ClientScopes) (ClientScopes, error) {
	rs := make(ClientScopes, 0)
	if eval.IsEmpty(items) {
		return rs, db.ErrorEmptyArguments
	}
	// build select values
	placeholders := addScopeQuery[m.driver].placeholder(len(items))
	args := make([]interface{}, 0)
	for _, item := range items {
		args = append(args, item.ClientID, item.Resource, pq.Array(item.Scopes), item.Disabled, item.Disabled)
	}
	q := addScopeQuery[m.driver].query(strings.Join(placeholders, ","))
	rows, err := m.dbw.Query(ctx, m.dbw.Rebind(q), args...)
	defer func() {
		rules.WhenTrue(rows != nil, func() {
			_ = rows.Close()
		})
	}()
	if err != nil {
		return rs, err
	}
	for rows.Next() {
		var item ClientScope
		if err = rows.StructScan(&item); err != nil {
			return rs, err
		}
		rs = append(rs, &item)
	}
	return rs, nil
}

func (m *clientManager) AddClientWithScopes(ctx context.Context, item ClientWithScopes) (*ClientWithScopes, error) {
	// start tx
	tx := m.dbw.MustBegin(ctx, nil)
	// add client
	ph := addClientQuery[m.driver].placeholder(1)
	q := addClientQuery[m.driver].query(strings.Join(ph, ","))
	var client Client
	err := tx.QueryRowx(tx.Rebind(q), item.Name, item.Secret, item.Disabled, item.ExpiredAt, item.Disabled).StructScan(&client)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	rs := &ClientWithScopes{Client: &client, Scopes: make(ClientScopes, 0)}
	// add scopes
	if item.Scopes == nil || len(item.Scopes) < 1 {
		return rs, tx.Commit()
	}
	ph = addScopeQuery[m.driver].placeholder(len(item.Scopes))
	q = addScopeQuery[m.driver].query(strings.Join(ph, ","))
	args := make([]interface{}, 0)
	for _, scope := range item.Scopes {
		args = append(args, client.ID, scope.Resource, pq.Array(scope.Scopes), scope.Disabled, scope.Disabled)
	}
	rows, errQ := tx.Queryx(tx.Rebind(q), args...)
	defer func() {
		rules.WhenTrue(rows != nil, func() {
			_ = rows.Close()
		})
	}()
	if errQ != nil {
		_ = tx.Rollback()
		return nil, errQ
	}
	for rows.Next() {
		var scope ClientScope
		if err = rows.StructScan(&scope); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		rs.Scopes = append(rs.Scopes, &scope)
	}
	return rs, tx.Commit()
}

func (m *clientManager) GetClientsWithScopesBy(ctx context.Context, by db.IHelper) (ClientsWithScopes, error) {
	qf, args := by.BuildSqlAndArgsWithWherePrefix()
	//goland:noinspection SqlResolve
	q := m.dbr.Rebind(getClientWithScopesByQuery[m.driver].query(qf))
	rs := make(ClientsWithScopes, 0)
	rows, err := m.dbr.Query(ctx, q, args...)
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
	rsMap := make(map[int]*ClientWithScopes)
	rsIdx := make([]int, 0)
	for rows.Next() {
		var (
			item      Client
			itemScope ClientScope
		)
		if err = rows.Scan(
			&item.ID, &item.Name, &item.Secret, &item.ExpiredAt, &item.Disabled,
			&item.CreatedAt, &item.UpdatedAt, &item.DisabledAt,
			&itemScope.ID, &itemScope.Resource, &itemScope.Scopes, &itemScope.Disabled,
			&itemScope.CreatedAt, &itemScope.UpdatedAt, &itemScope.DisabledAt,
		); err != nil {
			return rs, err
		}
		if _, ok := rsMap[item.ID]; !ok {
			rsMap[item.ID] = &ClientWithScopes{
				Client: &item,
				Scopes: make(ClientScopes, 0),
			}
			rsIdx = append(rsIdx, item.ID)
		}
		rsMap[item.ID].Scopes = append(rsMap[item.ID].Scopes, &itemScope)
	}
	// building result with origin sort from map
	for _, vId := range rsIdx {
		if v, ok := rsMap[vId]; ok {
			rs = append(rs, v)
		}
	}
	return rs, nil
}

func (m *clientManager) VoidClientsByIds(ctx context.Context, ids ...int) error {
	panic("implement me")
}

func (m *clientManager) VoidScopeByIds(ctx context.Context, ids ...int) error {
	//TODO implement me
	panic("implement me")
}

func (m *clientManager) IsAllowed(secret, resource string, scope Scope) (clientName string, allowed bool) {
	rules.WhenTrue(m.lazyLoad, func() {
		m.log.Info(m.loadData())
	})
	rm, ok := m.mapAuthorization[secret]
	if !ok {
		return "unknown", false
	}
	dt, okd := rm[resource]
	if !okd {
		return "unknown", false
	}
	if dt.Disabled {
		return dt.ClientName, false
	}
	if len(dt.Scopes) < 1 {
		return dt.ClientName, false
	}
	if dt.ExpiredAt != nil && ctime.Now().Before(*dt.ExpiredAt) {
		return dt.ClientName, false
	}
	return dt.ClientName, dt.Scopes.AllowedTo(scope)
}

func (m *clientManager) loadData() error {
	if len(m.mapAuthorization) > 0 {
		return nil
	}
	var (
		page  = 1
		limit = 20
	)
	for {
		items, err := m.GetClientsWithScopesBy(
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
			m.mapAuthorization[item.Secret] = make(map[string]clientDataForAuthorization)
			for _, scope := range item.Scopes {
				if scope == nil {
					continue
				}
				m.mapAuthorization[item.Secret][scope.Resource] = clientDataForAuthorization{
					ClientName: item.Name,
					Scopes:     scope.Scopes,
					Disabled:   rules.Iif(item.Disabled, item.Disabled, scope.Disabled),
					ExpiredAt:  item.ExpiredAt,
				}
			}
		}
		if len(items) < limit {
			break
		}
		page++
	}
	return nil
}

func (m *clientManager) migrate() {
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

func (m *clientManager) dbMigrate() *dbmate.DB {
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

func (m *clientManager) tableCheck(ctx context.Context, driver db.SupportedDriver) (int, error) {
	total := 0
	if !utils.KeyExistsInMap(clientTableExistenceCheckQuery, driver) {
		return total, errors.New("driver not supported by table existence check")
	}
	tableChecks := utils.GetValueOnMap(clientTableExistenceCheckQuery, driver, []string{})
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

func (m *clientManager) Init() error {
	ctx := context.Background()
	if err := m.initSchema(ctx); err != nil {
		return err
	}
	m.migrate()
	if !m.lazyLoad {
		if err := m.loadData(); err != nil {
			return err
		}
	}
	return nil
}

func (m *clientManager) initSchema(ctx context.Context) error {
	driver := m.dbw.Driver()
	flavor := getFlavorByDriver(driver)
	// to avoid unnecessary execution of schema scaffolding,
	// check the existence of tables -- should return total table of 2
	total, err := m.tableCheck(ctx, driver)
	if err != nil {
		return err
	}
	if total == 2 {
		return nil
	}
	// custom type
	if !utils.KeyExistsInMap(clientCustomDefinitions, driver) {
		return errors.New("driver not supported by custom definition")
	}
	if !utils.KeyExistsInMap(clientsColumnsDefinition, driver) {
		return errors.New("driver not supported by column definition")
	}
	clientsColumnDefinitions := utils.GetValueOnMap(clientsColumnsDefinition, driver, [][]string{})
	if len(clientsColumnDefinitions) < 1 {
		return errors.New("clients column definitions is empty")
	}
	clientsIndexDefinitions := utils.GetValueOnMap(clientsTableIndexDefinition, driver, []string{})
	if len(clientsIndexDefinitions) < 1 {
		return errors.New("clients index definitions is empty")
	}
	clientScopesColumnDefinitions := utils.GetValueOnMap(clientScopesColumnsDefinition, driver, [][]string{})
	if len(clientScopesColumnDefinitions) < 1 {
		return errors.New("client scopes column definitions is empty")
	}
	clientScopesIndexDefinitions := utils.GetValueOnMap(clientScopesTableIndexDefinition, driver, []string{})
	if len(clientScopesIndexDefinitions) < 1 {
		return errors.New("client scopes index definitions is empty")
	}
	builderClients := sqlbuilder.NewCreateTableBuilder().CreateTable(tableClients).IfNotExists()
	for _, definition := range clientsColumnDefinitions {
		builderClients = builderClients.Define(definition...)
	}
	builderClientScopes := sqlbuilder.NewCreateTableBuilder().CreateTable(tableClientScope).IfNotExists()
	for _, definition := range clientScopesColumnDefinitions {
		builderClientScopes = builderClientScopes.Define(definition...)
	}
	tx := m.dbw.MustBegin(ctx, &sql.TxOptions{})
	// create custom type

	for _, definition := range clientCustomDefinitions[driver] {
		_, err = tx.Exec(definition)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	// create clients table
	q, _ := builderClients.BuildWithFlavor(flavor)
	_, err = tx.Exec(q)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// execute index creation on clients table if not exists
	for _, definition := range clientsIndexDefinitions {
		_, err = tx.Exec(definition)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	// create client scopes table
	q, _ = builderClientScopes.BuildWithFlavor(flavor)
	_, err = tx.Exec(q)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// execute index creation on client scope table if not exists
	for _, definition := range clientScopesIndexDefinitions {
		_, err = tx.Exec(definition)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (m *clientManager) MustInit() ClientManager {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *clientManager) AddOptions(opts ...common.Option[clientManager]) ClientManager {
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}

func NewClientManager(dbm db.Manager, opts ...common.Option[clientManager]) ClientManager {
	m := &clientManager{
		dbw: dbm,
		dbr: dbm,
		driver: rules.WhenTrueR1(dbm == nil, func() db.SupportedDriver {
			return db.DriverPostgreSQL
		}, func() db.SupportedDriver {
			return dbm.Driver()
		}),
		log:              logger.NewNoop(),
		mapAuthorization: make(mapClientAuthorization),
		migrationEnabled: false,
		migrationDir:     []string{},
	}
	m.AddOptions(opts...)
	return m
}
