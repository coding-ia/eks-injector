FROM scratch

COPY eks-inject /

CMD ["/eks-inject"]
