package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var version = "development"

type AlertSourcePayload struct {
	Title            string                 `json:"title"`
	Status           string                 `json:"status"`
	DeduplicationKey string                 `json:"deduplication_key,omitempty"`
	Description      string                 `json:"description,omitempty"`
	SourceURL        string                 `json:"source_url,omitempty"`
	Metadata         map[string]interface{} `json:"metadata"`
}

func sendIncidentNotification(apiURL, token string, incidentData AlertSourcePayload) error {
	jsonData, err := json.Marshal(incidentData)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	fmt.Println("Sending incident.io notification with payload:")
	fmt.Println(string(jsonData))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Nagios Incident.io Plugin/"+version)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to send notification: %s", resp.Status)
	}

	fmt.Println("Incident.io notification sent successfully in status", incidentData.Status)
	return nil
}

func main() {
	// Define named command-line arguments
	apiURL := flag.String("api_url", "", "The Incident.io API URL")
	token := flag.String("token", "", "The API token for authorization")
	overrideDeduplicationKey := flag.String("deduplication_key", "", "(Optional) Deduplication key for the incident")
	overrideTitle := flag.String("title", "", "(Optional) Title of the incident")
	overrideDescription := flag.String("description", "", "(Optional) Description of the incident")
	overrideSourceURL := flag.String("source_url", "", "(Optional) Source URL for more information")

	// Metadata fields
	hostName := flag.String("host_name", "", "The name of the host")
	hostAddress := flag.String("host_address", "", "The address of the host")
	hostAlias := flag.String("host_alias", "", "The alias of the host")
	serviceDesc := flag.String("service_desc", "", "The description of the service")
	notificationType := flag.String("notification_type", "", "The type of notification")
	hostState := flag.String("host_state", "", "The state of the host")
	serviceState := flag.String("service_state", "", "The state of the service")
	serviceAttempt := flag.String("service_attempt", "", "The current attempt number for the service check")
	maxServiceAttempts := flag.String("max_service_attempts", "", "The maximum number of service check attempts")
	lastServiceState := flag.String("last_service_state", "", "The last state of the service or host")
	serviceOutput := flag.String("service_output", "", "The output of the service check")
	hostAttempt := flag.String("host_attempt", "", "The current attempt number for the host check")
	maxHostAttempts := flag.String("max_host_attempts", "", "The maximum number of host check attempts")
	lastHostState := flag.String("last_host_state", "", "The last state of the host")
	hostOutput := flag.String("host_output", "", "The output of the host check")
	serviceDuration := flag.String("service_duration", "", "The duration of the service check")
	hostDuration := flag.String("host_duration", "", "The duration of the host check")
	lastServiceCheck := flag.String("last_service_check", "", "The last state of the service")
	lastHostCheck := flag.String("last_host_check", "", "The last state of the host")
	serviceNotificationNumber := flag.String("service_notification_number", "", "The notification number for the service")
	hostNotificationNumber := flag.String("host_notification_number", "", "The notification number for the host")
	contactName := flag.String("contact_name", "", "The name of the contact")
	contactAlias := flag.String("contact_alias", "", "The alias of the contact")
	contactGroupAlias := flag.String("contact_group_alias", "", "The alias of the contact group")
	otherMetadata := flag.String("metadata", "", "Additional metadata in JSON format")

	// Parse the command-line arguments
	flag.Parse()

	// Check for required arguments
	if apiURL == nil || token == nil || *apiURL == "" || *token == "" {
		fmt.Println("Usage: notify_incident --api_url=<api_url> --token=<token> [--deduplication_key=<deduplication_key>] [--description=<description>] [--source_url=<source_url>] [--service_desc=<serviceDesc>] ...")
		os.Exit(2)
	}

	// Use the provided title, or construct one from the host and service
	title := ""
	if overrideTitle != nil && *overrideTitle != "" {
		title = *overrideTitle
	} else if hostName != nil && serviceDesc != nil && *hostName != "" && *serviceDesc != "" {
		title = fmt.Sprintf("%s: %s", *hostName, *serviceDesc)
	} else if hostName != nil && *hostName != "" {
		title = *hostName
	} else if serviceDesc != nil && *serviceDesc != "" {
		title = *serviceDesc
	} else {
		title = "Nagios Notification"
	}

	// Determine status based on available host and service states
	var status string
	
	// If both host and service are provided, both must be OK/UP for resolved status
	if hostState != nil && *hostState != "" && serviceState != nil && *serviceState != "" {
		if *hostState == "UP" && *serviceState == "OK" {
			status = "resolved"
		} else {
			status = "firing"
		}
	} else if hostState != nil && *hostState != "" {
		// Only host state is provided
		if *hostState == "UP" {
			status = "resolved"
		} else {
			status = "firing"
		}
	} else if serviceState != nil && *serviceState != "" {
		// Only service state is provided
		if *serviceState == "OK" {
			status = "resolved"
		} else {
			status = "firing"
		}
	} else {
		// Neither host nor service state is provided, default to firing
		status = "firing"
	}
	// Safely print states, handling nil pointers
	hostStateValue := "<nil>"
	if hostState != nil {
		hostStateValue = *hostState
	}
	serviceStateValue := "<nil>"
	if serviceState != nil {
		serviceStateValue = *serviceState
	}
	fmt.Print("States are: ", hostStateValue, " ", serviceStateValue, "\n")

	deduplicationKey := ""
	if overrideDeduplicationKey != nil && *overrideDeduplicationKey != "" {
		deduplicationKey = *overrideDeduplicationKey
	} else {
		defaultDeduplicationKey := ""
		// Generate a default deduplication key based on host and service if not provided
		if hostName != nil && serviceDesc != nil && *hostName != "" && *serviceDesc != "" {
			defaultDeduplicationKey = fmt.Sprintf("%s-%s", *hostName, *serviceDesc)
		} else if hostName != nil && *hostName != "" {
			defaultDeduplicationKey = *hostName
		} else if serviceDesc != nil && *serviceDesc != "" {
			defaultDeduplicationKey = *serviceDesc
		}
		deduplicationKey = defaultDeduplicationKey
	}

	// Build metadata from individual arguments
	metadata := make(map[string]interface{})
	if hostName != nil && *hostName != "" {
		metadata["host_name"] = *hostName
	}
	if hostAddress != nil && *hostAddress != "" {
		metadata["host_address"] = *hostAddress
	}
	if hostAlias != nil && *hostAlias != "" {
		metadata["host_alias"] = *hostAlias
	}
	if serviceDesc != nil && *serviceDesc != "" {
		metadata["service_desc"] = *serviceDesc
	}
	if notificationType != nil && *notificationType != "" {
		metadata["notification_type"] = *notificationType
	}
	if hostState != nil && *hostState != "" {
		metadata["host_state"] = *hostState
	}
	if serviceState != nil && *serviceState != "" {
		metadata["service_state"] = *serviceState
	}
	if serviceAttempt != nil && *serviceAttempt != "" {
		metadata["service_attempt"] = *serviceAttempt
	}
	if maxServiceAttempts != nil && *maxServiceAttempts != "" {
		metadata["max_service_attempts"] = *maxServiceAttempts
	}
	if lastServiceState != nil && *lastServiceState != "" {
		metadata["last_service_state"] = *lastServiceState
	}
	if serviceOutput != nil && *serviceOutput != "" {
		metadata["service_output"] = *serviceOutput
	}
	if hostAttempt != nil && *hostAttempt != "" {
		metadata["host_attempt"] = *hostAttempt
	}
	if maxHostAttempts != nil && *maxHostAttempts != "" {
		metadata["max_host_attempts"] = *maxHostAttempts
	}
	if lastHostState != nil && *lastHostState != "" {
		metadata["last_host_state"] = *lastHostState
	}
	if hostOutput != nil && *hostOutput != "" {
		metadata["host_output"] = *hostOutput
	}
	if serviceDuration != nil && *serviceDuration != "" {
		metadata["service_duration"] = *serviceDuration
	}
	if hostDuration != nil && *hostDuration != "" {
		metadata["host_duration"] = *hostDuration
	}
	if lastServiceCheck != nil && *lastServiceCheck != "" {
		metadata["last_service_check"] = *lastServiceCheck
	}
	if lastHostCheck != nil && *lastHostCheck != "" {
		metadata["last_host_check"] = *lastHostCheck
	}
	if serviceNotificationNumber != nil && *serviceNotificationNumber != "" {
		metadata["service_notification_number"] = *serviceNotificationNumber
	}
	if hostNotificationNumber != nil && *hostNotificationNumber != "" {
		metadata["host_notification_number"] = *hostNotificationNumber
	}
	if contactName != nil && *contactName != "" {
		metadata["contact_name"] = *contactName
	}
	if contactAlias != nil && *contactAlias != "" {
		metadata["contact_alias"] = *contactAlias
	}
	if contactGroupAlias != nil && *contactGroupAlias != "" {
		metadata["contact_group_alias"] = *contactGroupAlias
	}

	// Parse additional metadata from JSON string
	if otherMetadata != nil && *otherMetadata != "" {
		var otherMetadataMap map[string]interface{}
		if err := json.Unmarshal([]byte(*otherMetadata), &otherMetadataMap); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(2)
		}
		for key, value := range otherMetadataMap {
			metadata[key] = value
		}
	}

	description := ""
	if overrideDescription != nil && *overrideDescription != "" {
		description = *overrideDescription
	}

	sourceURL := ""
	if overrideSourceURL != nil && *overrideSourceURL != "" {
		sourceURL = *overrideSourceURL
	}

	// Construct the incident data
	incidentData := AlertSourcePayload{
		DeduplicationKey: deduplicationKey,
		Description:      description,
		Metadata:         metadata,
		SourceURL:        sourceURL,
		Status:           status,
		Title:            title,
	}

	// Send the notification
	if err := sendIncidentNotification(*apiURL, *token, incidentData); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
}
