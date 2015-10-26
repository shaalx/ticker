FROM golang

WORKDIR /gopath/app
ENV GOPATH /gopath/app
ADD . /gopath/app/

RUN go get github.com/shaalx/ticker
RUN go build -o ticker

EXPOSE 80
CMD ["/gopath/app/ticker"]
