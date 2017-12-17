package main

import (
	"github.com/NeuronFramework/restful"
	"github.com/NeuronUser/user/api-private/gen/restapi"
	"github.com/NeuronUser/user/api-private/gen/restapi/operations"
	"github.com/NeuronUser/user/cmd/user-private-api/handler"
	"github.com/go-openapi/loads"
	"net/http"
	"os"
)

func main() {
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "8086")

	restful.Run(func() (http.Handler, error) {
		h, err := handler.NewUserHandler()
		if err != nil {
			return nil, err
		}

		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return nil, err
		}

		api := operations.NewUserPrivateAPI(swaggerSpec)
		api.OauthStateHandler = operations.OauthStateHandlerFunc(h.OauthState)
		api.OauthJumpHandler = operations.OauthJumpHandlerFunc(h.OauthJump)
		api.RefreshTokenHandler = operations.RefreshTokenHandlerFunc(h.RefreshToken)
		api.LogoutHandler = operations.LogoutHandlerFunc(h.Logout)

		return api.Serve(nil), nil
	})
}
