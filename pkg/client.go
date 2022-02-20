/*
 *  Brown University, CS138, Spring 2022
 *
 *  Purpose: a LiteMiner client.
 */

package pkg

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Client represents a LiteMiner client
type Client struct {
	Pool     MiningConn // Currently connected pool
	PoolMtx  sync.Mutex // Mutex for connected pool
	TxResult chan int64 // Used to send results of transaction
}

// CreateClient creates a new client connected to the given pool address.
func CreateClient(addr string) *Client {
	c := &Client{
		TxResult: make(chan int64),
	}

	c.Connect(addr)

	return c
}

// Connect connects the client to the specified pool addresses.
func (c *Client) Connect(addr string) {
	conn, err := ClientConnect(addr)
	if err != nil {
		Err.Printf("Received error %v when connecting to pool %v\n", err, addr)
		return
	}

	c.PoolMtx.Lock()
	c.Pool = conn
	c.PoolMtx.Unlock()

	go c.processPool(conn)
}

// processPool handles incoming messages from the pool represented by conn.
func (c *Client) processPool(conn MiningConn) {
	for {
		msg, err := RecvMsg(conn)
		if err != nil {
			if _, ok := err.(*net.OpError); ok || err == io.EOF {
				Err.Printf("Lost connection to pool %v\n", conn.Conn.RemoteAddr())

				c.PoolMtx.Lock()
				c.Pool.Conn = nil
				c.PoolMtx.Unlock()

				c.TxResult <- -1  // -1 used to indicate error
				conn.Conn.Close() // Close the connection

				return
			}

			Err.Printf(
				"Received error %v when processing pool %v\n",
				err,
				conn.Conn.RemoteAddr(),
			)

			c.TxResult <- -1 // -1 used to indicate error

			continue
		}

		switch msg.Type {
		case BusyPool:
			Out.Printf("Pool %v is currently busy, disconnecting\n", conn.Conn.RemoteAddr())

			c.PoolMtx.Lock()
			c.Pool.Conn = nil
			c.PoolMtx.Unlock()

			conn.Conn.Close() // Close the connection

			return
		case ProofOfWork:
			Debug.Printf("Pool %v found nonce %v\n", conn.Conn.RemoteAddr(), msg.Nonce)

			c.TxResult <- int64(msg.Nonce)
		default:
			Err.Printf(
				"Received unexpected message of type %v from pool %v\n",
				msg.Type,
				conn.Conn.RemoteAddr(),
			)

			c.TxResult <- -1 // -1 used to indicate error
		}
	}
}

// Mine is given a transaction encoded as a string and an unsigned integer and returns
// the nonce calculated by the connected pool. This method should NOT be
// executed concurrently by the same client.
func (c *Client) Mine(data string, upperBound uint64) (int64, error) {
	c.PoolMtx.Lock()

	if c.Pool.Conn == nil {
		return -1, fmt.Errorf("not connected to any pools")
	}

	// Send transaction to connected pool(s)
	tx := TransactionMsg(data, upperBound)
	SendMsg(c.Pool, tx)
	c.PoolMtx.Unlock()

	nonce := <-c.TxResult

	return nonce, nil
}

// GetPool returns the address of the current pool or nil if there is no
// current pool.
func (p *Client) GetPool() net.Addr {
	p.PoolMtx.Lock()
	defer p.PoolMtx.Unlock()

	if p.Pool.Conn == nil {
		return nil
	}
	return p.Pool.Conn.RemoteAddr()
}
