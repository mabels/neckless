package kvpearl

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"time"

	"github.com/mabels/neckless/key"
	"github.com/mabels/neckless/pearl"
)

// JSONKVPearl the JSON resprestion of a KVPearl
type JSONKVPearl struct {
	// Tags []string
	Keys        []JSONKey // sorted by keys
	FingerPrint string    `json:"FingerPrint,omitempty"`
	Created     time.Time
}

// KVPearl which is stored in Gems
type KVPearl struct {
	Keys     keys
	Created  time.Time
	Pearl    *pearl.OpenPearl `json:"-"`
	orderRef *int
	order    int
}

// SetArg allow explict creating of a KVPearl
type SetArg struct {
	Key        string         // is set if plain Key
	Unresolved *FuncsAndParam // Unresolved is set value was resolved
	// Actions    *[]string      // Actions is set value was processed
	Val  string // Value is Set if an = is used
	Tags []string
}

// Set added the Args to a KVPearl
func (kvp *KVPearl) Set(a SetArg) *KVPearl {
	key, _ := kvp.Keys.getOrAdd(a.Key, kvp.orderRef)
	key.setValue(a.Unresolved, a.Val, tags2Map(a.Tags))
	return kvp
}

// FromJSON converts from json to a KVPearl
func FromJSON(jsStr []byte) (*KVPearl, error) {
	jskvp := JSONKVPearl{}
	err := json.Unmarshal(jsStr, &jskvp)
	if err != nil {
		return nil, err
	}
	kvp := CreateKVPearls().Add()
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

// AsJSON convert a KVPearl to a JSONKVPearl
func (kvp *KVPearl) AsJSON() *JSONKVPearl {
	keys := kvp.Keys.Sorted()
	fpr := ""
	if kvp.Pearl != nil {
		fpr = base64.StdEncoding.EncodeToString(kvp.Pearl.Closed.FingerPrint)
	}
	return &JSONKVPearl{
		Created:     kvp.Created,
		FingerPrint: fpr,
		Keys:        keys,
	}
}

// Type is the Typename of the KVPearl
const Type = "KVPearl"

// ClosePearl closes a Pearl with encryption
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

// OpenPearl tries to open the given pearl with the privatekey
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
		vals := key.Values.Ordered()
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
