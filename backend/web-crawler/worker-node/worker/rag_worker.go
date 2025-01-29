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
	fmt.Println("RagWorker Execute")

	ragParams, err := coordinator_client.CastParams[RagWorkerParams](task.Params)
	if err != nil {
		return fmt.Errorf("invalid params %+v", task.Params)
	}

	storedWebsite, err := w.storeWebsite(ragParams.Url)
	if err != nil {
		return err
	}

	fmt.Printf("markdown: %s\n\n\n", ragParams.Markdown)

	chunks, err := w.ragClient.ChunksFrom(ragParams.Markdown)
	if err != nil {
		return fmt.Errorf("error extracting chunks: %v", err)
	}

	contacts, err := w.ragClient.ContactsFrom(ragParams.Markdown)
	if err != nil {
		return fmt.Errorf("error extracting contacts: %v", err)
	}

	newSlice := make([]string, len(chunks)+len(contacts))
	copy(newSlice, chunks)

	for i, contact := range contacts {
		newSlice[len(chunks)+i] = contact.Context
	}

	embeddings, err := w.ragClient.EmbeddingsForAll(newSlice)
	if err != nil {
		return fmt.Errorf("error extracting embeddings: %v", err)
	}

	if err := w.storeChunks(chunks, embeddings, storedWebsite.ID); err != nil {
		return fmt.Errorf("error storing chunks: %v", err)
	}

	if err := w.storeContacts(contacts, embeddings[len(chunks):], storedWebsite.ID); err != nil {
		return fmt.Errorf("error storing contacts: %v", err)
	}

	return nil
}

func (w *RagWorker) storeWebsite(url string) (*storage.Website, error) {
	storedWebsite, err := storage.Store(w.store, storage.Website{URL: url})
	if err != nil {
		return nil, fmt.Errorf("error storing website: %v", err)
	}
	return storedWebsite, nil
}

func (w *RagWorker) storeChunks(chunks []string, embeddings [][]float32, websiteID int) error {

	var rags []storage.Rag
	for i, chunk := range chunks {

		rags = append(rags, storage.Rag{
			PosInDoc:  i,
			Embedding: embeddings[i],
			Source:    "WEBSITE",
			WebsiteID: websiteID,
			Text:      chunk,
		})
	}

	if _, err := storage.StoreAll(w.store, rags...); err != nil {
		return fmt.Errorf("error storing chunks: %v", err)
	}
	return nil
}

func (w *RagWorker) storeContacts(contacts []ragger.Contact, embeddings [][]float32, websiteID int) error {

	var contactsToStore []storage.Contact
	for i, contact := range contacts {

		contactsToStore = append(contactsToStore, storage.Contact{
			Context:       contact.Context,
			PosInContacts: i,
			Contact:       contact.Value,
			ContactType:   contact.Type,
			WebsiteID:     websiteID,
			Source:        "WEBSITE",
			Embedding:     embeddings[i],
		})
	}

	if _, err := storage.StoreAll(w.store, contactsToStore...); err != nil {
		return fmt.Errorf("error storing contacts: %v", err)
	}
	return nil
}

func (w *RagWorker) Cleanup(ctx context.Context, task *coordinator_client.Task) error {
	return w.coordinatorClient.SetProcessed(ctx, coordinator_client.CoordinatorClientTaskTopicRag, task)
}
