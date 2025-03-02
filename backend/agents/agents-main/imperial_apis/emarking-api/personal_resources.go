package emarking_api

type DeliverableResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ExerciseResponse struct {
	Year                  string                `json:"year"`
	ModuleCode            string                `json:"module_code"`
	Title                 string                `json:"title"`
	Start                 string                `json:"start"`
	End                   string                `json:"end"`
	SubmissionType        string                `json:"submission_type"`
	MaximumMark           int                   `json:"maximum_mark"`
	ExpectedHours         int                   `json:"expected_hours"`
	MarksPublished        string                `json:"marks_published,omitempty"`
	MarksPublishedBy      string                `json:"marks_published_by,omitempty"`
	Spec                  string                `json:"spec,omitempty"`
	Weight                int                   `json:"weight"`
	MarksHiddenToStudents bool                  `json:"marks_hidden_to_students,omitempty"`
	Deliverables          []DeliverableResponse `json:"deliverables"`
	ExtendedEnd           string                `json:"extended_end,omitempty"`
	Mark                  MarkResponse          `json:"mark,omitempty"`
	Submissions           []SubmissionResponse  `json:"submissions"`
	Feedback              FeedbackResponse      `json:"feedback"`
}

type MarkResponse struct {
	StudentUsername string `json:"student_username"`
	Marker          string `json:"marker"`
	Mark            int    `json:"mark"`
	Timestamp       string `json:"timestamp"`
}

type SubmissionResponse struct {
	Username                 string `json:"username"`
	Timestamp                string `json:"timestamp"`
	TargetSubmissionFileName string `json:"target_submission_file_name"`
}

type FeedbackResponse struct {
	Timestamp       string `json:"timestamp"`
	StudentUsername string `json:"student_username"`
	Marker          string `json:"marker"`
	ModuleCode      string `json:"module_code"`
}

type MemberResponse struct {
	Username string `json:"username"`
	IsLeader bool   `json:"is_leader"`
}

type SubmissionGroupResponse struct {
	Year       string           `json:"year"`
	ModuleCode string           `json:"module_code"`
	Members    []MemberResponse `json:"members"`
}

