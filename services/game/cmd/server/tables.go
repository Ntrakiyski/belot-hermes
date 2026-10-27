package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/Ntrakiyski/belot-hermes/services/game/internal/belot"
)

type seat struct {
	Name  string
	Token string
	Bot   bool
}

type table struct {
	mu    sync.Mutex
	ID    string
	Seats [4]*seat
	Game  *belot.Game
}

type store struct {
	mu     sync.Mutex
	tables map[string]*table
}

func newStore() *store { return &store{tables: map[string]*table{}} }

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *store) routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tables", s.create)
	mux.HandleFunc("POST /api/tables/{id}/join", s.join)
	mux.HandleFunc("POST /api/tables/{id}/bots", s.addBots)
	mux.HandleFunc("GET /api/tables/{id}", s.state)
	mux.HandleFunc("POST /api/tables/{id}/bid", s.bid)
	mux.HandleFunc("POST /api/tables/{id}/play", s.play)
	mux.HandleFunc("POST /api/tables/{id}/next", s.next)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func fail(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func cleanName(n string) string {
	n = strings.TrimSpace(n)
	if len(n) > 20 {
		n = n[:20]
	}
	if n == "" {
		n = "Player"
	}
	return n
}

func (s *store) get(id string) *table {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tables[strings.ToUpper(id)]
}

func (s *store) create(w http.ResponseWriter, r *http.Request) {
	var in struct{ Name string }
	_ = json.NewDecoder(r.Body).Decode(&in)
	t := &table{ID: strings.ToUpper(randHex(2))}
	s.mu.Lock()
	for s.tables[t.ID] != nil {
		t.ID = strings.ToUpper(randHex(2))
	}
	s.tables[t.ID] = t
	s.mu.Unlock()
	t.mu.Lock()
	defer t.mu.Unlock()
	tok := t.seatPlayer(0, cleanName(in.Name))
	writeJSON(w, 200, map[string]any{"table": t.ID, "seat": 0, "token": tok})
}

func (t *table) seatPlayer(i int, name string) string {
	tok := randHex(16)
	t.Seats[i] = &seat{Name: name, Token: tok}
	return tok
}

func (s *store) join(w http.ResponseWriter, r *http.Request) {
	t := s.get(r.PathValue("id"))
	if t == nil {
		fail(w, 404, "table not found")
		return
	}
	var in struct{ Name string }
	_ = json.NewDecoder(r.Body).Decode(&in)
	t.mu.Lock()
	defer t.mu.Unlock()
	for i := range t.Seats {
		if t.Seats[i] == nil {
			tok := t.seatPlayer(i, cleanName(in.Name))
			t.maybeStart()
			writeJSON(w, 200, map[string]any{"table": t.ID, "seat": i, "token": tok})
			return
		}
	}
	fail(w, 409, "table is full")
}

func (t *table) full() bool {
	for _, s := range t.Seats {
		if s == nil {
			return false
		}
	}
	return true
}

func (t *table) maybeStart() {
	if t.Game == nil && t.full() {
		t.Game = belot.NewGame(3, 151)
		_ = t.Game.StartDeal(nil)
		t.runBots()
	}
}

func (t *table) auth(r *http.Request) (int, bool) {
	tok := r.URL.Query().Get("token")
	if tok == "" {
		tok = r.Header.Get("X-Token")
	}
	for i, s := range t.Seats {
		if s != nil && !s.Bot && tok != "" && s.Token == tok {
			return i, true
		}
	}
	return -1, false
}

// DEV ONLY: practice bots. The product spec forbids bots in real games; this endpoint exists so
// two testers can try the engine before matchmaking exists. Tables show a visible "practice" flag.
func (s *store) addBots(w http.ResponseWriter, r *http.Request) {
	t := s.get(r.PathValue("id"))
	if t == nil {
		fail(w, 404, "table not found")
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.auth(r); !ok {
		fail(w, 403, "not seated at this table")
		return
	}
	n := 0
	for i := range t.Seats {
		if t.Seats[i] == nil {
			t.Seats[i] = &seat{Name: "Bot " + string(rune('A'+n)), Token: randHex(8), Bot: true}
			n++
		}
	}
	t.maybeStart()
	writeJSON(w, 200, map[string]int{"added": n})
}

func (t *table) hasBots() bool {
	for _, s := range t.Seats {
		if s != nil && s.Bot {
			return true
		}
	}
	return false
}

// runBots lets bots act until it is a human's turn. Simple strategy: pass, then first legal card.
func (t *table) runBots() {
	g := t.Game
	for i := 0; i < 200; i++ {
		if g == nil || (g.Phase != belot.PhaseBidding && g.Phase != belot.PhasePlaying) {
			return
		}
		s := t.Seats[g.Turn]
		if s == nil || !s.Bot {
			return
		}
		if g.Phase == belot.PhaseBidding {
			_ = g.Bid(g.Turn, "pass", "")
		} else {
			legal, _ := g.LegalCards(g.Turn)
			if len(legal) == 0 {
				return
			}
			_ = g.Play(g.Turn, legal[0])
		}
	}
}

type view struct {
	Table     string            `json:"table"`
	Practice  bool              `json:"practice"`
	Ruleset   string            `json:"ruleset"`
	You       int               `json:"you"`
	Seats     []seatView        `json:"seats"`
	Started   bool              `json:"started"`
	Phase     belot.Phase       `json:"phase"`
	Dealer    int               `json:"dealer"`
	Turn      int               `json:"turn"`
	Contract  belot.Contract    `json:"contract"`
	Bidder    int               `json:"bidder"`
	Mult      int               `json:"mult"`
	Hand      []belot.Card      `json:"hand"`
	Legal     []belot.Card      `json:"legal"`
	LegalWhy  string            `json:"legalWhy"`
	Trick     []belot.Played    `json:"trick"`
	LastTrick []belot.Played    `json:"lastTrick"`
	Scores    [2]int            `json:"scores"`
	Tricks    [2]int            `json:"tricks"`
	Target    int               `json:"target"`
	Result    *belot.DealResult `json:"result"`
	Winner    int               `json:"winner"`
	Log       []string          `json:"log"`
	CanBid    []string          `json:"canBid"`
}

type seatView struct {
	Name  string `json:"name"`
	Bot   bool   `json:"bot"`
	Cards int    `json:"cards"`
}

func (t *table) viewFor(me int) view {
	v := view{Table: t.ID, Practice: t.hasBots(), Ruleset: belot.RulesetVersion, You: me, Winner: -1}
	for _, s := range t.Seats {
		sv := seatView{}
		if s != nil {
			sv.Name, sv.Bot = s.Name, s.Bot
		}
		v.Seats = append(v.Seats, sv)
	}
	g := t.Game
	if g == nil {
		return v
	}
	v.Started, v.Phase, v.Dealer, v.Turn = true, g.Phase, g.Dealer, g.Turn
	v.Contract, v.Bidder, v.Mult = g.Contract, g.Bidder, g.Mult
	v.Trick, v.LastTrick, v.Scores, v.Tricks, v.Target = g.Trick, g.LastTrick, g.Scores, g.Tricks, g.Target
	v.Result, v.Winner, v.Log = g.Result, g.Winner, g.Log
	v.Hand = g.Hands[me] // ONLY the requester's own hand is ever serialized
	for i := range v.Seats {
		v.Seats[i].Cards = len(g.Hands[i])
	}
	if g.Phase == belot.PhasePlaying && g.Turn == me {
		v.Legal, v.LegalWhy = g.LegalCards(me)
	}
	if g.Phase == belot.PhaseBidding && g.Turn == me {
		v.CanBid = []string{"pass", "bid"}
		if g.Contract != "" {
			if g.Mult == 1 && belot.Team(me) != belot.Team(g.Bidder) {
				v.CanBid = append(v.CanBid, "double")
			}
			if g.Mult == 2 && belot.Team(me) == belot.Team(g.Bidder) {
				v.CanBid = append(v.CanBid, "redouble")
			}
		}
	}
	return v
}

func (s *store) state(w http.ResponseWriter, r *http.Request) {
	t := s.get(r.PathValue("id"))
	if t == nil {
		fail(w, 404, "table not found")
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	me, ok := t.auth(r)
	if !ok {
		fail(w, 403, "bad token")
		return
	}
	writeJSON(w, 200, t.viewFor(me))
}

func (s *store) act(w http.ResponseWriter, r *http.Request, f func(t *table, me int) error) {
	t := s.get(r.PathValue("id"))
	if t == nil {
		fail(w, 404, "table not found")
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	me, ok := t.auth(r)
	if !ok {
		fail(w, 403, "bad token")
		return
	}
	if t.Game == nil {
		fail(w, 409, "waiting for four players")
		return
	}
	if err := f(t, me); err != nil {
		fail(w, 400, err.Error())
		return
	}
	t.runBots()
	writeJSON(w, 200, t.viewFor(me))
}

func (s *store) bid(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Action   string
		Contract belot.Contract
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		fail(w, 400, "bad json")
		return
	}
	s.act(w, r, func(t *table, me int) error { return t.Game.Bid(me, in.Action, in.Contract) })
}

func (s *store) play(w http.ResponseWriter, r *http.Request) {
	var in belot.Card
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		fail(w, 400, "bad json")
		return
	}
	s.act(w, r, func(t *table, me int) error { return t.Game.Play(me, in) })
}

func (s *store) next(w http.ResponseWriter, r *http.Request) {
	s.act(w, r, func(t *table, me int) error {
		if t.Game.Phase == belot.PhaseMatch {
			t.Game = belot.NewGame(3, 151)
			return t.Game.StartDeal(nil)
		}
		if t.Game.Phase != belot.PhaseDealOver {
			return errors.New("deal is not finished")
		}
		return t.Game.NextDeal(nil)
	})
}
