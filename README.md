# slink

`slink` is a command line tool for calling APIs described with [Lexicon](https://atproto.com/specs/lexicon).

`slink` connects to remote services through a local [IO](https://agent.io/posts/io) which handles all routing and authentication.

To install `slink` on any system with Go installed:
1. clone the repo.
2. copy lexicons into a top-level directory called `lexicons`.
3. run `make all`.
