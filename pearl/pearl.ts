import * as uuid from 'uuid';
import { CreateShared } from '../asymmetric/asymmetric';
import { FromText, PrivateKey, PublicKey, RawKey } from '../key/key'
import { Fault, Ok, Result, uint8Equal } from '../utils/render'

export interface JsonStandardClaims {
aud?: string;
exp?: number;
jti?: string;
iat?: number;
iss?: string;
nbf?: number;
sub?: string;
}

export interface StandardClaims {
	Audience:  string
	ExpiresAt: number
	Id     :   string
	IssuedAt:  number
	Issuer:    string
	NotBefore: number
	Subject:   string
}

// PearlClaim is the enhanced structure to the jwt.StandardClaims
export interface PearlClaim extends StandardClaims {
	PayloadChecksum          :string // `json:"pcs"`
	EncryptedPayloadPassword :string // `json:"epp"`
}

// JWTokenPearlClaim the Claim encoded as a JWT-Token string
export type JWTokenPearlClaim = string

// CreatorOwners is the collection with Creators publickey and the JWT-Tokens
// sign and encrypted
export class CreatorOwners {
	constructor(
	public readonly Creator : PublicKey, // public-key
	public readonly Tokens:  JWTokenPearlClaim[]
	) {}

// AsJSON converts a CreatorOwners into the Json representation
    public AsJSON(): JSONCreatorOwners {
	return {
			Creator: this.Creator.Marshal(),
			Tokens:  this.Tokens
		}
	}
}

// JSONCreatorOwners is the json representation of the CreatorOwners
export interface JSONCreatorOwners {
	readonly Creator : string
	readonly Tokens  :JWTokenPearlClaim[]
}

// FromJSON converts from Json to CreatorOwners
function FromJSONCreatorOwners(jcp: JSONCreatorOwners): Result<CreatorOwners> {
	const ft = FromText(jcp.Creator, `SYN-${uuid.v4()}`)
	if (ft.isError) {
		return Fault(ft.error)
	}
	if (!ft.val.PublicKey) {
		return Fault("we need this as public key")
	}
	return Ok(new CreatorOwners(ft.val.PublicKey, jcp.Tokens))
}


// Pearl describes an Pearl which is part of an Necklace
export class Pearl {
	constructor(
	public FingerPrint: Uint8Array,
	public readonly Type        :string,
	public readonly Payload     : Uint8Array,
	public readonly Owners      : CreatorOwners
	) {}

// AsJSON converts a Pearl to a JSONPearl
 public AsJSON(): JSONPearl {
	return {
		FingerPrint: Buffer.from(this.FingerPrint).toString('base64'),
		Type:        this.Type,
		Payload:     Buffer.from(this.Payload).toString('base64'),
		Owners:      this.Owners.AsJSON(),
	}
}

}

// JSONPearl is the Json representation of a Pearl
export interface JSONPearl {
	FingerPrint :string
	Type        :string
	Payload     :string // base64
	Owners      :JSONCreatorOwners
}

interface argCalcFingerPrint {
	Signer  :PublicKey
	Payload :Uint8Array
	Type    :string
}

function calcFingerprint(cfp :argCalcFingerPrint) : Result<Uint8Array> {
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

function setFingerprint(cr :CloseRequestPearl, p :Pearl): Result<Pearl> {
	const sum = calcFingerprint({
		Signer:  cr.Owners.Signer.Public(),
		Payload: cr.Payload,
		Type:    cr.Type,
	})
	if (sum.isError) {
		return Fault(sum.error)
	}
	p.FingerPrint = sum.val
	return Ok(p)
}

function checkFingerprint(p :OpenPearl) : Result<OpenPearl> {
	const sum = calcFingerprint({
		Type:    p.Closed.Type,
		Signer:  p.Closed.Owners.Creator,
		Payload: p.Payload,
	})
	if (sum.isError) {
		return Fault(sum.error)
	}
	if (!uint8Equal(p.Closed.FingerPrint, sum.val)) {
		return Fault("checksum missmatch")
	}
	return Ok(p)
}

// FromJSON converts a json pearl to a Pearl
export function FromJSONPearl(jp: JSONPearl): Result<Pearl> {
	const payload = Buffer.from(jp.Payload, 'base64')
	if (!payload) {
		return Fault('decode error payload')
	}
	const fingerprint = Buffer.from(jp.FingerPrint, 'base64')
	if (!fingerprint) {
		return Fault('decode error fingerprint')
	}
	const owners = FromJSONCreatorOwners(jp.Owners)
	if (owners.isError) {
		return Fault(owners.error)
	}

	return Ok(new Pearl(
		fingerprint,
		       jp.Type,
		    payload,
		     owners.val
	))
}

export interface PearlOwner {
	readonly Signer :PrivateKey
	readonly Owners :PublicKey[]
}

export interface CloseRequestPearl {
	readonly Type    :string
	readonly Payload : Uint8Array
	readonly Owners  : PearlOwner
}

export interface CloseContainer {
readonly 	Checksum   : Uint8Array
	readonly PayloadKey :RawKey
}

function createJWTPearlClaim(signer: PrivateKey, pub :PublicKey, cl :CloseContainer) : Result<JWTokenPearlClaim> {
	const privPubKey = CreateShared(signer.Raw, pub.Raw)
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

}
