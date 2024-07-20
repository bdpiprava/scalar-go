package loader_test

import (
	"reflect"
	"testing"

	"github.com/bdpiprava/scalar-go/loader"
	"github.com/bdpiprava/scalar-go/model"
	"github.com/stretchr/testify/assert"
)

func Test_Load(t *testing.T) {
	spec, err := loader.LoadWithName("../data/loader", "pet-store.yml")

	assert.NoError(t, err)
	assert.NotNil(t, spec)
	assertBase(t, spec)
	assertSchema(t, spec)
	assertPaths(t, spec)
}

func Test_Load_MultipleFiles(t *testing.T) {
	spec, err := loader.LoadWithName("../data/loader-multiple-files", "api.yml")

	assert.NoError(t, err)
	assert.NotNil(t, spec)
	assertBase(t, spec)
	assertSchema(t, spec)
	assertPaths(t, spec)
}

func Test_LoadedFileShouldHaveIdenticalContent(t *testing.T) {
	specFromMultipleFiles, err := loader.LoadWithName("../data/loader-multiple-files", "api.yml")
	assert.NoError(t, err)

	specFromSingleFile, err := loader.LoadWithName("../data/loader", "pet-store.yml")
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(specFromMultipleFiles, specFromSingleFile))
}

func Test_Load_DocumentedPath(t *testing.T) {
	spec, err := loader.LoadWithName("../data/loader", "pet-store.yml")
	assert.NoError(t, err)

	paths := spec.DocumentedPaths()

	assert.Len(t, paths, 3)
	assert.Equal(t, model.DocumentedPath{Path: "/pets", Method: "get"}, paths[0])
	assert.Equal(t, model.DocumentedPath{Path: "/pets", Method: "post"}, paths[1])
	assert.Equal(t, model.DocumentedPath{Path: "/pets/{petId}", Method: "get"}, paths[2])
}

func getPathParam(params []interface{}, name string) model.GenericObject {
	for _, param := range params {
		if param.(model.GenericObject)["name"] == name {
			return param.(model.GenericObject)
		}
	}
	return nil
}

func assertBase(t *testing.T, spec *model.Spec) {
	assert.Equal(t, "3.0.0", spec.OpenAPI)

	assert.Equal(t, "1.0.0", spec.Info.Version)
	assert.Equal(t, "Swagger Petstore", spec.Info.Title)
	assert.Equal(t, "MIT", spec.Info.License.Name)

	assert.Len(t, spec.Servers, 1)
	assert.Equal(t, model.Server{URL: "http://petstore.swagger.io/v1"}, spec.Servers[0])

	assert.Len(t, spec.Tags, 2)
	assert.Equal(t, model.Tag{Name: "pets", Description: "Everything about your Pets"}, spec.Tags[0])
	assert.Equal(t, model.Tag{Name: "store", Description: "Access to Petstore"}, spec.Tags[1])
}

func assertPaths(t *testing.T, spec *model.Spec) {
	assert.Len(t, spec.Paths, 2)
	pathGetPets := spec.Paths["/pets"].(model.GenericObject)["get"].(model.GenericObject)
	assert.Equal(t, "List all pets", pathGetPets["summary"])
	assert.Equal(t, "listPets", pathGetPets["operationId"])
	assert.Equal(t, []interface{}{"pets"}, pathGetPets["tags"])

	nameParam := getPathParam(pathGetPets["parameters"].([]interface{}), "limit")
	assert.Equal(t, model.GenericObject{
		"description": "How many items to return at one time (max 100)",
		"in":          "query",
		"name":        "limit",
		"required":    false,
		"schema": model.GenericObject{
			"format":  "int32",
			"maximum": 100,
			"type":    "integer",
		},
	}, nameParam)

	assert.Equal(t, model.GenericObject{
		"200": model.GenericObject{
			"content": model.GenericObject{
				"application/json": model.GenericObject{
					"schema": model.GenericObject{
						"$ref": "#/components/schemas/Pets",
					},
				},
			},
			"description": "A paged array of pets",
			"headers": model.GenericObject{
				"x-next": model.GenericObject{
					"description": "A link to the next page of responses",
					"schema": model.GenericObject{
						"type": "string",
					},
				},
			},
		},
		"default": model.GenericObject{
			"content": model.GenericObject{
				"application/json": model.GenericObject{
					"schema": model.GenericObject{
						"$ref": "#/components/schemas/Error",
					},
				},
			},
			"description": "unexpected error",
		},
	}, pathGetPets["responses"])
}

func assertSchema(t *testing.T, spec *model.Spec) {
	assert.Len(t, spec.Components.Schemas, 3)
	petSchema := spec.Components.Schemas["Pet"].(model.GenericObject)
	assert.NotNil(t, petSchema)
	assert.Equal(t, "object", petSchema["type"])
	assert.Equal(t, []interface{}{"id", "name"}, petSchema["required"])
	assert.Equal(t, model.GenericObject{"format": "int64", "type": "integer"}, petSchema["properties"].(model.GenericObject)["id"])
	assert.Equal(t, model.GenericObject{"type": "string"}, petSchema["properties"].(model.GenericObject)["name"])
	assert.Equal(t, model.GenericObject{"type": "string"}, petSchema["properties"].(model.GenericObject)["tag"])

	petsSchema := spec.Components.Schemas["Pets"].(model.GenericObject)
	assert.NotNil(t, petsSchema)
	assert.Equal(t, "array", petsSchema["type"])
	assert.Equal(t, 100, petsSchema["maxItems"])
	assert.Equal(t, "#/components/schemas/Pet", petsSchema["items"].(model.GenericObject)["$ref"])

	errorSchema := spec.Components.Schemas["Error"].(model.GenericObject)
	assert.NotNil(t, errorSchema)
	assert.Equal(t, "object", errorSchema["type"])
	assert.Equal(t, []interface{}{"code", "message"}, errorSchema["required"])
	assert.Equal(t, model.GenericObject{"type": "integer", "format": "int32"}, errorSchema["properties"].(model.GenericObject)["code"])
	assert.Equal(t, model.GenericObject{"type": "string"}, errorSchema["properties"].(model.GenericObject)["message"])
}
