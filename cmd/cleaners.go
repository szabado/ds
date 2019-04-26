package cmd

import "github.com/sirupsen/logrus"

func noopCleaner(v interface{}) interface{} {
	return v
}

func stringMapKeyCleaner(v interface{}) interface{} {
	switch vv := v.(type) {
	case map[interface{}]interface{}:
		logrus.Debug("Checking if interface keyed map can be converted")
		// Check for any non-string keys
		for key := range vv {
			if _, ok := key.(string); !ok {
				logrus.Debug("Interface keyed map cannot be converted")
				return v
			}
		}

		stringMap := make(map[string]interface{})
		for key, value := range vv {
			stringMap[key.(string)] = stringMapKeyCleaner(value)
		}

		logrus.Debug("Successfully converted interface keyed map to string keyed map")
		return stringMap
	case []interface{}:
		sl := make([]interface{}, len(vv))
		for i, value := range vv {
			sl[i] = stringMapKeyCleaner(value)
		}

		return sl
	default:
		return v
	}
}
