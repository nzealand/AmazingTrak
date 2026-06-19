package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type Corridor struct {
	ID                    int64
	Name                  string
	Slug                  string
	Region                string
	Description           string
	OnTimePercent         sql.NullFloat64
	ServiceQualitySummary string
	HeroTrainID           sql.NullInt64
	HeroMediaID           sql.NullInt64
	ThumbnailMediaID      sql.NullInt64
	IsActive              bool
	SortOrder             int
	CreatedAt             string
	UpdatedAt             string
	TrainCount            int
	MediaCount            int
	HeroImageURL          string
	ThumbnailURL          string
	ScheduleURL           string
	ConductorUserID       sql.NullInt64
	ConductorUsername     string
}

type Train struct {
	ID               int64
	CorridorID       int64
	TrainNumber      string
	DisplayName      string
	Slug             string
	Direction        string
	Notes            string
	HeroMediaID      sql.NullInt64
	ThumbnailMediaID sql.NullInt64
	MapMediaID       sql.NullInt64
	IsActive         bool
	SortOrder        int
	CreatedAt        string
	UpdatedAt        string
	CorridorName     string
	CorridorSlug     string
	HeroImageURL     string
	HeroLat          sql.NullFloat64
	HeroLon          sql.NullFloat64
	ThumbnailURL     string
	MapImageURL      string
	MediaCount       int
	PendingCount     int
	VideoCount       int
	HasBestVideo     bool
}

func (t *Train) HasHeroGeo() bool {
	return t.HeroLat.Valid && t.HeroLon.Valid
}

type Media struct {
	ID               int64
	TrainID          sql.NullInt64
	CorridorID       sql.NullInt64
	MediaType        string
	SourceType       string
	URL              string
	LocalPath        string
	OriginalFilename string
	StoredFilename   string
	Title            string
	Caption          string
	Tags             string
	SourceDomain     string
	Latitude         sql.NullFloat64
	Longitude        sql.NullFloat64
	LocationName     string
	LocationSource   string
	IsPublished      bool
	IsBest           bool
	AddedBy          string
	UserID           sql.NullInt64
	ContributorName  string
	CreatedAt        string
	UpdatedAt        string
}

// Contributor returns the display name to credit for this media item: the
// registered username if attributed, otherwise a label derived from AddedBy.
func (m *Media) Contributor() string {
	if m.ContributorName != "" {
		return m.ContributorName
	}
	switch m.AddedBy {
	case "admin":
		return "Admin"
	case "approved_suggestion", "public_approved":
		return "Anonymous"
	default:
		return m.AddedBy
	}
}

func (m *Media) PublicURL() string {
	if m.LocalPath != "" {
		return "/uploads/" + strings.TrimPrefix(m.LocalPath, "/")
	}
	return m.URL
}

func (m *Media) HasGeo() bool {
	return m.Latitude.Valid && m.Longitude.Valid
}

func (m *Media) LatStr() string {
	if !m.Latitude.Valid {
		return ""
	}
	return fmt.Sprintf("%.6f", m.Latitude.Float64)
}

func (m *Media) LonStr() string {
	if !m.Longitude.Valid {
		return ""
	}
	return fmt.Sprintf("%.6f", m.Longitude.Float64)
}

type Stop struct {
	ID           int64
	CorridorID   int64
	Name         string
	StationCode  string
	Slug         string
	Latitude     sql.NullFloat64
	Longitude    sql.NullFloat64
	SortOrder    int
	CorridorName string
	CorridorSlug string
}

type TrainStop struct {
	ID                 int64
	TrainID            int64
	StopID             int64
	SortOrder          int
	ScheduledArrival   string
	ScheduledDeparture string
	RunsWeekday        bool
	RunsWeekend        bool
	StopName           string
	StationCode        string
	StopSlug           string
}

type SitePreferences struct {
	ID                  int64
	DefaultTheme        string
	NotificationEmail   string
	RatePerMinute       int
	RatePerHour         int
	RatePerDay          int
	RegisterRatePerHour int
	RegisterRatePerDay  int
	CommentRatePerHour  int
	CommentRatePerDay   int
	SiteName            string
	FaviconPath         string
	AdminTheme          string
}

type StationTrain struct {
	TrainID            int64
	DisplayName        string
	Slug               string
	Direction          string
	IsActive           bool
	CorridorName       string
	CorridorSlug       string
	ScheduledArrival   string
	ScheduledDeparture string
}

type Suggestion struct {
	ID                 int64
	TrainID            int64
	URL                string
	Title              string
	Caption            string
	Tags               string
	MediaType          string
	SourceDomain       string
	Status             string
	SubmitterIPHash    string
	SubmitterUserAgent string
	RejectionReason    string
	CreatedAt          string
	ReviewedAt         string
	AutoApproved       bool
	IsSpam             bool
	TrainName          string
	TrainSlug          string
}

// StatusLabel renders the status for display, distinguishing spam-flagged
// pending/rejected rows without needing a separate status value.
func (s *Suggestion) StatusLabel() string {
	if s.IsSpam {
		switch s.Status {
		case "pending":
			return "pending - spam"
		case "rejected":
			return "rejected - spam"
		}
	}
	return s.Status
}

