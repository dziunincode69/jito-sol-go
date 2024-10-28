package utils

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"

	json "github.com/json-iterator/go"
)

func ParseSimulateLogToJson(log, keyword string) (string, error) {
	reg := regexp.MustCompile(`{["\w:,]+}`)
	results := reg.FindAllString(log, -1)

	if results == nil || len(results) != 1 {
		return "", fmt.Errorf("simulate log failed to match json for keyword: %s", keyword)
	}

	return results[0], nil
}

func ParseSimulateValue(log, key string) (*big.Int, error) {
	reg := regexp.MustCompile(`"` + key + `":(\d+)`)
	results := reg.FindStringSubmatch(log)

	if results == nil || len(results) != 2 {
		return nil, errors.New("simulate log failed to match key: " + key)
	}
	value := StringToBig256(results[1])
	return value, nil
}

func ParseSimulateValueAsInt(jsonStr, key string) (int, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return 0, err
	}

	value, ok := data[key]
	if !ok {
		return 0, errors.New("key not found in JSON")
	}

	switch v := value.(type) {
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, errors.New("value is not a number")
	}
}
