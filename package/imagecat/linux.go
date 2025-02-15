//go:build linux
// +build linux

// Derivative work from https://github.com/danielgatis/imgcat

// MIT License

// Copyright (c) 2020 Daniel Gatis

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package imagecat

import (
	"log"
	"os"

	"golang.org/x/sys/unix"
)

type Echo struct{}

func NewEcho() *Echo {
	return &Echo{}
}

func (e *Echo) Disable() *unix.Termios {
	termios, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)
	if err != nil {
		log.Fatalf("failed to get the termios: %v", err)
	}

	newState := *termios
	newState.Lflag &^= unix.ECHO
	newState.Lflag |= unix.ICANON | unix.ISIG
	newState.Iflag |= unix.ICRNL
	if err := unix.IoctlSetTermios(int(os.Stdout.Fd()), unix.TCSETS, &newState); err != nil {
		log.Fatalf("failed to set the termios: %v", err)
	}

	return termios
}

func (e *Echo) Enable(termios *unix.Termios) {
	if err := unix.IoctlSetTermios(int(os.Stdout.Fd()), unix.TCSETS, termios); err != nil {
		log.Fatalf("failed to set the termios: %v", err)
	}
}
