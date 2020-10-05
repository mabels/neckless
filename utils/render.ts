export interface Result<T> {
    val: T
    error?: string
    isOk: boolean
    isError: boolean
}
export function Ok<T>(t: T): Result<T> {
    return {
        val: t,
        isOk: true,
        isError: false
    }
}

export function Fault<T>(val: string): Result<T> {
    return {
        val: undefined,
        error: val,
        isOk: false,
        isError: true
    }
}

export function uint8Equal(u1: Uint8Array, u2: Uint8Array) {
    if (u1.length != u2.length) { return false }
    for (let i = 0; i < u1.length; ++i) {
        if (u1[i] != u2[i]) {
            return false
        }
    }
    return true
}

export function toBufferUint8(b: Buffer): Uint8Array {
    const ret = new Uint8Array(b.length)
    for (let i = 0; i < b.length; i++) {
        ret[i] = b[i]
    }
    return ret;
}