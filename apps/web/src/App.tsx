import { useCallback, useEffect, useRef, useState } from 'react';

type Card = { suit: string; rank: string };
type Played = { seat: number; card: Card };
type SeatV = { name: string; bot: boolean; cards: number };
type Result = { contract: string; awarded: [number, number]; notes: string[]; allPassed: boolean; rawPoints: [number, number] };
type View = {
  table: string; practice: boolean; ruleset: string; you: number; seats: SeatV[]; started: boolean;
  phase: string; dealer: number; turn: number; contract: string; bidder: number; mult: number;
  hand: Card[]; legal: Card[] | null; legalWhy: string; trick: Played[] | null; lastTrick: Played[] | null;
  scores: [number, number]; tricks: [number, number]; target: number; result: Result | null;
  winner: number; log: string[]; canBid: string[] | null;
};
type Session = { table: string; token: string };

const SYM: Record<string, string> = { clubs: '♣', diamonds: '♦', hearts: '♥', spades: '♠' };
const CN: Record<string, string> = { clubs: '♣', diamonds: '♦', hearts: '♥', spades: '♠', notrump: 'Без коз', alltrump: 'Всичко коз' };
const red = (s: string) => s === 'hearts' || s === 'diamonds';
const same = (a: Card, b: Card) => a.suit === b.suit && a.rank === b.rank;

