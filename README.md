## Ask Ubuntu Chat Twitter Bot

A small web application that watches the Ask Ubuntu general room for messages that meet a certain criteria and allows them to tweeted at regular intervals. Some of the features that will eventually make their way into this app include:

- Customize criteria for discovering "relevant" messages
- Hand-select messages for tweeting
- Allow multiple users to manage tweets
- Send custom messages from the Twitter account
- Queue messages for tweeting at regular intervals

Due to its lengthy name, the app has been nicknamed "tbot".

### Building tbot

tbot is written in the Go programming language. It will run on any platform supported by Go and compiles to a single binary. Assuming you have Go installed and configured (`$GOPATH` is set), you can run:

    go get github.com/AskUbuntu/tbot

...and the resulting binary will be in `$GOPATH/bin`.
