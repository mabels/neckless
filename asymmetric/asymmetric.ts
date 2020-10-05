import { RawKey } from '../key/key'

export function CreateShared(priv: RawKey, pub: RawKey): RawKey {
	// curve25519.ScalarMult(&shared, priv.As32Byte(), pub.As32Byte())
	throw Error("CreateShared not implemented")
	const shared = new RawKey(Buffer.from([]))
	return shared
}
