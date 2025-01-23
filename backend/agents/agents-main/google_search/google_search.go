package googleSearch

type GoogleSearchClient interface {
	HtmlFromQuery(query string) (*string, error)
	HtmlFromURL(url string, actions ...Action) (*string, error)
}

type ActionType string

const (
	ActionTypeClick ActionType = "click"
)

type Action interface {
	Type() ActionType
}

type ClickAction struct {
	Element string
}

func (c *ClickAction) Type() ActionType {
	return ActionTypeClick
}

func NewClickAction(element string) *ClickAction {
	return &ClickAction{Element: element}
}
