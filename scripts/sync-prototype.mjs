import { readFile, writeFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import { resolve, dirname } from 'node:path';
const root = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const source = await readFile(resolve(root, 'prototype/index.html'));
const check = process.argv.includes('--check');
for (const env of ['demo', 'main']) {
  const filename = resolve(root, `${env}/index.html`);
  if (check) {
    const existing = await readFile(filename);
    if (!existing.equals(source)) throw new Error(`${filename} diverged from prototype source`);
  } else await writeFile(filename, source);
}
console.log(check ? 'Both snapshots match the canonical prototype.' : 'Updated both snapshots.');
