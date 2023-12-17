/**
 * @Author: steven
 * @Description:
 * @File: attribute
 * @Date: 17/12/23 23.45
 */

package telemetry

import (
	"encoding/json"
	"fmt"
	"github.com/evorts/kevlars/utils"
	"go.opentelemetry.io/otel/attribute"
)

type BytesMapAttributes interface {
	Exec() []attribute.KeyValue
}

// permanentlyMaskedFields fields that by default will be masked
var permanentlyMaskedFields = map[string]int8{
	"password": 1,
	"pass":     1,
	"pwd":      1,
	"ktp":      1,
	"nik":      1,
	"npwp":     1,
}

type bytesMapAttributes struct {
	bytesMap    []byte
	result      map[string]string
	excludeKeys map[string]int8
	attrKVs     []attribute.KeyValue
	prefix      string
}

func (a *bytesMapAttributes) Exec() []attribute.KeyValue {
	if a.excludeKeys == nil {
		a.excludeKeys = make(map[string]int8)
	}
	a.excludeKeys = utils.MapMerge(a.excludeKeys, permanentlyMaskedFields)
	// convert bytes into map
	var m map[string]interface{}
	_ = json.Unmarshal(a.bytesMap, &m)
	if m == nil {
		return a.attrKVs
	}
	for k, v := range m {
		key := k
		if len(a.prefix) > 0 {
			key = fmt.Sprintf("%s.%s", a.prefix, k)
		}
		if _, ok := a.excludeKeys[k]; ok {
			v = "********"
		}
		a.attrKVs = append(a.attrKVs, attribute.String(key, fmt.Sprintf("%v", v)))
	}
	return a.attrKVs
}

func UseBytesMapAttributes(prefix string, bytesMap []byte, excludeKeys map[string]int8) BytesMapAttributes {
	return &bytesMapAttributes{
		prefix:      prefix,
		bytesMap:    bytesMap,
		result:      make(map[string]string),
		excludeKeys: excludeKeys,
		attrKVs:     make([]attribute.KeyValue, 0),
	}
}
