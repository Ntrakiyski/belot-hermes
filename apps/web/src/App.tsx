import { useCallback, useEffect, useRef, useState } from 'react';

type Card = { suit: string; rank: string };
type Played = { seat: number; card: Card };
type Seat = { name: string; bot: boolean; cards: number };
type Declaration = { id: string; kind: 'sequence' | 'four'; name: string; points: number; cards: Card[] };
type Result = {
  contract: string; awarded: [number, number]; notes: string[]; allPassed: boolean;
  rawPoints: [number, number]; declarationPoints: [number, number]; belotPoints: [number, number];
  hanging: boolean; carried: number; carryAwarded: number; succeeded: boolean; multiplier: number;
  cardPoints: [number, number]; capotPoints: [number, number];
};
type View = {
  table: string; practice: boolean; ruleset: string; you: number; seats: Seat[]; started: boolean;
  seq: number; dealNumber: number; phase: string; dealer: number; turn: number; contract: string; bidder: number; mult: number;
  hand: Card[] | null; legal: Card[] | null; legalWhy: string; trick: Played[] | null; lastTrick: Played[] | null;
  scores: [number, number]; tricks: [number, number]; captured: [number, number]; target: number; result: Result | null;
  winner: number; log: string[] | null; canBid: string[] | null; collecting: boolean; trickWinner: number; trickPoints: number;
  nextActionAt: number; pending: number; bidHistory: { seat: number; action: string; contract: string }[] | null;
  history: { number: number; contract: string; awarded: [number, number]; scores: [number, number]; allPassed: boolean; hanging: boolean }[] | null;
  declarationOptions: Declaration[] | null; declared: { seat: number; kind: string; name: string; points: number }[] | null;
  canBelot: Card[] | null; belotPoints: [number, number];
};
type Session = { table: string; token: string };
type Health = { version: string; practiceEnabled: boolean };

const SYMBOL: Record<string, string> = { clubs: '♣', diamonds: '♦', hearts: '♥', spades: '♠' };
const SUIT: Record<string, string> = { clubs: 'спатия', diamonds: 'каро', hearts: 'купа', spades: 'пика' };
const CONTRACT: Record<string, string> = { clubs: 'Спатия', diamonds: 'Каро', hearts: 'Купа', spades: 'Пика', notrump: 'Без коз', alltrump: 'Всичко коз' };
const CONTRACTS = Object.keys(CONTRACT);
const isRed = (suit: string) => suit === 'hearts' || suit === 'diamonds';
const same = (a: Card, b: Card) => a.suit === b.suit && a.rank === b.rank;
const cardName = (card: Card) => `${card.rank} ${SUIT[card.suit] ?? card.suit}`;
function declarationChoices(options: Declaration[], selected: string[]) {
  const fourCards = options.filter(option => option.kind === 'four' && selected.includes(option.id)).flatMap(option => option.cards);
  return options.filter(option => option.kind === 'four' || !options.some(longer =>
    longer.kind === 'sequence' && longer.cards.length > option.cards.length &&
    option.cards.every(card => longer.cards.some(candidate => same(card, candidate))) &&
    !longer.cards.some(card => fourCards.some(candidate => same(card, candidate))),
  ));
}


