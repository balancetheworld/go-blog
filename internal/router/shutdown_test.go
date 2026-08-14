package router

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zyj/my-blog/pkg/constant"
)

const gracefulShutdownHelperEnv = "GRACEFUL_SHUTDOWN_TEST_HELPER"

func TestRouterGracefulShutdownOnSIGTERM(t *testing.T) {
	if os.Getenv(gracefulShutdownHelperEnv) == "1" {
		runGracefulShutdownHelper(t)
		return
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(
		os.Args[0],
		"-test.run=^TestRouterGracefulShutdownOnSIGTERM$",
	)
	command.Env = append(
		os.Environ(),
		gracefulShutdownHelperEnv+"=1",
		constant.EnvKeyPort+"="+strconv.Itoa(port),
		constant.EnvKeyShutdownTimeout+"=3",
		constant.EnvKeyTrustedProxyCIDRs+"=",
	)
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	finished := false
	t.Cleanup(func() {
		if !finished {
			_ = command.Process.Kill()
			_ = command.Wait()
		}
	})

	client := &http.Client{Timeout: 5 * time.Second}
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	waitForHTTPStatus(t, client, baseURL+"/healthz", http.StatusOK)

	responseDone := make(chan error, 1)
	go func() {
		response, err := client.Get(baseURL + "/shutdown-test")
		if err != nil {
			responseDone <- err
			return
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			responseDone <- err
			return
		}
		if response.StatusCode != http.StatusOK || string(body) != "done" {
			responseDone <- fmt.Errorf(
				"unexpected response: status=%d body=%q",
				response.StatusCode,
				body,
			)
			return
		}
		responseDone <- nil
	}()

	waitForHTTPStatus(t, client, baseURL+"/shutdown-test/started", http.StatusOK)
	startedAt := time.Now()
	if err := command.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-responseDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for delayed response")
	}

	if err := command.Wait(); err != nil {
		finished = true
		t.Fatalf("server exited with error: %v\n%s", err, output.String())
	}
	finished = true
	if elapsed := time.Since(startedAt); elapsed < 300*time.Millisecond {
		t.Fatalf("server exited before delayed request completed: %s", elapsed)
	}
}

func runGracefulShutdownHelper(t *testing.T) {
	h, err := newServer()
	if err != nil {
		t.Fatal(err)
	}

	var requestStarted atomic.Bool
	h.GET("/shutdown-test", func(ctx context.Context, c *app.RequestContext) {
		requestStarted.Store(true)
		time.Sleep(500 * time.Millisecond)
		c.String(http.StatusOK, "done")
	})
	h.GET("/shutdown-test/started", func(ctx context.Context, c *app.RequestContext) {
		if requestStarted.Load() {
			c.String(http.StatusOK, "started")
			return
		}
		c.String(http.StatusServiceUnavailable, "waiting")
	})

	h.Spin()
}

func waitForHTTPStatus(
	t *testing.T,
	client *http.Client,
	url string,
	expectedStatus int,
) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		response, err := client.Get(url)
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode == expectedStatus {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for %s to return %d", url, expectedStatus)
}
