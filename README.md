# static-file-server
Tiny, simple static file server using environment variables for configuration

Available on Docker Hub at https://hub.docker.com/r/halverneus/static-file-server/

Environment variables with defaults:
```
HOST=
PORT=8080
FOLDER=/web
```

### Without Docker
```
PORT=8888 FOLDER=. ./serve
```
Files can then be accessed by going to http://localhost:8888/my/file.txt

### With Docker
docker run -d -v /my/folder:/web -e PORT=9090 -p 9090:9090 halverneus/static-file-server:latest

This will serve the folder "/my/folder" over http://localhost:9090/my/file.txt
