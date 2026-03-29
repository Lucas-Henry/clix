package ipc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestServerGetState(t *testing.T) {
	srv := newTestServer(t, "/tmp/test-ipc", "selected.go")
	defer srv.Stop()

	resp := roundtrip(t, srv.socketPath, Request{Action: ActionGetState})
	if !resp.OK {
		t.Fatalf("get_state failed: %s", resp.Error)
	}
	if resp.CWD != "/tmp/test-ipc" {
		t.Errorf("CWD = %q, want %q", resp.CWD, "/tmp/test-ipc")
	}
	if resp.Selected != "selected.go" {
		t.Errorf("Selected = %q, want %q", resp.Selected, "selected.go")
	}
}

func TestServerNavigate(t *testing.T) {
	srv := newTestServer(t, "/tmp/start", "")
	defer srv.Stop()

	resp := roundtrip(t, srv.socketPath, Request{Action: ActionNavigate, Path: "/tmp/dest"})
	if !resp.OK {
		t.Fatalf("navigate failed: %s", resp.Error)
	}
	if resp.CWD != "/tmp/dest" {
		t.Errorf("CWD after navigate = %q, want /tmp/dest", resp.CWD)
	}
}

func TestServerSelect(t *testing.T) {
	selected := ""
	srv := newTestServerFull(t, "/tmp/cwd", "init", func(name string) {
		selected = name
	})
	defer srv.Stop()

	resp := roundtrip(t, srv.socketPath, Request{Action: ActionSelect, Name: "main.go"})
	if !resp.OK {
		t.Fatalf("select failed: %s", resp.Error)
	}
	if selected != "main.go" {
		t.Errorf("selectFn got %q, want main.go", selected)
	}
}

func TestServerNavigateMissingPath(t *testing.T) {
	srv := newTestServer(t, "/tmp/start", "")
	defer srv.Stop()

	resp := roundtrip(t, srv.socketPath, Request{Action: ActionNavigate})
	if resp.OK {
		t.Error("expected error for navigate without path")
	}
}

func TestServerSelectMissingName(t *testing.T) {
	srv := newTestServer(t, "/tmp/start", "")
	defer srv.Stop()

	resp := roundtrip(t, srv.socketPath, Request{Action: ActionSelect})
	if resp.OK {
		t.Error("expected error for select without name")
	}
}

func TestServerUnknownAction(t *testing.T) {
	srv := newTestServer(t, "/tmp/start", "")
	defer srv.Stop()

	resp := roundtrip(t, srv.socketPath, Request{Action: "bogus"})
	if resp.OK {
		t.Error("expected error for unknown action")
	}
}

func newTestServer(t *testing.T, cwd, selected string) *Server {
	t.Helper()
	return newTestServerFull(t, cwd, selected, func(string) {})
}

func newTestServerFull(t *testing.T, cwd, selected string, selectFn SelectFunc) *Server {
	t.Helper()
	sock := fmt.Sprintf("/tmp/clix-test-%d-%d.sock", os.Getpid(), time.Now().UnixNano())

	currentCWD := cwd
	srv := &Server{
		socketPath: sock,
		stateFn: func() (string, string) {
			return currentCWD, selected
		},
		navigateFn: func(path string) error {
			currentCWD = path
			return nil
		},
		selectFn: selectFn,
	}

	if err := srv.Start(); err != nil {
		t.Fatalf("server start: %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	return srv
}

func roundtrip(t *testing.T, sock string, req Request) Response {
	t.Helper()
	conn, err := net.Dial("unix", sock)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	if err := enc.Encode(req); err != nil {
		t.Fatalf("encode: %v", err)
	}

	var resp Response
	if err := dec.Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return resp
}
