FROM golang:1.10.3-stretch

LABEL MAINTAINER="Programm Eins <programm.eins@gmail.com>"

ARG BINDIR="/opt/cmd"

ADD bin/cmd/ $BINDIR/bin/cmd/
ADD public/ $BINDIR/public/
ADD cmd-webui $BINDIR/cmd-webui

ENV LD_LIBRARY_PATH $BINDIR/bin/cmd/

WORKDIR $BINDIR

CMD ["./cmd-webui"]
