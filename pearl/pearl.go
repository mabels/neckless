package pearl

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"neckless.adviser.com/asymmetric"
	"neckless.adviser.com/keys"
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
	Creator keys.RawKey // public-key
	Tokens  []JWTokenPearlClaim
}

type Pearl struct {
	Type    string
	Payload []byte
	Owners  CreatorOwners
}

type JsonPearl struct {
	Type    string
	Payload string // base64
	// Owners  []JWTokenPearlClaim
}

type CloseRequestPearl struct {
	Type    string
	Payload []byte
	Signer  keys.PrivateKey
	Owners  []keys.PublicKey
}

// func ownersToSeal(pks *map[string]keys.PublicKey) [][]byte {
// 	ret := make([][]byte, len(*pks))
// 	for i := range *pks {
// 		ret[i] = []byte(*pks[i].Marshal())
// 	}
// 	return ret
// }

type CloseContainer struct {
	Checksum   []byte
	PayloadKey *keys.RawKey
}

func createJWTPearlClaim(signer *keys.PrivateKey, pub *keys.PublicKey, cl *CloseContainer) (*JWTokenPearlClaim, error) {
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

func creatorOwners(pk *keys.PrivateKey, owners *[]keys.PublicKey, cl *CloseContainer) (*CreatorOwners, error) {
	jwted := make([]JWTokenPearlClaim, len(*owners))
	for i := range *owners {
		jwt, err := createJWTPearlClaim(pk, &(*owners)[i], cl)
		if err != nil {
			return nil, err
		}
		jwted[i] = *jwt
	}
	return &CreatorOwners{
		Creator: pk.Public().Key.Raw,
		Tokens:  jwted,
	}, nil
}

// Close a pearl with the EncryptedPayload and Owners
func Close(opa *CloseRequestPearl) (*Pearl, error) {
	payloadKey, err := keys.CreateRandomKey()
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
	co, err := creatorOwners(&opa.Signer, &opa.Owners, &CloseContainer{
		PayloadKey: payloadKey,
		Checksum:   sealed.Checksum,
	})
	if err != nil {
		return nil, err
	}
	return &Pearl{
		Type:    opa.Type,
		Payload: sealed.Payload,
		Owners:  *co,
	}, nil
}

func encryptPayloadKey(sk *keys.RawKey, csum []byte, b64 string) (*symmetric.OpenContainer, error) {
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
	Type    string
	Payload []byte
	Claim   PearlClaim
}

func (pea *Pearl) findByKeyId(pk *keys.PrivateKey) (*keys.RawKey, *jwt.Token, *PearlClaim, bool) {
	creatorPubKey := pea.Owners.Creator
	for i := range pea.Owners.Tokens {
		claims := PearlClaim{}
		sharedKey := asymmetric.CreateShared(&pk.Key.Raw, &creatorPubKey)
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
func Open(pk *keys.PrivateKey, pea *Pearl) (*OpenPearl, error) {
	sharedKey, _, claim, ok := pea.findByKeyId(pk)
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
	opc, err := symmetric.Open(keys.AsRawKey(payloadKey.Payload), &symmetric.SealedContainer{
		Checksum: payloadChecksum,
		Payload:  pea.Payload,
	}, symmetric.Verify)
	if err != nil {
		return nil, err
	}
	return &OpenPearl{
		Type:    pea.Type,
		Payload: opc.Payload,
		Claim:   *claim,
		// Owner:  me,
		// Owners: pea.Owners,
	}, nil

	// chkSum, err :=  base64.StdEncoding.DecodeString(claim.PayloadChecksum)
	// if err != nil {
	// 	return nil, err
	// }
	// _, pb, err := keys.FromText(claim.PublicKey)
	// if err != nil {
	// 	return nil, err
	// }
	// if pb == nil {
	// 	return nil, errors.New("no public key provided")
	// }

	// if bytes.Equal(chkSum, sha256.Sum256())

	// return nil
}
