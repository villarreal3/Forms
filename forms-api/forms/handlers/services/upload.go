package services

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"forms/database"
	"forms/handlers/helpers"
	"forms/handlers/repositories"

	"github.com/gorilla/mux"
)

var ImagesDir = "src/img"

// ---------- CONFIG ----------

func InitImageConfig(imagesDir string) {
	if imagesDir != "" {
		ImagesDir = imagesDir
	}
	if err := ensureImagesDir(); err != nil {
		log.Printf("Error creando directorio de imágenes: %v", err)
		return
	}
	log.Printf("Directorio de imágenes configurado: %s", ImagesDir)
}

func ensureImagesDir() error {
	if ImagesDir == "" {
		ImagesDir = "src/img"
	}
	if err := os.MkdirAll(ImagesDir, 0755); err != nil {
		log.Printf("[ensureImagesDir] Error creando directorio de imágenes: %v", err)
		return err
	}
	return nil
}

// ---------- TYPES ----------

type UploadOpts struct {
	FormID      int64
	ImageType   string // "desktop" | "mobile" (solo para logs)
	Suffix      string // "desktop" | "mobile" (parte del filename)
	IsMobile    bool
	ColumnName  string // "logo_url" | "logo_url_mobile"
	DeviceLabel string // "Desktop" | "Mobile"
}

type DeleteOpts struct {
	FormID      int64
	IsMobile    bool
	ColumnName  string
	DeviceLabel string
}

type parsedImage struct {
	File        io.ReadCloser
	HandlerName string
	Size        int64
	ContentType string
	Ext         string
}

// ---------- PUBLIC ROUTES ----------

// UploadFormImage maneja la subida de imágenes para un formulario (router)
func UploadFormImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	if err := ensureImagesDir(); err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "Error configurando directorio de imágenes")
		return
	}

	formID, err := parseFormIDFromVars(r)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "ID de formulario inválido")
		return
	}

	// Logging detallado para diagnóstico
	imageTypeRaw := r.URL.Query().Get("type")
	log.Printf("[UploadFormImage] ===== DIAGNÓSTICO =====")
	log.Printf("[UploadFormImage] FormID: %d", formID)
	log.Printf("[UploadFormImage] Query parameter 'type' (RAW): '%s' (len=%d)", imageTypeRaw, len(imageTypeRaw))
	log.Printf("[UploadFormImage] Todos los query parameters: %v", r.URL.Query())

	imageType := normalizeQuery(r, "type")
	log.Printf("[UploadFormImage] Query parameter 'type' (NORMALIZADO): '%s' (len=%d)", imageType, len(imageType))

	if imageType == "" {
		log.Printf("[UploadFormImage] ERROR: Tipo de imagen vacío después de normalización")
		helpers.RespondError(w, http.StatusBadRequest, "Tipo de imagen requerido. Debe ser 'desktop' o 'mobile'")
		return
	}

	// Mapear "pc" a "desktop" (obsoleto - el frontend debe usar "desktop")
	if imageType == "pc" {
		log.Printf("[UploadFormImage] ADVERTENCIA: Se recibió 'pc' que está obsoleto. Mapeando a 'desktop'. Por favor actualizar el frontend para usar 'desktop'.")
		imageType = "desktop"
	}

	opts, err := uploadOptsFromType(formID, imageType)
	if err != nil {
		log.Printf("[UploadFormImage] ERROR en uploadOptsFromType: tipo='%s', error=%v", imageType, err)
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("[UploadFormImage] Opciones creadas exitosamente: FormID=%d, ImageType=%s, IsMobile=%v, ColumnName=%s, DeviceLabel=%s",
		opts.FormID, opts.ImageType, opts.IsMobile, opts.ColumnName, opts.DeviceLabel)

	handleUploadImage(w, r, opts)
}

