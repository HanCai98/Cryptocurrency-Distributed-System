package test

import (
	liteminer "liteminer/pkg"
	"testing"
)

func TestClient1(t *testing.T) {
	p, err := liteminer.CreatePool("1112")
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

	if client.GetPool().String() != client.Pool.Conn.RemoteAddr().String() {
		t.Errorf("Connected pool address doesn't match, error")
	}

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
	if client.Pool.Conn.LocalAddr().String() != p.GetClient().String() {
		t.Errorf("Client address doesn't match, error")
	}

	if err != nil {
		t.Errorf("Received error %v when mining", err)
	} else {
		if nonce != int64(target) {
			t.Errorf("Expected nonce %d, but received %d", target, nonce)
		}
	}

}

func TestClient2(t *testing.T) {
	p, err := liteminer.CreatePool("1113")
	if err != nil {
		t.Errorf("Received error %v when creating pool", err)
	}

	poolAddr := p.Addr.String()
	client := liteminer.CreateClient("")
	client.Connect(poolAddr)

	if client.GetPool().String() != client.Pool.Conn.RemoteAddr().String() {
		t.Errorf("Connected pool address doesn't match")
	}
}

func TestClient3(t *testing.T) {
	p, err := liteminer.CreatePool("1114")
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

	if client.GetPool().String() != client.Pool.Conn.RemoteAddr().String() {
		t.Errorf("Connected pool address doesn't match")
	}

	data := "bbbb"
	upperbound := uint64(1000)
	// when client disconnect with the pool
	client.Pool.Conn = nil
	nonce, err := client.Mine(data, upperbound)
	target := int64(-1)

	if nonce != target {
		t.Errorf("Expected nonce %d, but received %d", target, nonce)
	}
}

func TestClient4(t *testing.T) {
	p, err := liteminer.CreatePool("1115")
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

	client1 := liteminer.CreateClient(poolAddr)
	client1.Pool.Conn.Close()
	client2 := liteminer.CreateClient(poolAddr)

	data := "aaa"
	upperbound := uint64(100)
	nonce, err := client2.Mine(data, upperbound)
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
