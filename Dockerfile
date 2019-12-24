FROM golang:1.10.3-stretch

LABEL MAINTAINER="Stephan Weintritt <45856463+sweintritt@users.noreply.github.com>"

ARG BINDIR="/opt/go-serve"

ADD build/ $BINDIR/

WORKDIR $BINDIR

CMD ["./go-serve"]
