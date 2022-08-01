package fcore

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/go-redis/redis/v8"
)

type RedisClientRing struct {
	List    []string
	current uint64
	Clients map[string]*redis.Client
	Ctx     context.Context
}

func NewRedisClientRing() *RedisClientRing {
	return &RedisClientRing{
		List:    []string{},
		Clients: make(map[string]*redis.Client),
		Ctx:     context.Background(),
	}
}

func (sr *RedisClientRing) GenerateOptions(redisAddr string, redisPoolSize, redisMinIdleConns, redisMaxRetries int) *redis.Options {
	return &redis.Options{
		Addr:         redisAddr,
		Password:     "",            // no password set
		DB:           0,             // use default DB
		PoolSize:     redisPoolSize, // Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
		MinIdleConns: redisMinIdleConns,
		MaxRetries:   redisMaxRetries,
	}
}

func (sr *RedisClientRing) GetRing(clientName string) *redis.Client {
	return sr.Clients[clientName]
}

func (sr *RedisClientRing) AddRing(clientName string, option *redis.Options) {
	sr.List = append(sr.List, clientName)

	newClient := redis.NewClient(option)
	sr.Clients[clientName] = newClient
}

func (sr *RedisClientRing) CloseRing(clientName string) {
	sr.Clients[clientName].Close()
}

func (sr *RedisClientRing) HasRing(clientName string) bool {
	if _, ok := sr.Clients[clientName]; ok {
		return true
	}

	return false
}

func (sr *RedisClientRing) HealthCheckedRing(clientName string) error {
	if _, err := sr.Clients[clientName].Ping(sr.Ctx).Result(); err != nil {
		return errors.New("redis server not response: " + clientName)
	}

	return nil
}

// Round Robin Implementation
func (sr *RedisClientRing) NextRing() string {
	next := int(atomic.AddUint64(&sr.current, uint64(1)) % uint64(len(sr.List)))
	return sr.List[next]
}

func (sr *RedisClientRing) GetCurrentClient() *redis.Client {
	return sr.Clients[sr.NextRing()]
}
