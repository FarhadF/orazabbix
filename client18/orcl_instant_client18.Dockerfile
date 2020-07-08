# LICENSE UPL 1.0
#
# Copyright (c) 2014, 2019, Oracle and/or its affiliates. All rights reserved.
#
# ORACLE DOCKERFILES PROJECT
# --------------------------
#
# Dockerfile template for Oracle Instant Client
#
# HOW TO BUILD THIS IMAGE
# -----------------------
#
# Run:
#      $ docker build --pull -t oracle/instantclient:18 .
#
#
FROM oraclelinux:7-slim

ENV PATH=$PATH:/usr/lib/oracle/${release}.${update}/client64/bin

ARG release=18
ARG update=5

RUN yum -y install oracle-release-el7 && \
    yum-config-manager --enable ol7_oracle_instantclient && \
    yum -y install oracle-instantclient${release}.${update}-basic \
        oracle-instantclient${release}.${update}-devel \
        oracle-instantclient${release}.${update}-sqlplus && \
    rm -rf /var/cache/yum /var/lib/yum/yumdb/* /usr/lib/udev/hwdb.d/* && \
    echo /usr/lib/oracle/${release}.${update}/client64/lib > /etc/ld.so.conf.d/oracle-instantclient${release}.${update}.conf && \
    ldconfig

CMD ["sqlplus", "-v"]