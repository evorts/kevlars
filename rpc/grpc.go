/**
 * @Author: steven
 * @Description:
 * @File: grpc
 * @Date: 29/09/23 10.48
 */

package rpc

import (
	"context"
	"crypto/tls"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/telemetry"
	"github.com/evorts/kevlars/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

type ClientManager interface {
	MustConnect() ClientManager
	Client() grpc.ClientConnInterface
	connect() error
	Teardown()
}

type grpcClientManager struct {
	conn                    *grpc.ClientConn
	creds                   credentials.TransportCredentials
	server                  string
	name                    string
	log                     logger.Manager
	tm                      telemetry.Manager
	metric                  telemetry.MetricsManager
	token                   string
	telemetryEnabled        bool
	metricEnabled           bool
	timeout                 *time.Duration
	logRequestPayload       bool
	logRequestPayloadInJson bool
	logResponse             bool
}

func (g *grpcClientManager) metricIsEnabled() bool {
	return g.metricEnabled && g.metric != nil
}

func (g *grpcClientManager) Client() grpc.ClientConnInterface {
	return g.conn
}

func (g *grpcClientManager) connect() error {
	var err error
	opts := make([]grpc.DialOption, 0)
	if g.creds == nil {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(g.creds))
	}
	if g.token != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(AuthorizationCredential(g.token)))
	}
	unaryInterceptors := make([]grpc.UnaryClientInterceptor, 0)
	utils.IfTrueThen(g.telemetryEnabled, func() {
		unaryInterceptors = append(unaryInterceptors, GrpcTraceInterceptor(g.tm.Tracer()))
	})
	utils.IfTrueThen(g.logRequestPayload, func() {
		unaryInterceptors = append(unaryInterceptors, GrpcLogRequestPayloadInterceptor(g.logRequestPayloadInJson, g.log.InfoWithProps))
	})
	ctx := context.Background()
	utils.IfTrueThen(g.timeout != nil && g.timeout.Milliseconds() > 0, func() {
		unaryInterceptors = append(unaryInterceptors, GrpcTimeoutInterceptor(*g.timeout))
	})
	utils.IfTrueThen(g.metricIsEnabled(), func() {
		unaryInterceptors = append(unaryInterceptors, GrpcMetricInterceptor(g.metric))
	})
	utils.IfTrueThen(len(unaryInterceptors) > 0, func() {
		opts = append(opts, grpc.WithChainUnaryInterceptor(unaryInterceptors...))
	})
	g.conn, err = grpc.DialContext(ctx, g.server, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (g *grpcClientManager) MustConnect() ClientManager {
	if err := g.connect(); err != nil {
		panic(err)
	}
	return g
}

func (g *grpcClientManager) Teardown() {
	err := g.conn.Close()
	if g.log == nil {
		return
	}
	g.log.WhenError(err)
}

func NewClient(server string, opts ...Option) ClientManager {
	m := &grpcClientManager{server: server}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m
}
