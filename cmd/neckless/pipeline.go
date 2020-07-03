package neckless

import (
	"time"
)

// // Export Key as PublicJWT respect Signing Key
// // Create Neckless
// neckless create [--file=neckless.private]
// // Export Neckless Key
// neckless export [--file=neckless.private]  --private
// // Export Neckless Public
// neckless export [--file=neckless.public]
// // SetAddSecret to Neckless
// neckless secret set [--file=neckless.public] --key= --value
// // RemoveSecret to Neckless
// neckless secret rm [--file=neckless.public] --key=
// // ListSecret from Neckless with Key
// neckless secret ls [--file=neckless.public]
// // AddUser to Neckless
// neckless user add [--file=neckless.public] --user=user.file --private=[id]

// PipelineArgs Command Args
// type  struct {
// 	name string
// 	// valid period.Period
// }

// // KeyPair is the pub and privkey container
// // KeyPair is the pub and privkey container
// type KeyPair struct {
// 	Priv PublicKey
// 	Publ PrivateKey
// }

// func (*KeyPair) PublAsKey() []byte {
// 	return []byte{}
// }

// func (*KeyPair) PrivAsKey() []byte {
// 	return []byte{}
// }

// Pipeline is The Created Entry for the PipeLine
type Pipeline struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	KeyPair    KeyType   `json:"keyPair"`
	ValidUntil time.Time `json:"valid"`
	Updated    time.Time `json:"updated"`
	Created    time.Time `json:"created"`
}

// Key is curve25519 key.
type Key [KeySize]byte

// // Create is used to Create a Pipeline
// func Create(arg PipelineArgs) *Pipeline {
// 	now := time.Now()
// 	pk, err := NewPrivateKey()
// 	if err != nil {
// 		log.Fatal("can not create new private key")
// 	}
// 	pkTxt, err := pk.MarshalText()
// 	if err != nil {
// 		log.Fatal("can not marshal private to text")
// 	}
// 	pubKey := pk.Public()
// 	pbTxt, err := pubKey.MarshalText()
// 	if err != nil {
// 		log.Fatal("can not marshal public to text")
// 	}
// 	return &Pipeline{
// 		Id:   uuid.New().String(),
// 		Name: arg.name,
// 		KeyPair: KeyPair{
// 			Priv: pkTxt,
// 			Publ: pbTxt,
// 		},
// 		ValidUntil: now.AddDate(5, 0, 0),
// 		Updated:    now,
// 		Created:    now,
// 	}
// }

// type MemberClaim struct {
// 	SignerName   string `json:"signerName"`
// 	SignerPubkey string `json:"signerPubKey"`
// 	jwt.StandardClaims
// }

// func SignMember(signer *Pipeline, member *Pipeline) string {
// 	claims := &MemberClaim{
// 		SignerName:   member.Name,
// 		SignerPubkey: string(member.KeyPair.Publ),
// 		StandardClaims: jwt.StandardClaims{
// 			// In JWT, the expiry time is expressed as unix milliseconds
// 			ExpiresAt: signer.ValidUntil.Unix(),
// 			Id:        uuid.New().String(),
// 			IssuedAt:  time.Now().Unix(),
// 			Subject:   "SignerClaim",
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(signer.KeyPair.Priv)
// 	if err != nil {
// 		log.Fatal("Siging failed")
// 	}
// 	return tokenString
// }

// func VerifyAndClaim(tknStr string, pl *Pipeline) (*MemberClaim, *jwt.Token, error) {
// 	claims := MemberClaim{}
// 	token, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
// 		return pl.KeyPair.Priv, nil
// 	})
// 	return &claims, token, err
// }

// func toByte32(p []byte) [32]byte {
// 	var ret [32]byte
// 	for i := range ret[:] {
// 		ret[i] = p[i]
// 	}
// 	return ret
// }

// func EncryptFor(t *testing.T, keyAlice *KeyPair, pubKeyBob []byte, msg string) []byte {
// 	prA := toByte32(keyAlice.PrivAsKey())
// 	puB := toByte32(pubKeyBob)
// 	var shared [32]byte
// 	curve25519.ScalarMult(&shared, &prA, &puB)
// 	t.Error(shared)
// 	aead, _ := chacha20poly1305.NewX(shared[:])
// 	nonce := make([]byte, chacha20poly1305.NonceSizeX)
// 	return aead.Seal(nil, nonce, []byte(msg), nil)
// }

// func DecryptFor(t *testing.T, keyBob *KeyPair, pubKeyAlice []byte, msg []byte) ([]byte, error) {
// 	prB := toByte32(keyBob.PrivAsKey())
// 	puA := toByte32(pubKeyAlice)
// 	var shared [32]byte
// 	curve25519.ScalarMult(&shared, &prB, &puA)
// 	t.Error(shared)
// 	aead, _ := chacha20poly1305.NewX(shared[:])
// 	nonce := make([]byte, chacha20poly1305.NonceSizeX)
// 	return aead.Open(nil, nonce, []byte(msg), nil)
// }
