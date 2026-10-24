/**
 * Deterministic matchmaking-policy example. NOT a server or an online queue.
 * A ticket is {id,accountIds:[one or two unique IDs], queuedAtMs, connected:true}.
 * Caller must atomically lock tickets in a real backend; this function only selects.
 */
export function selectMatch(tickets, nowMs, { mixedFallbackMinMs = 5000 } = {}) {
  const valid = tickets.filter(t => t.connected === true && Array.isArray(t.accountIds)
    && t.accountIds.length >= 1 && t.accountIds.length <= 2 &&
    new Set(t.accountIds).size === t.accountIds.length && Number.isFinite(t.queuedAtMs)
    && t.queuedAtMs <= nowMs).sort((a,b)=>a.queuedAtMs-b.queuedAtMs||String(a.id).localeCompare(String(b.id)));
  const solos = valid.filter(t=>t.accountIds.length===1), duos=valid.filter(t=>t.accountIds.length===2);
  const distinct = ts=>new Set(ts.flatMap(t=>t.accountIds)).size===ts.reduce((a,t)=>a+t.accountIds.length,0);
  const combos = [];
  // Small reference implementation: real backend must use bounded queue windows/indexes.
  for (let i=0;i<solos.length;i++) for(let j=i+1;j<solos.length;j++)
    for(let k=j+1;k<solos.length;k++) for(let l=k+1;l<solos.length;l++) {
      const ts=[solos[i],solos[j],solos[k],solos[l]];
      if (distinct(ts)) combos.push({kind:'four-solos',tickets:ts,teams:[
        [ts[0].accountIds[0],ts[2].accountIds[0]], [ts[1].accountIds[0],ts[3].accountIds[0]]
      ]});
    }
  for(let i=0;i<duos.length;i++) for(let j=i+1;j<duos.length;j++) {
    const ts=[duos[i],duos[j]];
    if(distinct(ts)) combos.push({kind:'duo-vs-duo',tickets:ts,teams:[ts[0].accountIds,ts[1].accountIds]});
  }
  if(combos.length===0) for(const d of duos) for(let i=0;i<solos.length;i++) for(let j=i+1;j<solos.length;j++) {
    const ts=[d,solos[i],solos[j]];
    if(!distinct(ts) || nowMs-Math.min(...ts.map(t=>t.queuedAtMs))<mixedFallbackMinMs)continue;
    combos.push({kind:'duo-vs-two-solos',tickets:ts,teams:[d.accountIds,[solos[i].accountIds[0],solos[j].accountIds[0]]]});
  }
  if(!combos.length)return null;
  combos.sort((a,b)=>Math.min(...a.tickets.map(t=>t.queuedAtMs))-Math.min(...b.tickets.map(t=>t.queuedAtMs))
    || a.tickets.map(t=>String(t.id)).sort().join('/').localeCompare(b.tickets.map(t=>String(t.id)).sort().join('/')));
  const selected=combos[0];
  return {...selected, ticketIds:selected.tickets.map(t=>t.id), accountIds:selected.teams.flat()};
}
