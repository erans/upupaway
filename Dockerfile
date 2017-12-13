FROM alpine:3.7
COPY bin/upupaway /upupaway
EXPOSE 8000
ENTRYPOINT ["/upupaway"]
