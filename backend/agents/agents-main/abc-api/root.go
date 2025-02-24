package abc_api

import (
	"time"
)

type YearsResponse struct {
	Years []int `json:"years"`
}

func (m *MockAbcApiClient) GetYears() YearsResponse {
	return YearsResponse{
		Years: []int{2018, 2019, 2020, 2021, 2022},
	}
}

type CohortResponse struct {
	DegreeYear  int    `json:"degree_year"`
	Degree      string `json:"degree"`
	DegreeShort string `json:"degree_short"`
}

func (m *MockAbcApiClient) GetCohorts() []CohortResponse {
	return []CohortResponse{
		{
			DegreeYear:  1,
			Degree:      "Bachelor of Science in Computer Science",
			DegreeShort: "BSc CS",
		},
		{
			DegreeYear:  1,
			Degree:      "Master of Computing",
			DegreeShort: "c1",
		},
		{
			DegreeYear:  2,
			Degree:      "Master of Computing",
			DegreeShort: "c2",
		},
		{
			DegreeYear:  3,
			Degree:      "Master of Computing",
			DegreeShort: "c3",
		},
		{
			DegreeYear:  4,
			Degree:      "Master of Computing",
			DegreeShort: "c4",
		},
		{
			DegreeYear:  1,
			Degree:      "Bachelors of Computing",
			DegreeShort: "c1",
		},
		{
			DegreeYear:  2,
			Degree:      "Bachelors of Computing",
			DegreeShort: "c2",
		},
		{
			DegreeYear:  3,
			Degree:      "Bachelors of Computing",
			DegreeShort: "c3",
		},
		{
			DegreeYear:  1,
			Degree:      "Joint Maths and Computing",
			DegreeShort: "j1",
		},
		{
			DegreeYear:  2,
			Degree:      "Joint Maths and Computing",
			DegreeShort: "j2",
		},
		{
			DegreeYear:  3,
			Degree:      "Joint Maths and Computing",
			DegreeShort: "j3",
		},
		{
			DegreeYear:  4,
			Degree:      "Joint Maths and Computing",
			DegreeShort: "j4",
		},
	}
}

type AcademicPeriodResponse struct {
	Name  string    `json:"name"`
	Weeks int       `json:"weeks"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (m *MockAbcApiClient) GetAcademicPeriods() []AcademicPeriodResponse {
	return []AcademicPeriodResponse{
		// 2021
		{"Autumn 2021", 16, time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2021, 12, 21, 23, 59, 59, 0, time.UTC)},
		{"Spring 2021", 14, time.Date(2021, 1, 11, 0, 0, 0, 0, time.UTC), time.Date(2021, 4, 16, 23, 59, 59, 0, time.UTC)},
		{"Summer 2021", 10, time.Date(2021, 5, 3, 0, 0, 0, 0, time.UTC), time.Date(2021, 7, 9, 23, 59, 59, 0, time.UTC)},

		// 2022
		{"Autumn 2022", 16, time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2022, 12, 21, 23, 59, 59, 0, time.UTC)},
		{"Spring 2022", 14, time.Date(2022, 1, 10, 0, 0, 0, 0, time.UTC), time.Date(2022, 4, 15, 23, 59, 59, 0, time.UTC)},
		{"Summer 2022", 10, time.Date(2022, 5, 2, 0, 0, 0, 0, time.UTC), time.Date(2022, 7, 8, 23, 59, 59, 0, time.UTC)},

		// 2023
		{"Autumn 2023", 16, time.Date(2023, 9, 4, 0, 0, 0, 0, time.UTC), time.Date(2023, 12, 22, 23, 59, 59, 0, time.UTC)},
		{"Spring 2023", 14, time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC), time.Date(2023, 4, 14, 23, 59, 59, 0, time.UTC)},
		{"Summer 2023", 10, time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 7, 7, 23, 59, 59, 0, time.UTC)},

		// 2024
		{"Autumn 2024", 16, time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC), time.Date(2024, 12, 20, 23, 59, 59, 0, time.UTC)},
		{"Spring 2024", 14, time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), time.Date(2024, 4, 12, 23, 59, 59, 0, time.UTC)},
		{"Summer 2024", 10, time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC), time.Date(2024, 7, 12, 23, 59, 59, 0, time.UTC)},

		// 2025
		{"Autumn 2025", 16, time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 12, 19, 23, 59, 59, 0, time.UTC)},
		{"Spring 2025", 14, time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC), time.Date(2025, 4, 11, 23, 59, 59, 0, time.UTC)},
		{"Summer 2025", 10, time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC), time.Date(2025, 7, 11, 23, 59, 59, 0, time.UTC)},
	}
}
