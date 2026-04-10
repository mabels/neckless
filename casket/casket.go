package casket

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/mabels/neckless/member"
)

type CreateArg struct {
	member.MemberArg
	DryRun bool    // if dryrun don't write
	Fname  *string //
}

type CasketAttribute struct {
	CasketFname *string   `json:"-"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type Casket struct {
	CasketAttribute
	Members map[string](*member.PrivateMember)
}

func getcasketFilename(fnames []string) (string, error) {
	var fname string
	if len(fnames) > 0 {
		fname = fnames[0]
	} else {
		fname = path.Join(os.Getenv("HOME"), ".neckless/casket")
	}
	err := os.MkdirAll(path.Dir(fname), 0700)
	return fname, err
}

func readcasket(fname string) (*Casket, error) {
	now := time.Now()
	jsonCasket := JsonCasket{
		CasketAttribute: CasketAttribute{
			Created:     now,
			Updated:     now,
			CasketFname: &fname,
		},
		Members: map[string]member.JsonPrivateMember{},
	}
	dat, err := os.ReadFile(fname)
	if err == nil {
		err = json.Unmarshal(dat, &jsonCasket)
		if err != nil {
			return nil, err
		}
	} else {
		err = nil
	}
	members := map[string]*member.PrivateMember{}
	for k := range jsonCasket.Members {
		jspm := jsonCasket.Members[k]
		pm, err := jspm.AsPrivateMember()
		if err != nil {
			return nil, err
		}
		// fmt.Printf("Ls:Key:%s\n", k)
		members[k] = pm
	}
	// fmt.Printf("Ls:%d\n", len(members))
	return &Casket{
		CasketAttribute: jsonCasket.CasketAttribute,
		Members:         members,
	}, err
}

type JsonCasket struct {
	Members map[string]member.JsonPrivateMember `json:"members"`
	CasketAttribute
}

func (casket *Casket) AsJSON() *JsonCasket {
	jsonMembers := map[string]member.JsonPrivateMember{}
	for i := range casket.Members {
		val := casket.Members[i]
		jsonMembers[i] = *val.AsJSON()
	}
	return &JsonCasket{
		CasketAttribute: casket.CasketAttribute,
		Members:         jsonMembers,
	}
}

func writecasket(casket *Casket) error {
	jsstr, err := json.MarshalIndent(casket.AsJSON(), "", "  ")
	if err != nil {
		return err
	}
	tmpf, err := os.CreateTemp(path.Dir(*casket.CasketFname), ".casket-*.tmp")
	if err != nil {
		return err
	}
	tmp := tmpf.Name()
	if err = os.Chmod(tmp, 0600); err != nil {
		os.Remove(tmp)
		return err
	}
	if _, err = tmpf.Write(jsstr); err != nil {
		tmpf.Close()
		os.Remove(tmp)
		return err
	}
	if err = tmpf.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err = os.Rename(tmp, *casket.CasketFname); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// UseCase Write the PrivateKey in den casket ~/.neckless/casket
// neckless casket create --name <name> [--device <name>] [--person|--device] [--file=~/.crazybee/casket]
func Create(ca CreateArg) (*Casket, *member.PrivateMember, error) {
	pk, err := member.MakePrivateMember(&member.PrivateMemberArg{
		Member: ca.MemberArg,
	})
	if err != nil {
		return nil, nil, err
	}
	var casket *Casket
	if ca.Fname == nil || len(*ca.Fname) == 0 {
		casket, err = Ls()
	} else {
		casket, err = Ls(*ca.Fname)
	}
	if err != nil {
		return nil, nil, err
	}
	casket.Members[pk.Id] = pk
	casket.Updated = time.Now()
	if !ca.DryRun {
		err = writecasket(casket)
		if err != nil {
			return nil, nil, err
		}
	}
	return casket, pk, nil
}

// UseCase List casket
// neckless casket ls
func Ls(fnames ...string) (*Casket, error) {
	_, present := os.LookupEnv("NECKLESS_PRIVKEY")
	if !present {
		fname, err := getcasketFilename(fnames)
		if err != nil {
			return nil, err
		}
		return readcasket(fname)
	}
	fname := "ENV:NECKLESS_PRIVKEY"
	return &Casket{
		CasketAttribute: CasketAttribute{
			CasketFname: &fname,
			Created:     time.Time{},
			Updated:     time.Time{},
		},
		Members: map[string]*member.PrivateMember{},
	}, nil
}

func (c *Casket) AsPrivateMembers() []*member.PrivateMember {
	ret := make([]*member.PrivateMember, len(c.Members))
	idx := 0
	for i := range c.Members {
		ret[idx] = c.Members[i]
		idx++
	}
	return ret
}

type RmArg struct {
	Ids    []string
	DryRun bool    // if dryrun don't write
	Fname  *string //
}

// UseCase Delete Key from casket
// neckless casket rm <id>
func Rm(rmarg RmArg) (*Casket, []*member.PrivateMember, error) {
	var ks *Casket
	var err error
	if rmarg.Fname != nil {
		ks, err = Ls(*rmarg.Fname)
	} else {
		ks, err = Ls()
	}
	if err != nil {
		return nil, nil, err
	}
	out := []*member.PrivateMember{}
	for i := range rmarg.Ids {
		id := rmarg.Ids[i]
		pk, ok := ks.Members[id]
		if ok {
			delete(ks.Members, id)
			out = append(out, pk)
		}
	}
	if !rmarg.DryRun {
		if err = writecasket(ks); err != nil {
			return nil, nil, err
		}
	}
	return ks, out, nil
}
