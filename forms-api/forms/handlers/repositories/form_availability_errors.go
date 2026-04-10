package repositories

// FormUnavailableReason representa por qué un formulario no acepta envíos.
type FormUnavailableReason string

const (
	FormUnavailableDraft        FormUnavailableReason = "draft"
	FormUnavailableNotOpenYet   FormUnavailableReason = "not_open_yet"
	FormUnavailableExpired      FormUnavailableReason = "expired"
	FormUnavailableClosed       FormUnavailableReason = "closed"
	FormUnavailableLimitReached FormUnavailableReason = "limit_reached"
)

// FormUnavailableError se devuelve cuando el formulario existe pero no está disponible para recibir envíos.
type FormUnavailableError struct {
	Reason     FormUnavailableReason
	Message    string
	StatusCode int
}

func (e *FormUnavailableError) Error() string {
	return e.Message
}
