FROM golang:1.19
COPY . /app
WORKDIR /app
RUN go build -o server .
EXPOSE 8000
ENTRYPOINT ["./server"]