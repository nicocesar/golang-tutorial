# golang-tutorial

This is a brief tutorial I created for a couple friends to show how to do some stuff in Go and Ethereum.

# Milestones

1. Create the environment to compile/execute in docker
1. Use abigen to compile an ERC20 contract into a Go 
1. Create an account in Infura and read information from the ERC20 token
1. Compile a UniswapV3 Pool and read swap() logs from it
   1. Process the logs into something more readeable print it out on screen.
1. Get the list of the most active pools from the Graph and monitor them all at once.

# Structure

* Each milestone will be a separate branch/tag (named milestone-nnn-xxxxxxx) so you can compare each step with a diff
* main branch will have the "final" code. 

# Running the code


First, you'll need a Infura API key. Get an account at https://infura.io/register.

This project is meant to be used as supporting material for a tutorial. Flawless execution is not the priority, but if you want a quick start, try:

```shell
$ ./make
go build -o golang-tutorial main.go getPools.go
$ export INFURA_API_KEY=....
$ ./golang-tutorial
Watching swaps for pool:  0x69D91B94f0AaF8e8A2586909fA77A5c2c89818d5
Watching swaps for pool:  0x2F62f2B4c5fcd7570a709DeC05D68EA19c82A9ec
(...)
```

You also can try out all this using Docker,

```shell
$ docker build . -t nicocesar/golang-tutorial
(...)
Successfully built 6b65760b3cd8
Successfully tagged nicocesar/golang-tutorial:latest
docker run -ti --env INFURA_API_KEY=... nicocesar/golang-tutorial
Watching swaps for pool:  0x69D91B94f0AaF8e8A2586909fA77A5c2c89818d5
Watching swaps for pool:  0x2F62f2B4c5fcd7570a709DeC05D68EA19c82A9ec
```

# Troubleshooting 

## INFURA_API_KEY is not set. 

If you see this:

```shell
$ ./golang-tutorial 
INFURA_API_KEY is not set. Get an account at https://infura.io/register
```

You'll need to create an INFURA_API_KEY from https://infura.io/register, the process is free of charge for low amount of requests. This tutorial can be addapted to use a local installation too, but currently is not the main objective of this repo.

## GLIBC_2.xx not found

If you see this: 

```shell
$ docker run -ti --env INFURA_API_KEY=...  nicocesar/golang-tutorial
./golang-tutorial: /lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.32' not found (required by ./golang-tutorial)
./golang-tutorial: /lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.34' not found (required by ./golang-tutorial)
```

It's very likely that there is a compilation in the local environment (Ubuntu/Debian, etc) that has a different glib from the alpine docker image. Solution:

```shell
$ make clean
rm -rf geth-* tools solidity_contracts golang-tutorial
$ docker build . -t nicocesar/golang-tutorial
(...)
$ docker run -ti --env INFURA_API_KEY=..  nicocesar/golang-tutorial
```


# Contributing

Feel free to create a Pull Request with suggestions
