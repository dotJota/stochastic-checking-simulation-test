package main

import (
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"log"
	"stochastic-checking-simulation/impl/parameters"
	"stochastic-checking-simulation/impl/protocols/bracha"
)

func main() {
	processCount := 256
	parameters := &parameters.Parameters{
		FaultyProcesses: 20,
	}

	system := actor.NewActorSystem()
	remoteConfig := remote.Configure("127.0.0.1", 8080)
	remoter := remote.NewRemote(system, remoteConfig)
	remoter.Start()

	pids := make([]*actor.PID, processCount)
	processes := make([]*bracha.Process, processCount)

	for i := 0; i < processCount; i++ {
		processes[i] = &bracha.Process{}
		pids[i] =
			system.Root.Spawn(
				actor.PropsFromProducer(
					func() actor.Actor {
						return processes[i]
					}),
			)
	}

	logger := log.Default()
	for i := 0; i < processCount; i++ {
		processes[i].InitProcess(pids[i], pids, parameters, logger)
	}

	processes[0].Broadcast(system.Root, 5)

	_, _ = console.ReadLine()
}
