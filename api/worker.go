package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func pidFilePath() string {
	return filepath.Join(os.TempDir(), "db_migrate_tool_worker_pids.json")
}

// runningCmds holds started commands in memory so the server can Wait() them
var (
	runningMu   sync.Mutex
	runningCmds = make(map[string]*exec.Cmd)
)

func savePIDs(p map[string]int) error {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pidFilePath(), b, 0644)
}

func loadPIDs() (map[string]int, error) {
	out := make(map[string]int)
	b, err := os.ReadFile(pidFilePath())
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return out, err
	}
	return out, nil
}

func stopAllProcess() error {
	pids, err := loadPIDs()
	if err != nil {
		return fmt.Errorf("load pids: %w", err)
	}

	// First try to stop in-memory commands (clean Wait and reap)
	runningMu.Lock()
	for name, cmd := range runningCmds {
		if cmd == nil || cmd.Process == nil {
			continue
		}
		_ = cmd.Process.Signal(syscall.SIGTERM)
		// give graceful time
		done := make(chan struct{})
		go func(name string, cmd *exec.Cmd) {
			_ = cmd.Wait()
			close(done)
		}(name, cmd)

		select {
		case <-done:
			// reaped
		case <-time.After(1 * time.Second):
			_ = cmd.Process.Kill()
		}
	}
	runningMu.Unlock()

	// Also attempt to kill any remaining pids from the pid file
	for _, pid := range pids {
		// kill process group to ensure child processes are stopped
		pgid := -pid
		_ = syscall.Kill(pgid, syscall.SIGTERM)
		time.Sleep(300 * time.Millisecond)
		_ = syscall.Kill(pgid, syscall.SIGKILL)
		_ = syscall.Kill(pid, syscall.SIGKILL)
		_ = syscall.Kill(pid, syscall.SIGTERM)
		_ = syscall.Kill(pgid, syscall.SIGTERM)
		_ = syscall.Kill(pgid, syscall.SIGKILL)
	}

	return nil
}

func (server *Server) GetWorkerStatus(c *gin.Context) {

	pids, err := loadPIDs()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: fmt.Sprintf("failed to load pid file: %v", err),
		})
		return
	}

	for name, pid := range pids {
		// check if process is running: on unix, Kill(pid, 0) returns
		// ESRCH when process does not exist, nil when it exists (or EPERM if no permission)
		err := syscall.Kill(pid, 0)
		if err != nil {
			if err == syscall.ESRCH {
				c.JSON(http.StatusOK, DefaultResponse{
					Status:  "error",
					Message: fmt.Sprintf("process %s (%d) not running", name, pid),
				})
				return
			}
			// EPERM or other errors -> we assume process exists but permission issue
		}
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "All worker processes are running",
	})

}

func (server *Server) RunWorker(c *gin.Context) {
	cmds := map[string][]string{
		"parser":   {"go", "run", "./cmd/parser"},
		"joiner":   {"go", "run", "./cmd/joiner"},
		"checker":  {"go", "run", "./cmd/checker"},
		"executor": {"go", "run", "./cmd/executor"},
	}

	pids := make(map[string]int)

	isAllRunning := true
	workerStartErr := error(nil)
	for name, argv := range cmds {
		cmd := exec.Command(argv[0], argv[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		if err := cmd.Start(); err != nil {
			workerStartErr = fmt.Errorf("start %s: %w", name, err)
			isAllRunning = false
			break
		}

		pids[name] = cmd.Process.Pid

		// register running cmd so we can Wait() later and avoid zombies
		runningMu.Lock()
		runningCmds[name] = cmd
		runningMu.Unlock()

		// spawn goroutine to Wait() and cleanup when process exits
		go func(name string, cmd *exec.Cmd) {
			_ = cmd.Wait()
			runningMu.Lock()
			delete(runningCmds, name)
			runningMu.Unlock()
		}(name, cmd)

		// small delay so processes have time to initialize
		time.Sleep(100 * time.Millisecond)
	}

	if isAllRunning {
		if err := savePIDs(pids); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   err.Error(),
				Status:  "error",
				Message: "Failed to save pid file",
			})
			return
		}

		c.JSON(http.StatusOK, DefaultResponse{
			Status:  "success",
			Message: "Worker processes started successfully",
		})

	} else {
		// delete all started processes
		_ = stopAllProcess()

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   workerStartErr.Error(),
			Status:  "error",
			Message: "Failed to start all worker processes",
		})
	}
}

func (server *Server) StopWorker(c *gin.Context) {
	_, err := loadPIDs()
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: fmt.Sprintf("failed to load pid file: %v", err),
		})
		return
	}

	// stop all processes
	err = stopAllProcess()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   err.Error(),
			Status:  "error",
			Message: "Failed to stop all worker processes",
		})
		return
	}

	c.JSON(http.StatusOK, DefaultResponse{
		Status:  "success",
		Message: "Worker processes stopped successfully",
	})
}
