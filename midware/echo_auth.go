package midware

import (
	"context"
	"github.com/evorts/kevlars/auth"
	"github.com/evorts/kevlars/contracts"
	"github.com/evorts/kevlars/requests"
	"github.com/evorts/kevlars/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApiKeySecretMap struct {
	Key      string `mapstructure:"key" json:"key"`
	ClientId string `mapstructure:"client_id" json:"client_id"`
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

func EchoWithAuthToken(am auth.OAuthManager, realm, clientID string, clientSecret string, permissionScope string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			token := tokenFromHeader(req)
			if len(token) < 1 {
				return c.JSON(contracts.NewResponseFail(http.StatusUnauthorized, "not eligible to access this resource", contracts.ErrorDetail{}))
			}

			//introspect token to keycloak -- check if token eligible to access the resource
			ctx := req.Context()

			//get resource nam from url
			resourceName := c.Request().URL.String()

			//get requesting party token from user token
			permissions, err := am.GetPermissions(ctx, token, realm)
			if err != nil {
				return c.JSON(contracts.NewResponseFail(http.StatusUnauthorized, "not permitted to access this resource due to failed to retrieve requesting party token ", contracts.ErrorDetail{}))
			}

			//check request authorization
			for _, permission := range permissions {
				resName := permission.ResourceName
				if resName == resourceName && permission.Scopes != nil {
					if utils.InArray(permission.Scopes, permissionScope) {
						return next(c)
					}

				}
			}
			return c.JSON(contracts.NewResponseFail(http.StatusUnauthorized, "not permitted to access this resource due to insufficient permission", contracts.ErrorDetail{}))
		}
	}
}
