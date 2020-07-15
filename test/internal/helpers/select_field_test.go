package helpers_test

import (
  "testing"

  helpers "github.com/shoplineapp/captin/internal/helpers"
  "github.com/stretchr/testify/assert"
)

func TestIncludeFields(t *testing.T) {
  fields := []string{"foo"}
  object := map[string]interface{}{"foo": "bar", "foo2": "bar2"}
  result := helpers.IncludeFields(object, fields)
  assert.Equal(t, map[string]interface{}{ "foo": "bar" }, result)
}

func TestExcludeFields(t *testing.T) {
  fields := []string{"foo"}
  object := map[string]interface{}{"foo": "bar", "foo2": "bar2"}
  result := helpers.ExcludeFields(object, fields)
  assert.Equal(t, map[string]interface{}{ "foo2": "bar2" }, result)

  // org object not modified
  assert.Equal(t, map[string]interface{}{"foo": "bar", "foo2": "bar2"}, object)
}

func TestNestedIncludeFields(t *testing.T) {
  fields := []string{"foo.deepfoo", "foo2.common"}
  object := map[string]interface{}{"foo": map[string]interface{}{"deepfoo": "deepbar", "else": "useless"}, "foo2": []interface{}{map[string]interface{}{"arrayfoo1": "arraybar1", "common": "a"}, map[string]interface{}{"arrayfoo2": "arraybar2", "common": "b"}}}
  result := helpers.IncludeFields(object, fields)
  assert.Equal(t, map[string]interface {}{"foo":map[string]interface {}{"deepfoo":"deepbar"}, "foo2":[]interface {}{map[string]interface {}{"common":"a"}, map[string]interface {}{"common":"b"}}}, result)
}

func TestNestedExcludeFields(t *testing.T) {
  fields := []string{"foo.else", "foo2.common"}
  object := map[string]interface{}{"foo": map[string]interface{}{"deepfoo": "deepbar", "else": "useless"}, "foo2": []interface{}{map[string]interface{}{"arrayfoo1": "arraybar1", "common": "a"}, map[string]interface{}{"arrayfoo2": "arraybar2", "common": "b"}}}
  result := helpers.ExcludeFields(object, fields)
  assert.Equal(t, map[string]interface {}{"foo":map[string]interface {}{"deepfoo":"deepbar"}, "foo2":[]interface {}{map[string]interface {}{"arrayfoo1":"arraybar1"}, map[string]interface {}{"arrayfoo2":"arraybar2"}}}, result)
}
