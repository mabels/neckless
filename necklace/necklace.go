package necklace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"neckless.adviser.com/pearl"
)

type jsonClosedNecklace []pearl.JsonPearl

type Necklace struct {
	FileName string
	Pearls   []*pearl.Pearl
}

func GetAndOpen(fname string) (Necklace, []error) {
	warns := []error{}
	dat, err := ioutil.ReadFile(fname)
	nl := Necklace{
		FileName: fname,
		Pearls:   []*pearl.Pearl{},
	}
	if err == nil {
		jsn := jsonClosedNecklace{}
		err := json.Unmarshal(dat, &jsn)
		if err != nil {
			warns = append(warns, err)
		}
		for i := range jsn {
			my, err := jsn[i].FromJson()
			if err != nil {
				warns = append(warns, err)
			} else {
				nl.Pearls = append(nl.Pearls, my)
			}
		}
	} else {
		warns = append(warns, err)
	}
	return nl, warns
}

func (nl *Necklace) Reset(p *pearl.Pearl, updateFprs ...[]byte) *Necklace {
	foundIt := [][]byte{}
	mapFpr := map[string]struct{}{}
	mapFpr[fmt.Sprintf("%x", p.FingerPrint)] = struct{}{}
	for i := range updateFprs {
		mapFpr[fmt.Sprintf("%x", updateFprs[i])] = struct{}{}
	}
	for i := range nl.Pearls {
		_, found := mapFpr[fmt.Sprintf("%x", nl.Pearls[i].FingerPrint)]
		// fmt.Println(found, fmt.Sprintf("%x", nl.Pearls[i].FingerPrint))
		if found {
			foundIt = append(foundIt, nl.Pearls[i].FingerPrint)
		}
	}
	// fmt.Println(p.FingerPrint, foundIt)
	nl.Rm(foundIt...)
	nl.Pearls = append(nl.Pearls, p)
	return nl
}

func (nl *Necklace) Rm(fprs ...[]byte) *Necklace {
	founds := []int{}
	mapFpr := map[string]struct{}{}
	for i := range fprs {
		mapFpr[fmt.Sprintf("%x", fprs[i])] = struct{}{}
	}
	for i := range nl.Pearls {
		_, found := mapFpr[fmt.Sprintf("%x", nl.Pearls[i].FingerPrint)]
		if found {
			founds = append(founds, i)
			break
		}
	}
	for i := len(founds) - 1; i >= 0; i-- {
		nl.Pearls = append(nl.Pearls[:founds[i]], nl.Pearls[founds[i]+1:]...)
	}
	return nl
}

func (nl *Necklace) Save(fnames ...string) ([]byte, error) {
	out := make([]*pearl.JsonPearl, len(nl.Pearls))
	for i := range nl.Pearls {
		my := nl.Pearls[i]
		out[i] = my.AsJson()
	}
	jsStr, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, err
	}
	if len(fnames) > 0 {
		err := ioutil.WriteFile(fnames[0], jsStr, 0644)
		if err != nil {
			return nil, err
		}
	}
	return jsStr, nil
}

func (nl *Necklace) FilterByType(typ string) []*pearl.Pearl {
	out := []*pearl.Pearl{}
	for i := range nl.Pearls {
		if strings.Compare(nl.Pearls[i].Type, typ) == 0 {
			out = append(out, nl.Pearls[i])
		}
	}
	return out
}

// 	out := []*pearl.OpenPearl{}
// 	for i := range closedNecklace {
// 		closedPearl := closedNecklace[i]
// 		for j := range pks {
// 			pk := pks[j]
// 			opearl, err := pearl.Open(pk, closedPearl)
// 			if err == nil {
// 				out = append(out, opearl)
// 				break
// 			}
// 		}
// 	}
// 	return out
// }