const ERRORS: Record<string, string> = {
  'table not found': 'Масата вече не съществува. Създай нова или въведи друг код.',
  'table is full': 'На тази маса вече има четирима играчи.',
  'bad token': 'Мястото ти вече не е достъпно. Влез отново с код или създай маса.',
  'not seated at this table': 'Това място вече не е достъпно.',
  'not your turn': 'Изчакай своя ред.',
  'waiting for four players': 'Чакаме четирима играчи.',
  'deal is not finished': 'Раздаването още не е приключило.',
  'not in bidding phase': 'Наддаването вече е приключило.',
  'not in playing phase': 'Сега не се играят карти.',
  'that card is not in your hand': 'Тази карта вече не е в ръката ти.',
  'bid must be higher than the current contract': 'Избери по-висока обява.',
  'cannot double now': 'В момента не може да се обяви контра.',
  'cannot redouble now': 'В момента не може да се обяви реконтра.',
};
const HINTS: Record<string, string> = {
  'You lead: any card.': 'Ти започваш. Избери карта.',
  'Follow the led suit.': 'Отговори на поисканата боя.',
  'No card of the led suit: play anything.': 'Нямаш от поисканата боя. Избери карта.',
  'Follow suit and beat the highest card if you can.': 'Отговори на боята и качи, ако можеш.',
  'Follow suit (you cannot beat the highest card).': 'Отговори на поисканата боя.',
  'Trump led: you must beat the highest trump if you can.': 'Поискан е коз. Качи, ако можеш.',
  'Trump led: follow with a trump.': 'Отговори с коз.',
  'No trumps: play anything.': 'Нямаш коз. Избери карта.',
  'Your partner is winning: play anything.': 'Партньорът ти печели. Можеш да избереш всяка карта.',
  'No led suit and no trumps: play anything.': 'Нямаш боя и коз. Избери карта.',
  'You must over-trump.': 'Трябва да надцакаш.',
  'You cannot over-trump: play anything.': 'Не можеш да надцакаш. Избери карта.',
  'You cannot follow suit: you must trump.': 'Нямаш от поисканата боя. Цакай с коз.',
};

class APIError extends Error {
  constructor(public status: number, message: string) { super(message); }
}
async function api(path: string, body?: unknown, token?: string) {
  const response = await fetch(`/api/tables${path}`, {
    method: body === undefined ? 'GET' : 'POST',
    headers: { 'Content-Type': 'application/json', ...(token ? { 'X-Token': token } : {}) },
    body: body === undefined ? undefined : JSON.stringify(body),
    signal: AbortSignal.timeout(8000),
  });
  const data = await response.json().catch(() => ({}));
  if (!response.ok) throw new APIError(response.status, ERRORS[data.error] ?? (typeof data.error === 'string' ? data.error : 'Действието не е прието. Обнови масата и опитай отново.'));
  return data;
}
function readSession(): Session | null {
  try {
    const value = JSON.parse(sessionStorage.getItem('belot-sess') ?? 'null');
    return value && typeof value.table === 'string' && typeof value.token === 'string' ? value : null;
  } catch { return null; }
}
function CardFace({ card, onClick, dim = false, winning = false }: { card: Card; onClick?: () => void; dim?: boolean; winning?: boolean }) {
  return <button type="button" className={`card ${isRed(card.suit) ? 'red' : ''} ${dim ? 'dim' : ''} ${winning ? 'winning-card' : ''}`}
    disabled={!onClick} onClick={onClick} aria-label={cardName(card)}>
    <span className="card-corner"><b>{card.rank}</b><span>{SYMBOL[card.suit]}</span></span>
    <span className="card-pip" aria-hidden="true">{SYMBOL[card.suit]}</span>
    <span className="card-corner reverse" aria-hidden="true"><b>{card.rank}</b><span>{SYMBOL[card.suit]}</span></span>
  </button>;
}

