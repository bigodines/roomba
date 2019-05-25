## Roomba

*"The annoying bot that keeps the house clean"*

### About this project

Roomba is a simple bot that queries GitHub and post relevant pending PullRequests to a slack room. Everything is configurable through environment variables and `yml` files.


### Quick start

* [Create a GitHub access token](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line)

* [Create a Slack incoming webhook](https://get.slack.help/hc/en-us/articles/115005265063-Incoming-WebHooks-for-Slack)

* Clone this repository

* Create a `development.yml` based off of `default.yml`:

`cp config/default.yml config/development.yml`

* Edit `development.yml` with relevant information (github org, repos, slack channel id and webhook)

* Compile and run:

`GITHUB_TOKEN=<your github token> make build && ./roomba`

### Dev

* Install linter

`GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0`
