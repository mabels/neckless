package kvpearl

import (
	"fmt"

	"github.com/xlzd/gotp"
)

// FuncsAndParam is the Parsed toResolve structure
type FuncsAndParam struct {
	Param string
	Funcs []string
}

// ActionFN is the type of an action function
type ActionFN func(string) (string, error)

func noopFN(val string) (string, error) {
	return val, nil
}

func lenFN(val string) (string, error) {
	return fmt.Sprintf("%d", len(val)), nil
}

func totpFN(val string) (string, error) {
	totp := gotp.NewDefaultTOTP(val)
	return totp.Now(), nil
}

var actions = map[string]ActionFN{
	"Noop": noopFN,
	"Len":  lenFN,
	"Totp": totpFN,
}

// RunFuncs applies the given function to the value
func (fap *FuncsAndParam) RunFuncs(inVal string) (string, error) {
	len := len(fap.Funcs) - 1
	val := inVal
	for o := range fap.Funcs {
		action := fap.Funcs[len-o]
		fn, found := actions[action]
		if !found {
			return inVal, fmt.Errorf("unknown action:%s", action)
		}
		var err error
		val, err = fn(val)
		if err != nil {
			return val, fmt.Errorf("action:%s error:%s", action, err)
		}
	}
	return val, nil
}
