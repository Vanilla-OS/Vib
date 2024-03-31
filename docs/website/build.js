const fs = require("fs");
const path = require("path");
const Handlebars = require("handlebars");
const crypto = require("crypto");

function calculateFileSignature(filePath) {
  const fileContent = fs.readFileSync(filePath);
  const hash = crypto.createHash("md5").update(fileContent).digest("hex");
  return hash;
}

class Builder {
  constructor(inputDir, outputDir) {
    this.inputDir = inputDir;
    this.outputDir = outputDir;
    this.fileSignatures = {};
    this.loadPartials();
  }

  loadPartials() {
    const partialsDir = path.join(this.inputDir, "partials");
    fs.readdirSync(partialsDir)
      .filter((file) => file.endsWith(".hbs"))
      .forEach((file) => {
        const partialName = path.basename(file, ".hbs");
        const partialContent = fs.readFileSync(
          path.join(partialsDir, file),
          "utf8"
        );
        Handlebars.registerPartial(partialName, partialContent);
      });
  }

  ensureDirectoryExists(directory) {
    if (!fs.existsSync(directory)) {
      fs.mkdirSync(directory, { recursive: true });
    }
  }

  cleanDirectory(directory, shouldClean) {
    if (shouldClean && fs.existsSync(directory)) {
      fs.rmdirSync(directory, { recursive: true });
    }
  }

  saveFileSignatures(filePath) {
    fs.writeFileSync(filePath, JSON.stringify(this.fileSignatures, null, 2));
  }

  loadFileSignatures(filePath) {
    if (fs.existsSync(filePath)) {
      const content = fs.readFileSync(filePath, "utf-8");
      this.fileSignatures = JSON.parse(content);
    }
  }

  updateFileSignature(filePath) {
    const fileContent = fs.readFileSync(filePath);
    const hash = crypto.createHash("md5").update(fileContent).digest("hex");
    this.fileSignatures[filePath] = hash;
  }

  shouldBuild(filePath) {
    if (!this.fileSignatures.hasOwnProperty(filePath)) {
      return true;
    }
    const currentSignature = calculateFileSignature(filePath);
    return currentSignature !== this.fileSignatures[filePath];
  }

  buildHandlebarsFiles() {
    const files = fs
      .readdirSync(this.inputDir)
      .filter((file) => file.endsWith(".hbs") && !file.includes("partials"));
    let hasPartialChanges = false;

    console.log("Building Handlebars files...");

    files.forEach((file) => {
      const filePath = path.join(this.inputDir, file);
      if (this.shouldBuild(filePath)) {
        hasPartialChanges = true;
        console.log(`├── ${file} has changed, rebuilding...`);
        this.compileAndWriteFile(filePath);
        this.updateFileSignature(filePath);
      }
    });

    if (hasPartialChanges) {
      console.log("Rebuilt due to changes.");
    } else {
      console.log("No changes detected.");
    }
  }

  compileAndWriteFile(filePath) {
    const template = fs.readFileSync(filePath, "utf8");
    const compiledTemplate = Handlebars.compile(template);
    const html = compiledTemplate({});
    const outputFilePath = path.join(
      this.outputDir,
      path.basename(filePath, ".hbs") + ".html"
    );
    fs.writeFileSync(outputFilePath, html);
  }

  copyAssets() {
    const assetsDir = path.join(this.inputDir, "assets");
    const outputAssetsDir = path.join(this.outputDir, "assets");

    this.copyDirectoryRecursive(assetsDir, outputAssetsDir);
    console.log("Copied assets.");
  }

  copyDirectoryRecursive(source, target) {
    this.ensureDirectoryExists(target);
    const items = fs.readdirSync(source, { withFileTypes: true });

    items.forEach((item) => {
      const sourcePath = path.join(source, item.name);
      const targetPath = path.join(target, item.name);

      if (item.isDirectory()) {
        this.copyDirectoryRecursive(sourcePath, targetPath);
      } else {
        fs.copyFileSync(sourcePath, targetPath);
      }
    });
  }

  buildAll() {
    const startTime = Date.now();
    const shouldClean = process.argv.includes("--no-cache");
    this.loadFileSignatures(path.join(this.outputDir, "signatures.json"));
    this.ensureDirectoryExists(this.outputDir);
    this.cleanDirectory(this.outputDir, shouldClean);
    this.buildHandlebarsFiles();
    this.copyAssets();
    this.saveFileSignatures(path.join(this.outputDir, "signatures.json"));
    const endTime = Date.now();
    console.log(
      `Total execution time: ${((endTime - startTime) / 1000).toFixed(
        2
      )} seconds`
    );
  }
}

const inputDir = path.join(__dirname, "src");
const outputDir = path.join(__dirname, "dist");
const builder = new Builder(inputDir, outputDir);
builder.buildAll();
