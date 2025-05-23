package output

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

func PrintJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func PrintYAML(v interface{}) {
	data, _ := yaml.Marshal(v)
	fmt.Println(string(data))
} 