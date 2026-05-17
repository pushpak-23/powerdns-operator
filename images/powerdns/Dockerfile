ARG BASE_IMAGE=powerdns/pdns-auth-48:latest
FROM ${BASE_IMAGE}

USER root
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 53/tcp 53/udp 8081/tcp

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]