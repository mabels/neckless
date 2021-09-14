const core = require("@actions/core");
const uuid = require("uuid");
const { spawn } = require("child_process");
const fs = require("fs").promises;
const path = require("path");
const http = require("http");
const https = require("https");
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
    const version = core.getInput("version");
    const url = core.getInput("url");
    const filename = core.getInput("filename");
    let os = core.getInput("os");
    let cpu = core.getInput("cpu");
    switch (process.platform) {
      case "darwin":
        os = "mac";
        cpu = undefined;
        break;
      case "linux":
        os = "linux";
        switch (process.arch) {
          case "x64":
            cpu = "amd64";
            break;
          case "arm":
            cpu = "arm-v7";
            break;
          case "arm64":
            cpu = "arm-v8";
            break;
        }
        break;
    }
    let necklessUrl = `${url}/${version}/${filename}-${os}`;
    if (cpu) {
      necklessUrl = `${necklessUrl}-${cpu}`;
    }
    core.info(`Fetch neckless from:[${necklessUrl}]`)
    const necklessBin = await download(necklessUrl, 0);
    const dir = path.join(getTempDirectory(), "neckless-bin");
    await fs.mkdir(dir, {
      recursive: true,
      mode: 0o755,
    });
    const necklessFname = path.join(dir, "neckless");
    await fs.writeFile(necklessFname, necklessBin);
    await fs.chmod(necklessFname, 0o755);
    core.exportVariable("NECKLESS_URL", necklessUrl);
    core.exportVariable("NECKLESS_FNAME", necklessFname);
    core.addPath(dir);
    core.info(`Installed neckless into:[${necklessFname}]`)
  } catch (e) {
    core.setFailed(e);
  }
}

main();
