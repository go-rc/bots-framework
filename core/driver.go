package bots

import (
	"net/http"
)

// The driver is doing initial request & final response processing
// That includes logging, creating input messages in a general format, sending response
type WebhookDriver interface {
	HandleWebhook(w http.ResponseWriter, r *http.Request, webhookHandler WebhookHandler)
}

type BotDriver struct {
	botHost BotHost
	router WebhooksRouter
}
var _ WebhookDriver = (*BotDriver)(nil) // Ensure BotDriver is implementing interface WebhookDriver

func NewBotDriver(host BotHost, router WebhooksRouter) WebhookDriver {
	return BotDriver{botHost: host, router: router}
}

func (d BotDriver) HandleWebhook(w http.ResponseWriter, r *http.Request, webhookHandler WebhookHandler) {
	log := d.botHost.GetLogger(r)
	log.Infof("HandleWebhook() => webhookHandler: %T", webhookHandler)

	entriesWithInputs, err := webhookHandler.GetEntryInputs(r)

	if err != nil {
		if _, ok := err.(AuthFailedError); ok {
			log.Warningf("Auth failed: %v", err)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		} else {
			log.Errorf("Failed to call webhookHandler.GetEntryInputs(r): %v", err)
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	log.Infof("Got %v entries", len(entriesWithInputs))
	for i, entryWithInputs := range entriesWithInputs {
		log.Infof("Entry[%v]: %v, %v inputs", i, entryWithInputs.Entry.GetID(), len(entryWithInputs.Inputs))
		for j, input := range entryWithInputs.Inputs {
			log.Infof("Input[%v]: %v", j, input)
		}
	}
}