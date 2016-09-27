FROM alpine:3.3
MAINTAINER Nathan Osman <nathan@quickmediasolutions.com>

WORKDIR /opt/tbot

# Add CA certificates to the container
RUN apk --update add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

# Copy the binary and data files to the container
COPY tbot .
COPY www/ www/

# Indicate the command to run
CMD ./tbot config.json

# Expose the default port
EXPOSE 8000
