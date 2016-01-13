FROM alpine:latest
COPY build/quorra /bin/quorra
ADD build/public /bin/public
WORKDIR /bin
ENTRYPOINT ["/bin/quorra"]
EXPOSE 8080
