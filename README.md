# Easyalert

Easyalert is a small application which enables you to send alerts in a simple and straightforward way. It is suitable especially for scripts where you want to notify yourself in case a certain event or failure occurs.

## Development

The easiest way to work on easyalert is by using docker.

1. Run `docker-compose build` to download all necessary images and to build the app image
2. Run `docker-compose up -d db` to start a database container
3. Run `docker-compose run app gom init` to initialize the migration table
4. Run `docker-compose run app gom migrate` to migrate the database

Now you can run:

`docker-compose up` to start the application. After changing something you can stop the container, run `docker-compose run app make` and start the container again with `docker-compose up`.

### Running tests

You can run tests by executing `docker-compose run app go test ./...`. You should have the database set up as instructed previously.