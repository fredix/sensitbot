FROM golang:1.9.2-alpine

RUN apk update && \
	apk add git

# Make the source code path
RUN mkdir -p /go/src/github.com/fredix/sensitbot

ADD . /go/src/github.com/fredix/sensitbot

RUN go get github.com/BurntSushi/toml && \
	go get github.com/julienschmidt/httprouter && \
	go get github.com/tucnak/telebot && \
	go get golang.org/x/oauth2 && \
	go get gopkg.in/mgo.v2

# Run the Go installer
RUN go build -o /go/bin/sensitbot /go/src/github.com/fredix/sensitbot/sensitbot.go

# Copy the application files (needed for production)
RUN mkdir /conf
ADD sensitbot.toml-example /conf/sensitbot.toml
ADD golang.jpg /go/bin/golang.jpg

# Expose the application on port 3000
EXPOSE 3000
ENTRYPOINT ["/go/bin/sensitbot", "/conf/sensitbot.toml"]
