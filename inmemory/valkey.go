/**
 * @Author: steven
 * @Description:
 * @File: valkey
 * @Date: 30/04/24 18.12
 */

package inmemory

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type valkeyManager struct {
	c valkey.Client

	addr string
	pwd  string
	db   int

	prefix string // act as prefix

	useTLS     bool
	certFile   string
	keyFile    string
	certB64    string
	keyB64     string
	serverName string

	tm    telemetry.Manager
	scope string
}

func (m *valkeyManager) MustConnect(ctx context.Context) Manager {
	if err := m.Connect(ctx); err != nil {
		panic(err)
	}
	return m
}

func (m *valkeyManager) Connect(ctx context.Context) error {
	var (
		tlsConfig *tls.Config
		err       error
		cert      tls.Certificate
	)
	if m.useTLS {
		// load cert when defined
		if len(m.certB64) > 0 {
			var cb, kb []byte
			cb, err = base64.StdEncoding.DecodeString(m.certB64)
			rules.WhenTrue(err == nil, func() {
				kb, err = base64.StdEncoding.DecodeString(m.keyB64)
			})
			rules.WhenTrue(err == nil, func() {
				cert, err = tls.X509KeyPair(cb, kb)
			})
		} else {
			cert, err = tls.LoadX509KeyPair(m.certFile, m.keyFile)
		}
		if err != nil {
			return err
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS11,
			ServerName:   m.serverName,
		}
	}
	m.c, err = valkey.NewClient(valkey.ClientOption{
		TLSConfig: tlsConfig,
		Password:  m.pwd,

		InitAddress: []string{
			m.addr,
		},
		ClientTrackingOptions: nil,
		SelectDB:              m.db,
	})
	return m.Ping()
}

func (m *valkeyManager) spanName(v string) string {
	return rules.WhenTrueR1(len(m.scope) > 0, func() string {
		return m.scope + ".inmemory." + v
	}, func() string {
		return "inmemory." + v
	})
}

func (m *valkeyManager) Set(ctx context.Context, key string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.Do(
			newCtx,
			m.c.B().Set().
				Key(injectPrefixWhenDefined(m.prefix, key)).
				Value(utils.CastToStringND(value)).Build(),
		).Error()
	})
}

func (m *valkeyManager) SetWhenNotExist(ctx context.Context, key string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.Do(
			newCtx,
			m.c.B().Setnx().
				Key(injectPrefixWhenDefined(m.prefix, key)).
				Value(utils.CastToStringND(value)).Build(),
		).Error()
	})
}

func (m *valkeyManager) SetWithExpire(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_ex"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.Do(
			newCtx,
			m.c.B().Setex().
				Key(injectPrefixWhenDefined(m.prefix, key)).
				Seconds(int64(expire.Seconds())).
				Value(utils.CastToStringND(value)).Build(),
		).Error()
	})
}

func (m *valkeyManager) SetWithExpireWhenNotExist(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_ex_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		if err := m.c.Do(newCtx, m.c.B().Get().Key(key).Build()).Error(); err != nil {
			return err
		}
		return m.c.Do(
			newCtx,
			m.c.B().Setex().
				Key(injectPrefixWhenDefined(m.prefix, key)).
				Seconds(int64(expire.Seconds())).
				Value(utils.CastToStringND(value)).Build(),
		).Error()
	})
}

func (m *valkeyManager) SetString(ctx context.Context, key, value string, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_str"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.Do(
			newCtx,
			m.c.B().Set().
				Key(injectPrefixWhenDefined(m.prefix, key)).
				Value(value).Build(),
		).Error()
	})
}

func (m *valkeyManager) Get(ctx context.Context, key string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("get"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		rs, err := m.c.Do(
			newCtx,
			m.c.B().Get().
				Key(injectPrefixWhenDefined(m.prefix, key)).
				Build(),
		).ToString()
		if err != nil {
			return err
		}
		return Scan([]byte(rs), bindTo)
	})
}

