# Contributing

## Making changes

1. Make your changes and ensure tests pass:
   ```sh
   make test
   ```

2. Build locally to verify:
   ```sh
   make build
   ./gitgogit sync --config config.yaml
   ```

## Releasing

Releases use `make release`, which builds all platform binaries and publishes a GitHub release via the `gh` CLI. The target enforces several preconditions — follow these steps in order:

1. Commit your changes and push to `main`:
   ```sh
   git add <files>
   git commit -m "your message"
   git push
   ```

2. Create and push a version tag (follows `vMAJOR.MINOR.PATCH`):
   ```sh
   git tag v1.0.2
   git push origin v1.0.2
   ```

3. Run the release:
   ```sh
   make release
   ```

The release target will fail if:
- There are uncommitted changes
- HEAD is not exactly on the tag
- The tag hasn't been pushed to the remote
- You are not on the `main` branch
