import test from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp, mkdir, readFile, writeFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { spawnSync } from 'node:child_process';

test('persona sync preserves the live file on missing input and backs up a complete replacement', async () => {
  const dir = await mkdtemp(join(tmpdir(), 'belot-profile-'));
  try {
    const profile = join(dir, 'profile'), source = join(dir, 'source');
    await mkdir(profile); await mkdir(source);
    await writeFile(join(profile, 'SOUL.md'), 'original');
    await writeFile(join(source, 'SOUL.md'), 'persona\n');
    const run = () => spawnSync('bash', ['ops/hermes/sync-profile.sh', source], {
      env: { ...process.env, HERMES_PROFILE_DIR: profile }, encoding: 'utf8',
    });
    assert.notEqual(run().status, 0);
    assert.equal(await readFile(join(profile, 'SOUL.md'), 'utf8'), 'original');
    await writeFile(join(source, 'SOUL-profile-appendix.md'), 'appendix\n');
    assert.equal(run().status, 0);
    assert.equal(await readFile(join(profile, 'SOUL.md'), 'utf8'), 'persona\nappendix\n');
    assert.equal(await readFile(join(profile, 'SOUL.md.bak'), 'utf8'), 'original');
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
});
