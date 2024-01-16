/**
 * @Author: steven
 * @Description:
 * @File: cache
 * @Date: 29/09/23 10.46
 */

package cache

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Manager interface {
	MustConnect(ctx context.Context) Manager
	Connect(ctx context.Context) error

	Set(ctx context.Context, key string, value interface{}) error
	SetWhenNotExist(ctx context.Context, key string, value interface{}) error
	SetWithExpire(ctx context.Context, key string, value interface{}, expire time.Duration) error
	SetWithExpireWhenNotExist(ctx context.Context, key string, value interface{}, expire time.Duration) error
	SetString(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string, bindTo interface{}) error
	GetString(ctx context.Context, key string) string
	Del(ctx context.Context, keys ...string) error

	HSet(ctx context.Context, key string, value ...interface{}) error
	HSetWhenNotExist(ctx context.Context, key, field string, value interface{}) error
	HGet(ctx context.Context, key, field string, bindTo interface{}) error
	HDel(ctx context.Context, key string, fields ...string) error
	HGetAll(ctx context.Context, key string, bindTo interface{}) error

	HMSet(ctx context.Context, key string, values ...interface{}) error
	HMGet(ctx context.Context, key, field string, bindTo interface{}) error

	Ping() error
}

type redisManager struct {
	r *redis.Client

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

func (m *redisManager) injectPrefixWhenDefined(key string) string {
	if len(m.prefix) < 1 {
		return key
	}
	return m.prefix + "_" + key
}

func (m *redisManager) injectPrefixIntoKeysWhenDefined(keys ...string) []string {
	if len(m.prefix) < 1 {
		return keys
	}
	rs := make([]string, len(keys))
	for _, key := range keys {
		rs = append(rs, m.prefix+"_"+key)
	}
	return rs
}

func wrapTelemetryTuple1[T any](ctx context.Context, tc trace.Tracer, spanName string, spanAttr []trace.SpanStartOption, task func(ctx context.Context) T) T {
	newCtx, span := tc.Start(ctx, spanName, append(spanAttr, trace.WithSpanKind(trace.SpanKindClient))...)
	defer span.End()
	return task(newCtx)
}

func (m *redisManager) spanName(v string) string {
	return utils.IfER(len(m.scope) > 0, func() string {
		return m.scope + ".cache." + v
	}, func() string {
		return "cache." + v
	})
}

func (m *redisManager) Ping() error {
	return m.r.Ping(context.Background()).Err()
}

func (m *redisManager) HSet(ctx context.Context, key string, value ...interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hset"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.r.HSet(newCtx, m.injectPrefixWhenDefined(key), value...).Err()
	})
}

func (m *redisManager) HSetWhenNotExist(ctx context.Context, key, field string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hset_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
		trace.WithAttributes(attribute.String("field", field)),
	}, func(newCtx context.Context) error {
		return m.r.HSetNX(newCtx, m.injectPrefixWhenDefined(key), field, value).Err()
	})
}

func (m *redisManager) HDel(ctx context.Context, key string, fields ...string) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hdel"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
		trace.WithAttributes(attribute.StringSlice("fields", fields)),
	}, func(newCtx context.Context) error {
		return m.r.HDel(newCtx, m.injectPrefixWhenDefined(key), fields...).Err()
	})
}

func (m *redisManager) HGet(ctx context.Context, key, field string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hget"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
		trace.WithAttributes(attribute.String("field", field)),
	}, func(newCtx context.Context) error {
		return m.r.HGet(newCtx, m.injectPrefixWhenDefined(key), field).Scan(bindTo)
	})
}

func (m *redisManager) HGetAll(ctx context.Context, key string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hgetall"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.r.HGetAll(newCtx, m.injectPrefixWhenDefined(key)).Scan(bindTo)
	})
}

func (m *redisManager) HMSet(ctx context.Context, key string, values ...interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hmset"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.r.HSet(newCtx, m.injectPrefixWhenDefined(key), values...).Err()
	})
}

func (m *redisManager) HMGet(ctx context.Context, key, field string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("hmget"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
		trace.WithAttributes(attribute.String("field", field)),
	}, func(newCtx context.Context) error {
		return m.r.HMGet(newCtx, m.injectPrefixWhenDefined(key), field).Scan(bindTo)
	})
}

func (m *redisManager) Set(ctx context.Context, key string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.r.Set(newCtx, m.injectPrefixWhenDefined(key), value, 0).Err()
	})
}

func (m *redisManager) SetWhenNotExist(ctx context.Context, key string, value interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.r.SetNX(newCtx, m.injectPrefixWhenDefined(key), value, 0).Err()
	})
}

func (m *redisManager) SetWithExpire(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_wx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.r.Set(newCtx, m.injectPrefixWhenDefined(key), value, expire).Err()
	})
}

func (m *redisManager) SetWithExpireWhenNotExist(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_wx_nx"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.r.SetNX(newCtx, m.injectPrefixWhenDefined(key), value, expire).Err()
	})
}

func (m *redisManager) Get(ctx context.Context, key string, bindTo interface{}) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("get"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		return m.r.Get(newCtx, m.injectPrefixWhenDefined(key)).Scan(bindTo)
	})
}

func (m *redisManager) Del(ctx context.Context, keys ...string) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("del"), []trace.SpanStartOption{}, func(newCtx context.Context) error {
		return m.r.Del(newCtx, m.injectPrefixIntoKeysWhenDefined(keys...)...).Err()
	})
}

func (m *redisManager) SetString(ctx context.Context, key, value string, expire time.Duration) error {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("set_str"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) error {
		_, err := m.r.Set(ctx, m.injectPrefixWhenDefined(key), value, expire).Result()
		return err
	})
}

func (m *redisManager) GetString(ctx context.Context, key string) (rs string) {
	return wrapTelemetryTuple1(ctx, m.tm.Tracer(), m.spanName("get_str"), []trace.SpanStartOption{
		trace.WithAttributes(attribute.String("key", key)),
	}, func(newCtx context.Context) string {
		var err error
		if rs, err = m.r.Get(newCtx, m.injectPrefixWhenDefined(key)).Result(); err != nil {
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
			utils.IfTrueThen(err == nil, func() {
				kb, err = base64.StdEncoding.DecodeString(m.keyB64)
			})
			utils.IfTrueThen(err == nil, func() {
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
	m.r = redis.NewClient(&redis.Options{
		Addr:      m.addr,
		Password:  m.pwd, // no password set
		DB:        m.db,  // use default DB
		TLSConfig: tlsConfig,
	})
	return m.r.Ping(ctx).Err()
}

func NewRedis(addr, pwd string, db int, opts ...Option) Manager {
	m := &redisManager{addr: addr, pwd: pwd, db: db}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}
