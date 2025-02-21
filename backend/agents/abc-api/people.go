package abc_api

/*
[
  {
    "login": "rbc",
    "year": "2223",
    "email": "rbc@ic.ac.uk",
    "firstname": "Rob",
    "lastname": "Chatley",
    "salutation": "Mr",
    "role_in_department": "staff",
    "roles_in_department": [
      "staff",
      "2nd Year Undergraduate Coordinator"
    ],
    "has_extension_clearance": true,
    "modules": [
      {
        "code": "50002",
        "title": "Software Engineering Design",
        "terms": [
          1,
          2
        ],
        "roles": [
          "Lecturer"
        ]
      }
    ]
  }
]
*/

type ModuleResponse struct {
	Level int    `json:"level"`
	Code  string `json:"code"`
	Title string `json:"title"`
	Terms []int  `json:"terms"`
}

type ModuleHelpedResponse struct {
	Code  string   `json:"code"`
	Title string   `json:"title"`
	Roles []string `json:"roles"`
	Terms []int    `json:"terms"`
}

type PersonalTutorResponse struct {
	Login     string `json:"login"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type StudentResponse struct {
	Login            string                 `json:"login"`
	Year             string                 `json:"year"`
	Email            string                 `json:"email"`
	Firstname        string                 `json:"firstname"`
	Lastname         string                 `json:"lastname"`
	Salutation       string                 `json:"salutation"`
	Cohort           string                 `json:"cohort"`
	DegreeYear       string                 `json:"degree_year"`
	RoleInDepartment string                 `json:"role_in_department"`
	Modules          []ModuleResponse       `json:"modules"`
	ModulesHelped    []ModuleHelpedResponse `json:"modules_helped"`
	PersonalTutor    PersonalTutorResponse  `json:"personal_tutor"`
}

type StaffPeopleResponse struct {
	Login                 string                 `json:"login"`
	Year                  string                 `json:"year"`
	Email                 string                 `json:"email"`
	Firstname             string                 `json:"firstname"`
	Lastname              string                 `json:"lastname"`
	Salutation            string                 `json:"salutation"`
	RoleInDepartment      string                 `json:"role_in_department"`
	RolesInDepartment     []string               `json:"roles_in_department"`
	HasExtensionClearance bool                   `json:"has_extension_clearance"`
	Modules               []ModuleHelpedResponse `json:"modules"`
}

func (m *MockAbcApiClient) GetStudents() ([]StudentResponse, error) {
	return []StudentResponse{
		{
			Login:            "as4522",
			Year:             "2324",
			Email:            "as4522@ic.ac.uk",
			Firstname:        "Anshul",
			Lastname:         "Sendil",
			Salutation:       "Mr",
			Cohort:           "c3",
			DegreeYear:       "meng3",
			RoleInDepartment: "student",
			Modules: []ModuleResponse{
				{
					Level: 3,
					Code:  "50009",
					Title: "Computer Vision",
					Terms: []int{1, 2},
				},
				{
					Level: 3,
					Code:  "50010",
					Title: "Network Security",
					Terms: []int{1},
				},
				{
					Level: 3,
					Code:  "50002",
					Title: "Software Engineering Group Project",
					Terms: []int{2},
				},
				{
					Level: 3,
					Code:  "50011",
					Title: "Japanese Language and Culture",
					Terms: []int{1, 2},
				},
			},
			ModulesHelped: []ModuleHelpedResponse{}, // No modules helped as he's struggling with his own studies
			PersonalTutor: PersonalTutorResponse{
				Login:     "dsmith",
				Firstname: "David",
				Lastname:  "Smith",
			},
		},
	}, nil
}

func (m *MockAbcApiClient) GetStaff() ([]StaffPeopleResponse, error) {
	return []StaffPeopleResponse{
		{
			Login:                 "rbc",
			Year:                  "2223",
			Email:                 "rbc@ic.ac.uk",
			Firstname:             "Rob",
			Lastname:              "Chatley",
			Salutation:            "Mr",
			RoleInDepartment:      "staff",
			RolesInDepartment:     []string{"staff", "2nd Year Undergraduate Coordinator"},
			HasExtensionClearance: true,
			Modules: []ModuleHelpedResponse{
				{
					Code:  "50002",
					Title: "Software Engineering Group Project",
					Roles: []string{"Lecturer"},
					Terms: []int{1, 2},
				},
			},
		},
	}, nil
}

func (m *MockAbcApiClient) GetAllStaffList() ([]StaffResponse, error) {
	return []StaffResponse{
		{
			Login:            "rbc",
			Email:            "rbc@ic.ac.uk",
			FirstName:        "Rob",
			LastName:         "Chatley",
			Salutation:       "Mr",
			RoleInDepartment: "staff",
		},
	}, nil
}
