package cmd

import "github.com/sirupsen/logrus"

// noopCleaner returns the value passed in.
func noopCleaner(v interface{}) interface{} {
	return v
}

// stringKeyedMapCleaner iterates over the provided value, and replaces all
// map[interface{}]interface{} with map[string]interface{} if possible. If
// there is a map[interface{}]interface{} with a non-string key, any maps
// within that map will not be processed.
func stringKeyedMapCleaner(v interface{}) interface{} {
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
			stringMap[key.(string)] = stringKeyedMapCleaner(value)
		}

		logrus.Debug("Successfully converted interface keyed map to string keyed map")
		return stringMap
	case []interface{}:
		sl := make([]interface{}, len(vv))
		for i, value := range vv {
			sl[i] = stringKeyedMapCleaner(value)
		}

		return sl
	default:
		return v
	}
}
