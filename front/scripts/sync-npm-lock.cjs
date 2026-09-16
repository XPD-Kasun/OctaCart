const fs = require("fs");
const os = require("os");
const path = require("path");
const { execSync } = require("child_process");

const projectDir = path.join(__dirname, "..");
const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "sync-lockfile-"));

try {
    // Copy only what npm needs to resolve versions — never touches
    // the real node_modules, so no Windows file-lock issues.
    fs.copyFileSync(
        path.join(projectDir, "package.json"),
        path.join(tmpDir, "package.json")
    );

    const existingLock = path.join(projectDir, "package-lock.json");
    if (fs.existsSync(existingLock)) {
        fs.copyFileSync(existingLock, path.join(tmpDir, "package-lock.json"));
    }

    execSync("npm install --package-lock-only", {
        cwd: tmpDir,
        stdio: "inherit",
    });

    fs.copyFileSync(
        path.join(tmpDir, "package-lock.json"),
        existingLock
    );

    console.log("package-lock.json updated.");
} finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
}