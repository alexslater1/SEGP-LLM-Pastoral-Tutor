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
			DegreeYear:  0,
			Degree:      "Bachelor of Science in Computer Science",
			DegreeShort: "BSc CS",
		},
		{
			DegreeYear:  1,
			Degree:      "Master of Engineering in Software Engineering",
			DegreeShort: "MEng SE",
		},
		{
			DegreeYear:  2,
			Degree:      "Doctor of Philosophy in Data Science",
			DegreeShort: "PhD DS",
		},
	}
}

type AcademicPeriodResponse struct {
	Name  string    `json:"name"`
	Weeks int       `json:"weeks"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (m *MockAbcApiClient) GetAcademicPeriods() [][]AcademicPeriodResponse {
	const layout = "2006-01-02"

	// Create a sample academic period.
	period, err := time.Parse(layout, "2021-10-01")
	if err != nil {
		panic(err)
	}
	periodEnd, err := time.Parse(layout, "2022-12-17")
	if err != nil {
		panic(err)
	}

	mockPeriod := AcademicPeriodResponse{
		Name:  "autumn term",
		Weeks: 11,
		Start: period,
		End:   periodEnd,
	}

	// Return a nested slice of AcademicPeriod.
	return [][]AcademicPeriodResponse{
		{mockPeriod},
	}
}
