maltmill
=======

[![Build Status](https://travis-ci.org/Songmu/maltmill.png?branch=master)][travis]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/maltmill?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/maltmill
[license]: https://github.com/Songmu/maltmill/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/maltmill

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
