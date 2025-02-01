package worker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	coordinator_client "github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/utils"
	"github.com/google/uuid"
)

type ScraperWorker struct {
	scraper           scraper.Scraper
	coordinatorClient coordinator_client.CoordinatorClient
	id                string
}

type ScraperWorkerParams struct {
	Url string
}

func (w *ScraperWorkerParams) WorkerType() WorkerType {
	return WorkerTypeScraper
}

type ScraperWorkerResult struct {
	Markdown string
	Url      string
}

func (w *ScraperWorker) WorkerType() WorkerType {
	return WorkerTypeScraper
}

func NewScraperWorker(scraper scraper.Scraper, coordinatorClient coordinator_client.CoordinatorClient) *ScraperWorker {
	id := uuid.New().String()
	return &ScraperWorker{scraper: scraper, coordinatorClient: coordinatorClient, id: id}
}

func (w *ScraperWorker) Id() string {
	return w.id
}

func (w *ScraperWorker) Execute(ctx context.Context, task *coordinator_client.Task) error {
	scraperParams, err := coordinator_client.CastParams[ScraperWorkerParams](task.Params)
	if err != nil {
		return fmt.Errorf("invalid params %+v", task.Params)
	}

	if scraperParams.Url == "" {
		return fmt.Errorf("url is required")
	}

	md, text, err := w.mdAndTextFromUrl(scraperParams.Url)
	if err != nil {
		return err
	}

	if md == "" {
		log.Printf("No markdown parsed for %s. No need to rag", scraperParams.Url)
		return nil
	}

	ragParams := RagWorkerParams{Markdown: md, Url: scraperParams.Url, InnerText: text}

	ragTask, err := coordinator_client.NewTask(uuid.New().String(), w.id, ragParams)
	if err != nil {
		return err
	}

	return w.coordinatorClient.CreateTask(ctx, coordinator_client.CoordinatorClientTaskTopicRag, ragTask)
}

func (w *ScraperWorker) Cleanup(ctx context.Context, task *coordinator_client.Task) error {
	return w.coordinatorClient.SetProcessed(ctx, coordinator_client.CoordinatorClientTaskTopicUrls, task)
}

func (w *ScraperWorker) mdAndTextFromUrl(url string) (string, string, error) {
	html, err := w.scraper.HtmlFromTag(url, "main")
	if err != nil {
		return "", "", err
	}

	h, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return "", "", err
	}

	md, err := utils.HtmlToMarkdown(html)
	if err != nil {
		return "", "", err
	}

	return md, h.Text(), nil
}
