/*
go:build integration
 +build integration
*/

/**
 * @Author: steven
 * @Description:
 * @File: client_test
 * @Date: 13/05/24 21.23
 */

package auth

import (
	"context"
	"embed"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"io/fs"
	"log"
	"testing"
	"time"
)

type clientAuthTestSuite struct {
	suite.Suite

	ctx       context.Context
	db        db.Manager
	dbname    string
	port      int
	user      string
	pass      string
	dsn       string
	container *postgres.PostgresContainer
	cm        ClientManager
}

//go:embed db/migrations/20240515015140_client_test.sql
var dataTest embed.FS

func (ts *clientAuthTestSuite) SetupTest() {
	var err error
	ts.dbname = "test_db"
	ts.user = "user"
	ts.pass = "secrets"
	ts.port = 55432
	ts.ctx = context.Background()
	ts.container, err = postgres.RunContainer(
		ts.ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		testcontainers.WithHostPortAccess(ts.port),
		postgres.WithDatabase(ts.dbname),
		postgres.WithUsername(ts.user),
		postgres.WithPassword(ts.pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	ts.Require().NoError(err)
	ts.dsn, err = ts.container.ConnectionString(ts.ctx, "sslmode=disable")
	ts.Require().NoError(err)
	ts.db = db.New(db.DriverPostgreSQL, ts.dsn).MustConnect(ts.ctx)
	//ts.cm = NewClientManager(ts.db).MustInit()
}

func (ts *clientAuthTestSuite) TestInstantiation() {
	type args struct {
		opts          []common.Option[clientManager]
		migrationFile fs.FS
		db            db.Manager
	}
	tests := map[string]struct {
		args      args
		wantErr   error
		wantPanic bool
	}{
		"instantiate with invalid db manager, should panic": {
			args: args{
				opts: []common.Option[clientManager]{},
			},
			wantPanic: true,
		},
		"instantiate but only init schema, should pass": {
			args: args{
				opts: []common.Option[clientManager]{},
				db:   ts.db,
			},
			wantErr:   nil,
			wantPanic: false,
		},
		"instantiate with data migration and load data into memory, should pass": {
			args: args{
				opts: []common.Option[clientManager]{
					ClientWithExecuteMigration(dataTest, func() bool {
						return true
					}),
				},
				db: ts.db,
			},
			wantErr:   nil,
			wantPanic: false,
		},
	}
	for name, tc := range tests {
		ts.Run(name, func() {
			cm := NewClientManager(tc.args.db, tc.args.opts...)
			if tc.wantPanic {
				assert.Panics(ts.T(), func() {
					cm.MustInit()
				}, "must init not panic")
			} else {
				assert.Equal(ts.T(), tc.wantErr, cm.Init())
			}
		})
	}
}

func (ts *clientAuthTestSuite) TestAddClient() {
	tests := map[string]struct {
		wantErr error
	}{
		"test init schema": {
			wantErr: nil,
		},
	}
	for name, tc := range tests {
		ts.Run(name, func() {
			err := ts.cm.AddClients(ts.ctx, Clients{})
			assert.Equal(ts.T(), tc.wantErr, err)
		})
	}
}

func (ts *clientAuthTestSuite) TearDownTest() {
	if err := ts.container.Terminate(ts.ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func TestClientAuthTestSuite(t *testing.T) {
	suite.Run(t, new(clientAuthTestSuite))
}
