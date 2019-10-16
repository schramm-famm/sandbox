FROM golang:latest AS builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/basic_system
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .
COPY static /static

FROM scratch
WORKDIR /
COPY --from=builder /app ./
COPY --from=builder /static ./static
EXPOSE 8080
ENTRYPOINT ["/app"]
