package memory

type ReActMemory struct {
	CognitiveCycle []ReActMemorySteps
}

type ReActMemorySteps struct {
	Thought     string
	Action      string
	Observation string
}

func (m *ReActMemory) Add(input ReActMemorySteps) error {
	m.CognitiveCycle = append(m.CognitiveCycle, input)
	return nil
}

func (m *ReActMemory) Get() ([]ReActMemorySteps, error) {
	return m.CognitiveCycle, nil
}