func (m *MockEmarkingApiClient) GetExercises() []ExerciseResponse {
	return []ExerciseResponse{
		{
			Year:                  "2425",
			ModuleCode:            "60007",
			Title:                 "Practice coursework",
			Start:                 "2024-10-21T12:00:00+00:00",
			End:                   "2024-11-01T17:00:00+00:00",
			SubmissionType:        "group",
			MaximumMark:           20,
			ExpectedHours:         10,
			MarksPublished:        "2024-11-13T19:00:00+00:00",
			MarksPublishedBy:      "azalea",
			Spec:                  "https://scientia.doc.ic.ac.uk/api/2425/60007/exercises/1/spec",
			Weight:                50,
			MarksHiddenToStudents: true,
			Deliverables: []DeliverableResponse{
				{
					Name: "practice.zip",
					Type: "file",
				},
			},
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker:          "azalea",
				Mark:            45,
				Timestamp:       "2024-11-13T19:00:00+00:00",
			},
			ExtendedEnd: "2024-11-06T17:00:00+00:00",
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2024-11-06T15:41:29+00:00",
					TargetSubmissionFileName: "practice.zip",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2024-11-13T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "azalea",
				ModuleCode:      "60007",
			},
		},
		{
			Year:                  "2425",
			ModuleCode:            "60007",
			Title:                 "CW: Theory Coursework",
			Start:                 "2024-11-01T12:00:00+00:00",
			End:                   "2024-11-22T17:00:00+00:00",
			SubmissionType:        "group",
			MaximumMark:           20,
			ExpectedHours:         10,
			MarksPublished:        "2024-11-25T19:00:00+00:00",
			MarksPublishedBy:      "azalea",
			Spec:                  "https://scientia.doc.ic.ac.uk/api/2425/60007/exercises/2/spec",
			Weight:                50,
			MarksHiddenToStudents: true,
			Deliverables: []DeliverableResponse{
				{
					Name: "answers.pdf",
					Type: "file",
				},
			},
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker:          "azalea",
				Mark:            52,
				Timestamp:       "2024-11-25T19:00:00+00:00",
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2024-11-22T15:41:29+00:00",
					TargetSubmissionFileName: "practice.zip",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2024-11-25T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "azalea",
				ModuleCode:      "60007",
			},
		},
		{
			Year:                  "2425",
			ModuleCode:            "60012",
			Title:                 "Decisions Trees",
			Start:                 "2024-10-14T12:00:00+00:00",
			End:                   "2024-11-01T19:00:00+00:00",
			SubmissionType:        "group",
			MaximumMark:           100,
			ExpectedHours:         12,
			MarksPublished:        "2024-11-13T19:00:00+00:00",
			MarksPublishedBy:      "jwang4",
			Spec:                  "https://scientia.doc.ic.ac.uk/api/2425/60012/exercises/1/spec",
			Weight:                40,
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
				Marker:          "jwang4",
				Mark:            80,
				Timestamp:       "2024-11-13T19:00:00+00:00",
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2024-10-22T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
				{
					Username:                 "as4522",
					Timestamp:                "2024-10-22T14:34:33+00:00",
					TargetSubmissionFileName: "source.zip",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2024-11-13T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "jwang4",
				ModuleCode:      "60012",
			},
		},
		{
			Year:             "2425",
			ModuleCode:       "60012",
			Title:            "Neural Networks",
			Start:            "2024-11-04T12:00:00+00:00",
			End:              "2024-11-22T19:00:00+00:00",
			SubmissionType:   "group",
			MaximumMark:      100,
			ExpectedHours:    12,
			MarksPublished:   "2024-11-24T19:00:00+00:00",
			MarksPublishedBy: "jwang4",
			Spec:             "https://scientia.doc.ic.ac.uk/api/2425/60012/exercises/2/spec",
			Weight:           60,
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
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker:          "jwang4",
				Mark:            63,
				Timestamp:       "2024-11-24T19:00:00+00:00",
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2024-11-20T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
				{
					Username:                 "as4522",
					Timestamp:                "2024-11-20T14:34:33+00:00",
					TargetSubmissionFileName: "source.zip",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2024-11-24T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "jwang4",
				ModuleCode:      "60012",
			},
		},
		{
			Year:             "2425",
			ModuleCode:       "60001",
			Title:            "The Coursework",
			Start:            "2024-10-29T12:00:00+00:00",
			End:              "2024-11-19T19:00:00+00:00",
			SubmissionType:   "individual",
			MaximumMark:      100,
			ExpectedHours:    8,
			MarksPublished:   "2024-11-24T19:00:00+00:00",
			MarksPublishedBy: "phjk",
			Spec:             "https://scientia.doc.ic.ac.uk/api/2425/60001",
			Weight:           100,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker:          "phjk",
				Mark:            50,
				Timestamp:       "2024-11-21T19:00:00+00:00",
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2024-11-18T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2024-11-21T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "phjk",
				ModuleCode:      "60001",
			},
		},
		{
			Year:             "2425",
			ModuleCode:       "70015",
			Title:            "The Coursework",
			Start:            "2024-11-04T12:00:00+00:00",
			End:              "2024-12-02T17:00:00+00:00",
			SubmissionType:   "group",
			MaximumMark:      100,
			ExpectedHours:    30,
			MarksPublished:   "2024-12-05T19:00:00+00:00",
			MarksPublishedBy: "rac101",
			Spec:             "https://scientia.doc.ic.ac.uk/api/2425/70015",
			Weight:           100,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker:          "rac101",
				Mark:            43,
				Timestamp:       "2024-12-05T19:00:00+00:00",
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2024-11-30T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2024-12-05T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "rac101",
				ModuleCode:      "70015",
			},
		},
		{
			Year:             "2425",
			ModuleCode:       "60006",
			Title:            "Coursework 1",
			Start:            "2025-01-16T12:00:00+00:00",
			End:              "2025-01-30T19:00:00+00:00",
			SubmissionType:   "individual",
			MaximumMark:      100,
			ExpectedHours:    4,
			MarksPublished:   "2025-02-05T19:00:00+00:00",
			MarksPublishedBy: "wbai",
			Spec:             "https://scientia.doc.ic.ac.uk/api/2425/60006",
			Weight:           40,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker:          "wbai",
				Mark:            58,
				Timestamp:       "2025-02-05T19:00:00+00:00",
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2025-01-29T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2025-02-05T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "wbai",
				ModuleCode:      "60006",
			},
		},
		{
			Year:           "2425",
			ModuleCode:     "60006",
			Title:          "Coursework 2",
			Start:          "2025-02-06T12:00:00+00:00",
			End:            "2025-02-21T19:00:00+00:00",
			SubmissionType: "individual",
			MaximumMark:    100,
			ExpectedHours:  6,
			Spec:           "https://scientia.doc.ic.ac.uk/api/2425/60006",
			Weight:         60,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2025-02-21T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
			},
		},
		{
			Year:             "2425",
			ModuleCode:       "60005",
			Title:            "Illumination and Shading",
			Start:            "2025-02-03T12:00:00+00:00",
			End:              "2025-02-14T19:00:00+00:00",
			SubmissionType:   "individual",
			MaximumMark:      100,
			ExpectedHours:    2,
			MarksPublished:   "2025-02-18T19:00:00+00:00",
			MarksPublishedBy: "abgh",
			Spec:             "https://scientia.doc.ic.ac.uk/api/2425/60005",
			Weight:           15,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker:          "abgh",
				Mark:            78,
				Timestamp:       "2025-02-18T19:00:00+00:00",
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2025-02-12T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2025-02-18T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "abgh",
				ModuleCode:      "60005",
			},
		},
		{
			Year:             "2425",
			ModuleCode:       "60005",
			Title:            "Texture",
			Start:            "2025-02-14T12:00:00+00:00",
			End:              "2025-02-21T19:00:00+00:00",
			SubmissionType:   "individual",
			MaximumMark:      100,
			ExpectedHours:    2,
			MarksPublished:   "2025-02-28T19:00:00+00:00",
			MarksPublishedBy: "abgh",
			Spec:             "https://scientia.doc.ic.ac.uk/api/2425/60005",
			Weight:           10,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Mark: MarkResponse{
				StudentUsername: "as4522",
				Marker:          "abgh",
				Mark:            64,
				Timestamp:       "2025-02-28T19:00:00+00:00",
			},
			Submissions: []SubmissionResponse{
				{
					Username:                 "as4522",
					Timestamp:                "2025-02-18T14:34:33+00:00",
					TargetSubmissionFileName: "report.pdf",
				},
			},
			Feedback: FeedbackResponse{
				Timestamp:       "2025-02-28T19:00:00+00:00",
				StudentUsername: "as4522",
				Marker:          "abgh",
				ModuleCode:      "60005",
			},
		},
		{
			Year:           "2425",
			ModuleCode:     "60005",
			Title:          "Raytracing",
			Start:          "2025-02-14T12:00:00+00:00",
			End:            "2025-03-07T19:00:00+00:00",
			SubmissionType: "individual",
			MaximumMark:    100,
			ExpectedHours:  6,
			Spec:           "https://scientia.doc.ic.ac.uk/api/2425/60005",
			Weight:         75,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Submissions: []SubmissionResponse{},
		},
		{
			Year:           "2425",
			ModuleCode:     "60008",
			Title:          "Custom Computing Assessed Coursework",
			Start:          "2025-01-27T12:00:00+00:00",
			End:            "2025-02-24T19:00:00+00:00",
			SubmissionType: "individual",
			MaximumMark:    100,
			ExpectedHours:  10,
			Spec:           "https://scientia.doc.ic.ac.uk/api/2425/60008",
			Weight:         100,
			Deliverables: []DeliverableResponse{
				{
					Name: "report.pdf",
					Type: "file",
				},
			},
			Submissions: []SubmissionResponse{},
		},
	}
}

