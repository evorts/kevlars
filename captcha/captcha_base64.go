/**
 * @Author: steven
 * @Description:
 * @File: captcha_base64
 * @Date: 15/01/24 21.48
 */

package captcha

import (
	"github.com/evorts/kevlars/common"
	"github.com/mojocn/base64Captcha"
)

type B64Type string

const (
	B64TypeDigit  B64Type = "digit"
	B64TypeString B64Type = "string"
)

const (
	defaultCaptchaWidth  = 170
	defaultCaptchaHeight = 50
	defaultCaptchaLength = 6

	defaultStringSource     = "1234567890qwertyuioplkjhgfdsazxcvbnm"
	defaultStringNoiseCount = 0
	defaultShowLineOptions  = base64Captcha.OptionShowSlimeLine
	defaultStringFont       = "RitaSmith.ttf"

	defaultDigitSkew     = 0.7
	defaultDigitDotCount = 70
)

type b64Manager struct {
	captchaType B64Type
	driver      base64Captcha.Driver
	store       base64Captcha.Store
	captcha     *base64Captcha.Captcha
}

func (m *b64Manager) Init() error {
	switch m.captchaType {
	case B64TypeDigit:
		m.driver = base64Captcha.NewDriverDigit(
			defaultCaptchaHeight,
			defaultCaptchaWidth,
			defaultCaptchaLength,
			defaultDigitSkew,
			defaultDigitDotCount,
		)
	default:
		m.driver = base64Captcha.NewDriverString(
			defaultCaptchaHeight,
			defaultCaptchaWidth,
			defaultStringNoiseCount,
			defaultShowLineOptions,
			defaultCaptchaLength,
			defaultStringSource,
			nil,
			base64Captcha.DefaultEmbeddedFonts,
			[]string{defaultStringFont},
		)
	}
	m.store = base64Captcha.DefaultMemStore
	m.captcha = base64Captcha.NewCaptcha(m.driver, m.store)
	return nil
}

func (m *b64Manager) MustInit() Manager[string] {
	if err := m.Init(); err != nil {
		panic(err)
	}
	return m
}

func (m *b64Manager) Generate() (id, answer, result string) {
	cid, b64, ans, err := m.captcha.Generate()
	if err != nil {
		return "", "", ""
	}
	return cid, ans, b64
}

func (m *b64Manager) Verify(id, value string, clear bool) bool {
	return m.store.Verify(id, value, clear)
}

func NewB64(opts ...common.Option[b64Manager]) Manager[string] {
	m := &b64Manager{
		captchaType: B64TypeString,
	}
	for _, opt := range opts {
		opt.Apply(m)
	}
	return m
}
