import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

// GitHub URL for the WASM file
const GITHUB_WASM_URL = 'https://github.com/Joshua-Beatty/online-go-interpreter/raw/refs/heads/main/src/build/main.wasm';

// Get __dirname equivalent in ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Path to the build directory
const buildDir = path.join(__dirname, '..', 'src');
const wasmPath = path.join(buildDir, 'build', 'main.wasm');
const htmlPath = path.join(buildDir, 'index.html');

console.log('Starting Cloudflare build process...');

// Delete main.wasm if it exists
if (fs.existsSync(wasmPath)) {
    fs.unlinkSync(wasmPath);
    console.log('✓ Deleted build/main.wasm');
} else {
    console.log('⚠ build/main.wasm not found, skipping deletion');
}

// Replace references to main.wasm in index.html
if (fs.existsSync(htmlPath)) {
    let htmlContent = fs.readFileSync(htmlPath, 'utf8');
    const originalContent = htmlContent;

    // Replace "build/main.wasm" with the GitHub URL
    htmlContent = htmlContent.replace(/["']build\/main\.wasm["']/g, `"${GITHUB_WASM_URL}"`);

    if (htmlContent !== originalContent) {
        fs.writeFileSync(htmlPath, htmlContent, 'utf8');
        console.log('✓ Replaced main.wasm references in index.html with GitHub URL');
    } else {
        console.log('⚠ No main.wasm references found to replace in index.html');
    }
} else {
    console.log('⚠ index.html not found');
}

console.log('Cloudflare build process complete!');
