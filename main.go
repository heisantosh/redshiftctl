package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

const (
	appId   = "com.github.heisantosh.redshiftctl"
	appName = "redshiftctl"
	version = "0.1.0"
)

const (
	// Redshift on/off state.
	stateOn  = "on"
	stateOff = "off"

	// Redshift state fields.
	stateTemperature = "temperature"
	stateState       = "state"

	// Supported commands.
	cmdHelp     = "help"
	cmdLoad     = "load"
	cmdIncrease = "increase"
	cmdDecrease = "decrease"
	cmdToggle   = "toggle"
	cmdSet      = "set"
	cmdGet      = "get"
)

func configDir() string {
	return os.Getenv("HOME") + "/.config/redshiftctl"
}

func configFile() string {
	return configDir() + "/config.json"
}

func help() {
	fmt.Println(appName + ` (` + appId + `) ` + version + `

Tool to manually control monitor color temperature using redshft.

USAGE
  redshiftctl COMMAND [ARG]

COMMANDS
  toggle [STATE]       toggle redshift to state on or off, if not provided toggele current state
  load                 load the state of the configuration file
  increase TEMP        increase the color temperature by TEMP
  decrease TEMP        decrease the color temperature by TEMP
  set TEMP             set the color temperature to TEMP
  get STATE            get the value of the state, STATE can be state or temperature
  help                 print this help information

CONFIGURATION
  Configuration file is ` + configFile() + `

  Keys are
  # current state on, off
  state=on
  # color temperature
  temperature=4500`)
}

type redshiftState struct {
	State       string `json:"state"`
	Temperature int    `json:"temperature"`
}

type cmdArgs struct {
	cmd              string
	toggleState      string
	temperatureDelta int
	temperature      int
	getState         string
}

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		os.Exit(1)
	}

	firstRunCheck()
	runCommand(args)
}

func firstRunCheck() {
	_, err := os.Stat(configFile())
	if os.IsNotExist(err) {
		os.MkdirAll(configDir(), 0744)
		stateStore(redshiftState{State: stateOn, Temperature: 4500})
	}
}

func parseArgs() (cmdArgs, error) {
	args := cmdArgs{}

	if len(os.Args) == 2 && (os.Args[1] == cmdHelp || os.Args[1] == cmdLoad || os.Args[1] == cmdToggle) {
		args.cmd = os.Args[1]
		return args, nil
	}

	if len(os.Args)-1 != 2 {
		return args, errors.New("insufficient args")
	}

	args.cmd = os.Args[1]

	switch args.cmd {
	case cmdDecrease, cmdIncrease:
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return args, errors.New("invalid " + os.Args[2] + " value")
		}
		args.temperatureDelta = v

	case cmdSet:
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return args, errors.New("invalid " + os.Args[2] + " value")
		}
		args.temperature = v

	case cmdGet:
		args.getState = os.Args[2]
		if !(args.getState == stateState || args.getState == stateTemperature) {
			return args, errors.New("invalid get arg" + args.getState)
		}

	case cmdToggle:
		args.toggleState = os.Args[2]
		if !(args.toggleState == stateOn || args.toggleState == stateOff) {
			return args, errors.New("invalid " + args.toggleState + " value")
		}

	default:
		return args, errors.New("unknown command " + args.cmd)
	}

	return args, nil
}

func stateLoad() (redshiftState, error) {
	state := redshiftState{}
	if b, err := ioutil.ReadFile(configFile()); err == nil {
		json.Unmarshal(b, &state)
	} else {
		return state, err
	}
	return state, nil
}

func stateStore(state redshiftState) {
	b, _ := json.Marshal(state)
	_ = os.WriteFile(configFile(), b, 0664)
}

func toggleOff() {
	exec.Command("redshift", "-o", "-x").Run()
}

func setTemperature(temp int) {
	exec.Command("redshift", "-P", "-o", "-O", strconv.Itoa(temp)).Run()
}

func runCommand(args cmdArgs) {
	state, _ := stateLoad()

	switch args.cmd {
	case cmdToggle:
		if (args.toggleState == "" && state.State == stateOn) || args.toggleState == stateOff {
			toggleOff()
			state.State = stateOff
		} else {
			setTemperature(state.Temperature)
			state.State = stateOn
		}

	case cmdIncrease:
		// TODO: upper limit check
		setTemperature(state.Temperature + args.temperatureDelta)
		state.Temperature = state.Temperature + args.temperatureDelta

	case cmdDecrease:
		// TODO: lower limit check
		setTemperature(state.Temperature - args.temperatureDelta)
		state.Temperature = state.Temperature - args.temperatureDelta

	case cmdSet:
		// TODO: limit check
		setTemperature(args.temperature)
		state.Temperature = args.temperature

	case cmdGet:
		switch args.getState {
		case stateState:
			fmt.Println(state.State)
		case stateTemperature:
			fmt.Println(state.Temperature)
		}

	case cmdLoad:
		if state.State == stateOff {
			toggleOff()
		} else {
			setTemperature(state.Temperature)
		}

	case cmdHelp:
		help()
		return
	}

	stateStore(state)
}
