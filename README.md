# OGame bot

Simple OGame bot that connects to your account and interacts with it. It uses some randomness to hide its bot face.

## Motivation

Simply for fun. I'm a former player of this addictive online game. And it was also an interesting challenge to kind of reverse engineer the login steps (with multiple redirection) knowing nothing about the web/http/... at the time :)

Still a lot to do! Feel free to join my efforts!

## Installation

- It requires Go language of course. You can set it up by downloading it here: https://golang.org/dl/
- Use go get or download the files directly from github to get the project
- Set your GOPATH (to the project location) and GOROOT (where Go is installed) environment variables.

## Build and usage

```
@ogbot $ go install ogbot
@ogbot $ bin/ogbot.exe -help
@ogbot $ bin/ogbot.exe -lang=fr -log=ogbot.log -login=YOUR_LOGIN -pass=YOUR_PWD -uni=YOUR_UNI -dump
```