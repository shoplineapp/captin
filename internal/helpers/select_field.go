package helpers

import (
  "fmt"
  "sort"
  "strings"
  "reflect"
  "github.com/thoas/go-funk"
  "github.com/mohae/deepcopy"
)

func IncludeFields(object map[string]interface{}, fields []string) interface{} {
  fieldTree := buildFieldTree(fields)
  clone := deepcopy.Copy(object)
  return filter(clone, fieldTree, "include")
}

func ExcludeFields(object map[string]interface{}, fields []string) interface{} {
  fieldTree := buildFieldTree(fields)
  clone := deepcopy.Copy(object)
  return filter(clone, fieldTree, "exclude")
}

func buildNode(ancestor map[string]interface{}, nodes []string, i int) {
  if(i >= len(nodes)) { return }

  node := nodes[i]

  if(ancestor[node] == nil) { ancestor[node] = map[string]interface{}{} }

  buildNode(ancestor[node].(map[string]interface{}), nodes, i + 1)
}

// adapted from open-api-node FieldFilter
func buildFieldTree(fields []string) map[string]interface{} {
  sort.Strings(fields)

  tree := map[string]interface{}{}

  for _, field := range fields {
    nodes := strings.Split(field, ".")
    buildNode(tree, nodes, 0)
  }

  return tree
}

func applyFilter(result map[string]interface{}, key string, value interface{}, mode string) {
  switch (mode) {
    case "include":
      result[key] = value
    case "exclude":
      delete(result, key)
    default:
      panic(fmt.Sprintf("unknown filter mode %s", "mode"))
  }
}

func filter(object interface{}, fieldTree map[string]interface{}, mode string) interface{} {
  var result interface{}

  switch mode {
    case "include":
      result = map[string]interface{}{}
    case "exclude":
      result = object
  }

  // do not filter primitive types (e.g, object contains array of strings)
  if(reflect.TypeOf(object).Kind() != reflect.Map) {
    return object
  }

  // eslint-disable-next-line no-restricted-syntax
  for key, value := range object.(map[string]interface{}) {
    // not in include/exclude tree, keep it as is
    if (!funk.Contains(funk.Keys(fieldTree), key)) {
      continue
    }

    // {} is wildcard, keep/discard everything
    if (reflect.TypeOf(fieldTree[key]).Kind() == reflect.Map && len(fieldTree[key].(map[string]interface{})) == 0){
      applyFilter(result.(map[string]interface{}), key, value, mode)
      continue
    }

    // nested object
    if (reflect.TypeOf(value).Kind() == reflect.Map) {
      result.(map[string]interface{})[key] = filter(value, fieldTree[key].(map[string]interface{}), mode)
      continue
    }

    // apply filter on each element of array
    if (reflect.TypeOf(value).Kind() == reflect.Slice || reflect.TypeOf(value).Kind() == reflect.Array) {
      result.(map[string]interface{})[key] = funk.Map(value, func(element interface{}) interface{} {
        return filter(element, fieldTree[key].(map[string]interface{}), mode)
      })
    }
  }

  return result
}
