package logging

// LogKeyValuePair is a helper struct to help organise the use of tflog. tflog
// uses pairs to print key-value items/args, so this just to help visualize the
// way the logging works
type LogKeyValuePair struct {
	// Key is what will appear next to the value, use it as a small hint/description
	Key interface{}
	// The value is what you want to log.
	Value interface{}
}

func MakePair(key, value interface{}) LogKeyValuePair {
	return LogKeyValuePair{
		Key:   key,
		Value: value,
	}
}

func UnpackPairs(pairs ...LogKeyValuePair) (args []interface{}) {
	args = make([]interface{}, 2*len(pairs))
	for _, pair := range pairs {
		args = append(args, pair.Key, pair.Value)
	}
	return args
}
