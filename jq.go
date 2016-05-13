package jq

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func JsonQuery(data []byte, query string) ([]byte, error) {
	if len(query) > 0 && query[0] == '.' { //最外层调用
		if len(query) == 1 {
			return data, nil
		}
		query = strings.Replace(query, "[", ".[", -1)
	}
	elems := strings.Split(strings.Trim(query, "."), ".")
	var reply interface{}
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, errors.New("query is null")
	}
	switch {
	case reflect.TypeOf(reply).Kind() == reflect.Map && !strings.Contains(elems[0], "["):
		replyMap := reply.(map[string]interface{})
		if child, ok := replyMap[elems[0]]; ok {
			if len(elems) > 1 {
				childjson, _ := json.Marshal(child)
				return JsonQuery(childjson, strings.Join(elems[1:], "."))
			} else {
				return json.Marshal(child)
			}
		} else {
			return nil, errors.New("No such element:" + elems[0])
		}
	case reflect.TypeOf(reply).Kind() == reflect.Slice && strings.Contains(elems[0], "["):
		if key, err := strconv.Atoi(strings.Trim(elems[0], "[]")); err != nil {
			return nil, err
		} else {
			replySlice := reply.([]interface{})
			if key > len(replySlice)-1 {
				return nil, errors.New("slice out of range")
			}
			if len(elems) > 1 {
				childjson, _ := json.Marshal(replySlice[key])
				return JsonQuery(childjson, strings.Join(elems[1:], "."))
			} else {
				return json.Marshal(replySlice[key])
			}
		}
	case reflect.TypeOf(reply).Kind() == reflect.String || reflect.TypeOf(reply).Kind() == reflect.Float64:
		if len(elems) > 0 {
			return nil, errors.New("string or number type has no child field!")
		} else {
			return data, nil
		}
	default:
		return nil, errors.New("Invalid json type:" + reflect.TypeOf(reply).Kind().String() + " with query: " + query)
	}
}

func String(data []byte, err error) (string, error) {
	if err != nil {
		return "", err
	}
	var tmp interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return "", err
	}
	switch tmp.(type) {
	case string:
		return strings.Trim(tmp.(string), "\""), nil
	case float64:
		return fmt.Sprintf("%f", tmp.(float64)), nil
	case map[string]interface{}, []interface{}:
		return string(data), nil
	case nil:
		return "<nil>", nil
	default:
		return "", errors.New("unknown json format:" + string(data))
	}
}

func Int64(data []byte, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	var ret int64
	if err := json.Unmarshal(data, &ret); err != nil {
		return 0, err
	}
	return ret, nil
}
