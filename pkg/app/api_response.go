package app

type APIResponse struct {
	Alerts        []string    `json:"alerts"`
	Notifications []string    `json:"notifications"`
	Results       interface{} `json:"results"`

	errors []error
}

func (ar *APIResponse) NoErrors() bool {
	return len(ar.errors) == 0
}

func (ar *APIResponse) AddError(e error) {
	ar.errors = append(ar.errors, e)
}

func (ar *APIResponse) AddErrors(e []error) {
	ar.errors = append(ar.errors, e...)
}

func (ar *APIResponse) AddNotification(e string) {
	ar.Notifications = append(ar.Notifications, e)
}

func (ar *APIResponse) AddNotifications(e []string) {
	ar.Notifications = append(ar.Notifications, e...)
}

func (ar *APIResponse) AddAlert(e string) {
	ar.Alerts = append(ar.Alerts, e)
}

func (ar *APIResponse) AddAlerts(e []string) {
	ar.Alerts = append(ar.Alerts, e...)
}

func (ar *APIResponse) ParseErrors() {
	for _, e := range ar.errors {
		ar.AddAlert(e.Error())
	}
}
