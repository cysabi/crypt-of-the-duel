package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/olahol/melody"
)

type Game struct {
	mu        sync.Mutex
	turnCount int
	// actedThisBeat map[PlayerId]bool
	actions map[PlayerId]PayloadAction // Accumulated actions for the current beat
}

func (g Game) New() Game {
	g.turnCount = 0
	// g.actedThisBeat = make(map[PlayerId]bool)
	g.actions = make(map[PlayerId]PayloadAction)
	return g
}

func (g *Game) ReceiveAction(pid PlayerId, action PayloadAction) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if action.turnCount != g.turnCount {
		// Ignore actions not for the current beat
		return
	}
	// Only accept the first action per player per beat
	if _, exists := g.actions[pid]; !exists {
		g.actions[pid] = action
	}
}

func (g *Game) ProcessBeat(m *melody.Melody) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Get all player IDs
	sessions, err := m.Sessions()
	if err != nil {
		panic(err)
	}

	for _, session := range sessions {
		pid := session.MustGet("pid").(PlayerId)
		if _, acted := g.actions[pid]; !acted {
			// Insert a skip action for players who didn't act
			g.actions[pid] = PayloadAction{
				pid:       pid,
				turnCount: g.turnCount,
				action:    "skip",
			}
		}
	}

	// Broadcast all actions for this beat
	for _, action := range g.actions {
		m.Broadcast([]byte(action.String()))
	}

	// Prepare for next beat
	g.actions = make(map[PlayerId]PayloadAction)
	g.turnCount++
}

func (g *Game) BroadcastMissing(m *melody.Melody) {
	sessions, err := m.Sessions()
	if err != nil {
		panic(err)
	}

	for _, session := range sessions {
		pid := session.MustGet("pid").(PlayerId)
		if !g.actedThisBeat[pid] {
			skipPayload := PayloadAction{
				pid:       pid,
				turnCount: g.turnCount,
				action:    "skip",
			}.String()
			m.Broadcast([]byte(skipPayload))
		}
	}
}

func GetOthers(m *melody.Melody, you PlayerId) string {
	sessions, err := m.Sessions()
	if err != nil {
		panic(err)
	}

	pids := []string{}
	for _, s := range sessions {
		pid := s.MustGet("pid").(PlayerId)
		if pid != you {
			pids = append(pids, string(pid))
		}
	}
	return strings.Join(pids, ",")
}

func main() {
	m := melody.New()
	game := Game{}.New()

	go func() {
		beatDuration := 200 * time.Millisecond // Set your beat duration here
		for {
			time.Sleep(beatDuration)
			game.ProcessBeat(m)
		}
	}()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleConnect(func(s *melody.Session) {
		fmt.Println("Got a connection")
		pid := GeneratePlayerId()
		s.Set("pid", pid)

		others := GetOthers(m, pid)
		payloadYou := PayloadYou{}.New(pid, []string{others})
		payloadThem := PayloadThem{}.New(pid, []string{})

		s.Write([]byte(payloadYou.String()))
		m.BroadcastOthers([]byte(payloadThem.String()), s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		pid := s.MustGet("pid").(PlayerId)
		payload := recievePayload(pid, string(msg))
		switch p := payload.(type) {
		case PayloadAction:
			game.ReceiveAction(pid, p)
		case PayloadStart:
			m.Broadcast([]byte(p.String()))
		}
	})

	const PORT = ":5174"
	log.Println("Server started on localhost" + PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
