import * as path from "path";
import * as fs from "fs";
import * as mkdirp from 'mkdirp';
import { JsonPrivateMember, MakePrivateMember, Member, MemberArg, PrivateMember } from '../member/member';
import { Result } from '../utils/render';

export interface CreateArg {
	Member: MemberArg
	DryRun: boolean    // if dryrun don't write
	Fname  :string //
}

export interface CasketAttribute {
 	CasketFname: string
 	Created: Date
 	Updated: Date
}

export class Casket implements CasketAttribute {
	constructor(
		public readonly CasketFname: string,
		public readonly Created: Date,
		public readonly Updated: Date,
		public readonly Members: Map<string, PrivateMember>
	) {}
	public AsJSON() : JsonCasket {
		return {
			CasketFname: this.CasketFname,
			Created: this.Created,
			Updated: this.Updated,
			Members:  Object.entries(this.Members).reduce((r, [k, v]) => {
				r[k] = v.AsJSON()
				return r
			}, {} as Record<string, JsonPrivateMember>)
		}
	}

public AsPrivateMembers(): PrivateMember[] {
	return Array.from(this.Members.values())
}
}

export interface JsonCasket extends CasketAttribute {
	Members: Record<string, JsonPrivateMember> // `json:"members"`
}


export function getcasketFilename(fname?: string): Promise<string> {
	if (!fname) {
		fname = path.join(process.env.HOME, ".neckless/casket");
	}
	await mkdirp(path.dirname(fname), {
		mode: 0o700
	});
	return Promise.resolve(fname);
}

export function readcasket(fname: string): Promise<Casket> {
	const now = new Date()

	const dat = await fs.promises.readFile(fname)
	const jsonCasket: JsonCasket = JSON.parse(dat.toString())
	const members = new Map<string, PrivateMember>()
	Object.values(jsonCasket.Members).forEach(jspm => {
		const pm = jspm.AsPrivateMember()
		if err != nil {
			return nil, err
		}
		// fmt.Printf("Ls:Key:%s\n", k)
		members[k] = pm
	})
	// fmt.Printf("Ls:%d\n", len(members))
	return new Casket{
		CasketAttribute: jsonCasket.CasketAttribute,
		Members:         members,
	}, err
}


function writecasket(casket: Casket): Promise<never> {
	const jsstr = JSON.stringify(casket.AsJSON(), null, 2)
	const tmp = path.join(path.dirname(casket.CasketFname),
		`.${process.pid}.${path.basename(casket.CasketFname)}`)
	await fs.promises.writeFile(tmp, jsstr, { mode: 0o600 })
	await fs.promises.rename(tmp, casket.CasketFname)
	return
}

// UseCase Write the PrivateKey in den casket ~/.neckless/casket
// neckless casket create --name <name> [--device <name>] [--person|--device] [--file=~/.crazybee/casket]
export interface CreateResult {
	readonly casket: Casket;
	readonly privateMember: PrivateMember;
}
export function Create(ca: CreateArg): Result<CreateResult> {
	const pk = MakePrivateMember(ca.Member)

	if err != nil {
		return nil, nil, err
	}
	var casket *Casket
	if ca.Fname == nil || len(*ca.Fname) == 0 {
		casket, err = Ls()
	} else {
		casket, err = Ls(*ca.Fname)
	}
	if err != nil {
		return nil, nil, err
	}
	casket.Members[pk.Id] = pk
	casket.Updated = time.Now()
	if !ca.DryRun {
		err = writecasket(casket)
		if err != nil {
			return nil, nil, err
		}
	}
	return casket, pk, nil
}

// UseCase List casket
// neckless casket ls
export function Ls(fname?: string): Promise<Casket> {
	fname = await getcasketFilename(fname)
	return readcasket(fname)
}


export interface RmArg {
	Ids    :string[]
	DryRun :boolean    // if dryrun don't write
	Fname?  :string //
}

// UseCase Delete Key from casket
// neckless casket rm <id>
export function Rm(rmarg: RmArg): (Casket, []*member.PrivateMember, error) {
	var ks *Casket
	var err error
	if rmarg.Fname != nil {
		ks, err = Ls(*rmarg.Fname)
	} else {
		ks, err = Ls()
	}
	if err != nil {
		return nil, nil, err
	}
	out := []*member.PrivateMember{}
	for i := range rmarg.Ids {
		id := rmarg.Ids[i]
		pk, ok := ks.Members[id]
		if ok {
			delete(ks.Members, id)
			out = append(out, pk)
		}
	}
	if !rmarg.DryRun {
		if err = writecasket(ks); err != nil {
			return nil, nil, err
		}
	}
	return ks, out, nil
}
