package test

import (
	liteminer "liteminer/pkg"
	"testing"
)

func TestBasic(t *testing.T) {
	p, err := liteminer.CreatePool("1111")
	if err != nil {
		t.Errorf("Received error %v when creating pool", err)
	}

	addr := p.Addr.String()

	numMiners := 2
	miners := make([]*liteminer.Miner, numMiners)
	for i := 0; i < numMiners; i++ {
		m, err := liteminer.CreateMiner(addr)
		if err != nil {
			t.Errorf("Received error %v when creating miner", err)
		}
		miners[i] = m
	}

	client := liteminer.CreateClient(addr)

	data := "data"
	upperbound := uint64(1)
	nonce, err := client.Mine(data, upperbound)
	if err != nil {
		t.Errorf("Received error %v when mining", err)
	} else {
		expected := int64(1)
		if nonce != expected {
			t.Errorf("Expected nonce %d, but received %d", expected, nonce)
		}
	}
}
