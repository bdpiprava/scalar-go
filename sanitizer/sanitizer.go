package sanitizer

import (
	"fmt"

	"github.com/bdpiprava/scalar-go/model"
)

// Sanitize removes any non-string keys from the provided spec.
func Sanitize(spec *model.Spec) {
	spec.Components.Schemas = sanitizeGenericObject(spec.Components.Schemas)
	spec.Components.Parameters = sanitizeGenericObject(spec.Components.Parameters)
	spec.Paths = sanitizeGenericObject(spec.Paths)
}

func sanitizeInterfaceArray[R any](in []R) []R {
	res := make([]R, len(in))
	for i, v := range in {
		res[i] = sanitizeInterface(v).(R)
	}
	return res
}

func sanitizeInterfaceMap(in map[any]any) map[string]any {
	res := make(map[string]any)
	for k, v := range in {
		res[convertToString(k)] = sanitizeInterface(v)
	}
	return res
}

func convertToString(k any) string {
	return fmt.Sprintf("%v", k)
}

func sanitizeGenericObject(in model.GenericObject) model.GenericObject {
	res := make(model.GenericObject)
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = sanitizeInterface(v)
	}
	return res
}

func sanitizeInterface(val any) any {
	switch v := val.(type) {
	case []any:
		return sanitizeInterfaceArray[any](v)
	case map[any]any:
		return sanitizeInterfaceMap(v)
	case model.GenericObject:
		return sanitizeGenericObject(v)
	case string:
		return v
	default:
		return val
	}
}
