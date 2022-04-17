FROM golang:1.18-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o server

FROM alpine
#RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server .
# Run the web service on container startup.
CMD ["/server"]
EXPOSE 8080