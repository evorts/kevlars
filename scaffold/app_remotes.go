/**
 * @Author: steven
 * @Description:
 * @File: app_remotes
 * @Date: 24/12/23 11.24
 */

package scaffold

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/evorts/kevlars/rest"
	"github.com/evorts/kevlars/rpc"
	"github.com/evorts/kevlars/soap"
	"github.com/evorts/kevlars/utils"
	semConv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc/credentials"
	"time"
)

type IRemotes interface {
	WithRestClients() IApplication
	WithSoapClients() IApplication
	WithGrpcClients() IApplication

	REST(context string) rest.Manager
	SOAP(context string) soap.Manager
	GRPC(context string) rpc.ClientManager
}

type authenticationType string

const (
	authTypeBasic   authenticationType = "basic"
	authTypeHeaders authenticationType = "headers"
	authTypeOAuth   authenticationType = "oauth"
)

func (at authenticationType) String() string {
	return string(at)
}

func (at authenticationType) KV() (string, string) {
	switch at {
	case authTypeBasic:
		return "username", "password"
	case authTypeHeaders:
		return "key", "value"
	case authTypeOAuth:
		return "endpoint", "scope"
	default:
		return "", ""
	}
}

type service struct {
	Name                    string `mapstructure:"name"`
	ServiceUrl              string `mapstructure:"service_url"`
	ProxyUrl                string `mapstructure:"proxy_url"`
	LogRequestPayload       bool   `mapstructure:"log_request_payload"`
	LogRequestPayloadInJson bool   `mapstructure:"log_request_payload_in_json"`
	LogResponse             bool   `mapstructure:"log_response"`
	CircuitBreaker          struct {
		Enabled    bool `mapstructure:"enabled" mapstructure:"enabled"`
		MaxRequest int  `mapstructure:"max_request" mapstructure:"max_request"`
		Interval   int  `mapstructure:"interval" mapstructure:"interval"`
		Timeout    int  `mapstructure:"timeout" mapstructure:"timeout"`
	} `mapstructure:"circuit_breaker"`
}

func (app *Application) REST(context string) rest.Manager {
	if c, ok := app.restClients[context]; ok {
		return c
	}
	return nil
}

func (app *Application) SOAP(context string) soap.Manager {
	if c, ok := app.soapClients[context]; ok {
		return c
	}
	return nil
}

func (app *Application) GRPC(context string) rpc.ClientManager {
	if c, ok := app.grpcClients[context]; ok {
		return c
	}
	return nil
}

