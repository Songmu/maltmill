maltmill
=======

[![test](https://github.com/Songmu/maltmill/actions/workflows/test.yaml/badge.svg)][GitHub Actions]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Songmu/maltmill)][PkgGoDev]

[GitHub Actions]: https://github.com/Songmu/maltmill/actions/workflows/test.yaml
[license]: https://github.com/Songmu/maltmill/blob/master/LICENSE
[PkgGoDev]: https://pkg.go.dev/github.com/Songmu/maltmill

Create and update Homebrew third party Formulae

## Description

maltmill retrieves artifacts from GitHub Releases and creates or updates
Homebrew Formulae. It automatically detects platform-specific archives
(darwin/linux, amd64/arm64), calculates SHA256 digests, and generates
the Formula file.

This is useful for maintaining your own [tap](https://docs.brew.sh/Taps)
repository (e.g. `homebrew-tap`) without manually writing or updating
Formula files.

The name "maltmill" comes from the malt mill — the equipment used at the
beginning of the beer brewing process to grind malt.

## Synopsis

### Create a new Formula

```console
% maltmill new -w Songmu/maltmill
```

This creates `maltmill.rb` from the latest GitHub Release.
You can also specify a tag explicitly:

```console
% maltmill new -w Songmu/maltmill@v1.0.0
```

Use `-o` to specify the output file name:

```console
% maltmill new -o Formula/maltmill.rb Songmu/maltmill
```

### Update existing Formulae

```console
% maltmill -w maltmill.rb
```

This fetches the latest release, updates the version, rewrites archive
URLs, and recalculates SHA256 digests. Multiple formula files can be
specified at once:

```console
% maltmill -w Formula/*.rb
```

If the repository uses prefixed tags (e.g. `my-product-v1.2.3`), use
`--tag-prefix`:

```console
% maltmill -w --tag-prefix my-product-v my-product.rb
```

If `--tag-prefix` is omitted, maltmill selects the latest release from
plain semver tags (e.g. `v1.2.3` or `1.2.3`).

If a release contains non-archive assets such as `.rpm` or `.deb`
packages alongside `.tar.gz` / `.zip` archives, maltmill may pick the
wrong file. Use `-asset` to narrow down which assets to consider:

```console
% maltmill -w -asset '\.(tar\.gz|zip)$' my-tool.rb
```

The value is a Go regular expression matched against each asset filename.

Without `-w`, the updated formula is printed to stdout.

## Install

### homebrew

```console
% brew install Songmu/tap/maltmill
```

### using [ghg](https://github.com/Songmu/ghg)

```console
% ghg get Songmu/maltmill
```

### using [aqua](https://aquaproj.github.io/)

```console
% aqua g -i Songmu/maltmill
```

### go install (for using HEAD)

```console
% go install github.com/Songmu/maltmill/cmd/maltmill@main
```

Built binaries are available on GitHub Releases.
https://github.com/Songmu/maltmill/releases

## GitHub Token

maltmill accesses the GitHub API. To avoid rate limits or access private
repositories, set a GitHub token via the `GITHUB_TOKEN` environment
variable or the `-token` flag. If neither is provided, maltmill falls
back to the token configured in git config (via `github.token` or
credential helpers).

## Options

| Flag | Description |
|------|-------------|
| `-w` | Write the result back to the source file instead of stdout |
| `-token` | GitHub API token (default: `$GITHUB_TOKEN`) |
| `-tag-prefix` | Tag prefix used to select releases (e.g. `my-product-v`) |
| `-asset` | Regexp pattern to select release assets by filename (e.g. `\.tar\.gz$`) |

### `new` subcommand options

| Flag | Description |
|------|-------------|
| `-w` | Write the result to `<name>.rb` instead of stdout |
| `-o` | Specify the output file path |
| `-token` | GitHub API token (default: `$GITHUB_TOKEN`) |
| `-tag-prefix` | Tag prefix used to select releases |
| `-asset` | Regexp pattern to select release assets by filename |

## Author

[Songmu](https://github.com/Songmu)
