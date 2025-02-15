//go:build windows
// +build windows

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

	"golang.org/x/sys/windows"
)

type Echo struct{}

func NewEcho() *Echo {
	return &Echo{}
}

func (we *Echo) Disable() uint32 {
	var st uint32

	if err := windows.GetConsoleMode(windows.Handle(os.Stdout.Fd()), &st); err != nil {
		log.Fatalf("failed to get the console state: %v", err)
	}

	newSt := st
	newSt = newSt &^ windows.ENABLE_ECHO_INPUT
	newSt |= windows.ENABLE_PROCESSED_INPUT
	newSt |= windows.ENABLE_LINE_INPUT
	newSt |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

	if err := windows.SetConsoleMode(windows.Handle(os.Stdout.Fd()), newSt); err != nil {
		log.Fatalf("failed to set the console state: %v", err)
	}

	return st
}

func (we *Echo) Enable(st uint32) {
	if err := windows.SetConsoleMode(windows.Handle(os.Stdout.Fd()), st); err != nil {
		log.Fatalf("failed to set the console state: %v", err)
	}
}
