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
        "/db/echoDB/check-update": {
            "get": {
                "description": "根据客户端提供的当前版本号，检查是否需要更新。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "update"
                ],
                "summary": "检查更新",
                "parameters": [
                    {
                        "type": "string",
                        "default": "v1.0.0",
                        "description": "当前版本号",
                        "name": "current_version",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功返回最新版本信息",
                        "schema": {
                            "$ref": "#/definitions/db.Response"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/db.Response"
                        }
                    }
                }
            }
        },
        "/db/echoDB/update-version": {
            "post": {
                "description": "更新服务器上记录的最新版本号。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "update"
                ],
                "summary": "更新应用版本",
                "parameters": [
                    {
                        "description": "新的版本号",
                        "name": "new_version",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "版本更新成功",
                        "schema": {
                            "$ref": "#/definitions/db.Response"
                        }
                    },
                    "400": {
                        "description": "无效的输入数据",
                        "schema": {
                            "$ref": "#/definitions/db.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {
                    "$ref": "#/definitions/db.UpdateData"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "db.UpdateData": {
            "type": "object",
            "properties": {
                "current_version": {
                    "type": "string"
                },
                "needs_update": {
                    "type": "boolean"
                },
                "newest_version": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
