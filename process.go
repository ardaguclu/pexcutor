package pexcutor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// Process stores pexcutor related internal fields like command, retry-recovery procedure, etc.
type Process struct {
	path   string          // execution path which can be defined $PATH, or relative path being searched with lookpath
	args   []string        // args passed to path
	envs   []string        // optional environment variables
	rc     int             // retry count when the process crashes
	crc    int             // current retry count which is incremented for each retry
	cmd    *exec.Cmd       // external process command
	ctx    context.Context // storing context in struct field is not the best practice, but in retry cases, in order to guarantee that same context is used, best option is storing context in here.
	stdOut io.Reader
	stdErr io.Reader
}

// New creates command and returns Process with given initialized values.
func New(ctx context.Context, retryCount int, path string, args ...string) *Process {
	return &Process{
		ctx:  ctx,
		path: path,
		args: args,
		rc:   retryCount,
		crc:  0,
	}
}

// SetEnv sets environment variables for command
func (p *Process) SetEnv(envs ...string) {
	p.envs = append(os.Environ(), envs...)
}

// Start starts command in process
func (p *Process) Start() error {
	p.cmd = exec.CommandContext(p.ctx, p.path, p.args...)
	if p.envs != nil {
		p.cmd.Env = p.envs
	}

	so, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	se, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}

	p.stdOut = so
	p.stdErr = se

	err = p.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

// GetResult returns the result of process has been started
func (p *Process) GetResult() (string, string, error) {
	if p.cmd == nil {
		return "", "", nil
	}

	var stdOut, stdErr string

	var wg sync.WaitGroup
	sOut := bufio.NewScanner(p.stdOut)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for sOut.Scan() {
			if sOut.Text() != "" {
				stdOut += fmt.Sprintf("%s\n", sOut.Text())
			}
		}
	}()

	sErr := bufio.NewScanner(p.stdErr)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for sErr.Scan() {
			if sErr.Text() != "" {
				stdErr += fmt.Sprintf("%s\n", sErr.Text())
			}
		}
	}()

	wg.Wait()

	if err := p.cmd.Wait(); err != nil {
		if eErr, ok := err.(*exec.ExitError); ok {
			st := eErr.ProcessState.Sys().(syscall.WaitStatus)
			if st.CoreDump() && p.crc < p.rc {
				log.Println("crashed ", eErr, " process will be restarted")
				p.crc++
				err = p.Start()
				if err != nil {
					return stdOut, stdErr, err
				}

				return p.GetResult()
			}

			return stdOut, stdErr, eErr
		}
	}

	return stdOut, stdErr, nil
}

// Stop sends sigterm signal the already running external process
func (p *Process) Stop() error {
	if p.cmd == nil {
		return nil
	}

	// Instead terminating process with SIGKILL, shutting down with sigterm gives the process cleaning up.
	if err := p.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	log.Println("processed stopped by caller.")
	return nil
}

// Signal sends signal to the external process.
// This functions can be used as a pipe for relaying signals captured by caller and send to the callee.
func (p *Process) Signal(sig os.Signal) error {
	if p.cmd == nil {
		return nil
	}

	return p.cmd.Process.Signal(sig)
}
