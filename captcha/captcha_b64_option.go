/**
 * @Author: steven
 * @Description:
 * @File: captcha_b64_option
 * @Date: 16/01/24 00.13
 */

package captcha

import "github.com/evorts/kevlars/common"

func B64WithType(v B64Type) common.Option[b64Manager] {
	return common.OptionFunc[b64Manager](func(m *b64Manager) {
		m.captchaType = v
	})
}
