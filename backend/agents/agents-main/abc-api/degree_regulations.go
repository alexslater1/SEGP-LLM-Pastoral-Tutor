package abc_api

var (
	examContrib     = 60
	cwContrib       = 40
	totalQuestions  = 4
	neededQuestions = 2
)

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

func (m *MockAbcApiClient) GetDegreeRegulations() DegreeRegulationsResponse {
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
						Code:                   "50002",
						Title:                  "Software Engineering Design",
						Terms:                  []int{1, 2},
						ApplicableCohorts:      []string{"c3", "j3"},
						ExamContribution:       &examContrib,
						CourseworkContribution: &cwContrib,
						ExamQuestionsTotal:     &totalQuestions,
						ExamQuestionsNeeded:    &neededQuestions,
					},
				},
			},
		},
	}
}
