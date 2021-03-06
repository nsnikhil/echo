FROM golang:alpine as builder
WORKDIR /echo
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o echo *.go

FROM scratch
COPY --from=builder /echo/echo .
ENTRYPOINT ["./echo"]