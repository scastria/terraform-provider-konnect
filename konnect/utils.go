package konnect

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
)

func convertSetToArray(set *schema.Set) []string {
	setList := set.List()
	retVal := []string{}
	for _, s := range setList {
		retVal = append(retVal, s.(string))
	}
	return retVal
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func removeNulls(m map[string]interface{}) {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(m, e.String())
			continue
		}
		switch t := v.Interface().(type) {
		// If key is a JSON object (Go Map), use recursion to go deeper
		case map[string]interface{}:
			removeNulls(t)
		}
	}
}

func copyMap[K string, V interface{}](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}
