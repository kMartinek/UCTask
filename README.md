# UCTask

Project consists of two parts:
  - application "app.js" written in NodeJS
  - service "service.go" written in golang
  
## Setup
First, start nats-server via docker.
```
docker run -d -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats-server nats:latest
```

Then run `service.go`

In the `app` folder run `npm ci` to get required node_modules
```
npm ci
```

And finally run `app.js` with path argument.
```
node app.js --path="full_file_path_here"
```

