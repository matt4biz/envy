package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

type ExecCommand struct {
	*App
}

func (cmd *ExecCommand) Run() int {
	if len(cmd.args) < 2 {
		cmd.usage()
		return 1
	}

	var (
		m    []string
		k, v string
		err  error
	)

	parts := strings.Split(cmd.args[0], "/")

	if len(parts) == 1 {
		m, err = cmd.FetchAsVarList(cmd.args[0])
	} else {
		k = parts[1]
		v, err = cmd.Get(parts[0], parts[1])
	}

	if err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	if v != "" {
		m = append(m, k+"="+v)
	}

	done := make(chan os.Signal, 1)
	sub := exec.Command(cmd.args[1], cmd.args[2:]...)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sub.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	sub.Stdout = cmd.stdout
	sub.Stderr = cmd.stderr
	sub.Env = append(os.Environ(), m...)

	go func() {
		s := <-done

		if err := sub.Process.Signal(s); err != nil {
			fmt.Fprintln(cmd.stderr, "can't send signal", s)

			if err = syscall.Kill(-sub.Process.Pid, syscall.SIGKILL); err != nil {
				fmt.Fprintf(cmd.stderr, "failed to stop emulator pid=%d\n", sub.Process.Pid)
			}
		}
	}()

	if err := sub.Start(); err != nil {
		fmt.Fprintln(cmd.stderr, err)
		return -1
	}

	if err := sub.Wait(); err != nil {
		fmt.Fprintln(cmd.stderr, "can't wait", err)
	}

	return sub.ProcessState.ExitCode()
}
