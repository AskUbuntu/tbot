FROM alpine:3.3
MAINTAINER Nathan Osman <nathan@quickmediasolutions.com>

WORKDIR /opt/tbot

# Copy the binary, data files, and default configuration to the container
COPY tbot .
COPY www/ www/
COPY config.json.default config.json

# Indicate the command to run
CMD ./tbot config.json

# Expose the default port
EXPOSE 8000
