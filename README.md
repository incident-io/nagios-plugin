# Nagios incident.io Alert Source Plugin

This repo contains a Nagios plugin to send Nagios notifications to incident.io. 

## Installation

1. **Download the Binary**:
   Download the relevant `notify_incident_io-*` binary for your OS from the releases page.

2. **Move the Binary to the Nagios Plugins Directory**:
   Place the `notify_incident_io` binary in your Nagios plugins directory, typically located at `/usr/local/nagios/libexec/`. You may need to use `sudo` to copy the binary to this directory.

   ```bash
   sudo cp notify_incident_io /usr/local/nagios/libexec/
   sudo chmod +x /usr/local/nagios/libexec/notify_incident_io

3. Add a new command for notifying incident-io using the binary
    ```
    define command {
    command_name    notify_incident_io
    command_line    /usr/local/nagios/libexec/notify_incident_io \
                    --api_url="<your_api_url>" \
                    --token="<your_api_token>" \
                    --metadata="{}" \
                    --host_name="$HOSTNAME$" \
                    --host_address="$HOSTADDRESS$" \
                    --host_alias="$HOSTALIAS$" \
                    --service_desc="$SERVICEDESC$" \
                    --notification_type="$NOTIFICATIONTYPE$" \
                    --host_state="$HOSTSTATE$" \
                    --service_state="$SERVICESTATE$" \
                    --service_attempt="$SERVICEATTEMPT$" \
                    --max_service_attempts="$MAXSERVICEATTEMPTS$" \
                    --last_service_state="$LASTSERVICESTATE$" \
                    --service_output="$SERVICEOUTPUT$" \
                    --host_attempt="$HOSTATTEMPT$" \
                    --max_host_attempts="$MAXHOSTATTEMPTS$" \
                    --last_host_state="$LASTHOSTSTATE$" \
                    --host_output="$HOSTOUTPUT$" \
                    --service_duration="$SERVICEDURATION$" \
                    --host_duration="$HOSTDURATION$" \
                    --last_service_check="$LASTSERVICECHECK$" \
                    --contact_name="$CONTACTNAME$" \
                    --content_alias="$CONTENTALIAS$" \
                    --last_host_check="$LASTHOSTCHECK$" \
                    --service_notification_number="$SERVICENOTIFICATIONNUMBER$" \
                    --host_notification_number="$HOSTNOTIFICATIONNUMBER$"
   }
   ```

4. By default, a title, alert status and description will be generated based on the variables provided. But you can override those with the flags `--title`, etc.
5. Add a new contact for incident.io
    ```
    define contact {
        contact_name                    incident_io
        alias                           incident.io
        service_notification_period     24x7
        host_notification_period        24x7
        service_notification_options    w,u,c,r
        host_notification_options       d,r
        service_notification_commands   notify_incident_io
        host_notification_commands      notify_incident_io
    }
    ```
6. Add the contact to a contact group
    ```
    define contactgroup {
        contactgroup_name       admins
        alias                   Nagios Administrators
        members                 root,incident_io
    }
    ```
7. Add the contact group to a service or host
    ```
    define service {
        use                     generic-service
        host_name               localhost
        service_description     PING
        check_command           check_ping!100.0,20%!500.0,60%
        contact_groups          admins
    }
    ```
8. Restart Nagios to apply the changes
    ```
    sudo systemctl restart nagios
    ```

## Building the Plugin

### Requirements

- [Go](https://golang.org/) installed (v1.16+ recommended)
- Nagios Core installation with access to configure command definitions

To build the plugin, clone the repository and run `make`:

```bash
git clone <repo-url>
cd <repo-folder>
make
