# static-file-server

<a href="https://github.com/sponsors/halverneus" style="background-color:#fff;color:#000;padding:3px 12px;font-size:12px;border-color:#000;border:1px solid;border-radius:6px;box-sizing:border-box;line-height:20px;display:inline-block;">
<svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" style="vertical-align: middle; margin-right:4px;color:#f00;">
<path fill-rule="evenodd" style="color:#f00;fill: currentColor;" d="M4.25 2.5c-1.336 0-2.75 1.164-2.75 3 0 2.15 1.58 4.144 3.365 5.682A20.565 20.565 0 008 13.393a20.561 20.561 0 003.135-2.211C12.92 9.644 14.5 7.65 14.5 5.5c0-1.836-1.414-3-2.75-3-1.373 0-2.609.986-3.029 2.456a.75.75 0 01-1.442 0C6.859 3.486 5.623 2.5 4.25 2.5zM8 14.25l-.345.666-.002-.001-.006-.003-.018-.01a7.643 7.643 0 01-.31-.17 22.075 22.075 0 01-3.434-2.414C2.045 10.731 0 8.35 0 5.5 0 2.836 2.086 1 4.25 1 5.797 1 7.153 1.802 8 3.02 8.847 1.802 10.203 1 11.75 1 13.914 1 16 2.836 16 5.5c0 2.85-2.045 5.231-3.885 6.818a22.08 22.08 0 01-3.744 2.584l-.018.01-.006.003h-.002L8 14.25zm0 0l.345.666a.752.752 0 01-.69 0L8 14.25z"></path>
</svg>
<span style="color:#000">Buy me a Smoothie</span>
</a>

## Introduction

Tiny, simple static file server using environment variables for configuration.
Install from any of the following locations:

- Docker Hub: https://hub.docker.com/r/halverneus/static-file-server/
- GitHub: https://github.com/halverneus/static-file-server

## Configuration

### Environment Variables

Default values are shown with the associated environment variable.

```bash
# Enables resource access from any domain.
CORS=false

# Enable debugging for troubleshooting. If set to 'true' this prints extra
# information during execution. IMPORTANT NOTE: The configuration summary is
# printed to stdout while logs generated during execution are printed to stderr.
DEBUG=false

# Optional Hostname for binding. Leave black to accept any incoming HTTP request
# on the prescribed port.
HOST=

# If assigned, must be a valid port number.
PORT=8080

# Automatically serve the index file for a given directory (default). If set to
# 'false', URLs ending with a '/' will return 'NOT FOUND'.
SHOW_LISTING=true

# Folder with the content to serve.
FOLDER=/web

# URL path prefix. If 'my.file' is in the root of $FOLDER and $URL_PREFIX is
# '/my/place' then file is retrieved with 'http://$HOST:$PORT/my/place/my.file'.
URL_PREFIX=

# Paths to the TLS certificate and key. If one is set then both must be set. If
# both set then files are served using HTTPS. If neither are set then files are
# served using HTTP.
TLS_CERT=
TLS_KEY=

# If TLS certificates are set then the minimum TLS version may also be set. If
# the value isn't set then the default minimum TLS version is 1.0. Allowed
# values include "TLS10", "TLS11", "TLS12" and "TLS13" for TLS1.0, TLS1.1,
# TLS1.2 and TLS1.3, respectively. The value is not case-sensitive.
TLS_MIN_VERS=

# List of accepted HTTP referrers. Return 403 if HTTP header `Referer` does not
# match prefixes provided in the list.
# Examples:
#   'REFERRERS=http://localhost,https://...,https://another.name'
#   To accept missing referrer header, add a blank entry (start comma):
#   'REFERRERS=,http://localhost,https://another.name'
REFERRERS=
```

### YAML Configuration File

YAML settings are individually overridden by the corresponding environment
variable. The following is an example configuration file with defaults. Pass in
the path to the configuration file using the command line option
('-c', '-config', '--config').

```yaml
cors: false
debug: false
folder: /web
host: ""
port: 8080
referrers: []
show-listing: true
tls-cert: ""
tls-key: ""
tls-min-vers: ""
url-prefix: ""
```

Example configuration with possible alternative values:

```yaml
debug: true
folder: /var/www
port: 80
referrers:
    - http://localhost
    - https://mydomain.com
```

## Deployment

### Without Docker

```bash
PORT=8888 FOLDER=. ./serve
```

Files can then be accessed by going to http://localhost:8888/my/file.txt

### With Docker

```bash
docker run -d \
    -v /my/folder:/web \
    -p 8080:8080 \
    halverneus/static-file-server:latest
```

This will serve the folder "/my/folder" over http://localhost:8080/my/file.txt

Any of the variables can also be modified:

```bash
docker run -d \
    -v /home/me/dev/source:/content/html \
    -v /home/me/dev/files:/content/more/files \
    -e FOLDER=/content \
    -p 8080:8080 \
    halverneus/static-file-server:latest
```

### Getting Help

```bash
./serve help
# OR
docker run -it halverneus/static-file-server:latest help
```
