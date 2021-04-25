package gojsoncompare

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
)

type ArrayLessSort func(a, b interface{}, parentKey string) bool

func isSamesies(a, b interface{}, asort ArrayLessSort, pkey string) bool {
	if a == nil || b == nil {
		if a == nil && b == nil {
			return true
		}
		return false
	}

	ka := reflect.TypeOf(a).Kind()
	kb := reflect.TypeOf(b).Kind()
	if ka != kb {
		return false
	}

	switch ka {
	case reflect.Bool:
		if a.(bool) != b.(bool) {
			return false
		}
	case reflect.String:
		switch aa := a.(type) {
		case json.Number:
			bb, ok := b.(json.Number)
			if !ok || aa != bb {
				return false
			}
		case string:
			bb, ok := b.(string)
			if !ok || aa != bb {
				return false
			}
		}
	case reflect.Slice:
		sa, sb := a.([]interface{}), b.([]interface{})
		salen, sblen := len(sa), len(sb)
		if salen != sblen {
			return false
		}

		if asort != nil {
			sort.Slice(sa, func(i, j int) bool {
				return asort(sa[i], sa[j], pkey)
			})
			sort.Slice(sb, func(i, j int) bool {
				return asort(sb[i], sb[j], pkey)
			})
		}

		for i := 0; i < salen; i++ {
			same := isSamesies(sa[i], sb[i], asort, "")
			if !same {
				return false
			}
		}
		return true
	case reflect.Map:
		ma, mb := a.(map[string]interface{}), b.(map[string]interface{})
		keysMap := make(map[string]bool)
		for k := range ma {
			keysMap[k] = true
		}
		for k := range mb {
			keysMap[k] = true
		}
		keys := make([]string, 0, len(keysMap))
		for k := range keysMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			va, aok := ma[k]
			vb, bok := mb[k]

			if aok && bok {
				same := isSamesies(va, vb, asort, k)
				if !same {
					return false
				}
			}
		}
		return true
	}

	return true
}

// Does a deep equal comparison on two json documents. You can optionally pass in a
// less function to sort arrays.
//
// Returned true if the json documents are equal.
func DeepEqual(a, b []byte, asort ArrayLessSort) bool {
	var av, bv interface{}
	da := json.NewDecoder(bytes.NewReader(a))
	da.UseNumber()
	db := json.NewDecoder(bytes.NewReader(b))
	db.UseNumber()
	errA := da.Decode(&av)
	errB := db.Decode(&bv)
	if errA != nil || errB != nil {
		return false
	}

	return isSamesies(av, bv, asort, "")
}
