# Chassit on Repeat
### The better and more simple ListenOnRepeat

This is a website for looping video with optional start/end times and tracking the amount of loops.

## Getting started
Copy the `.env.example` file and rename to `.env`,
in this file you write your MongoDB credentials to connect the site to your database.

### Docker

```shell
docker pull ghcr.io/ravenholdt/chassitonrepeat:latest
```

Example [docker-compose.yml](docker-compose.yml) on how to run this.

### Env variables


| Name          | Type   | Required | Default     |
|---------------|--------|:--------:|-------------|
| MONGODB_URI   | string |    ✓     |             |
| PORT          | int    |    x     | 8080        |
| FILE_PATH     | string |    x     | "/files"    |
| CONFIG_PATH   | string |    x     | "/config"   |
| ENABLE_PROXY  | bool   |    x     | false       |
| TRUSTED_PROXY | string |    x     | "127.0.0.1" |
| LOG_LEVEL     | string |    x     | "info"      |


## License
ISC