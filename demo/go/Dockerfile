FROM iron/go:dev

WORKDIR /app

ENV SRC_DIR=/go/src/github.com/lmolkova/go-app/
ADD . $SRC_DIR

RUN go get -u go.opencensus.io
RUN go get -u contrib.go.opencensus.io/exporter/ocagent
RUN go get -u golang.org/x/time/rate
RUN go get -u github.com/prometheus/client_golang/prometheus
RUN go get -u github.com/prometheus/client_golang/prometheus/promhttp
RUN cd $SRC_DIR; go build -o goapp; cp goapp /app/

EXPOSE 8080

ENTRYPOINT ["./goapp"]