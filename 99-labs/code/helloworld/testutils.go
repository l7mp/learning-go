package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func testServerAvailable(t *testing.T, ctx context.Context, port int) {
// 	if _, err := http.Get(fmt.Sprintf("http://:%d", port)); err != nil {
// 		// Server is not running

// 		cmd := exec.CommandContext(ctx, "go", "run", "main.go")
// 		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
// 		go cmd.Run()

// 		// Wait for the server to start
// 		time.Sleep(500 * time.Millisecond)
// 	}
// }

func execCmd(t *testing.T, cmd string, args ...string) (string, string) {
	return execCmdContext(context.Background(), t, cmd, args...)
}

func execCmdContext(ctx context.Context, t *testing.T, cmd string, args ...string) (string, string) {
	p, err := exec.LookPath(cmd)
	if errors.Is(err, exec.ErrDot) {
		err = nil
	}
	assert.NoError(t, err, fmt.Sprintf("find command %q in PATH", cmd))

	log.Print("Executing:\t", cmd, " ", strings.Join(args, " "))

	e := exec.CommandContext(ctx, p, args...)
	var outb, errb bytes.Buffer
	e.Stdout = &outb
	e.Stderr = &errb
	assert.NoError(t, e.Run(), fmt.Sprintf("run command %q", cmd))

	log.Print("StdOut:\t", outb.String())
	log.Print("StdErr:\t ", errb.String())

	return outb.String(), errb.String()
}
