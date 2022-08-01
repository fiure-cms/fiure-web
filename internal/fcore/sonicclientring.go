package fcore

import (
	"context"
	"errors"
	"github.com/uretgec/go-sonic/sonic"
	"sync/atomic"
)

type SonicClientRing struct {
	List    []string
	current uint64
	Clients map[string]*sonic.Client
	Ctx     context.Context
}

func NewSonicClientRing() *SonicClientRing {
	return &SonicClientRing{
		List:    []string{},
		Clients: make(map[string]*sonic.Client),
		Ctx:     context.Background(),
	}
}

func (sr *SonicClientRing) GenerateOptions(channelMode, authPassword, sonicAddr string, sonicPoolSize, sonicMinIdleConns, sonicMaxRetries int) *sonic.Options {
	return &sonic.Options{
		Addr:         sonicAddr,
		AuthPassword: authPassword,
		PoolSize:     sonicPoolSize, // Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
		MinIdleConns: sonicMinIdleConns,
		MaxRetries:   sonicMaxRetries,
		ChannelMode:  channelMode,
	}
}

func (sr *SonicClientRing) GetRing(clientName string) *sonic.Client {
	return sr.Clients[clientName]
}

func (sr *SonicClientRing) AddRing(clientName string, option *sonic.Options) {
	sr.List = append(sr.List, clientName)

	newClient := sonic.NewClient(option)
	sr.Clients[clientName] = newClient
}

func (sr *SonicClientRing) CloseRing(clientName string) {
	sr.Clients[clientName].Close()
}

func (sr *SonicClientRing) HasRing(clientName string) bool {
	if _, ok := sr.Clients[clientName]; ok {
		return true
	}

	return false
}

func (sr *SonicClientRing) HealthCheckedRing(clientName string) error {
	if err := sr.Clients[clientName].Ping(sr.Ctx).Err(); err != nil {
		return errors.New("sonic server not response: " + clientName)
	}

	return nil
}

// Round Robin Implementation
func (sr *SonicClientRing) NextRing() string {
	next := int(atomic.AddUint64(&sr.current, uint64(1)) % uint64(len(sr.List)))
	return sr.List[next]
}

func (sr *SonicClientRing) GetCurrentClient() *sonic.Client {
	return sr.Clients[sr.NextRing()]
}
