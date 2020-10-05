import { FromText, NewPrivateKey, PrivateKey, PublicKey } from "../key/key";
import { Fault, Ok, Result } from "../utils/render";

export enum MemberType {
  Device = "Device",
  Person = "Person",
}

export interface MemberArg {
  Id: string;
  Type: MemberType;
  Name: string;
  Device: string;
  Email: string;
  ValidUntil: Date;
  Updated: Date;
  Created: Date;
}

export interface MemberBase {
  type: MemberType;
  name: string;
  email: string;
  device: string;
  validUntil: Date;
  updated: Date;
  created: Date;
}

export interface Member extends MemberBase {
  id: string;
}

function extractMember(o: Member): Member {
  return {
    id: o.id,
    type: o.type,
    name: o.name,
    email: o.email,
    device: o.device,
    validUntil: o.validUntil,
    updated: o.updated,
    created: o.created,
  };
}

// type JsonMember struct {
// 	MemberBase
// 	Id string `json:"id"`
// }

export interface PrivateMemberArg {
  Member: MemberArg;
  PrivateKey: PrivateKey;
}

export class PrivateMember {
  constructor(
    public readonly PrivateKey: PrivateKey,
    public readonly Member: Member
  ) {}
  public Public(): PublicMember {
    return new PublicMember(this.PrivateKey.Public(), this.Member);
  }
  public AsJSON(): JsonPrivateMember {
    return {
      ...this.Member,
      clazz: "JsonPrivateMember",
      PrivateKey: this.PrivateKey.Marshal(),
    };
  }
}

export class PublicMember {
  constructor(
    public readonly PublicKey: PublicKey,
    public readonly Member: Member
  ) {}

  public AsJSON(): JsonPublicMember {
    return {
      ...this.Member,
      clazz: "JsonPublicMember",
      PublicKey: this.PublicKey.Marshal(),
    };
  }
}

export function NewMember(m: Partial<MemberArg>): Result<Member> {
  const ret: Partial<Member> = {};
  ret.id = m.Id;
  // if len(m.Id) == 0 {
  // ret.Id = uuid.New().String()
  // } else {
  // ret.Id = m.Id
  // }

  ret.type = m.Type;
  if (m.Name) {
    return Fault("Require name");
  }
  ret.name = m.Name;
  if (m.Email) {
    ret.email = m.Email;
  }
  if (m.Device) {
    ret.device = m.Device;
  }
  const now = new Date();
  if (m.ValidUntil) {
    ret.validUntil = m.ValidUntil;
  } else {
    const plus5 = new Date();
    plus5.setFullYear(now.getFullYear() + 5);
    ret.validUntil = plus5;
  }
  if (m.Updated) {
    ret.updated = m.Updated;
  } else {
    ret.updated = now;
  }
  if (m.Created) {
    ret.created = m.Created;
  } else {
    ret.created = now;
  }
  return Ok(ret as Member);
}

export function MakePrivateMember(pm: PrivateMemberArg): Result<PrivateMember> {
  const m = NewMember(pm.Member);
  if (m.isError) {
    return Fault(m.error);
  }
  const pk = NewPrivateKey(pm.PrivateKey, m.val.id);
  return Ok(new PrivateMember(pk, m.val));
}

export function MakePublicMember(pm: PublicMember): Result<PublicMember> {
  return Ok(new PublicMember(pm.PublicKey, pm.Member));
}

export interface JsonPublicMember extends Member {
  clazz: string; // `json:"clazz"`
  PublicKey: string; // `json:"publicKey"`
}

export function JsonPublicMemberValueBy(
  p1: JsonPublicMember,
  p2: JsonPublicMember
): boolean {
  return p1.id < p2.id;
}

export interface JsonPrivateMember extends Member {
  clazz: string; // `json:"clazz"`
  PrivateKey: string; // `json:"privatekey"`
}

export function AsPrivateMember(jpm: JsonPrivateMember): Result<PrivateMember> {
  const my = FromText(jpm.PrivateKey, jpm.id);
  if (my.isError) {
    return Fault(my.error);
  }
  if (!my.val.PrivateKey) {
    return Fault("need to be an PK");
  }
  return Ok(new PrivateMember(my.val.PrivateKey, extractMember(jpm)));
}

export function ToJsonPrivateMember(
  pkms: PrivateMember[]
): JsonPrivateMember[] {
  return pkms.map((i) => i.AsJSON());
}

export function ToJsonPublicMember(pkms: PrivateMember[]): JsonPublicMember[] {
  return pkms.map((i) => i.Public().AsJSON());
}

export interface JsonPrivatePublicMember extends Member {
  PrivateKey: string;
  PublicKey: string;
}

export function Matcher(...args: string[]): (p: PrivateMember) => boolean {
  return (p: PrivateMember): boolean => {
    return (
      args.length == 0 ||
      !!args.find((i) => {
        if (
          p.Member.id.includes(i) ||
          p.Member.name.includes(i) ||
          p.Member.email.includes(i)
        ) {
          return true;
        }
      })
    );
  };
}

export function Filter(
  pkms: PrivateMember[],
  ...filters: ((pm: PrivateMember) => boolean)[]
): PrivateMember[] {
  let filter = (_: PrivateMember) => true;
  if (filters.length > 0) {
    filter = filters[0];
  }
  return pkms.filter((i) => filter(i));
}

export function FilterByType(
  pkms: PrivateMember[],
  ...typs: MemberType[]
): PrivateMember[] {
  const mtyps = typs.reduce((r, i) => {
    r.add(i);
    return r;
  }, new Set<MemberType>());
  return pkms.filter((pkm) => {
    return mtyps.size == 0 || mtyps.has(pkm.Member.type);
  });
}

export function JsToPublicMember(
  jspub: JsonPublicMember
): Result<PublicMember> {
  const m = FromText(jspub.PublicKey, jspub.id);
  if (m.isError) {
    return Fault(m.error);
  }
  if (!m.val.PublicKey) {
    return Fault("we need a publickey");
  }
  return Ok(new PublicMember(m.val.PublicKey, extractMember(jspub)));
}

export function JsToPrivateMember(
  jspriv: JsonPrivateMember
): Result<PrivateMember> {
  const m = FromText(jspriv.PrivateKey, jspriv.id);
  if (m.isError) {
    return Fault(m.error);
  }
  if (!m.val.PrivateKey) {
    return Fault("we need a privatekey");
  }
  return Ok(new PrivateMember(m.val.PrivateKey, extractMember(jspriv)));
}

export interface PrivPub {
  private?: PrivateMember;
  public?: PublicMember;
}

export function FromJSON(str: Buffer | string): Result<PrivPub> {
  const jppm = JSON.parse(str.toString());
  if (jppm.PrivateKey) {
    const privk = JsToPrivateMember({
      ...jppm.Member,
      PrivateKey: jppm.PrivateKey,
    });
    if (privk.isError) {
      return Fault(privk.error);
    }
    return Ok({ private: privk.val });
  }
  if (jppm.PublicKey) {
    const pubk = JsToPublicMember({
      ...jppm.Member,
      PublicKey: jppm.PublicKey,
    });
    if (pubk.isError) {
      return Fault(pubk.error);
    }
    return Ok({ public: pubk.val });
  }
  return Fault("No Pub or Priv Key");
}

export function ToPrivateKeys(pkms: PrivateMember[]): PrivateKey[] {
  return pkms.map((i) => i.PrivateKey);
}

export function ToPublicKeys(pkms: PublicMember[]): PublicKey[] {
  return pkms.map((i) => i.PublicKey);
}
