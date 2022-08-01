package services

import (
	"github.com/fiure-cms/fiure-web/internal/fcore"
	"github.com/fiure-cms/fiure-web/loggers"
)

var Ss *fcore.SonicClientRing

func SetupSonicSearchClient(serviceName string, sonicClients, sonicAddrs, sonicChannelMode, sonicAuthPassword fcore.ArrayFlagString, sonicPoolSizes, sonicMinIdleConns, sonicMaxRetries fcore.ArrayFlagInt) {
	if len(sonicClients) == 0 {
		loggers.Sugar.Fatal("sonic client names not found")
	}

	Ss = fcore.NewSonicClientRing()
	for index, clientName := range sonicClients {

		ringOption := Ss.GenerateOptions(sonicChannelMode[index], sonicAuthPassword[index], sonicAddrs[index], sonicPoolSizes[index], sonicMinIdleConns[index], sonicMaxRetries[index])
		Ss.AddRing(sonicChannelMode[index], ringOption)
		//defer ss.CloseRing(clientName)

		err := Ss.HealthCheckedRing(sonicChannelMode[index])
		if err != nil {
			loggers.Sugar.With("error", err, "clientName", clientName).Fatal("ring health check error")
		}
	}
}
