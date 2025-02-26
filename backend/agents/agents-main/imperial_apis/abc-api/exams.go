package abc_api

type ExamResponse struct {
	Period            string `json:"period"`
	Term              int    `json:"term"`
	ExamDate          string `json:"exam_date"`
	StartTime         string `json:"start_time"`
	Duration          int    `json:"duration"`
	Title             string `json:"title"`
	ModuleCode        string `json:"module_code"`
	ComputerBased     bool   `json:"computer_based_exam"`
	AnswerbookExamURL string `json:"answerbook_exam_url"`
}

func (m *MockAbcApiClient) GetExams() []ExamResponse {
	return []ExamResponse{
		{
			Period:            "2024-2025",
			Term:              1,
			ExamDate:          "2024-12-09",
			StartTime:         "14:30",
			Duration:          90,
			Title:             "Introduction to Machine Learning (Term 1)",
			ModuleCode:        "60012",
			ComputerBased:     true,
			AnswerbookExamURL: "https://exams.doc.ic.ac.uk/60012/2024",
		},
		{
			Period:            "2024-2025",
			Term:              1,
			ExamDate:          "2024-12-11",
			StartTime:         "10:00",
			Duration:          90,
			Title:             "Mathematics for Machine Learning",
			ModuleCode:        "70015",
			ComputerBased:     false,
			AnswerbookExamURL: "https://exams.doc.ic.ac.uk/70015/2024",
		},
		{
			Period:            "2024-2025",
			Term:              1,
			ExamDate:          "2024-12-11",
			StartTime:         "14:30",
			Duration:          120,
			Title:             "Advanced Computer Architecture",
			ModuleCode:        "60001",
			ComputerBased:     true,
			AnswerbookExamURL: "https://exams.doc.ic.ac.uk/60001/2024",
		},
		{
			Period:            "2024-2025",
			Term:              1,
			ExamDate:          "2024-12-12",
			StartTime:         "10:00",
			Duration:          120,
			Title:             "The Theory and Practice of Concurrent Programming",
			ModuleCode:        "60007",
			ComputerBased:     true,
			AnswerbookExamURL: "https://exams.doc.ic.ac.uk/60007/2024",
		},
		{
			Period:            "2024-2025",
			Term:              2,
			ExamDate:          "2025-03-17",
			StartTime:         "14:00",
			Duration:          120,
			Title:             "Graphics",
			ModuleCode:        "60005",
			ComputerBased:     false,
			AnswerbookExamURL: "https://exams.doc.ic.ac.uk/60005/2025",
		},
		{
			Period:            "2024-2025",
			Term:              2,
			ExamDate:          "2025-03-18",
			StartTime:         "10:00",
			Duration:          120,
			Title:             "Custom Computing",
			ModuleCode:        "60008",
			ComputerBased:     true,
			AnswerbookExamURL: "https://exams.doc.ic.ac.uk/60008/2025",
		},
		{
			Period:            "2024-2025",
			Term:              2,
			ExamDate:          "2025-03-19",
			StartTime:         "14:00",
			Duration:          90,
			Title:             "Computer Vision",
			ModuleCode:        "60006",
			ComputerBased:     true,
			AnswerbookExamURL: "https://exams.doc.ic.ac.uk/60006/2025",
		},
	}
}
