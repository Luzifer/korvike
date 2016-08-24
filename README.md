![license](https://badges.fyi/github/license/Luzifer/badge-gen)
[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/korvike)](https://goreportcard.com/report/github.com/Luzifer/korvike)
[![Download on GoBuilder.me](https://badges.fyi/static/Download on/GoBuilder.me)](https://gobuilder.me/github.com/Luzifer/korvike)

# Luzifer / korvike

`korvike` is the finnish translation to the word "replacer" and that is what it does: It takes a Go template and executes it.

## Available functions

- Take key-value pairs from the CLI and replace them inside the template:  
`echo "{{ .foo }}" | korvike -v foo=bar => "bar"`
- Read environment variables and replace them inside the template:  
`export FOO=bar; echo '{{env "FOO"}}' | korvike => "bar"`
- Read a file and place it inside the template:  
`echo "Hello World" > hello; echo '{{file "hello"}}' | korvike => "Hello World"`
