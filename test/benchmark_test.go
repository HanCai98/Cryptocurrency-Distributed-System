package test

import (
	liteminer "liteminer/pkg"
	"testing"
)

func BenchmarkTest1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p, err := liteminer.CreatePool("")
		if err != nil {
			b.Errorf("Received error %v when creating pool", err)
		}

		poolAddr := p.Addr.String()
		numMiners := 10
		miners := make([]*liteminer.Miner, numMiners)

		for i := 0; i < numMiners; i++ {
			m, err := liteminer.CreateMiner(poolAddr)
			if err != nil {
				b.Errorf("Received error %v when creating miner", err)
			}
			miners[i] = m
		}

		client := liteminer.CreateClient(poolAddr)

		data := "aaaa"
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
			b.Errorf("Received error %v when mining", err)
		} else {
			if nonce != int64(target) {
				b.Errorf("Expected nonce %d, but received %d", target, nonce)
			}
		}
	}
}
