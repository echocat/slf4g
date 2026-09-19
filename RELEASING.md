# Releasing

Root and native releases use the same canonical `v0.x.y` or `v1.x.y` version.
The Windows Event Log module is unfinished and is not released.

## Before publishing

1. Run `mise run check:full` on the intended release commit.
2. Confirm that Continuous Integration and CodeQL succeeded for that commit on
   `main`.
3. Prepare release notes that call out compatibility changes, including changes
   to the minimum supported Go version.

## Publish

Create and publish a GitHub release for the root tag, for example `v1.9.0`.
Prereleases and module major versions above v1 are intentionally rejected.

Publishing starts the Release workflow. It performs the following steps before
publishing the native module tag:

1. Resolves and checks out the root release tag.
2. Removes the local root-module replacement from `native/go.mod` and requires
   the exact root release version.
3. Runs the native module race tests with Go 1.23 and the released root module.
4. Creates a detached child commit and pushes `native/vX.Y.Z` to it without
   creating a branch or changing `main`.

Re-running the workflow is safe when the native tag has the expected root parent
and file tree. A tag with different contents is treated as an error and is never
moved automatically.

If verification fails, do not move or delete a published tag. Correct the
problem and publish a new patch version.
