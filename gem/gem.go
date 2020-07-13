package gem

import (
	"encoding/json"
	"sort"
	"strings"

	"neckless.adviser.com/key"
	"neckless.adviser.com/member"
	"neckless.adviser.com/pearl"
)

type Gem struct {
	PubKeys map[string]*member.PublicMember
	Pearl   *pearl.OpenPearl
}

func Create(gems ...*member.PublicMember) *Gem {
	gem := Gem{
		PubKeys: map[string]*member.PublicMember{},
	}
	return gem.Add(gems...)
}

func (gem *Gem) Add(gems ...*member.PublicMember) *Gem {
	for i := range gems {
		gem.PubKeys[gems[i].Id] = gems[i]
	}
	return gem
}

func (gem *Gem) Rm(ids ...string) *Gem {
	for i := range ids {
		delete(gem.PubKeys, ids[i])
	}
	return gem
}

func (gem *Gem) Ls(ids ...string) []*member.PublicMember {
	mapIds := map[string]string{}
	for i := range ids {
		mapIds[ids[i]] = ids[i]
	}
	ret := []*member.PublicMember{}
	for i := range gem.PubKeys {
		if len(mapIds) == 0 {
			ret = append(ret, gem.PubKeys[i])
		} else {
			_, found := mapIds[gem.PubKeys[i].Id]
			if found {
				ret = append(ret, gem.PubKeys[i])
			}
		}
	}
	return ret
}

func (gem *Gem) LsByType(typ member.MemberType, ids ...string) []*member.PublicMember {
	m := gem.Ls(ids...)
	out := []*member.PublicMember{}
	for i := range m {
		if strings.Compare(string(typ), string(m[i].Type[:])) == 0 {
			out = append(out, m[i])
		}
	}
	return out
}

type JsonGem struct {
	PubKeys []*member.JsonPublicMember
}

func (kvp *Gem) AsJSON() *JsonGem {
	pubks := make([]*member.JsonPublicMember, len(kvp.PubKeys))
	idx := 0
	for i := range kvp.PubKeys {
		member := kvp.PubKeys[i]
		pubks[idx] = member.AsJSON()
		idx++
	}
	sort.Sort(&member.JsonPublicMemberSorter{
		Values: pubks,
		By:     member.JsonPublicMemberValueBy,
	})
	return &JsonGem{
		PubKeys: pubks,
	}
}

func ToJsonGems(gs ...*Gem) []*JsonGem {
	ret := make([]*JsonGem, len(gs))
	for i := range gs {
		ret[i] = gs[i].AsJSON()
	}
	return ret
}

const Type = "Gem"

func (gem *Gem) ClosePearl(owners *pearl.PearlOwner) (*pearl.Pearl, error) {
	jsonStr, err := json.Marshal(gem.AsJSON())
	if err != nil {
		return nil, err
	}
	return pearl.Close(&pearl.CloseRequestPearl{
		Type:    Type,
		Payload: jsonStr,
		Owners:  *owners,
	})
}

func FromJSON(jsStr []byte) (*Gem, error) {
	jsGem := JsonGem{}
	err := json.Unmarshal(jsStr, &jsGem)
	if err != nil {
		return nil, err
	}
	gem := Create()
	// kvp.Tags = uniqStrings(kvp.Tags)
	for i := range jsGem.PubKeys {
		key, err := member.JsToPublicMember(jsGem.PubKeys[i])
		if err != nil {
			return nil, err
		}
		gem.PubKeys[key.Id] = key
	}
	return gem, nil
}

func OpenPearl(pks []*key.PrivateKey, prl *pearl.Pearl) (*Gem, error) {
	op, err := pearl.Open(pks, prl)
	if err != nil {
		return nil, err
	}
	gem, err := FromJSON(op.Payload)
	if err != nil {
		return nil, err
	}
	gem.Pearl = op
	return gem, err
}
