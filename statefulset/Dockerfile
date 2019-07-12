FROM golang:latest
WORKDIR $GOPATH/src/github.com/rodrigo-albuquerque/app
RUN apt-get update && apt-get upgrade -y && apt-get install golang -y
COPY app-pet-v2.go .
RUN go get -d -v ./... 
RUN go install -v ./...
EXPOSE 8080
ENTRYPOINT ["app"]
