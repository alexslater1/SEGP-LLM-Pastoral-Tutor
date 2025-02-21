package abc_api

type ModuleResponse3 struct {
	Code                   string                     `json:"code"`
	Title                  string                     `json:"title"`
	Terms                  []int                      `json:"terms"`
	ApplicableCohorts      []string                   `json:"applicable_cohorts"`
	CohortRegulations      []CohortRegulationResponse `json:"cohort_regulations"`
	ExamContribution       int                        `json:"exam_contribution"`
	CourseworkContribution int                        `json:"coursework_contribution"`
	ExamQuestionsTotal     int                        `json:"exam_questions_total"`
	ExamQuestionsNeeded    int                        `json:"exam_questions_needed"`
	ECTS                   int                        `json:"ects"`
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
	Salutation       string `json:"salutation"`
	RoleInDepartment string `json:"role_in_department"`
	Department       string `json:"department"`
	CID              string `json:"cid"`
}

type HelperResponse struct {
	Login            string `json:"login"`
	Email            string `json:"email"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	RoleInDepartment string `json:"role_in_department"`
	Department       string `json:"department"`
}

type EnrolledStudentResponse struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Level     int    `json:"level"`
	Status    string `json:"status"`
	Cohort    string `json:"cohort"`
}

// For first term modules
func (m *MockAbcApiClient) GetModules() ([]*ModuleResponse3, error) {
	return []*ModuleResponse3{
		{
			Code:        "60001",
			Title:       "Advanced Computer Architecture",
			ECTS:        5,
			Terms:       []int{1},
			ApplicableCohorts: []string{"c3", "c4", "i3", "i4", "j3", "j4", "o3", "x1", "x3", "x5"},
			CohortRegulations: []CohortRegulationResponse{
				{Cohort: "c3", PassMark: 40},
				{Cohort: "i3", PassMark: 40},
				{Cohort: "j3", PassMark: 40},
				{Cohort: "o3", PassMark: 40},
				{Cohort: "x3", PassMark: 40},
				{Cohort: "c4", PassMark: 40},
				{Cohort: "i4", PassMark: 40},
				{Cohort: "j4", PassMark: 40},
				{Cohort: "x1", PassMark: 40},
				{Cohort: "x5", PassMark: 40},
			},
			ExamContribution: 80,
			CourseworkContribution: 20,
			ExamQuestionsTotal: 3,
			ExamQuestionsNeeded: 3,
			Staff: []StaffResponse{
				{
					Login:          "phjk",
					Email:          "p.kelly@imperial.ac.uk",
					RoleInDepartment: "staff",
					LastName:       "Kelly",
					FirstName:      "Paul",
					Department:     "DoC",
					CID:            "00003206",
				},
			},
			Helpers: []HelperResponse{
				{
					Login:    "lp721",
					Email:    "l.panayi21@imperial.ac.uk",
					LastName: "Panayi",
					FirstName: "Luke",
					RoleInDepartment:    "Tutorial Helper",
					Department:       "DoC",
				},
			},
		},
		{
			Code:        "60007",
			Title:       "The Theory and Practice of Concurrent Programming",
			ECTS:        5,
			Terms:       []int{1},
			ApplicableCohorts: []string{"c3", "j3"},
			CohortRegulations: []CohortRegulationResponse{
				{Cohort: "c3", PassMark: 40},
				{Cohort: "j3", PassMark: 40},
			},
			ExamContribution: 80,
			CourseworkContribution: 20,
			ExamQuestionsTotal: 2,
			ExamQuestionsNeeded: 2,
			Staff: []StaffResponse{
				{
					Login:          "azalea",
					Email:          "azalea.raad@imperial.ac.uk",
					RoleInDepartment: "staff",
					LastName:       "Raad",
					FirstName:      "Azalea",
					Department:     "DoC",
					CID:            "00483298",
				},
			},
			Helpers: []HelperResponse{
				{
					Login:    "sh2221",
					Email:    "shinghin.ho21@imperial.ac.uk",
					LastName: "Ho",
					FirstName: "Shing",
					RoleInDepartment:    "Tutorial Helper",
					Department:       "DoC",
				},
			},
		},
		{
			Code:        "60012",
			Title:       "Introduction to Machine Learning",
			ECTS:        5,
			Terms:       []int{1},
			ApplicableCohorts: []string{"c3", "i3", "j3", "o3"},
			CohortRegulations: []CohortRegulationResponse{
				{Cohort: "c3", PassMark: 40},
				{Cohort: "i3", PassMark: 40},
				{Cohort: "j3", PassMark: 40},
				{Cohort: "o3", PassMark: 40},
			},
			ExamContribution: 70,
			CourseworkContribution: 30,
			ExamQuestionsTotal: 3,
			ExamQuestionsNeeded: 3,
			Staff: []StaffResponse{
				{
					Login:          "jwang4",
					Email:          "josiah.wang@imperial.ac.uk",
					RoleInDepartment: "staff",
					LastName:       "Wang",
					FirstName:      "Josiah",
					Department:     "DoC",
					CID:            "01030000",
				},
			},
			Helpers: []HelperResponse{
				{
					Login:    "ad5518",
					Email:    "adam.dejl18@imperial.ac.uk",
					LastName: "Dejl",
					FirstName: "Adam",
					RoleInDepartment:    "Tutorial Helper",
					Department:       "DoC",
				},
			},
		},
		{
		Code:   "70015",
		Title:  "Mathematics for Machine Learning",
		Terms:  []int{1},
		ApplicableCohorts: []string{
			"a5", "c3", "c4", "i4", "o3", "q5", "r6", "s5", "t5",
		},
		CohortRegulations: []CohortRegulationResponse{
			{"r6", 50},
			{"t5", 50},
			{"c3", 50},
			{"c4", 50},
			{"i4", 50},
			{"s5", 50},
			{"a5", 50},
			{"q5", 50},
			{"o3", 50},
		},
		ExamContribution:       70,
		CourseworkContribution: 30,
		ExamQuestionsTotal:     3,
		ExamQuestionsNeeded:    3,
		ECTS:                   5,
		Staff: []StaffResponse{
			{
				Login:            "rac101",
				Email:            "robert.craven@imperial.ac.uk",
				FirstName:        "Robert",
				LastName:         "Craven",
				RoleInDepartment: "staff",
				Department:       "DoC",
				CID:              "00343970",
			},
		},
		Helpers: []HelperResponse{
			{
				Login:            "zo122",
				Email:            "z.ou22@imperial.ac.uk",
				FirstName:        "Zijing",
				LastName:         "Ou",
				RoleInDepartment: "Tutorial Helper",
				Department:       "DoC",
			},
		},
	},
	}, nil

}

// For Intro to ML
func (m *MockAbcApiClient) GetModule() (*ModuleResponse3, error) {
	return &ModuleResponse3{
			Code:        "60012",
			Title:       "Introduction to Machine Learning",
			ECTS:        5,
			Terms:       []int{1},
			ApplicableCohorts: []string{"c3", "i3", "j3", "o3"},
			CohortRegulations: []CohortRegulationResponse{
				{Cohort: "c3", PassMark: 40},
				{Cohort: "i3", PassMark: 40},
				{Cohort: "j3", PassMark: 40},
				{Cohort: "o3", PassMark: 40},
			},
			ExamContribution: 70,
			CourseworkContribution: 30,
			ExamQuestionsTotal: 3,
			ExamQuestionsNeeded: 3,
			Staff: []StaffResponse{
				{
					Login:          "jwang4",
					Email:          "josiah.wang@imperial.ac.uk",
					RoleInDepartment: "staff",
					LastName:       "Wang",
					FirstName:      "Josiah",
					Department:     "DoC",
					CID:            "01030000",
				},
			},
			Helpers: []HelperResponse{
				{
					Login:    "ad5518",
					Email:    "adam.dejl18@imperial.ac.uk",
					LastName: "Dejl",
					FirstName: "Adam",
					RoleInDepartment:    "Tutorial Helper",
					Department:       "DoC",
				},
			},
		}, nil
}

// Truncated for SEGP
func (m *MockAbcApiClient) GetEnrolledStudents2() ([]*EnrolledStudentResponse, error) {
	return []*EnrolledStudentResponse{
		{
			Login:     "jw5322",
			Email:     "james.watling22@imperial.ac.uk",
			FirstName: "James",
			LastName:  "Watling",
			Level:     3,
			Status:    "Normal",
			Cohort:    "c3",
		},
		{
			Login:     "ap3022",
			Email:     "alexey.popov22@imperial.ac.uk",
			FirstName: "Alexey",
			LastName:  "Popov",
			Level:     3,
			Status:    "Normal",
			Cohort:    "c3",
		},
		{
			Login:     "dw922",
			Email:     "daniel.wait22@imperial.ac.uk",
			FirstName: "Daniel",
			LastName:  "Wait",
			Level:     3,
			Status:    "Normal",
			Cohort:    "c3",
		},
		{
			Login:     "mp1822",
			Email:     "mann.patira22@imperial.ac.uk",
			FirstName: "Mann",
			LastName:  "Patira",
			Level:     3,
			Status:    "Normal",
			Cohort:    "c3",
		},
		{
			Login:     "as4522",
			Email:     "anshul.sendil22@imperial.ac.uk",
			FirstName: "Anshul",
			LastName:  "Sendil",
			Level:     3,
			Status:    "Normal",
			Cohort:    "c3",
		},
		}, nil
}
