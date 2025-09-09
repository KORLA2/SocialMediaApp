FROM golang:1.23-alpine AS build 
WORKDIR /app
COPY go.mod go.sum ./  
RUN go mod download  
COPY . .
RUN go build -o social ./cmd/api


FROM alpine:3.20
WORKDIR /app
COPY --from=build app/social .
EXPOSE 8008
CMD ["./social"]


