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
	"database/sql"
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/ctime"
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
	ts.cm = NewClientManager(ts.db).MustInit()
}

func (ts *clientAuthTestSuite) TestInstantiation() {
	type args struct {
		opts          []common.Option[clientManager]
		migrationFile fs.FS
		db            db.Manager
	}
	tests := []struct {
		name      string
		args      args
		wantErr   error
		wantPanic bool
	}{
		{
			name: "instantiate with invalid db manager, should panic",
			args: args{
				opts: []common.Option[clientManager]{},
			},
			wantPanic: true,
		},
		{
			name: "instantiate but only init schema, should pass",
			args: args{
				opts: []common.Option[clientManager]{},
				db:   ts.db,
			},
			wantErr:   nil,
			wantPanic: false,
		},
		{
			name: "instantiate with data migration and load data into memory, should pass",
			args: args{
				opts: []common.Option[clientManager]{
					ClientWithExecuteMigration(true, "./sample_migration"),
				},
				db: ts.db,
			},
			wantErr:   nil,
			wantPanic: false,
		},
	}
	for _, tc := range tests {
		ts.Run(tc.name, func() {
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
	type args struct {
		items Clients
	}
	tests := []struct {
		name               string
		args               args
		shouldProduceError bool
		expectResultCount  int
	}{
		{
			name:               "given empty items should produce error",
			args:               args{items: nil},
			shouldProduceError: true,
			expectResultCount:  0,
		},
		{
			name: "given valid single items should success added to table",
			args: args{items: Clients{
				{
					Name:      "Add Client 1",
					Secret:    "Add Client Secret 1",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(10 * time.Hour), Valid: true},
					Disabled:  false,
				},
			}},
			shouldProduceError: false,
			expectResultCount:  1,
		},
		{
			name: "given valid multiple of new items should success save to table",
			args: args{items: Clients{
				{
					Name:      "Add Client 2",
					Secret:    "Add Client Secret 2",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(12 * time.Hour), Valid: true},
					Disabled:  true,
				},
				{
					Name:      "Add Client 3",
					Secret:    "Add Client Secret 3",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(11 * time.Hour), Valid: true},
					Disabled:  false,
				},
				{
					Name:      "Add Client 4",
					Secret:    "Add Client Secret 4",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(11 * time.Hour), Valid: true},
					Disabled:  false,
				},
				{
					Name:      "Add Client 5",
					Secret:    "Add Client Secret 5",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(11 * time.Hour), Valid: true},
					Disabled:  true,
				},
			}},
			shouldProduceError: false,
			expectResultCount:  4,
		},
	}
	for _, tc := range tests {
		ts.Run(tc.name, func() {
			rc, err := ts.cm.AddClient(ts.ctx, tc.args.items)
			if tc.shouldProduceError {
				assert.Error(ts.T(), err)
			} else {
				assert.NoError(ts.T(), err)
			}
			assert.Equal(ts.T(), tc.expectResultCount, len(rc))
		})
	}
}

