FROM golang:1.19-alpine

WORKDIR .

COPY . .

ENV GOPATH=/

RUN go get

RUN go build main.go

CMD [ "./main" ]