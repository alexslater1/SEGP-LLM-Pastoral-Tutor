package worker

import (
	"context"
	"fmt"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/storage"
	"github.com/google/uuid"
)

type RagWorker struct {
	id                string
	ragClient         ragger.Ragger
	coordinatorClient coordinator_client.CoordinatorClient
	store             storage.Storage
}

type RagWorkerParams struct {
	Markdown string `json:"markdown"`
	Url      string `json:"url"`
}

func NewRagWorker(ragClient ragger.Ragger, coordinatorClient coordinator_client.CoordinatorClient, store storage.Storage) *RagWorker {
	id := uuid.New().String()
	return &RagWorker{id: id, ragClient: ragClient, coordinatorClient: coordinatorClient, store: store}
}

func (w *RagWorker) WorkerType() WorkerType {
	return WorkerTypeRag
}

func (w *RagWorker) Id() string {
	return w.id
}

func (w *RagWorker) Execute(ctx context.Context, task *coordinator_client.Task) error {
	ragParams, err := coordinator_client.CastParams[RagWorkerParams](task.Params)
	if err != nil {
		return fmt.Errorf("invalid params %+v", task.Params)
	}

	storedWebsite, err := w.storeWebsite(ragParams.Url)
	if err != nil {
		return err
	}

	if err := w.processAndStoreChunks(ragParams.Markdown, storedWebsite.ID); err != nil {
		return err
	}

	return w.processAndStoreContacts(ragParams.Markdown, storedWebsite.ID)
}

func (w *RagWorker) storeWebsite(url string) (*storage.Website, error) {
	storedWebsite, err := storage.Store(w.store, storage.Website{URL: url})
	if err != nil {
		return nil, fmt.Errorf("error storing website: %v", err)
	}
	return storedWebsite, nil
}

func (w *RagWorker) processAndStoreChunks(markdown string, websiteID int) error {
	chunks, err := w.ragClient.ChunksFrom(markdown)
	if err != nil {
		return fmt.Errorf("error extracting chunks: %v", err)
	}

	var rags []storage.Rag
	for i, chunk := range chunks {
		embeddings, err := w.ragClient.EmbeddingsFor(chunk)
		if err != nil {
			return fmt.Errorf("error extracting embeddings: %v", err)
		}

		rags = append(rags, storage.Rag{
			PosInDoc:  i,
			Embedding: embeddings,
			Source:    "WEBSITE",
			WebsiteID: websiteID,
			Text:      chunk,
		})
	}

	if _, err = storage.StoreAll(w.store, rags...); err != nil {
		return fmt.Errorf("error storing chunks: %v", err)
	}
	return nil
}

func (w *RagWorker) processAndStoreContacts(markdown string, websiteID int) error {
	contacts, err := w.ragClient.ContactsFrom(markdown)
	if err != nil {
		return fmt.Errorf("error extracting contacts: %v", err)
	}

	var contactsToStore []storage.Contact
	for i, contact := range contacts {
		embedding, err := w.ragClient.EmbeddingsFor(contact.Value)
		if err != nil {
			return fmt.Errorf("error extracting embeddings: %v", err)
		}

		contactsToStore = append(contactsToStore, storage.Contact{
			Context:       contact.Context,
			PosInContacts: i,
			Contact:       contact.Value,
			ContactType:   contact.Type,
			WebsiteID:     websiteID,
			Source:        "WEBSITE",
			Embedding:     embedding,
		})
	}

	if _, err = storage.StoreAll(w.store, contactsToStore...); err != nil {
		return fmt.Errorf("error storing contacts: %v", err)
	}
	return nil
}

func (w *RagWorker) Cleanup(ctx context.Context, task *coordinator_client.Task) error {
	return w.coordinatorClient.SetProcessed(ctx, coordinator_client.CoordinatorClientTaskTopicRag, task)
}
