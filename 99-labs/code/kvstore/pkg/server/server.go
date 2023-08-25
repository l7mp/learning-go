package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	translog "transactionlog"

	"kvstore/pkg/api"
)

const timeout = 10 * time.Minute

type Server struct {
	store  map[string]api.VersionedValue
	mu     sync.RWMutex
	server *http.Server
	logger translog.TransactionLogger
}

func NewServer(logFile string) (*Server, error) {
	s := Server{store: make(map[string]api.VersionedValue)}

	log.Printf("Using transaction log in %q", logFile)
	fileLogger, err := translog.NewFileTransactionLogger(logFile)
	if err != nil {
		fmt.Errorf("could not open transaction log: %s", err.Error())
	}
	s.logger = fileLogger

	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		log.Println("reset")

		s.reset()
		// Write record to transaction log
		s.logger.Write("reset", `""`)

	})

	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {
		key := ""
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("get: key=%s\n", key)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.get(key))

	})

	http.HandleFunc("/api/put", func(w http.ResponseWriter, r *http.Request) {
		vkv := api.VersionedKeyValue{}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&vkv); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if vkv.Key == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("put: key=%s, value=%s, version=%d\n", vkv.Key, vkv.Value, vkv.Version)
		if err := s.put(vkv.Key, vkv.VersionedValue); err != nil {
			w.WriteHeader(http.StatusPreconditionRequired) // 418
			return
		}

		// Write record to transaction log
		json, _ := json.Marshal(vkv)
		s.logger.Write("put", string(json))

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		log.Println("list")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.list())
	})

	return &s, nil
}

func (s *Server) reset() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.store = make(map[string]api.VersionedValue)
}

func (s *Server) get(key string) api.VersionedValue {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store[key]
}

func (s *Server) put(key string, v api.VersionedValue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.store[key]
	if !ok {
		v.Version = 1
		s.store[key] = v
		return nil
	}

	if val.Value == v.Value {
		return nil
	}

	if val.Version != v.Version {
		return fmt.Errorf("put: version mismatch: %d != %d",
			val.Version, v.Version)
	}

	val.Value = v.Value
	val.Version += 1
	s.store[key] = val

	return nil
}

func (s *Server) list() []api.VersionedKeyValue {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ret := []api.VersionedKeyValue{}
	for k, v := range s.store {
		ret = append(ret, api.VersionedKeyValue{k, v})
	}

	return ret
}

func (s *Server) initTransactionLog() error {
	events, err := s.logger.ReadEvents()
	if err != nil {
		return fmt.Errorf("could not read transaction log: %w", err)
	}

	for _, e := range events {
		switch e.Type {
		case "put":
			var vkv api.VersionedKeyValue
			if err := json.Unmarshal([]byte(e.Value), &vkv); err != nil {
				log.Printf("error decoding transaction log value %q: %s",
					e.Value, err.Error())
				continue
			}

			if err := s.put(vkv.Key, vkv.VersionedValue); err != nil {
				log.Printf("error registering %#v from transaction log value: %s",
					vkv, err.Error())
				continue
			}
		case "reset":
			s.reset()
		default:
			log.Print("unknown key in transaction log")
		}
	}

	return nil
}

func (s *Server) Run(ctx context.Context, addr string) error {
	if err := s.initTransactionLog(); err != nil {
		return err
	}
	s.logger.Run(ctx)

	log.Printf("Starting HTTP server at %s", addr)
	s.server = &http.Server{Addr: addr}
	go func() {
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()

	return s.server.Close()
}
