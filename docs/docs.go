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
        "/": {
            "get": {
                "description": "Check if the application is running and connected to the database and cache",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "default"
                ],
                "summary": "Check if the application is running",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/auth/token": {
            "post": {
                "description": "login in the application and get a token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "default"
                ],
                "summary": "login in the application",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_Jonattas-21_loan-engine_internal_api_dto.TokenResponse_dto"
                        }
                    }
                }
            }
        },
        "/v1/loanconditions": {
            "get": {
                "description": "Get all conditions",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "conditions"
                ],
                "summary": "Show the list of loan conditions, fees by age group",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_Jonattas-21_loan-engine_internal_domain_entities.LoanCondition"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "update a loan condition by name tier1, tier2, tier3, tier4",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "conditions"
                ],
                "summary": "update a loan condition by name",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/loansimulations": {
            "post": {
                "description": "Get a plenty of loan simulations",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "simulation"
                ],
                "summary": "Get a plenty of loan simulations",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_Jonattas-21_loan-engine_internal_domain_entities.LoanSimulation"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_Jonattas-21_loan-engine_internal_api_dto.TokenResponse_dto": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "integer"
                },
                "refresh_token": {
                    "type": "string"
                },
                "token_type": {
                    "type": "string"
                }
            }
        },
        "github_com_Jonattas-21_loan-engine_internal_domain_entities.Installment": {
            "type": "object",
            "properties": {
                "currency": {
                    "type": "string"
                },
                "installment_amount": {
                    "type": "number"
                },
                "installment_fee_amount": {
                    "type": "number"
                },
                "installment_number": {
                    "type": "integer"
                }
            }
        },
        "github_com_Jonattas-21_loan-engine_internal_domain_entities.LoanCondition": {
            "type": "object",
            "properties": {
                "interest_rate": {
                    "type": "number"
                },
                "max_age": {
                    "type": "integer"
                },
                "min_age": {
                    "type": "integer"
                },
                "modified_date": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "github_com_Jonattas-21_loan-engine_internal_domain_entities.LoanSimulation": {
            "type": "object",
            "properties": {
                "amount_fee_to_be_paid": {
                    "type": "number"
                },
                "amount_to_be_paid": {
                    "type": "number"
                },
                "currency": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fee_amount_percentage": {
                    "type": "number"
                },
                "installments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_Jonattas-21_loan-engine_internal_domain_entities.Installment"
                    }
                },
                "loan_amount": {
                    "type": "number"
                },
                "simulation_date": {
                    "type": "string"
                },
                "total_installments": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Loan Engine API",
	Description:      "This project It's a credit simulator which allows users to consult loan conditions, based in some payments conditions.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
