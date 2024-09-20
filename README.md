# URLShortner

Simple URLShortner written in GoLang
Day 2 of Gophercises

## Usage

Running the binary starts a server on port 8080

Route details

## shorten

Accepted Method(s): POST

Takes in a json body of the following format:
```json
{
	"name": "UrlAlias",
	"url": "https://example.com"
}
```

Returns shortned url

## go

Accepted Method(s): GET

If name exists redirects to the url or else returns that name is invalid

Example: localhost:8080/go/UrlAlias will redirect you to https://example.com
