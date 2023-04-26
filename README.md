# Chassit on Repeat
### The better and more simple ListenOnRepeat

This is a website for looping video with optional start/end times and tracking the amount of loops.


## API
The api doc is available in [API.md](API.md).


## Getting started running

### Run/Build from source
Copy the `.env.example` file and rename to `.env`,
in this file you write your MongoDB credentials to connect the site to your database.

#### RUN
```shell
go run main.go
```

#### Build
```shell
go build -o chassitonrepeat
```
```shell
./chassitonrepeat
```

### Docker

Fetch the image
```shell
docker pull ghcr.io/ravenholdt/chassitonrepeat:latest
```

#### Docker run
```shell
docker run -e MONGODB_URI=mongodb://user:password@url/ \
    -v ./files:/files \
    -p 8080 \
    ghcr.io/ravenholdt/chassitonrepeat:latest 
```

#### Docker compose

```yaml
version: '3'

services:
  repeat:
    image: ghcr.io/ravenholdt/chassitonrepeat:latest
    ports:
      - "8080"
    volumes:
      - "./files:/files"
    environment:
      MONGODB_URI: 'mongodb+srv://user:password@url'
    restart: unless-stopped
```

## Env variables

| Name          | Type   | Required | Default     | Description                                                                                 |
|---------------|--------|:--------:|-------------|---------------------------------------------------------------------------------------------|
| MONGODB_URI   | string |    ✓     |             | URI to connect to the database                                                              |
| PORT          | int    |    x     | 8080        | The port to use                                                                             |
| FILES_PATH    | string |    x     | "/files"    | The path to the directory of videos to use                                                  |
| CONFIG_PATH   | string |    x     | "/config"   | The path to the directory of configs, used by the overrides                                 |
| ENABLE_PROXY  | bool   |    x     | false       | Enable this is the server is running behind a proxy                                         |
| TRUSTED_PROXY | string |    x     | "127.0.0.1" | The ip of the proxy                                                                         |
| LOG_LEVEL     | string |    x     | "info"      | The log level used, valid entries [trace, debug, info, warn, error, fatal, panic, disabled] |

## Files
This folder holds all the videos used by the app.
New files are auto-detected when put into this folder or removed.

## Overrides
By adding a files named `overrides.csv` in the `CONFIG_PATH` folder you can override random videos.
This file is auto reloaded when updated.

The format is the following:
```csv
key_to_override,new_key
other_key,other_new_key
```
The keys are the value in the url after the `/` when viewing a video, for example for this url `http://localhost:8080/Y9EKzvTo3g0` the key is `Y9EKzvTo3g0`

## License
ISC