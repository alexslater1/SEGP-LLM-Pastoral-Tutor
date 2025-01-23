package googleSearch

import "time"

type GoogleSearchClient interface {
	HtmlFromQuery(query string) (*string, error)
	HtmlFromURL(url string, actions ...Action) (*string, error)
}

type ActionType string

const (
	ActionTypeClick ActionType = "click"
	ActionTypeWait  ActionType = "wait"
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

func NewClickActionCloseGoogleCookies() *ClickAction {
	return &ClickAction{Element: ".UywwFc-LgbsSe.UywwFc-LgbsSe-OWXEXe-dgl2Hf.XWZjwc"}
}

type WaitAction struct {
	Element  *string
	Duration *time.Duration
}

func NewWaitAction(element string, duration time.Duration) *WaitAction {
	return &WaitAction{Element: &element, Duration: &duration}
}

func NewWaitDurationAction(duration time.Duration) *WaitAction {
	return &WaitAction{Duration: &duration}
}

func NewWaitElementAction(element string) *WaitAction {
	return &WaitAction{Element: &element}
}

func (w *WaitAction) Type() ActionType {
	return ActionTypeWait
}