func (m *valkeyManager) GetString(ctx context.Context, key string) string {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("get_str"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) string {
		rs, err := m.c.Do(
			newCtx,
			m.c.B().Get().
				Key(injectPrefixWhenDefined(m.prefix, key)).
				Build(),
		).ToString()
		if err != nil {
			return ""
		}
		return rs
	})
}

func (m *valkeyManager) Del(ctx context.Context, keys ...string) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("del"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.StringSlice("keys", keys)),
	}, func(newCtx context.Context) error {
		return m.c.Do(
			newCtx,
			m.c.B().Del().
				Key(injectPrefixIntoKeysWhenDefined(m.prefix, keys...)...).
				Build(),
		).Error()
	})
}

func (m *valkeyManager) HSet(ctx context.Context, key string, value ...interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_set"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		if len(value)%2 != 0 {
			return errors.New("value must be an even length of field -> value")
		}
		fv := m.c.B().Hset().
			Key(injectPrefixWhenDefined(m.prefix, key)).
			FieldValue()
		for i := 0; i < len(value); i += 2 {
			fv.FieldValue(utils.CastToStringND(value[i]), utils.CastToStringND(value[i+1]))
		}
		return m.c.Do(
			newCtx, fv.Build(),
		).Error()
	})
}

func (m *valkeyManager) HSetWhenNotExist(ctx context.Context, key, field string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_set_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		fv := m.c.B().Hsetnx().
			Key(injectPrefixWhenDefined(m.prefix, key)).Field(field).Value(utils.CastToStringND(value))
		return m.c.Do(
			newCtx, fv.Build(),
		).Error()
	})
}

func (m *valkeyManager) HGet(ctx context.Context, key, field string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_get"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		fv := m.c.B().Hget().
			Key(injectPrefixWhenDefined(m.prefix, key)).Field(field)
		rs, err := m.c.Do(
			newCtx, fv.Build(),
		).ToString()
		if err != nil {
			return err
		}
		return Scan([]byte(rs), bindTo)
	})
}

func (m *valkeyManager) HDel(ctx context.Context, key string, fields ...string) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_del"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		fv := m.c.B().Hdel().
			Key(injectPrefixWhenDefined(m.prefix, key)).Field(fields...)
		return m.c.Do(
			newCtx, fv.Build(),
		).Error()
	})
}

func (m *valkeyManager) HGetAll(ctx context.Context, key string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_get_all"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		fv := m.c.B().Hgetall().
			Key(injectPrefixWhenDefined(m.prefix, key))
		rs, err := m.c.Do(
			newCtx, fv.Build(),
		).ToString()
		if err != nil {
			return err
		}
		return Scan([]byte(rs), bindTo)
	})
}

func (m *valkeyManager) HMSet(ctx context.Context, key string, values ...interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hm_set"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		if len(values)%2 != 0 {
			return errors.New("value must be an even length of field -> value")
		}
		fv := m.c.B().Hmset().
			Key(injectPrefixWhenDefined(m.prefix, key)).FieldValue()
		for i := 0; i < len(values); i += 2 {
			fv.FieldValue(utils.CastToStringND(values[i]), utils.CastToStringND(values[i+1]))
		}
		return m.c.Do(
			newCtx, fv.Build(),
		).Error()
	})
}

func (m *valkeyManager) HMGet(ctx context.Context, key, field string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hm_get"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		fv := m.c.B().Hmget().
			Key(injectPrefixWhenDefined(m.prefix, key)).Field(field)
		rs, err := m.c.Do(
			newCtx, fv.Build(),
		).ToString()
		if err != nil {
			return err
		}
		return Scan([]byte(rs), bindTo)
	})
}

func (m *valkeyManager) Ping() error {
	return m.c.Do(context.Background(), m.c.B().Ping().Build()).Error()
}

func NewValKey(addr string, opts ...Option[valkeyManager]) Manager {
	m := &valkeyManager{addr: addr}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}
