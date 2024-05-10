/**
 * @Author: steven
 * @Description:
 * @File: cache
 * @Date: 29/09/23 10.46
 */

package inmemory

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"github.com/evorts/kevlars/rules"
	"github.com/evorts/kevlars/telemetry"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type redisManager struct {
	c *redis.Client

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

func (m *redisManager) spanName(v string) string {
	return rules.WhenTrueR1(len(m.scope) > 0, func() string {
		return m.scope + ".inmemory." + v
	}, func() string {
		return "inmemory." + v
	})
}

func (m *redisManager) Ping() error {
	return m.c.Ping(context.Background()).Err()
}

func (m *redisManager) HSet(ctx context.Context, key string, value ...interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_set"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.HSet(newCtx, injectPrefixWhenDefined(m.prefix, key), value...).Err()
	})
}

func (m *redisManager) HSetWhenNotExist(ctx context.Context, key, field string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_set_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
		trace.WithAttributes(attribute.String("field", field)),
	}, func(newCtx context.Context) error {
		return m.c.HSetNX(newCtx, injectPrefixWhenDefined(m.prefix, key), field, value).Err()
	})
}

func (m *redisManager) HDel(ctx context.Context, key string, fields ...string) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_del"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
		trace.WithAttributes(attribute.StringSlice("fields", fields)),
	}, func(newCtx context.Context) error {
		return m.c.HDel(newCtx, injectPrefixWhenDefined(m.prefix, key), fields...).Err()
	})
}

func (m *redisManager) HGet(ctx context.Context, key, field string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_get"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
		trace.WithAttributes(attribute.String("field", field)),
	}, func(newCtx context.Context) error {
		return m.c.HGet(newCtx, injectPrefixWhenDefined(m.prefix, key), field).Scan(bindTo)
	})
}

func (m *redisManager) HGetAll(ctx context.Context, key string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("h_get_all"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.HGetAll(newCtx, injectPrefixWhenDefined(m.prefix, key)).Scan(bindTo)
	})
}

func (m *redisManager) HMSet(ctx context.Context, key string, values ...interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hm_set"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.HSet(newCtx, injectPrefixWhenDefined(m.prefix, key), values...).Err()
	})
}

func (m *redisManager) HMGet(ctx context.Context, key, field string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hm_get"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
		trace.WithAttributes(attribute.String("field", field)),
	}, func(newCtx context.Context) error {
		return m.c.HMGet(newCtx, injectPrefixWhenDefined(m.prefix, key), field).Scan(bindTo)
	})
}

func (m *redisManager) Set(ctx context.Context, key string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.Set(newCtx, injectPrefixWhenDefined(m.prefix, key), value, 0).Err()
	})
}

func (m *redisManager) SetWhenNotExist(ctx context.Context, key string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.SetNX(newCtx, injectPrefixWhenDefined(m.prefix, key), value, 0).Err()
	})
}

func (m *redisManager) SetWithExpire(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_ex"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.Set(newCtx, injectPrefixWhenDefined(m.prefix, key), value, expire).Err()
	})
}

func (m *redisManager) SetWithExpireWhenNotExist(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_ex_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.SetNX(newCtx, injectPrefixWhenDefined(m.prefix, key), value, expire).Err()
	})
}

func (m *redisManager) Get(ctx context.Context, key string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("get"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.c.Get(newCtx, injectPrefixWhenDefined(m.prefix, key)).Scan(bindTo)
	})
}

func (m *redisManager) Del(ctx context.Context, keys ...string) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("del"), []trace.SpanStartOption{}, func(newCtx context.Context) error {
		return m.c.Del(newCtx, injectPrefixIntoKeysWhenDefined(m.prefix, keys...)...).Err()
	})
}

func (m *redisManager) SetString(ctx context.Context, key, value string, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_str"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		_, err := m.c.Set(ctx, injectPrefixWhenDefined(m.prefix, key), value, expire).Result()
		return err
	})
}

func (m *redisManager) GetString(ctx context.Context, key string) (rs string) {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("get_str"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) string {
		var err error
		if rs, err = m.c.Get(newCtx, injectPrefixWhenDefined(m.prefix, key)).Result(); err != nil {
			return ""
		}
		return rs
	})
}

func (m *redisManager) MustConnect(ctx context.Context) Manager {
	if err := m.Connect(ctx); err != nil {
		panic(err)
	}
	return m
}

func (m *redisManager) Connect(ctx context.Context) error {
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
	m.c = redis.NewClient(&redis.Options{
		Addr:      m.addr,
		Password:  m.pwd, // no password set
		DB:        m.db,  // use default DB
		TLSConfig: tlsConfig,
	})
	return m.c.Ping(ctx).Err()
}

func NewRedis(addr string, opts ...Option[redisManager]) Manager {
	m := &redisManager{addr: addr}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}
