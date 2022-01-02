FROM alpine:edge AS build

RUN apk update
RUN apk upgrade
RUN apk add --update go=1.17.4-r0 git make build-base

WORKDIR /app
ADD ./ /app
RUN go build

FROM alpine:edge
WORKDIR /app
RUN cd /app
COPY --from=build /app/mcytbot /app/mcytbot
CMD ["./mcytbot"]
