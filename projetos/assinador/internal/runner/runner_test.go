package runner

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestParseResponseSuccess(t *testing.T) {
	response, err := parseResponse([]byte(`{"success":true,"data":"SIGNATURE"}`))
	if err != nil {
		t.Fatalf("parseResponse returned error: %v", err)
	}
	if !response.Success || response.Data != "SIGNATURE" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestParseResponseStructuredError(t *testing.T) {
	response, err := parseResponse([]byte(`{"success":false,"errorCode":"INVALID","errorMessage":"bad input"}`))
	if err != nil {
		t.Fatalf("parseResponse returned error: %v", err)
	}
	if response.Success || response.ErrorCode != "INVALID" || response.ErrorMessage != "bad input" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestParseResponseRejectsInvalidJSON(t *testing.T) {
	if _, err := parseResponse([]byte(`not-json`)); err == nil {
		t.Fatal("expected invalid JSON to be rejected")
	}
}

func TestNormalizePortUsesDefaultWhenUnset(t *testing.T) {
	if got := normalizePort(0); got != defaultServerPort {
		t.Fatalf("normalizePort(0) = %d, want %d", got, defaultServerPort)
	}
	if got := normalizePort(9090); got != 9090 {
		t.Fatalf("normalizePort(9090) = %d, want 9090", got)
	}
}

func TestStartServerRejectsNegativeTimeout(t *testing.T) {
	if _, err := StartServer(8080, -1); err == nil {
		t.Fatal("expected StartServer to reject negative timeout")
	}
}

func TestStartServerFailsClearlyWhenPortIsOccupiedByAnotherProcess(t *testing.T) {
	useTempHome(t)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen returned error: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	_, err = StartServer(port, 0)
	if err == nil {
		t.Fatal("expected StartServer to reject occupied port")
	}
	if !strings.Contains(err.Error(), "porta") || !strings.Contains(err.Error(), "indisponivel") {
		t.Fatalf("error = %q, want clear occupied-port message", err)
	}
}

func TestInvokeSignUsesHTTPWhenServerIsActive(t *testing.T) {
	port, closeServer, seen := startFakeAssinadorHTTP(t)
	defer closeServer()
	useTempHome(t)

	if err := writeServerState(&ServerState{
		PID:       1234,
		Port:      port,
		JavaPath:  "java",
		JarPath:   "assinador.jar",
		StartedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("writeServerState returned error: %v", err)
	}

	restore := replaceInvokeLocal(t, func(string, map[string]interface{}) (*Response, error) {
		t.Fatal("invokeLocal should not be called when HTTP server is active")
		return nil, nil
	})
	defer restore()

	response, err := InvokeSignWithOptions(map[string]interface{}{"bundle": "FHIR"}, InvokeOptions{Port: port})
	if err != nil {
		t.Fatalf("InvokeSignWithOptions returned error: %v", err)
	}
	if !response.Success || response.Data != "HTTP-SIGNATURE" {
		t.Fatalf("unexpected response: %+v", response)
	}
	if *seen != "/sign" {
		t.Fatalf("HTTP endpoint = %q, want /sign", *seen)
	}
}

func TestInvokeValidateFallsBackToLocalWhenServerIsUnavailable(t *testing.T) {
	useTempHome(t)

	restore := replaceInvokeLocal(t, func(command string, payload map[string]interface{}) (*Response, error) {
		if command != "validate" {
			t.Fatalf("command = %q, want validate", command)
		}
		return &Response{Success: true, Data: "LOCAL-VALID"}, nil
	})
	defer restore()

	response, err := InvokeValidateWithOptions(map[string]interface{}{"signatureData": "abc"}, InvokeOptions{Port: 1})
	if err != nil {
		t.Fatalf("InvokeValidateWithOptions returned error: %v", err)
	}
	if !response.Success || response.Data != "LOCAL-VALID" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestInvokeLocalOptionBypassesActiveHTTPServer(t *testing.T) {
	port, closeServer, seen := startFakeAssinadorHTTP(t)
	defer closeServer()
	useTempHome(t)

	if err := writeServerState(&ServerState{Port: port, StartedAt: time.Now().UTC()}); err != nil {
		t.Fatalf("writeServerState returned error: %v", err)
	}

	restore := replaceInvokeLocal(t, func(command string, payload map[string]interface{}) (*Response, error) {
		if command != "sign" {
			t.Fatalf("command = %q, want sign", command)
		}
		return &Response{Success: true, Data: "LOCAL-SIGNATURE"}, nil
	})
	defer restore()

	response, err := InvokeSignWithOptions(map[string]interface{}{"bundle": "FHIR"}, InvokeOptions{Local: true, Port: port})
	if err != nil {
		t.Fatalf("InvokeSignWithOptions returned error: %v", err)
	}
	if !response.Success || response.Data != "LOCAL-SIGNATURE" {
		t.Fatalf("unexpected response: %+v", response)
	}
	if *seen != "" {
		t.Fatalf("HTTP endpoint = %q, want no HTTP call", *seen)
	}
}

func TestStopServerKillsRegisteredProcessAndRemovesState(t *testing.T) {
	useTempHome(t)

	cmd := exec.Command(os.Args[0], "-test.run=TestHelperLongRunningProcess")
	cmd.Env = append(os.Environ(), "ASSINADOR_TEST_HELPER_PROCESS=1")
	if err := cmd.Start(); err != nil {
		t.Fatalf("cmd.Start returned error: %v", err)
	}
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_, _ = cmd.Process.Wait()
		}
	})

	state := &ServerState{
		PID:       cmd.Process.Pid,
		Port:      19080,
		JavaPath:  "java",
		JarPath:   "assinador.jar",
		StartedAt: time.Now().UTC(),
	}
	if err := writeServerState(state); err != nil {
		t.Fatalf("writeServerState returned error: %v", err)
	}

	stopped, err := StopServer(state.Port)
	if err != nil {
		t.Fatalf("StopServer returned error: %v", err)
	}
	if stopped.PID != state.PID {
		t.Fatalf("stopped PID = %d, want %d", stopped.PID, state.PID)
	}
	if _, err := os.Stat(serverStatePath(state.Port)); !os.IsNotExist(err) {
		t.Fatalf("state file should be removed, stat error = %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("process was not stopped")
	}
}

func TestStopServerReturnsErrorWhenNoStateExists(t *testing.T) {
	useTempHome(t)

	if _, err := StopServer(19180); err == nil {
		t.Fatal("expected StopServer to fail without registered state")
	}
}

func TestConcurrentStartServerOnSamePortHandlesRaceCondition(t *testing.T) {
	useTempHome(t)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen returned error: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	const goroutines = 3
	var mu sync.Mutex
	reuseCount := 0
	errCount := 0

	for i := 0; i < goroutines; i++ {
		go func() {
			state, err := StartServer(port, 0)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errCount++
				return
			}
			if state.Reused {
				reuseCount++
			}
		}()
	}

	time.Sleep(3 * time.Second)

	mu.Lock()
	if reuseCount+errCount < 2 {
		t.Fatalf("expected at least 2 goroutines to reuse or error, got reuseCount=%d errCount=%d", reuseCount, errCount)
	}
	mu.Unlock()

	if state, err := readServerState(port); err == nil {
		if state.PID > 0 {
			StopServer(port)
		}
	}
}

func TestHelperLongRunningProcess(t *testing.T) {
	if os.Getenv("ASSINADOR_TEST_HELPER_PROCESS") != "1" {
		return
	}

	for {
		time.Sleep(time.Second)
	}
}

func replaceInvokeLocal(t *testing.T, fn func(string, map[string]interface{}) (*Response, error)) func() {
	t.Helper()
	previous := invokeLocalFunc
	invokeLocalFunc = fn
	return func() {
		invokeLocalFunc = previous
	}
}

func useTempHome(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")
	if err := os.MkdirAll(filepath.Join(home, ".hubsaude"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
}

func startFakeAssinadorHTTP(t *testing.T) (int, func(), *string) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen returned error: %v", err)
	}

	seen := ""
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	})
	mux.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Path
		assertJSONRequest(t, r)
		writeJSON(t, w, Response{Success: true, Data: "HTTP-SIGNATURE"})
	})
	mux.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Path
		assertJSONRequest(t, r)
		writeJSON(t, w, Response{Success: true, Data: "HTTP-VALID"})
	})

	server := &http.Server{Handler: mux}
	go func() {
		_ = server.Serve(listener)
	}()

	port := listener.Addr().(*net.TCPAddr).Port
	return port, func() {
		_ = server.Close()
	}, &seen
}

func assertJSONRequest(t *testing.T, r *http.Request) {
	t.Helper()
	if r.Method != http.MethodPost {
		t.Fatalf("method = %q, want POST", r.Method)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatalf("request body is not valid JSON: %v", err)
	}
}

func writeJSON(t *testing.T, w http.ResponseWriter, response Response) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.Fatalf("json.Encode returned error: %v", err)
	}
}
