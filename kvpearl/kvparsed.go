package kvpearl

import (
	"errors"
	"fmt"
	"regexp"
)

// KVParsed is the Parsed KeyValue from the CommandLine
type KVParsed struct {
	Key       *string        // is set if plain Key
	KeyRegex  *regexp.Regexp // is allways set
	ToResolve *FuncsAndParam // ToResolv is set if Needs to Resolv
	// Actions   *[]string      // set if a Action is given
	Val  *string // Value is Set if an = is used
	Tags Tags
}

// ResolvFN is the Resolver Function Reference
type ResolvFN func(key string, fparam FuncsAndParam) (*string, error)

var plainRegex = regexp.MustCompile("^[A-Za-z0-9_]+$")

func isPlain(val string) *string {
	if plainRegex.MatchString(val) {
		return &val
	}
	return nil
}

var wildCardRegex = regexp.MustCompile(".*")

func asMatch(val string) (*regexp.Regexp, error) {
	if len(val) == 0 {
		return wildCardRegex, nil
	} else if isPlain(val) != nil {
		return regexp.Compile(fmt.Sprintf("^%s$", val))
	} else {
		return regexp.Compile(val)
	}
}

var isKeyValue = regexp.MustCompile(`^([^=@]+)([=@])(.*)$`)

// func processResolv(toResolv string, key string, value string) (string, *string, error) {
// 	if toResolv == "@" {
// 		if len(value) == 0 {
// 			return value, nil, nil
// 		}
// 		val, err := resolv(key, value)
// 		if err != nil {
// 			return value, nil, err
// 		}
// 		if val != nil {
// 			return *val, &value, nil
// 		}
// 	}
// 	return value, nil, nil
// }

var funcMatch = regexp.MustCompile("([^\\(]+)\\((.*)\\)")

// ParseFuncsAndParams parses a String to FuncsAndParam
func ParseFuncsAndParams(toParse string) *FuncsAndParam {
	ret := FuncsAndParam{
		Param: "",
		Funcs: []string{},
	}
	return interalParseFuncsAndParams(toParse, &ret)
}
func interalParseFuncsAndParams(toParse string, ret *FuncsAndParam) *FuncsAndParam {
	m := funcMatch.FindStringSubmatch(toParse)
	if len(m) != 3 {
		ret.Param = toParse
		return ret
	}
	ret.Funcs = append(ret.Funcs, m[1])
	return interalParseFuncsAndParams(m[2], ret)
}

func isToResolve(isAt string, value string) *FuncsAndParam {
	if isAt == "@" {
		return ParseFuncsAndParams(value)
	}
	return nil
}

func isValue(isAt string, value string) *string {
	if isToResolve(isAt, value) == nil {
		return &value
	}
	return nil
}

func parseComma(arg string) (*KVParsed, error) {
	split := isKeyValue.FindStringSubmatch(arg)
	if len(split) != 4 {
		return nil, errors.New("found no key value")
	}
	commas := tagstring2Map(split[3])
	value := ""
	tags := Tags{}
	if len(commas) == 1 {
		value = commas.byOrder()[0]
	}
	if len(commas) > 1 {
		value = commas.byOrder()[0]
		tags = tags2Map(commas.byOrder()[1:])
	}
	// value, unresolved, err := processResolv(split[2], split[1], value)
	// if unresolved != nil {
	// fmt.Println("parseComma:", value, *unresolved, err)
	// }
	// if err != nil {
	// return nil, err
	// }
	keyRegex, err := asMatch(split[1])
	if err != nil {
		return nil, err
	}
	parsed := KVParsed{
		Key:      isPlain(split[1]),
		KeyRegex: keyRegex,

		ToResolve: isToResolve(split[2], value),
		Val:       isValue(split[2], value),

		Tags: tags,
	}
	return &parsed, nil
}

var argRegex = regexp.MustCompile(`^([^=@]+)([=@])([^\[]*)(\[([^\]]*)\])$`)

func parseBrackets(arg string) (*KVParsed, error) {
	split := argRegex.FindStringSubmatch(arg)
	if len(split) != 6 {
		return nil, fmt.Errorf("no matching kv:[%s]", arg)
	}
	keyRegex, err := asMatch(split[1])
	if err != nil {
		return nil, err
	}
	if len(split) > 3 && len(split[5]) > 0 {
		// fmt.Println("split set tag")
		return &KVParsed{
			Key:       isPlain(split[1]),
			KeyRegex:  keyRegex,
			ToResolve: isToResolve(split[2], split[3]),
			Val:       isValue(split[2], split[3]),

			Tags: tagstring2Map(split[5]),
		}, nil
	}
	return &KVParsed{
		Key:      isPlain(split[1]),
		KeyRegex: keyRegex,

		ToResolve: isToResolve(split[2], split[3]),
		Val:       isValue(split[2], split[3]),

		Tags: Tags{},
	}, nil
}

// Resolv is the Shim t the ResolvFN
func (kvp *KVParsed) Resolv(fn ResolvFN) (*KVParsed, error) {
	if kvp.ToResolve == nil && kvp.Val != nil {
		return kvp, nil
	}
	res, err := fn(*kvp.Key, *kvp.ToResolve)
	if err != nil {
		return nil, err
	}
	kvp.Val = res
	return kvp, nil
}

// ToSetArgs converts KVParsed to ToSetArgs
func (kvp *KVParsed) ToSetArgs() (*SetArg, error) {
	if kvp.Key == nil {
		return nil, errors.New("ToSetArgs need a key")
	}
	if kvp.Val == nil {
		return nil, errors.New("ToSetArgs need a val")
	}
	return &SetArg{
		Key:        *kvp.Key,      // is set if plain Key
		Unresolved: kvp.ToResolve, // Unresolved is set value was resolved
		// Actions:    kvp.Actions,   // Actions is set value was processed
		Val:  *kvp.Val, // Value is Set if an = is used
		Tags: kvp.Tags.toArray(),
	}, nil
}

// Parse the given string to a KVParsed type
func Parse(arg string) (*KVParsed, error) {
	parsed, err := parseBrackets(arg)
	if err != nil {
		parsed, err = parseComma(arg)
	}
	return parsed, err
}

// Match a key and value agains a KVParsed
func (kvp *KVParsed) Match(key string, val *JSONValue) (*KVParsed, bool) {
	// findKey := false
	// for ikvps := range kvps {
	// kvp := kvps[ikvps]
	if kvp.KeyRegex.MatchString(key) {
		// fmt.Printf("Matched:%s:%d", kvp.Key, len(kvp.Tags))
		if len(kvp.Tags) == 0 {
			return kvp, true
		}
		for i := range val.Tags {
			// fmt.Printf("%s:%s\n", tag, kvp.Tags)
			_, found := kvp.Tags[val.Tags[i]]
			if found {
				return kvp, true
			}
		}
	}
	return nil, false
}
