/*
 *  Brown University, CS138, Spring 2022
 *
 *  Purpose: a LiteMiner mining pool.
 */

package pkg

import (
	"encoding/gob"
	"io"
	"math"
	"net"
	"sync"
	"time"

	"go.uber.org/atomic"
)

// HeartbeatTimeout is the time duration at which a pool considers a miner 'dead'
const HeartbeatTimeout = 3 * HeartbeatFreq
const IntervalRatio = 5

// Pool represents a LiteMiner mining pool
type Pool struct {
	Addr      net.Addr
	Miners    map[net.Addr]MiningConn // Currently connected miners
	MinersMtx sync.Mutex              // Mutex for concurrent access to Miners
	Client    MiningConn              // The current client
	ClientMtx sync.Mutex              // Mutex for concurrent access to Client
	busy      *atomic.Bool            // True when processing a transaction

	Nonces            chan Message
	Intervals         chan Interval
	data              *atomic.String
	totalIntervals    *atomic.Uint64
	completeIntervals *atomic.Uint64
}

// CreatePool creates a new pool at the specified port.
func CreatePool(port string) (*Pool, error) {
	p := &Pool{
		busy:              atomic.NewBool(false),
		data:              atomic.NewString(""),
		completeIntervals: atomic.NewUint64(0),
		Miners:            make(map[net.Addr]MiningConn),
		Intervals:         make(chan Interval),
	}

	// TODO: Students should (if necessary) initialize any additional members
	// to the Pool struct here.
	// p.Channel = make(chan Interval)
	err := p.startListener(port)
	return p, err
}

// startListener starts listening for new connections.
func (p *Pool) startListener(port string) error {
	listener, portID, err := OpenListener(port)
	if err != nil {
		return err
	}

	Out.Printf("Listening on port %v\n", portID)

	p.Addr = listener.Addr()

	// Listen for and accept connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				Err.Printf("Received error %v when listening for connections\n", err)
				continue
			}

			go p.handleConnection(conn)
		}
	}()

	return nil
}

// handleConnection handles an incoming connection and delegates to
// handleMinerConnection or handleClientConnection.
func (p *Pool) handleConnection(nc net.Conn) {
	// Set up connection
	conn := MiningConn{}
	conn.Conn = nc
	conn.Enc = gob.NewEncoder(nc)
	conn.Dec = gob.NewDecoder(nc)

	// Wait for Hello message
	msg, err := RecvMsg(conn)
	if err != nil {
		Err.Printf(
			"Received error %v when processing Hello message from %v\n",
			err,
			conn.Conn.RemoteAddr(),
		)
		conn.Conn.Close() // Close the connection
		return
	}

	switch msg.Type {
	case MinerHello:
		p.handleMinerConnection(conn)
	case ClientHello:
		p.handleClientConnection(conn)
	default:
		Err.Printf("Pool received unexpected message type %v (msg=%v)", msg.Type, msg)
		SendMsg(conn, ErrorMsg("Unexpected message type"))
	}
}

// handleClientConnection handles a connection from a client.
func (p *Pool) handleClientConnection(conn MiningConn) {
	Debug.Printf("Received client connection from %v", conn.Conn.RemoteAddr())

	p.ClientMtx.Lock()
	if p.Client.Conn != nil {
		Debug.Printf(
			"Busy with client %v, sending BusyPool message to client %v",
			p.Client.Conn.RemoteAddr(),
			conn.Conn.RemoteAddr(),
		)
		SendMsg(conn, BusyPoolMsg())
		p.ClientMtx.Unlock()
		return
	}
	p.ClientMtx.Unlock()

	if p.busy.Load() {
		Debug.Printf(
			"Busy with previous transaction, sending BusyPool message to client %v",
			conn.Conn.RemoteAddr(),
		)
		SendMsg(conn, BusyPoolMsg())
		return
	}
	p.ClientMtx.Lock()
	p.Client = conn
	p.ClientMtx.Unlock()

	// Listen for and handle incoming messages
	for {
		msg, err := RecvMsg(conn)
		if err != nil {
			if _, ok := err.(*net.OpError); ok || err == io.EOF {
				Out.Printf("Client %v disconnected\n", conn.Conn.RemoteAddr())

				conn.Conn.Close() // Close the connection

				p.ClientMtx.Lock()
				p.Client.Conn = nil
				p.ClientMtx.Unlock()

				return
			}
			Err.Printf(
				"Received error %v when processing message from client %v\n",
				err,
				conn.Conn.RemoteAddr(),
			)
			return
		}

		if msg.Type != Transaction {
			SendMsg(conn, ErrorMsg("Expected Transaction message"))
			continue
		}

		Debug.Printf(
			"Received transaction from client %v with data %v and upper bound %v",
			conn.Conn.RemoteAddr(),
			msg.Data,
			msg.Upper,
		)

		p.MinersMtx.Lock()
		if len(p.Miners) == 0 {
			SendMsg(conn, ErrorMsg("No miners connected"))
			p.MinersMtx.Unlock()
			continue
		}
		p.MinersMtx.Unlock()

		// TODO: Students should handle an incoming transaction from a client. A
		// pool may process one transaction at a time – thus, if you receive
		// another transaction while busy, you should send a BusyPool message.
		// Otherwise, you should let the miners do their jobs. Note that miners
		// are handled in separate go routines (`handleMinerConnection`). To notify
		// the miners, consider using a shared data structure.

		if p.busy.Load() {
			SendMsg(conn, BusyPoolMsg())
		} else {
			p.busy.Store(true)
			if msg.Upper > 0 {
				p.Nonces = make(chan Message)
				//Debug.Printf("Processing transaction from %v", conn.Conn.RemoteAddr())
				p.data.Store(msg.Data)
				p.MinersMtx.Lock()
				n := len(p.Miners)
				p.MinersMtx.Unlock()
				// if number of miners * 5 < upper, numIntervals = number of miners * 5
				// if number of miners * 5 > upper, numIntervals = upper
				numIntervals := int(math.Min(float64(n*IntervalRatio), float64(msg.Upper)))
				intervals := GenerateIntervals(msg.Upper, numIntervals)
				p.totalIntervals = atomic.NewUint64(uint64(numIntervals))

				go func() {
					for _, interval := range intervals {
						p.Intervals <- interval
					}
				}()

				go p.findSmallestNonce(msg.Data)
			} else {
				SendMsg(p.Client, ProofOfWorkMsg(p.data.Load(), 0, Hash(msg.Data, 0)))
				p.busy.Store(false)
			}

		}
	}
}

