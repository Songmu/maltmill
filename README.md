maltmill
=======

[![test](https://github.com/Songmu/maltmill/actions/workflows/test.yaml/badge.svg)][GitHub Actions]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Songmu/maltmill)][PkgGoDev]

[GitHub Actions]: https://github.com/Songmu/maltmill/actions/workflows/test.yaml
[license]: https://github.com/Songmu/maltmill/blob/master/LICENSE
[PkgGoDev]: https://pkg.go.dev/github.com/Songmu/maltmill

create and update Homebrew thrid party Formulae

## Synopsis

### new

```console
% maltmill new -w Songmu/maltmill
```

### update

```console
% maltmill -w maltmill.rb
```

## Description

The maltmill retrieve artifacts from GitHub Releases and create or update
the Homebrew Formulae.

## Install

### homebrew

```console
% brew install Songmu/tap/maltmill
```

### using [ghg](https://github.com/Songmu/ghg)

```console
% ghg get Songmu/maltmill
```

### go get (for using HEAD)

```console
% go get github.com/Songmu/maltmill/cmd/maltmill
```

Built binaries are available on GitHub Releases.
https://github.com/Songmu/maltmill/releases

## Author

[Songmu](https://github.com/Songmu)