func (app *Application) WithRestClients() IApplication {
	var restServices []service
	err := app.Config().UnmarshalTo("rest.services", &restServices)
	if err != nil || len(restServices) < 1 {
		panic("asking rest clients but no configuration defined")
	}
	var (
		telemetryEnabled = app.Config().GetBool("rest.telemetry_enabled")
		metricEnabled    = app.Config().GetBool("rest.metric_enabled")
		name             = app.Config().GetStringOrElse("rest.name", "rest.client")
	)
	app.restClients = make(map[string]rest.Manager)
	for _, svc := range restServices {
		contextName := fmt.Sprintf("%s.%s", name, svc.Name)
		opts := []rest.Option{
			rest.WithName(contextName),
			rest.WithBaseUrl(svc.ServiceUrl),
			rest.WithDebugging(app.Config().GetBool("rest.debug_enabled")),
			rest.WithTracing(app.Config().GetBool("rest.trace_enabled")),
			rest.WithLogger(app.Log()),
		}
		utils.IfTrueThen(telemetryEnabled, func() {
			tOps := []rest.TransportOption{
				rest.TransportWithTelemetry(app.Telemetry()),
				rest.TransportWithName(contextName),
				rest.TransportWithAttributes(semConv.ServiceNameKey.String(app.name)),
			}
			opts = append(opts, rest.WithTransport(tOps...))
		})
		utils.IfTrueThen(metricEnabled, func() {
			opts = append(opts, rest.WithMetrics(metricEnabled, app.Metrics()))
		})
		utils.IfTrueThen(len(svc.ProxyUrl) > 0, func() {
			opts = append(opts, rest.WithProxy(svc.ProxyUrl))
		})
		utils.IfTrueThen(svc.LogRequestPayload, func() {
			opts = append(opts, rest.WithLogRequest(true))
		})
		utils.IfTrueThen(svc.LogResponse, func() {
			opts = append(opts, rest.WithLogResponse(true))
		})
		app.restClients[svc.Name] = rest.New(opts...)
		utils.IfTrueThen(svc.CircuitBreaker.Enabled, func() {
			app.restClients[svc.Name].WithCircuitBreaker(
				uint32(svc.CircuitBreaker.MaxRequest),
				time.Duration(svc.CircuitBreaker.Interval)*time.Millisecond,
				time.Duration(svc.CircuitBreaker.Timeout)*time.Millisecond,
			)
		})
		// check whether there's a need for authentication
		var svcAuth map[string]interface{}
		svcAuthCfgKey := fmt.Sprintf("rest.auth.%s", svc.Name)
		err = app.Config().UnmarshalTo(svcAuthCfgKey, &svcAuth)
		utils.IfTrueThen(
			utils.AND(
				err == nil,
				utils.CastToBoolND(utils.GetValueOnMap(svcAuth, "enabled", false)),
			),
			func() {
				use := utils.CastToStringND(utils.GetValueOnMap(svcAuth, "use", ""))
				switch use {
				case authTypeHeaders.String():
					k, v := authTypeHeaders.KV()
					kk := app.Config().GetString(fmt.Sprintf("%s.%s.%s", svcAuthCfgKey, use, k))
					vv := app.Config().GetString(fmt.Sprintf("%s.%s.%s", svcAuthCfgKey, use, v))
					utils.IfTrueThen(
						utils.AND(
							len(kk) > 0,
							len(vv) > 0,
						),
						func() {
							app.restClients[svc.Name].WithDefaultHeaders(map[string]string{
								kk: vv,
							})
						},
					)
				}
			},
		)
	}
	return app
}

func (app *Application) WithSoapClients() IApplication {
	var soapServices []service
	err := app.Config().UnmarshalTo("soap.services", &soapServices)
	if err != nil || len(soapServices) < 1 {
		panic("asking soap managers but no configuration defined")
	}
	var (
		telemetryEnabled = app.Config().GetBool("soap.telemetry_enabled")
		metricEnabled    = app.Config().GetBool("soap.metric_enabled")
		name             = app.Config().GetStringOrElse("soap.name", "soap.client")
		largeFileTimeout = app.Config().GetIntOrElse("soap.large_file_timeout", 300) // default 5 minute
	)
	app.soapClients = make(map[string]soap.Manager)
	for _, svc := range soapServices {
		contextName := fmt.Sprintf("%s.%s", name, svc.Name)
		opts := []soap.Option{
			soap.WithServiceUrl(svc.ServiceUrl),
			soap.WithLogger(app.Log()),
			soap.WithName(contextName),
		}
		if largeFileTimeout > 0 {
			opts = append(opts, soap.WithLargeContentTimeout(time.Duration(largeFileTimeout)*time.Second))
		}
		utils.IfTrueThen(telemetryEnabled, func() {
			opts = append(opts,
				soap.WithTransport(
					soap.WithTelemetry(app.Telemetry()),
					soap.WithTransportName(contextName),
				),
			)
		})
		utils.IfTrueThen(metricEnabled, func() {
			opts = append(opts, soap.WithMetrics(metricEnabled, app.Metrics()))
		})
		authType := app.Config().GetString(fmt.Sprintf("soap.auth.%s.use", svc.Name))
		if authType == authTypeBasic.String() {
			opts = append(opts,
				soap.WithBasicAuth(
					app.Config().GetString(fmt.Sprintf("soap.auth.%s.%s.username", svc.Name, authType)),
					app.Config().GetString(fmt.Sprintf("soap.auth.%s.%s.password", svc.Name, authType)),
				),
			)
		}
		app.soapClients[svc.Name] = soap.New(opts...)
		if svc.CircuitBreaker.Enabled {
			app.soapClients[svc.Name].WithCircuitBreaker(
				uint32(svc.CircuitBreaker.MaxRequest),
				time.Duration(svc.CircuitBreaker.Interval)*time.Millisecond,
				time.Duration(svc.CircuitBreaker.Timeout)*time.Millisecond,
			)
		}
	}
	return app
}

