FROM        scratch

COPY        imgserv /imgserv

EXPOSE      9000
ENTRYPOINT  [ "/imgserv" ]
