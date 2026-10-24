import { useEffect, useState } from 'react';

type Health = { status: string; service: string; multiplayerReady: boolean };

export default function App() {
  const [health, setHealth] = useState<Health | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch('/healthz').then(r => { if (!r.ok) throw new Error('unavailable'); return r.json(); })
      .then(setHealth).catch(() => setError('Go server is not running'));
  }, []);

  return <main className="shell">
    <header><strong>BELOT</strong><span>v0.2.0 · development scaffold</span></header>
    <section className="hero">
      <div className="suits" aria-hidden="true">♣ ♦ ♥ ♠</div>
      <h1>Играй белот с истински хора</h1>
      <p>Публични мачове с четирима играчи и частни стаи с приятели.</p>
      <button disabled title="Matchmaking is not yet implemented">Старт на случайна игра</button>
      <button className="secondary" disabled title="Private rooms are not yet implemented">Създай стая</button>
      <small>Двата режима ще бъдат активирани, когато сървърът поддържа реални мачове.</small>
    </section>
    <footer>
      Go engine: {health ? `${health.status} · no multiplayer yet` : error ?? 'connecting…'}
    </footer>
  </main>;
}
