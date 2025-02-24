package abc_api

type DegreeRegulationsResponse struct {
	DegreeCode     string           `json:"degree_code"`
	Cohort         string           `json:"cohort"`
	CreditCriteria string           `json:"credit_criteria"`
	DegreeTitle    string           `json:"degree_title"`
	DegreeYear     int              `json:"degree_year"`
	Regulations    []RegulationItem `json:"regulations"`
}

type RegulationItem struct {
	Label            string       `json:"label"`
	MaximumSelection int          `json:"maximum_selection"`
	MinimumSelection int          `json:"minimum_selection"`
	Modules          []ModuleItem `json:"modules"`
}

type ModuleItem struct {
	Code                   string   `json:"code"`
	Title                  string   `json:"title"`
	Terms                  []int    `json:"terms"`
	ApplicableCohorts      []string `json:"applicable_cohorts"`
	ExamContribution       *int     `json:"exam_contribution,omitempty"`
	CourseworkContribution *int     `json:"coursework_contribution,omitempty"`
	ExamQuestionsTotal     *int     `json:"exam_questions_total,omitempty"`
	ExamQuestionsNeeded    *int     `json:"exam_questions_needed,omitempty"`
}

func intPtr(i int) *int {
	return &i
}

func (m *MockAbcApiClient) GetDegreeRegulations() DegreeRegulationsResponse {

	examContr80 := 80
	cwContr20 := 20
	examContr70 := 70
	cwContr30 := 30
	examQTotal3 := 3
	examQNeeded3 := 3
	examQTotal2 := 2
	examQNeeded2 := 2

	return DegreeRegulationsResponse{
		DegreeCode:     "beng",
		Cohort:         "c3",
		CreditCriteria: "Count",
		DegreeTitle:    "B.Eng Computing",
		DegreeYear:     3,
		Regulations: []RegulationItem{
			{
				Label:            "Extracurricular",
				MaximumSelection: 0,
				MinimumSelection: 0,
				Modules: []ModuleItem{
					{
						Code:              "COMPM0804",
						Title:             "Student Support and Wellbeing",
						Terms:             []int{1, 2, 3},
						ApplicableCohorts: []string{"c3", "j3", "i3"},
					},
				},
			},
			{
				Label:            "Core",
				MaximumSelection: 4,
				MinimumSelection: 4,
				Modules: []ModuleItem{
					{
						Code:                   "60001",
						Title:                  "Advanced Computer Architecture",
						Terms:                  []int{1},
						ApplicableCohorts:      []string{"c3", "c4", "i3", "i4", "j3", "j4", "o3", "x1", "x3", "x5"},
						ExamContribution:       intPtr(examContr80),
						CourseworkContribution: intPtr(cwContr20),
						ExamQuestionsTotal:     intPtr(examQTotal3),
						ExamQuestionsNeeded:    intPtr(examQNeeded3),
					},
					{
						Code:                   "60005",
						Title:                  "Graphics",
						Terms:                  []int{2},
						ApplicableCohorts:      []string{"c3", "c4", "i3", "i4", "j3", "j4", "o3", "x1", "x3", "x5"},
						ExamContribution:       intPtr(examContr80),
						CourseworkContribution: intPtr(cwContr20),
						ExamQuestionsTotal:     intPtr(examQTotal3),
						ExamQuestionsNeeded:    intPtr(examQNeeded3),
					},
					{
						Code:                   "60006",
						Title:                  "Computer Vision",
						Terms:                  []int{2},
						ApplicableCohorts:      []string{"c3", "c4", "i3", "i4", "j3", "j4", "o3", "x1", "x3", "x5"},
						ExamContribution:       intPtr(examContr70),
						CourseworkContribution: intPtr(cwContr30),
						ExamQuestionsTotal:     intPtr(examQTotal3),
						ExamQuestionsNeeded:    intPtr(examQNeeded3),
					},
					{
						Code:                   "60007",
						Title:                  "The Theory and Practice of Concurrent Programming",
						Terms:                  []int{1},
						ApplicableCohorts:      []string{"c3", "j3"},
						ExamContribution:       intPtr(examContr80),
						CourseworkContribution: intPtr(cwContr20),
						ExamQuestionsTotal:     intPtr(examQTotal2),
						ExamQuestionsNeeded:    intPtr(examQNeeded2),
					},
					{
						Code:                   "60008",
						Title:                  "Custom Computing",
						Terms:                  []int{2},
						ApplicableCohorts:      []string{"c3", "c4", "i3", "i4", "j3", "j4", "o3", "x1", "x3", "x5"},
						ExamContribution:       intPtr(examContr80),
						CourseworkContribution: intPtr(cwContr20),
						ExamQuestionsTotal:     intPtr(examQTotal3),
						ExamQuestionsNeeded:    intPtr(examQNeeded3),
					},
					{
						Code:                   "60012",
						Title:                  "Introduction to Machine Learning",
						Terms:                  []int{1},
						ApplicableCohorts:      []string{"c3", "i3", "j3", "o3"},
						ExamContribution:       intPtr(examContr70),
						CourseworkContribution: intPtr(cwContr30),
						ExamQuestionsTotal:     intPtr(examQTotal3),
						ExamQuestionsNeeded:    intPtr(examQNeeded3),
					},
					{
						Code:                   "70015",
						Title:                  "Mathematics for Machine Learning",
						Terms:                  []int{1},
						ApplicableCohorts:      []string{"a5", "c3", "c4", "i4", "o3", "q5", "r6", "s5", "t5"},
						ExamContribution:       intPtr(examContr70),
						CourseworkContribution: intPtr(cwContr30),
						ExamQuestionsTotal:     intPtr(examQTotal3),
						ExamQuestionsNeeded:    intPtr(examQNeeded3),
					},
				},
			},
		},
	}
}