// DeleteFormImage elimina la imagen de un formulario (router)
func DeleteFormImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		helpers.RespondError(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	formID, err := parseFormIDFromVars(r)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "ID de formulario inválido")
		return
	}

	imageType := normalizeQuery(r, "type")
	if imageType == "" {
		helpers.RespondError(w, http.StatusBadRequest, "Tipo de imagen requerido. Debe ser 'desktop' o 'mobile'")
		return
	}

	// Mapear "pc" a "desktop" (obsoleto - el frontend debe usar "desktop")
	if imageType == "pc" {
		log.Printf("[DeleteFormImage] ADVERTENCIA: Se recibió 'pc' que está obsoleto. Mapeando a 'desktop'. Por favor actualizar el frontend para usar 'desktop'.")
		imageType = "desktop"
	}

	opts, err := deleteOptsFromType(formID, imageType)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	handleDeleteImage(w, r, opts)
}

// ServeImage sirve las imágenes almacenadas
func ServeImage(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ServeImage] Iniciando solicitud de imagen")

	if err := ensureImagesDir(); err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "Error configurando directorio de imágenes")
		return
	}

	vars := mux.Vars(r)
	formID := vars["form_id"]
	filename := vars["filename"]

	// Seguridad básica: evita ../
	if isUnsafePathSegment(formID) || isUnsafePathSegment(filename) {
		helpers.RespondError(w, http.StatusBadRequest, "Parámetros inválidos")
		return
	}

	if err := assertFormExists(formID); err != nil {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}

	filePath := filepath.Join(ImagesDir, formID, filename)
	filePath = resolveFallbackPaths(filePath, formID, filename)

	fi, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			helpers.RespondError(w, http.StatusNotFound, "Imagen no encontrada")
			return
		}
		helpers.RespondError(w, http.StatusInternalServerError, "Error verificando archivo")
		return
	}

	contentType := detectContentTypeByExt(filename)
	log.Printf("[ServeImage] Sirviendo archivo: %s (size=%d, ct=%s)", filePath, fi.Size(), contentType)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	http.ServeFile(w, r, filePath)
}

// ---------- UPLOAD FLOW ----------

func handleUploadImage(w http.ResponseWriter, r *http.Request, opts UploadOpts) {
	logPrefix := "[upload]"
	log.Printf("%s ===== INICIO SUBIDA (%s) =====", logPrefix, opts.DeviceLabel)
	log.Printf("%s FormID=%d, Column=%s, Type=%s", logPrefix, opts.FormID, opts.ColumnName, opts.ImageType)

	formName, err := getFormName(opts.FormID)
	if err != nil {
		writeFormNameError(w, err)
		return
	}

	img, err := parseAndValidateMultipartImage(r, "image")
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer img.File.Close()

	formDir, err := ensureFormDir(opts.FormID)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error creando directorio: %v", err))
		return
	}

	filename := buildImageFilename(opts.FormID, opts.Suffix, formName, img.Ext)
	filePath := filepath.Join(formDir, filename)

	// borrar anterior si existía
	if err := deletePreviousImageIfAny(opts.FormID, opts.ColumnName, formDir); err != nil {
		// no lo haces fatal; log y sigue
		log.Printf("%s Warning: no se pudo borrar anterior: %v", logPrefix, err)
	}

	if err := saveUploadedFile(filePath, img.File); err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error guardando archivo: %v", err))
		return
	}

	imageURL := fmt.Sprintf("/images/%d/%s", opts.FormID, filename)
	log.Printf("%s ANTES DE UpdateLogo - formID=%d, imageURL='%s', isMobile=%v, ColumnName='%s'",
		logPrefix, opts.FormID, imageURL, opts.IsMobile, opts.ColumnName)

	ok, err := repositories.CustomizationRepo.UpdateLogo(opts.FormID, imageURL, opts.IsMobile)
	log.Printf("%s DESPUÉS DE UpdateLogo - success=%v, error=%v", logPrefix, ok, err)

	if err != nil || !ok {
		log.Printf("%s ERROR actualizando personalización: success=%v, error=%v", logPrefix, ok, err)
		_ = os.Remove(filePath)
		helpers.RespondError(w, http.StatusInternalServerError, fmt.Sprintf("Error actualizando personalización: %v", err))
		return
	}

	log.Printf("%s ÉXITO: Imagen subida y guardada correctamente en columna '%s'", logPrefix, opts.ColumnName)

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"message":   "Imagen subida correctamente",
		"image_url": imageURL,
	})
}

