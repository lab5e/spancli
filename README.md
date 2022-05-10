# Span command line client

The Span command line client aims to be a convenient way to manage
your [Span](https://console.lab5e.com/) resources.  This is a work in
progress so it might be a good idea to check back and update your
install of this command.

## Installing

You can install this utility by issuing the following command.

    go get github.com/lab5e/spancli/cmd/span@latest

## Usage

The general usage of `span` is

    span [options] <subcommand> [subcommand options]

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

## TAB Completion

In order to get tab completion you can add this to your
`.bash_profile` or wherever you put your completion settings.

    _completion_span() {
        # All arguments except the first one
        args=("${COMP_WORDS[@]:1:$COMP_CWORD}")

        # Only split on newlines
        local IFS=$'\n'

        # Call completion (note that the first element of COMP_WORDS is
        # the executable itself)
        COMPREPLY=($(GO_FLAGS_COMPLETION=1 ${COMP_WORDS[0]} "${args[@]}"))
        return 0
    }

    complete -F _completion_span span


## Requirements

This was written in Go 1.18, so you should have Go 1.18 or a newer
version of Go installed if you want to edit the code.





