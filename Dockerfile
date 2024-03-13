FROM alpine:3.19

ARG USER=gouser

RUN adduser -D $USER

COPY eks-injector /

USER $USER

CMD ["/eks-injector"]