func (p *Pool) findSmallestNonce(data string) {
	firstNonce := true
	var sendNonce uint64
	var sendHash uint64
	for msg := range p.Nonces {
		if firstNonce {
			sendNonce = msg.Nonce
			sendHash = msg.Hash
			firstNonce = false
		} else {
			if msg.Hash < sendHash {
				sendNonce = msg.Nonce
				sendHash = msg.Hash
			}
		}
	}
	SendMsg(p.Client, ProofOfWorkMsg(data, sendNonce, sendHash))
	p.busy.Store(false)
}

// handleMinerConnection handles a connection from a miner.
func (p *Pool) handleMinerConnection(conn MiningConn) {
	Debug.Printf("Received miner connection from %v", conn.Conn.RemoteAddr())

	p.MinersMtx.Lock()
	p.Miners[conn.Conn.RemoteAddr()] = conn
	p.MinersMtx.Unlock()

	msgChan := make(chan Message)
	go p.receiveFromMiner(conn, msgChan)

	// TODO: Students should handle a miner connection. If a miner does not
	// send a StatusUpdate message every `HeartbeatTimeout` while mining,
	// any work assigned to them should be redistributed and they should be
	// disconnected and removed from `p.Miners`.
	// For maintaining a queue of jobs yet to be taken, consider using a go channel.

	for interval := range p.Intervals {
		timer := time.After(HeartbeatTimeout)
		message := MineRequestMsg(p.data.Load(), interval.Lower, interval.Upper)
		SendMsg(conn, message)
		isMining := true
		for isMining {
			select {
			case msg := <-msgChan:
				switch msg.Type {
				case ProofOfWork:
					p.Nonces <- msg
					p.completeIntervals.Add(1)
					if p.totalIntervals.Load() == p.completeIntervals.Load() {
						p.completeIntervals.Store(0)
						p.totalIntervals.Store(0)
						close(p.Nonces)
					}
					isMining = false
					break

				case StatusUpdate:
					break

				default:
					break
				}
				timer = time.After(HeartbeatTimeout)
				break

			case <-timer:
				Debug.Printf("Find a slow miner AT: %v, reassigned job.", conn.Conn.RemoteAddr())
				delete(p.Miners, conn.Conn.RemoteAddr())
				p.Intervals <- interval
				return

			default:
				break
			}
		}
	}
}

// receiveFromMiner waits for messages from the miner specified by conn and
// forwards them over msgChan.
func (p *Pool) receiveFromMiner(conn MiningConn, msgChan chan Message) {
	for {
		msg, err := RecvMsg(conn)
		if err != nil {
			if _, ok := err.(*net.OpError); ok || err == io.EOF {
				Out.Printf("Miner %v disconnected\n", conn.Conn.RemoteAddr())

				p.MinersMtx.Lock()
				delete(p.Miners, conn.Conn.RemoteAddr())
				p.MinersMtx.Unlock()

				conn.Conn.Close() // Close the connection

				return
			}
			Err.Printf(
				"Received error %v when processing message from miner %v\n",
				err,
				conn.Conn.RemoteAddr(),
			)
			continue
		}
		msgChan <- msg
	}
}

// GetMiners returns the addresses of any connected miners.
func (p *Pool) GetMiners() []net.Addr {
	p.MinersMtx.Lock()
	defer p.MinersMtx.Unlock()

	miners := []net.Addr{}
	for _, m := range p.Miners {
		miners = append(miners, m.Conn.RemoteAddr())
	}
	return miners
}

// GetClient returns the address of the current client or nil if there is no
// current client.
func (p *Pool) GetClient() net.Addr {
	p.ClientMtx.Lock()
	defer p.ClientMtx.Unlock()

	if p.Client.Conn == nil {
		return nil
	}
	return p.Client.Conn.RemoteAddr()
}
