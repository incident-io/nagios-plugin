# Releasing

As this tool is installed into Nagios, there is no automatic update mechanism. We do however track usage of versions via `User-Agent`, so you can work out who is using what version. 

Creating a new release will bump the version number in the user agent and creates a new release on GitHub.

## Steps

When you want to cut a new release, you can:

1. Merge any of the changes you want in the release to master.
1. Create a new commit on master that adjusts the CHANGELOG so all unreleased
   changes appear under the new version.
1. Tag that commit with whatever your release version should be, and push everything.

That will trigger the CI pipeline that will publish a new GitHub release with the updated binary.

That ends up looking like this:

```
git commit -m "Changelog for v1.2.3"
git tag v1.2.3
git push --tags
```
