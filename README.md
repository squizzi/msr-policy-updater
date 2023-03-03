# MSR Mirroring Policy Password Update Tool
Use this tool to update all poll and push mirroring policies affiliated with
all repositories across an MSR with a single username/password.

## Who is this for?
Mirroring in MSR is configured on a per repository basis instead of registry
wide, some users of MSR have the desire to configure a single service account
which manages mirroring but require that service account's password to rotate
periodically.  This tool can be used to update the username and password of
an account used across all mirroring jobs.

## Usage
Provide an MSR URL and associated credentials via `msr-username` and
`msr-password` to authenticate with the target MSR, then specify a new
username/password combo to update policies with via `username` and `password`.

Use the `poll-mirroring` and `push-mirroring` flags to specify whether to update
either push or poll mirroring policies or both.  One option is required.

```bash
docker run --rm -it squizzi/msr-policy-updater \
    --username <username> \
    --password <password> \
    --msr-url <msr-url> \
    --msr-username <msr-username> \
    --msr-password <msr-password>
    --poll-mirroring \
    --push-mirroring
```

Optionally the `--log-level debug` flag can be set for debug logging and the
number of repositories to query for in a batch can be modified from it's
default value of 100 with `-b, --batch-size`.
