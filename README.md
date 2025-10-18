# cli

`cli` is a minimalist, ergonomic, scalable library for building Go
commands.

## Usage

Simple [example](example/main.go) demonstrates subcommands, flag struct tags,
error handling, and a simple CLI code management pattern in a single small Go
file. No dependencies outside Go stdlib.

## Why

There are lots of cli libraries out there, I've worked with many and I don't
like any of them.

- Go's stdlib `flag` only works reasonably for the most rudimentary things
  and is quite difficult to get basic things working like being able to 
  append a flag to a command line that has arguments already.
- Cobra/Viper is huge and extremely verbose and seems to constantlly yell
  "you should ..." when it's just a distraction.  It is also much harder
  to customize.
- Other attempts at finding a balance between minimalism and expressivity
  don't seem to hit the mark for expressivity.


## Status

`cli` is young but it works.  The general design and organisation of the data types
has mostly settled but may evolve slightly.  May still have some bugs in the corners.


