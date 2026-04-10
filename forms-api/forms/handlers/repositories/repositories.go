package repositories

// Repositorios globales (se inicializan en main.go o en un init)
var (
	FormRepo          FormRepository
	SubmissionRepo    SubmissionRepository
	DashboardRepo     DashboardRepository
	EmailRepo         EmailRepository
	CustomizationRepo CustomizationRepository
	PublicRepo        PublicRepository
	ExportRepo        ExportRepository
	RoleRepo          RoleRepository
)

// Init inicializa todos los repositories
func Init() {
	FormRepo = NewFormRepository()
	SubmissionRepo = NewSubmissionRepository()
	DashboardRepo = NewDashboardRepository()
	EmailRepo = NewEmailRepository()
	CustomizationRepo = NewCustomizationRepository()
	PublicRepo = NewPublicRepository()
	ExportRepo = NewExportRepository()
	RoleRepo = NewRoleRepository()
}














