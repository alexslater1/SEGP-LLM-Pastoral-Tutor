type PublicCourseResponse struct {
    Code  string `json:"code"`
    Title string `json:"title"`
    ECTS  int    `json:"ects"`
}

type PublicModuleTypesResponse struct {
    Title         string `json:"title"`
    Code          string `json:"code"`
    SyllabusLabel string `json:"syllabus_label"`
}

func (m *MockAbcApiClient) GetPublicCourses() (*PublicCourseResponse, error) {
	return &PublicCourseResponse{
		Code:  "50002",
		Title: "Software Engineering Design",
		ECTS:  5,
	}, nil
}

func (m *MockAbcApiClient) GetPublicModuleTypes() (*PublicModuleTypesResponse, error) {
	return &PublicModuleTypesResponse{
		Title:         "Software Engineering Design",
		Code:          "50002",
		SyllabusLabel: "Software Engineering Design",
	}, nil
}