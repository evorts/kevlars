package midware

import (
	"context"
	"github.com/evorts/kevlars/auth"
	"github.com/evorts/kevlars/contracts"
	"github.com/evorts/kevlars/jwe"
	"github.com/evorts/kevlars/logger"
	"github.com/evorts/kevlars/requests"
	"github.com/evorts/kevlars/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApiKeySecretMap struct {
	Key      string `mapstructure:"key" json:"key"`
	ClientId string `mapstructure:"client_id" json:"client_id"`
}

func EchoWithUserAuthorization(aum auth.UserManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			token := tokenFromHeader(req)
			if len(token) < 1 {
				return c.JSON(contracts.NewResponseFail(http.StatusUnauthorized, "not eligible to access this resource", contracts.ErrorDetail{}))
			}
			ctx := req.Context()

			//introspect token
			var claim jwe.Claim
			if err := aum.Introspect(ctx, token, &claim); err != nil {
				return c.JSON(contracts.NewResponseFail(http.StatusUnauthorized, "invalid token", contracts.ErrorDetail{}))
			}

			//get resource name from url
			resourceName := c.Request().URL.Path
			scope := auth.Scope("").FromHttpMethod(req.Method)

			//get requesting party token from user token
			//introspect token and get user id from within it
			permitted, err := aum.IsAllowed(ctx, claim.ID, resourceName, scope)
			if err != nil || !permitted {
				return c.JSON(contracts.NewResponseFail(http.StatusUnauthorized, "not permitted to access this resource due to insufficient permission", contracts.ErrorDetail{}))
			}

			return next(c)
		}
	}
}

func EchoWithClientAuthorization(ac auth.ClientManager, log logger.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			resource := req.URL.Path
			method := req.Method
			key := apiKeyFromHeader(req)
			scope := auth.Scope("").FromHttpMethod(req.Method)
			cm, allowed := ac.IsAllowed(key, resource, scope)
			if !allowed {
				log.ErrorWithProps(map[string]interface{}{
					"cid":    cm,
					"method": method,
					"path":   resource,
				}, "request not allowed")
				return c.JSON(contracts.NewResponseFail(http.StatusUnauthorized, "Not authorized to access this resource", contracts.ErrorDetail{
					Code: "ERR:NOK:AUTH",
					Errors: map[string]string{
						"err": "key not acceptable",
					},
				}))
			}
			return next(c)
		}
	}
}

func echoWithAuthApiKeySecretsEligibleClients(maps []ApiKeySecretMap, compareFunc func(items []string, item string) bool, clientIds ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			reqKey := apiKeyFromHeader(req)
			if len(reqKey) > 0 {
				reqCtx := req.Context()
				hasClients := len(clientIds) > 0
				for _, v := range maps {
					if hasClients && !compareFunc(clientIds, v.ClientId) {
						continue
					}
					if v.Key == reqKey {
						// set client id in both echo context and native context
						c.Set(requests.ContextClientId.String(), v.ClientId)
						c.SetRequest(req.WithContext(context.WithValue(reqCtx, requests.ContextClientId, v.ClientId)))
						return next(c)
					}
				}
			}
			return c.JSON(contracts.NewResponseFail(http.StatusUnauthorized, "Not authorized to access this resource", contracts.ErrorDetail{
				Code: "ERR:NOK:AUTH",
				Errors: map[string]string{
					"err": "key not acceptable",
				},
			}))
		}
	}
}

func EchoWithAuthApiKeySecretsWithWhitelistClient(maps []ApiKeySecretMap, clientIds ...string) echo.MiddlewareFunc {
	return echoWithAuthApiKeySecretsEligibleClients(maps, func(items []string, item string) bool {
		return utils.InArray(items, item)
	}, clientIds...)
}

func EchoWithAuthApiKeySecretsWithBlacklistClient(maps []ApiKeySecretMap, clientIds ...string) echo.MiddlewareFunc {
	return echoWithAuthApiKeySecretsEligibleClients(maps, func(items []string, item string) bool {
		return !utils.InArray(items, item)
	}, clientIds...)
}

func EchoWithAuthApiKeySecrets(maps []ApiKeySecretMap) echo.MiddlewareFunc {
	return EchoWithAuthApiKeySecretsWithWhitelistClient(maps)
}
