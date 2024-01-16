/**
 * @Author: steven
 * @Description:
 * @File: bot
 * @Date: 15/01/24 20.49
 */

package messaging

import (
	"github.com/evorts/kevlars/common"
)

type Sender interface {
	SendMessage(message string) error
	common.Init[Sender]
}
