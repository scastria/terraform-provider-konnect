package konnect

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
)

func convertSetToArray(set *schema.Set) []string {
	setList := set.List()
	retVal := []string{}
	for _, s := range setList {
		line := ""
		if s != nil {
			line = s.(string)
		}
		retVal = append(retVal, line)
	}
	return retVal
}

func unionArrays(l1 []string, l2 []string) []string {
	uniqueMap := make(map[string]bool)
	retVal := []string{}
	for _, item := range l1 {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			retVal = append(retVal, item)
		}
	}
	for _, item := range l2 {
		if !uniqueMap[item] {
			uniqueMap[item] = true
			retVal = append(retVal, item)
		}
	}
	return retVal
}

func subtractArrays(l1 []string, l2 []string) []string {
	removeMap := make(map[string]bool)
	for _, item := range l2 {
		removeMap[item] = true
	}
	retVal := []string{}
	for _, item := range l1 {
		if !removeMap[item] {
			retVal = append(retVal, item)
		}
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

func copyMapByJSON[K string, V interface{}](m map[K]V) map[K]V {
	bytes, _ := json.Marshal(m)
	result := make(map[K]V)
	json.Unmarshal(bytes, &result)
	return result
}
