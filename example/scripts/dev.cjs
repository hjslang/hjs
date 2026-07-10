const path = require("path");
const fs = require("fs");
const chokidar = require("chokidar");
const { execFileSync } = require("child_process");
const { rimrafSync } = require("rimraf");
const liveServer = require("live-server");

const COMPILER = "hjs"
const INDEX_PATH = "./index.html";
const SRC_DIR = "./src";
const DIST_DIR = "./dist";

function compile(inputPath) {
  // compile filename
  const stdout = execFileSync(COMPILER, [inputPath]);

  // write output
  const { dir, name } = path.parse(inputPath);
  const outputDir = path.join(DIST_DIR, dir);
  const outputPath = path.join(outputDir, name + ".js");
  fs.mkdirSync(outputDir, { recursive: true });
  fs.writeFileSync(outputPath, stdout, "utf-8");
}

function compileDir(dirPath) {
  const files = fs.readdirSync(dirPath);
  for (const file of files) {
    const filePath = path.join(dirPath, file);
    if (fs.lstatSync(filePath).isDirectory()) {
      compileDir(filePath);
      continue;
    }
    compile(filePath);
  }
}

function main() {
  rimrafSync(DIST_DIR);
  fs.mkdirSync(DIST_DIR);
  fs.copyFileSync(INDEX_PATH, path.join(DIST_DIR, INDEX_PATH));

  const w = chokidar.watch(SRC_DIR, {
    alwaysStat: true,
    ignored: (path, stats) => stats?.isFile() && !path.endsWith(".hjs"),
  });
  w.on("change", (inputPath) => compile(inputPath));
  w.on("add", (inputPath) => compile(inputPath));
  w.on("unlink", (inputPath) => {
    const { dir, name } = path.parse(inputPath);
    const outputPath = path.join(DIST_DIR, dir, name + ".js");
    if (fs.existsSync(outputPath)) fs.unlinkSync(outputPath);
  });

  liveServer.start({
    root: DIST_DIR,
    open: false,
  });
}

main();
