FROM alpine:3.20
COPY bin/upupaway /upupaway
EXPOSE 8000
ENTRYPOINT ["/upupaway"]
