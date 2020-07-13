package necklace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"neckless.adviser.com/pearl"
)

type jsonClosedNecklace []pearl.JSONPearl

// Necklace is a chain of Pearls
type Necklace struct {
	FileName string
	Pearls   []*pearl.Pearl
}

// Read reads a Necklace
func Read(fname string) (Necklace, []error) {
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
			my, err := jsn[i].FromJSON()
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

// Reset Set or Add a Pearl to the Necklace. UpdateFprs enables
// to replace an existing Pearl in a Necklace Chain
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

// Rm removes Pearls from an Necklace
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

// Save saves a Neckless to the given filename. If no file is provide
// it won't right any file
func (nl *Necklace) Save(fnames ...string) ([]byte, error) {
	out := make([]*pearl.JSONPearl, len(nl.Pearls))
	for i := range nl.Pearls {
		my := nl.Pearls[i]
		out[i] = my.AsJSON()
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

// FilterByType filters out of the Necklace with corresponding type
// to a chain of Pearls
func (nl *Necklace) FilterByType(typ string) []*pearl.Pearl {
	out := []*pearl.Pearl{}
	for i := range nl.Pearls {
		if strings.Compare(nl.Pearls[i].Type, typ) == 0 {
			out = append(out, nl.Pearls[i])
		}
	}
	return out
}
