FROM  go-oci8:latest AS builder
ADD . /go/src/orazabbix/
WORKDIR /go/src/orazabbix/
#RUN dep ensure -vendor-only
#RUN go mod download

RUN go get ./... 
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /orazabbix
RUN go build -o /orazabbix
#RUN go build -o /orazabbix

FROM store/oracle/database-instantclient:12.2.0.1

COPY --from=builder /orazabbix /orazabbix
ADD ./docker/* /

CMD ["/orazabbix.sh"]