func (ts *clientAuthTestSuite) TestAddScopes() {
	type args struct {
		items ClientScopes
	}
	tests := []struct {
		name               string
		args               args
		dependsOn          func(ctx context.Context, scopes ClientScopes)
		shouldProduceError bool
		expectResultCount  int
	}{
		{
			name:               "given empty items should produce error",
			args:               args{items: nil},
			shouldProduceError: true,
			dependsOn: func(ctx context.Context, scopes ClientScopes) {
				// do nothing
			},
		},
		{
			name: "given valid single items should success added to table",
			args: args{items: ClientScopes{
				{
					Resource: "/path/to/resource",
					Scopes:   Scopes{ScopeRead, ScopeWrite},
					Disabled: false,
				},
			}},
			shouldProduceError: false,
			dependsOn: func(ctx context.Context, scopes ClientScopes) {
				rs, _ := ts.cm.AddClient(ctx, Clients{{
					Name:      "Add Client for Scope 1",
					Secret:    "Add Client Secret for Scope 1",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(10 * time.Hour), Valid: true},
					Disabled:  false,
				}})
				for _, v := range scopes {
					v.ClientID = rs[0].ID
				}
			},
			expectResultCount: 1,
		},
		{
			name: "given multiple new items should success save to table",
			args: args{items: ClientScopes{
				{
					Resource: "/path/to/resource/updated",
					Scopes:   Scopes{ScopeRead, ScopeDelete},
					Disabled: false,
				},
				{
					Resource: "/path/to/resource/read",
					Scopes:   Scopes{ScopeRead},
					Disabled: false,
				},
			}},
			shouldProduceError: false,
			dependsOn: func(ctx context.Context, scopes ClientScopes) {
				rs, _ := ts.cm.AddClient(ctx, Clients{{
					Name:      "Add Client for Scope 2",
					Secret:    "Add Client Secret for Scope 2",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(10 * time.Hour), Valid: true},
					Disabled:  false,
				}})
				for _, scope := range scopes {
					scope.ClientID = rs[0].ID
				}
			},
			expectResultCount: 2,
		},
	}
	for _, tc := range tests {
		ts.Run(tc.name, func() {
			tc.dependsOn(ts.ctx, tc.args.items)
			rc, err := ts.cm.AddClientScope(ts.ctx, tc.args.items)
			if tc.shouldProduceError {
				assert.Error(ts.T(), err)
			} else {
				assert.NoError(ts.T(), err)
			}
			assert.Equal(ts.T(), tc.expectResultCount, len(rc))
		})
	}
}

func (ts *clientAuthTestSuite) TestAddClientWithScopes() {
	type args struct {
		items ClientWithScopes
	}
	tests := []struct {
		name               string
		args               args
		shouldProduceError bool
	}{
		{
			name: "given client with empty scopes should only create client data",
			args: args{
				items: ClientWithScopes{
					Client: &Client{
						Name:      "Add Client With Scope: Client 1",
						Secret:    "Add Client With Scope: Secret 1",
						ExpiredAt: sql.NullTime{Time: ctime.NowAdd(15 * time.Hour), Valid: true},
						Disabled:  false,
					},
					Scopes: nil,
				}},
			shouldProduceError: false,
		},
		{
			name: "given valid single items should success added to table",
			args: args{items: ClientWithScopes{
				Client: &Client{
					Name:      "Add Client With Scope: Client 2",
					Secret:    "Add Client With Scope: Secret 2",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(14 * time.Hour), Valid: true},
					Disabled:  true,
				},
				Scopes: ClientScopes{
					{
						Resource: "/path/to/resource",
						Scopes:   Scopes{ScopeRead, ScopeWrite},
						Disabled: false,
					},
				},
			}},
			shouldProduceError: false,
		},
		{
			name: "given multiple scope items should success save to table",
			args: args{items: ClientWithScopes{
				Client: &Client{
					Name:      "Add Client With Scope: Client 3",
					Secret:    "Add Client With Scope: Secret 3",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(17 * time.Hour), Valid: true},
					Disabled:  false,
				},
				Scopes: ClientScopes{
					{
						Resource: "/path/to/resource/updated",
						Scopes:   Scopes{ScopeRead, ScopeDelete},
						Disabled: false,
					},
					{
						Resource: "/path/to/resource/read",
						Scopes:   Scopes{ScopeRead},
						Disabled: true,
					},
				},
			}},
			shouldProduceError: false,
		},
	}
	for _, tc := range tests {
		ts.Run(tc.name, func() {
			_, err := ts.cm.AddClientWithScopes(ts.ctx, tc.args.items)
			assert.NoError(ts.T(), err)
		})
	}
}

