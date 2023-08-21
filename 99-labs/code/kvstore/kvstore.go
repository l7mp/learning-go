package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	translog "transactionlog"
)

var logFile = "/tmp/translog.log"
var logger translog.TransactionLogger

type VersionedValue struct {
	Value   string `json:"value"`
	Version int    `json:"version"`
}

type VersionedKeyValue struct {
	Key string `json:"key"`
	VersionedValue
}

var store = make(map[string]VersionedValue)
var mu sync.RWMutex

func reset() {
	mu.RLock()
	defer mu.RUnlock()
	store = make(map[string]VersionedValue)
}

func get(key string) VersionedValue {
	mu.RLock()
	defer mu.RUnlock()
	return store[key]
}

func put(key string, v VersionedValue) error {
	mu.Lock()
	defer mu.Unlock()

	val, ok := store[key]
	if !ok {
		v.Version = 1
		store[key] = v
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
	store[key] = val

	return nil
}

func list() []VersionedKeyValue {
	mu.RLock()
	defer mu.RUnlock()

	ret := []VersionedKeyValue{}
	for k, v := range store {
		ret = append(ret, VersionedKeyValue{k, v})
	}

	return ret
}

func initTransactionLog() {
	fileLogger, err := translog.NewFileTransactionLogger(logFile)
	if err != nil {
		log.Fatalf("could not open transaction log: %s", err.Error())
	}
	logger = fileLogger

	events, err := logger.ReadEvents()
	if err != nil {
		log.Fatalf("could not read transaction log: %s", err.Error())
	}

	for _, e := range events {
		switch e.Type {
		case "put":
			var vkv VersionedKeyValue
			if err := json.Unmarshal([]byte(e.Value), &vkv); err != nil {
				log.Printf("error decoding transaction log value %q: %s",
					e.Value, err.Error())
				continue
			}

			if err := put(vkv.Key, vkv.VersionedValue); err != nil {
				log.Printf("error registering %#v from transaction log value: %s",
					vkv, err.Error())
				continue
			}
		case "reset":
			reset()
		default:
			log.Print("unknown key in transaction log")
		}
	}
}

func main() {
	// Include date, time and filename in the log messages
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Using transaction log in %q", logFile)
	initTransactionLog()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Run(ctx)

	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		log.Println("reset")

		reset()
		// Write record to transaction log
		logger.Write("reset", `""`)

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
		json.NewEncoder(w).Encode(get(key))

	})
	http.HandleFunc("/api/put", func(w http.ResponseWriter, r *http.Request) {
		vkv := VersionedKeyValue{}
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
		if err := put(vkv.Key, vkv.VersionedValue); err != nil {
			w.WriteHeader(http.StatusPreconditionRequired) // 418
			return
		}

		// Write record to transaction log
		json, _ := json.Marshal(vkv)
		logger.Write("put", string(json))

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		log.Println("list")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list())
	})

	log.Println("Starting HTTP server at :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