// RarityCount returns how many rarity tags this submission carries.
func (s *Suggestion) RarityCount() int {
	return rarityCount(s.Tags)
}

// Comment is a registered-user comment on a train, moderated the same way as
// suggestions: it enters as 'pending' and an admin approves or rejects it.
type Comment struct {
	ID              int64
	TrainID         int64
	UserID          int64
	Body            string
	Status          string
	SubmitterIPHash string
	RejectionReason string
	CreatedAt       string
	ReviewedAt      string
	Username        string
	TrainName       string
	TrainSlug       string
}

type AdminUser struct {
	ID                 int64
	Username           string
	PasswordHash       string
	MustChangePassword bool
	PermissionLevel    int
	CreatedAt          string
	LastLoginAt        string
}

// IsCentral reports whether this admin account is the immutable central admin
// (permission_level 0 = full access, cannot be deleted or have its role changed).
func (u *AdminUser) IsCentral() bool { return u.PermissionLevel == 0 }

// PermissionLabel renders a sub-admin's cumulative level for display.
func (u *AdminUser) PermissionLabel() string {
	switch u.PermissionLevel {
	case 0:
		return "Central admin (full access)"
	case 1:
		return "L1 — Suggestions"
	case 2:
		return "L2 — Suggestions, Comments"
	case 3:
		return "L3 — Suggestions, Comments, Trains"
	case 4:
		return "L4 — Suggestions, Comments, Trains, Corridors"
	case 5:
		return "L5 — Suggestions, Comments, Trains, Corridors, Settings"
	case 6:
		return "L6 — Suggestions, Comments, Trains, Corridors, Settings, Users"
	}
	return fmt.Sprintf("L%d", u.PermissionLevel)
}

type User struct {
	ID              int64
	Username        string
	Email           string
	PasswordHash    string
	Status          string
	EmailConfirmed  bool
	ConfirmToken    string
	CreatedAt       string
	LastLoginAt     sql.NullString
	SubmissionCount int
	IsSpammer       bool
}

// StatusLabel renders a user's account status for display.
func (u *User) StatusLabel() string {
	if u.IsSpammer {
		return "Spammer"
	}
	switch u.Status {
	case "pending":
		return "Pending"
	case "confirmed":
		return "Confirmed (email verified)"
	case "approved":
		return "Approved"
	case "auto_approved":
		return "Auto-approved"
	}
	return u.Status
}

func (u *User) IsApproved() bool {
	return u.Status == "approved" || u.Status == "auto_approved"
}

type Session struct {
	ID          string
	AdminUserID int64
	CSRFToken   string
	ExpiresAt   string
}

// ----- Corridor queries -----

func allCorridors(db *sql.DB, activeOnly bool) ([]Corridor, error) {
	q := `SELECT c.id, c.name, c.slug, COALESCE(c.region,''), COALESCE(c.description,''),
		c.on_time_percent, COALESCE(c.service_quality_summary,''),
		c.hero_train_id, c.hero_media_id, c.thumbnail_media_id,
		c.is_active, c.sort_order, c.created_at, c.updated_at,
		COUNT(DISTINCT t.id),
		(SELECT COUNT(*) FROM media WHERE corridor_id=c.id),
		COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=c.hero_media_id), ''),
		COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=c.thumbnail_media_id), ''),
		COALESCE(c.schedule_url,''),
		c.conductor_user_id, COALESCE((SELECT username FROM users WHERE id=c.conductor_user_id), '')
		FROM corridors c
		LEFT JOIN trains t ON t.corridor_id = c.id`
	if activeOnly {
		q += ` WHERE c.is_active = 1`
	}
	q += ` GROUP BY c.id ORDER BY c.sort_order, c.name`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCorridors(rows)
}

func corridorBySlug(db *sql.DB, slug string) (Corridor, error) {
	var c Corridor
	err := db.QueryRow(`SELECT c.id, c.name, c.slug, COALESCE(c.region,''), COALESCE(c.description,''),
		c.on_time_percent, COALESCE(c.service_quality_summary,''),
		c.hero_train_id, c.hero_media_id, c.thumbnail_media_id,
		c.is_active, c.sort_order, c.created_at, c.updated_at,
		(SELECT COUNT(*) FROM trains WHERE corridor_id=c.id),
		(SELECT COUNT(*) FROM media WHERE corridor_id=c.id),
		COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=c.hero_media_id), ''),
		COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=c.thumbnail_media_id), ''),
		COALESCE(c.schedule_url,''),
		c.conductor_user_id, COALESCE((SELECT username FROM users WHERE id=c.conductor_user_id), '')
		FROM corridors c WHERE c.slug=?`, slug).
		Scan(&c.ID, &c.Name, &c.Slug, &c.Region, &c.Description,
			&c.OnTimePercent, &c.ServiceQualitySummary,
			&c.HeroTrainID, &c.HeroMediaID, &c.ThumbnailMediaID,
			&c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt,
			&c.TrainCount, &c.MediaCount, &c.HeroImageURL, &c.ThumbnailURL, &c.ScheduleURL,
			&c.ConductorUserID, &c.ConductorUsername)
	return c, err
}

