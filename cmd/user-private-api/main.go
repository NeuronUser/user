package main

import (
	"github.com/NeuronFramework/restful"
	"github.com/NeuronUser/user/api/gen/restapi"
	"github.com/NeuronUser/user/api/gen/restapi/operations"
	"github.com/NeuronUser/user/cmd/user-private-api/handler"
	"github.com/go-openapi/loads"
	"net/http"
)

func main() {
	restful.Run(func() (http.Handler, error) {
		h, err := handler.NewUserHandler()
		if err != nil {
			return nil, err
		}

		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return nil, err
		}

		api := operations.NewUserAPI(swaggerSpec)
		api.BearerAuth = h.BearerAuth
		api.GetUserInfoHandler = operations.GetUserInfoHandlerFunc(h.GetUserInfo)

		return api.Serve(nil), nil
	})
}
