package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"forms/database"
	"forms/models"
	"forms/schema"
)

type mysqlDashboardRepository struct{}

func NewDashboardRepository() DashboardRepository {
	return &mysqlDashboardRepository{}
}

func (r *mysqlDashboardRepository) GetDashboardStats() (*models.DashboardStats, error) {
	var stats models.DashboardStats
	err := database.DB.QueryRow("SELECT * FROM v_dashboard_stats").Scan(
		&stats.TotalForms, &stats.ActiveForms, &stats.ExpiredForms,
		&stats.TotalUsers, &stats.TotalSubmissions, &stats.SubmissionsToday,
		&stats.SubmissionsWeek, &stats.SubmissionsMonth,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *mysqlDashboardRepository) GetFormStats() ([]models.FormStats, error) {
	rows, err := database.DB.Query("SELECT * FROM v_form_stats ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formStats []models.FormStats
	for rows.Next() {
		var stat models.FormStats
		var lastSubmission sql.NullTime
		var isClosed sql.NullBool
		err := rows.Scan(
			&stat.ID, &stat.FormName, &stat.CreatedAt, &stat.ExpiresAt,
			&isClosed, &stat.SubmissionCount, &stat.UniqueUsers, &lastSubmission, &stat.Status,
		)
		if err != nil {
			log.Printf("Error escaneando estadísticas: %v", err)
			continue
		}
		if lastSubmission.Valid {
			stat.LastSubmission = lastSubmission.Time
		}
		formStats = append(formStats, stat)
	}
	return formStats, nil
}

func (r *mysqlDashboardRepository) GetRecentSubmissions(limit int) ([]struct {
	ID          int64
	SubmittedAt time.Time
	FormName    string
	FirstName   string
	LastName    string
	Email       string
	IDNumber    string
}, error) {
	// form_submissions + forms; nombres desde form_{id}_responses vía query por fila
	rows, err := database.DB.Query(`
		SELECT fs.response_id, fs.submitted_at, f.id, f.form_name
		FROM form_submissions fs
		INNER JOIN forms f ON f.id = fs.form_id
		ORDER BY fs.submitted_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type rowMeta struct {
		responseID  int64
		submittedAt time.Time
		formID      int64
		formName    string
	}
	var metas []rowMeta
	for rows.Next() {
		var m rowMeta
		if err := rows.Scan(&m.responseID, &m.submittedAt, &m.formID, &m.formName); err != nil {
			continue
		}
		metas = append(metas, m)
	}

	var submissions []struct {
		ID          int64
		SubmittedAt time.Time
		FormName    string
		FirstName   string
		LastName    string
		Email       string
		IDNumber    string
	}
	for _, m := range metas {
		sub := struct {
			ID          int64
			SubmittedAt time.Time
			FormName    string
			FirstName   string
			LastName    string
			Email       string
			IDNumber    string
		}{
			ID:          m.responseID,
			SubmittedAt: m.submittedAt,
			FormName:    m.formName,
		}
		table := schema.ResponseTableName(m.formID)
		var nombre, apellido, email, cedula sql.NullString
		err := database.DB.QueryRow(
			fmt.Sprintf("SELECT nombre, apellido, email, cedula FROM `%s` WHERE id = ?", table),
			m.responseID,
		).Scan(&nombre, &apellido, &email, &cedula)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("GetRecentSubmissions row scan: %v", err)
		}
		if nombre.Valid {
			sub.FirstName = nombre.String
		}
		if apellido.Valid {
			sub.LastName = apellido.String
		}
		if email.Valid {
			sub.Email = email.String
		}
		if cedula.Valid {
			sub.IDNumber = cedula.String
		}
		submissions = append(submissions, sub)
	}
	return submissions, nil
}

func (r *mysqlDashboardRepository) GetResponseTimeline(days int) ([]models.ResponseTimeline, error) {
	// Agregado global desde form_submissions
	rows, err := database.DB.Query(`
		SELECT DATE(submitted_at) as date, COUNT(*) as count
		FROM form_submissions
		WHERE submitted_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(submitted_at)
		ORDER BY date ASC`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var timeline []models.ResponseTimeline
	for rows.Next() {
		var item models.ResponseTimeline
		if err := rows.Scan(&item.Date, &item.Count); err != nil {
			continue
		}
		timeline = append(timeline, item)
	}
	return timeline, nil
}

func (r *mysqlDashboardRepository) GetTopActiveForms(limit int) ([]models.TopForm, error) {
	rows, err := database.DB.Query(`
		SELECT f.form_name, COUNT(fs.id) as submission_count
		FROM forms f
		LEFT JOIN form_submissions fs ON f.id = fs.form_id
		GROUP BY f.id, f.form_name
		ORDER BY submission_count DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var forms []models.TopForm
	for rows.Next() {
		var form models.TopForm
		if err := rows.Scan(&form.FormName, &form.SubmissionCount); err != nil {
			continue
		}
		forms = append(forms, form)
	}
	return forms, nil
}

func (r *mysqlDashboardRepository) GetGeographicDistribution() ([]models.ProvinceDistribution, error) {
	forms, err := database.DB.Query("SELECT id FROM forms")
	if err != nil {
		return []models.ProvinceDistribution{}, nil
	}
	defer forms.Close()
	var counts = make(map[string]int64)
	for forms.Next() {
		var formID int64
		if err := forms.Scan(&formID); err != nil {
			continue
		}
		table := schema.ResponseTableName(formID)
		rrows, err := database.DB.Query("SELECT provincia FROM `" + table + "` WHERE provincia IS NOT NULL AND provincia != ''")
		if err != nil {
			continue
		}
		for rrows.Next() {
			var prov sql.NullString
			if err := rrows.Scan(&prov); err != nil {
				continue
			}
			if prov.Valid && prov.String != "" {
				counts[prov.String]++
			}
		}
		rrows.Close()
	}
	var distributions []models.ProvinceDistribution
	for prov, c := range counts {
		distributions = append(distributions, models.ProvinceDistribution{Province: prov, Count: int(c)})
	}
	return distributions, nil
}

func (r *mysqlDashboardRepository) GetUserMetrics() (*models.UserMetrics, error) {
	var metrics models.UserMetrics
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM users
		WHERE MONTH(created_at) = MONTH(CURDATE()) AND YEAR(created_at) = YEAR(CURDATE())`,
	).Scan(&metrics.NewUsersMonth)
	if err != nil {
		return nil, err
	}
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT user_id FROM form_submissions
			GROUP BY user_id
			HAVING COUNT(*) > 1
		) t`,
	).Scan(&metrics.ReturningUsers)
	if err != nil {
		return nil, err
	}
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT u.id)
		FROM users u
		WHERE EXISTS (SELECT 1 FROM form_submissions fs WHERE fs.user_id = u.id)
		AND NOT EXISTS (
			SELECT 1 FROM form_submissions fs
			WHERE fs.user_id = u.id
			AND fs.submitted_at >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
		)`,
	).Scan(&metrics.InactiveUsers)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}
