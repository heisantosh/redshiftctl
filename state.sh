#/bin/bash

state=$(redshiftctl get state)
if [ "${state}" = "on" ]; then
    temp=$(redshiftctl get temperature)
    echo "%{F#EC7875} %{F-}${temp}"
else
    echo "%{F#42A5F5} %{F-}"
fi