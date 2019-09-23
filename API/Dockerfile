FROM golang:latest
WORKDIR "/app"
COPY . .
RUN go get
CMD ["go", "run", "."]