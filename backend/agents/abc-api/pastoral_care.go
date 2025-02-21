package abc_api

type Person struct {
	Login     string `json:"login"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
}

type TutorialGroup struct {
	Number  int      `json:"number"`
	Type    string   `json:"type"`
	Tutor   Person   `json:"tutor"`
	UTA     Person   `json:"uta"`
	Members []Person `json:"members"`
}

type TutorialGroupsResponse []TutorialGroup

type Tutee struct {
	Login     string `json:"login"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Cohort    string `json:"cohort"`
}

type TutorTuteeRelation struct {
	Tutor Person `json:"tutor"`
	Tutee Tutee  `json:"tutee"`
}

type PersonalTuteesForTutorResponse []TutorTuteeRelation

func (m *MockAbcApiClient) GetTutorialGroups() (*TutorialGroupsResponse, error) {
	return &TutorialGroupsResponse{
		Number: 3,
		Type:   "MEng Computing",
		Tutor: Person{
			Login:     "tutor123",
			Lastname:  "Smith",
			Firstname: "John",
		},
		UTA: Person{
			Login:     "uta456",
			Lastname:  "Johnson",
			Firstname: "Sarah",
		},
		Members: []Person{
			{
				Login:     "as4522",
				Lastname:  "Kumar",
				Firstname: "Anshul",
			},
		},
	}, nil
}

func (m *MockAbcApiClient) GetPersonalTuteesForTutor() (*PersonalTuteesForTutorResponse, error) {
	return &PersonalTuteesForTutorResponse{
		{
			Tutor: Person{
				Login:     "tutor123",
				Lastname:  "Smith",
				Firstname: "John",
			},
			Tutee: Tutee{
				Login:     "as4522",
				Lastname:  "Kumar",
				Firstname: "Anshul",
				Cohort:    "2023/24",
			},
		},
		{
			Tutor: Person{
				Login:     "tutor123",
				Lastname:  "Smith",
				Firstname: "John",
			},
			Tutee: Tutee{
				Login:     "jane22",
				Lastname:  "Doe",
				Firstname: "Jane",
				Cohort:    "2023/24",
			},
		},
	}, nil
}