async function api(path: string, body?: unknown, token?: string) {
  const r = await fetch(`/api/tables${path}${token ? `?token=${token}` : ''}`, {
    method: body === undefined ? 'GET' : 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  const j = await r.json();
  if (!r.ok) throw new Error(j.error ?? 'error');
  return j;
}

function CardFace({ c, onClick, dim, hot }: { c: Card; onClick?: () => void; dim?: boolean; hot?: boolean }) {
  return <button className={`card ${red(c.suit) ? 'red' : ''} ${dim ? 'dim' : ''} ${hot ? 'hot' : ''}`} disabled={!onClick} onClick={onClick}>
    <b>{c.rank}</b><span>{SYM[c.suit]}</span>
  </button>;
}

export default function App() {
  const [name, setName] = useState(localStorage.getItem('belot-name') ?? '');
  const [code, setCode] = useState('');
  const [sess, setSess] = useState<Session | null>(() => { try { return JSON.parse(sessionStorage.getItem('belot-sess') ?? 'null'); } catch { return null; } });
  const [v, setV] = useState<View | null>(null);
  const [err, setErr] = useState('');
  const busy = useRef(false);

  const save = (s: Session | null) => { setSess(s); s ? sessionStorage.setItem('belot-sess', JSON.stringify(s)) : sessionStorage.removeItem('belot-sess'); };
  const refresh = useCallback(async () => {
    if (!sess || busy.current) return;
    try { setV(await api(`/${sess.table}`, undefined, sess.token)); } catch { /* keep last view */ }
  }, [sess]);
  useEffect(() => { refresh(); const i = setInterval(refresh, 1000); return () => clearInterval(i); }, [refresh]);

  const act = async (path: string, body: unknown) => {
    if (!sess) return;
    busy.current = true; setErr('');
    try { setV(await api(`/${sess.table}${path}`, body, sess.token)); } catch (e) { setErr((e as Error).message); }
    busy.current = false;
  };
  const enter = async (path: string, body: object) => {
    setErr(''); localStorage.setItem('belot-name', name);
    try { const r = await api(path, body); save({ table: r.table, token: r.token }); } catch (e) { setErr((e as Error).message); }
  };

  if (!sess || !v) return <main className="shell">
    <header><strong>BELOT</strong><span>тестов вариант · {sess ? 'свързване…' : 'без акаунти'}</span></header>
    <section className="hero">
      <div className="suits" aria-hidden="true">♣ ♦ ♥ ♠</div>
      <h1>Игра на белот</h1>
      <input placeholder="Твоето име" value={name} maxLength={20} onChange={e => setName(e.target.value)} />
      <button onClick={() => enter('', { name })}>Създай маса</button>
      <div className="row">
        <input placeholder="Код на маса" value={code} maxLength={4} onChange={e => setCode(e.target.value.toUpperCase())} />
        <button className="secondary" disabled={code.length < 4} onClick={() => enter(`/${code}/join`, { name })}>Влез</button>
      </div>
      {err && <p className="err">{err}</p>}
      <small>Споделете кода с трима приятели. Истински мачмейкинг и профили още не са готови.</small>
    </section></main>;

  const me = v.you, rel = (s: number) => (s - me + 4) % 4; // 0 bottom, 1 left, 2 top, 3 right
  const pos = ['bottom', 'left', 'top', 'right'];
  const myTurn = v.turn === me;
  const lastWinner = v.lastTrick ? undefined : undefined; void lastWinner;
  const teamName = (t: number) => (t === me % 2 ? 'Ние' : 'Те');
  const empty = v.seats.filter(s => !s.name).length;

  return <main className="shell table">
    <header><strong>BELOT</strong><span>Маса <b>{v.table}</b>{v.practice ? ' · ПРАКТИКА с ботове' : ''}</span></header>

    {!v.started && <section className="hero">
      <h1>Чакаме играчи ({4 - empty}/4)</h1>
      <p className="code">{v.table}</p>
      <ul className="seats">{v.seats.map((s, i) => <li key={i}>{s.name || '— свободно —'}{i === me ? ' (ти)' : ''}</li>)}</ul>
      {empty > 0 && <button className="secondary" onClick={() => act('/bots', {})}>Допълни с ботове (само за тест)</button>}
      <small>Ботовете са само за изпробване на двигателя. Истинските игри са с четирима хора.</small>
    </section>}

    {v.started && <>
      <div className="score">
        <span>Ние <b>{v.scores[me % 2]}</b></span>
        <span className="mid">{v.contract ? `${CN[v.contract]}${v.mult > 1 ? ` ×${v.mult}` : ''} · ${v.seats[v.bidder]?.name ?? ''}` : 'Наддаване'} · до {v.target}</span>
        <span>Те <b>{v.scores[1 - me % 2]}</b></span>
      </div>
      <div className="board">
        {[2, 1, 3].map(r => { const s = (me + r) % 4; return <div key={r} className={`seat ${pos[r]} ${v.turn === s && v.phase !== 'dealover' ? 'turn' : ''}`}>
          <span>{v.seats[s].name}{v.dealer === s ? ' ⓓ' : ''}</span><em>{v.seats[s].cards} 🂠</em></div>; })}
        <div className="felt">
          {(v.trick ?? []).map(p => <div key={p.seat} className={`played ${pos[rel(p.seat)]}`}><CardFace c={p.card} /></div>)}
          {!(v.trick?.length) && v.phase === 'playing' && v.lastTrick && <small className="faint">Последна ръка: {v.lastTrick.map(p => `${p.card.rank}${SYM[p.card.suit]}`).join(' ')}</small>}
        </div>
      </div>

      {v.phase === 'bidding' && <div className="panel">
        {myTurn ? <>
          <p>Твой ред за наддаване{v.contract ? ` · текущо: ${CN[v.contract]}` : ''}</p>
          <div className="bids">
            {(['clubs', 'diamonds', 'hearts', 'spades', 'notrump', 'alltrump'] as const).map(c => <button key={c} className={red(c) ? 'redbtn' : ''} onClick={() => act('/bid', { action: 'bid', contract: c })}>{CN[c]}</button>)}
          </div>
          <div className="bids">
            <button className="secondary" onClick={() => act('/bid', { action: 'pass' })}>Пас</button>
            {v.canBid?.includes('double') && <button className="secondary" onClick={() => act('/bid', { action: 'double' })}>Контра</button>}
            {v.canBid?.includes('redouble') && <button className="secondary" onClick={() => act('/bid', { action: 'redouble' })}>Реконтра</button>}
          </div></> : <p>Наддава {v.seats[v.turn].name}…</p>}
      </div>}

      {v.phase === 'playing' && <p className="hint">{myTurn ? `Твой ред. ${v.legalWhy}` : `Играе ${v.seats[v.turn].name}…`}</p>}

      {(v.phase === 'dealover' || v.phase === 'matchover') && v.result && <div className="panel">
        <h2>{v.phase === 'matchover' ? (v.winner === me % 2 ? '🏆 Печелите мача!' : 'Загубихте мача') : 'Край на ръката'}</h2>
        {v.result.allPassed ? <p>{v.result.notes[0]}</p> : <>
          <p>{teamName(me % 2)}: <b>+{v.result.awarded[me % 2]}</b> · {teamName(1 - me % 2)}: <b>+{v.result.awarded[1 - me % 2]}</b></p>
          <ul>{v.result.notes.map((n, i) => <li key={i}>{n}</li>)}</ul></>}
        <button onClick={() => act('/next', {})}>{v.phase === 'matchover' ? 'Нов мач' : 'Следваща ръка'}</button>
      </div>}

      <div className={`hand ${myTurn && v.phase === 'playing' ? 'active' : ''}`}>
        {v.hand.map(c => {
          const ok = v.phase === 'playing' && myTurn && (v.legal ?? []).some(l => same(l, c));
          return <CardFace key={c.suit + c.rank} c={c} hot={ok} dim={v.phase === 'playing' && myTurn && !ok} onClick={v.phase === 'playing' && myTurn ? () => act('/play', c) : undefined} />;
        })}
      </div>
      <details><summary>Дневник на играта · правила {v.ruleset}</summary><ol>{[...v.log].reverse().slice(0, 15).map((l, i) => <li key={i}>{l}</li>)}</ol></details>
    </>}
    {err && <p className="err">{err}</p>}
    <footer><a href="#" onClick={e => { e.preventDefault(); save(null); setV(null); }}>Напусни масата</a></footer>
  </main>;
}
