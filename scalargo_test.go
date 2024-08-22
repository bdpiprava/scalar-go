package scalargo

import (
	"testing"

	"github.com/bdpiprava/scalar-go/model"
	"github.com/stretchr/testify/assert"
)

func Test_ShouldCallOverrideHandler_WhenProvider(t *testing.T) {
	var called bool
	content, err := New(
		"./data/loader",
		WithBaseFileName("pet-store.yml"),
		WithSpecModifier(func(spec *model.Spec) *model.Spec {
			called = true
			assert.Equal(t, "Swagger Petstore", spec.Info.Title)

			spec.Info.Title = "PetStore API"
			return spec
		}),
	)

	assert.NoError(t, err)
	assert.True(t, called)
	assert.NotNil(t, content)
	assert.Contains(t, content, "<title>PetStore API</title>")
}
