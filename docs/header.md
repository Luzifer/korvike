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
