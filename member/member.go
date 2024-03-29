package member

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/mabels/neckless/key"
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
	Email      string
	ValidUntil *time.Time
	Updated    *time.Time `json:"updated"`
	Created    *time.Time `json:"created"`
}

type MemberBase struct {
	Type       MemberType `json:"type"`
	Name       string     `json:"name"`
	Email      string     `json:"email,omitempty"`
	Device     string     `json:"device,omitempty"`
	ValidUntil time.Time  `json:"validUntil"`
	Updated    time.Time  `json:"updated"`
	Created    time.Time  `json:"created"`
}

type Member struct {
	MemberBase
	Id string `json:"id"`
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
	if len(m.Email) != 0 {
		ret.Email = m.Email
	}
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
	Clazz string `json:"clazz"`
	Member
	PublicKey string `json:"publicKey"`
}

func JsonPublicMemberValueBy(p1, p2 *JsonPublicMember) bool {
	return strings.Compare(p1.Member.Id, p2.Member.Id) < 0
}

type JsonPublicMemberSorter struct {
	Values [](*JsonPublicMember)
	By     func(p1, p2 *JsonPublicMember) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *JsonPublicMemberSorter) Len() int {
	return len(s.Values)
}

// Swap is part of sort.Interface.
func (s *JsonPublicMemberSorter) Swap(i, j int) {
	s.Values[i], s.Values[j] = s.Values[j], s.Values[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *JsonPublicMemberSorter) Less(i, j int) bool {
	return JsonPublicMemberValueBy(s.Values[i], s.Values[j])
}

// func (mb *Member) AsJSON() *JsonMember {
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

func (pm *PublicMember) AsJSON() *JsonPublicMember {
	return &JsonPublicMember{
		Clazz:     "JsonPublicMember",
		Member:    pm.Member,
		PublicKey: pm.PublicKey.Marshal(),
	}
}

func (pm *JsonPublicMember) String() ([]byte, error) {
	return json.Marshal(pm)
}

type JsonPrivateMember struct {
	Clazz string `json:"clazz"`
	Member
	PrivateKey string `json:"privatekey"`
}

func (pm *PrivateMember) AsJSON() *JsonPrivateMember {
	return &JsonPrivateMember{
		Clazz:      "JsonPrivateMember",
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

func ToJsonPrivateMember(pkms ...*PrivateMember) []*JsonPrivateMember {
	ret := make([]*JsonPrivateMember, len(pkms))
	for i := range pkms {
		ret[i] = pkms[i].AsJSON()
	}
	return ret
}

func ToJsonPublicMember(pkms ...*PrivateMember) []*JsonPublicMember {
	ret := make([]*JsonPublicMember, len(pkms))
	for i := range pkms {
		ret[i] = pkms[i].Public().AsJSON()
	}
	return ret
}

type JsonPrivatePublicMember struct {
	Member
	PrivateKey *string
	PublicKey  *string
}

func Matcher(args ...string) func(p *PrivateMember) bool {
	return func(p *PrivateMember) bool {
		for i := range args {
			if strings.Contains(p.Id, args[i]) ||
				strings.Contains(p.Name, args[i]) ||
				strings.Contains(p.Email, args[i]) {
				return true
			}
		}
		return len(args) == 0
	}
}

func Filter(pkms []*PrivateMember, filters ...func(*PrivateMember) bool) []*PrivateMember {
	ret := []*PrivateMember{}
	/*
		mids := map[string](struct{}){}
		for i := range ids {
			mids[ids[i]] = struct{}{}
		}
	*/
	filter := func(*PrivateMember) bool { return true }
	if len(filters) > 0 {
		filter = filters[0]
	}

	for i := range pkms {
		pkm := pkms[i]
		if filter(pkm) {
			// _, found := mids[pkm.Id]
			// fmt.Println(pkm.Id, ids)
			// if len(mids) == 0 || found {
			ret = append(ret, pkm)
		}
	}
	return ret
}

func FilterByType(pkms []*PrivateMember, typs ...MemberType) []*PrivateMember {
	ret := []*PrivateMember{}
	mtyps := map[MemberType](struct{}){}
	for i := range typs {
		mtyps[typs[i]] = struct{}{}
	}
	for i := range pkms {
		pkm := pkms[i]
		_, found := mtyps[pkm.Type]
		if len(mtyps) == 0 || found {
			ret = append(ret, pkm)
		}
	}
	return ret
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

func FromJSON(str []byte) (*PrivateMember, *PublicMember, error) {
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

func ToPrivateKeys(pkms []*PrivateMember) []*key.PrivateKey {
	out := make([]*key.PrivateKey, len(pkms))
	for i := range pkms {
		out[i] = &pkms[i].PrivateKey
	}
	return out
}

func ToPublicKeys(pkms []*PublicMember) []*key.PublicKey {
	out := make([]*key.PublicKey, len(pkms))
	for i := range pkms {
		out[i] = &pkms[i].PublicKey
	}
	return out
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
