FROM golang:latest
COPY . /src
RUN cd /src &&\
    go build -o /app cmd/main.go &&\
    cd / &&\
    rm -rf /src
WORKDIR /

# TODO @kris-nova set these from the systemd unit file
ENV KEY value

CMD [ "/app" ]