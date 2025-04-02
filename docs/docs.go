// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/url": {
            "post": {
                "description": "Creates a short alias for the provided URL",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "url"
                ],
                "summary": "Save URL",
                "parameters": [
                    {
                        "description": "URL shortening request data",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_url_save.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal_http-server_handlers_url_save.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_popvaleks_url-shortener_internal_lib_api_response.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_popvaleks_url-shortener_internal_lib_api_response.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "internal_http-server_handlers_url_save.Request": {
            "description": "Request to create a short URL",
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "alias": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "internal_http-server_handlers_url_save.Response": {
            "description": "Response with the created alias",
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Url shortener",
	Description:      "Shortener service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
