package abc_api

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

func (m *MockAbcApiClient) GetStudents() StudentResponse {
	return StudentResponse{
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
	}
}

func (m *MockAbcApiClient) GetStaff() StaffPeopleResponse {
	return StaffPeopleResponse{
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
	}
}

func (m *MockAbcApiClient) GetAllStaffList() []StaffPeopleResponse {
	return []StaffPeopleResponse{
		{
			Login:            "rbc",
			Email:            "rbc@ic.ac.uk",
			Firstname:        "Rob",
			Lastname:         "Chatley",
			Salutation:       "Mr",
			RoleInDepartment: "staff",
		},
	}
}

type ProfileImageResponse struct {
	Image string `json:"image"`
}

func (m *MockAbcApiClient) GetProfileImage() (*ProfileImageResponse, error) {
	return &ProfileImageResponse{
		Image: "https://imperialimages.com/as4522/profile.jpg",
	}, nil
}

type TotalEnrolledStudentsResponse struct {
	Total int `json:"total"`
}

func (m *MockAbcApiClient) GetTotalEnrolledStudents() TotalEnrolledStudentsResponse {
	return TotalEnrolledStudentsResponse{
		Total: 100,
	}
}

type EnrolledStudentsResponse struct {
	Modules []ModuleResponse2 `json:"modules"`
}

type ModuleResponse2 struct {
	ModuleCode string   `json:"module_code"`
	Students   []string `json:"students"`
}

func (m *MockAbcApiClient) GetModulesEnrolledStudents() EnrolledStudentsResponse {
	return EnrolledStudentsResponse{
		Modules: []ModuleResponse2{
			{ModuleCode: "50002", Students: []string{"anshul"}},
		},
	}
}

type IdentityResponse struct {
	IsStaff bool   `json:"is_staff"`
	Login   string `json:"login"`
	Email   string `json:"email"`
}

func (m *MockAbcApiClient) GetIdentity() IdentityResponse {
	return IdentityResponse{
		IsStaff: false,
		Login:   "anshul",
		Email:   "as522@ic.ac.uk",
	}
}

type StudentDetailsResponse struct {
	Login            string `json:"login"`
	Email            string `json:"email"`
	Lastname         string `json:"lastname"`
	Firstname        string `json:"firstname"`
	Salutation       string `json:"salutation"`
	Year             string `json:"year"`
	RoleInDepartment string `json:"role_in_department"`
	Cohort           string `json:"cohort"`
	CID              string `json:"cid"`
	ExamClass        string `json:"exam_class"`
	DegreeCode       string `json:"degree_code"`
	DegreeYear       string `json:"degree_year"`
	StudentStatus    string `json:"student_status"`
	EntryYear        int    `json:"entry_year"`
	FeeStatus        string `json:"fee_status"`
	PersonalTutor    struct {
		Login     string `json:"login"`
		Lastname  string `json:"lastname"`
		Firstname string `json:"firstname"`
	} `json:"personal_tutor"`
}

func (m *MockAbcApiClient) GetStudentDetails() StudentDetailsResponse {
	return StudentDetailsResponse{
		Login:            "anshul",
		Email:            "anshul@university.edu",
		Lastname:         "Sharma",
		Firstname:        "Anshul",
		Salutation:       "Mr.",
		Year:             "3",
		RoleInDepartment: "Student",
		Cohort:           "2021",
		CID:              "A12345678",
		ExamClass:        "2025",
		DegreeCode:       "MEng Computing",
		DegreeYear:       "2024",
		StudentStatus:    "Active",
		EntryYear:        2022,
		FeeStatus:        "Home",
		PersonalTutor: struct {
			Login     string `json:"login"`
			Lastname  string `json:"lastname"`
			Firstname string `json:"firstname"`
		}{
			Login:     "ad321",
			Lastname:  "Donaldson",
			Firstname: "Alistair",
		},
	}
}

func (m *MockAbcApiClient) GetAllStudents() []StudentDetailsResponse {
	studentDetails := m.GetStudentDetails()
	return []StudentDetailsResponse{studentDetails}
}
