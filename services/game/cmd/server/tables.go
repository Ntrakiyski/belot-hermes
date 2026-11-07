package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Ntrakiyski/belot-hermes/services/game/internal/belot"
)

type seat struct {
	Name  string
	Token string
	Bot   bool
}

type table struct {
	mu           sync.Mutex
	ID           string
	Seats        [4]*seat
	Game         *belot.Game
	Seq          int
	DealNumber   int
	CollectUntil time.Time
	NextAction   time.Time
	TrickWinner  int
	TrickPoints  int
	BidHistory   []bidEvent
	History      []dealSummary
}

type store struct {
	mu              sync.Mutex
	tables          map[string]*table
	practiceEnabled bool
}

func newStore() *store {
	return &store{tables: map[string]*table{}, practiceEnabled: os.Getenv("ENABLE_PRACTICE_BOTS") == "1"}
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *store) routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tables", s.create)
	mux.HandleFunc("POST /api/tables/{id}/join", s.join)
	if s.practiceEnabled { // DEMO only; never on MAIN
		mux.HandleFunc("POST /api/tables/{id}/bots", s.addBots)
	}
	mux.HandleFunc("GET /api/tables/{id}", s.state)
	mux.HandleFunc("POST /api/tables/{id}/bid", s.bid)
	mux.HandleFunc("POST /api/tables/{id}/play", s.play)
	mux.HandleFunc("POST /api/tables/{id}/next", s.next)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func fail(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func cleanName(n string) string {
	n = strings.TrimSpace(n)
	runes := []rune(n)
	if len(runes) > 20 {
		n = string(runes[:20])
	}
	if n == "" {
		n = "Играч"
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
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&in) != nil {
		fail(w, 400, "Невалидни данни.")
		return
	}
	t := &table{ID: strings.ToUpper(randHex(2)), TrickWinner: -1}
	tok := t.seatPlayer(0, cleanName(in.Name))
	s.mu.Lock()
	for s.tables[t.ID] != nil {
		t.ID = strings.ToUpper(randHex(2))
	}
	s.tables[t.ID] = t
	s.mu.Unlock()
	writeJSON(w, 200, map[string]any{"table": t.ID, "seat": 0, "token": tok})
}

func (t *table) seatPlayer(i int, name string) string {
	tok := randHex(16)
	t.Seats[i] = &seat{Name: name, Token: tok}
	t.Seq++
	return tok
}

func (s *store) join(w http.ResponseWriter, r *http.Request) {
	t := s.get(r.PathValue("id"))
	if t == nil {
		fail(w, 404, "Масата вече не съществува. Създай нова игра.")
		return
	}
	var in struct{ Name string }
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&in) != nil {
		fail(w, 400, "Невалидни данни.")
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for i := range t.Seats {
		if t.Seats[i] == nil {
			tok := t.seatPlayer(i, cleanName(in.Name))
			if err := t.maybeStart(); err != nil {
				fail(w, 500, "Неуспешно раздаване.")
				return
			}
			writeJSON(w, 200, map[string]any{"table": t.ID, "seat": i, "token": tok})
			return
		}
	}
	fail(w, 409, "Масата е пълна.")
}

func (t *table) full() bool {
	for _, s := range t.Seats {
		if s == nil {
			return false
		}
	}
	return true
}

func (t *table) maybeStart() error {
	if t.Game == nil && t.full() {
		g := belot.NewGame(3, 151)
		if err := g.StartDeal(nil); err != nil {
			return err
		}
		t.Game = g
		t.DealNumber = 1
		t.Seq++
		t.schedule(time.Now())
	}
	return nil
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
		fail(w, 404, "Масата вече не съществува. Създай нова игра.")
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	me, ok := t.auth(r)
	if !ok {
		fail(w, 403, "Нямаш място на тази маса.")
		return
	}
	n := 0
	for i := range t.Seats {
		if t.Seats[i] == nil {
			t.Seats[i] = &seat{Name: []string{"Борис", "Елена", "Виктор"}[n], Token: randHex(8), Bot: true}
			n++
		}
	}
	if n > 0 {
		t.Seq++
	}
	if err := t.maybeStart(); err != nil {
		fail(w, 500, "Неуспешно раздаване.")
		return
	}
	v := t.viewFor(me)
	v.PracticeEnabled = s.practiceEnabled
	writeJSON(w, 200, v)
}

func (t *table) hasBots() bool {
	for _, s := range t.Seats {
		if s != nil && s.Bot {
			return true
		}
	}
	return false
}

type view struct {
	Version            string                  `json:"version"`
	PracticeEnabled    bool                    `json:"practiceEnabled"`
	Seq                int                     `json:"seq"`
	DealNumber         int                     `json:"dealNumber"`
	Collecting         bool                    `json:"collecting"`
	TrickWinner        int                     `json:"trickWinner"`
	TrickPoints        int                     `json:"trickPoints"`
	NextActionAt       int64                   `json:"nextActionAt"`
	Captured           [2]int                  `json:"captured"`
	Pending            int                     `json:"pending"`
	BidHistory         []bidEvent              `json:"bidHistory"`
	History            []dealSummary           `json:"history"`
	DeclarationOptions []belot.Declaration     `json:"declarationOptions"`
	Declared           []belot.DeclaredSummary `json:"declared"`
	CanBelot           []belot.Card            `json:"canBelot"`
	BelotPoints        [2]int                  `json:"belotPoints"`
	Table              string                  `json:"table"`
	Practice           bool                    `json:"practice"`
	Ruleset            string                  `json:"ruleset"`
	You                int                     `json:"you"`
	Seats              []seatView              `json:"seats"`
	Started            bool                    `json:"started"`
	Phase              belot.Phase             `json:"phase"`
	Dealer             int                     `json:"dealer"`
	Turn               int                     `json:"turn"`
	Contract           belot.Contract          `json:"contract"`
	Bidder             int                     `json:"bidder"`
	Mult               int                     `json:"mult"`
	Hand               []belot.Card            `json:"hand"`
	Legal              []belot.Card            `json:"legal"`
	LegalWhy           string                  `json:"legalWhy"`
	Trick              []belot.Played          `json:"trick"`
	LastTrick          []belot.Played          `json:"lastTrick"`
	Scores             [2]int                  `json:"scores"`
	Tricks             [2]int                  `json:"tricks"`
	Target             int                     `json:"target"`
	Result             *belot.DealResult       `json:"result"`
	Winner             int                     `json:"winner"`
	Log                []string                `json:"log"`
	CanBid             []string                `json:"canBid"`
}

type seatView struct {
	Name  string `json:"name"`
	Bot   bool   `json:"bot"`
	Cards int    `json:"cards"`
}

func (t *table) viewFor(me int) view {
	v := view{Version: appVersion, Table: t.ID, Practice: t.hasBots(), Ruleset: belot.RulesetVersion, You: me, Winner: -1,
		Seq: t.Seq, DealNumber: t.DealNumber, Collecting: !t.CollectUntil.IsZero(), TrickWinner: t.TrickWinner, TrickPoints: t.TrickPoints,
		BidHistory: t.BidHistory, History: t.History}
	if !t.NextAction.IsZero() {
		v.NextActionAt = t.NextAction.UnixMilli()
	}
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
	v.Captured, v.Pending, v.BelotPoints = g.Captured, g.Pending, g.BelotPoints
	v.Declared = g.PublicDeclarations()
	v.DeclarationOptions = g.DeclarationOptions(me)
	v.Hand = g.Hands[me] // ONLY the requester's own hand is ever serialized
	for i := range v.Seats {
		v.Seats[i].Cards = len(g.Hands[i])
	}
	if !v.Collecting && g.Phase == belot.PhasePlaying && g.Turn == me {
		v.Legal, v.LegalWhy = g.LegalCards(me)
		for _, card := range v.Legal {
			if g.CanAnnounceBelot(me, card) {
				v.CanBelot = append(v.CanBelot, card)
			}
		}
	}
	if !v.Collecting && g.Phase == belot.PhaseBidding && g.Turn == me {
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
		fail(w, 404, "Масата вече не съществува. Създай нова игра.")
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	me, ok := t.auth(r)
	if !ok {
		fail(w, 403, "Сесията е изтекла. Влез отново.")
		return
	}
	v := t.viewFor(me)
	v.PracticeEnabled = s.practiceEnabled
	writeJSON(w, 200, v)
}

func (s *store) act(w http.ResponseWriter, r *http.Request, f func(t *table, me int) error) {
	t := s.get(r.PathValue("id"))
	if t == nil {
		fail(w, 404, "Масата вече не съществува. Създай нова игра.")
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	me, ok := t.auth(r)
	if !ok {
		fail(w, 403, "Сесията е изтекла. Влез отново.")
		return
	}
	if t.Game == nil {
		fail(w, 409, "Очакваме четирима играчи.")
		return
	}
	if !t.CollectUntil.IsZero() {
		fail(w, 409, "Изчакай прибирането на взятката.")
		return
	}
	beforeTricks := t.Game.Tricks[0] + t.Game.Tricks[1]
	beforePoints := t.Game.Captured[0] + t.Game.Captured[1]
	if err := f(t, me); err != nil {
		fail(w, 400, err.Error())
		return
	}
	t.progress(time.Now(), beforeTricks, beforePoints)
	v := t.viewFor(me)
	v.PracticeEnabled = s.practiceEnabled
	writeJSON(w, 200, v)
}

func (s *store) bid(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Action   string
		Contract belot.Contract
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&in) != nil {
		fail(w, 400, "Невалидни данни.")
		return
	}
	s.act(w, r, func(t *table, me int) error {
		if err := t.Game.Bid(me, in.Action, in.Contract); err != nil {
			return err
		}
		t.BidHistory = append(t.BidHistory, bidEvent{Seat: me, Action: in.Action, Contract: in.Contract})
		return nil
	})
}

func (s *store) play(w http.ResponseWriter, r *http.Request) {
	var in struct {
		belot.Card
		Declarations []string `json:"declarations"`
		Belot        bool     `json:"belot"`
	}
	if json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&in) != nil {
		fail(w, 400, "Невалидни данни.")
		return
	}
	s.act(w, r, func(t *table, me int) error { return t.Game.PlayDeclared(me, in.Card, in.Declarations, in.Belot) })
}

func (s *store) next(w http.ResponseWriter, r *http.Request) {
	s.act(w, r, func(t *table, me int) error {
		if t.Game.Phase == belot.PhaseMatch {
			game := belot.NewGame(3, 151)
			if err := game.StartDeal(nil); err != nil {
				return err
			}
			t.Game = game
			t.DealNumber = 1
			t.History = nil
			t.BidHistory = nil
			t.TrickWinner = -1
			t.TrickPoints = 0
			return nil
		}
		if t.Game.Phase != belot.PhaseDealOver {
			return errors.New("Раздаването още не е завършило.")
		}
		if err := t.Game.NextDeal(nil); err != nil {
			return err
		}
		t.DealNumber++
		t.BidHistory = nil
		t.TrickWinner = -1
		t.TrickPoints = 0
		return nil
	})
}
