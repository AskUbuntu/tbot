## Ask Ubuntu Chat Twitter Bot

A small web application that watches the Ask Ubuntu general room for messages that meet a certain criteria and allows them to be tweeted at regular intervals. Some of the features that will eventually make their way into this app include:

- Customize criteria for discovering "relevant" messages
- Hand-select messages for tweeting
- Allow multiple users to manage tweets
- Send custom messages from the Twitter account
- Queue messages for tweeting at regular intervals

Due to its lengthy name, the app has been nicknamed "tbot".

### Screenshots

You can see the application in action here:
https://tbot.quickmediasolutions.com

![](http://i.stack.imgur.com/sdCmW.png)

*Home page displaying a single tweet*

### Building tbot

tbot is written in the Go programming language. It will run on any platform supported by Go and compiles to a single binary. Assuming you have Go installed and configured (`$GOPATH` is set), you can run:

    go get github.com/AskUbuntu/tbot

...and the resulting binary will be in `$GOPATH/bin`.

### Running the Application

Configuration is split into two categories:

- Information needed to start the application (address, etc.)
- Settings that control application behavior (search terms, etc.)

The information needed to start the application is provided in a JSON file as a single argument. An example configuration is provided in `config.json.default`. Copy this file to `config.json` and customize it as necessary. Then launch the application with:

    $(GOPATH)/bin/tbot config.json

Assuming the default configuration, point your browser to http://127.0.0.1:8000 and you're good to go!

### Running in Docker

To run tbot in Docker, pull the `askubuntu/tbot` image:

    docker pull askubuntu/tbot

Then run the container with:

    docker \
        run \
        -d \
        -p 8000:8000 \
        --env ADMIN_PASSWORD="" \
        --env SECRET_KEY="" \
        --env TWITTER_CONSUMER_KEY="" \
        --env TWITTER_CONSUMER_SECRET="" \
        --env TWITTER_ACCESS_TOKEN="" \
        --env TWITTER_ACCESS_SECRET="" \
        -v /data/tbot:/path/for/data \
        askubuntu/tbot

Two important things to note:

1. Each of the environment variables listed in the command above must have a valid value set.
2. To persist settings and data, mount a volume at `/data/tbot`.
