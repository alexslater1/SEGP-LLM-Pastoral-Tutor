package abc_api

type ModuleConsolidationResponse struct {
	AcademicYear           string `json:"academic_year"`
	ExamClass              string `json:"exam_class"`
	ExamCode               string `json:"exam_code"`
	PassMark               int    `json:"pass_mark"`
	ExamContribution       int    `json:"exam_contribution"`
	CourseworkContribution int    `json:"coursework_contribution"`
	ExamQuestionsTotal     int    `json:"exam_questions_total"`
	ExamQuestionsNeeded    int    `json:"exam_questions_needed"`
}

func (m *MockAbcApiClient) GetModuleConsolidation() (*[]ModuleConsolidationResponse, error) {
	return &[]ModuleConsolidationResponse{
		{
			AcademicYear:           "2324",
			ExamClass:              "bm2",
			ExamCode:               "COMP50002",
			PassMark:               40,
			ExamContribution:       80,
			CourseworkContribution: 20,
			ExamQuestionsTotal:     3,
			ExamQuestionsNeeded:    3,
		},
	}, nil
}

type ExamRegistrationResponse struct {
	ExamCode string   `json:"exam_code"`
	Students []string `json:"students"`
}

func (m *MockAbcApiClient) GetExamRegistrations() (*[]ExamRegistrationResponse, error) {
	return &[]ExamRegistrationResponse{
		{
			ExamCode: "COMP50002",
			Students: []string{"123456", "456789"},
		},
	}, nil
}
