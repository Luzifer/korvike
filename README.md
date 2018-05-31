![license](https://badges.fyi/github/license/Luzifer/badge-gen)
[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/korvike)](https://goreportcard.com/report/github.com/Luzifer/korvike)

# Luzifer / korvike

`korvike` is the finnish translation to the word "replacer" and that is what it does: It takes a Go template and executes it.

## Available functions

- `{{ .<variable name> }}`  
  Take key-value pairs from the CLI and replace them inside the template
  ```bash
  # echo "{{ .foo }}" | korvike -v foo=bar
  bar
  ```
- `{{ b64encode <string> }}`  
  Encodes the string with base64 [StdEncoding](https://golang.org/pkg/encoding/base64/#pkg-variables)
  ```console
  # echo '{{ b64encode "Hello World" }}' | korvike
  SGVsbG8gV29ybGQ=
  ```
- `{{ env <variable name> [default value] }}`  
  Read environment variables and replace them inside the template
  ```bash
  # export FOO=bar
  # echo '{{env "FOO"}}' | korvike
  bar
  ```
- `{{ file <file name> [default value] }}`  
  Read a file and place it inside the template
  ```bash
  # echo "Hello World" > hello
  # echo '{{file "hello"}}' | korvike
  Hello World
  ```
- `{{ now <format string> }}`  
  Format the current date into the template (uses [Go time format](https://golang.org/pkg/time/#Time.Format))
  ```bash
  # echo '{{now "2006-01-02 15:04:05"}}' | korvike
  2017-04-17 16:27:34
  ```
- `{{ vault <path> <key> [default value] }}`  
  Read a key from Vault using `VAULT_ADDR` and `VAULT_TOKEN` environment variables (or `~/.vault-token` file) for authentication.
  ```bash
  # vault write secret/test foo=bar
  # echo '{{vault "secret/test" "foo"}}' | korvike
  bar
  ```

----

![project status](https://d2o84fseuhwkxk.cloudfront.net/korvike.svg)
