# huproxy

A golang version of the hue-remote repo since it isn't all that necessary now

## Setup

Pull down the repo, `git clone https://github.com/yashdalfthegray/huproxy.git`

Run `go mod tidy` and `make build`

## Running

This application requires some environment variables to be set before it can do anything. Check the section below.

Once you have the environment variables and you've built the binary, run `make run` or just execute the binary directly by running `bin/huproxy`.

There are two endpoints exposed by this thing

`/ping` will give you the status of the server

`/page` will make the hue lights specified by the `GROUPED_LIGHT_ID` blink between `START_COLOR` and `JUMP_COLOR` for `DURATION` seconds.

## Running under Docker

You can also run this thing as a Docker container. Use `docker build -t huproxy .` to build the container image and then use `docker run -d -p 9090:9090 --env-file .env --name myhuproxy huproxy:latest` to run it as a container.

## Environment variables

## Environment Variables

| Variable             | Description                                    | Default   | Required |
| -------------------- | ---------------------------------------------- | --------- | -------- |
| `HUE_BRIDGE_ADDRESS` | IP address of the Hue Bridge                   |           | Yes      |
| `GROUPED_LIGHT_ID`   | ID of the grouped light resource               |           | Yes      |
| `HUE_USERNAME`       | Username for accessing the Hue API             |           | Yes      |
| `START_COLOR`        | Starting color in hex format (e.g., `#ff5722`) | `#ff5722` | No       |
| `JUMP_COLOR`         | Jump color in hex format (e.g., `#ff0000`)     | `#ff0000` | No       |
| `DURATION_SECONDS`   | Duration of the effect in seconds              | `15`      | No       |
