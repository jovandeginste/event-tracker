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

	c, err := a.AllCategories()
	if err != nil {
		return err
	}

	at, err := e.AutoTags(c)
	if err != nil {
		return err
	}

	e.AICategories = at

	return a.UpdateEvent(e)
}

func (e *Event) AutoTags(cats []string) ([]string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	noStream := false
	c := fmt.Sprintf(
		"Summary: %s\nDescription: %s\nLocation: %s\nTags: %s",
		e.Summary, e.Description(),
		e.Location(), strings.Join(e.Categories, ", "),
	)

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		Model: "mistral",
		System: `You will be provided with details for an event. The event details will be formatted as "key: value". Classify the event with any of the existing tags, if any are relevant.
Suggest addiditional new tags if possible in Dutch.
Provide your output in json format with the keys: tags and new_tags.
Skip tags already assigned.

Existing tags:
` + strings.Join(cats, "\n"),
		Prompt: c,
		Stream: &noStream,
	}

	var res struct {
		Tags    []string `json:"tags"`
		NewTags []string `json:"new_tags"`
	}

	respFunc := func(resp api.GenerateResponse) error {
		log.Print("Summary: ", e.Summary)
		log.Print("Response: ", resp.Response)

		return json.Unmarshal([]byte(resp.Response), &res)
	}

	if err := client.Generate(context.Background(), req, respFunc); err != nil {
		return nil, err
	}

	return slices.Concat(res.Tags, res.NewTags), nil
}
