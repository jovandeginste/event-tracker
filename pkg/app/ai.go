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

func (a *App) AddAITags(e *Event) error {
	if len(e.AICategories) > 0 {
		return nil
	}

	at, err := a.AutoTags(e)
	if err != nil {
		return err
	}

	e.AICategories = at

	return a.UpdateEvent(e)
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

	div := "####"
	noStream := false
	c := fmt.Sprintf(
		div+"\nNew event:\nSummary: %s\nDescription: %s\nLocation: %s\nTags: %s\n"+div,
		e.Summary, e.Description(),
		e.Location(), strings.Join(e.Categories, ", "),
	)

	p := `You will be provided with details for all existing events, and one new event, all divided with ` + div + `. The event details will be formatted as "key: value". Classify the new event with any tags from the previous events, if any are relevant.
Suggest additional new tags if no existing tags are relevant.
Use the same language as the event summary and description.
Provide only json output with the following keys: tags and new_tags. Don't add any comments or other text to the output.

Existing events:
`

	for _, ev := range evs {
		p += fmt.Sprintf(
			div+"\nSummary: %s\nDescription: %s\nLocation: %s\nTags: %s\n",
			ev.Summary, ev.Description(),
			ev.Location(), strings.Join(ev.Categories, ", "),
		)
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
