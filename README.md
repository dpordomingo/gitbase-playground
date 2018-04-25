[![Build Status](https://travis-ci.org/src-d/gitbase-playground.svg)](https://travis-ci.org/src-d/gitbase-playground)
[![codecov.io](https://codecov.io/github/src-d/gitbase-playground/coverage.svg)](https://codecov.io/github/src-d/go-git)
![unstable](https://svg-badge.appspot.com/badge/stability/unstable?a)

# Gitbase Playground

Web application powered by [gitbase](https://github.com/src-d/gitbase), were you can query git repositories with a MySQL.

<!-- ![Screenshot](.github/screenshot.png?raw=true) //-->


# Installation

```bash
go get github.com/src-d/gitbase-playground/...
```


# Usage

## Dependencies

It is needed to have a running [gitbase](https://github.com/src-d/gitbase) and [bblfsh](https://github.com/bblfsh/bblfshd) servers

## Run the playground

```bash
$ gitbase-playground -gitbase-addr <gitbase-address>
```

open http://localhost:8080


# Contribute

[Contributions](https://github.com/src-d/gitbase-playground/issues) are more than welcome, if you are interested please take a look to our [Contributing Guidelines](CONTRIBUTING.md). You have more information on how to run it locally for [development purposes here](CONTRIBUTING.md#Development).


# Code of Conduct

All activities under source{d} projects are governed by the [source{d} code of conduct](https://github.com/src-d/guide/blob/master/.github/CODE_OF_CONDUCT.md).


## License

GPL v3.0, see [LICENSE](LICENSE)
