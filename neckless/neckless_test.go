package neckless

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"neckless.adviser.com/pearl"
)

func TestGetAndOpenEmpty(t *testing.T) {
	fname := uuid.New().String()
	nl := GetAndOpen(fname)
	if len(nl.Pearls) != 0 {
		t.Error("should not happend")
	}
	if strings.Compare(nl.FileName, fname) != 0 {
		t.Error("hello happend")
	}
}

func TestClose(t *testing.T) {
	nl := GetAndOpen(uuid.New().String())
	js, err := nl.Save()
	if err != nil {
		t.Error("should not happend", err)
	}
	if !bytes.Equal([]byte("[]"), js) {
		t.Error("should not happend", js)
	}
}

func TestFileSideEffect(t *testing.T) {
	nl := GetAndOpen(uuid.New().String())
	nl.Save(nl.FileName)
	nl = GetAndOpen(nl.FileName)
	if len(nl.Pearls) != 0 {
		t.Error("illegal len")
	}
	nl.Reset(&pearl.Pearl{
		FingerPrint: []byte("Id1"),
	})
	nl.Reset(&pearl.Pearl{
		FingerPrint: []byte("Id2"),
	})
	nl.Save(nl.FileName)
	nl = GetAndOpen(nl.FileName)
	if len(nl.Pearls) != 2 {
		t.Error("illegal len", nl)
	}
	nl.Rm([]byte("Id2"))
	nl.Save(nl.FileName)
	nl = GetAndOpen(nl.FileName)
	if len(nl.Pearls) != 1 {
		t.Error("illegal len")
	}
	os.Remove(nl.FileName)
}

func TestSet(t *testing.T) {
	nl := GetAndOpen(uuid.New().String())
	if len(nl.Pearls) != 0 {
		t.Error("not the right len")
	}
	nl.Reset(&pearl.Pearl{
		FingerPrint: []byte("Id1"),
	})
	if len(nl.Pearls) != 1 {
		t.Error("not the right len")
	}
	nl.Reset(&pearl.Pearl{
		FingerPrint: []byte("Id2"),
	})
	if len(nl.Pearls) != 2 {
		t.Error("not the right len")
	}
	nl.Reset(&pearl.Pearl{
		FingerPrint: []byte("Id1"),
	})
	if len(nl.Pearls) != 2 {
		t.Error("not the right len", nl.Pearls)
	}

	nl.Reset(&pearl.Pearl{
		FingerPrint: []byte("Id9"),
	}, []byte("Id1"))

	if len(nl.Pearls) != 2 {
		t.Error("not the right len", nl.Pearls)
	}

	nl.Rm([]byte("IdX"))
	if len(nl.Pearls) != 2 {
		t.Error("not the right len")
	}
	nl.Rm([]byte("Id2"))
	if len(nl.Pearls) != 1 {
		t.Error("not the right len")
	}
	nl.Rm([]byte("Id1"))
	if len(nl.Pearls) != 0 {
		t.Error("not the right len")
	}

}
