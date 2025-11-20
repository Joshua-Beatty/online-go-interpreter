import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Configuration
const CHUNK_SIZE = 15 * 1024 * 1024; // 15MB in bytes
const WASM_FILE = path.join(__dirname, '../src/build/main.wasm');
const OUTPUT_DIR = path.join(__dirname, '../src/build');
const OUTPUT_PREFIX = 'main.wasm.part';

console.log('Starting WASM file splitting...');

// Read the WASM file
if (!fs.existsSync(WASM_FILE)) {
    console.error(`Error: WASM file not found at ${WASM_FILE}`);
    process.exit(1);
}

const wasmBuffer = fs.readFileSync(WASM_FILE);
const totalSize = wasmBuffer.length;
const numChunks = Math.ceil(totalSize / CHUNK_SIZE);

console.log(`WASM file size: ${(totalSize / (1024 * 1024)).toFixed(2)}MB`);
console.log(`Chunk size: ${(CHUNK_SIZE / (1024 * 1024)).toFixed(2)}MB`);
console.log(`Number of chunks: ${numChunks}`);

// Delete existing chunk files
const existingChunks = fs.readdirSync(OUTPUT_DIR).filter(f => f.startsWith(OUTPUT_PREFIX));
existingChunks.forEach(chunk => {
    fs.unlinkSync(path.join(OUTPUT_DIR, chunk));
    console.log(`Deleted existing chunk: ${chunk}`);
});

// Split the file into chunks
for (let i = 0; i < numChunks; i++) {
    const start = i * CHUNK_SIZE;
    const end = Math.min(start + CHUNK_SIZE, totalSize);
    const chunk = wasmBuffer.slice(start, end);

    const chunkFileName = `${OUTPUT_PREFIX}${i}`;
    const chunkFilePath = path.join(OUTPUT_DIR, chunkFileName);

    fs.writeFileSync(chunkFilePath, chunk);
    console.log(`Created ${chunkFileName} (${(chunk.length / (1024 * 1024)).toFixed(2)}MB)`);
}

console.log('WASM file splitting complete!');
