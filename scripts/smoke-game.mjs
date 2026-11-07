// Creates disposable tables: run only against a local or explicitly authorized DEMO instance.
import assert from 'node:assert/strict';

const base = process.env.BELOT_SMOKE_URL ?? 'http://127.0.0.1:3000';
const practice = process.argv.includes('--practice');
async function request(path, body, token, status = 200) {
  const response = await fetch(`${base}${path}`, {
    method: body === undefined ? 'GET' : 'POST',
    headers: { 'Content-Type': 'application/json', ...(token ? { 'X-Token': token } : {}) },
    body: body === undefined ? undefined : JSON.stringify(body),
    signal: AbortSignal.timeout(5000),
  });
  assert.equal(response.status, status, `${path}: unexpected HTTP status`);
  return response.headers.get('content-type')?.includes('application/json') ? response.json() : response.text();
}
const health = await request('/healthz');
assert.equal(health.status, 'ok');
assert.equal(health.multiplayerReady, false);
assert.equal(health.version, '0.3.0');
assert.equal(health.practiceEnabled, practice);
assert.match(await request('/'), /<div id="root">/);
const owner = await request('/api/tables', { name: 'Smoke 0' });
const path = `/api/tables/${owner.table}`;
await request(path, undefined, undefined, 403);
await request(`${path}/bid`, { action: 'pass' }, owner.token, 409);
if (!practice) {
  await request(`${path}/bots`, {}, owner.token, 404);
}
const seats = [owner];
for (let i = 1; i < 4; i++) seats.push(await request(`${path}/join`, { name: `Smoke ${i}` }));
await request(`${path}/join`, { name: 'Overflow' }, undefined, 409);
const views = await Promise.all(seats.map(s => request(path, undefined, s.token)));
assert.equal(new Set(views.flatMap(v => v.hand.map(c => `${c.rank}:${c.suit}`))).size, 20);
for (const [i, v] of views.entries()) {
  assert.ok((v.declarationOptions ?? []).every(d => d.cards.every(c => v.hand.some(h => h.suit === c.suit && h.rank === c.rank))));
  assert.equal(v.you, i);
  assert.equal(v.hand.length, 5);
  assert.equal(v.practice, false);
  assert.ok(!('hands' in v));
  assert.ok(v.seats.every(s => !('token' in s) && !('hand' in s)));
}
let state = views[0];
await request(`${path}/bid`, { action: 'pass' }, seats[(state.turn + 1) % 4].token, 400);
state = await request(`${path}/bid`, { action: 'bid', contract: 'clubs' }, seats[state.turn].token);
for (let i = 0; i < 3; i++) state = await request(`${path}/bid`, { action: 'pass' }, seats[state.turn].token);
assert.equal(state.phase, 'playing');
for (let i = 0; i < 32; i++) {
  const token = seats[state.turn].token;
  state = await request(path, undefined, token);
  assert.ok(state.legal.length > 0);
  state = await request(`${path}/play`, state.legal[0], token);
  if (state.collecting) {
    assert.equal(state.lastTrick.length, 4);
    assert.ok(state.trickWinner >= 0 && state.trickWinner < 4);
    await request(`${path}/play`, { suit: 'clubs', rank: '7' }, seats[state.turn].token, 409);
    const deadline = Date.now() + 5000;
    while (state.collecting && Date.now() < deadline) {
      await new Promise(resolve => setTimeout(resolve, 100));
      state = await request(path, undefined, token);
    }
    assert.equal(state.collecting, false, 'completed trick must be collected');
  }
}
assert.ok(['dealover', 'matchover'].includes(state.phase));
assert.equal(state.tricks[0] + state.tricks[1], 8);
assert.ok(state.result);
state = await request(`${path}/next`, {}, owner.token);
assert.equal(state.phase, 'bidding');
if (practice) {
  const p = await request('/api/tables', { name: 'Practice smoke' });
  const filled = await request(`/api/tables/${p.table}/bots`, {}, p.token);
  assert.equal(filled.table, p.table);
  assert.equal(filled.hand.length, 5);
  assert.equal(filled.seats.length, 4);
  const v = await request(`/api/tables/${p.table}`, undefined, p.token);
  assert.equal(v.practice, true);
  assert.equal(v.seats.filter(s => s.bot).length, 3);
  const bid = await request(`/api/tables/${p.table}/bid`, { action: 'bid', contract: 'alltrump' }, p.token);
  assert.equal(bid.phase, 'bidding');
  assert.equal(bid.turn, 1, 'first bot turn must remain visible');
  let pv = bid;
  let sawCollection = false;
  const deadline = Date.now() + 90000;
  while (Date.now() < deadline) {
    pv = await request(`/api/tables/${p.table}`, undefined, p.token);
    if (pv.collecting) {
      sawCollection = true;
      assert.equal(pv.lastTrick.length, 4);
    } else if (pv.phase === 'dealover' || pv.phase === 'matchover') {
      break;
    } else if (pv.turn === 0 && pv.phase === 'playing') {
      assert.ok(pv.legal.length > 0);
      await request(`/api/tables/${p.table}/play`, pv.legal[0], p.token);
    } else if (pv.turn === 0 && pv.phase === 'bidding') {
      await request(`/api/tables/${p.table}/bid`, { action: 'pass' }, p.token);
    }
    await new Promise(resolve => setTimeout(resolve, 150));
  }
  assert.ok(['dealover', 'matchover'].includes(pv.phase), 'paced practice must complete a deal');
  assert.ok(sawCollection, 'practice must display complete tricks');
  assert.equal(pv.tricks[0] + pv.tricks[1], 8);
  assert.ok(pv.result);
  assert.equal(pv.history.length, 1);
}
console.log(`PASS: health, UI, seat isolation, rejected actions, visible 8-trick deal, collection gate, next deal, bots ${practice ? 'enabled explicitly' : 'disabled'}`);
