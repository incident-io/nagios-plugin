## 1.0.3

- Fixes an issue where events with only a service or only a host in a resolved state would result in an incorrect firing state.

## 1.0.2

- Executable is no longer linked to glibc, so will run in any environment.

## 1.0.1

- Adds the ability to set additional metadata, such as a team when invoking the command. E.g. `--metadata='{"team":"myteam"}'`
- Allows the plugin to capture contact and contact group names from alerts, which may be useful in capturing the owning team of an alert.

## 1.0.0

- Initial release
