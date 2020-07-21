FROM brgtsisdt3ptf001/oracleinstantclient:18 as build

ENV LD_LIBRARY_PATH=/usr/lib/oracle/18.5/client64/lib
ENV PKG_CONFIG_PATH=/usr/lib/oracle/18.5/client64/lib
ENV CGO_CFLAGS="-I/usr/include/oracle/18.5/client64"
ENV CGO_LDFLAGS="-L/usr/lib/oracle/18.5/client64/lib/ -lstdc++ -lclntsh"

ADD orazabbix.sh /
ADD . /root/go/src/orazabbix/
ADD oci8.pc /usr/lib/oracle/18.5/client64/lib/

WORKDIR /root/go/src/orazabbix/

RUN yum -y install \
  gzip \
  tar \
  wget \
  gcc \
  libstdc++.x86_64 \
  compat-libstdc++-33.x86_64 \
  gcc-c++.x86_64 && \
  yum -y clean all && \
  rm -rf /var/cache/yum /var/lib/yum/yumdb/* /usr/lib/udev/hwdb.d/* && \
  wget https://dl.google.com/go/go1.13.7.linux-amd64.tar.gz && \
  mkdir -p /usr/local/go && \
  tar xzf go1.13.7.linux-amd64.tar.gz -C /usr/local/ && \
  PATH=$PATH:/usr/local/go/bin && \
  go mod vendor && go build -mod=vendor -tags noPkgConfig -o /orazabbix main.go && \
  rm -rf go1.13.7.linux-amd64.tar.gz

FROM brgtsisdt3ptf001/oracleinstantclient:18

ENV LD_LIBRARY_PATH=/usr/lib/oracle/18.5/client64/lib
ENV PKG_CONFIG_PATH=/usr/lib/oracle/18.5/client64/lib
ENV CGO_CFLAGS="-I/usr/include/oracle/18.5/client64"
ENV CGO_LDFLAGS="-L/usr/lib/oracle/18.5/client64/lib/ -lstdc++ -lclntsh"

COPY --from=build /orazabbix /

ADD orazabbix.sh /
ADD oci8.pc /usr/lib/oracle/18.5/client64/lib/

RUN adduser orazabbix && \
    chown orazabbix:orazabbix /orazabbix /orazabbix.sh

USER orazabbix

CMD ["/orazabbix.sh"]
