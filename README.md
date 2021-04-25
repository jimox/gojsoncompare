# GoJSONCompare

Allows you to perform a deep equal on two json documents. You can optionally pass in a less function that allows you to sort arrays. 

## Usage
```
a := `{ "AnArray": [{ "Node": 1 }, { "Node": 2 }]}`
b := `{ "AnArray": [{ "Node": 2 }, { "Node": 1 }]}`

sorter := func(a, b interface{}, parentKey string) bool {
    if parentKey == "AnArray" {
        ai, aok := a.(map[string]interface{})
        bi, bok := b.(map[string]interface{})
        if aok && bok {
            av, aok := ai["Node"]
            bv, bok := bi["Node"]
            if aok && bok {
                return av.(json.Number) < bv.(json.Number)
            }
        }
    }
    return true
}

isEqual := gojsoncompare.DeepEqual([]byte(a), []byte(b), sorter)
fmt.Printf("Is Equal: %v\n", isEqual)
```