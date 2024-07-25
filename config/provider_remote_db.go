/**
 * @Author: steven
 * @Description:
 * @File: provider_db
 * @Date: 21/07/24 05.36
 */

package config

import "github.com/evorts/kevlars/common"

type dbProvider struct {
	tableName     string
	contextPrefix []string
	data          map[string]interface{}
}

func (d *dbProvider) Init() error {
	//TODO implement me
	panic("implement me")
}

func (d *dbProvider) GetData() map[string]interface{} {
	return d.data
}

func NewRemoteDB(contextPrefix []string, opts ...common.Option[dbProvider]) Provider {
	p := &dbProvider{
		tableName:     defaultTableName,
		contextPrefix: contextPrefix,
		data:          map[string]interface{}{},
	}
	for _, opt := range opts {
		opt.Apply(p)
	}
	return p
}
