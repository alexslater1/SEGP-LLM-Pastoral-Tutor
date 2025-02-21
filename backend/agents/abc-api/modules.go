package abc_api

/*
[
  {
    "code": "50002",
    "title": "Software Engineering Design",
    "terms": [
      1,
      2
    ],
    "applicable_cohorts": [
      "c3",
      "j3"
    ],
    "cohort_regulations": [
      {
        "cohort": "c3",
        "pass_mark": 40
      },
      {
        "cohort": "j3",
        "pass_mark": 40
      }
    ],
    "exam_contribution": 60,
    "coursework_contribution": 40,
    "exam_questions_total": 4,
    "exam_questions_needed": 2,
    "ects": 5,
    "staff": [
      {
        "login": "rbc",
        "email": "rbc@ic.ac.uk",
        "firstname": "Rob",
        "lastname": "Chatley",
        "role_in_department": "staff",
        "department": "DoC",
        "cid": "123456"
      },
      {
        "login": "ip914",
        "email": "ip914@ic.ac.uk",
        "firstname": "Ivan",
        "lastname": "Procaccini",
        "role_in_department": "staff",
        "department": "DoC",
        "cid": "654321"
      }
    ],
    "helpers": [
      {
        "login": "hgranger",
        "email": "hgranger@ic.ac.uk",
        "firstname": "Hermione",
        "lastname": "Granger",
        "role_in_department": "student",
        "department": "DoC"
      }
    ]
  }


*/

type CourseResponse struct {
	Code                   string                     `json:"code"`
	Title                  string                     `json:"title"`
	Terms                  []int                      `json:"terms"`
	ApplicableCohorts      []string                   `json:"applicable_cohorts"`
	CohortRegulations      []CohortRegulationResponse `json:"cohort_regulations"`
	ExamContribution       int                        `json:"exam_contribution"`
	CourseworkContribution int                        `json:"coursework_contribution"`
	ExamQuestionsTotal     int                        `json:"exam_questions_total"`
	ExamQuestionsNeeded    int                        `json:"exam_questions_needed"`
	Ects                   int                        `json:"ects"`
	Staff                  []StaffResponse            `json:"staff"`
	Helpers                []HelperResponse           `json:"helpers"`
}

type CohortRegulationResponse struct {
	Cohort   string `json:"cohort"`
	PassMark int    `json:"pass_mark"`
}

type StaffResponse struct {
	Login            string `json:"login"`
	Email            string `json:"email"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	RoleInDepartment string `json:"role_in_department"`
	Department       string `json:"department"`
	Cid              string `json:"cid"`
}

type HelperResponse struct {
	Login            string `json:"login"`
	Email            string `json:"email"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	RoleInDepartment string `json:"role_in_department"`
	Department       string `json:"department"`
}

func (m *MockAbcApiClient) GetModules() (*CourseResponse, error) {
	return &CourseResponse{
		Code:              "50002",
		Title:             "Software Engineering Design",
		Terms:             []int{1, 2},
		ApplicableCohorts: []string{"c3", "j3"},
		CohortRegulations: []CohortRegulationResponse{
			{Cohort: "c3", PassMark: 40},
			{Cohort: "j3", PassMark: 40},
		},
		ExamContribution:       60,
		CourseworkContribution: 40,
		ExamQuestionsTotal:     4,
		ExamQuestionsNeeded:    2,
		Ects:                   5,
		Staff: []StaffResponse{
			{
				Login:            "rbc",
				Email:            "rbc@ic.ac.uk",
				FirstName:        "Rob",
				LastName:         "Chatley",
				RoleInDepartment: "staff",
				Department:       "DoC",
				Cid:              "123456",
			},
			{
				Login:            "ip914",
				Email:            "ip914@ic.ac.uk",
				FirstName:        "Ivan",
				LastName:         "Procaccini",
				RoleInDepartment: "staff",
				Department:       "DoC",
				Cid:              "654321",
			},
		},
		Helpers: []HelperResponse{
			{
				Login:            "hgranger",
				Email:            "hgranger@ic.ac.uk",
				FirstName:        "Hermione",
				LastName:         "Granger",
				RoleInDepartment: "student",
				Department:       "DoC",
			},
		},
	}, nil
}
