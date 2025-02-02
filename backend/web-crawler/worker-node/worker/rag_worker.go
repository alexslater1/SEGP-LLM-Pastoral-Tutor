package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/storage"
	"github.com/ethanhosier/worker-node/utils"
	"github.com/google/uuid"
)

type RagWorker struct {
	id                string
	ragClient         ragger.Ragger
	coordinatorClient coordinator_client.CoordinatorClient
	store             storage.Storage
}

type RagWorkerParams struct {
	Markdown  string `json:"markdown"`
	Url       string `json:"url"`
	InnerText string `json:"text"`
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

	storedRagSource, err := w.storeRagSource(ragParams.Url, "WEBSITE")
	if err != nil {
		return err
	}

	chunks, err := w.ragClient.ChunksFrom(utils.CleanText(ragParams.InnerText))
	if err != nil {
		return fmt.Errorf("error extracting chunks: %v", err)
	}

	contacts, err := w.ragClient.ContactsFrom(utils.CleanText(ragParams.Markdown))
	if err != nil {
		return fmt.Errorf("error extracting contacts: %v", err)
	}

	newSlice := make([]string, len(chunks)+len(contacts))
	copy(newSlice, chunks)

	for i, contact := range contacts {
		newSlice[len(chunks)+i] = contact.Context
	}

	log.Printf("RAG: generating embeddings for %d chunks and %d contacts\n", len(chunks), len(contacts))
	embeddings, err := w.ragClient.EmbeddingsForAll(newSlice)
	if err != nil {
		return fmt.Errorf("error extracting embeddings: %v", err)
	}

	if err := w.storeChunks(chunks, embeddings, storedRagSource.ID); err != nil {
		return fmt.Errorf("error storing chunks: %v", err)
	}

	if err := w.storeContacts(contacts, embeddings[len(chunks):], storedRagSource.ID); err != nil {
		return fmt.Errorf("error storing contacts: %v", err)
	}

	return nil
}

func (w *RagWorker) storeRagSource(url string, typ string) (*storage.RagSource, error) {
	storedRagSource, err := storage.Store(w.store, storage.RagSource{URL: url, Type: typ})
	if err != nil {
		return nil, fmt.Errorf("error storing rag source: %v", err)
	}
	return storedRagSource, nil
}

func (w *RagWorker) storeChunks(chunks []string, embeddings [][]float32, ragSourceId int) error {

	var rags []storage.RagChunk
	for i, chunk := range chunks {

		rags = append(rags, storage.RagChunk{
			PosInSource: i,
			Embedding:   embeddings[i],
			RagSourceId: ragSourceId,
			Text:        chunk,
		})
	}

	if _, err := storage.StoreAll(w.store, rags...); err != nil {
		return fmt.Errorf("error storing chunks: %v", err)
	}
	return nil
}

func (w *RagWorker) storeContacts(contacts []ragger.Contact, embeddings [][]float32, ragSourceId int) error {

	var contactsToStore []storage.RagContact
	for i, contact := range contacts {

		contactsToStore = append(contactsToStore, storage.RagContact{
			Context:     contact.Context,
			PosInSource: i,
			Contact:     contact.Value,
			ContactType: string(contact.Type),
			RagSourceId: ragSourceId,
			Embedding:   embeddings[i],
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
