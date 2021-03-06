FROM golang:1.12.2 as builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/madappgang/identifo
COPY . ./
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /identifo .
COPY server-config.yaml ./
COPY jwt/*.pem ./jwt/
COPY static ./static
COPY cmd/import/apps.json ./cmd/import/apps.json
COPY cmd/import/users.json ./cmd/import/users.json
COPY cmd/demo/init-config.yaml ./init-config.yaml

ENTRYPOINT ["./identifo"]