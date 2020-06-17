package miuires

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// FilterConfig holds filter rules
type FilterConfig struct {
	StringsKeyRules   map[string][]FilterRules `yaml:"strings_key_rules"`
	ArraysKeyRules    map[string][]FilterRules `yaml:"arrays_key_rules"`
	PluralsKeyRules   map[string][]FilterRules `yaml:"plurals_key_rules"`
	StringsValueRules map[string][]FilterRules `yaml:"strings_value_rules"`
	ArraysValueRules  map[string][]FilterRules `yaml:"arrays_value_rules"`
}

// FilterRules holds rules used to filter keys and/or values
type FilterRules struct {
	Match string `yaml:"match"`
	Mode  string `yaml:"mode"`
}

// GetFilterConfigFromFile returns a new FilterConfig from a YAML file
func GetFilterConfigFromFile(r *os.File) (*FilterConfig, error) {

	// Read data from reader
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Initialise FilterConfig
	var fc FilterConfig

	// Unmarshall data in HCL to config
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return nil, err
	}

	return &fc, nil
}
