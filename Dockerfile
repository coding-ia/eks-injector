FROM alpine:3.19

COPY eks-inject /

CMD ["/eks-inject"]
