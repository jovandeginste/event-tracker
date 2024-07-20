package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/ollama/ollama/api"
)

const (
	AIDiv    = "####"
	AIPrompt = `You will be provided with details for all existing events, and one new event, all divided with ` + AIDiv + `. The event details will be formatted as "key: value".
Classify the new event with any relevant tags from the other events.
Suggest two additional new tags.
Use the same language as the event's summary and description.
Provide only json output with the following keys: tags and new_tags. Don't add any comments or other text to the output.

Existing events:
`
)

func (a *App) AddAITags(e *Event) error {
	a.logger.Info("Adding tags to event: " + e.Summary)
	a.logger.Info("Current tags:" + strings.Join(e.AICategories, ", "))

	if len(e.AICategories) > 0 {
		a.logger.Info("skipping.")
		return nil
	}

	at, err := a.AutoTags(e)
	if err != nil {
		return err
	}

	e.AICategories = at

	return a.UpdateEvent(e)
}

func (e *Event) AIFormat() string {
	return fmt.Sprintf(
		AIDiv+"\nNew event:\nSummary: %s\nDescription: %s\nLocation: %s\nTags: %s\n"+AIDiv,
		e.Summary, e.Description(),
		e.Location(), strings.Join(e.Categories, ", "),
	)
}

func (a *App) AutoTags(e *Event) ([]string, error) {
	log.Print("Generating categories for: ", e.Summary)

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	evs, err := a.AllEvents()
	if err != nil {
		return nil, err
	}

	noStream := false
	c := e.AIFormat()

	p := AIPrompt

	for _, ev := range evs {
		p += ev.AIFormat()
	}

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		Model:  "mistral",
		System: p,
		Prompt: c,
		Stream: &noStream,
	}

	var res struct {
		Tags    []string `json:"tags"`
		NewTags []string `json:"new_tags"`
	}

	respFunc := func(resp api.GenerateResponse) error {
		log.Print("Response: ", resp.Response)

		return json.Unmarshal([]byte(resp.Response), &res)
	}

	if err := client.Generate(context.Background(), req, respFunc); err != nil {
		return nil, err
	}

	return slices.Concat(res.Tags, res.NewTags), nil
}
