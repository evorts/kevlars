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
	"github.com/evorts/kevlars/db"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/ts"
	"github.com/evorts/kevlars/utils"
	"github.com/huandu/go-sqlbuilder"
	"net/url"
)

type ClientManager interface {
	AddClients(ctx context.Context, items Clients) error
	AddScopes(ctx context.Context, items ClientScopes) error
	AddClientsWithScopes(ctx context.Context, items ClientWithScopes) error

	GetClientsWithScopesBy(ctx context.Context, by db.IHelper) (ClientsWithScopes, error)

	VoidClientsByIds(ctx context.Context, ids ...int) error
	VoidScopeByIds(ctx context.Context, id ...int) error

	IsAllowed(secret, resource, scope string) (clientName string, allowed bool)

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
	tableExistenceCheckQuery = map[db.SupportedDriver][]string{
		db.DriverPostgreSQL: {
			fmt.Sprintf(`select count(table_name) as tableCount from information_schema.tables ist
				       where ist.table_name in ('%s','%s')`, tableClients, tableClientScope),
		},
	}
	clientsColumnsDefinition = map[db.SupportedDriver][][]string{
		db.DriverPostgreSQL: {
			{"id", "serial", "primary key"},
			{"name", "varchar(25)", "not null"},
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
			{"resource", "varchar(150)", "not null"},
			{"scopes", "jsonb", "default '[]'::jsonb"},
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

func (m *clientManager) AddClients(ctx context.Context, items Clients) error {
	//TODO implement me
	panic("implement me")
}

func (m *clientManager) AddScopes(ctx context.Context, items ClientScopes) error {
	//TODO implement me
	panic("implement me")
}

func (m *clientManager) AddClientsWithScopes(ctx context.Context, items ClientWithScopes) error {
	//TODO implement me
	panic("implement me")
}

func (m *clientManager) GetClientsWithScopesBy(ctx context.Context, by db.IHelper) (ClientsWithScopes, error) {
	qf, args := by.BuildSqlAndArgsWithWherePrefix()
	//goland:noinspection SqlResolve
	q := fmt.Sprintf(`
		select 
		    c.id, c.name, c.secret, c.expired_at, c.disabled, 
		    c.created_at, c.updated_at, c.disabled_at,
		    cs.id as scope_id, cs.resource, cs.scopes, cs.disabled as scope_disabled,
		    cs.created_at as scope_created_at, cs.updated_at as scope_updated_at,
		    cs.disabled_at as scope_disabled_at
		from %s c 
		join %s cs on cs.client_id = c.id
		%s`, tableClients, tableClientScope, qf)
	rs := make(ClientsWithScopes, 0)
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

func (m *clientManager) IsAllowed(secret, resource, scope string) (clientName string, allowed bool) {
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
	if dt.ExpiredAt != nil && ts.Now().Before(*dt.ExpiredAt) {
		return dt.ClientName, false
	}
	return dt.ClientName, dt.Scopes.AllowedTo(Scope(scope))
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

func (m *clientManager) getFlavorByDriver(driver db.SupportedDriver) sqlbuilder.Flavor {
	if driver == db.DriverPostgreSQL {
		return sqlbuilder.PostgreSQL
	}
	panic("unsupported flavor")
}

func (m *clientManager) tableCheck(ctx context.Context, driver db.SupportedDriver) (int, error) {
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
	// to avoid unnecessary execution of schema scaffolding,
	// check the existence of tables -- should return total table of 2
	total, err := m.tableCheck(ctx, driver)
	if err != nil {
		return err
	}
	if total == 2 {
		return nil
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
	// create clients table
	q := builderClients.String()
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
	_, err = tx.Exec(builderClientScopes.String())
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

func NewClientManager(db db.Manager, opts ...common.Option[clientManager]) ClientManager {
	m := &clientManager{
		dbw:              db,
		dbr:              db,
		log:              logger.NewNoop(),
		mapAuthorization: make(mapClientAuthorization),
		migrationEnabled: false,
		migrationDir:     []string{},
	}
	m.AddOptions(opts...)
	return m
}
