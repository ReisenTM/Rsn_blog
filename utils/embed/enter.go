package embed

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

// LoadMappingFromJSON 把json字符串转为es的types.TypeMapping
func LoadMappingFromJSON(jsonStr string) (*types.TypeMapping, error) {
	var mapping types.TypeMapping
	if err := json.Unmarshal([]byte(jsonStr), &mapping); err != nil {
		return nil, err
	}
	return &mapping, nil
}
