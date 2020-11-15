# Span command line client

The Span command line client aims to be a convenient way to manage
your Span resources.  This is a work in progress so it might be a good
idea to check back and update your install of this command.


## Installing

You can install this utility by issuing the following command.

    go get -u github.com/lab5e/spancli/cmd/span

## Usage

The general usage of `span` is

    span [options] <subcommand> [subcommand options]
	
The subcommands are:

    collection  collection management commands
    data        list data from device or collection
    device      device management commands
    listen      listen for messages from Span

*There is a ws subcommand right now as well, but this will go away in
the future.*

For more help please refer to the `-h` option:

    span -h
	span <subcommand> -h


## Environment variables

In order to make life a bit easier, two of the required command line
options can be set from the environment.  Also, it is a good idea to
not use the `--token` option if you can avoid it since this makes your
token end up in your command line history.

- `--token` option can be set in `SPAN_API_TOKEN`
- `--collection-id` option can be set in `SPAN_COLLECTION_ID`

If you have these environment variables set you can omit their
respective options.  If you do specify the command line options they
will override what is set in the environment.

## Requirements

This was written in Go 1.15, so you should have Go 1.15 or a newer
version of Go installed.