func corridorByID(db *sql.DB, id int64) (Corridor, error) {
	var c Corridor
	err := db.QueryRow(`SELECT c.id, c.name, c.slug, COALESCE(c.region,''), COALESCE(c.description,''),
		c.on_time_percent, COALESCE(c.service_quality_summary,''),
		c.hero_train_id, c.hero_media_id, c.thumbnail_media_id,
		c.is_active, c.sort_order, c.created_at, c.updated_at,
		(SELECT COUNT(*) FROM trains WHERE corridor_id=c.id),
		(SELECT COUNT(*) FROM media WHERE corridor_id=c.id),
		COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=c.hero_media_id), ''),
		COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=c.thumbnail_media_id), ''),
		COALESCE(c.schedule_url,''),
		c.conductor_user_id, COALESCE((SELECT username FROM users WHERE id=c.conductor_user_id), '')
		FROM corridors c WHERE c.id=?`, id).
		Scan(&c.ID, &c.Name, &c.Slug, &c.Region, &c.Description,
			&c.OnTimePercent, &c.ServiceQualitySummary,
			&c.HeroTrainID, &c.HeroMediaID, &c.ThumbnailMediaID,
			&c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt,
			&c.TrainCount, &c.MediaCount, &c.HeroImageURL, &c.ThumbnailURL, &c.ScheduleURL,
			&c.ConductorUserID, &c.ConductorUsername)
	return c, err
}

