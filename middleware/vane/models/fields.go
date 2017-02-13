package models

// Values maps a string key to a list of values.
// It is typically used for query parameters and form values.
type Values map[string][]interface{}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v Values) Get(key string) interface{} {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (v Values) Set(key string, value interface{}) {
	v[key] = []interface{}{value}
}

func (v Values) SetAll(key string, values []interface{}) {
	v[key] = values
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v Values) Add(key string, value interface{}) {
	v[key] = append(v[key], value)
}

// Del deletes the values associated with key.
func (v Values) Del(key string) {
	delete(v, key)
}
