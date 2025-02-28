FROM golang:1.24-alpine AS build

RUN apk update && apk upgrade --no-cache
RUN apk add --no-cache --update go gcc g++
WORKDIR /src
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build

FROM alpine:3.21
ARG USER=gouser

RUN adduser -D $USER
COPY --from=build /src/eks-injector /eks-injector

USER $USER

CMD ["/eks-injector"]