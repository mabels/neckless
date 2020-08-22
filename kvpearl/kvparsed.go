package kvpearl

import (
	"errors"
	"fmt"
	"regexp"
)

type KVParsed struct {
	Key       *string        // is set if plain Key
	KeyRegex  *regexp.Regexp // is allways set
	ToResolve *string        // ToResolv is set if Needs to Resolv
	Val       *string        // Value is Set if an = is used
	Tags      Tags
}

type ResolvFN func(key string, fname string) (*string, error)

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

func isToResolve(isAt string, value string) *string {
	if isAt == "@" {
		return &value
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
		return nil, errors.New(fmt.Sprintf("no matching kv:[%s]", arg))
	}
	keyRegex, err := asMatch(split[1])
	if err != nil {
		return nil, err
	}
	if len(split) > 3 && len(split[5]) > 0 {
		// fmt.Println("split set tag")
		return &KVParsed{
			Key:      isPlain(split[1]),
			KeyRegex: keyRegex,

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

func (p *KVParsed) Resolv(fn ResolvFN) (*KVParsed, error) {
	res, err := fn(*p.Key, *p.ToResolve)
	if err != nil {
		return nil, err
	}
	p.Val = res
	return p, nil
}

func (p *KVParsed) ToSetArgs() (*SetArg, error) {
	if p.Key == nil {
		return nil, errors.New("ToSetArgs need a key")
	}
	if p.Val == nil {
		return nil, errors.New("ToSetArgs need a val")
	}
	return &SetArg{
		Key:        *p.Key,      // is set if plain Key
		Unresolved: p.ToResolve, // Unresolved is set value was resolved
		Val:        *p.Val,      // Value is Set if an = is used
		Tags:       p.Tags.toArray(),
	}, nil
}

func Parse(arg string) (*KVParsed, error) {
	parsed, err := parseBrackets(arg)
	if err != nil {
		parsed, err = parseComma(arg)
	}
	return parsed, err
}