func (m *MockEmarkingApiClient) GetFeedback() []FeedbackResponse {
	return []FeedbackResponse{
		{
			Timestamp:       "2024-11-13T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "azalea",
			ModuleCode:      "60007",
		},
		{
			Timestamp:       "2024-11-25T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "azalea",
			ModuleCode:      "60007",
		},
		{
			Timestamp:       "2024-11-13T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "jwang4",
			ModuleCode:      "60012",
		},
		{
			Timestamp:       "2024-11-24T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "jwang4",
			ModuleCode:      "60012",
		},
		{
			Timestamp:       "2024-11-21T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "phjk",
			ModuleCode:      "60001",
		},
		{
			Timestamp:       "2024-12-05T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "rac101",
			ModuleCode:      "70015",
		},
		{
			Timestamp:       "2025-02-05T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "wbai",
			ModuleCode:      "60006",
		},
		{
			Timestamp:       "2025-02-18T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "abgh",
			ModuleCode:      "60005",
		},
		{
			Timestamp:       "2025-02-28T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "abgh",
			ModuleCode:      "60005",
		},
	}
}

func (m *MockEmarkingApiClient) GetExerciseSummary() ExerciseResponse {
	return ExerciseResponse{
		Year:                  "2425",
		ModuleCode:            "60012",
		Title:                 "Decisions Trees",
		Start:                 "2024-10-14T12:00:00+00:00",
		End:                   "2024-11-01T19:00:00+00:00",
		SubmissionType:        "group",
		MaximumMark:           100,
		ExpectedHours:         12,
		MarksPublished:        "2024-11-13T19:00:00+00:00",
		MarksPublishedBy:      "jwang4",
		Spec:                  "https://scientia.doc.ic.ac.uk/api/2425/60012/exercises/1/spec",
		Weight:                40,
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
			Marker:          "jwang4",
			Mark:            80,
			Timestamp:       "2024-11-13T19:00:00+00:00",
		},
		Submissions: []SubmissionResponse{
			{
				Username:                 "as4522",
				Timestamp:                "2024-10-22T14:34:33+00:00",
				TargetSubmissionFileName: "report.pdf",
			},
			{
				Username:                 "as4522",
				Timestamp:                "2024-10-22T14:34:33+00:00",
				TargetSubmissionFileName: "source.zip",
			},
		},
		Feedback: FeedbackResponse{
			Timestamp:       "2024-11-13T19:00:00+00:00",
			StudentUsername: "as4522",
			Marker:          "jwang4",
			ModuleCode:      "60012",
		},
	}
}

func (m *MockEmarkingApiClient) GetSubmissionGroup() []SubmissionGroupResponse {
	return []SubmissionGroupResponse{
		{
			Year:       "2425",
			ModuleCode: "60012",
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
		},
		{
			Year:       "2425",
			ModuleCode: "60007",
			Members: []MemberResponse{
				{
					Username: "as4522",
					IsLeader: false,
				},
				{
					Username: "dbs21",
					IsLeader: true,
				},
				{
					Username: "th1522",
					IsLeader: false,
				},
			},
		},
		{
			Year:       "2425",
			ModuleCode: "70015",
			Members: []MemberResponse{
				{
					Username: "as4522",
					IsLeader: false,
				},
				{
					Username: "th1522",
					IsLeader: true,
				},
				{
					Username: "eh1322",
					IsLeader: false,
				},
			},
		},
	}
}
