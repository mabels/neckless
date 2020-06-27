package pipeline

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func checkPipeline(pl *Pipeline, now time.Time, t *testing.T) {
	if len(pl.Id) != 36 {
		t.Errorf("Create do not pass name: %s", pl.Id)
	}
	if pl.Name != "Test" {
		t.Errorf("Create do not pass name: %s", pl.Name)
	}
	if pl.Created.Before(now) {
		t.Errorf("created has to be younger: %s", pl.Created.String())
	}
	if pl.Updated.Before(now) {
		t.Errorf("created has to be younger: %s", pl.Updated.String())
	}
	if pl.ValidUntil.Before(now.AddDate(5, 0, 0)) {
		t.Errorf("created has to be younger: %s", pl.Updated.String())
	}
	if len(pl.KeyPair.Priv) != 72 {
		t.Errorf("created key: %s", pl.KeyPair.Priv)
	}
	if len(pl.KeyPair.Publ) != 64 {
		t.Errorf("created key: %s", pl.KeyPair.Publ)
	}
}

func TestPipeLineCreate(t *testing.T) {
	p := PipelineArgs{name: "Test"}
	now := time.Now()
	pl := Create(p)
	checkPipeline(pl, now, t)
	jsonStr, _ := json.Marshal(pl)
	out := Pipeline{}
	json.Unmarshal(jsonStr, &out)
	checkPipeline(&out, now, t)
}

func TestSignedToken(t *testing.T) {
	p := PipelineArgs{name: "Test"}
	pl := Create(p)
	jwt := SignMember(pl, pl)
	// t.Error(jwt)
	claim, _, err := VerifyAndClaim(jwt, pl)
	if err != nil {
		t.Error(err)
	}
	if claim.SignerName != pl.Name {
		t.Error(claim.SignerName)
	}
	if claim.SignerPubkey != string(pl.KeyPair.Publ) {
		t.Error(claim.SignerPubkey)
	}
}

func TestEncrypt(t *testing.T) {
	bob := Create(PipelineArgs{name: "Bob"})
	o := toByte32(bob.KeyPair.Priv)
	if bytes.Compare(o[:], bob.KeyPair.Priv) != 0 {
		t.Error("toByte32:", string(o[:]), string(bob.KeyPair.Priv))
	}
	alice := Create(PipelineArgs{name: "Alice"})
	// t.Error(pl.KeyPair.Publ, len(pl.KeyPair.Publ))
	secMsg := EncryptFor(t, &bob.KeyPair, alice.KeyPair.Publ, "superGeheim")
	t.Error(secMsg)
	usecMsg, err := DecryptFor(t, &alice.KeyPair, bob.KeyPair.Publ, secMsg)
	t.Error(err)
	if "superGeheim" != string(usecMsg) {
		t.Error("No matching", usecMsg)
	}
	if strings.Contains(string(secMsg), string(usecMsg)) {
		t.Error("should not happend", secMsg, usecMsg)
	}

}
