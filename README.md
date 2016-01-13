# Quorra

![](/quorra.png?raw=true)

# Prequisites
To get a local setup you will need [Go](http://golang.org) and 
[npm](https://www.npmjs.com/).  On Linux, it is recommended to use 
[nvm](https://github.com/creationix/nvm).  On Mac, you should be able to 
`brew install npm`.

# Local Setup
After getting the prereqs, you can run the following to get a local dev
setup:

`make dev-setup`

This will take a while depending on your system to compile all of the Node 
dependencies needed to build the frontend application.

# Build
Use the `Makefile`

To build everything, simply run:

- `make`

To build the media only:

- `make media`

To build the Go binary only:

- `make build`

# Database
RethinkDB is used for the database.  You can use Docker to run a local
database.  Quorra will create an initial admin user and will output
a random password in the console.  To run a local RethinkDB instance using
Docker:

`docker run -ti --restart=always --name db -P rethinkdb`

You can then run `docker ps` to get the port for Quorra.  You will need
the allocated port for the container port `28015`.

# Controller
Once the database is started you can run the app.  Use the high port from
the RethinkDB container as the argument to the controller:

`./cmd/quorra/quorra -D serve`

You should then be able to access Quorra on port 8080.
