# GCAL-CLI

```shell
 ██████╗  ██████╗ █████╗ ██╗       ██████╗██╗     ██╗
██╔════╝ ██╔════╝██╔══██╗██║      ██╔════╝██║     ██║
██║  ███╗██║     ███████║██║█████╗██║     ██║     ██║
██║   ██║██║     ██╔══██║██║╚════╝██║     ██║     ██║
╚██████╔╝╚██████╗██║  ██║███████╗ ╚██████╗███████╗██║
 ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝  ╚═════╝╚══════╝╚═╝
```

a command line interface for printing events from user google calendar

## requirements

- go installed in your machine

## installation

via **go install**

```shell
go install github.com/agungfir98/gcal-cli@latest

```

## setup

Run the following command to create `./config/gcal-cli` in your user home directory, this directory is where you can store the credentials.json and the oauth2 token

```shell
gcal-cli config
```

### obtain the credentials

Google calendar apis require user to setup a Google Cloud Project and obtain the credentials to be able to access the api, I won't provide the google cloud platform credentials so you may obtain it by following the documentation ->
[setup gcp](https://developers.google.com/calendar/api/quickstart/go)

once you obtain the credentials.json you have to put it in the `(your home dir)/.config/gcal-cli`

## Authenticating

in order for you to authenticate is by execute any command that access the google calendar api like `calw` or `events`, it works by reading the `token.json` in your config file. When you execute the command for the first time it will give you an auth url to the google auth consent screen. once you done the consent screen don't immediately close it, you need to copy the token params in the url and paste it to the terminal.

## Usage

```shell
Available Commands:
  calw        get a week event calendar
  completion  Generate the autocompletion script for the specified shell
  config      create a .config/ directory at user home directory to store credentials and token json
  events      shows event in your google calendar
  help        Help about any command

Flags:
  -h, --help      help for gcal-cli
  -v, --version   version for gcal-cli

```
