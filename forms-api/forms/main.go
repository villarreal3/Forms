package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"forms/config"
	"forms/database"
	"forms/handlers"
	"forms/middleware"

	"github.com/gorilla/mux"
)

//go:embed docs/*
var docsContent embed.FS

// CORS middleware simple
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "*")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		h.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// helper para registrar rutas (método + OPTIONS)
func addRoute(r *mux.Router, method, path string, h http.HandlerFunc) {
	r.HandleFunc(path, h).Methods(method, http.MethodOptions)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer database.Close()
	fmt.Println("✅ DB conectada:", cfg.DBName)

	middleware.SetJWTSecret(cfg.JWTSecret)
	middleware.SetLoginServiceURL(cfg.LoginServiceURL)
	fmt.Println("✅ JWT Secret configurado (compartido con users-api)")
	if cfg.UsersAPIURL != "" {
		fmt.Println("✅ Users API URL:", cfg.UsersAPIURL)
	}

	handlers.InitEmailConfig(cfg.EmailsDir)
	fmt.Println("✅ EmailsDir:", cfg.EmailsDir)

	handlers.InitImageConfig(cfg.ImagesDir)
	fmt.Println("✅ ImagesDir:", cfg.ImagesDir)

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// Documentación OpenAPI (go-swagger / Swagger UI) — antes del subrouter /api para no ser absorbido
	docsSub, err := fs.Sub(docsContent, "docs")
	if err != nil {
		log.Fatalf("docs embed: %v", err)
	}
	r.HandleFunc("/api/docs", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/api/docs/", http.StatusFound)
	})
	r.PathPrefix("/api/docs/").Handler(http.StripPrefix("/api/docs/", http.FileServer(http.FS(docsSub))))

	// API base
	api := r.PathPrefix("/api").Subrouter()

	// Públicos
	addRoute(api, http.MethodPost, "/forms/submit", handlers.SubmitForm)
	addRoute(api, http.MethodPost, "/forms/auto-submit", handlers.AutoSubmitFromPrevious)
	addRoute(api, http.MethodGet, "/public/forms", handlers.GetPublicForms)
	addRoute(api, http.MethodGet, "/public/forms/{id}", handlers.GetPublicForm)
	addRoute(api, http.MethodGet, "/public/forms/{id}/sections", handlers.GetPublicFormSections)
	addRoute(api, http.MethodGet, "/public/forms/{id}/customization", handlers.GetFormCustomization)
	addRoute(api, http.MethodGet, "/public/forms/{id}/logo", handlers.GetPublicFormLogo)
	addRoute(api, http.MethodGet, "/images/{form_id}/{filename}", handlers.ServeImage)
	addRoute(api, http.MethodPost, "/users/get", handlers.GetUserByCredentials)
	addRoute(api, http.MethodGet, "/health", handlers.HealthCheck)

	// Admin (protegido)
	protected := api.PathPrefix("/admin").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Formularios
	addRoute(protected, http.MethodPost, "/forms", handlers.CreateForm)
	addRoute(protected, http.MethodGet, "/forms", handlers.GetForms)
	addRoute(protected, http.MethodGet, "/forms/{id}", handlers.GetForm)
	addRoute(protected, http.MethodPatch, "/forms/{id}", handlers.PatchFormMetadata)
	addRoute(protected, http.MethodPut, "/forms/{id}/schema", handlers.UpdateFormSchema)
	addRoute(protected, http.MethodPost, "/forms/{id}/close", handlers.CloseForm)

	// Publicar / pasar a borrador (editor, admin)
	visibilityRouter := protected.PathPrefix("/forms").Subrouter()
	visibilityRouter.Use(middleware.RoleMiddleware("editor", "admin"))
	addRoute(visibilityRouter, http.MethodPost, "/{id}/publish", handlers.PublishForm)
	addRoute(visibilityRouter, http.MethodPost, "/{id}/draft", handlers.SetFormDraft)
	addRoute(visibilityRouter, http.MethodPost, "/{id}/open", handlers.OpenForm)

	addRoute(protected, http.MethodGet, "/fields/templates", handlers.GetCommonFieldTemplates)
	addRoute(protected, http.MethodGet, "/forms/{form_id}/sections", handlers.GetFormSections)

	// Dashboard
	addRoute(protected, http.MethodGet, "/dashboard/stats", handlers.GetDashboardStats)
	addRoute(protected, http.MethodGet, "/dashboard/forms", handlers.GetFormStats)
	addRoute(protected, http.MethodGet, "/dashboard/submissions", handlers.GetRecentSubmissions)
	addRoute(protected, http.MethodGet, "/dashboard/responses-timeline", handlers.GetResponseTimeline)
	addRoute(protected, http.MethodGet, "/dashboard/top-active-forms", handlers.GetTopActiveForms)
	addRoute(protected, http.MethodGet, "/dashboard/geographic-distribution", handlers.GetGeographicDistribution)
	addRoute(protected, http.MethodGet, "/dashboard/user-metrics", handlers.GetUserMetrics)

	// Submissions
	addRoute(protected, http.MethodGet, "/submissions", handlers.GetSubmissions)

	// Asistencia (editor, admin)
	attendanceRouter := protected.PathPrefix("/attendance").Subrouter()
	attendanceRouter.Use(middleware.RoleMiddleware("editor", "admin"))
	addRoute(attendanceRouter, http.MethodPost, "/update", handlers.UpdateAttendance)

	// Export (editor, admin)
	exportRouter := protected.PathPrefix("/export").Subrouter()
	exportRouter.Use(middleware.RoleMiddleware("editor", "admin"))
	addRoute(exportRouter, http.MethodPost, "/submissions", handlers.ExportSubmissions)

	// Customización (editor, admin)
	customRouter := protected.PathPrefix("/customization").Subrouter()
	customRouter.Use(middleware.RoleMiddleware("editor", "admin"))
	addRoute(customRouter, http.MethodGet, "/forms/{id}", handlers.GetFormCustomization)
	addRoute(customRouter, http.MethodPost, "/forms/{id}", handlers.CreateOrUpdateCustomization)
	addRoute(customRouter, http.MethodPut, "/forms/{id}", handlers.CreateOrUpdateCustomization)

	// Imágenes (editor, admin)
	imageRouter := protected.PathPrefix("/forms").Subrouter()
	imageRouter.Use(middleware.RoleMiddleware("editor", "admin"))
	addRoute(imageRouter, http.MethodPost, "/{id}/image", handlers.UploadFormImage)
	addRoute(imageRouter, http.MethodDelete, "/{id}/image", handlers.DeleteFormImage)

	// Emails masivos (admin)
	emailRouter := protected.PathPrefix("/email").Subrouter()
	emailRouter.Use(middleware.RoleMiddleware("admin"))
	addRoute(emailRouter, http.MethodPost, "/bulk", handlers.SendBulkEmail)

	// Root
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "🚀 API de Formularios funcionando correctamente forms-app")
	})

	// Logging de endpoints
	addr := ":" + cfg.ServerPort
	base := "http://172.28.125.76" + addr

	publicEPs := []struct {
		m, p string
	}{
		{"POST", "/api/forms/submit"},
		{"GET ", "/api/public/forms"},
		{"GET ", "/api/public/forms/{id}"},
		{"GET ", "/api/public/forms/{id}/sections"},
		{"GET ", "/api/public/forms/{id}/customization"},
		{"GET ", "/api/public/forms/{id}/logo"},
		{"GET ", "/api/images/{form_id}/{filename}"},
		{"POST", "/api/users/get"},
		{"GET ", "/api/health"},
	}

	adminEPs := []struct {
		m, p string
	}{
		{"POST", "/api/admin/forms"},
		{"GET ", "/api/admin/forms"},
		{"GET ", "/api/admin/forms/{id}"},
		{"PUT ", "/api/admin/forms/{id}/schema"},
		{"POST", "/api/admin/forms/{id}/close"},
		{"GET ", "/api/admin/fields/templates"},
		{"GET ", "/api/admin/forms/{form_id}/sections"},
		{"GET ", "/api/admin/dashboard/stats"},
		{"GET ", "/api/admin/dashboard/forms"},
		{"GET ", "/api/admin/dashboard/submissions"},
		{"GET ", "/api/admin/dashboard/responses-timeline"},
		{"GET ", "/api/admin/dashboard/top-active-forms"},
		{"GET ", "/api/admin/dashboard/geographic-distribution"},
		{"GET ", "/api/admin/dashboard/user-metrics"},
		{"GET ", "/api/admin/submissions"},
		{"POST", "/api/admin/attendance/update"},
		{"POST", "/api/admin/export/submissions"},
		{"GET ", "/api/admin/customization/forms/{id}"},
		{"POST", "/api/admin/customization/forms/{id}"},
		{"PUT ", "/api/admin/customization/forms/{id}"},
		{"POST", "/api/admin/forms/{id}/image"},
		{"DELETE", "/api/admin/forms/{id}/image"},
		{"POST", "/api/admin/email/bulk"},
	}

	// Endpoints nuevos (table-per-form): esquema en JSON, sin CRUD de campos/secciones por API
	newEPs := []struct {
		m, p string
	}{
		{"PUT ", "/api/admin/forms/{id}/schema"},
		{"GET ", "/api/admin/forms/{id}"},
		{"GET ", "/api/admin/forms/{form_id}/sections"},
		{"POST", "/api/admin/forms"},
	}

	fmt.Printf("🌐 Servidor en http://0.0.0.0%s (accesible en %s)\n", addr, base)
	fmt.Println("📝 Endpoints nuevos (table-per-form):")
	for _, e := range newEPs {
		fmt.Printf("   %s %s%s\n", e.m, base, e.p)
	}
	fmt.Println("📝 Endpoints públicos:")
	for _, e := range publicEPs {
		fmt.Printf("   %s %s%s\n", e.m, base, e.p)
	}
	fmt.Println("📝 Endpoints de administración:")
	for _, e := range adminEPs {
		fmt.Printf("   %s %s%s\n", e.m, base, e.p)
	}

	if err := http.ListenAndServe("0.0.0.0"+addr, r); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
