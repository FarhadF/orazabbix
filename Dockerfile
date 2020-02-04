FROM oracle/instantclient:19

RUN yum install -y gzip tar git wget gcc libstdc++.x86_64 compat-libstdc++-33.x86_64 gcc-c++.x86_64
RUN wget https://dl.google.com/go/go1.13.7.linux-amd64.tar.gz
RUN mkdir -p /usr/local/go && tar xzf go1.13.7.linux-amd64.tar.gz -C /usr/local/
ADD . /root/go/src/orazabbix/
WORKDIR /root/go/src/orazabbix/
ADD oci8.pc /usr/lib/oracle/19.5/client64/lib/
RUN PATH=$PATH:/usr/local/go/bin && \
export LD_LIBRARY_PATH=/usr/lib/oracle/19.5/client64/lib/ && \
export PKG_CONFIG_PATH=/usr/lib/oracle/19.5/client64/lib/ && \
export CGO_CFLAGS="-I/usr/include/oracle/19.5/client64/" && \
export CGO_LDFLAGS="-L/usr/lib/oracle/19.5/client64/lib/ -lstdc++ -lclntsh"  && \
go mod vendor && go build -mod=vendor -tags noPkgConfig -o /orazabbix main.go
ADD orazabbix.sh /

CMD ["/orazabbix.sh"]
