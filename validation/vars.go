/**
 * @Author: steven
 * @Description:
 * @File: vars
 * @Date: 15/11/23 13.56
 */

package validation

import "regexp"

var (
	dateFormatRegex              = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	alphabet                     = regexp.MustCompile("[a-zA-Z]+")
	alphaNumeric                 = regexp.MustCompile("[\\w]+")
	alphaNumericWithSpace        = regexp.MustCompile("[\\w\\s]+")
	alphaNumericWithDashAndSpace = regexp.MustCompile("[\\w\\s-]+")
)
