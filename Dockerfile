FROM debian:stable-slim

COPY go-blog-aggregator /bin/go-blog-aggregator

ENV PORT 8080

CMD [ "/bin/go-blog-aggregator" ]