import * as crypto from 'crypto'
import { CreateRandomKey, RawKey } from '../key/key'
import { Fault, Ok, Result, toBufferUint8, uint8Equal } from '../utils/render'

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

/*
// example 96-bit nonce
let nonce = Buffer.alloc(12, 0xff);

// example 256-bit key
let key = Buffer.alloc(32, 0x01);

// some associated data
let assocData = Buffer.alloc(16, 0xaa);

// some data to encrypt
let data = Buffer.alloc(64, 0xbb);

// construct the cipher
let cipher = crypto.createCipheriv('chacha20-poly1305', key, nonce, { authTagLength: 16 });

// add associated data to cipher
cipher.setAAD(assocData);

// encrypt the data which will return an encrypted Buffer
// that is of equal length to the overall input to the
// stream cipher
cipher.update(data);
// 25805b670d5834ecb8a018ea87b6ff864117762481880fc723690d0e2d0cfd08a43c144291eb2df148b0d6981b66ca101344ea27c7a0860c2e5f1a7eed1e70eb

// finalize cipher which allows us to calculate the MAC
cipher.final();

// obtain the 128-bit poly1305 MAC that includes our associated data
cipher.getAuthTag()
// 9ef622cec7a5719261031e9ca91049d4
*/

export function Checksum(sr:SealRequest) :SealRequest {
	const my = crypto.createHash('sha256')
	my.update(sr.Payload)
	my.update(sr.Key.As32Byte())
	sr.Checksum = toBufferUint8(my.digest())
	return sr
}

function nonce(aead: cipher.AEAD, n : Uint8Array) : Uint8Array {
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
	const my = crypto.createHash('sha256')
	my.update(open)
	my.update(key.As32Byte())
	return uint8Equal(toBufferUint8(my.digest()), csum)
}

export function SkipVerify(csum : Uint8Array, key : RawKey, open : Uint8Array) : boolean {
	return true
}

class xchacha20poly1305 {
	key = new Uint32Array(8)
}

// NewX returns a XChaCha20-Poly1305 AEAD that uses the given 256-bit key.
//
// XChaCha20-Poly1305 is a ChaCha20-Poly1305 variant that takes a longer nonce,
// suitable to be generated randomly without risk of collisions. It should be
// preferred when nonce uniqueness cannot be trivially ensured, or whenever
// nonces are randomly generated.
function newX(key: Uint8Array): Result<xchacha20poly1305> {
	const ret = new xchacha20poly1305()
	if (key.length != ret.key.length) {
		return Fault("chacha20poly1305: bad key length")
	}
	const buf = Buffer.from(key)
	ret.key[0] = buf.readInt32LE(0)
	ret.key[1] = buf.readInt32LE(4)
	ret.key[2] = buf.readInt32LE(8)
	ret.key[3] = buf.readInt32LE(12)
	ret.key[4] = buf.readInt32LE(16)
	ret.key[5] = buf.readInt32LE(20)
	ret.key[6] = buf.readInt32LE(24)
	ret.key[7] = buf.readInt32LE(28)
	return Ok(ret)
}

export function Open(key: RawKey, sc :SealedContainer, verify: VerifyFN): Result<OpenContainer> {
	const aead = newX(key.As32Byte())
	// aead, err := chacha20poly1305.NewX(key[:])
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
	const plain = Buffer.from(pp.Payload, 'base64')
	if (!plain) {
		return Fault("decoding error of Payload")
	}
	return Open(key, {
		Checksum: pp.Checksum,
		Payload:  plain,
	}, verify)
}
