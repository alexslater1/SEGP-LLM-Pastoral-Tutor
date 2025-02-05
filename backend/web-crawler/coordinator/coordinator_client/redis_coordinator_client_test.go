package coordinator_client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	os.Exit(m.Run())
}

func TestRedisCoordinatorClientCreateTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	taskParams := map[string]string{
		"url": "https://ethanhosier.com",
	}

	task, err := NewTask("2b665be2-80b7-40d4-9117-a6e9794afe97", "asdasdasd", taskParams)
	if err != nil {
		t.Fatalf("Failed to    create ta sk:   %v", err)
	}

	err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task      : %v", err)
	}
}

func TestRedisCoordinatorClientCreate100Tasks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	for i := 0; i < 100; i++ {
		params := map[string]string{
			"url": "https://ethanhosier.com?test=" + strconv.Itoa(i),
		}

		task, err := NewTask(uuid.New().String(), "asdasdasd", params)
		if err != nil {
			t.Fatalf("Failed to  create  task: %v", err)
		}

		err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
		if err != nil {
			t.Fatalf("Failed to crea te task: %v", err)
		}
	}
}

func TestRedisCoordinatorClientCreate100UrlTasks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip(" Skipping test in CICD                      ")
	}

	links := `https://www.imperial.ac.uk/
https://www.imperial.ac.uk/study/
https://www.imperial.ac.uk/study/courses/
https://www.imperial.ac.uk/study/subjects/
https://www.imperial.ac.uk/study/apply/
https://www.imperial.ac.uk/study/apply/undergraduate/
https://www.imperial.ac.uk/study/apply/undergraduate/entry-requirements/
https://www.imperial.ac.uk/study/apply/undergraduate/entry-requirements/accepted-qualifications/
https://www.imperial.ac.uk/study/apply/undergraduate/process/
https://www.imperial.ac.uk/study/apply/undergraduate/process/choose-course/
https://www.imperial.ac.uk/study/apply/undergraduate/process/deadlines/
https://www.imperial.ac.uk/study/apply/undergraduate/process/personal-statement/
https://www.imperial.ac.uk/study/apply/undergraduate/process/reference/
https://www.imperial.ac.uk/study/apply/undergraduate/process/reference/teachers/
https://www.imperial.ac.uk/study/apply/undergraduate/process/admissions-schemes/
https://www.imperial.ac.uk/study/apply/undergraduate/process/admissions-tests/
https://www.imperial.ac.uk/study/apply/undergraduate/process/admissions-tests/esat/
https://www.imperial.ac.uk/study/apply/undergraduate/process/admissions-tests/tmua/
https://www.imperial.ac.uk/study/apply/undergraduate/process/admissions-tests/tmua/important-information-tmua-october-2024-sitting/
https://www.imperial.ac.uk/study/apply/undergraduate/process/admissions-tests/ucat/
https://www.imperial.ac.uk/study/apply/undergraduate/process/interviews/
https://www.imperial.ac.uk/study/apply/undergraduate/process/selection/
https://www.imperial.ac.uk/study/apply/undergraduate/process/results/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/next-steps/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/next-steps/deferred-entry/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/life-on-campus/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/accommodation/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/fees-and-funding/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/support/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/after-you-graduate/
https://www.imperial.ac.uk/study/apply/undergraduate/offer-holders/international-students/
https://www.imperial.ac.uk/study/apply/undergraduate/advisers/
https://www.imperial.ac.uk/study/apply/undergraduate/advisers/upcoming-events/
https://www.imperial.ac.uk/study/apply/undergraduate/advisers/application-advice/
https://www.imperial.ac.uk/study/apply/undergraduate/advisers/visits/
https://www.imperial.ac.uk/study/apply/undergraduate/advisers/counsellor-faqs/
https://www.imperial.ac.uk/study/apply/undergraduate/advisers/academic-resources/
https://www.imperial.ac.uk/study/apply/undergraduate/advisers/advisers-newsletter/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/entry-requirements/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/entry-requirements/accepted-qualifications/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/application-process/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/application-process/choose-course/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/application-process/personal-statement/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/application-process/deadlines/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/application-process/application-fee/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/application-process/application-fee/application-fee-waiver/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/application-process/reference/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/application-process/interviews/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/next-steps/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/next-steps/submit-offer-conditions/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/next-steps/deferred-entry/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/life-on-campus/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/accommodation/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/fees-and-funding/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/support/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/after-you-graduate/
https://www.imperial.ac.uk/study/apply/postgraduate-taught/offer-holders/international-students/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/entry-requirements/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/entry-requirements/accepted-qualifications/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/choose-course/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/choose-course/phd/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/choose-course/split-phd/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/choose-course/professional-doctorate/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/choose-course/integrated-phd/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/choose-course/pri-scheme/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/choose-course/advanced-standing/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/reference/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/supervisor/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/research-proposal/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/application-process/interview/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/next-steps/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/next-steps/submit-offer-conditions/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/next-steps/deferred-entry/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/next-steps/submit-offer-conditions/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/life-on-campus/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/accommodation/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/fees-and-funding/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/support/
https://www.imperial.ac.uk/study/apply/postgraduate-doctoral/offer-holders/international-students/
https://www.imperial.ac.uk/study/apply/english-language/
https://www.imperial.ac.uk/study/apply/english-language/english-language-exemption/
https://www.imperial.ac.uk/study/apply/visiting-students/
https://www.imperial.ac.uk/study/apply/contact/
https://www.imperial.ac.uk/study/fees-and-funding/
https://www.imperial.ac.uk/study/fees-and-funding/undergraduate/
https://www.imperial.ac.uk/study/fees-and-funding/undergraduate/tuition-fees/
https://www.imperial.ac.uk/study/fees-and-funding/undergraduate/bursaries-grants-scholarships/
https://www.imperial.ac.uk/study/fees-and-funding/undergraduate/bursaries-grants-scholarships/imperial-bursary/
https://www.imperial.ac.uk/study/fees-and-funding/undergraduate/bursaries-grants-scholarships/sanctuary/
https://www.imperial.ac.uk/study/fees-and-funding/undergraduate/bursaries-grants-scholarships/nhs-bursary/
https://www.imperial.ac.uk/study/fees-and-funding/scholarships-search/
https://www.imperial.ac.uk/study/fees-and-funding/undergraduate/bursaries-grants-scholarships/ib-excellence/
https://www.imperial.ac.uk/study/fees-and-funding/undergraduate/bursaries-grants-scholarships/presidential-scholarships-black-heritage-students/`

	client := NewRedisCoordinatorClient(context.Background(), "18.133.156.65:6379", "password", 0)
	urls := strings.Split(links, "\n")

	for _, url := range urls {
		params := map[string]string{
			"url": url,
		}

		task, err := NewTask(uuid.New().String(), "asdasdsasd", params)
		if err != nil {
			t.Fatalf("Failed  to create task: %v", err)
		}

		err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
		if err != nil {
			t.Fatalf("Failed to   create task: %v", err)
		}
	}
}

func TestRedisCoordinatorClientGetTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}
	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	task, err := client.GetTask(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	t.Logf("Task: %+v", task)
}

func TestRedisCoordinatorClientGetTaskAndSetProcessing(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	task, err := client.GetTaskAndSetProcessing(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	t.Logf("Task: %+v", task)
}

func TestRedisCoordinatorClientSetProcessed(t *testing.T) {
	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", "test")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD ")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	client.SetProcessed(context.Background(), CoordinatorClientTaskTopicUrls, task)
}

func TestRedisCoordinatorClientCreateGetTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	type TestStruct struct {
		Number int    `json:"number"`
		Name   string `json:"name"`
	}

	params := TestStruct{Number: 1, Name: "a name"}

	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", params)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	task, err = client.GetTask(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	parsedParams, err := CastParams[TestStruct](task.Params)
	if err != nil {
		t.Fatalf("Failed to parse params: %v", err)
	}

	assert.Equal(t, parsedParams.Number, params.Number)
	assert.Equal(t, parsedParams.Name, params.Name)
}

func TestRedisCoordinatorClientStoreError(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	params := map[string]string{
		"url": "https://ethanhosier.com",
	}

	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", params)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = client.StoreError(context.Background(), CoordinatorClientTaskTopicUrls, task, fmt.Errorf("an error"))
	if err != nil {
		t.Fatalf("Failed to store error : %v", err)
	}
}

func TestRedisCoordinatorClientCreateTasks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	params1 := map[string]string{
		"url": "https://ethanhosier.com",
	}
	task1, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", params1)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	params2 := map[string]string{
		"url": "https://ethanhosier.com/blog",
	}
	task2, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", params2)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	client.CreateTasks(context.Background(), CoordinatorClientTaskTopicUrls, []*Task{task1, task2})
}

func TestRedisCoordinatorClientNumTasks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	numTasks, err := client.NumTasks(context.Background(), CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get number of tasks: %v", err)
	}

	assert.Equal(t, numTasks, 5)
}

func TestRedisCoordinatorClientNumProcessingTasks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	numProcessingTasks, err := client.NumProcessingTasks(context.Background(), CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get number of processing tasks: %v", err)
	}

	assert.Equal(t, numProcessingTasks, 0)
}

func TestRedisCoordinatorClientGetErrors(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	errors, err := client.GetErrors(context.Background(), CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get errors: %v", err)
	}

	assert.Equal(t, len(errors), 0)
	params := map[string]string{
		"url": "https://ethanhosier.com",
	}
	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", params)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if err := client.StoreError(context.Background(), CoordinatorClientTaskTopicUrls, task, fmt.Errorf("an error")); err != nil {
		t.Fatalf("Failed to store error: %v", err)
	}

	errors, err = client.GetErrors(context.Background(), CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get errors: %v", err)
	}

	assert.Equal(t, len(errors), 1)
	assert.Equal(t, errors[0].Error, "an error")
}
