[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/korvike)](https://goreportcard.com/report/github.com/Luzifer/korvike)
![](https://badges.fyi/github/license/Luzifer/korvike)
![](https://badges.fyi/github/downloads/Luzifer/korvike)
![](https://badges.fyi/github/latest-release/Luzifer/korvike)

# Luzifer / korvike

`korvike` is the finnish translation to the word "replacer" and that is what it does: It takes a Go template and executes it.

## Available functions

- `{{ .<variable name> }}`  
  Take key-value pairs from the CLI and replace them inside the template
  ```console
  $ echo "{{ .foo }}" | korvike -v foo=bar
  bar
  ```
- `{{ b64decode <string> }}`  
  Decodes the string with base64 [StdEncoding](https://golang.org/pkg/encoding/base64/#pkg-variables)
  ```console
  $ echo '{{ b64decode "SGVsbG8gV29ybGQ=" }}' | korvike
  Hello World
  ```
- `{{ b64encode <string> }}`  
  Encodes the string with base64 [StdEncoding](https://golang.org/pkg/encoding/base64/#pkg-variables)
  ```console
  $ echo '{{ b64encode "Hello World" }}' | korvike
  SGVsbG8gV29ybGQ=
  ```
- `{{ env <variable name> [default value] }}`  
  Read environment variables and replace them inside the template
  ```console
  $ export FOO=bar
  $ echo '{{ env "FOO" }}' | korvike
  bar
  ```
- `{{ file <file name> [default value] }}`  
  Read a file and place it inside the template
  ```console
  $ echo "Hello World" > hello
  $ echo '{{ file "hello" }}' | korvike
  Hello World
  ```
- `{{ hash <algo> <string> }}`  
  Hash string with given algorithm (supported algorithms: md5, sha1, sha256, sha512)
  ```console
  $ echo '{{ hash "sha256" "this is a test" }}' | korvike
  2e99758548972a8e8822ad47fa1017ff72f06f3ff6a016851f45c398732bc50c
  ```
- `{{ markdown <source> }}`  
  Format the source using a markdown parser
  ```console
  $ echo '{{ markdown "# headline" }}' | korvike
  <h1>headline</h1>
  ```
- `{{ now <format string> }}`  
  Format the current date into the template (uses [Go time format](https://golang.org/pkg/time/#Time.Format))
  ```console
  $ echo '{{ now "2006-01-02 15:04:05" }}' | korvike
  2017-04-17 16:27:34
  ```
- `{{ tplexec (file "my.tpl") }}`  
  Execute the given template with the same function set and variables as the parent template.
  ```console
  $ export FOO=bar
  $ echo '{{ env "FOO" }}' >my.tpl
  $ echo '{{ tplexec (file "my.tpl") }}' | korvike
  bar
  ```
- `{{ urlescape <input string> }}`  
  Do an URL escape to use the input string inside an query parameter in an URL
  ```console
  $ echo '{{ urlescape "Hellö Wörld@Golang" }}' | korvike 
  Hell%C3%B6+W%C3%B6rld%40Golang
  ```
- `{{ vault <path> <key> [default value] }}`  
  Read a key from Vault using `VAULT_ADDR` and `VAULT_TOKEN` environment variables (or `~/.vault-token` file) for authentication.
  ```console
  $ vault write secret/test foo=bar
  $ echo '{{ vault "secret/test" "foo" }}' | korvike
  bar
  ```
