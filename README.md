This program gets livestreams from different hosts (e.g twitch, hitbox) and adds them to a database.

# Packages
## Main
Contains the main function which uses the `api` package to get a number of livestreams
for a specific game. The number of streams and which game is specified at compile-time.

### Flags
The host to query and which database environment (e.g prod) is specified using input flags.

* -site (shorthand -s) specifies which host you want to query (e.g twitch).
* -env (shorthand -e) specified which database environment to target.

## Api
The package `api` specifies the `api.Api` interface which is the layer between the actual
host-specific implementations and the main function. The `api` package also contains a 
struct `api.Streams` which should correspond with your database. The job of the host-specific
implmenetations is then to implement the functions for the interface and convert the
result of an actual api call into the `api.Streams` struct.

### Util functions
The api package also provides the function `api.ApiCall(...)` which performs a GET
request to the given url using the (optional) authentication value-pair and an arbitrary
amount of parameters. The authentication and parameters are given as an `api.Pair` struct
(which is a simple string: string) pair. If no authentication is requred simply provide empty strings.

## Host-specific implementation
To implement a new host simply create a new package corresponding to the hostname (e.g package azubu) and
create a struct which will implement the `api.Api` interface. The details of this implementation is 
up to you. A straightforward way to handle it (which also allows you to use the utility functions
in the `api` package) is the create a struct that corresponds to the json-response from the host,
the the `api.ApiCall()` function to unmarshal the response into this struct and then convert the
struct to an `api.Streams` struct. The already implemented `twitch` and `hitbox` packages can be used as 
reference for this method (twitch requiring authentication while hitbox does not).

Remember to update the init function that sets up the api connection in `update-streams.go`
to correctly handle your implementation.
