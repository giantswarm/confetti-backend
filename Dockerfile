FROM quay.io/giantswarm/alpine:3.12

USER root

RUN apk add --no-cache ca-certificates

RUN mkdir -p /opt
ADD ./confetti-backend /opt/confetti-backend

USER giantswarm

EXPOSE 7777
ENTRYPOINT ["/opt/confetti-backend"]