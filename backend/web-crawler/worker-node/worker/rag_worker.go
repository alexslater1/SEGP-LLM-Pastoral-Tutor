package worker

type RagWorkerParams struct {
	Markdown string `json:"markdown"`
	Url      string `json:"url"`
}
