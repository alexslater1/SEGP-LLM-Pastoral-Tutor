package emarking_api



type DeliverableResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ExerciseResponse struct {
	Year                      string    `json:"year"`
	ModuleCode                string    `json:"module_code"`
	Title                     string    `json:"title"`
	Start                     string `json:"start"`
	End                       string `json:"end"`
	SubmissionType            string    `json:"submission_type"`
	MaximumMark               int       `json:"maximum_mark"`
	ExpectedHours             int       `json:"expected_hours"`
	MarksPublished            string `json:"marks_published,omitempty"`
	MarksPublishedBy          string    `json:"marks_published_by,omitempty"`
	Spec                      string `json:"spec,omitempty"`
	Weight                    int       `json:"weight"`
	MarksHiddenToStudents     bool      `json:"marks_hidden_to_students,omitempty"`
	Deliverables              []DeliverableResponse `json:"deliverables"`
	Mark                      MarkResponse      `json:"mark,omitempty"`
	Submissions               []SubmissionResponse `json:"submissions,omitempty"`
	Feedback                  FeedbackResponse `json:"feedback,omitempty"`
}

type MarkResponse struct {
	StudentUsername string    `json:"student_username"`
	Marker        string    `json:"marker"`
	Mark          int       `json:"mark"`
	Timestamp     string `json:"timestamp"`
}



type SubmissionResponse struct {
	Username                  string    `json:"username"`
	Timestamp                 string `json:"timestamp"`
	TargetSubmissionFileName  string    `json:"target_submission_file_name"`
}

type FeedbackResponse struct {
	Timestamp       string    `json:"timestamp"`
	StudentUsername string    `json:"student_username"`
	Marker          string    `json:"marker"`
	ModuleCode	    string    `json:"module_code"`
}

type MemberResponse struct {
	Username string `json:"username"`
	IsLeader bool `json:"is_leader"`
}

type SubmissionGroupResponse struct {
	Year                      string    `json:"year"`
	ModuleCode                string    `json:"module_code"`
	Members				  []MemberResponse  `json:"members"`
}

