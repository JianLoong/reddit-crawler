FROM golang:1.20.0 AS builder

WORKDIR /app/

COPY go.mod .
COPY go.sum .

COPY . .

RUN go mod download

COPY main.go .

# This method of building is needed
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o main .

# CMD ["/app/main"]


# # Build final image
FROM scratch

# RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main ./

EXPOSE 8080

CMD [ "./main" ]
