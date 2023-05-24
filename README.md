## Installation Guide

### Set environment variables.
```bash
# You can use the predefined database config.

cp .env.example .env

# Edit the database config (optional)

nano .env
```

### Start docker image for PostGIS db server.
```bash
docker-compose up
```

### Run the web server application directly from the project or build/run binary
```bash
# Run directly from the project

go run .

# Build the project and run binary

go build .
./outdoorsy
```

### Run test
```bash
go test
```