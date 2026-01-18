# slink

`slink` is a command line tool for calling APIs described with [Lexicon](https://atproto.com/specs/lexicon).

`slink` was created to connect to remote services through a local [IO](https://agent.io/posts/io) which  handles all routing and authentication, but `slink` can also be configured to make direct calls to XRPC hosts.

## Installing slink

To install `slink` on any system with Go installed:
1. clone the repo.
2. copy the [lexicons](https://github.com/bluesky-social/atproto/tree/main/lexicons) directory to the root of the repo.
3. run `make all`. This will generate XRPC command handlers and a CLI in the `gen` directory. It then builds everything to create `slink`.

## Using slink

`slink` makes it very easy to call XRPC APIs.

Set the host to be called with the `SLINK_HOST` environment variable (be sure to include the `https://` or `http://` prefix).
```
export SLINK_HOST=$YOUR_PDS_URL
```

With that, you can use `slink` to make unauthenticated calls.
```
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

To make authenticated calls, use `SLINK_AUTH` to specify the authentication header. To make calls as a PDS admin, set `SLINK_AUTH` to the Basic Auth credentials for the admin.
```
export SLINK_AUTH="Basic $(echo -n "admin:$ADMIN_PASSWORD" | base64)"
```

Now you can make calls as the admin.
```
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

To make calls as a PDS user, set `SLINK_AUTH` to a Bearer token for the user. Here's one way to do that:
```
export SLINK_AUTH="Bearer $(slink call com.atproto.server create-session --identifier $HANDLE --password $PASSWORD | jq .accessJwt -r)"
```

With this, you can make calls as the user.
```
$ slink call com.atproto.server get-session
{
  "active": true,
  "did": "did:plc:...",
  "email": "...",
  "emailConfirmed": true,
  "handle": "..."
}
```

## Using slink with a sidecar proxy

`slink` can be configured to use a sidecar proxy like [IO](https://agent.io/posts/io). This moves authentication and routing to the sidecar and can allow credential owners to keep their secrets secure while allowing these secrets to be used to call specific (allow-listed) XRPC methods.

If the sidecar proxy listens on a TCP socket, configuration is trivial: just set `SLINK_HOST` to the port where the proxy is listening. If the sidecar proxy listens on a Unix domain socket, it's also really easy: set `SLINK_HOST` to `unix:` followed by the socket name, e.g.
```
export SLINK_HOST=unix:@io-calling-pds
```
Here `@io-calling-pds` is the name of a Linux abstract socket that a sidecar provides for calling a PDS.

## Warning!

`slink` is pre-release software. Not all aspects of [Lexicon](https://atproto.com/specs/lexicon) are currently supported, but it's coming.

## Copyright

Copyright 2026, Agent IO (Tim Burks).

## License

`slink` is released under the [AGPL](https://www.gnu.org/licenses/agpl-3.0.html). 