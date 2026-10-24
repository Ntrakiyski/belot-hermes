import test from 'node:test';
import assert from 'node:assert/strict';
import { selectMatch } from '../packages/matchmaker/select-match.mjs';
const solo=(id,t=0)=>({id,accountIds:[id],queuedAtMs:t,connected:true});
const duo=(id,a,b,t=0)=>({id,accountIds:[a,b],queuedAtMs:t,connected:true});
test('four solo users may form a match with four distinct accounts',()=>{
 const m=selectMatch([solo('a'),solo('b'),solo('c'),solo('d')],0);assert.equal(m.kind,'four-solos');assert.equal(new Set(m.accountIds).size,4);
});
test('two duos form opposing teams intact',()=>{
 const m=selectMatch([duo('ab','a','b'),duo('cd','c','d')],0);assert.equal(m.kind,'duo-vs-duo');assert.deepEqual(m.teams,[['a','b'],['c','d']]);
});
test('mixed matchmaking waits 5 seconds but not longer than necessary',()=>{
 const ts=[duo('ab','a','b'),solo('c'),solo('d')];
 assert.equal(selectMatch(ts,4999),null);
 assert.equal(selectMatch(ts,5000).kind,'duo-vs-two-solos');
 assert.equal(selectMatch(ts,15000).kind,'duo-vs-two-solos');
});
test('fewer than four humans never yields a match or bot',()=>{
 assert.equal(selectMatch([duo('ab','a','b'),solo('c')],150000),null);
});
test('prefer balanced candidate even if mixed is available',()=>{
 const ts=[duo('ab','a','b'),solo('c'),solo('d'),duo('ef','e','f')];
 assert.equal(selectMatch(ts,15000).kind,'duo-vs-duo');
});
test('same account across tickets cannot form a valid four-human match',()=>{
 assert.equal(selectMatch([solo('a'),solo('a'),solo('b'),solo('c')],50000),null);
});
test('disconnected user is excluded',()=>{
 const ts=[solo('a'),solo('b'),{...solo('c'),connected:false},solo('d')];
 assert.equal(selectMatch(ts,50000),null);
});
