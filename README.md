# MSR Mirroring Policy Password Update Tool
Use this tool to update all poll and push mirroring policies affiliated with
all repositories across an MSR with a single username/password.

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
