package kvpearl

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"time"

	"neckless.adviser.com/key"
	"neckless.adviser.com/pearl"
)

// // fmt.Println("XXXX", len(key.Values), len(mergedValue), mergedValue)
// if len(key.Values) < len(mergedValue) {
// 	for i := len(key.Values); i < len(mergedValue); i++ {
// 		key.Values = append(key.Values, &Value{
// 			Value: "",
// 			Tags:  []string{},
// 		})
// 		// fmt.Println("APPEND-To-Values", len(key.Values))
// 	}
// }
// // key.Values = make([]Value, len(mergedValue))
// idx := 0
// for i := range mergedValue {
// 	key.Values[idx] = mergedValue[i]
// 	idx++
// }

type JsonKVPearl struct {
	// Tags []string
	Keys        []*JsonKey // sorted by keys
	FingerPrint string     `json:"FingerPrint,omitempty"`
	Created     time.Time
}

type KVPearl struct {
	// Tags []string
	Keys    keys
	Created time.Time
	// Seq     int64            `json:"-"` // windows timer is not good enough
	Pearl *pearl.OpenPearl `json:"-"`
}

func Create( /* tags ...string */ ) *KVPearl {
	return &KVPearl{
		// Tags: uniqStrings(tags),
		Keys:    keys{},
		Created: time.Now(),
	}
}

type SetArg struct {
	Key        string  // is set if plain Key
	Unresolved *string // Unresolved is set value was resolved
	Val        string  // Value is Set if an = is used
	Tags       []string
}

// type Map map[string]*kvpearl.SetArg,
// type Map map[string]*kvpearl.,
// type Map map[string]*kvpearl.SetArg,toKVpearls

func (kvp *KVPearl) Set(a SetArg) *KVPearl {
	key, _ := kvp.Keys.getOrAdd(a.Key)
	key.setValue(a.Unresolved, a.Val, tags2Map(a.Tags))
	return kvp
}

// func (sa *SetArg) ToKVParsed() *KVParsed {
// 	return &KVParsed{
// 		Key:        sa.Key,
// 		Unresolved: sa.Unresolved,
// 		Val:        sa.Val,
// 		Tags:       stringArray2Map(sa.Tags),
// 	}
// }

func FromJSON(jsStr []byte) (*KVPearl, error) {
	jskvp := JsonKVPearl{}
	err := json.Unmarshal(jsStr, &jskvp)
	if err != nil {
		return nil, err
	}
	kvp := Create()
	kvp.Created = jskvp.Created
	// kvp.Tags = uniqStrings(kvp.Tags)
	for i := range jskvp.Keys {
		key := jskvp.Keys[i]
		for j := range key.Values {
			kvp.Set(SetArg{
				Key:        key.Key,
				Unresolved: key.Values[j].Unresolved,
				Val:        key.Values[j].Value,
				Tags:       key.Values[j].Tags,
			})

		}
		// kvp.Keys[key.Key] = &Key{
		// 	Key:    key.Key,
		// 	Values: kvp.get(key.Key),
		// }
	}
	return kvp, nil
}

func (kvp *KVPearl) AsJSON() *JsonKVPearl {
	keys := kvp.Keys.Sorted().asJson()
	fpr := ""
	if kvp.Pearl != nil {
		fpr = base64.StdEncoding.EncodeToString(kvp.Pearl.Closed.FingerPrint)
	}
	return &JsonKVPearl{
		Created:     kvp.Created,
		FingerPrint: fpr,
		Keys:        keys,
	}
}

const Type = "KVPearl"

func (kvp *KVPearl) ClosePearl(owners *pearl.PearlOwner) (*pearl.Pearl, error) {
	jsonStr, err := json.Marshal(kvp.AsJSON())
	if err != nil {
		return nil, err
	}
	return pearl.Close(&pearl.CloseRequestPearl{
		Type:    Type,
		Payload: jsonStr,
		Owners:  *owners,
	})
}

func OpenPearl(pks []*key.PrivateKey, prl *pearl.Pearl) (*KVPearl, error) {
	op, err := pearl.Open(pks, prl)
	if err != nil {
		return nil, err
	}
	kvp, err := FromJSON(op.Payload)
	if err != nil {
		return nil, err
	}
	kvp.Pearl = op
	return kvp, err
}

func findTag(t1 Tags, t2 []string) bool {
	if len(t2) == 0 {
		return true
	}
	if len(t1) == 0 {
		return false
	}
	var ref []string
	var toSet []string
	if len(t1) > len(t2) {
		ref = t2
		toSet = t1.sorted()
	} else {
		ref = t1.sorted()
		toSet = t2
	}
	set := map[string](struct{}){}
	for i := range toSet {
		set[toSet[i]] = struct{}{}
	}
	for i := range ref {
		_, found := set[ref[i]]
		if found {
			// fmt.Println("FOUND:", t1, t2, set)
			return true
		}
	}
	// fmt.Println("!FOUND:", t1, t2, set)
	return false
}

func (kvp *KVPearl) matchTag(tags []string) []KeyValue {
	// if len(kvp.Tags) != 0 {
	// 	if !findTag(kvp.Tags, tags) {
	// 		return []KeyValue{}
	// 	}
	// }
	mapRet := map[string]KeyValue{}
	for i := range kvp.Keys {
		key := kvp.Keys[i]
		vals := key.Values.revOrdered()
		for j := range *vals {
			if findTag((*vals)[j].Tags, tags) {
				mapRet[key.Key] = KeyValue{
					Key:   key.Key,
					Value: (*vals)[j].Value,
				}
				break
			}
		}

	}
	ret := make(KeyValues, len(mapRet))
	idx := 0
	for i := range mapRet {
		ret[idx] = mapRet[i]
		idx++
	}
	sort.Sort(&ret)
	return ret
}
