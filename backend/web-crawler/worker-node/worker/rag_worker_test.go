package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/storage"
	"github.com/stretchr/testify/assert"
)

func TestRagWorkerWorkerType(t *testing.T) {
	mock := NewRagWorker(nil, nil, nil)
	if mock.WorkerType() != WorkerTypeRag {
		t.Errorf("WorkerType should be WorkerTypeRag")
	}
}

func TestRagWorkerId(t *testing.T) {
	ragWorker := NewRagWorker(nil, nil, nil)
	if ragWorker.Id() == "" {
		t.Errorf("Id should not be empty")
	}
}

func TestStoreWebsite(t *testing.T) {
	var (
		memoryStorage = storage.NewMemoryStorage()
		ragWorker     = NewRagWorker(nil, nil, memoryStorage)
		url           = "https://example.com"
	)

	storedRagSource, err := ragWorker.storeRagSource(url, "WEBSITE")
	if err != nil {
		t.Errorf("Error storing rag source: %v", err)
	}

	assert.Equal(t, storedRagSource.URL, url)
	assert.Equal(t, storedRagSource.Type, "WEBSITE")
}

func TestRagWorkerStoreChunks(t *testing.T) {
	var (
		memoryStorage = storage.NewMemoryStorage()
		ragWorker     = NewRagWorker(nil, nil, memoryStorage)
		chunks        = []string{"Hello, world!1", "Hello, world!2", "Hello, world!3"}
		embeddings    = [][]float32{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}, {7.0, 8.0, 9.0}}
	)

	ragWorker.storeChunks(chunks, embeddings, 1)

	rags, err := storage.GetAll[storage.RagChunk](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting chunks: %v", err)
	}

	for _, rag := range rags {
		assert.Equal(t, rag.Text, chunks[rag.PosInSource])
		assert.Equal(t, rag.Embedding, embeddings[rag.PosInSource])
		assert.Equal(t, rag.RagSourceId, 1)
	}
}

func TestRagWorkerStoreContacts(t *testing.T) {
	var (
		memoryStorage = storage.NewMemoryStorage()
		ragWorker     = NewRagWorker(nil, nil, memoryStorage)
	)

	ragWorker.storeContacts([]ragger.Contact{
		{Context: "Hello, world!", Value: "John Doe", Type: "person"},
	}, [][]float32{{1.0, 2.0, 3.0}}, 1)

	contacts, err := storage.GetAll[storage.RagContact](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting contacts: %v", err)
	}

	assert.Equal(t, len(contacts), 1)
	assert.Equal(t, contacts[0].Context, "Hello, world!")
	assert.Equal(t, contacts[0].Contact, "John Doe")
	assert.Equal(t, contacts[0].ContactType, "person")
	assert.Equal(t, contacts[0].Embedding, []float32{1.0, 2.0, 3.0})
	assert.Equal(t, contacts[0].RagSourceId, 1)
}

