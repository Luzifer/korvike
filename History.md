# 1.0.4 / 2025-01-28

  * Update go dependencies

# 1.0.3 / 2024-12-12

  * Update Go dependencies

# 1.0.2 / 2024-04-19

  * Update dependencies

# 1.0.1 / 2024-03-08

  * Update dependencies

# 1.0.0 / 2024-02-29

  * Breaking: Add sprig functions, replace some internal ones
  * Replace old build-system

**Breaking changes:**

- Function `env` no longer takes a default, use `env "MYVAR" | default "..."`
- Function `file` no longer takes a default, use `file "[filename]" | default "..."`
- Function `now` returns `time.Time`, use `now | date "[format]"`
- Function `split` now has reversed parameters `split <sep> <str>`
- Function `vault` no longet takes a default, use `vault "key" "field" | default "..."`
- Removed function `b64decode`, use `b64dec`
- Removed function `b64encode`, use `b64enc`
- Removed function `hash`, use `sha1sum` / `sha256sum` / `sha512sum`

# 0.13.0 / 2022-03-30

  * Add basic string manipulation `join` and `split`
  * Upgrade dependencies, update mod-files

# 0.12.0 / 2021-12-11

  * Add urlescape function

# 0.11.0 / 2021-08-30

  * Add support for b64decode

# 0.10.0 / 2021-04-11

  * Add "hash" template function

# 0.9.0 / 2021-03-22

  * Add support for executing sub-templates

# 0.8.2 / 2021-03-09

  * Update dependencies

# 0.8.1 / 2020-11-06

  * Fix: Update go.sum file in main dir

# 0.8.0 / 2020-11-06

  * Update blackfriday

# 0.7.1 / 2020-04-08

  * Fix tests

# 0.7.0 / 2020-04-08

  * Make blackfriday markdown available

# 0.6.1 / 2019-07-28

  * Fix: Update repo-runner config for go modules

# 0.6.0 / 2019-07-28

  * Add AppRole support into `vault` function
  * Add go modules support (remove old vendoring)
  * README updates
  * Replace build image

# 0.5.0 / 2018-05-31

  * Add b64encode as a function
  * Fix Copyright line in LICENSE

# 0.4.1 / 2017-06-20

  * Fix: New dependencies were missing in vendoring

# 0.4.0 / 2017-06-20

  * Add `vault` template function

# 0.3.0 / 2017-04-17

  * Add 'now' function and function tests

# 0.2.1 / 2017-02-14

  * Add builder config
  * Add dependencies
  * Add export for template functions

# 0.2.0 / 2016-08-24

  * Add ability to include files

# 0.1.0 / 2016-07-31

  * initial version
