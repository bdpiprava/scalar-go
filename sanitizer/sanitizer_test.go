package sanitizer_test

import (
	"testing"

	"github.com/bdpiprava/scalar-go/model"
	"github.com/bdpiprava/scalar-go/sanitizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *model.Spec
		validate func(t *testing.T, got *model.Spec)
	}{
		{
			name: "should sanitize spec with map[any]any in schemas",
			input: &model.Spec{
				Components: model.Components{
					Schemas: map[string]any{
						"User": map[any]any{
							"type": "object",
							123:    "numeric_key",
							true:   "bool_key",
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				userSchema := got.Components.Schemas["User"].(map[string]any)
				assert.Equal(t, "object", userSchema["type"])
				assert.Equal(t, "numeric_key", userSchema["123"])
				assert.Equal(t, "bool_key", userSchema["true"])
			},
		},
		{
			name: "should sanitize spec with map[any]any in paths",
			input: &model.Spec{
				Paths: map[string]any{
					"/users": map[any]any{
						"get": map[any]any{
							"summary":     "Get users",
							"operationId": "getUsers",
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				usersPath := got.Paths["/users"].(map[string]any)
				getOp := usersPath["get"].(map[string]any)
				assert.Equal(t, "Get users", getOp["summary"])
				assert.Equal(t, "getUsers", getOp["operationId"])
			},
		},
		{
			name: "should sanitize spec with nested map[any]any structures",
			input: &model.Spec{
				Components: model.Components{
					Schemas: map[string]any{
						"Pet": map[any]any{
							"properties": map[any]any{
								"id": map[any]any{
									"type":   "integer",
									"format": "int64",
								},
								"name": map[any]any{
									"type": "string",
								},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				petSchema := got.Components.Schemas["Pet"].(map[string]any)
				properties := petSchema["properties"].(map[string]any)
				idProp := properties["id"].(map[string]any)
				assert.Equal(t, "integer", idProp["type"])
				assert.Equal(t, "int64", idProp["format"])
			},
		},
		{
			name: "should sanitize spec with arrays containing map[any]any",
			input: &model.Spec{
				Paths: map[string]any{
					"/pets": map[any]any{
						"get": map[any]any{
							"parameters": []any{
								map[any]any{
									"name": "limit",
									"in":   "query",
								},
								map[any]any{
									"name": "offset",
									"in":   "query",
								},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				petsPath := got.Paths["/pets"].(map[string]any)
				getOp := petsPath["get"].(map[string]any)
				params := getOp["parameters"].([]any)
				require.Len(t, params, 2)

				param0 := params[0].(map[string]any)
				assert.Equal(t, "limit", param0["name"])
				assert.Equal(t, "query", param0["in"])

				param1 := params[1].(map[string]any)
				assert.Equal(t, "offset", param1["name"])
			},
		},
		{
			name: "should handle empty spec",
			input: &model.Spec{
				Components: model.Components{
					Schemas:    map[string]any{},
					Parameters: map[string]any{},
				},
				Paths: map[string]any{},
			},
			validate: func(t *testing.T, got *model.Spec) {
				assert.NotNil(t, got.Components.Schemas)
				assert.NotNil(t, got.Components.Parameters)
				assert.NotNil(t, got.Paths)
				assert.Empty(t, got.Components.Schemas)
				assert.Empty(t, got.Components.Parameters)
				assert.Empty(t, got.Paths)
			},
		},
		{
			name: "should preserve string keys unchanged",
			input: &model.Spec{
				Components: model.Components{
					Schemas: map[string]any{
						"User": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"id":   map[string]any{"type": "integer"},
								"name": map[string]any{"type": "string"},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				userSchema := got.Components.Schemas["User"].(map[string]any)
				assert.Equal(t, "object", userSchema["type"])
				properties := userSchema["properties"].(map[string]any)
				assert.Contains(t, properties, "id")
				assert.Contains(t, properties, "name")
			},
		},
		{
			name: "should sanitize parameters with map[any]any",
			input: &model.Spec{
				Components: model.Components{
					Parameters: map[string]any{
						"limitParam": map[any]any{
							"name":     "limit",
							"in":       "query",
							"required": false,
							"schema":   map[any]any{"type": "integer"},
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				limitParam := got.Components.Parameters["limitParam"].(map[string]any)
				assert.Equal(t, "limit", limitParam["name"])
				assert.Equal(t, "query", limitParam["in"])
				assert.Equal(t, false, limitParam["required"])
				schema := limitParam["schema"].(map[string]any)
				assert.Equal(t, "integer", schema["type"])
			},
		},
		{
			name: "should handle mixed types in arrays",
			input: &model.Spec{
				Paths: map[string]any{
					"/test": map[any]any{
						"post": map[any]any{
							"tags": []any{"users", "admin", 123, true},
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				testPath := got.Paths["/test"].(map[string]any)
				postOp := testPath["post"].(map[string]any)
				tags := postOp["tags"].([]any)
				assert.Equal(t, "users", tags[0])
				assert.Equal(t, "admin", tags[1])
				assert.Equal(t, 123, tags[2])
				assert.Equal(t, true, tags[3])
			},
		},
		{
			name: "should handle nil values in map",
			input: &model.Spec{
				Components: model.Components{
					Schemas: map[string]any{
						"Test": map[any]any{
							"nullable":    true,
							"description": nil,
							"default":     nil,
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				testSchema := got.Components.Schemas["Test"].(map[string]any)
				assert.Equal(t, true, testSchema["nullable"])
				assert.Nil(t, testSchema["description"])
				assert.Nil(t, testSchema["default"])
			},
		},
		{
			name: "should handle deeply nested structures",
			input: &model.Spec{
				Paths: map[string]any{
					"/api": map[any]any{
						"post": map[any]any{
							"responses": map[any]any{
								"200": map[any]any{
									"content": map[any]any{
										"application/json": map[any]any{
											"schema": map[any]any{
												"$ref": "#/components/schemas/Response",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, got *model.Spec) {
				apiPath := got.Paths["/api"].(map[string]any)
				postOp := apiPath["post"].(map[string]any)
				responses := postOp["responses"].(map[string]any)
				response200 := responses["200"].(map[string]any)
				content := response200["content"].(map[string]any)
				jsonContent := content["application/json"].(map[string]any)
				schema := jsonContent["schema"].(map[string]any)
				assert.Equal(t, "#/components/schemas/Response", schema["$ref"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := sanitizer.Sanitize(tc.input)

			require.NotNil(t, got)
			tc.validate(t, got)
		})
	}
}

func TestSanitize_ComplexRealWorldScenario(t *testing.T) {
	t.Parallel()

	// Simulate a spec that might come from YAML unmarshaling with map[any]any
	input := &model.Spec{
		OpenAPI: "3.0.0",
		Info: model.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]any{
			"/users": map[any]any{
				"get": map[any]any{
					"summary":     "List users",
					"operationId": "listUsers",
					"parameters": []any{
						map[any]any{
							"name":     "limit",
							"in":       "query",
							"required": false,
							"schema": map[any]any{
								"type":    "integer",
								"default": 10,
							},
						},
					},
					"responses": map[any]any{
						"200": map[any]any{
							"description": "Success",
							"content": map[any]any{
								"application/json": map[any]any{
									"schema": map[any]any{
										"type": "array",
										"items": map[any]any{
											"$ref": "#/components/schemas/User",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Components: model.Components{
			Schemas: map[string]any{
				"User": map[any]any{
					"type":     "object",
					"required": []any{"id", "name"},
					"properties": map[any]any{
						"id": map[any]any{
							"type":   "integer",
							"format": "int64",
						},
						"name": map[any]any{
							"type": "string",
						},
						"email": map[any]any{
							"type":   "string",
							"format": "email",
						},
					},
				},
			},
			Parameters: map[string]any{
				"limitParam": map[any]any{
					"name": "limit",
					"in":   "query",
					"schema": map[any]any{
						"type":    "integer",
						"default": 10,
					},
				},
			},
		},
	}

	got := sanitizer.Sanitize(input)

	require.NotNil(t, got)

	// Verify paths are sanitized
	usersPath := got.Paths["/users"].(map[string]any)
	getOp := usersPath["get"].(map[string]any)
	assert.Equal(t, "List users", getOp["summary"])

	// Verify nested parameters
	params := getOp["parameters"].([]any)
	require.Len(t, params, 1)
	param := params[0].(map[string]any)
	assert.Equal(t, "limit", param["name"])

	// Verify responses
	responses := getOp["responses"].(map[string]any)
	response200 := responses["200"].(map[string]any)
	assert.Equal(t, "Success", response200["description"])

	// Verify schemas
	userSchema := got.Components.Schemas["User"].(map[string]any)
	assert.Equal(t, "object", userSchema["type"])
	properties := userSchema["properties"].(map[string]any)
	idProp := properties["id"].(map[string]any)
	assert.Equal(t, "integer", idProp["type"])
	assert.Equal(t, "int64", idProp["format"])

	// Verify parameters
	limitParam := got.Components.Parameters["limitParam"].(map[string]any)
	assert.Equal(t, "limit", limitParam["name"])
	assert.Equal(t, "query", limitParam["in"])
}

func TestSanitize_PreservesOriginalSpec(t *testing.T) {
	t.Parallel()

	input := &model.Spec{
		OpenAPI: "3.0.0",
		Info: model.Info{
			Title:   "Original Title",
			Version: "1.0.0",
		},
		Components: model.Components{
			Schemas: map[string]any{
				"Test": map[any]any{
					"type": "object",
				},
			},
		},
	}

	got := sanitizer.Sanitize(input)

	// Verify original spec metadata is preserved
	assert.Equal(t, "3.0.0", got.OpenAPI)
	assert.Equal(t, "Original Title", got.Info.Title)
	assert.Equal(t, "1.0.0", got.Info.Version)

	// Verify returned spec is the same instance (sanitize modifies in place)
	assert.Same(t, input, got)
}