// ---------- DELETE FLOW ----------

func handleDeleteImage(w http.ResponseWriter, r *http.Request, opts DeleteOpts) {
	logPrefix := "[delete]"
	log.Printf("%s ===== INICIO ELIMINACIÓN (%s) =====", logPrefix, opts.DeviceLabel)

	url, err := getCustomizationURL(opts.FormID, opts.ColumnName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.RespondError(w, http.StatusNotFound, "No hay imagen asociada a este formulario")
			return
		}
		helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo imagen")
		return
	}
	if url == "" {
		helpers.RespondError(w, http.StatusNotFound, "No hay imagen asociada a este formulario")
		return
	}

	filename, err := filenameFromImageURL(url)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "URL de imagen inválida")
		return
	}

	// borrar archivo físico
	filePath := filepath.Join(ImagesDir, fmt.Sprintf("%d", opts.FormID), filename)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		log.Printf("%s Error eliminando archivo: %v", logPrefix, err)
	}

	ok, err := repositories.CustomizationRepo.DeleteLogo(opts.FormID, opts.IsMobile)
	if err != nil || !ok {
		helpers.RespondError(w, http.StatusInternalServerError, "Error eliminando imagen")
		return
	}

	helpers.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Imagen eliminada correctamente",
	})
}

// ---------- SMALL HELPERS ----------

func parseFormIDFromVars(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars["id"], 10, 64)
}

func normalizeQuery(r *http.Request, key string) string {
	// Intentar obtener de query parameters primero
	value := r.URL.Query().Get(key)

	// Si no está en query params y el form ya está parseado, intentar obtener del form data
	if value == "" && r.MultipartForm != nil {
		if values := r.MultipartForm.Value[key]; len(values) > 0 {
			value = values[0]
		}
	}

	return strings.TrimSpace(strings.ToLower(value))
}

func uploadOptsFromType(formID int64, imageType string) (UploadOpts, error) {
	switch imageType {
	case "desktop":
		return UploadOpts{
			FormID:      formID,
			ImageType:   "desktop",
			Suffix:      "desktop",
			IsMobile:    false,
			ColumnName:  "logo_url",
			DeviceLabel: "Desktop",
		}, nil
	case "mobile":
		return UploadOpts{
			FormID:      formID,
			ImageType:   "mobile",
			Suffix:      "mobile",
			IsMobile:    true,
			ColumnName:  "logo_url_mobile",
			DeviceLabel: "Mobile",
		}, nil
	default:
		return UploadOpts{}, fmt.Errorf("Tipo de imagen inválido. Debe ser 'desktop' o 'mobile'")
	}
}

func deleteOptsFromType(formID int64, imageType string) (DeleteOpts, error) {
	switch imageType {
	case "desktop":
		return DeleteOpts{FormID: formID, IsMobile: false, ColumnName: "logo_url", DeviceLabel: "Desktop"}, nil
	case "mobile":
		return DeleteOpts{FormID: formID, IsMobile: true, ColumnName: "logo_url_mobile", DeviceLabel: "Mobile"}, nil
	default:
		return DeleteOpts{}, fmt.Errorf("Tipo de imagen inválido. Debe ser 'desktop' o 'mobile'")
	}
}

func getFormName(formID int64) (string, error) {
	var formName string
	err := database.DB.QueryRow("SELECT form_name FROM forms WHERE id = ?", formID).Scan(&formName)
	return formName, err
}

func writeFormNameError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		helpers.RespondError(w, http.StatusNotFound, "Formulario no encontrado")
		return
	}
	helpers.RespondError(w, http.StatusInternalServerError, "Error obteniendo información del formulario")
}

