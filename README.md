# denon-to-mqtt-gw

[![mqtt-smarthome](https://img.shields.io/badge/mqtt-smarthome-blue.svg)](https://github.com/mqtt-smarthome/mqtt-smarthome)

Maintain a topic with volume and power state. 
Tested with AVR X1200W

# Messages

## Example message

```json
{
  "power":"ON",
  "volume":24.5,
  "volumeMax": 98.0
}
```

## Example configuration

```json
{
  "mqtt": {
    "url": "tcp://192.168.0.1:1883",
    "retain": true,
    "topic": "home/denon",
    "qos": 2
  },
  "denon": {
    "ip": "127.0.0.1"
  }
}
```

# Bridge status

The bridge maintains a status topic:

## Topic: `.../bridge/state`

| Value     | Description                          |
| --------- | ------------------------------------ |
| `online`  | The bridge is started                |
| `offline` | The bridge is currently not started. |

# run

Copy the `config-example.json` to `/production/config/config.json`

```
cd ./production
docker-compose up -d
```

# build

## Build container

Build the docker container using `build.sh`.
