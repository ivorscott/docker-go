FROM golang:1.13

RUN mkdir /app

COPY . /app

WORKDIR /app

# compile the binary executable
RUN go build -o main ./cmd/web
# run the binary executable
CMD ["/app/main"]
