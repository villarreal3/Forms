package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
	RoleKey     contextKey = "role"
)

var jwtSecret = []byte("forms-secret-key-change-in-production") // Cambiar en producción
var loginServiceURL = "http://localhost:8080"                           // URL del servicio de login

// SetJWTSecret establece la clave secreta para JWT
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// SetLoginServiceURL establece la URL del servicio de login
func SetLoginServiceURL(url string) {
	loginServiceURL = url
}

// Claims representa los claims del JWT (formato legacy de forms-api)
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// UsersAPIClaims representa los claims del JWT emitidos por users-api
type UsersAPIClaims struct {
	Sub    interface{} `json:"sub"`    // userId puede ser string o número
	Email  string      `json:"email"`   // email del usuario
	Roles  []string    `json:"roles"`  // array de roles
	jwt.RegisteredClaims
}

// GenerateToken genera un token JWT para un usuario
func GenerateToken(userID int64, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "forms",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken valida un token JWT (formato legacy de forms-api)
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// ValidateUsersAPIToken valida un token JWT emitido por users-api
func ValidateUsersAPIToken(tokenString string) (*UsersAPIClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UsersAPIClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea el esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UsersAPIClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// mapUsersAPIRole mapea roles de users-api a roles de forms-api
func mapUsersAPIRole(usersAPIRoles []string) string {
	// Prioridad: super_admin > admin > editor > auditor > user
	for _, role := range usersAPIRoles {
		switch role {
		case "super_admin", "admin":
			return "admin"
		case "editor":
			return "editor"
		case "auditor", "user":
			return "viewer"
		}
	}
	// Si no hay roles reconocidos, usar viewer por defecto
	return "viewer"
}

// extractUserID extrae el user ID del campo sub (puede ser string UUID o número)
func extractUserID(sub interface{}) (int64, error) {
	switch v := sub.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		// Si es un UUID (formato con guiones), usar hash del UUID como ID
		// Esto permite mantener compatibilidad con el sistema que espera int64
		if len(v) > 10 && (v[8] == '-' || v[4] == '-') {
			// Es un UUID, generar un hash numérico del UUID
			// Usar los primeros caracteres numéricos del UUID o un hash simple
			var hash int64
			for _, char := range v {
				if char >= '0' && char <= '9' {
					hash = hash*10 + int64(char-'0')
					if hash > 1000000000 { // Limitar el tamaño
						hash = hash % 1000000000
					}
				}
			}
			// Si no hay suficientes dígitos, usar un hash basado en la suma de caracteres
			if hash == 0 {
				for _, char := range v {
					hash += int64(char)
				}
				hash = hash % 1000000000
			}
			return hash, nil
		}
		// Intentar convertir string numérico a int64
		var id int64
		_, err := fmt.Sscanf(v, "%d", &id)
		if err != nil {
			// Si no es numérico, generar un hash del string
			var hash int64
			for _, char := range v {
				hash += int64(char)
			}
			return hash % 1000000000, nil
		}
		return id, nil
	default:
		return 0, fmt.Errorf("tipo de sub no soportado: %T", v)
	}
}

// validateTokenViaLoginService valida un token llamando al servicio de login
func validateTokenViaLoginService(tokenString string) (*Claims, error) {
	// Preparar request
	reqBody := map[string]string{"token": tokenString}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error codificando request: %v", err)
	}

	// Hacer request al servicio de login
	resp, err := http.Post(loginServiceURL+"/api/auth/validate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error llamando servicio de login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token inválido")
	}

	// Decodificar respuesta
	var result struct {
		Success  bool   `json:"success"`
		Valid    bool   `json:"valid"`
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %v", err)
	}

	if !result.Success || !result.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	// Crear claims desde la respuesta
	claims := &Claims{
		UserID:   result.UserID,
		Username: result.Username,
		Role:     result.Role,
	}

	return claims, nil
}

// AuthMiddleware verifica la autenticación del usuario
// Ahora valida tokens directamente emitidos por users-api
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir peticiones OPTIONS (preflight de CORS) sin autenticación
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Token de autenticación requerido")
			return
		}

		// Formato: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, http.StatusUnauthorized, "Formato de token inválido. Use: Bearer <token>")
			return
		}

		tokenString := parts[1]
		
		// Intentar validar token como token de users-api primero
		usersAPIClaims, err := ValidateUsersAPIToken(tokenString)
		if err == nil && usersAPIClaims != nil {
			// Token válido de users-api
			userID, err := extractUserID(usersAPIClaims.Sub)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "Error procesando token: "+err.Error())
				return
			}

			// Mapear roles de users-api a roles de forms-api
			mappedRole := mapUsersAPIRole(usersAPIClaims.Roles)

			// Agregar información del usuario al contexto
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UsernameKey, usersAPIClaims.Email)
			ctx = context.WithValue(ctx, RoleKey, mappedRole)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Si falla, intentar validar como token legacy de forms-api (para compatibilidad)
		legacyClaims, err := ValidateToken(tokenString)
		if err == nil && legacyClaims != nil {
			// Token válido legacy
			ctx := context.WithValue(r.Context(), UserIDKey, legacyClaims.UserID)
			ctx = context.WithValue(ctx, UsernameKey, legacyClaims.Username)
			ctx = context.WithValue(ctx, RoleKey, legacyClaims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Si ambos fallan, token inválido
		respondError(w, http.StatusUnauthorized, "Token inválido o expirado")
	})
}

// RoleMiddleware verifica que el usuario tenga uno de los roles permitidos
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Permitir peticiones OPTIONS (preflight de CORS) sin verificación de rol
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}
			
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				respondError(w, http.StatusForbidden, "No se pudo determinar el rol del usuario")
				return
			}

			allowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					allowed = true
					break
				}
			}

			if !allowed {
				respondError(w, http.StatusForbidden, "No tienes permisos para acceder a este recurso")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID obtiene el ID del usuario del contexto
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// GetUsername obtiene el username del contexto
func GetUsername(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

// GetRole obtiene el rol del usuario del contexto
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": message,
	})
}

