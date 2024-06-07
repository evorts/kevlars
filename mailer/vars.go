/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 28/05/24 17.05
 */

package mailer

import "github.com/evorts/kevlars/utils"

type Provider string

const (
	ProviderGMail Provider = "gmail"
)

var (
	validProviders = []Provider{ProviderGMail}
)

func (p Provider) String() string {
	return string(p)
}

func ValidProvider(v Provider) bool {
	return utils.InArray(validProviders, v)
}
