FROM golang:1.14-alpine AS builder
ADD . /poster
WORKDIR /poster
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o poster cmd/poster/*.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /poster ./
RUN chmod +x poster
CMD ./poster
