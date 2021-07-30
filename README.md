# redshiftctl
Tool to manually control monitor color temperature.

## Installation
### Pre requisites
Make sure to have redshift installed for your Linux system.

### Using go
```bash
go install github.com/heisantosh/redshiftctl
```

## Usage
```text
redshiftctl (com.github.heisantosh.redshiftctl) 0.1.0

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
  Configuration file is ~/.config/redshiftctl/config.json

  Keys are
  # current state on, off
  state=on
  # color temperature
  temperature=4500

```

### Examples

#### Toggle current state
```bash
redshiftctl toggle
```
#### Toggle state on
```bash
redshiftctl toggle on
```

#### Set color temperature
```bash
redshiftctl set 4500
```

#### Decrease color temperature
```bash
redshiftctl decrease 500
```

#### Get the current color temperature
```bash
redshiftctl get temperature
```

## Polybar 

This was initially designed for use in a polybar custom script module. The goal was to support mappings for scroll up and scroll down actions to decrease and increase color temperatures respectively as well as toggling redshift on/off. 

### redshift module
```ini
[module/redshift]
type = custom/ipc
# When starting, restore state from last time.
hook-0 = redshiftctl load && ~/.config/polybar/default/bin/redshift/state.sh
# Toggle state on/off.
hook-1 = redshiftctl toggle && ~/.config/polybar/default/bin/redshift/state.sh
# Increase color temperature.
hook-2 = redshiftctl toggle on && redshiftctl increase 500 && ~/.config/polybar/default/bin/redshift/state.sh
# Decrease color temperature.
hook-3 = redshiftctl toggle on && redshiftctl decrease 500 && ~/.config/polybar/default/bin/redshift/state.sh

initial = 1
# Toggle on/off on left click.
click-left = polybar-msg -p %pid% hook redshift 2
# Increase color temperature on scroll up.
scroll-up = polybar-msg -p %pid% hook redshift 3
# Decrease color temperature on scroll down.
scroll-down = polybar-msg -p %pid% hook redshift 4
```

### `state.sh` for module content 
```bash
#/bin/bash

state=$(redshiftctl get state)

if [ "${state}" = "on" ]; then
    temperature=$(redshiftctl get temperature)
    echo "%{F#EC7875} %{F-}${temperature}"
else
    echo "%{F#42A5F5} %{F-}"
fi
```