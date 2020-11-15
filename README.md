# Span command line client

The Span command line client aims to be a convenient way to manage
your Span resources.  This is a work in progress so it might be a good
idea to check back and update your install of this command.

The general usage of `span` is

    span [options] <subcommand> [subcommand options]
	
The subcommands are:

    collection  collection management commands
    data        list data from device or collection
    device      device management commands
    listen      listen for messages from Span

*There is a ws subcommand right now as well, but this will go away in
the future.*

## Installing

You can install this utility by issuing the following command.

    go get -u github.com/lab5e/spancli/cmd/span

## Requirements

This was written in Go 1.15, so you should have Go 1.15 or a newer
version of Go installed.




