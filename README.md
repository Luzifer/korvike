![](https://badges.fyi/github/license/Luzifer/korvike)
![](https://badges.fyi/github/downloads/Luzifer/korvike)
![](https://badges.fyi/github/latest-release/Luzifer/korvike)

# Luzifer / korvike

`korvike` is the Finnish translation to the word "replacer" and that is what it does: It takes a Go template and executes it.

## Available functions

Starting with `v1.0.0` Korvike is based on the [sprig functions collection](https://masterminds.github.io/sprig/) with some additions:

- `{{ .<variable name> }}`  
  Take key-value pairs from the CLI and replace them inside the template
  ```console
  $ echo "{{ .foo }}" | korvike -v foo=bar
  bar
  ```
- `{{ file <file name> }}` / `{{ mustFile <file name> }}`  
  Read a file and place it inside the template, `file` returns an empty string on error, `mustFile` an error
  ```console
  $ echo "Hello World" > hello
  $ echo '{{ file "hello" }}' | korvike
  Hello World
  ```
- `{{ markdown <source> }}`  
  Format the source using a markdown parser
  ```console
  $ echo '{{ markdown "# headline" }}' | korvike
  <h1>headline</h1>
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
- `{{ vault <path> <key> }}` / `{{ mustVault <path> <key> }}`  
  Read a key from Vault using `VAULT_ADDR` and `VAULT_TOKEN` environment variables (or `~/.vault-token` file) for authentication. `vault` returns an empty string on error, `mustVault` an error
  ```console
  $ vault write secret/test foo=bar
  $ echo '{{ vault "secret/test" "foo" }}' | korvike
  bar
  ```
