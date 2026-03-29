package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type Action string

const (
	ActionNavigate Action = "navigate"
	ActionGetState Action = "get_state"
	ActionSelect   Action = "select"
)

type Request struct {
	Action Action `json:"action"`
	Path   string `json:"path,omitempty"`
	Name   string `json:"name,omitempty"`
}

type Response struct {
	OK       bool   `json:"ok"`
	CWD      string `json:"cwd,omitempty"`
	Selected string `json:"selected,omitempty"`
	Error    string `json:"error,omitempty"`
}

type StateFunc    func() (cwd string, selected string)
type NavigateFunc func(path string) error
type SelectFunc   func(name string)

type Server struct {
	socketPath string
	listener   net.Listener
	stateFn    StateFunc
	navigateFn NavigateFunc
	selectFn   SelectFunc
	mu         sync.Mutex
}

func SocketPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("clix-%s.sock", strconv.Itoa(os.Getpid())))
}

func NewServer(stateFn StateFunc, navigateFn NavigateFunc, selectFn SelectFunc) *Server {
	return &Server{
		socketPath: SocketPath(),
		stateFn:    stateFn,
		navigateFn: navigateFn,
		selectFn:   selectFn,
	}
}

func (s *Server) Start() error {
	os.Remove(s.socketPath)
	ln, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("ipc listen: %w", err)
	}
	s.listener = ln
	go s.accept()
	return nil
}

func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
	os.Remove(s.socketPath)
}

func (s *Server) accept() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	enc := json.NewEncoder(conn)

	for scanner.Scan() {
		var req Request
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			enc.Encode(Response{Error: "invalid json"})
			continue
		}
		resp := s.dispatch(req)
		enc.Encode(resp)
	}
}

func (s *Server) dispatch(req Request) Response {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch req.Action {
	case ActionGetState:
		cwd, sel := s.stateFn()
		return Response{OK: true, CWD: cwd, Selected: sel}

	case ActionNavigate:
		if req.Path == "" {
			return Response{Error: "path required"}
		}
		if err := s.navigateFn(req.Path); err != nil {
			return Response{Error: err.Error()}
		}
		cwd, sel := s.stateFn()
		return Response{OK: true, CWD: cwd, Selected: sel}

	case ActionSelect:
		if req.Name == "" {
			return Response{Error: "name required"}
		}
		s.selectFn(req.Name)
		cwd, sel := s.stateFn()
		return Response{OK: true, CWD: cwd, Selected: sel}

	default:
		return Response{Error: fmt.Sprintf("unknown action: %s", req.Action)}
	}
}
