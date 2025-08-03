# Get the latest version of the shoehorn image
FROM ghcr.io/opensourcererprime/shoehorn:latest AS shoehorn

FROM docker.io/adguard/adguardhome:latest

# Copy the shoehorn binary
COPY --from=shoehorn /shoehorn/shoehorn /shoehorn/
# Copy the shoehorn configuration file
COPY shoehorn.yaml /shoehorn/
# Override original entrypoint, the original should already be configured in the shoehorn config file.
ENTRYPOINT ["/shoehorn/shoehorn", "/shoehorn/shoehorn.yaml"]
# Override original command, which is not needed now that it is embedded in the shoehorn config file.
CMD []
