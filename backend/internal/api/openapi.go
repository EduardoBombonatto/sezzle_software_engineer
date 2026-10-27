package api

import (
	"encoding/json"
	"net/http"
)

const (
	ContentTypeHTML = "text/html; charset=utf-8"
	OpenAPIVersion  = "3.0.3"
)

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Sezzle Calculator API — OpenAPI Documentation</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
  <style>
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: '/api/v1/openapi.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis
        ]
      });
    };
  </script>
</body>
</html>`

func (h *Handler) handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed; use GET")
		return
	}

	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(buildOpenAPISpec())
}

func (h *Handler) handleOpenAPIDocs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed; use GET")
		return
	}

	w.Header().Set(ContentTypeHeader, ContentTypeHTML)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerUIHTML))
}

func buildOpenAPISpec() map[string]interface{} {
	return map[string]interface{}{
		"openapi": OpenAPIVersion,
		"info": map[string]interface{}{
			"title":       "Sezzle Full-Stack Calculator API",
			"version":     "1.0.0",
			"description": "RESTful microservice for basic and advanced arithmetic operations with strict pre-calculation validation and standardized JSON response envelopes.",
		},
		"servers": []map[string]interface{}{
			{
				"url":         "/",
				"description": "Current API Server",
			},
		},
		"paths": map[string]interface{}{
			"/api/v1/health": buildHealthPathItem(),
			"/api/v1/calculate/add": buildBinaryOperationPathItem(
				"Addition (a + b)",
				"Computes the sum of two finite float64 numbers.",
				OperationAdd,
				12.5, 7.5, 20.0,
			),
			"/api/v1/calculate/subtract": buildBinaryOperationPathItem(
				"Subtraction (a - b)",
				"Computes the difference between two finite float64 numbers.",
				OperationSubtract,
				25.0, 10.0, 15.0,
			),
			"/api/v1/calculate/multiply": buildBinaryOperationPathItem(
				"Multiplication (a * b)",
				"Computes the product of two finite float64 numbers.",
				OperationMultiply,
				6.0, 7.0, 42.0,
			),
			"/api/v1/calculate/divide": buildBinaryOperationPathItem(
				"Division (a / b)",
				"Computes the quotient of two finite float64 numbers. Fails fast with 422 if b is zero.",
				OperationDivide,
				45.0, 6.0, 7.5,
			),
			"/api/v1/calculate/power": buildBinaryOperationPathItem(
				"Exponentiation (a ^ b)",
				"Computes base a raised to exponent b. Rejects 0^(negative) and negative bases with fractional exponents.",
				OperationPower,
				2.0, 8.0, 256.0,
			),
			"/api/v1/calculate/sqrt": buildUnaryOperationPathItem(
				"Square Root (√a)",
				"Computes the non-negative square root of operand a. Rejects negative numbers and any second operand b.",
				OperationSqrt,
				81.0, 9.0,
			),
			"/api/v1/calculate/percentage": buildBinaryOperationPathItem(
				"Percentage (a% of b)",
				"Computes a percent of b ((a * b) / 100) while preventing spurious intermediate float64 overflow.",
				OperationPercentage,
				20.0, 150.0, 30.0,
			),
		},
		"components": map[string]interface{}{
			"schemas": buildOpenAPISchemas(),
		},
	}
}

func buildHealthPathItem() map[string]interface{} {
	return map[string]interface{}{
		"get": map[string]interface{}{
			"tags":        []string{"System"},
			"summary":     "Health Check",
			"description": "Returns the operational status of the calculator microservice.",
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Service is healthy",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"example": map[string]interface{}{
								"erro":      false,
								"message":   "Service is healthy",
								"data":      map[string]interface{}{"status": "healthy"},
								"httpcode":  200,
								"timestamp": "2026-10-26T21:00:00Z",
							},
						},
					},
				},
			},
		},
	}
}

func buildBinaryOperationPathItem(
	summary string,
	description string,
	operation string,
	exampleA float64,
	exampleB float64,
	exampleResult float64,
) map[string]interface{} {
	return map[string]interface{}{
		"post": map[string]interface{}{
			"tags":        []string{"Calculator Operations"},
			"summary":     summary,
			"description": description,
			"requestBody": map[string]interface{}{
				"required": true,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"$ref": "#/components/schemas/BinaryCalculationRequest",
						},
						"example": map[string]interface{}{
							"a": exampleA,
							"b": exampleB,
						},
					},
				},
			},
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Calculation succeeded",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/CalculationSuccessEnvelope",
							},
							"example": map[string]interface{}{
								"erro":    false,
								"message": "Operation completed successfully",
								"data": map[string]interface{}{
									"operation": operation,
									"a":         exampleA,
									"b":         exampleB,
									"result":    exampleResult,
								},
								"httpcode":  200,
								"timestamp": "2026-10-26T21:00:00Z",
							},
						},
					},
				},
				"400": buildErrorResponseRef("Invalid request payload or missing required operands"),
				"422": buildErrorResponseRef("Mathematical domain error (e.g. division by zero or float64 overflow)"),
			},
		},
	}
}

func buildUnaryOperationPathItem(
	summary string,
	description string,
	operation string,
	exampleA float64,
	exampleResult float64,
) map[string]interface{} {
	return map[string]interface{}{
		"post": map[string]interface{}{
			"tags":        []string{"Calculator Operations"},
			"summary":     summary,
			"description": description,
			"requestBody": map[string]interface{}{
				"required": true,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"$ref": "#/components/schemas/UnaryCalculationRequest",
						},
						"example": map[string]interface{}{
							"a": exampleA,
						},
					},
				},
			},
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Calculation succeeded",
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]interface{}{
								"$ref": "#/components/schemas/CalculationSuccessEnvelope",
							},
							"example": map[string]interface{}{
								"erro":    false,
								"message": "Operation completed successfully",
								"data": map[string]interface{}{
									"operation": operation,
									"a":         exampleA,
									"result":    exampleResult,
								},
								"httpcode":  200,
								"timestamp": "2026-10-26T21:00:00Z",
							},
						},
					},
				},
				"400": buildErrorResponseRef("Invalid request payload, missing operand 'a', or unexpected operand 'b'"),
				"422": buildErrorResponseRef("Mathematical domain error (e.g. negative square root)"),
			},
		},
	}
}

func buildErrorResponseRef(description string) map[string]interface{} {
	return map[string]interface{}{
		"description": description,
		"content": map[string]interface{}{
			"application/json": map[string]interface{}{
				"schema": map[string]interface{}{
					"$ref": "#/components/schemas/ErrorEnvelope",
				},
			},
		},
	}
}

func buildOpenAPISchemas() map[string]interface{} {
	return map[string]interface{}{
		"BinaryCalculationRequest": map[string]interface{}{
			"type":                 "object",
			"required":             []string{"a", "b"},
			"additionalProperties": false,
			"properties": map[string]interface{}{
				"a": map[string]interface{}{
					"type":        "number",
					"format":      "double",
					"description": "First finite float64 operand",
				},
				"b": map[string]interface{}{
					"type":        "number",
					"format":      "double",
					"description": "Second finite float64 operand",
				},
			},
		},
		"UnaryCalculationRequest": map[string]interface{}{
			"type":                 "object",
			"required":             []string{"a"},
			"additionalProperties": false,
			"properties": map[string]interface{}{
				"a": map[string]interface{}{
					"type":        "number",
					"format":      "double",
					"description": "Single finite float64 operand",
				},
			},
		},
		"CalculationData": map[string]interface{}{
			"type":     "object",
			"required": []string{"operation", "a", "result"},
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
					"enum": []string{
						OperationAdd,
						OperationSubtract,
						OperationMultiply,
						OperationDivide,
						OperationPower,
						OperationSqrt,
						OperationPercentage,
					},
				},
				"a": map[string]interface{}{
					"type":   "number",
					"format": "double",
				},
				"b": map[string]interface{}{
					"type":   "number",
					"format": "double",
				},
				"result": map[string]interface{}{
					"type":   "number",
					"format": "double",
				},
			},
		},
		"CalculationSuccessEnvelope": map[string]interface{}{
			"type":     "object",
			"required": []string{"erro", "message", "data", "httpcode", "timestamp"},
			"properties": map[string]interface{}{
				"erro": map[string]interface{}{
					"type":    "boolean",
					"example": false,
				},
				"message": map[string]interface{}{
					"type":    "string",
					"example": "Operation completed successfully",
				},
				"data": map[string]interface{}{
					"$ref": "#/components/schemas/CalculationData",
				},
				"httpcode": map[string]interface{}{
					"type":    "integer",
					"example": 200,
				},
				"timestamp": map[string]interface{}{
					"type":   "string",
					"format": "date-time",
				},
			},
		},
		"ErrorEnvelope": map[string]interface{}{
			"type":     "object",
			"required": []string{"erro", "message", "data", "httpcode", "timestamp"},
			"properties": map[string]interface{}{
				"erro": map[string]interface{}{
					"type":    "boolean",
					"example": true,
				},
				"message": map[string]interface{}{
					"type":    "string",
					"example": "division by zero is not allowed",
				},
				"data": map[string]interface{}{
					"type":     "object",
					"nullable": true,
					"example":  nil,
				},
				"httpcode": map[string]interface{}{
					"type":    "integer",
					"example": 422,
				},
				"timestamp": map[string]interface{}{
					"type":   "string",
					"format": "date-time",
				},
			},
		},
	}
}
