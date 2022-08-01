package fcore

import (
	"log"
	"testing"
)

func TestRingMehod(t *testing.T) {
	redisClients := []string{"moderatestore", "livestore"}
	redisAddrs := []string{"localhost:6479", "localhost:6480"}

	// RedisClientRing Conn
	sr := NewRedisClientRing()

	for index, clientName := range redisClients {

		ringOption := sr.GenerateOptions(redisAddrs[index], 1000, 2, 2)
		sr.AddRing(clientName, ringOption)
		defer sr.CloseRing(clientName)

		err := sr.HealthCheckedRing(clientName)
		if err != nil {
			log.Fatalf("storage ring healthcheck error: %v\n", err)
		}
	}

	for i := 0; i < 10; i++ {
		nextClient := sr.NextRing()
		log.Println("next client", nextClient)

		resp, err := sr.Clients[nextClient].Ping(sr.Ctx).Result()
		if err != nil {
			log.Fatalf("redis server not response: %s\n", nextClient)
		}

		log.Printf("%s : %s\n", nextClient, resp)
	}
}