export default function App() {
  const [name, setName] = useState(() => localStorage.getItem('belot-name') ?? '');
  const [code, setCode] = useState(() => new URLSearchParams(location.search).get('table')?.toUpperCase().slice(0, 4) ?? '');
  const [session, setSession] = useState<Session | null>(readSession);
  const [view, setView] = useState<View | null>(null);
  const [health, setHealth] = useState<Health | null>(null);
  const [error, setError] = useState('');
  const [offline, setOffline] = useState(false);
  const [busy, setBusy] = useState(false);
  const [selectedDeclarations, setSelectedDeclarations] = useState<string[]>([]);
  const [announceBelot, setAnnounceBelot] = useState(false);
  const [copied, setCopied] = useState(false);
  const busyRef = useRef(false);
  const pollingRef = useRef(false);
  const requestRef = useRef(0);
  const sessionRef = useRef(session);

  const saveSession = useCallback((value: Session | null) => {
    requestRef.current++;
    sessionRef.current = value;
    setSession(value); setView(null); setOffline(false); setCopied(false); setSelectedDeclarations([]); setAnnounceBelot(false);
    if (value) sessionStorage.setItem('belot-sess', JSON.stringify(value));
    else sessionStorage.removeItem('belot-sess');
  }, []);

  useEffect(() => {
    let active = true;
    const check = () => fetch('/healthz', { signal: AbortSignal.timeout(8000) }).then(response => {
      if (!response.ok) throw new Error();
      return response.json();
    }).then(data => { if (active) { setHealth(data); if (!sessionRef.current) setOffline(false); } }).catch(() => { if (active) setOffline(true); });
    void check();
    const interval = setInterval(check, 15000);
    return () => { active = false; clearInterval(interval); };
  }, []);

  const refresh = useCallback(async () => {
    if (!session || busyRef.current || pollingRef.current) return;
    pollingRef.current = true;
    const request = ++requestRef.current;
    try {
      const data: View = await api(`/${session.table}`, undefined, session.token);
      if (request !== requestRef.current || sessionRef.current !== session) return;
      setView(data); setOffline(false);
    } catch (failure) {
      if (request !== requestRef.current || sessionRef.current !== session) return;
      if (failure instanceof APIError && (failure.status === 403 || failure.status === 404)) {
        saveSession(null); setError(failure.message);
      } else { setOffline(true); }
    } finally { pollingRef.current = false; }
  }, [session, saveSession]);

  useEffect(() => {
    void refresh();
    const interval = setInterval(refresh, 350);
    return () => { clearInterval(interval); requestRef.current++; };
  }, [refresh]);
  useEffect(() => { setSelectedDeclarations([]); setAnnounceBelot(false); }, [view?.table, view?.dealNumber]);

  const act = async (path: string, body: object) => {
    if (!session || busyRef.current) return;
    const activeSession = session;
    const request = ++requestRef.current;
    busyRef.current = true; setBusy(true); setError('');
    try {
      const data: View = await api(`/${session.table}${path}`, body, session.token);
      if (request === requestRef.current && sessionRef.current === activeSession) {
        setView(data); setOffline(false);
        if (path === '/play') { setSelectedDeclarations([]); setAnnounceBelot(false); }
      }
    } catch (failure) {
      if (sessionRef.current !== activeSession) return;
      if (failure instanceof APIError && (failure.status === 403 || failure.status === 404)) saveSession(null);
      setError(failure instanceof APIError ? failure.message : 'Няма връзка със сървъра. Опитай отново.');
    } finally { busyRef.current = false; setBusy(false); }
  };
  const enter = async (practice = false, join = false) => {
    if (busyRef.current) return;
    const entryRequest = ++requestRef.current;
    busyRef.current = true; setBusy(true); setError('');
    localStorage.setItem('belot-name', name.trim());
    try {
      const nextSession: Session = await api(join ? `/${code}/join` : '', { name: name.trim() });
      if (entryRequest !== requestRef.current) return;
      saveSession(nextSession);
      if (practice) {
        const data: View = await api(`/${nextSession.table}/bots`, {}, nextSession.token);
        if (sessionRef.current === nextSession) setView(data);
      }
    } catch (failure) {
      setError(failure instanceof APIError ? failure.message : 'Няма връзка със сървъра. Опитай отново.');
    } finally { busyRef.current = false; setBusy(false); }
  };
  const leave = () => { saveSession(null); setError(''); };
  const copyInvite = async () => {
    if (!view) return;
    try {
      await navigator.clipboard.writeText(`${location.origin}/?table=${view.table}`);
      setCopied(true);
    } catch { setError(`Кодът на масата е ${view.table}. Сподели го с приятелите си.`); }
  };
  const version = health?.version ?? '0.3.0';
  const status = offline ? <div className="connection-warning" role="status">Връзката прекъсна. Опитваме отново… {session && <button className="text-button" onClick={() => void refresh()}>Опитай сега</button>}</div> : null;
  const alert = error && <div className="error-message" role="alert">{error}<button type="button" className="text-button" onClick={() => setError('')}>Разбрах</button></div>;
  const header = <header className="app-header"><a className="wordmark" href="/" onClick={event => { event.preventDefault(); leave(); }}>БЕЛОТ<span className="wordmark-suit" aria-hidden="true">♣</span></a><span className="version">v{version} · ДЕМО</span></header>;

  if (!session) return <main className="shell home-shell">
    {header}{status}
    <section className="home-content">
      <div className="home-cards" aria-hidden="true"><span className="sample-card">A<span>♠</span></span><span className="sample-card red">K<span>♥</span></span><span className="sample-card red">Q<span>♥</span></span></div>
      <h1>Една маса.<br />Още една игра.</h1>
      <p className="home-description">Български белот до 151. Тренирай с ботове или покани трима приятели.</p>
      <label className="field-label" htmlFor="player-name">Как да те наричаме?</label>
      <input id="player-name" placeholder="Твоето име" autoComplete="nickname" maxLength={20} value={name} onChange={event => setName(event.target.value)} />
      <button className="primary-button" disabled={busy || !health?.practiceEnabled} onClick={() => void enter(true)}>{busy ? 'Свързваме те…' : 'Практика с ботове'}<span aria-hidden="true">→</span></button>
      {!health?.practiceEnabled && health && <p className="field-help">Практиката не е включена на този сървър.</p>}
      <button className="secondary-button" disabled={busy} onClick={() => void enter()}>Създай маса с приятели</button>
      <div className="join-section"><label className="field-label" htmlFor="table-code">Имаш код за маса?</label><form className="join-form" onSubmit={event => { event.preventDefault(); if (code.length === 4) void enter(false, true); }}><input id="table-code" placeholder="Код на маса" maxLength={4} autoCapitalize="characters" value={code} onChange={event => setCode(event.target.value.toUpperCase().replace(/[^A-Z0-9]/g, ''))} /><button className="secondary-button" disabled={busy || code.length !== 4}>Влез</button></form></div>
      {alert}
      <p className="honest-note">Гост маси без регистрация. Публично намиране на играчи още не е включено. Ботовете участват само в ясно обозначена практика.</p>
    </section>
    <footer className="home-footer">Четирима на масата. Два отбора. До 151 точки.</footer>
  </main>;

  if (!view) return <main className="shell home-shell">{header}{status}<section className="loading-screen" role="status"><span className="loading-suit" aria-hidden="true">♣</span><h1>Свързваме те с масата</h1><p>Възстановяваме мястото ти…</p><button className="secondary-button" onClick={leave}>Обратно към началото</button></section>{alert}</main>;

  const me = view.you, myTeam = me % 2, otherTeam = 1 - myTeam;
  const seatName = (seat: number) => seat === me ? 'Ти' : view.seats[seat]?.name || 'Свободно място';
  const positions = ['bottom', 'right', 'top', 'left'];
  const relative = (seat: number) => (seat - me + 4) % 4;
  const myTurn = view.turn === me && !view.collecting;
  const hand = view.hand ?? [];
  const done = view.phase === 'dealover' || view.phase === 'matchover';
  const currentTrick = view.collecting ? view.lastTrick ?? [] : view.trick ?? [];
  const canPlay = myTurn && view.phase === 'playing' && !busy && !offline;
  const declarations = [...(view.declarationOptions ?? [])].sort((a, b) => b.points - a.points || b.cards.length - a.cards.length);
  const offeredDeclarations = declarationChoices(declarations, selectedDeclarations);
  const showBelot = canPlay && (view.canBelot?.length ?? 0) > 0;
  const currentBid = CONTRACTS.indexOf(view.contract);
  const turnText = view.collecting ? `${seatName(view.trickWinner)} печели взятката` : done ? 'Раздаването приключи' : myTurn ? (view.phase === 'bidding' ? 'Ти обявяваш' : 'Твой ред') : `${seatName(view.turn)} ${view.phase === 'bidding' ? 'обявява' : 'играе'}…`;

  return <main className="shell game-shell">
    {header}{status}
    <div className="table-meta"><span>{view.practice ? 'Практика с ботове' : `Маса ${view.table}`}</span><span>Раздаване {view.dealNumber || 1} · {Math.min(view.tricks[0] + view.tricks[1] + (view.collecting ? 0 : 1), 8)}/8 взятки</span><button className="text-button" onClick={leave}>Напусни</button></div>
    {!view.started ? <section className="waiting-room">
      <span className="section-caption">МАСА С ПРИЯТЕЛИ</span><h1>Запази им място.</h1><p>Сподели кода или поканата. Играта започва, когато се присъединят четирима.</p>
      <div className="room-code" aria-label={`Код на масата ${view.table}`}>{view.table}</div><button className="secondary-button" onClick={() => void copyInvite()}>{copied ? 'Поканата е копирана' : 'Копирай покана'}</button>
      <ol className="waiting-seats">{view.seats.map((seat, index) => <li key={index}><span className={`seat-avatar team-${index % 2}`} aria-hidden="true">{seat.name ? seat.name.slice(0, 1).toUpperCase() : '+'}</span><span>{seatName(index)}<small>{index % 2 === myTeam ? 'Твоят отбор' : 'Другият отбор'}</small></span>{seat.name && <span className="joined-mark">На масата</span>}</li>)}</ol>
      <p className="honest-note">{view.seats.filter(seat => seat.name).length}/4 места заети · Тази маса е за хора. Практиката с ботове се стартира отделно от началото.</p>{alert}
    </section> : <>
      <section className="scoreboard" aria-label="Резултат на мача"><div className="team-score"><span>Ние</span><strong>{view.scores[myTeam]}</strong><small>{view.tricks[myTeam]} взятки · {view.captured?.[myTeam] ?? 0} т.</small></div><div className="contract-display"><span>{view.contract ? CONTRACT[view.contract] : 'Наддаване'}</span>{view.contract && <small>{SYMBOL[view.contract] ?? ''} {seatName(view.bidder)}{view.mult > 1 ? ` · ×${view.mult}` : ''}</small>}<span className="target-label">до {view.target}</span></div><div className="team-score opponents"><span>Те</span><strong>{view.scores[otherTeam]}</strong><small>{view.tricks[otherTeam]} взятки · {view.captured?.[otherTeam] ?? 0} т.</small></div></section>
      <section className={`board ${view.collecting ? 'collecting' : ''}`} aria-label="Игрална маса">
        <div className="felt" aria-hidden="true"><span>БЕЛОТ</span></div>
        {[2, 1, 3].map(offset => { const seatIndex = (me + offset) % 4, seat = view.seats[seatIndex]; return <div key={offset} className={`player-seat ${positions[offset]} ${view.turn === seatIndex && !done && !view.collecting ? 'active-seat' : ''}`}><div className="seat-avatar" aria-hidden="true">{seat.name.slice(0, 1).toUpperCase()}</div><strong>{seat.name}</strong><span>{seatIndex % 2 === myTeam ? 'Партньор' : 'Съперник'}{view.dealer === seatIndex ? ' · раздава' : ''}</span><div className="card-backs" aria-label={`${seat.cards} карти`}>{Array.from({ length: Math.min(seat.cards, 8) }, (_, index) => <i key={index} />)}</div></div>; })}
        <div className="trick-cards">{currentTrick.map(played => <div key={`${played.seat}-${played.card.suit}-${played.card.rank}`} className={`played ${positions[relative(played.seat)]}`}><CardFace card={played.card} winning={view.collecting && view.trickWinner === played.seat} /></div>)}</div>
        {view.collecting && <div className="trick-result" role="status"><strong>{seatName(view.trickWinner)}</strong><span>+{view.trickPoints} точки</span></div>}
        {!currentTrick.length && view.phase === 'bidding' && <div className="table-message">Пет карти.<br /><span>Избери своята игра.</span></div>}
      </section>
      <div className={`turn-status ${myTurn ? 'your-turn' : ''}`} role="status"><span className="turn-dot" aria-hidden="true" /><strong>{turnText}</strong>{myTurn && view.phase === 'playing' && <span>{HINTS[view.legalWhy] ?? view.legalWhy ?? 'Избери една от осветените карти.'}</span>}</div>
      {view.phase === 'bidding' && <section className="bidding-panel" aria-label="Наддаване"><div className="bid-history">{(view.bidHistory ?? []).slice(-4).map((bid, index) => <span key={index}><strong>{seatName(bid.seat)}</strong> {bid.action === 'pass' ? 'пас' : bid.action === 'double' ? 'контра' : bid.action === 'redouble' ? 'реконтра' : CONTRACT[bid.contract]}</span>)}</div><div className="bid-options">{CONTRACTS.map((contract, index) => <button className={`bid-button ${isRed(contract) ? 'red-bid' : ''}`} key={contract} disabled={!myTurn || busy || offline || index <= currentBid || !view.canBid?.includes('bid')} onClick={() => void act('/bid', { action: 'bid', contract })}><span>{SYMBOL[contract] ?? (contract === 'notrump' ? 'БК' : 'ВК')}</span>{CONTRACT[contract]}</button>)}</div><div className="bid-actions"><button className="secondary-button" disabled={!myTurn || busy || offline || !view.canBid?.includes('pass')} onClick={() => void act('/bid', { action: 'pass' })}>Пас</button>{view.canBid?.includes('double') && <button className="secondary-button" disabled={busy || offline} onClick={() => void act('/bid', { action: 'double' })}>Контра</button>}{view.canBid?.includes('redouble') && <button className="secondary-button" disabled={busy || offline} onClick={() => void act('/bid', { action: 'redouble' })}>Реконтра</button>}</div></section>}
      {canPlay && declarations.length > 0 && <fieldset className="declarations"><legend>Избери анонси, после изиграй карта <span>по избор</span></legend>{offeredDeclarations.map(declaration => { const selected = selectedDeclarations.includes(declaration.id); const conflict = !selected && declarations.some(other => selectedDeclarations.includes(other.id) && other.cards.some(card => declaration.cards.some(candidate => same(card, candidate)))); return <label key={declaration.id} className={conflict ? 'disabled-declaration' : ''}><input type="checkbox" disabled={conflict} checked={selected} onChange={() => setSelectedDeclarations(value => { const next = selected ? value.filter(id => id !== declaration.id) : [...value, declaration.id]; const available = declarationChoices(declarations, next); return next.filter(id => available.some(option => option.id === id)); })} /><span><strong>{declaration.name} · {declaration.points} т.</strong><small>{declaration.cards.map(card => `${card.rank}${SYMBOL[card.suit]}`).join(' ')}</small></span></label>; })}</fieldset>}
      {showBelot && <label className="belot-option"><input type="checkbox" checked={announceBelot} onChange={event => setAnnounceBelot(event.target.checked)} />Обяви белот с поп или дама коз</label>}
      {(view.declared?.length ?? 0) > 0 && <div className="public-declarations" aria-label="Обявени анонси">{view.declared?.map((declaration, index) => <span key={`${declaration.seat}-${index}`}>{seatName(declaration.seat)}: {declaration.name} · {declaration.points}</span>)}</div>}
      {done && !view.collecting && view.result && <section className="result-panel" aria-label="Резултат от раздаването"><span className="section-caption">{view.phase === 'matchover' ? 'КРАЙ НА МАЧА' : `РАЗДАВАНЕ ${view.dealNumber}`}</span><h2>{view.phase === 'matchover' ? (view.winner === myTeam ? 'Мачът е ваш!' : 'Победа за съперниците') : view.result.allPassed ? 'Всички пас. Раздаваме отново.' : view.result.hanging ? 'Висяща игра' : 'Точките са на масата.'}</h2>{!view.result.allPassed && <><div className="result-totals"><span>Ние <strong>+{view.result.awarded[myTeam]}</strong></span><span>Те <strong>+{view.result.awarded[otherTeam]}</strong></span></div><dl className="score-ledger"><div><dt>Карти и последна взятка</dt><dd>{view.result.cardPoints?.[myTeam] ?? 0} / {view.result.cardPoints?.[otherTeam] ?? 0}</dd></div><div><dt>Капо</dt><dd>{view.result.capotPoints?.[myTeam] ?? 0} / {view.result.capotPoints?.[otherTeam] ?? 0}</dd></div><div><dt>Анонси</dt><dd>{view.result.declarationPoints?.[myTeam] ?? 0} / {view.result.declarationPoints?.[otherTeam] ?? 0}</dd></div><div><dt>Белот</dt><dd>{view.result.belotPoints?.[myTeam] ?? 0} / {view.result.belotPoints?.[otherTeam] ?? 0}</dd></div><div className="ledger-total"><dt>Общо точки преди закръгляне</dt><dd>{view.result.rawPoints[myTeam]} / {view.result.rawPoints[otherTeam]}</dd></div></dl><p className="contract-outcome">{view.result.hanging ? 'Обявата виси.' : view.result.succeeded ? 'Обявилият отбор изпълни договора.' : 'Обявилият отбор е вътре.'}{view.result.multiplier > 1 ? ` · ×${view.result.multiplier}` : ''}</p>{view.result.hanging && <p>Точките на обявилия отбор остават за следващото раздаване.</p>}{view.result.carryAwarded > 0 && <p>Разпределени висящи точки: {view.result.carryAwarded}.</p>}<details className="score-explanation"><summary>Как се получи резултатът?</summary><ul>{view.result.notes.map((note, index) => <li key={index}>{note}</li>)}</ul></details></>}<button className="primary-button" disabled={busy || view.collecting || offline} onClick={() => void act('/next', {})}>{view.phase === 'matchover' ? 'Нов мач' : 'Следващо раздаване'}<span aria-hidden="true">→</span></button></section>}
      {hand.length > 0 && <section className="hand-section" aria-label="Твоите карти"><div className="hand-label"><span>Твоята ръка</span><span>{view.phase === 'bidding' ? 'Още 3 карти след наддаването' : canPlay ? 'Докосни карта, за да я изиграеш' : 'Изчакай своя ред'}</span></div><div className={`hand ${canPlay ? 'hand-active' : ''}`} style={{ gridTemplateColumns: `repeat(${hand.length}, minmax(0, 1fr))`, width: `min(100%, calc(${hand.length} * var(--card-width) + 25px))` }}>{hand.map(card => { const legal = canPlay && (view.legal ?? []).some(candidate => same(candidate, card)); return <CardFace key={`${card.suit}-${card.rank}`} card={card} dim={canPlay && !legal} onClick={legal ? () => void act('/play', { ...card, declarations: selectedDeclarations.filter(id => offeredDeclarations.some(option => option.id === id)), belot: announceBelot && (view.canBelot ?? []).some(candidate => same(candidate, card)) }) : undefined} />; })}</div></section>}
      {view.pending > 0 && <p className="pending-points">Висящи точки за следващото раздаване: <strong>{view.pending}</strong></p>}
      {alert}
      <details className="history-panel"><summary>Резултати по раздавания{view.history?.length ? ` (${view.history.length})` : ''}</summary>{view.history?.length ? <table><thead><tr><th scope="col">№</th><th scope="col">Игра</th><th scope="col">Ние</th><th scope="col">Те</th></tr></thead><tbody>{view.history.map(deal => <tr key={deal.number}><th scope="row">{deal.number}</th><td>{deal.allPassed ? 'Всички пас' : `${CONTRACT[deal.contract]}${deal.hanging ? ' · висяща' : ''}`}</td><td>+{deal.awarded[myTeam]}</td><td>+{deal.awarded[otherTeam]}</td></tr>)}</tbody></table> : <p>Резултатите ще се появят след първото раздаване.</p>}</details>
      <footer className="game-footer"><span>{view.practice ? 'Практика с ботове' : 'Гост маса с хора'} · правила {view.ruleset}</span><span>v{version}</span></footer>
    </>}
  </main>;
}
