# AdGuardHome Example

The configuration for AdGuardHome is contained within a single file.
This file contains mostly non-sensitive items that could easily be stored in, for example, a Kubernetes `ConfigMap`.
However it also contains the user credentials which could be stored seperately in, for example, a Kubernetes `Secret`.

To separate these concepts so that these non-sensitive items can be easily managed declaratively we can augment the official container image with `shoehorn`.

## Shoehorn Configuration

First we will need some details about AdGuard works and how currently is set up within the official container image.

- Config File: `/opt/adguardhome/conf/AdGuardHome.yaml`
- Original Entrypoint: `/opt/adguardhome/AdGuardHome`
- Original CMD: `["--no-check-update", "-c", "/opt/adguardhome/conf/AdGuardHome.yaml", "-w", "/opt/adguardhome/work"]`
- Supported Configuration Reload Method: `SIGHUP`

Then we need to make some decisions about where we want to place the separate pieces of the configuration we will be mounting.

- Main Config: `/configs/adguard/config.yaml` (This will be most of what is in `AdGuardHome.yaml`)
- Users Config: `/secrets/users/users.yaml`

Based on this information and our choices we end up with a `shoehorn.yaml` that looks like:

```yaml
generate:
  - name: AdGuardHome.yaml
    path: /opt/adguardhome/conf
    strategy: append
    inputs:
      - name: main
        path: /configs/adguard/config.yaml
      - name: users
        path: /secrets/users/users.yaml
process:
  path: /opt/adguardhome/AdGuardHome
  reload:
    enabled: true
    method: signal
    signal: SIGHUP
  args:
    [
      "--no-check-update",
      "-c",
      "/opt/adguardhome/conf/AdGuardHome.yaml",
      "-w",
      "/opt/adguardhome/work",
    ]
```

That will monitor for changes to `/configs/adguard/config.yaml` and `/secrets/users/users.yaml`. In response it will regenerate `/opt/adguardhome/conf/AdGuardHome.yaml`.
After the config file is regenerated it will send a `SIGHUP` to the AgGuard process to tell it to reload its config.

## Dockerfile

The next step is to create a `Dockerfile` to add our shoehorn. This Dockerfile will accomplish a couple things:

- Get the `shoehorn` binary
- Build on top of the official AdGuardHome container image
- Copy the shoehorn binary and config into the container image
- Override the `ENTRYPOINT` to run `shoehorn`
- Override the `CMD` to avoid duplicating args

This will result in a `Dockerfile` that looks like:

```Dockerfile
# Get the a version of the shoehorn image
FROM ghcr.io/opensourcererprime/shoehorn:<verson> AS shoehorn

FROM docker.io/adguard/adguardhome:<version>

# Copy the shoehorn binary
COPY --from=shoehorn /shoehorn/shoehorn /shoehorn/
# Copy the shoehorn configuration file
COPY shoehorn.yaml /shoehorn/
# Override original entrypoint, the original should already be configured in the shoehorn config file.
ENTRYPOINT ["/shoehorn/shoehorn", "/shoehorn/shoehorn.yaml"]
# Override original command, which is not needed now that it is embedded in the shoehorn config file.
CMD []
```

## Deploy

Now its time to deploy this using your chosen method. This could be something as simple as `docker-compose` where you would how have 3 mounts.
One for the work directory, one for mounting configs to `/configs`, and one for mounting secrets to `/secrets`. As long as the monitored file paths exist shoehorn will use the to generate the composite config file and start (or reload) the AdGuard process.

For something more complex like Kubernetes, that would require something like a helm chart or a kustomization which is a completely separate topic and effort.
