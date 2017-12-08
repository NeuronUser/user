// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

// SwaggerJSON embedded version of the swagger document used at generation time
var SwaggerJSON json.RawMessage

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http",
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "title": "User Private API",
    "contact": {
      "name": "mars"
    },
    "version": "v1"
  },
  "basePath": "/api-private/v1/users",
  "paths": {
    "/logout": {
      "post": {
        "operationId": "Logout",
        "parameters": [
          {
            "type": "string",
            "name": "jwt",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "ok"
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/logoutDefaultBody"
            }
          }
        }
      }
    },
    "/oauth/jump": {
      "post": {
        "operationId": "OauthJump",
        "parameters": [
          {
            "type": "string",
            "name": "authorizationCode",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "state",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "jwt",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/oauthJumpDefaultBody"
            }
          }
        }
      }
    },
    "/oauth/state": {
      "get": {
        "operationId": "GetOauthState",
        "responses": {
          "200": {
            "description": "state",
            "schema": {
              "type": "string"
            }
          },
          "default": {
            "description": "Error response",
            "schema": {
              "$ref": "#/definitions/getOauthStateDefaultBody"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "getOauthStateDefaultBody": {
      "type": "object",
      "properties": {
        "code": {
          "description": "Error code",
          "type": "string"
        },
        "errors": {
          "$ref": "#/definitions/getOauthStateDefaultBodyErrors"
        },
        "message": {
          "description": "Error message",
          "type": "string"
        },
        "status": {
          "type": "string",
          "format": "int32",
          "default": "Http status"
        }
      },
      "x-go-gen-location": "operations"
    },
    "getOauthStateDefaultBodyErrors": {
      "description": "Errors",
      "type": "array",
      "items": {
        "$ref": "#/definitions/getOauthStateDefaultBodyErrorsItems"
      },
      "x-go-gen-location": "operations"
    },
    "getOauthStateDefaultBodyErrorsItems": {
      "type": "object",
      "properties": {
        "code": {
          "description": "error code",
          "type": "string"
        },
        "field": {
          "description": "field name",
          "type": "string"
        },
        "message": {
          "description": "error message",
          "type": "string"
        }
      },
      "x-go-gen-location": "operations"
    },
    "logoutDefaultBody": {
      "type": "object",
      "properties": {
        "code": {
          "description": "Error code",
          "type": "string"
        },
        "errors": {
          "$ref": "#/definitions/getOauthStateDefaultBodyErrors"
        },
        "message": {
          "description": "Error message",
          "type": "string"
        },
        "status": {
          "type": "string",
          "format": "int32",
          "default": "Http status"
        }
      },
      "x-go-gen-location": "operations"
    },
    "logoutDefaultBodyErrors": {
      "description": "Errors",
      "type": "array",
      "items": {
        "$ref": "#/definitions/getOauthStateDefaultBodyErrorsItems"
      },
      "x-go-gen-location": "operations"
    },
    "logoutDefaultBodyErrorsItems": {
      "type": "object",
      "properties": {
        "code": {
          "description": "error code",
          "type": "string"
        },
        "field": {
          "description": "field name",
          "type": "string"
        },
        "message": {
          "description": "error message",
          "type": "string"
        }
      },
      "x-go-gen-location": "operations"
    },
    "oauthJumpDefaultBody": {
      "type": "object",
      "properties": {
        "code": {
          "description": "Error code",
          "type": "string"
        },
        "errors": {
          "$ref": "#/definitions/getOauthStateDefaultBodyErrors"
        },
        "message": {
          "description": "Error message",
          "type": "string"
        },
        "status": {
          "type": "string",
          "format": "int32",
          "default": "Http status"
        }
      },
      "x-go-gen-location": "operations"
    },
    "oauthJumpDefaultBodyErrors": {
      "description": "Errors",
      "type": "array",
      "items": {
        "$ref": "#/definitions/getOauthStateDefaultBodyErrorsItems"
      },
      "x-go-gen-location": "operations"
    },
    "oauthJumpDefaultBodyErrorsItems": {
      "type": "object",
      "properties": {
        "code": {
          "description": "error code",
          "type": "string"
        },
        "field": {
          "description": "field name",
          "type": "string"
        },
        "message": {
          "description": "error message",
          "type": "string"
        }
      },
      "x-go-gen-location": "operations"
    }
  },
  "responses": {
    "ErrorResponse": {
      "description": "Error response",
      "schema": {
        "type": "object",
        "properties": {
          "code": {
            "description": "Error code",
            "type": "string"
          },
          "errors": {
            "$ref": "#/definitions/getOauthStateDefaultBodyErrors"
          },
          "message": {
            "description": "Error message",
            "type": "string"
          },
          "status": {
            "type": "string",
            "format": "int32",
            "default": "Http status"
          }
        }
      }
    }
  },
  "securityDefinitions": {
    "Basic": {
      "type": "basic"
    }
  }
}`))
}