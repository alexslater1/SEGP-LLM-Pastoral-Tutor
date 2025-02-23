package abc_api

type PersonResponse struct {
	Login     string `json:"login"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
}

type TutorialGroupResponse struct {
	Number  int              `json:"number"`
	Type    string           `json:"type"`
	Tutor   PersonResponse   `json:"tutor"`
	UTA     PersonResponse   `json:"uta"`
	Members []PersonResponse `json:"members"`
}

type TutorialGroupsResponse []TutorialGroupResponse

type TuteeResponse struct {
	Login     string `json:"login"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Cohort    string `json:"cohort"`
}

type TutorTuteeRelation struct {
	Tutor PersonResponse `json:"tutor"`
	Tutee TuteeResponse  `json:"tutee"`
}

type PersonalTuteesForTutorResponse []TutorTuteeRelation

func (m *MockAbcApiClient) GetTutorialGroups() TutorialGroupsResponse {
	return TutorialGroupsResponse{
		{
			Number: 3,
			Type:   "MEng Computing",
			Tutor: PersonResponse{
				Login:     "ad321",
				Lastname:  "Donaldson",
				Firstname: "Alistair",
			},
			UTA: PersonResponse{
				Login:     "yw2023",
				Lastname:  "Wong",
				Firstname: "Yuki",
			},
			Members: []PersonResponse{
				{
					Login:     "as4522",
					Lastname:  "Kumar",
					Firstname: "Anshul",
				},
			},
		},
	}
}

func (m *MockAbcApiClient) GetPersonalTuteesForTutor() PersonalTuteesForTutorResponse {
	return PersonalTuteesForTutorResponse{
		{
			Tutor: PersonResponse{
				Login:     "ad321",
				Lastname:  "Donaldson",
				Firstname: "Alistair",
			},
			Tutee: TuteeResponse{
				Login:     "as4522",
				Lastname:  "Kumar",
				Firstname: "Anshul",
				Cohort:    "2023/24",
			},
		},
		{
			Tutor: PersonResponse{
				Login:     "ad321",
				Lastname:  "Donaldson",
				Firstname: "Alistair",
			},
			Tutee: TuteeResponse{
				Login:     "zl214",
				Lastname:  "Liu",
				Firstname: "Zhang",
				Cohort:    "2023/24",
			},
		},
	}
}
