type ExamResponse struct {
    Period            string    `json:"period"`
    Term              int       `json:"term"`
    ExamDate          string    `json:"exam_date"`
    StartTime         string    `json:"start_time"`
    Duration          int       `json:"duration"`
    Title            string    `json:"title"`
    ModuleCode       string    `json:"module_code"`
    ComputerBased    bool      `json:"computer_based_exam"`
    AnswerbookExamURL string    `json:"answerbook_exam_url"`
}

func (m *MockAbcApiClient) GetExams() (*ExamResponse, error) {
	return &ExamResponse{
		Period:            "2024-2025",
		Term:              2,
		ExamDate:          "2025-05-15",
		StartTime:         "14:00",
		Duration:          120,
		Title:             "Software Engineering Design",
		ModuleCode:       "50002",
		ComputerBased:    true,
		AnswerbookExamURL: "https://exams.doc.ic.ac.uk/50002/2025",
	}, nil
}
