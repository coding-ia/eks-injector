FROM alpine:3.19

ARG USER=gouser

RUN adduser -D $USER

COPY eks-inject /

USER $USER

CMD ["/eks-inject"]
