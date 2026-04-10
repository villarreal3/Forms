package helpers

import (
	"strings"

	"forms/models"
)

// sqlInjectionPatterns son subcadenas que se consideran indicativas de intento de SQL injection.
// Se evitan comillas sueltas para no bloquear texto legítimo (ej. apellido "O'Brien").
var sqlInjectionPatterns = []string{
	";",
	"--",
	"/*",
	"*/",
	" union ",
	" select ",
	" insert ",
	" update ",
	" delete ",
	" drop ",
	"exec(",
	"execute(",
	"xp_",
	"char(",
	"nchar(",
	"concat(",
	"benchmark(",
	"sleep(",
	" into ",
	" from ",
	" where ",
	" or ",
	" and ",
	" like ",
}

// ContainsSQLInjectionPattern indica si la cadena contiene algún patrón típico de SQL injection.
func ContainsSQLInjectionPattern(s string) bool {
	lower := strings.ToLower(s)
	// Normalizar con espacios para que patrones " palabra " coincidan al inicio/fin del valor
	normalized := " " + lower + " "
	for _, p := range sqlInjectionPatterns {
		if strings.Contains(normalized, p) {
			return true
		}
	}
	return false
}

// ValidateNoSQLInjectionPatterns revisa todos los campos string del request de envío de formulario
// y rechaza si alguno contiene patrones considerados peligrosos.
// Devuelve (true, "") si todo es válido, (false, mensaje) si se detecta un patrón.
func ValidateNoSQLInjectionPatterns(req *models.SubmitFormRequest) (ok bool, message string) {
	msg := "Caracteres o patrones no permitidos en uno o más campos"
	check := func(v string) bool {
		return !ContainsSQLInjectionPattern(v)
	}
	if !check(req.FirstName) || !check(req.LastName) || !check(req.IDNumber) || !check(req.Email) {
		return false, msg
	}
	if !check(req.Phone) || !check(req.Province) || !check(req.Gender) {
		return false, msg
	}
	if !check(req.BusinessName) || !check(req.BusinessRegistration) {
		return false, msg
	}
	if !check(req.Instagram) || !check(req.Facebook) || !check(req.TikTok) || !check(req.Twitter) {
		return false, msg
	}
	for _, a := range req.Answers {
		if !check(a.FieldValue) {
			return false, msg
		}
	}
	return true, ""
}
