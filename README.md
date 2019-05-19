## Roomba

*"The annoying bot that keeps the house clean"*

### About this project

Roomba is a simple bot that queries GitHub and post relevant pending PullRequests to a slack room. Everything is configurable through environment variables and `yml` files. ~You don't even need `Go`. Just download one of the release binaries and start using it.~ (not yet)


### Quick start

* [Create a GitHub access token](https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line)

* [Create a Slack incoming webhook](https://get.slack.help/hc/en-us/articles/115005265063-Incoming-WebHooks-for-Slack)

* Clone this repository

* Create a `development.yml` based off of `default.yml`:

`cp config/default.yml config/development.yml`

* Edit `development.yml` with relevant information (github org, repos, slack channel id and webhook)

* Compile and run:

`GITHUB_TOKEN=<your github token> make build && ./roomba`


### Roadmap

Roomba is a weekend project thus development might be slow but here are a few things I plan to work on:
* Add tests (seriously, I'm working on it)
* Turn roomba into a slack app so users can interact with it
* Message users as their PRs get reviewed
