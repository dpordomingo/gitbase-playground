# Quickstart

## bblfsh and gitbase dependencies

It is recommended to read about `bblfsh` and `gitbase` from its own documentation, but here is a small guide about how to run both easily:

Launch [bblfshd](https://github.com/bblfsh/bblfshd) and install the drivers. More info in the [bblfshd documentation](https://doc.bblf.sh/user/getting-started.html):

```bash
$ docker run -d --name bblfshd --privileged -p 9432:9432 -v /var/lib/bblfshd:/var/lib/bblfshd bblfsh/bblfshd
docker exec -it bblfshd bblfshctl driver install --all
```

Install [gitbase](https://github.com/src-d/gitbase):

```bash
$ go get github.com/src-d/gitbase/...
$ cd $GOPATH/src/github.com/src-d/gitbase
$ make dependencies
```

Populate a directory with git repositories and start gitbase reading from that directory:

```bash
$ REPOS_PATH=~/gitbase/repos
$ mkdir -p ${REPOS_PATH}
$ git clone git@github.com:src-d/go-git-fixtures.git ${REPOS_PATH}/go-git-fixtures
$ go run cli/gitbase/main.go server -v --git=${REPOS_PATH}
```

## gitbase-playground

Once bblfsh and gitbase are running and accessible, you can build and serve the playground:

```bash
$ make dependencies
$ make serve
```

Once the server is running &ndash;with its default values&ndash;, it will be accessible through: http://localhost:8080

You have more information on the playground architecture and development guides in the [CONTRIBUTING.md#development](CONTRIBUTING.md#development).


## run a query

You will find more info about how to run queries using the playground API on the [rest-api guide](rest-api.md)
