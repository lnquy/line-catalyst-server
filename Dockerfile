FROM golang:1.19.12-alpine3.18 AS go_builder

WORKDIR /app
COPY . /app
RUN rm -f crm-core*.bin

RUN go build -v -o catalyst .

FROM alpine:3.18

WORKDIR /app

RUN mkdir -p /app/_misc

COPY --from=go_builder /app/catalyst /app/catalyst
COPY _misc/jokes.json /app/_misc/jokes.json
#COPY _misc/zoneinfo.zip /app/_misc/zoneinfo.zip
RUN chmod +x /app/catalyst

#ENV TZ=/app/_misc/zoneinfo.zip

CMD ["/app/catalyst"]