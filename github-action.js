const core = require("@actions/core");
const fs = require("fs").promises;
const path = require("path");
const http = require("http");
const https = require("https");
const tar = require("tar");
const { ok } = require("assert");

function download(url, cnt) {
  return new Promise((rs, rj) => {
    if (cnt > 8) {
      rj(Error(`loop detected ${url}`));
      return;
    }
    let client = http;
    if (url.toString().indexOf("https") === 0) {
      client = https;
    }
    client
      .request(url, (response) => {
	if (300<=response.statusCode && response.statusCode < 400) {
		const url = response.headers['location'];
		console.log(`url ==> ${url}`);
		rs(download(url, cnt + 1));
		return;
	}
	if (!(200<=response.statusCode && response.statusCode < 300)) {
		rj(`fetch error ${url}:${response.statusCode}`);
	}
        response.on("error", (e) => {
          rj(e);
        });
        const data = [];
        response.on("data", function (chunk) {
          data.push(Buffer.from(chunk));
        });
        response.on("end", function () {
          rs(Buffer.concat(data));
        });
      })
      .end();
  });
}

function getTempDirectory() {
  const tempDirectory = process.env["RUNNER_TEMP"] || ".";
  ok(tempDirectory, "Expected RUNNER_TEMP to be defined");
  return tempDirectory;
}

async function main() {
  try {
    const params = {
	    version: core.getInput("version"),
	    url: core.getInput("url"),
	    filename: core.getInput("filename"),
	    os: core.getInput("os"), // no default
	    suffix: core.getInput("suffix"),
	    cpu: core.getInput("cpu") // no default
    };
    if (params.os === "") {
	    switch (process.platform) {
	      case "darwin":
		params.os = "Darwin";
		break;
	      case "linux":
		params.os = "Linux";
	    }
     }
     if (params.cpu === "") {
	switch (process.arch) {
	  case "i386":
	    params.cpu = "i386";
	    break;
	  case "x64":
	    params.cpu = "x86_64";
	    break;
	  case "arm":
	    params.cpu = "arm7";
	    break;
	  case "arm64":
	    params.cpu = "arm64";
	    break;
	}
    }
    core.info(`Fetch neckless choosen cpu [${params.cpu}][${process.arch}]`)
    const plainVersion = params.version.replace(/^v/, '');
    const necklessUrl = `${params.url}/${params.version}/${params.filename}_${plainVersion}_${params.os}_${params.cpu}${params.suffix}`;
    core.info(`Fetch neckless from:[${necklessUrl}]`)
    const necklessBin = await download(necklessUrl, 0);
    const necklessBinDir = path.join(getTempDirectory(), "neckless-bin");
    // const dir = path.join(getTempDirectory()); // , `neckless${params.suffix}`);
    await fs.mkdir(necklessBinDir, {
      recursive: true,
      mode: 0o755,
    });
    const necklessFnameTar = path.join(getTempDirectory(), `neckless${params.suffix}`);
    await fs.writeFile(necklessFnameTar, necklessBin);
    await tar.x({
	    file: necklessFnameTar,
	    cwd: necklessBinDir
    });
    // await fs.chmod(necklessFname, 0o755);
    core.exportVariable("NECKLESS_URL", necklessUrl);
    const necklessFname = path.join(necklessBinDir, "neckless");
    core.exportVariable("NECKLESS_FNAME", necklessFname);
    core.addPath(necklessBinDir);
    await fs.unlink(necklessFnameTar);
    core.info(`Installed neckless into:[${necklessFname}]`)
  } catch (e) {
    core.setFailed(e);
  }
}

main();
