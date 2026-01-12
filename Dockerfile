# Build
FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Run
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/server /server
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
