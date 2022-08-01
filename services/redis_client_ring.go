package services

import (
	"github.com/fiure-cms/fiure-web/internal/fcore"
	"github.com/fiure-cms/fiure-web/loggers"
)

var Sr *fcore.RedisClientRing

func SetupRedisClientRing(redisClients, redisAddrs fcore.ArrayFlagString, redisPoolSizes, redisMinIdleConns, redisMaxRetries fcore.ArrayFlagInt) {
	if len(redisClients) == 0 {
		loggers.Sugar.Fatal("redis client names not found")
	}

	Sr = fcore.NewRedisClientRing()
	for index, clientName := range redisClients {

		ringOption := Sr.GenerateOptions(redisAddrs[index], redisPoolSizes[index], redisMinIdleConns[index], redisMaxRetries[index])
		Sr.AddRing(clientName, ringOption)
		//defer Sr.CloseRing(clientName)

		err := Sr.HealthCheckedRing(clientName)
		if err != nil {
			loggers.Sugar.With("error", err).Fatal("ring health check error")
		}
	}
}
