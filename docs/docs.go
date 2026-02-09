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
        "/subscriptions": {
            "get": {
                "description": "Get a list of all subscriptions",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Get all subscriptions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.SubscriptionResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new subscription record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Create a new subscription",
                "parameters": [
                    {
                        "description": "Subscription data",
                        "name": "subscription",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SubscriptionRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SubscriptionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/subscriptions/aggregate": {
            "get": {
                "description": "Calculate total price of subscriptions with optional filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Aggregate total price of subscriptions",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID to filter by",
                        "name": "user_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Service name to filter by",
                        "name": "service_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date (MM-YYYY) to filter by",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date (MM-YYYY) to filter by",
                        "name": "end_date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AggregateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/subscriptions/{id}": {
            "get": {
                "description": "Get subscription details by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Get a subscription by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SubscriptionResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Update subscription details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Update a subscription",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated subscription data",
                        "name": "subscription",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SubscriptionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SubscriptionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a subscription by ID",
                "tags": [
                    "subscriptions"
                ],
                "summary": "Delete a subscription",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AggregateResponse": {
            "type": "object",
            "properties": {
                "total_price": {
                    "type": "integer"
                }
            }
        },
        "models.SubscriptionRequest": {
            "type": "object",
            "required": [
                "price",
                "service_name",
                "start_date",
                "user_id"
            ],
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "price": {
                    "type": "integer",
                    "minimum": 0
                },
                "service_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "models.SubscriptionResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "service_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        }
    }
}`

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
