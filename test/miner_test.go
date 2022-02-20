package test

import (
	liteminer "liteminer/pkg"
	"testing"
)

func TestMiner1(t *testing.T) {
	p, err := liteminer.CreatePool("1116")
	if err != nil {
		t.Errorf("Received error %v when creating pool", err)
	}

	poolAddr := p.Addr.String()
	numMiners := 5
	miners := make([]*liteminer.Miner, numMiners)

	for i := 0; i < numMiners; i++ {
		m, err := liteminer.CreateMiner(poolAddr)
		if err != nil {
			t.Errorf("Received error %v when creating miner", err)
		}
		miners[i] = m
	}

	client := liteminer.CreateClient(poolAddr)

	data := "aaa"
	upperbound := uint64(100)
	nonce, err := client.Mine(data, upperbound)
	target := uint64(0)
	minHash := liteminer.Hash(data, uint64(0))

	for i := uint64(1); i < upperbound; i++ {
		curHash := liteminer.Hash(data, i)
		if curHash < minHash {
			minHash = curHash
			target = i
		}
	}

	if len(p.GetMiners()) != numMiners {
		t.Errorf("Number of miners doesn't match, error!")
	}

	if err != nil {
		t.Errorf("Received error %v when mining", err)
	} else {
		if nonce != int64(target) {
			t.Errorf("Expected nonce %d, but received %d", target, nonce)
		}
	}
}

func TestMiner2(t *testing.T) {
	p, err := liteminer.CreatePool("1117")
	if err != nil {
		t.Errorf("Received error %v when creating pool", err)
	}

	poolAddr := p.Addr.String()
	numMiners := 5
	miners := make([]*liteminer.Miner, numMiners)

	for i := 0; i < numMiners; i++ {
		m, err := liteminer.CreateMiner(poolAddr)
		if err != nil {
			t.Errorf("Received error %v when creating miner", err)
		}
		miners[i] = m
	}

	shutdown := miners[0].IsShutdown

	if shutdown.Load() != false {
		t.Errorf("Miner 1 should not be shoutdown but is shutdown now, error!")
	}

	miners[0].Shutdown()
	shutdown = miners[0].IsShutdown
	if shutdown.Load() != true {
		t.Errorf("Miner 1 should be shoutdown but is not shutdown now, error!")
	}
}

func TestMiner3(t *testing.T) {
	p, err := liteminer.CreatePool("1118")
	if err != nil {
		t.Errorf("Received error %v when creating pool", err)
	}

	poolAddr := p.Addr.String()
	numMiners := 5
	miners := make([]*liteminer.Miner, numMiners)

	for i := 0; i < numMiners; i++ {
		m, err := liteminer.CreateMiner(poolAddr)
		if err != nil {
			t.Errorf("Received error %v when creating miner", err)
		}
		miners[i] = m
	}

	miners[0].Shutdown()
	client := liteminer.CreateClient(poolAddr)

	data := "aaa"
	upperbound := uint64(100)
	nonce, err := client.Mine(data, upperbound)
	target := uint64(0)
	minHash := liteminer.Hash(data, uint64(0))

	for i := uint64(1); i < upperbound; i++ {
		curHash := liteminer.Hash(data, i)
		if curHash < minHash {
			minHash = curHash
			target = i
		}
	}

	if err != nil {
		t.Errorf("Received error %v when mining", err)
	} else {
		if nonce != int64(target) {
			t.Errorf("Expected nonce %d, but received %d", target, nonce)
		}
	}
}
