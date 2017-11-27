FROM centurylink/ca-certs
COPY bin/upupaway /upupaway
ENTRYPOINT ["/upupaway"]
