# Reddit Crawler

Simple Reddit parser using Reddit JSON API written in Go

# Endpoints

/api/askreddit/indexes.json


# Running this Project

1. Via GO installation

```
# GO111MODULE should be switched on
# On Windows via
go env -w GO111MODULE=on 

go mod tidy

# Running
go run main.go [subreddit-name] [no-of-post]
```

2. Via Docker Compose

```
docker compose build
docker compose up
```