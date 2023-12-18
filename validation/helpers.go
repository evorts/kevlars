/**
 * @Author: steven
 * @Description:
 * @File: helpers
 * @Date: 15/11/23 13.50
 */

package validation

func IsStandardDateFormat(v string) bool {
	return dateFormatRegex.MatchString(v)
}

func IsAlphabet(v string) bool {
	return alphabet.MatchString(v)
}

func IsAlphaNumeric(v string) bool {
	return alphaNumeric.MatchString(v)
}

func IsAlphaNumericWithSpace(v string) bool {
	return alphaNumericWithSpace.MatchString(v)
}

func IsAlphaNumericWithDashAndSpace(v string) bool {
	return alphaNumericWithDashAndSpace.MatchString(v)
}