func TestRagWorkerExecute(t *testing.T) {
	// given
	var (
		memoryStorage     = storage.NewMemoryStorage()
		ragClient         = ragger.NewMockRagClient()
		coordinatorClient = coordinator_client.NewMockCoordinatorClient()
		ragWorker         = NewRagWorker(ragClient, coordinatorClient, memoryStorage)

		websiteUrl = "https://example.com"
		markdown   = "Hello, world!"
		text       = "Hello, world!"
		chunks     = []string{"Hello, world!1"}
		embeddings = [][]float32{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}}

		contacts = []ragger.Contact{{Context: markdown, Value: "John Doe", Type: "person"}}
	)

	task, err := coordinator_client.NewTask("1", "test", RagWorkerParams{
		Markdown:  markdown,
		Url:       websiteUrl,
		InnerText: text,
	})
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	}

	ragClient.SetChunksFor(markdown, chunks)
	ragClient.SetContactsFor(markdown, contacts)

	newSlice := make([]string, len(chunks)+len(contacts))
	copy(newSlice, chunks)

	for i, contact := range contacts {
		newSlice[len(chunks)+i] = contact.Context
	}

	ragClient.SetEmbeddingsForAll(newSlice, embeddings)

	// when
	err = ragWorker.Execute(context.TODO(), task)
	if err != nil {
		t.Errorf("Error executing task: %v", err)
	}

	// then
	ragSources, err := storage.GetAll[storage.RagSource](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting rag sources: %v", err)
	}

	assert.Equal(t, len(ragSources), 1)
	assert.Equal(t, ragSources[0].URL, websiteUrl)

	rags, err := storage.GetAll[storage.RagChunk](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting chunks: %v", err)
	}

	assert.Equal(t, len(rags), 1)
	assert.Equal(t, rags[0].Text, chunks[0])
	assert.Equal(t, rags[0].Embedding, embeddings[0])
	assert.Equal(t, rags[0].RagSourceId, ragSources[0].ID)

	storedContacts, err := storage.GetAll[storage.RagContact](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting contacts: %v", err)
	}

	assert.Equal(t, len(storedContacts), 1)
	assert.Equal(t, storedContacts[0].Context, markdown)
	assert.Equal(t, storedContacts[0].Contact, "John Doe")
	assert.Equal(t, storedContacts[0].ContactType, "person")
	assert.Equal(t, storedContacts[0].Embedding, []float32{4.0, 5.0, 6.0})
	assert.Equal(t, storedContacts[0].RagSourceId, ragSources[0].ID)
}

func TestRagWorkerCleanup(t *testing.T) {
	// given
	var (
		coordinatorClient = coordinator_client.NewMockCoordinatorClient()
		ragWorker         = NewRagWorker(nil, coordinatorClient, nil)
	)

	task, err := coordinator_client.NewTask("1", "test", RagWorkerParams{
		Markdown: "Hello, world!",
		Url:      "https://example.com",
	})
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	}

	coordinatorClient.CreateTask(context.TODO(), coordinator_client.CoordinatorClientTaskTopicRag, task)
	if _, err := coordinatorClient.GetTaskAndSetProcessing(context.TODO(), 1*time.Second, coordinator_client.CoordinatorClientTaskTopicRag); err != nil {
		t.Errorf("Error getting task: %v", err)
	}

	// when
	err = ragWorker.Cleanup(context.TODO(), task)

	// then
	if err != nil {
		t.Errorf("Error cleaning up: %v", err)
	}

	if err := coordinatorClient.SetProcessed(context.TODO(), coordinator_client.CoordinatorClientTaskTopicRag, task); err != coordinator_client.ErrNoTasksCompleted {
		t.Errorf("There should be an error setting task to processed: %v", err)
	}
}

