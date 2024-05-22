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
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

type userAuthTestSuite struct {
	suite.Suite

	ctx       context.Context
	db        db.Manager
	dbname    string
	port      int
	user      string
	pass      string
	dsn       string
	container *postgres.PostgresContainer
	um        userManager
}

func (ts *userAuthTestSuite) SetupTest() {
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
}

func (ts *userAuthTestSuite) TestInstantiation() {
	type args struct {
		opts []common.Option[userManager]
		db   db.Manager
	}
	tests := map[string]struct {
		args      args
		wantErr   error
		wantPanic bool
	}{
		"instantiate with invalid db manager, should panic": {
			args: args{
				opts: []common.Option[userManager]{},
			},
			wantPanic: true,
		},
		"instantiate but only init schema, should pass": {
			args: args{
				opts: []common.Option[userManager]{},
				db:   ts.db,
			},
			wantErr:   nil,
			wantPanic: false,
		},
	}
	for name, tc := range tests {
		ts.Run(name, func() {
			um := NewUserAuthManager(tc.args.db, tc.args.opts...)
			if tc.wantPanic {
				assert.Panics(ts.T(), func() {
					um.MustInit()
				}, "must init not panic")
			} else {
				assert.Equal(ts.T(), tc.wantErr, um.Init())
			}
		})
	}
}

/*
	func (ts *userAuthTestSuite) TestAddUser() {
		tests := map[string]struct {
			wantErr error
		}{
			"test init schema": {
				wantErr: nil,
			},
		}
		for name, tc := range tests {
			ts.Run(name, func() {
				err := ts.um.Add(ts.ctx, UserAuthRecord{})
				assert.Equal(ts.T(), tc.wantErr, err)
			})
		}
	}
*/
func (ts *userAuthTestSuite) TearDownTest() {
	if err := ts.container.Terminate(ts.ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func TestUserAuthTestSuite(t *testing.T) {
	suite.Run(t, new(userAuthTestSuite))
}
