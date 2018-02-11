FROM golang:1.8

WORKDIR /go/src/github.com/chrisfisher/jubilant-waffle/server
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"] # ["app"]

# docker build -t waffle .
# docker run -it --rm -p 8080:8080 waffle