func TestRagWorkerASDASDAS(t *testing.T) {
	if os.Getenv("CICD") != "" {
		t.Skip("Skipping test in CICD")
	}

	f := `{"id":"735ad0a5-9541-4d94-aef1-ca3530383962","created_by":"3899939c-80cd-45bd-a9f4-75a42320ad90","params":{"markdown":"In this section\n\n- [Study](/study/)\n\n- [Submit your offer conditions](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/)\n- [Request a deferral](/study/apply/undergraduate/offer-holders/next-steps/deferred-entry/)\n\n[Study](/study/)\n\n- [Imperial Home](/)\n- [Study](/study/)\n- [Apply](/study/apply/)\n- [Undergraduate](/study/apply/undergraduate/)\n- [Offer holders](/study/apply/undergraduate/offer-holders/)\n- [Next steps](/study/apply/undergraduate/offer-holders/next-steps/)\n\n# Next steps\n\n[Offer holders](/study/apply/undergraduate/offer-holders/)\n\n- [Submit your offer conditions](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/)\n- [Request a deferral](/study/apply/undergraduate/offer-holders/next-steps/deferred-entry/)\n\n![](https://pxl-imperialacuk.terminalfour.net/fit-in/1079x305/prod01/channel_3/media/study---restricted-images/offer-holder-hub/next-steps/nextsteps_content_header_3000X850.jpg)\n\nCongratulations on your offer! What's next?\n\n## Find out your next steps\n\nWe've put together key information to help you understand your offer, how to reply and all the key dates and deadlines you need to know.\n\nYou can also find information about accommodation, funding your studies and how to send us evidence of any additional qualifications you need to meet your offer conditions.\n\n## Next steps\n\n- [1 – Reply to your offer](#tabpanel_1456127)\n- [2 – Consider your accommodation options](#tabpanel_1456128)\n- [3 – Explore your funding options](#tabpanel_1456129)\n- [4 – Tell us about additional qualifications](#tabpanel_1456130)\n- [5 – Welcome Season](#tabpanel_1456131)\n\n1 – Reply to your offer\n\n\nAfter receiving decisions on all your applications, it’s time to reply to your offers by the deadline.\n\nIf you receive all of your offers by **Wednesday 14 May 2025**, you must reply to UCAS by **Wednesday** **4 June 2025**.\n\nIf you decide to accept your offer, you can choose to make Imperial either your firm or insurance choice.\n\n#### Firm\n\nThis means we’re your first choice. If your offer is unconditional, then congratulations – you’re in! If it’s conditional, your place will be confirmed after you’ve met the conditions of your offer.\n\n#### Insurance\n\nThis means we’re your backup choice. You’ll be offered a place here if you don’t meet the conditions of your firm choice offer but meet the conditions of your Imperial offer.\n\nFind out how to reply to your offer on [UCAS](https://www.ucas.com/undergraduate/after-you-apply/types-offer/replying-your-ucas-undergraduate-offers \"UCAS\").\n\n2 – Consider your accommodation options\n\n\nImperial accommodation is well-maintained with round-the-clock support, and the rent includes all bills. If you pick Imperial as your firm choice, we offer a guaranteed place in our accommodation if you want it.\\*\n\nIf Imperial is your firm choice, you can apply for accommodation in late May. You’ll be invited to apply in late July if we're your insurance choice. You can choose your preferences based on your preferred hall, room type, and price.\n\nWe aim to make our accommodation the best choice for as many students as possible, but we know it won’t be suitable for everyone. [Our Accommodation Office](/students/accommodation/private-accommodation/) can also offer advice about renting privately and how to find the right place for you.\n\n[Visit the Accommodation Office website to learn more.](/students/accommodation/)\n\n\\*You’re guaranteed a place in our accommodation if you make Imperial your firm choice, submit your accommodation application by **18 July 2025**, meet all conditions of your offer and hold an unconditional offer by (date to be confirmed; as a guide, the unconditional offer deadline for 2024-25 entry was 23 August 2024).\n\n3 – Explore your funding options\n\n\nLearn more about tuition fees and funding options at Imperial, including our Imperial Bursary, student loans and scholarships.\n\nFind out more on our Offer Holder [fees and funding page](/study/apply/undergraduate/offer-holders/fees-and-funding/).\n\n4 – Tell us about additional qualifications\n\n\nUCAS will send us your A-level, International Baccalaureate (IB), Irish Leaving Certificates, Pre-U Certificate, Scottish Advanced Highers and STEP grades on your results day. However, you may need to send us the results for other qualifications.\n\nFind out if that [includes your qualifications](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/).\n\n5 – Welcome Season\n\n\nThe first month of our new term is Welcome Season, which **Welcome Week** is packed with opportunities to engage with the Imperial community for the first time.\n\nIt’s organised by Imperial and our Students’ Union and gives you the chance to make friends, get involved in loads of events and activities, and get to know your new home at Imperial.\n\n## Key dates\n\n- [May 2025](#tabpanel_1456133)\n- [June 2025](#tabpanel_1456134)\n- [July 2025](#tabpanel_1456135)\n- [August 2025](#tabpanel_1456136)\n- [September 2025](#tabpanel_1456137)\n\nMay 2025\n\n\n**Late May 2025 -** If you’ve made Imperial your firm choice on UCAS, you’ll be invited to apply for accommodation.\n\nJune 2025\n\n\n**5 June 2025 -** If you receive all of your offers by 14 May 2025, you must reply to UCAS by 5 June 2025.\n\nJuly 2025\n\n\n**5 July 2025 -** International Baccalaureate results day\n\n**18 July 2025** **-** Accommodation guarantee deadline\n\n**24 July 2025** **-** English Language and completed qualifications deadline\n\n**Late July** **-** If Imperial is your insurance choice on UCAS, you’ll be invited to apply for accommodation.\n\nAugust 2025\n\n\n**14 August 2025** **-** A-level results day\n\n**22 August 2025** **-** Deadline for receiving results that need additional, non-UCAS verification\n\n**31 August 2025** **-** Deadline to submit amended grades as a result of an appeal\n\nSeptember 2025\n\n\n**Early September 2025** **-** Accommodation offers sent out\n\n**27 September 2025** **-** Term starts\n\n[![](https://pxl-imperialacuk.terminalfour.net/fit-in/720x462/prod01/channel_3/media/study---restricted-images/offer-holder-hub/Contact-180227_southken_campus_snow_Pano.jpg)\\\n\\\n**Contact us** \\\n\\\nIf you have any questions about your offer, you can contact the relevant Admissions team.\\\n\\\nContact us](/study/apply/contact/)\n\n[![](https://pxl-imperialacuk.terminalfour.net/fit-in/720x462/prod01/channel_3/media/study---restricted-images/offer-holder-hub/Submit-offer-conditions.jpg)\\\n\\\n**Submit your offer conditions** \\\n\\\nWe may ask you to submit your results/certificates of qualifications that aren't sent to us via UCAS.\\\n\\\nFind out more](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/)\n\n[![](https://pxl-imperialacuk.terminalfour.net/fit-in/720x462/prod01/channel_3/media/study---restricted-images/offer-holder-hub/Request-a-deferral.jpg)\\\n\\\n**Request a deferral** \\\n\\\nYou may apply to defer your admission to the next academic year, under exceptional circumstances. \\\n\\\nFind out more](/study/apply/undergraduate/offer-holders/next-steps/deferred-entry/)\n\n## Offer holders\n\n[Next steps](/study/apply/undergraduate/offer-holders/next-steps/)\n\n[Life on campus](/study/apply/undergraduate/offer-holders/life-on-campus/)\n\n[Accommodation](/study/apply/undergraduate/offer-holders/accommodation/)\n\n[Fees and funding](/study/apply/undergraduate/offer-holders/fees-and-funding/)\n\n[Support and FAQs](/study/apply/undergraduate/offer-holders/support/)\n\n[After you graduate](/study/apply/undergraduate/offer-holders/after-you-graduate/)\n\n[International students](/study/apply/undergraduate/offer-holders/international-students/)\n\n[Return to offer holder home](/study/apply/undergraduate/offer-holders/)\n\n[Accept your offer](https://www.ucas.com/undergraduate/after-you-apply/ucas-undergraduate-types-offer)\n\n[Submit your offer conditions](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/)","url":"https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/next-steps/"}}`

	var task coordinator_client.Task
	if err := json.Unmarshal([]byte(f), &task); err != nil {
		t.Errorf("Error unmarshalling task: %v", err)
	}

	ragClient := ragger.NewRAGClient("../model/model.onnx", "../libonnxruntime.so.1.20.1", "../model/tokenizer.json")
	storageClinet := storage.NewMemoryStorage()
	cooridnatorClient := coordinator_client.NewMockCoordinatorClient()
	ragWorker := NewRagWorker(ragClient, cooridnatorClient, storageClinet)

	err := ragWorker.Execute(context.TODO(), &task)
	if err != nil {
		t.Errorf("Error executing task: %v", err)
	}

	fmt.Printf("%+v\n", task)
}