func (m *MockEmarkingApiClient) GetExercises() []ExerciseResponse {
	return []ExerciseResponse{
		{
			Year:              "2425",
			ModuleCode:        "60007",
			Title:             "Practice coursework",
			Start:             "2024-10-21T12:00:00+00:00",
			End:               "2024-11-01T17:00:00+00:00",
			SubmissionType:    "group",
			MaximumMark:       20,
			ExpectedHours:     10,
			Spec: "https://scientia.doc.ic.ac.uk/api/2425/60007/exercises/1/spec",
			Weight:            50,
			MarksHiddenToStudents: true,
			Deliverables: []DeliverableResponse{
				{
					Name: "practice.zip",
					Type: "file",
				},
			},
			Submissions:           []SubmissionResponse{},
		},
		{
			Year:              "2425",
			ModuleCode:        "60007",
			Title:             "CW: Theory Coursework",
			Start:             "2024-11-01T12:00:00+00:00",
			End:               "2024-11-22T17:00:00+00:00",
			SubmissionType:    "group",
			MaximumMark:       20,
			ExpectedHours:     10,
			Spec: "https://scientia.doc.ic.ac.uk/api/2425/60007/exercises/2/spec",
			Weight:            50,
			MarksHiddenToStudents: true,
			Deliverables: []DeliverableResponse{
				{
					Name: "answers.pdf",
					Type: "file",
				},
			},
			Submissions:           []SubmissionResponse{},
		},
		{
			Year:              "2425",
			ModuleCode:        "60012",
			Title:             "Decisions Trees",
			Start:             "2024-10-14T12:00:00+00:00",
			End:               "2024-11-01T19:00:00+00:00",
			SubmissionType:    "group",
			MaximumMark:       100,
			ExpectedHours:     12,
			MarksPublished:    "2024-11-13T19:00:00+00:00",
			MarksPublishedBy:  "jwang4",
			Spec: 			   "https://scientia.doc.ic.ac.uk/api/2425/60012/exercises/1/spec",
			Weight:            40,
			MarksHiddenToStudents: false,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
				{
					Name: "source.zip",
					Type: "file",
				},
			},
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker: "jwang4",
				Mark: 80,
				Timestamp: "2024-11-13T19:00:00+00:00",
			},
			Submissions:           []SubmissionResponse{
				{
					Username: "as4522",
					Timestamp: "2024-10-22T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
				{
					Username: "as4522",
					Timestamp: "2024-10-22T14:34:33+00:00",
					TargetSubmissionFileName: "source.zip",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp: "2024-11-13T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker: "jwang4",
				ModuleCode: "60012",
			},
		},
		{
			Year:              "2425",
			ModuleCode:        "60012",
			Title:             "Neural Networks",
			Start:             "2024-11-04T12:00:00+00:00",
			End:               "2024-11-22T19:00:00+00:00",
			SubmissionType:    "group",
			MaximumMark:       100,
			ExpectedHours:     12,
			Spec: "https://scientia.doc.ic.ac.uk/api/2425/60012/exercises/2/spec",
			Weight:            60,
			Deliverables: []DeliverableResponse{
				{
					Name: "gitlab hash (via labts)",
					Type: "hash",
				},
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Submissions:           []SubmissionResponse{},
		},
		{
			Year:              "2425",
			ModuleCode:        "60001",
			Title:             "The Coursework",
			Start:             "2024-10-29T12:00:00+00:00",
			End:               "2024-11-19T19:00:00+00:00",
			SubmissionType:    "individual",
			MaximumMark:       100,
			ExpectedHours:     8,
			Spec: "https://scientia.doc.ic.ac.uk/api/2425/60001",
			Weight:            100,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Submissions:           []SubmissionResponse{},
		},
		{
			Year:              "2425",
			ModuleCode:        "70015",
			Title:             "The Coursework",
			Start:             "2024-11-04T12:00:00+00:00",
			End:               "2024-12-02T17:00:00+00:00",
			SubmissionType:    "group",
			MaximumMark:       100,
			ExpectedHours:     30,
			Spec: "https://scientia.doc.ic.ac.uk/api/2425/70015",
			Weight:            100,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Submissions:           []SubmissionResponse{},
		},
	}
}

func (m *MockEmarkingApiClient) GetFeedback() []FeedbackResponse {
	return []FeedbackResponse{
		{
			Timestamp: "2024-11-13T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker: "jwang4",
			ModuleCode: "60012",
		},
	}
}

func (m *MockEmarkingApiClient) GetExerciseSummary() ExerciseResponse {
	return ExerciseResponse{
		Year:              "2425",
		ModuleCode:        "60012",
		Title:             "Decisions Trees",
		Start:             "2024-10-14T12:00:00+00:00",
		End:               "2024-11-01T19:00:00+00:00",
		SubmissionType:    "group",
		MaximumMark:       100,
		ExpectedHours:     12,
		MarksPublished:    "2024-11-13T19:00:00+00:00",
		MarksPublishedBy:  "jwang4",
		Spec: "https://scientia.doc.ic.ac.uk/api/2425/60012/exercises/1/spec",
		Weight:            40,
		MarksHiddenToStudents: false,
		Deliverables: []DeliverableResponse{
			{
				Name: "report.pdf",
				Type: "file",
			},
			{
				Name: "source.zip",
				Type: "file",
			},
		},
		Mark: MarkResponse{
			StudentUsername: "as4522",
			Marker: "jwang4",
			Mark: 80,
			Timestamp: "2024-11-13T19:00:00+00:00",
		},
		Submissions:           []SubmissionResponse{
			{
				Username: "as4522",
				Timestamp: "2024-10-22T14:34:33+00:00",
				TargetSubmissionFileName: "report.pdf",
			},
			{
				Username: "as4522",
				Timestamp: "2024-10-22T14:34:33+00:00",
				TargetSubmissionFileName: "source.zip",
			},
		},
		Feedback: FeedbackResponse{
			Timestamp: "2024-11-13T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker: "jwang4",
			ModuleCode: "60012",
		},
	}
}

func (m *MockEmarkingApiClient) GetSubmissionGroup() SubmissionGroupResponse {
	return SubmissionGroupResponse{
		Year:              "2425",
		ModuleCode:        "60012",
		Members: []MemberResponse{
			{
				Username: "as4522",
				IsLeader: true,
			},
			{
				Username: "dbs21",
				IsLeader: false,
			},
			{
				Username: "eh1322",
				IsLeader: false,
			},
		},
	}
}



