import { execSync } from "child_process";
import fs from "fs";
import path from "path";

// ask Go for GOROOT
const goroot = execSync("go env GOROOT").toString().trim();
const src = path.join(goroot, "lib", "wasm", "wasm_exec.js");
const dest = "./src/build/wasm_exec.js";

// copy the file
fs.copyFileSync(src, dest);

console.log("Copied:", src, "→", dest);
