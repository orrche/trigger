FROM fedora

RUN dnf update -y
RUN mkdir -p /opt/trigger/state

ADD config.toml /opt/trigger/
ADD trigger /opt/trigger/
RUN adduser trigger

EXPOSE 8080

ENTRYPOINT ["/opt/trigger/trigger"]
