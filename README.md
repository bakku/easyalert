# Easyalert
[![Build Status](https://travis-ci.org/bakku/easyalert.svg?branch=master)](https://travis-ci.org/bakku/easyalert)

Easyalert is a small application which enables you to send alerts in a simple and straightforward way. It is suitable especially for scripts where you want to notify yourself in case a certain event or failure occurs.

## Development

The easiest way to work on easyalert is by using docker.

First run `make docker_build` to build and pull all the images you need and then run `make init` to initialize you local docker environment. Afterwards you can run `docker-compose up` to start the application.

In case you make library changes or somehow need to rebuild your environment you can just execute `make reset` and again `make docker_build` and `make init` and you are back up.

It also makes sense to run `make docker_clean` occasionally  which will delete all containers that were exited.

### Compiling

After code changes you can stop the container and execute `make build` to build the executable in docker. Afterwards you can start the application again using `docker-compose up`.

### Running tests

You can run tests by executing `make test`. You should have the database set up as instructed previously.

### Applying migrations

Easyalert is using [gom](https://github.com/bakku/gom) for migrations. To apply the latest migrations locally you can run `make migrate`.
