FROM centurylink/ca-certs
COPY bin/upupaway /upupaway
EXPOSE 8000
ENTRYPOINT ["/upupaway"]
