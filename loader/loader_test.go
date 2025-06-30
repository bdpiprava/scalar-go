package loader_test

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bdpiprava/scalar-go/loader"
	"github.com/bdpiprava/scalar-go/model"
)

func Test_Load(t *testing.T) {
	spec, err := loader.LoadFromDir("../data/loader", "pet-store.yml")

	require.NoError(t, err)
	require.NotNil(t, spec)
	requireBase(t, spec)
	requireSchema(t, spec)
	requirePaths(t, spec)
}

func Test_Load_JsonFile(t *testing.T) {
	spec, err := loader.LoadFromDir("../data/loader", "pet-store.json")

	require.NoError(t, err)
	require.NotNil(t, spec)
}

func Test_Load_InvalidFileType(t *testing.T) {
	spec, err := loader.LoadFromDir("../data/loader", "pet-store.xyz")

	require.ErrorContains(t, err, `file '../data/loader/pet-store.xyz' is not a YAML or JSON file, supported extensions are [yml|yaml|json]`)
	require.Nil(t, spec)
}

func Test_Load_MultipleFiles(t *testing.T) {
	spec, err := loader.LoadFromDir("../data/loader-multiple-files", "api.yml")

	require.NoError(t, err)
	require.NotNil(t, spec)
	requireBase(t, spec)
	requireSchema(t, spec)
	requirePaths(t, spec)
}

func Test_LoadedFileShouldHaveIdenticalContent(t *testing.T) {
	specFromMultipleFiles, err := loader.LoadFromDir("../data/loader-multiple-files", "api.yml")
	require.NoError(t, err)

	specFromSingleFile, err := loader.LoadFromDir("../data/loader", "pet-store.yml")
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(specFromMultipleFiles, specFromSingleFile))
}

func Test_Load_DocumentedPath(t *testing.T) {
	spec, err := loader.LoadFromDir("../data/loader", "pet-store.yml")
	require.NoError(t, err)

	paths := spec.DocumentedPaths()

	sort.Slice(paths, func(i, j int) bool {
		return paths[i].String() <= paths[j].String()
	})

	require.Len(t, paths, 3)
	require.Equal(t, model.DocumentedPath{Path: "/pets", Method: "get"}, paths[0])
	require.Equal(t, model.DocumentedPath{Path: "/pets/{petId}", Method: "get"}, paths[1])
	require.Equal(t, model.DocumentedPath{Path: "/pets", Method: "post"}, paths[2])
}

func Test_Load_XTagGroups(t *testing.T) {
	spec, err := loader.LoadFromDir("../data/xTagGroups", "withXTagGroups.yaml")
	require.NoError(t, err)

	require.Len(t, spec.TagsGroup, 2)
	require.Equal(t, model.TagGroup{
		Name:        "GroupOne",
		Description: "These are the GroupOne APIs",
		Tags:        []string{"SubGroup1.1", "SubGroup1.2"},
	}, spec.TagsGroup[0])

	require.Equal(t, model.TagGroup{
		Name:        "GroupTwo",
		Description: "These are the GroupTwo APIs",
		Tags:        []string{"SubGroup2.1"},
	}, spec.TagsGroup[1])
}

func Test_LoadFromBytes(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
	}{
		{
			name:     "YAML",
			filePath: "../data/loader/pet-store.yml",
		},
		{
			name:     "JSON",
			filePath: "../data/loader/pet-store.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := os.ReadFile(tc.filePath)
			require.NoError(t, err)

			spec, err := loader.LoadFromBytes(content)
			require.NoError(t, err)
			require.NotNil(t, spec)

			// test JSON file does not have expected schema
			if tc.name != "JSON" {
				requireBase(t, spec)
				requireSchema(t, spec)
				requirePaths(t, spec)
			}
		})
	}
}

func getPathParam(params []interface{}, name string) model.GenericObject {
	for _, param := range params {
		if param.(model.GenericObject)["name"] == name {
			return param.(model.GenericObject)
		}
	}
	return nil
}

func requireBase(t *testing.T, spec *model.Spec) {
	require.Equal(t, "3.0.0", spec.OpenAPI)

	require.Equal(t, "1.0.0", spec.Info.Version)
	require.Equal(t, "Swagger Petstore", spec.Info.Title)
	require.Equal(t, "MIT", spec.Info.License.Name)

	require.Len(t, spec.Servers, 1)
	require.Equal(t, model.Server{URL: "http://petstore.swagger.io/v1"}, spec.Servers[0])

	require.Len(t, spec.Tags, 2)
	require.Equal(t, model.Tag{Name: "pets", Description: "Everything about your Pets"}, spec.Tags[0])
	require.Equal(t, model.Tag{Name: "store", Description: "Access to Petstore"}, spec.Tags[1])
}

func requirePaths(t *testing.T, spec *model.Spec) {
	require.Len(t, spec.Paths, 2)
	pathGetPets := spec.Paths["/pets"].(model.GenericObject)["get"].(model.GenericObject)
	require.Equal(t, "List all pets", pathGetPets["summary"])
	require.Equal(t, "listPets", pathGetPets["operationId"])
	require.Equal(t, []interface{}{"pets"}, pathGetPets["tags"])

	nameParam := getPathParam(pathGetPets["parameters"].([]interface{}), "limit")
	require.Equal(t, model.GenericObject{
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

	require.Equal(t, model.GenericObject{
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
			"$ref": "#/components/responses/Error",
		},
	}, pathGetPets["responses"])
}

func requireSchema(t *testing.T, spec *model.Spec) {
	require.Len(t, spec.Components.Schemas, 3)
	petSchema := spec.Components.Schemas["Pet"].(model.GenericObject)
	require.NotNil(t, petSchema)
	require.Equal(t, "object", petSchema["type"])
	require.Equal(t, []interface{}{"id", "name"}, petSchema["required"])
	require.Equal(t, model.GenericObject{"format": "int64", "type": "integer"}, petSchema["properties"].(model.GenericObject)["id"])
	require.Equal(t, model.GenericObject{"type": "string"}, petSchema["properties"].(model.GenericObject)["name"])
	require.Equal(t, model.GenericObject{"type": "string"}, petSchema["properties"].(model.GenericObject)["tag"])

	petsSchema := spec.Components.Schemas["Pets"].(model.GenericObject)
	require.NotNil(t, petsSchema)
	require.Equal(t, "array", petsSchema["type"])
	require.Equal(t, 100, petsSchema["maxItems"])
	require.Equal(t, "#/components/schemas/Pet", petsSchema["items"].(model.GenericObject)["$ref"])

	errorSchema := spec.Components.Schemas["Error"].(model.GenericObject)
	require.NotNil(t, errorSchema)
	require.Equal(t, "object", errorSchema["type"])
	require.Equal(t, []interface{}{"code", "message"}, errorSchema["required"])
	require.Equal(t, model.GenericObject{"type": "integer", "format": "int32"}, errorSchema["properties"].(model.GenericObject)["code"])
	require.Equal(t, model.GenericObject{"type": "string"}, errorSchema["properties"].(model.GenericObject)["message"])
}
