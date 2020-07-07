package member

import (
	"encoding/json"
	"errors"
	"time"

	"neckless.adviser.com/key"
)

type MemberType string

const (
	Device = MemberType("Device")
	Person = MemberType("Person")
)

type MemberArg struct {
	Id         string
	Type       MemberType
	Name       string
	Device     string
	ValidUntil *time.Time
	Updated    *time.Time `json:"updated"`
	Created    *time.Time `json:"created"`
}

type MemberBase struct {
	Type       MemberType `json:"type"`
	Name       string     `json:"name"`
	Device     string     `json:"device,omitempty"`
	ValidUntil time.Time  `json:"validUntil"`
	Updated    time.Time  `json:"updated"`
	Created    time.Time  `json:"created"`
}

type Member struct {
	MemberBase
	Id string
}

// type JsonMember struct {
// 	MemberBase
// 	Id string `json:"id"`
// }

type PrivateMemberArg struct {
	Member     MemberArg
	PrivateKey *key.PrivateKey
}

type PrivateMember struct {
	Member
	PrivateKey key.PrivateKey
}

type PublicMember struct {
	Member
	PublicKey key.PublicKey
}

func NewMember(m *MemberArg) (*Member, error) {
	ret := Member{}
	ret.Id = m.Id
	// if len(m.Id) == 0 {
	// ret.Id = uuid.New().String()
	// } else {
	// ret.Id = m.Id
	// }

	ret.Type = m.Type
	if len(m.Name) == 0 {
		return &ret, errors.New("Require name")
	}
	ret.Name = m.Name
	if len(m.Device) != 0 {
		ret.Device = m.Device
	}
	now := time.Now()
	if m.ValidUntil != nil {
		ret.ValidUntil = *m.ValidUntil
	} else {
		ret.ValidUntil = now.AddDate(5, 0, 0)
	}
	if m.Updated != nil {
		ret.Updated = *m.Updated
	} else {
		ret.Updated = now
	}
	if m.Created != nil {
		ret.Created = *m.Created
	} else {
		ret.Created = now
	}
	return &ret, nil
}

func MakePrivateMember(pm *PrivateMemberArg) (*PrivateMember, error) {
	m, err := NewMember(&pm.Member)
	if err != nil {
		return nil, err
	}
	pk, err := key.NewPrivateKey(pm.PrivateKey)
	if err != nil {
		return nil, err
	}
	if len(m.Id) > 0 {
		pk.Key.Id = m.Id
	} else {
		m.Id = pk.Key.Id
	}
	return &PrivateMember{
		Member:     *m,
		PrivateKey: *pk,
	}, nil
}

func (pm *PrivateMember) Public() *PublicMember {
	return &PublicMember{
		Member:    pm.Member,
		PublicKey: *pm.PrivateKey.Public(),
	}
}

func MakePublicMember(pm *PublicMember) (*PublicMember, error) {
	return &PublicMember{
		Member:    pm.Member,
		PublicKey: pm.PublicKey,
	}, nil
}

type JsonPublicMember struct {
	Member
	PublicKey string `json:"publicKey"`
}

// func (mb *Member) AsJson() *JsonMember {
// 	return &JsonMember{
// 		MemberBase: mb.MemberBase,
// 		Id:         mb.Id.Id,
// 	}
// }

// func (mb *JsonMember) AsMember() *Member {
// 	return &Member{
// 		MemberBase: mb.MemberBase,
// 		Id: mb.I
// 	}
// }

func (pm *PublicMember) AsJson() *JsonPublicMember {
	return &JsonPublicMember{
		Member:    pm.Member,
		PublicKey: pm.PublicKey.Marshal(),
	}
}

func (pm *JsonPublicMember) String() ([]byte, error) {
	return json.Marshal(pm)
}

type JsonPrivateMember struct {
	Member
	PrivateKey string /* json:privatekey */
}

func (pm *PrivateMember) AsJson() *JsonPrivateMember {
	return &JsonPrivateMember{
		Member:     pm.Member,
		PrivateKey: pm.PrivateKey.Marshal(),
	}
}

func (pm *JsonPrivateMember) String() ([]byte, error) {
	return json.Marshal(pm)
}

func (pm *JsonPrivateMember) AsPrivateMember() (*PrivateMember, error) {
	pk, _, err := key.FromText(pm.PrivateKey, pm.Id)
	if err != nil {
		return nil, err
	}
	if pk == nil {
		return nil, errors.New("need to be an PK")
	}
	return &PrivateMember{
		Member:     pm.Member,
		PrivateKey: *pk,
	}, nil
}

type JsonPrivatePublicMember struct {
	Member
	PrivateKey *string
	PublicKey  *string
}

func JsToPublicMember(jspub *JsonPublicMember) (*PublicMember, error) {
	_, pb, err := key.FromText(jspub.PublicKey, jspub.Id)
	if err != nil {
		return nil, err
	}
	if pb == nil {
		return nil, errors.New("we need a publickey")
	}
	return &PublicMember{
		Member:    jspub.Member,
		PublicKey: *pb,
	}, nil
}

func JsToPrivateMember(jspriv *JsonPrivateMember) (*PrivateMember, error) {
	pk, _, err := key.FromText(jspriv.PrivateKey, jspriv.Id)
	if err != nil {
		return nil, err
	}
	if pk == nil {
		return nil, errors.New("we need a privatekey")
	}
	return &PrivateMember{
		Member:     jspriv.Member,
		PrivateKey: *pk,
	}, nil
}

func FromJson(str []byte) (*PrivateMember, *PublicMember, error) {
	// pk := string("")
	// pb := string("")
	jppm := JsonPrivatePublicMember{
		// PrivateKey: &pk,
		// PublicKey:  &pb,
	}
	err := json.Unmarshal(str, &jppm)
	if err != nil {
		return nil, nil, err
	}
	if jppm.PrivateKey != nil {
		privk, err := JsToPrivateMember(&JsonPrivateMember{
			Member:     jppm.Member,
			PrivateKey: *jppm.PrivateKey,
		})
		if err != nil {
			return nil, nil, err
		}
		return privk, nil, nil
	}
	if jppm.PublicKey != nil {
		pubk, err := JsToPublicMember(&JsonPublicMember{
			Member:    jppm.Member,
			PublicKey: *jppm.PublicKey,
		})
		if err != nil {
			return nil, nil, err
		}
		return nil, pubk, nil
	}
	return nil, nil, errors.New("No Pub or Priv Key")
}

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
// 	pbTxt, err := pubbase.MarshalText()
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
