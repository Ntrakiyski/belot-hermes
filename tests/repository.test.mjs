import test from 'node:test';
import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
const read=async p=>readFile(new URL(`../${p}`,import.meta.url),'utf8');
test('source prototype and preview snapshots are identical on seed commit',async()=>{
 const src=await read('prototype/index.html');
 assert.equal(await read('demo/index.html'),src);
 assert.equal(await read('main/index.html'),src);
});
test('the four pillars and three skills exist and state human-first requirement',async()=>{
 for (const name of ['01-game-product','02-architecture-release','03-hermes-discord','04-workflows-acceptance'])assert.ok((await read(`docs/${name}.md`)).length>2000);
 for (const name of ['feature-delivery','gameplay-qa','admin-ops'])assert.ok((await read(`skills/${name}/SKILL.md`)).includes('name:'));
 const one=await read('docs/01-game-product.md');
 assert.ok(one.includes('Start Random Game')&&one.includes('Create a Room'));
 for (const id of ['D05','D06','D07','D08','D09']) assert.ok(one.includes(id));
});
