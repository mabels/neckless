package pearl

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"neckless.adviser.com/asymmetric"
	"neckless.adviser.com/key"
	"neckless.adviser.com/symmetric"
)

type PearlClaim struct {
	PayloadChecksum          string `json:"pcs"`
	EncryptedPayloadPassword string `json:"epp"`
	jwt.StandardClaims
}

type ClosedPearlAttribute struct {
}

type JWTokenPearlClaim string

type CreatorOwners struct {
	Creator key.PublicKey // public-key
	Tokens  []JWTokenPearlClaim
}

type JsonCreatorOwners struct {
	Creator string
	Tokens  []JWTokenPearlClaim
}

func (co *CreatorOwners) AsJson() *JsonCreatorOwners {
	return &JsonCreatorOwners{
		Creator: co.Creator.Marshal(),
		Tokens:  co.Tokens,
	}
}
func (jcp *JsonCreatorOwners) FromJson() (*CreatorOwners, error) {
	_, creator, err := key.FromText(jcp.Creator, fmt.Sprintf("SYN-%s", uuid.New().String()))
	if err != nil {
		return nil, err
	}
	if creator == nil {
		return nil, errors.New("we need this as public key")
	}
	return &CreatorOwners{
		Creator: *creator,
		Tokens:  jcp.Tokens,
	}, nil
}

type Pearl struct {
	FingerPrint []byte
	Type        string
	Payload     []byte
	Owners      CreatorOwners
}

type JsonPearl struct {
	FingerPrint string
	Type        string
	Payload     string // base64
	Owners      JsonCreatorOwners
}

func (p *Pearl) AsJson() *JsonPearl {
	return &JsonPearl{
		FingerPrint: base64.StdEncoding.EncodeToString(p.FingerPrint),
		Type:        p.Type,
		Payload:     base64.StdEncoding.EncodeToString(p.Payload),
		Owners:      *p.Owners.AsJson(),
	}
}

type CalcFingerPrint struct {
	Signer  *key.PublicKey
	Payload []byte
	Type    string
}

func calcFingerprint(cfp *CalcFingerPrint) ([]byte, error) {
	sum := sha256.New()
	_, err := sum.Write([]byte(cfp.Type))
	if err != nil {
		return nil, err
	}
	_, err = sum.Write(cfp.Payload)
	if err != nil {
		return nil, err
	}
	_, err = sum.Write(cfp.Signer.Key.Raw[:])
	if err != nil {
		return nil, err
	}
	return sum.Sum(nil), nil
}

func setFingerprint(cr *CloseRequestPearl, p *Pearl) (*Pearl, error) {
	sum, err := calcFingerprint(&CalcFingerPrint{
		Signer:  cr.Owners.Signer.Public(),
		Payload: cr.Payload,
		Type:    cr.Type,
	})
	if err != nil {
		return nil, err
	}
	p.FingerPrint = sum
	return p, nil
}

func checkFingerprint(p *OpenPearl) (*OpenPearl, error) {
	sum, err := calcFingerprint(&CalcFingerPrint{
		Type:    p.Closed.Type,
		Signer:  &p.Closed.Owners.Creator,
		Payload: p.Payload,
	})
	if err != nil {
		return nil, err
	}
	if bytes.Compare(p.Closed.FingerPrint, sum) != 0 {
		return nil, errors.New("checksum missmatch")
	}
	return p, nil
}

func (jp *JsonPearl) FromJson() (*Pearl, error) {
	payload, err := base64.StdEncoding.DecodeString(jp.Payload)
	if err != nil {
		return nil, err
	}
	fingerprint, err := base64.StdEncoding.DecodeString(jp.FingerPrint)
	if err != nil {
		return nil, err
	}
	owners, err := jp.Owners.FromJson()
	if err != nil {
		return nil, err
	}

	return &Pearl{
		FingerPrint: fingerprint,
		Type:        jp.Type,
		Payload:     payload,
		Owners:      *owners,
	}, nil
}

type PearlOwner struct {
	Signer *key.PrivateKey
	Owners []*key.PublicKey
}

type CloseRequestPearl struct {
	Type    string
	Payload []byte
	Owners  PearlOwner
}

// func ownersToSeal(pks *map[string]key.PublicKey) [][]byte {
// 	ret := make([][]byte, len(*pks))
// 	for i := range *pks {
// 		ret[i] = []byte(*pks[i].Marshal())
// 	}
// 	return ret
// }

type CloseContainer struct {
	Checksum   []byte
	PayloadKey *key.RawKey
}

