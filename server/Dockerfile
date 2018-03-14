FROM golang:1.10

RUN go get github.com/pilu/fresh

WORKDIR /go/src/github.com/chrisfisher/jubilant-waffle/server
COPY . .

RUN go get
RUN go build

CMD ["fresh"]

# docker build -t waffle .
# docker run -it --rm -p 8080:8080 waffle