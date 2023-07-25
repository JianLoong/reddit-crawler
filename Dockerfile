FROM golang:1.20.0 AS builder

WORKDIR /app/

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download

COPY main.go .

# This method of building is needed
RUN CGO_ENABLED=1 go build -a -installsuffix cgo -o main .

# # Build final image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

# RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main ./

# CMD [ "./main" ]
ENTRYPOINT ["./main"]

