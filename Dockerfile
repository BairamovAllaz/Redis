FROM golang:latest
RUN go version
ENV GOPATH=/
COPY ./ ./
# build go app
RUN go mod download
RUN go build -o golangcrud ./cmd/main.go
CMD ["./Redis"]