func scanCorridors(rows *sql.Rows) ([]Corridor, error) {
	var out []Corridor
	for rows.Next() {
		var c Corridor
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Region, &c.Description,
			&c.OnTimePercent, &c.ServiceQualitySummary,
			&c.HeroTrainID, &c.HeroMediaID, &c.ThumbnailMediaID,
			&c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt,
			&c.TrainCount, &c.MediaCount, &c.HeroImageURL, &c.ThumbnailURL, &c.ScheduleURL,
			&c.ConductorUserID, &c.ConductorUsername); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ----- Conductor queries -----

// ConductorRequest is a user's pending/decided request to maintain a corridor.
type ConductorRequest struct {
	ID           int64
	CorridorID   int64
	UserID       int64
	Status       string
	Message      string
	CreatedAt    string
	ReviewedAt   sql.NullString
	CorridorName string
	CorridorSlug string
	Username     string
}

// corridorsConductedBy returns the corridors a user is the conductor of.
func corridorsConductedBy(db *sql.DB, userID int64) ([]Corridor, error) {
	q := `SELECT c.id, c.name, c.slug, COALESCE(c.region,''), COALESCE(c.description,''),
		c.on_time_percent, COALESCE(c.service_quality_summary,''),
		c.hero_train_id, c.hero_media_id, c.thumbnail_media_id,
		c.is_active, c.sort_order, c.created_at, c.updated_at,
		(SELECT COUNT(*) FROM trains WHERE corridor_id=c.id),
		(SELECT COUNT(*) FROM media WHERE corridor_id=c.id),
		COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=c.hero_media_id), ''),
		COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=c.thumbnail_media_id), ''),
		COALESCE(c.schedule_url,''),
		c.conductor_user_id, COALESCE((SELECT username FROM users WHERE id=c.conductor_user_id), '')
		FROM corridors c WHERE c.conductor_user_id=? ORDER BY c.sort_order, c.name`
	rows, err := db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCorridors(rows)
}

// isConductorOf reports whether userID is the conductor of corridorID.
func isConductorOf(db *sql.DB, userID, corridorID int64) (bool, error) {
	var n int
	err := db.QueryRow(`SELECT COUNT(*) FROM corridors WHERE id=? AND conductor_user_id=?`, corridorID, userID).Scan(&n)
	return n > 0, err
}

func setCorridorConductor(db *sql.DB, corridorID, userID int64) error {
	_, err := db.Exec(`UPDATE corridors SET conductor_user_id=? WHERE id=?`, userID, corridorID)
	return err
}

func clearCorridorConductor(db *sql.DB, corridorID int64) error {
	_, err := db.Exec(`UPDATE corridors SET conductor_user_id=NULL WHERE id=?`, corridorID)
	return err
}

// pendingConductorRequest reports whether the user already has a pending request
// for the given corridor.
func pendingConductorRequest(db *sql.DB, corridorID, userID int64) (bool, error) {
	var n int
	err := db.QueryRow(`SELECT COUNT(*) FROM conductor_requests WHERE corridor_id=? AND user_id=? AND status='pending'`, corridorID, userID).Scan(&n)
	return n > 0, err
}

func createConductorRequest(db *sql.DB, corridorID, userID int64, message string) error {
	_, err := db.Exec(`INSERT INTO conductor_requests (corridor_id, user_id, message) VALUES (?, ?, ?)`, corridorID, userID, message)
	return err
}

func conductorRequestByID(db *sql.DB, id int64) (ConductorRequest, error) {
	var cr ConductorRequest
	err := db.QueryRow(`SELECT cr.id, cr.corridor_id, cr.user_id, cr.status, cr.message, cr.created_at, cr.reviewed_at,
		c.name, c.slug, u.username
		FROM conductor_requests cr
		JOIN corridors c ON c.id=cr.corridor_id
		JOIN users u ON u.id=cr.user_id
		WHERE cr.id=?`, id).
		Scan(&cr.ID, &cr.CorridorID, &cr.UserID, &cr.Status, &cr.Message, &cr.CreatedAt, &cr.ReviewedAt,
			&cr.CorridorName, &cr.CorridorSlug, &cr.Username)
	return cr, err
}

// conductorRequestsByStatus lists requests with the given status, newest first.
func conductorRequestsByStatus(db *sql.DB, status string) ([]ConductorRequest, error) {
	rows, err := db.Query(`SELECT cr.id, cr.corridor_id, cr.user_id, cr.status, cr.message, cr.created_at, cr.reviewed_at,
		c.name, c.slug, u.username
		FROM conductor_requests cr
		JOIN corridors c ON c.id=cr.corridor_id
		JOIN users u ON u.id=cr.user_id
		WHERE cr.status=? ORDER BY cr.created_at DESC`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ConductorRequest
	for rows.Next() {
		var cr ConductorRequest
		if err := rows.Scan(&cr.ID, &cr.CorridorID, &cr.UserID, &cr.Status, &cr.Message, &cr.CreatedAt, &cr.ReviewedAt,
			&cr.CorridorName, &cr.CorridorSlug, &cr.Username); err != nil {
			return nil, err
		}
		out = append(out, cr)
	}
	return out, rows.Err()
}

// decideConductorRequest sets a request's status and review metadata.
func decideConductorRequest(db *sql.DB, id int64, status string, adminUserID int64) error {
	_, err := db.Exec(`UPDATE conductor_requests SET status=?, reviewed_at=CURRENT_TIMESTAMP, reviewed_by=? WHERE id=?`,
		status, adminUserID, id)
	return err
}

// rejectOtherPendingRequests rejects all other pending requests for a corridor
// (used when one request is approved or a conductor is assigned directly).
func rejectOtherPendingRequests(db *sql.DB, corridorID, exceptID int64, adminUserID int64) error {
	_, err := db.Exec(`UPDATE conductor_requests SET status='rejected', reviewed_at=CURRENT_TIMESTAMP, reviewed_by=?
		WHERE corridor_id=? AND status='pending' AND id!=?`, adminUserID, corridorID, exceptID)
	return err
}

// ----- Train queries -----

const trainSelectBase = `SELECT t.id, t.corridor_id, t.train_number, t.display_name, t.slug,
	COALESCE(t.direction,''), COALESCE(t.notes,''),
	t.hero_media_id, t.thumbnail_media_id, t.map_media_id,
	t.is_active, t.sort_order, t.created_at, t.updated_at,
	c.name, c.slug,
	COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=t.hero_media_id), ''),
	(SELECT latitude FROM media WHERE id=t.hero_media_id),
	(SELECT longitude FROM media WHERE id=t.hero_media_id),
	COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=t.thumbnail_media_id), ''),
	COALESCE((SELECT CASE WHEN local_path!='' AND local_path IS NOT NULL THEN '/uploads/'||local_path ELSE url END FROM media WHERE id=t.map_media_id), ''),
	(SELECT COUNT(*) FROM media WHERE train_id=t.id AND is_published=1),
	(SELECT COUNT(*) FROM suggestions WHERE train_id=t.id AND status='pending'),
	(SELECT COUNT(*) FROM media WHERE train_id=t.id AND is_published=1 AND media_type='video'),
	COALESCE((SELECT 1 FROM media WHERE train_id=t.id AND is_best=1 LIMIT 1), 0)
	FROM trains t JOIN corridors c ON c.id=t.corridor_id`

func scanTrains(rows *sql.Rows) ([]Train, error) {
	var out []Train
	for rows.Next() {
		var t Train
		if err := rows.Scan(&t.ID, &t.CorridorID, &t.TrainNumber, &t.DisplayName, &t.Slug,
			&t.Direction, &t.Notes,
			&t.HeroMediaID, &t.ThumbnailMediaID, &t.MapMediaID,
			&t.IsActive, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt,
			&t.CorridorName, &t.CorridorSlug,
			&t.HeroImageURL, &t.HeroLat, &t.HeroLon,
			&t.ThumbnailURL, &t.MapImageURL,
			&t.MediaCount, &t.PendingCount, &t.VideoCount, &t.HasBestVideo); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func trainsByCorridorID(db *sql.DB, corridorID int64, activeOnly bool) ([]Train, error) {
	q := trainSelectBase + ` WHERE t.corridor_id=?`
	if activeOnly {
		q += ` AND t.is_active=1`
	}
	q += ` ORDER BY t.sort_order, t.train_number`
	rows, err := db.Query(q, corridorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTrains(rows)
}

func allTrains(db *sql.DB) ([]Train, error) {
	q := trainSelectBase + ` ORDER BY c.sort_order, t.sort_order, t.train_number`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTrains(rows)
}

func trainBySlug(db *sql.DB, slug string) (Train, error) {
	var t Train
	err := db.QueryRow(trainSelectBase+` WHERE t.slug=?`, slug).
		Scan(&t.ID, &t.CorridorID, &t.TrainNumber, &t.DisplayName, &t.Slug,
			&t.Direction, &t.Notes,
			&t.HeroMediaID, &t.ThumbnailMediaID, &t.MapMediaID,
			&t.IsActive, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt,
			&t.CorridorName, &t.CorridorSlug,
			&t.HeroImageURL, &t.HeroLat, &t.HeroLon,
			&t.ThumbnailURL, &t.MapImageURL,
			&t.MediaCount, &t.PendingCount, &t.VideoCount, &t.HasBestVideo)
	return t, err
}

func trainByID(db *sql.DB, id int64) (Train, error) {
	var t Train
	err := db.QueryRow(trainSelectBase+` WHERE t.id=?`, id).
		Scan(&t.ID, &t.CorridorID, &t.TrainNumber, &t.DisplayName, &t.Slug,
			&t.Direction, &t.Notes,
			&t.HeroMediaID, &t.ThumbnailMediaID, &t.MapMediaID,
			&t.IsActive, &t.SortOrder, &t.CreatedAt, &t.UpdatedAt,
			&t.CorridorName, &t.CorridorSlug,
			&t.HeroImageURL, &t.HeroLat, &t.HeroLon,
			&t.ThumbnailURL, &t.MapImageURL,
			&t.MediaCount, &t.PendingCount, &t.VideoCount, &t.HasBestVideo)
	return t, err
}

// ----- Media queries -----

func mediaByID(db *sql.DB, id int64) (Media, error) {
	var m Media
	err := db.QueryRow(`SELECT m.id, m.train_id, m.corridor_id, m.media_type, m.source_type,
		COALESCE(m.url,''), COALESCE(m.local_path,''), COALESCE(m.original_filename,''), COALESCE(m.stored_filename,''),
		COALESCE(m.title,''), COALESCE(m.caption,''), COALESCE(m.tags,''), COALESCE(m.source_domain,''),
		m.latitude, m.longitude, COALESCE(m.location_name,''), COALESCE(m.location_source,'unknown'),
		m.is_published, m.is_best, m.added_by, m.user_id, COALESCE(u.username,''), m.created_at, m.updated_at
		FROM media m LEFT JOIN users u ON u.id=m.user_id WHERE m.id=?`, id).
		Scan(&m.ID, &m.TrainID, &m.CorridorID, &m.MediaType, &m.SourceType,
			&m.URL, &m.LocalPath, &m.OriginalFilename, &m.StoredFilename,
			&m.Title, &m.Caption, &m.Tags, &m.SourceDomain,
			&m.Latitude, &m.Longitude, &m.LocationName, &m.LocationSource,
			&m.IsPublished, &m.IsBest, &m.AddedBy, &m.UserID, &m.ContributorName, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func mediaByTrainID(db *sql.DB, trainID int64, publishedOnly bool) ([]Media, error) {
	q := `SELECT m.id, m.train_id, m.corridor_id, m.media_type, m.source_type,
		COALESCE(m.url,''), COALESCE(m.local_path,''), COALESCE(m.original_filename,''), COALESCE(m.stored_filename,''),
		COALESCE(m.title,''), COALESCE(m.caption,''), COALESCE(m.tags,''), COALESCE(m.source_domain,''),
		m.latitude, m.longitude, COALESCE(m.location_name,''), COALESCE(m.location_source,'unknown'),
		m.is_published, m.is_best, m.added_by, m.user_id, COALESCE(u.username,''), m.created_at, m.updated_at
		FROM media m LEFT JOIN users u ON u.id=m.user_id WHERE m.train_id=?`
	if publishedOnly {
		q += ` AND m.is_published=1`
	}
	q += ` ORDER BY m.created_at DESC`
	rows, err := db.Query(q, trainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMedia(rows)
}

func mediaByCorridorID(db *sql.DB, corridorID int64, publishedOnly bool) ([]Media, error) {
	q := `SELECT m.id, m.train_id, m.corridor_id, m.media_type, m.source_type,
		COALESCE(m.url,''), COALESCE(m.local_path,''), COALESCE(m.original_filename,''), COALESCE(m.stored_filename,''),
		COALESCE(m.title,''), COALESCE(m.caption,''), COALESCE(m.tags,''), COALESCE(m.source_domain,''),
		m.latitude, m.longitude, COALESCE(m.location_name,''), COALESCE(m.location_source,'unknown'),
		m.is_published, m.is_best, m.added_by, m.user_id, COALESCE(u.username,''), m.created_at, m.updated_at
		FROM media m LEFT JOIN users u ON u.id=m.user_id WHERE m.corridor_id=?`
	if publishedOnly {
		q += ` AND m.is_published=1`
	}
	q += ` ORDER BY m.created_at DESC`
	rows, err := db.Query(q, corridorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMedia(rows)
}

func scanMedia(rows *sql.Rows) ([]Media, error) {
	var out []Media
	for rows.Next() {
		var m Media
		if err := rows.Scan(&m.ID, &m.TrainID, &m.CorridorID, &m.MediaType, &m.SourceType,
			&m.URL, &m.LocalPath, &m.OriginalFilename, &m.StoredFilename,
			&m.Title, &m.Caption, &m.Tags, &m.SourceDomain,
			&m.Latitude, &m.Longitude, &m.LocationName, &m.LocationSource,
			&m.IsPublished, &m.IsBest, &m.AddedBy, &m.UserID, &m.ContributorName, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// ----- Stop queries -----

func stopsByCorridorID(db *sql.DB, corridorID int64) ([]Stop, error) {
	rows, err := db.Query(`SELECT s.id, s.corridor_id, s.name, COALESCE(s.station_code,''), COALESCE(s.slug,''),
		s.latitude, s.longitude, s.sort_order, c.name, c.slug
		FROM stops s JOIN corridors c ON c.id=s.corridor_id
		WHERE s.corridor_id=? ORDER BY s.sort_order`, corridorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Stop
	for rows.Next() {
		var s Stop
		if err := rows.Scan(&s.ID, &s.CorridorID, &s.Name, &s.StationCode, &s.Slug,
			&s.Latitude, &s.Longitude, &s.SortOrder, &s.CorridorName, &s.CorridorSlug); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func stopsByTrainID(db *sql.DB, trainID int64) ([]TrainStop, error) {
	rows, err := db.Query(`SELECT ts.id, ts.train_id, ts.stop_id, ts.sort_order,
		COALESCE(ts.scheduled_arrival,''), COALESCE(ts.scheduled_departure,''),
		ts.runs_weekday, ts.runs_weekend,
		s.name, COALESCE(s.station_code,''), COALESCE(s.slug,'')
		FROM train_stops ts JOIN stops s ON s.id=ts.stop_id
		WHERE ts.train_id=? ORDER BY ts.sort_order`, trainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TrainStop
	for rows.Next() {
		var ts TrainStop
		if err := rows.Scan(&ts.ID, &ts.TrainID, &ts.StopID, &ts.SortOrder,
			&ts.ScheduledArrival, &ts.ScheduledDeparture,
			&ts.RunsWeekday, &ts.RunsWeekend,
			&ts.StopName, &ts.StationCode, &ts.StopSlug); err != nil {
			return nil, err
		}
		out = append(out, ts)
	}
	return out, rows.Err()
}

func stopsBySlug(db *sql.DB, slug string) ([]Stop, error) {
	rows, err := db.Query(`SELECT s.id, s.corridor_id, s.name, COALESCE(s.station_code,''), COALESCE(s.slug,''),
		s.latitude, s.longitude, s.sort_order, c.name, c.slug
		FROM stops s JOIN corridors c ON c.id=s.corridor_id
		WHERE s.slug=? ORDER BY c.sort_order, s.sort_order`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Stop
	for rows.Next() {
		var s Stop
		if err := rows.Scan(&s.ID, &s.CorridorID, &s.Name, &s.StationCode, &s.Slug,
			&s.Latitude, &s.Longitude, &s.SortOrder, &s.CorridorName, &s.CorridorSlug); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func trainsByStopID(db *sql.DB, stopID int64) ([]StationTrain, error) {
	rows, err := db.Query(`SELECT t.id, t.display_name, t.slug, COALESCE(t.direction,''), t.is_active,
		c.name, c.slug,
		COALESCE(ts.scheduled_arrival,''), COALESCE(ts.scheduled_departure,'')
		FROM train_stops ts
		JOIN trains t ON t.id=ts.train_id
		JOIN corridors c ON c.id=t.corridor_id
		WHERE ts.stop_id=? AND t.is_active=1
		ORDER BY ts.scheduled_departure, ts.scheduled_arrival, t.train_number`, stopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []StationTrain
	for rows.Next() {
		var st StationTrain
		if err := rows.Scan(&st.TrainID, &st.DisplayName, &st.Slug, &st.Direction, &st.IsActive,
			&st.CorridorName, &st.CorridorSlug, &st.ScheduledArrival, &st.ScheduledDeparture); err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}

func getSitePrefs(db *sql.DB) (SitePreferences, error) {
	var p SitePreferences
	err := db.QueryRow(`SELECT id, default_theme, COALESCE(notification_email,''), COALESCE(rate_per_minute,1), COALESCE(rate_per_hour,5), COALESCE(rate_per_day,20), COALESCE(register_rate_per_hour,5), COALESCE(register_rate_per_day,20), COALESCE(comment_rate_per_hour,10), COALESCE(comment_rate_per_day,50), COALESCE(site_name,'AmazingTrak'), COALESCE(favicon_path,''), COALESCE(admin_theme,'default') FROM site_preferences WHERE id=1`).
		Scan(&p.ID, &p.DefaultTheme, &p.NotificationEmail, &p.RatePerMinute, &p.RatePerHour, &p.RatePerDay, &p.RegisterRatePerHour, &p.RegisterRatePerDay, &p.CommentRatePerHour, &p.CommentRatePerDay, &p.SiteName, &p.FaviconPath, &p.AdminTheme)
	if err == sql.ErrNoRows {
		return SitePreferences{DefaultTheme: "auto", RatePerMinute: 1, RatePerHour: 5, RatePerDay: 20, RegisterRatePerHour: 5, RegisterRatePerDay: 20, CommentRatePerHour: 10, CommentRatePerDay: 50, SiteName: "AmazingTrak", AdminTheme: "default"}, nil
	}
	return p, err
}

func getAdminUser(db *sql.DB, id int64) (AdminUser, error) {
	var u AdminUser
	err := db.QueryRow(`SELECT id, username, password_hash, must_change_password, permission_level FROM admin_users WHERE id=?`, id).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.MustChangePassword, &u.PermissionLevel)
	return u, err
}

// allAdminUsers returns every admin account (central admin + sub-admins),
// ordered with the central admin first.
func allAdminUsers(db *sql.DB) ([]AdminUser, error) {
	rows, err := db.Query(`SELECT id, username, password_hash, must_change_password, permission_level, created_at, COALESCE(last_login_at,'')
		FROM admin_users ORDER BY permission_level, username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AdminUser
	for rows.Next() {
		var u AdminUser
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.MustChangePassword, &u.PermissionLevel, &u.CreatedAt, &u.LastLoginAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// ----- User (registered public account) queries -----

const userSelectCols = `id, username, COALESCE(email,''), password_hash, status, email_confirmed, COALESCE(confirm_token,''), created_at, last_login_at, is_spammer`

func scanUser(row interface{ Scan(...interface{}) error }) (User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Status, &u.EmailConfirmed, &u.ConfirmToken, &u.CreatedAt, &u.LastLoginAt, &u.IsSpammer)
	return u, err
}

func userByID(db *sql.DB, id int64) (User, error) {
	return scanUser(db.QueryRow(`SELECT `+userSelectCols+` FROM users WHERE id=?`, id))
}

func userByUsername(db *sql.DB, username string) (User, error) {
	return scanUser(db.QueryRow(`SELECT `+userSelectCols+` FROM users WHERE username=?`, username))
}

func userByConfirmToken(db *sql.DB, token string) (User, error) {
	return scanUser(db.QueryRow(`SELECT `+userSelectCols+` FROM users WHERE confirm_token=? AND confirm_token != ''`, token))
}

// allUsers returns every registered user with a count of their submissions.
func allUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT ` + userSelectCols + `,
		(SELECT COUNT(*) FROM suggestions WHERE user_id = users.id)
		FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Status, &u.EmailConfirmed, &u.ConfirmToken, &u.CreatedAt, &u.LastLoginAt, &u.IsSpammer, &u.SubmissionCount); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// submissionsByUserID returns all suggestions tied to a registered user.
func submissionsByUserID(db *sql.DB, userID int64) ([]Suggestion, error) {
	q := `SELECT s.id, s.train_id, s.url, COALESCE(s.title,''), COALESCE(s.caption,''), COALESCE(s.tags,''), s.media_type, COALESCE(s.source_domain,''),
		s.status, COALESCE(s.submitter_ip_hash,''), COALESCE(s.submitter_user_agent,''),
		COALESCE(s.rejection_reason,''), s.created_at, COALESCE(s.reviewed_at,''), s.auto_approved, s.is_spam,
		t.display_name, t.slug
		FROM suggestions s JOIN trains t ON t.id=s.train_id WHERE s.user_id=? ORDER BY s.created_at DESC`
	rows, err := db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSuggestions(rows)
}

// ----- Suggestion queries -----

func allSuggestions(db *sql.DB, status string) ([]Suggestion, error) {
	q := `SELECT s.id, s.train_id, s.url, COALESCE(s.title,''), COALESCE(s.caption,''), COALESCE(s.tags,''), s.media_type, COALESCE(s.source_domain,''),
		s.status, COALESCE(s.submitter_ip_hash,''), COALESCE(s.submitter_user_agent,''),
		COALESCE(s.rejection_reason,''), s.created_at, COALESCE(s.reviewed_at,''), s.auto_approved, s.is_spam,
		t.display_name, t.slug
		FROM suggestions s JOIN trains t ON t.id=s.train_id`
	if status != "" {
		q += ` WHERE s.status=?`
	}
	q += ` ORDER BY s.created_at DESC`
	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = db.Query(q, status)
	} else {
		rows, err = db.Query(q)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSuggestions(rows)
}

func suggestionsByTrainID(db *sql.DB, trainID int64, status string) ([]Suggestion, error) {
	q := `SELECT s.id, s.train_id, s.url, COALESCE(s.title,''), COALESCE(s.caption,''), COALESCE(s.tags,''), s.media_type, COALESCE(s.source_domain,''),
		s.status, COALESCE(s.submitter_ip_hash,''), COALESCE(s.submitter_user_agent,''),
		COALESCE(s.rejection_reason,''), s.created_at, COALESCE(s.reviewed_at,''), s.auto_approved, s.is_spam,
		t.display_name, t.slug
		FROM suggestions s JOIN trains t ON t.id=s.train_id WHERE s.train_id=?`
	args := []interface{}{trainID}
	if status != "" {
		q += ` AND s.status=?`
		args = append(args, status)
	}
	q += ` ORDER BY s.created_at DESC`
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSuggestions(rows)
}

func suggestionByID(db *sql.DB, id int64) (Suggestion, error) {
	var s Suggestion
	err := db.QueryRow(`SELECT s.id, s.train_id, s.url, COALESCE(s.title,''), COALESCE(s.caption,''), COALESCE(s.tags,''), s.media_type, COALESCE(s.source_domain,''),
		s.status, COALESCE(s.submitter_ip_hash,''), COALESCE(s.submitter_user_agent,''),
		COALESCE(s.rejection_reason,''), s.created_at, COALESCE(s.reviewed_at,''), s.auto_approved, s.is_spam,
		t.display_name, t.slug
		FROM suggestions s JOIN trains t ON t.id=s.train_id WHERE s.id=?`, id).
		Scan(&s.ID, &s.TrainID, &s.URL, &s.Title, &s.Caption, &s.Tags, &s.MediaType, &s.SourceDomain,
			&s.Status, &s.SubmitterIPHash, &s.SubmitterUserAgent,
			&s.RejectionReason, &s.CreatedAt, &s.ReviewedAt, &s.AutoApproved, &s.IsSpam,
			&s.TrainName, &s.TrainSlug)
	return s, err
}

func scanSuggestions(rows *sql.Rows) ([]Suggestion, error) {
	var out []Suggestion
	for rows.Next() {
		var s Suggestion
		if err := rows.Scan(&s.ID, &s.TrainID, &s.URL, &s.Title, &s.Caption, &s.Tags, &s.MediaType, &s.SourceDomain,
			&s.Status, &s.SubmitterIPHash, &s.SubmitterUserAgent,
			&s.RejectionReason, &s.CreatedAt, &s.ReviewedAt, &s.AutoApproved, &s.IsSpam,
			&s.TrainName, &s.TrainSlug); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// ----- Comment queries -----

const commentSelectBase = `SELECT c.id, c.train_id, c.user_id, c.body, c.status,
	COALESCE(c.submitter_ip_hash,''), COALESCE(c.rejection_reason,''),
	c.created_at, COALESCE(c.reviewed_at,''),
	u.username, t.display_name, t.slug
	FROM comments c
	JOIN users u ON u.id=c.user_id
	JOIN trains t ON t.id=c.train_id`

func scanComments(rows *sql.Rows) ([]Comment, error) {
	var out []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.TrainID, &c.UserID, &c.Body, &c.Status,
			&c.SubmitterIPHash, &c.RejectionReason,
			&c.CreatedAt, &c.ReviewedAt,
			&c.Username, &c.TrainName, &c.TrainSlug); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// commentsByTrainID returns a train's comments, optionally filtered by status,
// newest first.
func commentsByTrainID(db *sql.DB, trainID int64, status string) ([]Comment, error) {
	q := commentSelectBase + ` WHERE c.train_id=?`
	args := []interface{}{trainID}
	if status != "" {
		q += ` AND c.status=?`
		args = append(args, status)
	}
	q += ` ORDER BY c.created_at DESC`
	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanComments(rows)
}

// allComments returns every comment, optionally filtered by status, newest first.
func allComments(db *sql.DB, status string) ([]Comment, error) {
	q := commentSelectBase
	var rows *sql.Rows
	var err error
	if status != "" {
		q += ` WHERE c.status=? ORDER BY c.created_at DESC`
		rows, err = db.Query(q, status)
	} else {
		q += ` ORDER BY c.created_at DESC`
		rows, err = db.Query(q)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanComments(rows)
}

// commentsByUserID returns all approved comments by a user, newest first.
func commentsByUserID(db *sql.DB, userID int64) ([]Comment, error) {
	q := commentSelectBase + ` WHERE c.user_id=? AND c.status='approved' ORDER BY c.created_at DESC`
	rows, err := db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanComments(rows)
}

func commentByID(db *sql.DB, id int64) (Comment, error) {
	var c Comment
	err := db.QueryRow(commentSelectBase+` WHERE c.id=?`, id).
		Scan(&c.ID, &c.TrainID, &c.UserID, &c.Body, &c.Status,
			&c.SubmitterIPHash, &c.RejectionReason,
			&c.CreatedAt, &c.ReviewedAt,
			&c.Username, &c.TrainName, &c.TrainSlug)
	return c, err
}

// ----- Util -----

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			b.WriteRune('-')
		}
	}
	result := b.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}