func (app *Application) WithGrpcClients() IApplication {
	var grpcServices []service
	err := app.Config().UnmarshalTo("grpc.services", &grpcServices)
	if err != nil || len(grpcServices) < 1 {
		panic("asking grpc managers but no configuration defined")
	}
	var (
		telemetryEnabled = app.Config().GetBool("grpc.telemetry_enabled")
		metricEnabled    = app.Config().GetBool("grpc.metric_enabled")
		name             = app.Config().GetStringOrElse("grpc.name", "grpc.client")
	)
	app.grpcClients = make(map[string]rpc.ClientManager)
	for _, svc := range grpcServices {
		contextName := fmt.Sprintf("%s.%s", name, svc.Name)
		opts := []rpc.Option{
			rpc.WithName(contextName),
			rpc.WithLogger(app.Log()),
		}
		utils.IfTrueThen(telemetryEnabled, func() {
			opts = append(opts, rpc.WithTelemetry(app.Telemetry()), rpc.WithTelemetryEnabled(telemetryEnabled))
		})
		utils.IfTrueThen(svc.CircuitBreaker.Timeout > 0, func() {
			opts = append(opts, rpc.WithTimeout(time.Duration(svc.CircuitBreaker.Timeout)*time.Millisecond))
		})
		utils.IfTrueThen(metricEnabled, func() {
			opts = append(opts, rpc.WithMetrics(metricEnabled, app.Metrics()))
		})
		utils.IfTrueThen(svc.LogRequestPayload, func() {
			opts = append(opts, rpc.WithLoggingRequestPayload(svc.LogRequestPayload, svc.LogRequestPayloadInJson))
		})
		authCfgKey := "grpc.auth." + svc.Name
		utils.IfTrueThen(app.Config().GetBool(authCfgKey+".enabled"), func() {
			authType := app.Config().GetString(authCfgKey + ".use")
			switch authType {
			case "token":
				token := app.Config().GetString(authCfgKey + "." + authType)
				opts = append(opts, rpc.WithAuthorizationToken(token))
			case "ssl_tls":
				var cert tls.Certificate
				b64Creds := app.Config().GetString(authCfgKey + "." + authType + ".cert_b64")
				b64Key := app.Config().GetString(authCfgKey + "." + authType + ".key_b64")
				if len(b64Creds) > 0 {
					var cb, kb []byte
					cb, err = base64.StdEncoding.DecodeString(b64Creds)
					utils.IfTrueThen(err == nil, func() {
						kb, err = base64.StdEncoding.DecodeString(b64Key)
					})
					utils.IfTrueThen(err == nil, func() {
						cert, err = tls.X509KeyPair(cb, kb)
					})
					if err != nil {
						return
					}
					tlsCreds := credentials.NewTLS(&tls.Config{
						Certificates: []tls.Certificate{cert},
						MinVersion:   tls.VersionTLS11,
						ServerName:   app.name,
					})
					opts = append(opts, rpc.WithTransportCreds(tlsCreds))
				}
			default:
				opts = append(opts, rpc.WithTransportCreds(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
			}
		})
		app.grpcClients[svc.Name] = rpc.NewClient(svc.ServiceUrl, opts...).MustConnect()
	}
	return app
}
