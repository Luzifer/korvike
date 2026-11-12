![License](https://img.shields.io/github/license/Luzifer/korvike)
![Downloads](https://img.shields.io/github/downloads/Luzifer/korvike/total)
![Latest release](https://img.shields.io/github/v/release/Luzifer/korvike)

# Luzifer / korvike

`korvike` is Finnish for "substitute." It is a command-line tool for
rendering Go templates from files or standard input. Values can be supplied
through command-line flags, while Sprig and Korvike-specific functions
provide common transformations and integrations.

## Getting started

Download a binary from the
[GitHub releases](https://github.com/Luzifer/korvike/releases),
or install with Go 1.26 or newer:

```console
$ go install github.com/Luzifer/korvike@latest
```

Make sure your Go binary directory (`GOBIN`, or `$(go env GOPATH)/bin`)
is in `PATH`.

Pass variables with `-v` and read a template from standard input:

```console
$ echo '{{ .name }}' | korvike -v name=World
World
```

By default, `korvike` reads from standard input and writes to standard
output. Use `--in` (`-i`) and `--out` (`-o`) for files, `--key-value`
(`-v`) for template variables, and `--version` to print the version.

Use `--safe` for untrusted templates. It loads Sprig's hermetic
functions and Korvike's safe functions. Functions that read files
(`file`, `mustFile`), access Vault (`vault`, `mustVault`), or execute
nested templates (`tplexec`) are unavailable in this mode.

## Available functions

Starting with `v1.0.0` Korvike is based on the
[sprig functions collection](https://masterminds.github.io/sprig/)
with some additions:

> [!NOTE]
> Entries marked with 🔒 are available in `--safe` mode.
> Entries marked as **Deprecated** will be removed in a later release.

- `{{ .<variable name> }}` 🔒\
  Take key-value pairs from the CLI and replace them inside the template
  ```console
  $ korvike -v 'foo=bar' <<-EOF
  {{ .foo }}
  EOF
  bar
  ```
- `{{ extractMarkdownContent <source> }}` 🔒\
  Return the Markdown body without frontmatter. Missing frontmatter is valid and returns the original Markdown. Invalid frontmatter returns an error.
  ```console
  $ echo -e '---\ntitle: Example\n---\n# Heading' > page.md
  $ korvike <<-EOF
  {{ extractMarkdownContent (file "page.md") }}
  EOF
  # Heading
  ```
- `{{ extractMarkdownFrontmatter <source> }}` 🔒\
  Extract frontmatter as a map. Missing frontmatter is valid and returns an empty map. Invalid frontmatter returns an error.
  ```console
  $ echo -e '---\ntitle: Example\n---\n# Heading' > page.md
  $ korvike <<-EOF
  {{ index (extractMarkdownFrontmatter (file "page.md")) "title" }}
  EOF
  Example
  ```
- `{{ file <file name> }}`\
  Read a file and place it inside the template. Returns an empty string on error.
  ```console
  $ echo "Hello World" > hello
  $ korvike <<-EOF
  {{ file "hello" }}
  EOF
  Hello World
  ```
- `{{ fromYaml <source> }}` 🔒\
  Parse YAML into a template value.
  ```console
  $ korvike <<-EOF
  {{ index (fromYaml "answer: 42") "answer" }}
  EOF
  42
  ```
- `{{ markdown <source> }}` 🔒\
  Format the source using a Markdown parser.
  ```console
  $ korvike <<-EOF
  {{ markdown "# headline" }}
  EOF
  <h1>headline</h1>
  ```
- `{{ mustFile <file name> }}` **(Deprecated)**\
  Read a file and place it inside the template. Returns an error when the file cannot be read.
  ```console
  $ echo "Hello World" > hello
  $ korvike <<-EOF
  {{ mustFile "hello" }}
  EOF
  Hello World
  ```
- `{{ mustVault <path> <key> }}` **(Deprecated)**\
  Read a key from Vault at `VAULT_ADDR`. Authenticate with `VAULT_TOKEN`, `VAULT_ROLE_ID` with optional `VAULT_SECRET_ID`, or a `~/.vault-token` file. Returns an error when the value cannot be read.
  ```console
  $ vault write secret/test foo=bar
  $ korvike <<-EOF
  {{ mustVault "secret/test" "foo" }}
  EOF
  bar
  ```
- `{{ toYaml <value> }}` 🔒\
  Serialize a value as YAML with two-space indentation.
  ```console
  $ korvike <<-EOF
  {{ toYaml (dict "answer" 42) }}
  EOF
  answer: 42
  ```
- `{{ tplexec (file "my.tpl") }}`\
  Execute the given template with the same function set and variables as the parent template.
  ```console
  $ echo '{{ env "FOO" }}' >my.tpl
  $ export 'FOO=bar'
  $ korvike <<-EOF
  {{ tplexec (file "my.tpl") }}
  EOF
  bar
  ```
- `{{ urlescape <input string> }}` 🔒\
  Escape input for use inside a URL query parameter.
  ```console
  $ korvike <<-EOF
  {{ urlescape "Hellö Wörld@Golang" }}
  EOF
  Hell%C3%B6+W%C3%B6rld%40Golang
  ```
- `{{ vault <path> <key> }}`\
  Read a key from Vault at `VAULT_ADDR`. Authenticate with `VAULT_TOKEN`, `VAULT_ROLE_ID` with optional `VAULT_SECRET_ID`, or a `~/.vault-token` file. Returns an empty string on error.
  ```console
  $ vault write secret/test foo=bar
  $ korvike <<-EOF
  {{ vault "secret/test" "foo" }}
  EOF
  bar
  ```
