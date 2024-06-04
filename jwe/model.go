/**
 * @Author: steven
 * @Description:
 * @File: model
 * @Date: 01/06/24 07.42
 */

package jwe

import (
	"github.com/mitchellh/mapstructure"
	"time"
)

const (
	ISSUER            = "evorts.com"
	DefaultExpiration = 1 * time.Hour //in hour
)

type Metadata map[string]any

func (m Metadata) ToStruct(bindTo any) error {
	return mapstructure.Decode(m, bindTo)
}

type Claim struct {
	Issuer   string   `json:"issuer"`
	ClientID string   `json:"client_id"`
	ID       int64    `json:"id"`
	Meta     Metadata `json:"meta"`
}

func NewClaim(clientID string, id int64, meta Metadata) *Claim {
	return &Claim{
		Issuer:   ISSUER,
		ClientID: clientID,
		ID:       id,
		Meta:     meta,
	}
}
