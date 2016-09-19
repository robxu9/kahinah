package conf

import (
	"github.com/pelletier/go-toml"
	"github.com/robxu9/kahinah/common/klog"
)

var (
	// Tree represents the TOML configuration tree for the application
	Tree *toml.TomlTree
)

func init() {
	var err error
	Tree, err = toml.LoadFile("app.toml")
	if err != nil {
		klog.Fatalf("conf: unable to load app.toml: %v", err)
	}
}

// Get returns the value for the key specified.
func Get(key string) interface{} {
	return Tree.Get(key)
}

// GetDefault returns the value for the key specified, or "def" if not found.
func GetDefault(key string, def interface{}) interface{} {
	return Tree.GetDefault(key, def)
}

// Has checks whether the key specified was found in the tree.
func Has(key string) bool {
	return Tree.Has(key)
}

// Keys return the list of keys in the TOML tree.
func Keys() []string {
	return Tree.Keys()
}

// Query queries the TOML tree and returns its result.
func Query(query string) (*toml.QueryResult, error) {
	return Tree.Query(query)
}

// Set sets a value in the TOML tree.
func Set(key string, value interface{}) {
	Tree.Set(key, value)
}

// String returns the TOML tree as a string.
func String() string {
	return Tree.String()
}

// ToMap returns the TOML tree as a map.
func ToMap() map[string]interface{} {
	return Tree.ToMap()
}
