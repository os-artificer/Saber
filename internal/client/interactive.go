/**
 * Copyright 2025 Saber authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package client

import (
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// stdio ties stdin and stdout into a single io.ReadWriter for use with term.NewTerminal.
type stdio struct{}

func (*stdio) Read(p []byte) (n int, err error)  { return os.Stdin.Read(p) }
func (*stdio) Write(p []byte) (n int, err error) { return os.Stdout.Write(p) }

func setupGracefulShutdown(svr *Service) {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigC
		svr.Close()
		os.Exit(0)
	}()
}

// runInteractive uses golang.org/x/term.Terminal for line editing, echo, and Ctrl+L clear.
// See: https://pkg.go.dev/golang.org/x/term#Terminal
func runInteractive() error {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		// Not a TTY: fall back to simple line reading without raw mode.
		return runInteractiveStdin()
	}
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState)

	t := term.NewTerminal(&stdio{}, "SABER> ")
	for {
		line, err := t.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			if err == term.ErrPasteIndicator {
				continue
			}
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch strings.ToLower(line) {
		case "exit", "quit":
			return nil
		case "clear":
			clearScreenTerminal(t)
		default:
			// placeholder: unknown command, can be extended later
		}
	}
	return nil
}

// clearScreenTerminal sends VT100 clear-screen and home cursor to the terminal.
func clearScreenTerminal(t *term.Terminal) {
	t.Write([]byte("\x1b[2J\x1b[H"))
}

// runInteractiveStdin is used when stdin is not a TTY (e.g. pipe); no raw mode, no Ctrl+L.
func runInteractiveStdin() error {
	buf := make([]byte, 0, 256)
	for {
		os.Stdout.WriteString("SABER> ")
		buf = buf[:0]
		for {
			var b [1]byte
			n, err := os.Stdin.Read(b[:])
			if err != nil {
				if err == io.EOF && len(buf) > 0 {
					break
				}
				return err
			}
			if n == 0 {
				continue
			}
			c := b[0]
			if c == '\n' || c == '\r' {
				break
			}
			buf = append(buf, c)
			os.Stdout.Write(b[:])
		}
		line := strings.TrimSpace(string(buf))
		if line == "" {
			continue
		}
		switch strings.ToLower(line) {
		case "exit", "quit":
			return nil
		case "clear":
			os.Stdout.WriteString("\x1b[2J\x1b[H")
		default:
		}
	}
}
