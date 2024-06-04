/**
 * @Author: steven
 * @Description:
 * @File: audit_noop
 * @Date: 01/06/24 07.13
 */

package audit

import "context"

type noop struct{}

func (m *noop) Add(ctx context.Context, records ...Record) error {
	return nil
}

func (m *noop) Init() error {
	return nil
}

func (m *noop) MustInit() Manager {
	return m
}

func NewNoop() Manager {
	return &noop{}
}