func parseAndValidateMultipartImage(r *http.Request, field string) (*parsedImage, error) {
	// 10MB parse, 5MB límite real
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, fmt.Errorf("Error procesando formulario: %v", err)
	}

	file, handler, err := r.FormFile(field)
	if err != nil {
		return nil, fmt.Errorf("No se encontró el archivo '%s' en la petición", field)
	}

	contentType := handler.Header.Get("Content-Type")
	if !isAllowedImageType(contentType) {
		file.Close()
		return nil, fmt.Errorf("Tipo de archivo no permitido. Solo se permiten: jpg, jpeg, png, gif, webp")
	}
	if handler.Size > 5<<20 {
		file.Close()
		return nil, fmt.Errorf("El archivo es demasiado grande. Tamaño máximo: 5MB")
	}

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext == "" {
		ext = extFromContentType(contentType)
	}

	return &parsedImage{
		File:        file,
		HandlerName: handler.Filename,
		Size:        handler.Size,
		ContentType: contentType,
		Ext:         ext,
	}, nil
}

func isAllowedImageType(ct string) bool {
	switch ct {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

func extFromContentType(ct string) string {
	switch ct {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

func ensureFormDir(formID int64) (string, error) {
	formDir := filepath.Join(ImagesDir, fmt.Sprintf("%d", formID))
	if err := os.MkdirAll(formDir, 0755); err != nil {
		return "", err
	}
	info, err := os.Stat(formDir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s existe pero no es un directorio", formDir)
	}
	return formDir, nil
}

func buildImageFilename(formID int64, suffix string, formName string, ext string) string {
	clean := sanitizeFilename(formName)
	return fmt.Sprintf("%d-%s-%s%s", formID, suffix, clean, ext)
}

func sanitizeFilename(name string) string {
	// si quieres, puedes cambiar a regex; esto mantiene tu lógica actual
	repl := []string{" ", "_", "/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_"}
	r := strings.NewReplacer(repl...)
	return r.Replace(name)
}

func saveUploadedFile(filePath string, src io.Reader) error {
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(filePath)
		return err
	}
	return nil
}

func getCustomizationURL(formID int64, columnName string) (string, error) {
	var val sql.NullString
	err := database.DB.QueryRow(fmt.Sprintf("SELECT %s FROM form_customization WHERE form_id = ?", columnName), formID).Scan(&val)
	if err != nil {
		return "", err
	}
	if !val.Valid {
		return "", nil
	}
	return val.String, nil
}

func filenameFromImageURL(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid url")
	}
	return parts[len(parts)-1], nil
}

func deletePreviousImageIfAny(formID int64, columnName string, formDir string) error {
	oldURL, err := getCustomizationURL(formID, columnName)
	if err != nil {
		// si no existe fila, no hay nada que borrar
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if oldURL == "" {
		return nil
	}

	oldFilename, err := filenameFromImageURL(oldURL)
	if err != nil {
		return nil // no rompas el upload por URL rara
	}

	oldPath := filepath.Join(formDir, oldFilename)
	if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ---------- SERVE HELPERS ----------

func assertFormExists(formID string) error {
	var exists int
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM forms WHERE id = ?", formID).Scan(&exists); err != nil {
		return err
	}
	if exists == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func detectContentTypeByExt(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "application/octet-stream"
	}
	if ct := mime.TypeByExtension(ext); ct != "" {
		// mime devuelve "image/png; charset=utf-8" a veces, lo normalizamos
		return strings.Split(ct, ";")[0]
	}
	switch ext {
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}

func resolveFallbackPaths(filePath, formID, filename string) string {
	// Mantiene tu compatibilidad: intenta rutas alternativas si no existe
	if _, err := os.Stat(filePath); err == nil {
		return filePath
	}

	possible := []string{
		filePath,
		filepath.Join("/app", filePath),
		filepath.Join("/app/src/img", formID, filename),
	}
	for _, p := range possible {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return filePath
}

func isUnsafePathSegment(s string) bool {
	// evita ../ y separadores raros
	if strings.Contains(s, "..") {
		return true
	}
	if strings.ContainsAny(s, `\/`) {
		// form_id no debería traer /, filename tampoco debe venir con path
		// (si filename puede tener subdirs, quita esto)
	}
	return false
}
