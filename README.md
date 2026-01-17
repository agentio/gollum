# slink

`slink` is a command line tool for calling APIs described with [Lexicon](https://atproto.com/specs/lexicon).

`slink` is designed to connect to remote services through a local [IO](https://agent.io/posts/io) which  handles all routing and authentication, but `slink` can also be configured to make direct calls to XRPC hosts.

## Installing slink

To install `slink` on any system with Go installed:
1. clone the repo.
2. copy the [lexicons](https://github.com/bluesky-social/atproto/tree/main/lexicons) directory to the root of the repo.
3. run `make all`. This will generate XRPC command handlers and a CLI in the `gen` directory. It then builds everything to create `slink`.

## Using slink

With `slink`, it's very easy to call XRPC APIs that don't require authentication.
```
$ export SLINK_HOST=(your PDS url)
$ slink call com.atproto.sync list-repos
{
  "cursor": "...",
  "repos": [
    {
      "active": true,
      "did": "did:plc:...",
      "head": "...",
      "rev": "..."
    },
    {
      "active": true,
      "did": "did:plc:...",
      "head": "...",
      "rev": "..."
    }
  ]
}
```

The `slink` authentication header can be set using `SLINK_AUTH`.

Here's how you can use `slink` to make calls as a PDS admin:
```
$ export SLINK_HOST=(your PDS url)
$ export SLINK_AUTH="Basic $(echo -n "admin:$ADMIN_PASSWORD" | base64)"
$ slink call com.atproto.admin get-invite-codes
{
  "codes": [
    {
      "available": 1,
      "code": "...",
      "createdAt": "...",
      "createdBy": "admin",
      "forAccount": "admin"
    },
    ...
  ]
  "cursor": "..."
}
```

Here's how you can use `slink` to make calls as a PDS user:
```
$ export SLINK_HOST=(your PDS url)
$ slink call com.atproto.server create-session --identifier $HANDLE --password $PASSWORD
{
  "accessJwt": "...
  "active": true,
  "did": "did:plc:...
  "email": "...
  "emailConfirmed": true,
  "handle": "...",
  "refreshJwt": "..."
}
$ export SLINK_AUTH="Bearer $(slink call com.atproto.server create-session --identifier $HANDLE --password $PASSWORD | jq .accessJwt -r)"
$ slink call com.atproto.server get-session
{
  "active": true,
  "did": "did:plc:...",
  "email": "...",
  "emailConfirmed": true,
  "handle": "..."
}
```

## Warning!

`slink` is pre-release software. Not all aspects of [Lexicon](https://atproto.com/specs/lexicon) are currently supported, but that's the goal.

## Copyright

Copyright 2026, Agent IO (Tim Burks).

## License

`slink` is released under the [AGPL](https://www.gnu.org/licenses/agpl-3.0.html). 