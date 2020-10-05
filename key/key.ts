import * as uuid from 'uuid';
import { Fault, Ok, Result } from '../utils/render';
// package key

// import (
// 	"bytes"
// 	"crypto/rand"
// 	"crypto/subtle"
// 	"encoding/base64"
// 	"errors"
// 	"fmt"
// 	"strings"

// 	"github.com/google/uuid"
// 	"golang.org/x/crypto/curve25519"
// )

// KeySize is Default curve25519 key size
const RawKeyLen = 32;

export class RawKey {
	public readonly bytes = new Uint8Array(RawKeyLen)
	constructor(my: Uint8Array | Buffer) {
		const ret = this.bytes
		if (my.length < ret.length) {
			for (let i = 0; i < my.length; ++i) {
				ret[i] = my[i]
			}
		} else {
			for (let i = 0; i < ret.length; ++i) {
				ret[i] = my[i]
			}
		}
	}
	public As32Byte(): Uint8Array {
		return this.bytes
	}
	public Marshal(): string {
		const out = Buffer.from(this.bytes);
		return out.toString('base64')
	}

}

export enum KeyStyle {
	Private = "Private",
	Public  = "Public"
}

export interface JsonKeyType {
	Id:    string
	Style: KeyStyle
	Raw:   string
}

export function ToKeyType(jsk: JsonKeyType): KeyType | undefined {
	const raw = Buffer.from(jsk.Raw, 'base64')
	return new KeyType(jsk.Id, jsk.Style, raw);
}


export class KeyType {
	public readonly Raw:   RawKey
	constructor(
	public readonly Id:    string,
	public readonly Style: KeyStyle,
	key: Uint8Array|Buffer
	) {
		this.Raw = new RawKey(key)
	}

	public AsJSON(): JsonKeyType {
		let buff = Buffer.from(this.Raw.As32Byte());
		return {
			Id:    this.Id,
			Style: this.Style,
			Raw:   buff.toString('base64')
		}
	}
}



export class PublicKey extends KeyType {
	readonly Style = KeyStyle.Public
	public Marshal(): string {
		return this.Raw.Marshal()
	}

}

export class PrivateKey extends KeyType {
	readonly Style = KeyStyle.Private
	// MarshalText create a String of the Private Key
	public Marshal(): string {
		return `privkey:${this.Raw.Marshal()}`
	}
	// Public computes the public key matching this curve25519 secret key.
	public Public(): PublicKey {
	if (isZero(this.Raw)) {
		throw Error("Tried to generate emptyPrivateKey.Public()")
	}
	const pub = new Uint8Array(RawKeyLen)
	throw Error("not ready")
	// curve25519.ScalarBaseMult(&pub, k.Key.Raw.As32Byte())
	return MakePublicKey(pub, this.Id)
}

}

export function MakePublicKey(key: Uint8Array, id: string): PublicKey {
	return new PublicKey( id, KeyStyle.Public, key)
}

export function MakePrivateKey(key :Uint8Array, id?: string) :PrivateKey {
	return new PrivateKey(id || uuid.v4(), KeyStyle.Private, key)
}

// CreateRandomKey generates a new random key.
export function CreateRandomKey(): Result<RawKey> {
	return Fault("missing implementation")
	// return new RawKey()
}

// NewPrivateKey generates a new curve25519 secret key.
// It conforms to the format described on https://cr.yp.to/ecdh.html.
export function NewPrivateKey(pk?: PrivateKey, id?: string): PrivateKey {
	let k: Uint8Array
	if (pk && pk.Raw.bytes[0] &&
		pk.Style == KeyStyle.Private &&
		pk.Raw.bytes.length === RawKeyLen) {
		k = pk.Raw.bytes
	} else {
		k = CreateRandomKey().val.bytes
	}
	k[0] &= 248
	k[31] = (k[31] & 127) | 64
	return MakePrivateKey(k, id)
}

export function isZero(k: RawKey): boolean {
	return !k.bytes.find(i => i !== 0)
}



export function fromText(pkstr: string): RawKey {
	return new RawKey(Buffer.from(pkstr, 'base64'))
}

export interface FromTextResult {
	PrivateKey?: PrivateKey,
	PublicKey?: PublicKey,
}
export function FromText(pkstr: string, id: string): Result<FromTextResult> {
	if (pkstr.startsWith("privkey:")) {
		const keyBytes = fromText(pkstr.substr("privkey:".length))
		return Ok({
			PrivateKey: MakePrivateKey(keyBytes.bytes, id)
		})
	}
	const keyBytes = fromText(pkstr)
	return Ok({
		PublicKey: MakePublicKey(keyBytes.bytes, id)
	})
}
