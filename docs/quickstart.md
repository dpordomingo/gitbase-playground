# Quickstart

## bblfsh and gitbase dependencies

It is recommended to read about `bblfsh` and `gitbase` from its own documentation, but here is a small guide about how to run both easily:

Launch [bblfshd](https://github.com/bblfsh/bblfshd) and install the drivers. More info in the [bblfshd documentation](https://doc.bblf.sh/user/getting-started.html):

```bash
$ BBLFSHD_CONTAINER=bblfsh
$ BBLFSHD_IMAGE=bblfsh/bblfshd;
$ docker run --privileged
    --publish 9432:9432
    --volume /var/lib/bblfshd:/var/lib/bblfshd
    --name ${BBLFSHD_CONTAINER}
    ${BBLFSHD_IMAGE};
$ docker exec -it ${BBLFSHD_CONTAINER}
    bblfshctl driver install --recommended;
```

[gitbase](https://github.com/src-d/gitbase) will serve git repositories, so it is needed to populate a directory with them:

```bash
$ REPOS_PATH=~/gitbase/repos
$ mkdir -p ${REPOS_PATH}
$ git clone git@github.com:src-d/go-git-fixtures.git ${REPOS_PATH}/go-git-fixtures
```

Install and run [gitbase](https://github.com/src-d/gitbase):

```bash
$ REPOS_PATH=~/gitbase/repos
$ GITBASE_CONTAINER=gitbase;
$ BBLFSHD_CONTAINER=bblfsh;
$ GITBASE_IMAGE=dpordomingo/gitbase:latest;
$ docker run
    --publish 3306:3306
    --link ${BBLFSHD_CONTAINER}
    --volume ${REPOS_PATH}:/opt/repos
    --env BBLFSH_ENDPOINT=${BBLFSHD_CONTAINER}:9432
    --name ${GITBASE_CONTAINER}
    ${GITBASE_IMAGE};
```


## gitbase-playground

Once bblfsh and gitbase are running and accessible, you can serve the playground:

```bash
$ GITBASE_CONTAINER=gitbase;
$ PLAYGROUND_CONTAINER=gitbasePlayground;
$ PLAYGROUND_IMAGE=dpordomingo/gitbase-playground:latest;
$ docker run -d
    --publish 8080:8080
    --link ${GITBASE_CONTAINER}
    --env GITBASEPG_ENV=dev
    --env GITBASEPG_DB_CONNECTION="gitbase@tcp(${GITBASE_CONTAINER}:3306)/none?maxAllowedPacket=4194304"
    --name ${PLAYGROUND_CONTAINER}
    ${PLAYGROUND_IMAGE};
```

Once the server is running &ndash;with its default values&ndash;, it will be accessible through: http://localhost:8080

You have more information on the playground architecture and development guides in the [CONTRIBUTING.md#development](CONTRIBUTING.md#development).


## run a query

You will find more info about how to run queries using the playground API on the [rest-api guide](rest-api.md)
