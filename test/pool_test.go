package test

import (
	liteminer "liteminer/pkg"
	"testing"
)

func TestPool1(t *testing.T) {
	_, err1 := liteminer.CreatePool("2222")
	if err1 != nil {
		t.Errorf("Received error %v when creating pool", err1)
	}

	_, err2 := liteminer.CreatePool("2222")
	if err2 == nil {
		t.Errorf("Should receive error for duplicate port use")
	}
}

func TestPool2(t *testing.T) {
	p, err := liteminer.CreatePool("2223")
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
	p = nil

	liteminer.SendMsg(client.Pool, liteminer.ErrorMsg("Error"))

	data := "aaa"
	upperbound := uint64(100)
	nonce, err := client.Mine(data, upperbound)

	if err != nil {
		t.Errorf("Receive error %v when mining", err)
	}

	if nonce != int64(-1) {
		t.Errorf("Expected nonce %d, but received %d", int64(-1), nonce)
	}
}

func TestPool3(t *testing.T) {
	_, err := liteminer.CreatePool("err")
	if err == nil {
		t.Errorf("LiteMiner: Expect error, but get none when creating pool")
	}
}

func TestPool4(t *testing.T) {
	p, err := liteminer.CreatePool("2224")
	if err != nil {
		t.Errorf("LiteMiner: Receive error %v when creating pool", err)
	}

	addr := p.Addr.String()
	numMiners := 20
	for i := 0; i < numMiners; i++ {
		_, err := liteminer.CreateMiner(addr + "fake")
		if err == nil {
			t.Errorf("LiteMiner: Expect error for wrong pool connection, but get null")
		}
	}
}
