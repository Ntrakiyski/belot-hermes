package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/Ntrakiyski/belot-hermes/services/game/internal/belot"
)

func tableRequest(t *testing.T, mux *http.ServeMux, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("POST", path, bytes.NewReader(raw))
	r.Header.Set("X-Token", token)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w
}

func TestPracticeBidDoesNotSkipThreeBotTurns(t *testing.T) {
	t.Setenv("ENABLE_PRACTICE_BOTS", "1")
	s := newStore()
	tbl := &table{ID: "TEST"}
	token := tbl.seatPlayer(0, "Tester")
	s.tables[tbl.ID] = tbl
	mux := http.NewServeMux()
	s.routes(mux)
	if w := tableRequest(t, mux, "/api/tables/TEST/bots", token, map[string]any{}); w.Code != 200 {
		t.Fatal(w.Body.String())
	}
	w := tableRequest(t, mux, "/api/tables/TEST/bid", token, map[string]string{"action": "bid", "contract": "clubs"})
	var v view
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || v.Phase != belot.PhaseBidding || v.Turn != 1 {
		t.Fatalf("human must see first bot's bidding turn, status=%d phase=%s turn=%d", w.Code, v.Phase, v.Turn)
	}
}

func TestCompletedTrickRemainsVisibleAndBlocksFifthCard(t *testing.T) {
	s := newStore()
	tbl := &table{ID: "TEST"}
	tokens := [4]string{}
	for i := range tokens {
		tokens[i] = tbl.seatPlayer(i, "Tester")
	}
	tbl.maybeStart()
	s.tables[tbl.ID] = tbl
	_ = tbl.Game.Bid(0, "bid", belot.Clubs)
	for i := 1; i < 4; i++ {
		_ = tbl.Game.Bid(i, "pass", "")
	}
	mux := http.NewServeMux()
	s.routes(mux)
	var last *httptest.ResponseRecorder
	for i := 0; i < 4; i++ {
		legal, _ := tbl.Game.LegalCards(tbl.Game.Turn)
		last = tableRequest(t, mux, "/api/tables/TEST/play", tokens[tbl.Game.Turn], legal[0])
		if last.Code != 200 {
			t.Fatal(last.Body.String())
		}
	}
	var v map[string]any
	if err := json.Unmarshal(last.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if v["collecting"] != true {
		t.Fatal("completed four-card trick must remain in collection state")
	}
	cards, ok := v["lastTrick"].([]any)
	if !ok || len(cards) != 4 {
		t.Fatal("all four completed cards must be visible")
	}
	legal, _ := tbl.Game.LegalCards(tbl.Game.Turn)
	before := tbl.Game.Seq
	w := tableRequest(t, mux, "/api/tables/TEST/play", tokens[tbl.Game.Turn], legal[0])
	if w.Code != 409 || tbl.Game.Seq != before {
		t.Fatal("fifth play during collection must be rejected without mutation")
	}
}

func TestBotSchedulerAdvancesOnlyOneActionPerDeadline(t *testing.T) {
	tbl := &table{ID: "TEST"}
	for i := 0; i < 4; i++ {
		tbl.seatPlayer(i, "Tester")
		tbl.Seats[i].Bot = i != 0
	}
	if err := tbl.maybeStart(); err != nil {
		t.Fatal(err)
	}
	if err := tbl.Game.Bid(0, "bid", belot.AllTrump); err != nil {
		t.Fatal(err)
	}
	now := time.Unix(1000, 0)
	tbl.schedule(now)
	before := tbl.Game.Seq
	tbl.advance(now.Add(botBidDelay - time.Millisecond))
	if tbl.Game.Seq != before {
		t.Fatal("bot acted before deadline")
	}
	tbl.advance(now.Add(botBidDelay))
	if tbl.Game.Turn != 2 || tbl.Game.Seq != before+1 {
		t.Fatal("deadline must advance exactly one bot")
	}
	tbl.advance(now.Add(botBidDelay))
	if tbl.Game.Turn != 2 {
		t.Fatal("same timestamp skipped another bot")
	}
}

func TestPacedPracticeMatchesCannotStall(t *testing.T) {
	// Virtual time exercises actual bot/declaration choices without slowing the test suite.
	for n := 0; n < 40; n++ {
		tbl := &table{ID: "TEST"}
		for i := 0; i < 4; i++ {
			tbl.seatPlayer(i, "Tester")
			tbl.Seats[i].Bot = i != 0
		}
		tbl.Game = belot.NewGame(3, 151)
		rng := rand.New(rand.NewSource(int64(n + 1)))
		if err := tbl.Game.StartDeal(rng); err != nil {
			t.Fatal(err)
		}
		tbl.DealNumber = 1
		now := time.Unix(1000, 0)
		complete := false
		for step := 0; step < 20000; step++ {
			g := tbl.Game
			if !tbl.CollectUntil.IsZero() {
				now = tbl.CollectUntil
				tbl.advance(now)
				continue
			}
			if g.Phase == belot.PhaseMatch {
				complete = true
				break
			}
			if !tbl.NextAction.IsZero() {
				now = tbl.NextAction
				tbl.advance(now)
				continue
			}
			beforeTricks, beforePoints := g.Tricks[0]+g.Tricks[1], g.Captured[0]+g.Captured[1]
			if g.Phase == belot.PhaseDealOver {
				if err := g.NextDeal(rng); err != nil {
					t.Fatal(err)
				}
				tbl.DealNumber++
				tbl.BidHistory = nil
			} else {
				if g.Turn != 0 {
					t.Fatalf("seed=%d bot stalled phase=%s turn=%d seq=%d", n, g.Phase, g.Turn, g.Seq)
				}
				if g.Phase == belot.PhaseBidding {
					action, c := "pass", belot.Contract("")
					if g.Contract == "" {
						action, c = "bid", belot.AllTrump
					}
					if err := g.Bid(0, action, c); err != nil {
						t.Fatal(err)
					}
				} else {
					legal, _ := g.LegalCards(0)
					if len(legal) == 0 {
						t.Fatal("no human legal move")
					}
					if err := g.PlayDeclared(0, legal[0], nil, g.CanAnnounceBelot(0, legal[0])); err != nil {
						t.Fatal(err)
					}
				}
			}
			tbl.progress(now, beforeTricks, beforePoints)
		}
		if !complete {
			t.Fatalf("seed=%d practice match did not finish", n)
		}
	}
}

func TestInvalidTableInputAndBulgarianName(t *testing.T) {
	s := newStore()
	mux := http.NewServeMux()
	s.routes(mux)
	r := httptest.NewRequest("POST", "/api/tables", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != 400 || len(s.tables) != 0 {
		t.Fatal("malformed create must not create a table")
	}
	name := "Български играч с дълго име"
	clean := cleanName(name)
	if !utf8.ValidString(clean) || len([]rune(clean)) != 20 {
		t.Fatalf("invalid Unicode name %q", clean)
	}
}

func TestSeatProjectionOnlyExposesOwnDeclarationCards(t *testing.T) {
	s := newStore()
	tbl := &table{ID: "TEST"}
	tokens := [4]string{}
	for i := range tokens {
		tokens[i] = tbl.seatPlayer(i, "Tester")
	}
	tbl.Game = belot.NewGame(3, 151)
	_ = tbl.Game.StartDealWith(belot.NewDeck())
	_ = tbl.Game.Bid(0, "bid", belot.Hearts)
	for i := 1; i < 4; i++ {
		_ = tbl.Game.Bid(i, "pass", "")
	}
	s.tables[tbl.ID] = tbl
	mux := http.NewServeMux()
	s.routes(mux)
	for me := 0; me < 4; me++ {
		r := httptest.NewRequest("GET", "/api/tables/TEST", nil)
		r.Header.Set("X-Token", tokens[me])
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("private response must not be cached")
		}
		var v view
		if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
			t.Fatal(err)
		}
		for _, d := range v.DeclarationOptions {
			for _, card := range d.Cards {
				own := false
				for _, c := range v.Hand {
					if c == card {
						own = true
					}
				}
				if !own {
					t.Fatal("declaration leaks another seat's card")
				}
			}
		}
		for _, secret := range tokens {
			if bytes.Contains(w.Body.Bytes(), []byte(secret)) {
				t.Fatal("seat token leaked")
			}
		}
	}
}