func createJWTPearlClaim(signer *key.PrivateKey, pub *key.PublicKey, cl *CloseContainer) (*JWTokenPearlClaim, error) {
	privPubKey := asymmetric.CreateShared(&signer.Key.Raw, &pub.Key.Raw)
	// fmt.Printf("Close:%x:%x=>%x\n", signer.Key.Raw, pub.Key.Raw, privPubKey)
	sealedPwd, err := symmetric.Seal(&symmetric.SealRequest{
		Checksum: cl.Checksum,
		Key:      privPubKey,
		Payload:  cl.PayloadKey[:],
	})
	if err != nil {
		return nil, err
	}
	ownerClaim := PearlClaim{
		PayloadChecksum:          base64.StdEncoding.EncodeToString(cl.Checksum),
		EncryptedPayloadPassword: base64.StdEncoding.EncodeToString(sealedPwd.Payload),
		StandardClaims: jwt.StandardClaims{
			Id:       pub.Key.Id,
			Issuer:   signer.Key.Id,
			IssuedAt: time.Now().Unix(),
			Subject:  "PearlClaim",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ownerClaim)
	tokenString, err := token.SignedString(privPubKey[:])
	if err != nil {
		return nil, err
	}
	jwtToken := JWTokenPearlClaim(tokenString)
	return &jwtToken, err

}

func creatorOwners(pk *key.PrivateKey, owners []*key.PublicKey, cl *CloseContainer) (*CreatorOwners, error) {
	jwted := make([]JWTokenPearlClaim, len(owners))
	for i := range owners {
		jwt, err := createJWTPearlClaim(pk, owners[i], cl)
		if err != nil {
			return nil, err
		}
		jwted[i] = *jwt
	}
	return &CreatorOwners{
		Creator: *pk.Public(),
		Tokens:  jwted,
	}, nil
}

// Close a pearl with the EncryptedPayload and Owners
func Close(opa *CloseRequestPearl) (*Pearl, error) {
	payloadKey, err := key.CreateRandomKey()
	if err != nil {
		return nil, err
	}
	sealed, err := symmetric.Seal(symmetric.Checksum(&symmetric.SealRequest{
		Key:     *payloadKey,
		Payload: opa.Payload,
	}))

	if err != nil {
		return nil, err
	}
	co, err := creatorOwners(opa.Owners.Signer, opa.Owners.Owners, &CloseContainer{
		PayloadKey: payloadKey,
		Checksum:   sealed.Checksum,
	})
	if err != nil {
		return nil, err
	}
	return setFingerprint(opa, &Pearl{
		Type:    opa.Type,
		Payload: sealed.Payload,
		Owners:  *co,
	})
}

func encryptPayloadKey(sk *key.RawKey, csum []byte, b64 string) (*symmetric.OpenContainer, error) {
	encryptedPayloadKey, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return symmetric.Open(sk, &symmetric.SealedContainer{
		Checksum: csum,
		Payload:  encryptedPayloadKey,
	}, symmetric.SkipVerify)
}

type OpenPearl struct {
	Closed  Pearl
	Payload []byte
	Claim   PearlClaim
}

func (pea *Pearl) tryOpenWithKey(pk *key.PrivateKey) (*key.RawKey, *jwt.Token, *PearlClaim, bool) {
	creatorPubKey := pea.Owners.Creator
	for i := range pea.Owners.Tokens {
		claims := PearlClaim{}
		sharedKey := asymmetric.CreateShared(&pk.Key.Raw, &creatorPubKey.Key.Raw)
		token, err := jwt.ParseWithClaims(string(pea.Owners.Tokens[i]), &claims,
			func(token *jwt.Token) (interface{}, error) {
				return sharedKey[:], nil
			})
		// fmt.Printf("Open:%x:%x=>%x:%s\n", pk.Key.Raw, creatorPubKey, sharedKey, err)
		if err == nil {
			return &sharedKey, token, &claims, true
		}
	}
	return nil, nil, nil, false
}

// Close creates a pearl with the EncryptedPayload and Owners
func Open(pks []*key.PrivateKey, pea *Pearl) (*OpenPearl, error) {
	errs := []string{}
	for i := range pks {
		op, err := OpenOne(pks[i], pea)
		if err == nil {
			return op, nil
		}
		errs = append(errs, err.Error())
	}
	return nil, errors.New(fmt.Sprintf("can't open this pearl:%x:[%s]", pea.FingerPrint, strings.Join(errs, "],[")))
}

func OpenOne(pk *key.PrivateKey, pea *Pearl) (*OpenPearl, error) {
	sharedKey, _, claim, ok := pea.tryOpenWithKey(pk)
	if !ok {
		return nil, fmt.Errorf("id not found in owners:[%x]", pk.Key.Id)
	}
	payloadChecksum, err := base64.StdEncoding.DecodeString(claim.PayloadChecksum)
	if err != nil {
		return nil, err
	}
	payloadKey, err := encryptPayloadKey(sharedKey, payloadChecksum, claim.EncryptedPayloadPassword)
	if err != nil {
		return nil, err
	}
	opc, err := symmetric.Open(key.AsRawKey(payloadKey.Payload), &symmetric.SealedContainer{
		Checksum: payloadChecksum,
		Payload:  pea.Payload,
	}, symmetric.Verify)
	if err != nil {
		return nil, err
	}
	return checkFingerprint(&OpenPearl{
		Closed:  *pea,
		Payload: opc.Payload,
		Claim:   *claim,
	})

	// chkSum, err :=  base64.StdEncoding.DecodeString(claim.PayloadChecksum)
	// if err != nil {
	// 	return nil, err
	// }
	// _, pb, err := key.FromText(claim.PublicKey)
	// if err != nil {
	// 	return nil, err
	// }
	// if pb == nil {
	// 	return nil, errors.New("no public key provided")
	// }

	// if bytes.Equal(chkSum, sha256.Sum256())

	// return nil
}
