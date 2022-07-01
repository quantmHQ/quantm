package webhooks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/temporal/common"
)

// ConsumeGithubInstallationEvent handles GitHub installation events
func GithubWebhook(response http.ResponseWriter, request *http.Request) {
	id := request.Header.Get("X-GitHub-Delivery")
	signature := request.Header.Get("X-Hub-Signature")

	if signature == "" {
		handleError(id, ErrorMissingHeaderGithubSignature, http.StatusUnauthorized, response)
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	if err := verifySignature(body, signature); err != nil {
		handleError(id, err, http.StatusUnauthorized, response)
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		handleError(id, ErrorMissingHeaderGithubEvent, http.StatusBadRequest, response)
		return
	}

	event := GithubEvent(headerEvent)

	switch event {
	case GithubInstallationEvent:
		var payload common.GithubInstallationEventPayload
		err := json.Unmarshal(body, &payload)

		if err != nil {
			handleError(id, err, http.StatusBadRequest, response)
			return
		}

		consumeGithubInstallationEvent(payload, response)

	case GithubAppAuthorizationEvent:
		var payload common.GithubAppAuthorizationEventPayload
		err := json.Unmarshal(body, &payload)

		if err != nil {
			handleError(id, err, http.StatusBadRequest, response)
			return
		}

		consumeGithubAppAuthorizationEvent(payload, response)

	default:
		conf.Logger.Error("Unsupported event: " + headerEvent)
		handleError(id, ErrorInvalidEvent, http.StatusBadRequest, response)
	}
}
