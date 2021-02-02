## We specify the base image we need for our
## go application
FROM golang:buster
RUN mkdir /app
ADD . /app
WORKDIR /app
ENV TEST_SERVER=http://host.docker.internal:8080
RUN go mod download
CMD ["go", "test", "./...", "-v"]
