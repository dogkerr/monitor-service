// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Lintang BS"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/monitors/dashboards/logs": {
            "get": {
                "description": "GetUserLogsDashboard",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "monitor"
                ],
                "summary": "Mendapatkan Dashboard Logs containers milik User",
                "operationId": "logs_dashboard",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/rest.logsDashboardRes"
                        }
                    },
                    "500": {
                        "description": "internal server error (bug/error di kode)",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        },
        "/monitors/dashboards/monitors": {
            "get": {
                "description": "GetUserMonitorDashboard",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "monitor"
                ],
                "summary": "Mendapatkan Dashboard Container metrics milik User",
                "operationId": "monitor_dashboard",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/rest.dashboardRes"
                        }
                    },
                    "500": {
                        "description": "internal server error (bug/error di kode)",
                        "schema": {
                            "$ref": "#/definitions/rest.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Dashboard": {
            "description": "ini data dashboard (isinya id, owner, uid, type)",
            "type": "object",
            "properties": {
                "id": {
                    "description": "id dashboard di database",
                    "type": "string"
                },
                "owner": {
                    "description": "owner /pemilik dashboard",
                    "type": "string"
                },
                "type": {
                    "description": "type dashboard",
                    "type": "string"
                },
                "uid": {
                    "description": "uid dashboard di grafana",
                    "type": "string"
                }
            }
        },
        "rest.ResponseError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "rest.dashboardRes": {
            "description": "Response saat get metrics dashboard milik user",
            "type": "object",
            "properties": {
                "cpu_usage_link": {
                    "description": "link dashboard cpu usage per contaainer",
                    "type": "string"
                },
                "dashboard": {
                    "description": "data dashboard milik user (isinya uid, owner, type, id)",
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.Dashboard"
                        }
                    ]
                },
                "memory_swap_per_container_link": {
                    "description": "link memory swap per container",
                    "type": "string"
                },
                "memory_usage_not_graph": {
                    "description": "link memory usage per container gak pake graph cuma angka",
                    "type": "string"
                },
                "memory_usage_per_container_link": {
                    "description": "link memory usage per container pake graph",
                    "type": "string"
                },
                "overall_cpu_usage": {
                    "description": "link overal cpu usage untuk semua container milik user",
                    "type": "string"
                },
                "received_network_link": {
                    "description": "link dashboard metrics received network per contaainer",
                    "type": "string",
                    "example": "http://127.0.0.1/d-solo/VWUxxYP3/vwuxxyp3?orgId=1\u0026refresh=5s\u0026from=now-5m\u0026theme=light\u0026to=now\u0026panelId=8"
                },
                "send_network_link": {
                    "description": "link dashboard metrics send networks per contaainer",
                    "type": "string"
                },
                "total_container": {
                    "description": "jumlah container yang dijalankan user di dogker",
                    "type": "string"
                }
            }
        },
        "rest.logsDashboardRes": {
            "description": "Response saat get logs dashboard milik user",
            "type": "object",
            "properties": {
                "dashboard": {
                    "description": "data dashboard milik user (isinya uid, owner, type, id)",
                    "allOf": [
                        {
                            "$ref": "#/definitions/domain.Dashboard"
                        }
                    ]
                },
                "logs_dashboard_link": {
                    "description": "link dashboard logs yang diembed di frontend",
                    "type": "string",
                    "example": "http://localhost:3000/d/YwXYwNAj/ywxywnaj?orgId=1\u0026var-search_filter=\u0026var-Levels=info\u0026var-container_name=go_container_log2\u0026var-Method=GET\u0026from=1714796971638\u0026to=1714797271638\u0026theme=light"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "locahost:9191",
	BasePath:         "/api/v1/",
	Schemes:          []string{},
	Title:            "Dogker Monitor Service",
	Description:      "Monitor Servicee buat nampilin logs dashboard & conainer metrics dashboard milik user",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
