# Luzifer / korvike

`korvike` is the finnish translation to the word "replacer" and that is what it does: It takes a Go template and executes it.

## Available functions

- Take key-value pairs from the CLI and replace them inside the template:  
`echo "{{ .foo }}" | korvike -v foo=bar => "bar"`
- Read environment variables and replace them inside the template:  
`export FOO=bar; echo '{{env "FOO"}}' | korvike => "bar"`
