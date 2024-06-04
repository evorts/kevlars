/**
 * @Author: steven
 * @Description:
 * @File: jwe_noop
 * @Date: 01/06/24 17.05
 */

package jwe

type jweNoop struct{}

func (j *jweNoop) Encode(v Claim) (token string, err error) {
	return "", nil
}

func (j *jweNoop) Decode(token string) (v Claim, err error) {
	return Claim{}, nil
}

func (j *jweNoop) Init() error {
	return nil
}

func (j *jweNoop) MustInit() Manager {
	return j
}

func NewNoop() Manager {
	return &jweNoop{}
}
