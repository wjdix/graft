package graft

import (
	"github.com/benmills/quiz"
	"github.com/wjdix/tiktok"
	"testing"
	"time"
)

type SpyServer struct {
	electionStarted bool
	electionCount   int
}

func FakeTicker(d time.Duration) Tickable {
	return tiktok.NewTicker(d)
}

func (server *SpyServer) StartElection() {
	server.electionStarted = true
	server.electionCount += 1
}

func TestTimerTellsServerToStartElectionWhenReceivingOnTimeoutChannel(t *testing.T) {
	test := quiz.Test(t)

	spyServer := &SpyServer{electionStarted: false, electionCount: 0}
	timer := NewElectionTimer(1, spyServer)
	timer.ElectionChannel <- 1
	test.Expect(spyServer.electionStarted).ToBeTrue()
	timer.ShutDown()
}

func TestTimer(t *testing.T) {
	test := quiz.Test(t)

	spyServer := &SpyServer{electionStarted: false, electionCount: 0}
	timer := NewElectionTimer(1, spyServer)
	timer.tickerBuilder = FakeTicker
	defer tiktok.ClearTickers()

	timer.StartTimer()

	tiktok.Tick(1)

	timer.ShutDown()
	test.Expect(spyServer.electionStarted).ToBeTrue()
}

func TestTimerStartsMultipleElections(t *testing.T) {
	test := quiz.Test(t)

	spyServer := &SpyServer{electionStarted: false, electionCount: 0}
	timer := NewElectionTimer(2, spyServer)
	timer.tickerBuilder = FakeTicker
	defer tiktok.ClearTickers()

	timer.StartTimer()

	tiktok.Tick(10)

	timer.ShutDown()

	test.Expect(spyServer.electionCount).ToBeGreaterThan(1)
}

func TestResetTimerDoesNotTick(t *testing.T) {
	test := quiz.Test(t)

	spyServer := &SpyServer{electionStarted: false, electionCount: 0}
	timer := NewElectionTimer(5, spyServer)
	timer.tickerBuilder = FakeTicker
	defer tiktok.ClearTickers()

	timer.StartTimer()

	tiktok.Tick(3)
	timer.Reset()

	tiktok.Tick(2)
	timer.ShutDown()
	test.Expect(spyServer.electionCount).ToEqual(0)
}
