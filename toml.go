package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type tomlTrigger struct {
	Id           string
	Message      string
	Queue        string
	GitHubSecret string
}

type tomlState struct {
	Version string
	Trigger []tomlTrigger
}

func SaveState(state State) {
	var ttomlState tomlState
	ttomlState.Version = "1.0"
	for _, trigger := range state.Triggers {
		tTrigger := tomlTrigger{trigger.Id, trigger.Message, trigger.Queue, trigger.GitHubSecret}

		ttomlState.Trigger = append(ttomlState.Trigger, tTrigger)
	}

	file, err := os.Create("/opt/trigger/state/state.toml")
	failOnError(err, "Unable to save state file")
	defer file.Close()

	enc := toml.NewEncoder(file)
	enc.Encode(ttomlState)
}

func LoadState(config Config) *State {
	var state State

	var ttomlState tomlState
	_, err := toml.DecodeFile("/opt/trigger/state/state.toml", &ttomlState)
	log.Print(err, "Unable to decode state file")

	for _, tTrigger := range ttomlState.Trigger {
		var trigger Trigger
		trigger.Id = tTrigger.Id
		trigger.Message = tTrigger.Message
		trigger.Queue = tTrigger.Queue
		trigger.GitHubSecret = tTrigger.GitHubSecret

		trigger.Init(config)

		state.Triggers = append(state.Triggers, &trigger)
	}

	return &state
}
