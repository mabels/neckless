import * as crypto from 'crypto'
import { CreateRandomKey, RawKey } from '../key/key'
import { Fault, Ok, Result } from '../utils/render'

export class KeyAndNonce {
	constructor(
	public readonly Key: Uint8Array,
	public readonly Nouce: Uint8Array
	){}
	public AsJSON() : JsonKeyAndNonce {
		return new JsonKeyAndNonce(
			 Buffer.from(this.Key).toString('base64'),
			Buffer.from(this.Nouce).toString('base64'),
		)
	}
}

export class JsonKeyAndNonce {
	constructor(
		public readonly Key:   string,
		public readonly Nouce: string
	) {}
	public AsKeyAndNonce(): Result<KeyAndNonce> {
	const key = Buffer.from(this.Key, 'base64')
	if (!key) {
		return Fault('base64 decode faild for Key')
	}
	const nonce = Buffer.from(this.Nouce, 'base64')
	if (!nonce) {
		return Fault('base64 decode faild for Nouce')
	}
	return Ok(new KeyAndNonce(
		  key,
		nonce,
	))
}
}

export interface SealedContainer {
	Checksum : Uint8Array
	Payload  : Uint8Array
}

export interface SealRequest {
	Key: RawKey,
	Payload  : Uint8Array
	Checksum : Uint8Array
}

function Checksum(sr:SealRequest) :SealRequest {
	const my = crypto.createHash('sha256')
	my.update(sr.Payload)
	my.update(sr.Key.As32Byte())
	sr.Checksum = toBufferUint8(my.digest())
	// fmt.Printf("Seal:%x=>%x\n", sr.Checksum, data)
	return sr
}

function nonce(aead cipher.AEAD, n : Uint8Array) : Uint8Array {
	csum := make([]byte, aead.NonceSize())
	if len(n) < aead.NonceSize() {
		copy(csum, n)
	} else {
		copy(csum, n[:aead.NonceSize()])
	}
	return csum
}

function Seal(sr:SealRequest): Result<SealedContainer> {
	aead, err := chacha20poly1305.NewX(sr.Key[:])
	if err != nil {
		return nil, err
	}
	nonce := nonce(aead, sr.Checksum)
	sealed := aead.Seal(nil, nonce, sr.Payload, nil)
	// fmt.Printf("Seal:%d:%x=>%x:%x\n", len(nonce), nonce, sr.Key, sealed)
	return &SealedContainer{
		Checksum: sr.Checksum,
		Payload:  sealed,
	}, nil
}

export interface OpenContainer {
	Checksum : Uint8Array
	Payload  : Uint8Array
}

export type VerifyFN = (x: Uint8Array, raw: RawKey, u: Uint8Array) => boolean;

export function Verify(csum: Uint8Array, key: RawKey, open: Uint8Array): boolean  {
	data := append(*open, key[:]...)
	// fmt.Printf("Verify:%x=>%x", csum, data)
	tmp := sha256.Sum256(data)
	return bytes.Equal(tmp[:], csum)
}

export function SkipVerify(csum : Uint8Array, key : RawKey, open : Uint8Array) : boolean {
	return true
}

export function Open(key: RawKey, sc :SealedContainer, verify: VerifyFN): Result<OpenContainer> {
	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, err
	}
	nonce := nonce(aead, sc.Checksum)
	// fmt.Printf("Open:%d:%x=>%x:%x\n", len(nonce), nonce, key, sc.Payload)
	open, err := aead.Open(nil, nonce, sc.Payload, nil)
	// fmt.Printf("Open:%x:%x\n", key[:], sc.Checksum[:aead.NonceSize()])
	if err != nil {
		return nil, err
	}
	if !verify(sc.Checksum, key, &open) {
		return nil, errors.New("checksum error")
	}
	return &OpenContainer{
		Checksum: sc.Checksum,
		Payload:  open,
	}, nil
}

export interface Base64SealContainer extends KeyAndNonce {
	Checksum: Uint8Array
	Payload:  string
}


export function OpenBase64(key: RawKey, pp: Base64SealContainer, verify: VerifyFN): Result<OpenContainer> {
	plain, err := base64.RawStdEncoding.DecodeString(pp.Payload)
	if err != nil {
		return nil, err
	}
	return Open(key, &SealedContainer{
		Checksum: pp.Checksum,
		Payload:  plain,
	}, verify)
}
