/*
 *  Brown University, CS138, Spring 2022
 *
 *  Purpose: a CLI for LiteMiner clients.
 */

package main

import (
	"flag"
	liteminer "liteminer/pkg"

	"strconv"
	"strings"

	"github.com/abiosoft/ishell"
)

func main() {
	var addr string
	var debug bool

	flag.StringVar(&addr, "connect", "", "Address of the mining pool(s) to connect to.")
	flag.StringVar(&addr, "c", "", "Address of the mining pool(s) to connect to.")

	flag.BoolVar(&debug, "debug", false, "Turn debug message printing on or off – defaults to off.")
	flag.BoolVar(&debug, "d", false, "Turn debug message printing on or off – defaults to off.")

	flag.Parse()

	liteminer.SetDebug(debug)

	// Kick off shell
	shell := ishell.New()

	// Connect to the mining pools
	client := liteminer.CreateClient(addr)

	shell.AddCmd(&ishell.Cmd{
		Name: "connect",
		Help: "Connect to the specified pool(s)",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 1 {
				shell.Println("Usage: connect <pool addresses>")
				return
			}

			if client.Pool.Conn != nil {
				shell.Printf("Client already connected to pool at address %v\n", client.Pool.Conn.RemoteAddr())
				return
			}

			client.Connect(c.Args[0])
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "mine",
		Help: "Send a mine request to the connected pool(s)",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				shell.Println("Usage: mine <data> <upper bound on nonce>")
				return
			}

			upperBound, err := strconv.ParseUint(c.Args[1], 10, 64)
			if err != nil {
				shell.Println("Usage: mine <data> <upper bound on nonce>")
				return
			}

			nonce, err := client.Mine(c.Args[0], upperBound)
			if err != nil {
				shell.Println(err.Error())
			} else {
				shell.Printf("Result from pool: %v\n", nonce)
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "debug",
		Help: "Turn debug statements on or off",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				shell.Println("Usage: debug <on|off>")
				return
			}

			debugState := strings.ToLower(c.Args[0])

			switch debugState {
			case "on":
				liteminer.SetDebug(true)
				shell.Println("Debug turned on")
			case "off":
				liteminer.SetDebug(false)
				shell.Println("Debug turned off")
			default:
				shell.Println("Usage: debug <on|off>")
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "pool",
		Help: "Print the pool that the client is currently connected to",
		Func: func(c *ishell.Context) {
			shell.Println(client.GetPool())
		},
	})

	shell.Println(shell.HelpText())
	shell.Run()
}
