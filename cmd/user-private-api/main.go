package main

import (
	"github.com/NeuronFramework/log"
	"github.com/NeuronFramework/restful"
	"github.com/NeuronUser/user/api-private/gen/restapi"
	"github.com/NeuronUser/user/api-private/gen/restapi/operations"
	"github.com/NeuronUser/user/cmd/user-private-api/handler"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	log.Init(true)

	middleware.Debug = false

	logger := zap.L().Named("main")

	var bindAddr string

	cmd := cobra.Command{}
	cmd.PersistentFlags().StringVar(&bindAddr, "bind-addr", ":8085", "api server bind addr")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return err
		}
		api := operations.NewUserPrivateAPI(swaggerSpec)

		h, err := handler.NewUserHandler()
		if err != nil {
			return err
		}

		api.NewOauthStateHandler = operations.NewOauthStateHandlerFunc(h.NewOauthState)
		api.OauthJumpHandler = operations.OauthJumpHandlerFunc(h.OauthJump)
		api.RefreshTokenHandler = operations.RefreshTokenHandlerFunc(h.RefreshToken)
		api.LogoutHandler = operations.LogoutHandlerFunc(h.Logout)

		logger.Info("Start server", zap.String("addr", bindAddr))
		err = http.ListenAndServe(bindAddr,
			restful.Recovery(cors.AllowAll().Handler(api.Serve(nil))))
		if err != nil {
			return err
		}

		return nil
	}
	cmd.Execute()
}