func (ts *clientAuthTestSuite) TestModifyClient() {
	type args struct {
		item *Client
	}
	tests := []struct {
		name               string
		args               args
		dependOn           func(item *Client)
		shouldProduceError bool
	}{
		{
			name: "given invalid item (missing id) should produce an error",
			args: args{
				item: &Client{
					Name:      "Modify Client: Client 1",
					Secret:    "Modify Client: Secret 1",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(15 * time.Hour), Valid: true},
				},
			},
			shouldProduceError: true,
			dependOn:           func(item *Client) {},
		},
		{
			name: "given valid item should modified successfully",
			args: args{
				item: &Client{
					ID:        1,
					Name:      "Modify Client: Client 1",
					Secret:    "Modify Client: Secret 1",
					ExpiredAt: sql.NullTime{Time: ctime.NowAdd(15 * time.Hour), Valid: true},
				},
			},
			shouldProduceError: false,
			dependOn: func(item *Client) {
				rs, _ := ts.cm.AddClient(ts.ctx, Clients{{
					ID:        1,
					Name:      "New Client 1",
					Secret:    "New Secret 1",
					ExpiredAt: sql.NullTime{},
				}})
				item.ID = rs[0].ID
			},
		},
		{
			name: "given valid item without expired date should modified successfully",
			args: args{
				item: &Client{
					ID:        2,
					Name:      "Modify Client: Client 2",
					Secret:    "Modify Client: Secret 2",
					ExpiredAt: sql.NullTime{},
				},
			},
			shouldProduceError: false,
			dependOn: func(item *Client) {
				rs, _ := ts.cm.AddClient(ts.ctx, Clients{{
					ID:        2,
					Name:      "New Client 2",
					Secret:    "New Secret 2",
					ExpiredAt: sql.NullTime{},
				}})
				item.ID = rs[0].ID
			},
		},
	}
	for _, tc := range tests {
		ts.Run(tc.name, func() {
			tc.dependOn(tc.args.item)
			err := ts.cm.ModifyClient(ts.ctx, *tc.args.item)
			if tc.shouldProduceError {
				assert.Error(ts.T(), err)
			} else {
				assert.NoError(ts.T(), err)
			}
		})
	}
}

func (ts *clientAuthTestSuite) TestModifyClientScope() {
	type args struct {
		item *ClientScope
	}
	tests := []struct {
		name               string
		args               args
		dependOn           func(item *ClientScope)
		shouldProduceError bool
	}{
		{
			name: "given invalid item (missing id) should produce an error",
			args: args{
				item: &ClientScope{
					ClientID: 1,
					Resource: "resource 1",
					Scopes:   nil,
				},
			},
			shouldProduceError: true,
			dependOn:           func(item *ClientScope) {},
		},
		{
			name: "given valid item should modified successfully",
			args: args{
				item: &ClientScope{
					ClientID: 1,
					Resource: "resource 1",
					Scopes:   nil,
				},
			},
			shouldProduceError: false,
			dependOn: func(item *ClientScope) {
				rs, _ := ts.cm.AddClient(ts.ctx, Clients{{
					ID:        1,
					Name:      "New Client 1",
					Secret:    "New Secret 1",
					ExpiredAt: sql.NullTime{},
				}})
				item.ClientID = rs[0].ID
				rss, _ := ts.cm.AddClientScope(ts.ctx, ClientScopes{{
					ClientID: rs[0].ID,
					Resource: "/path/to/resource",
					Scopes:   Scopes{ScopeRead, ScopeWrite},
				}})
				item.ID = rss[0].ID
			},
		},
		{
			name: "given valid item with scopes should modified successfully",
			args: args{
				item: &ClientScope{
					ClientID: 1,
					Resource: "resource 1",
					Scopes:   Scopes{ScopeWrite},
				},
			},
			shouldProduceError: false,
			dependOn: func(item *ClientScope) {
				rs, _ := ts.cm.AddClient(ts.ctx, Clients{{
					ID:        2,
					Name:      "New Client 2",
					Secret:    "New Secret 2",
					ExpiredAt: sql.NullTime{},
				}})
				item.ClientID = rs[0].ID
				rss, _ := ts.cm.AddClientScope(ts.ctx, ClientScopes{{
					ClientID: rs[0].ID,
					Resource: "/path/to/resource",
					Scopes:   Scopes{ScopeRead, ScopeWrite},
				}})
				item.ID = rss[0].ID
			},
		},
	}
	for _, tc := range tests {
		ts.Run(tc.name, func() {
			tc.dependOn(tc.args.item)
			err := ts.cm.ModifyClientScope(ts.ctx, *tc.args.item)
			if tc.shouldProduceError {
				assert.Error(ts.T(), err)
			} else {
				assert.NoError(ts.T(), err)
			}
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
