# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-digimap-backend"
LABEL REPO="https://github.com/hhung06/digimap-backend"

ENV PROJPATH=/go/src/github.com/hhung06/digimap-backend

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/hhung06/digimap-backend
WORKDIR /go/src/github.com/hhung06/digimap-backend

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/hhung06/digimap-backend"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/digimap-backend/bin

WORKDIR /opt/digimap-backend/bin

COPY --from=build-stage /go/src/github.com/hhung06/digimap-backend/bin/digimap-backend /opt/digimap-backend/bin/
RUN chmod +x /opt/digimap-backend/bin/digimap-backend

# Create appuser
RUN adduser -D -g '' digimap-backend
USER digimap-backend

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/digimap-backend/bin/digimap-backend"]
