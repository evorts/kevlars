/**
 * @Author: steven
 * @Description:
 * @File: app_audit
 * @Date: 13/05/24 18.28
 */

package scaffold

import (
	"github.com/evorts/kevlars/common"
	"github.com/evorts/kevlars/crypt"
)

type ISecurity interface {
	WithCrypts() IApplication
	WithHasher() IApplication

	HasCrypts() bool
	HasDefaultCrypts() bool
	Crypt(key string) crypt.Manager
	DefaultCrypt() crypt.Manager

	Hasher() crypt.Hasher
}

func (app *Application) WithCrypts() IApplication {
	if app.HasCrypts() {
		return app
	}
	crypts := app.Config().GetStringMap("crypt")
	if len(crypts) == 0 {
		panic("No crypt configured")
	}
	for k, vc := range crypts {
		item, ok := vc.(map[string]interface{})
		if !ok {
			continue
		}
		cipher, key, iv, pk := "", "", "", ""
		if v, exist := item["cipher"]; exist {
			cipher = v.(string)
		}
		if v, exist := item["key"]; exist {
			key = v.(string)
		}
		if v, exist := item["iv"]; exist {
			iv = v.(string)
		}
		if v, exist := item["pk"]; exist {
			pk = v.(string)
		}
		app.cryptMap[k] = crypt.New()
		if len(cipher) > 0 {
			app.cryptMap[k].AddOptions(crypt.WithCipher(crypt.Cipher(cipher)))
		}
		if len(key) > 0 {
			app.cryptMap[k].AddOptions(crypt.WithKey(common.Bytes(key)))
		}
		if len(iv) > 0 {
			app.cryptMap[k].AddOptions(crypt.WithIV(common.Bytes(iv)))
		}
		if len(pk) > 0 {
			app.cryptMap[k].AddOptions(crypt.WithPrivateKey(pk))
		}
		app.cryptMap[k].MustInit()
	}
	return app
}

func (app *Application) WithHasher() IApplication {
	app.hasher = crypt.NewHasher()
	return app
}

func (app *Application) HasCrypts() bool {
	return app.cryptMap != nil && len(app.cryptMap) > 0
}
func (app *Application) HasDefaultCrypts() bool {
	return app.HasCrypts() && app.DefaultCrypt() != nil
}
func (app *Application) Crypt(key string) crypt.Manager {
	if !app.HasCrypts() {
		return nil
	}
	if v, ok := app.cryptMap[key]; ok {
		return v
	}
	return nil
}

func (app *Application) DefaultCrypt() crypt.Manager {
	return app.Crypt(DefaultKey)
}
