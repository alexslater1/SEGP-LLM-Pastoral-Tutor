package abc_api

type Supervisor struct {
    Login        string `json:"login"`
    FirstName    string `json:"firstname"`
    LastName     string `json:"lastname"`
    EmailAddress string `json:"email_address"`
    Role         string `json:"role"`
}

type PhdStudentsForSupervisorResponse struct {
    Login        string       `json:"login"`
    FirstName    string       `json:"firstname"`
    LastName     string       `json:"lastname"`
    EmailAddress string       `json:"email_address"`
    StudentStatus string      `json:"studentstatus"`
    StartDate    string       `json:"start_date"`
    EndDate      string       `json:"end_date"`
    CID          string       `json:"cid"`
    Supervisors  []Supervisor `json:"supervisors"`
}

type PhdStudentInformationResponse struct {
    Login         string       `json:"login"`
    FirstName     string       `json:"firstname"`
    LastName      string       `json:"lastname"`
    EmailAddress  string       `json:"email_address"`
    StudentStatus string       `json:"studentstatus"`
    StartDate     string       `json:"start_date"`
    EndDate       string       `json:"end_date"`
    CID           string       `json:"cid"`
    Supervisors   []Supervisor `json:"supervisors"`
}

func (m *MockAbcApiClient) GetPhdStudentsForSupervisor() PhdStudentsForSupervisorResponse {
    return PhdStudentsForSupervisorResponse{
        Login:        "azhang435",
        FirstName:    "Aisha",
        LastName:     "Zhang",
        EmailAddress: "azhang435@ic.ac.uk",
        StudentStatus: "Active",
        StartDate:    "2022-10-01",
        EndDate:      "2025-09-30",
        CID:          "23857149",
        Supervisors: []Supervisor{
            {
                Login:        "rgupta",
                FirstName:    "Raj",
                LastName:     "Gupta",
                EmailAddress: "rgupta@ic.ac.uk",
                Role:         "Primary Supervisor",
            },
            {
                Login:        "ekowalski",
                FirstName:    "Elena",
                LastName:     "Kowalski",
                EmailAddress: "ekowalski@ic.ac.uk",
                Role:         "Secondary Supervisor",
            },
        },
    }
}

func (m *MockAbcApiClient) GetPhdStudentInformation() PhdStudentInformationResponse {
    return PhdStudentInformationResponse{
        Login:         "msmith123",
        FirstName:     "Maria",
        LastName:      "Smith",
        EmailAddress:  "msmith123@ic.ac.uk",
        StudentStatus: "Active",
        StartDate:     "2021-09-01",
        EndDate:       "2024-08-31",
        CID:           "21934567",
        Supervisors: []Supervisor{
            {
                Login:        "jchen",
                FirstName:    "James",
                LastName:     "Chen",
                EmailAddress: "jchen@ic.ac.uk",
                Role:         "Primary Supervisor",
            },
            {
                Login:        "lpatel",
                FirstName:    "Leela",
                LastName:     "Patel",
                EmailAddress: "lpatel@ic.ac.uk",
                Role:         "Secondary Supervisor",
            },
        },
    }
}
