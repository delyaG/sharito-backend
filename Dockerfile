# build the binary
FROM golang:1.15 AS build

RUN useradd -u 10001 gopher

# cache go deps
RUN mkdir /sharito
WORKDIR /sharito
COPY go.mod go.sum ./

# download dependencies if go.sum changed
RUN go mod download
COPY . .

RUN make build

# run the binary
FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

USER gopher

COPY --from=build /sharito/migrations /migrations
COPY --from=build /sharito/bin/sharito /sharito

EXPOSE $PORT
CMD ["./sharito"]
