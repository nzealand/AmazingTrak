package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed templates static
var embedFS embed.FS

var staticVersion string

var (
	siteNameMu  sync.RWMutex
	siteNameVal = "AmazingTrak"

	faviconMu      sync.RWMutex
	faviconPathVal string
	faviconVerVal  string

	adminThemeMu   sync.RWMutex
	adminThemeVal  = "default"
	adminCompactMu sync.RWMutex
	adminCompactOn bool
)

type themeSpec struct {
	navVars string // CSS vars for the admin nav bar (--anb, --anf, --anl)
	// publicLight / publicDark are full sets of public-site CSS variable
	// overrides for the two color schemes. Each is injected verbatim inside
	// a :root{} block (light) or a dark-mode selector (dark).
	publicLight string
	publicDark  string
}

// publicVars builds the CSS variables that control the public site palette.
// accent = the main hue; accentFg = text on that hue; navBg = site nav bg;
// navFg = site nav text; tagBg/tagFg = video-tag pill colors.
func publicVars(accent, accentFg, navBg, navFg, tagBg, tagFg string) string {
	return "--accent:" + accent + ";--accent-fg:" + accentFg +
		";--nav-bg:" + navBg + ";--nav-fg:" + navFg +
		";--tag-bg:" + tagBg + ";--tag-fg:" + tagFg
}

// siteThemes defines all available color themes.
var siteThemes = map[string]themeSpec{
	// default — dark navy admin nav, classic blue accent
	"default": {
		navVars:     "--anb:#1e293b;--anf:#f1f5f9;--anl:#94a3b8",
		publicLight: publicVars("#1d4ed8", "#fff", "#1e3a5f", "#f0f4ff", "#dbeafe", "#1e40af"),
		publicDark:  publicVars("#58a6ff", "#0d1117", "#010409", "#e6edf3", "#1c2a3f", "#58a6ff"),
	},
	// grey — charcoal nav, neutral grey accent
	"grey": {
		navVars:     "--anb:#374151;--anf:#f9fafb;--anl:#9ca3af",
		publicLight: publicVars("#374151", "#fff", "#1f2937", "#f9fafb", "#f3f4f6", "#374151"),
		publicDark:  publicVars("#9ca3af", "#111827", "#111827", "#f9fafb", "#1f2937", "#9ca3af"),
	},
	// white — white/light admin nav; in dark mode everything inverts to near-black
	"white": {
		navVars:     "--anb:#ffffff;--anf:#111827;--anl:#6b7280;--an-border:1px solid #e5e7eb",
		publicLight: publicVars("#111827", "#fff", "#ffffff", "#111827", "#f3f4f6", "#111827"),
		publicDark:  publicVars("#f9fafb", "#111827", "#0a0a0a", "#f9fafb", "#1a1a1a", "#f9fafb"),
	},
	// blue — royal blue
	"blue": {
		navVars:     "--anb:#1e40af;--anf:#eff6ff;--anl:#bfdbfe",
		publicLight: publicVars("#1e40af", "#fff", "#1e3a8a", "#eff6ff", "#dbeafe", "#1e40af"),
		publicDark:  publicVars("#60a5fa", "#0d1117", "#0c1a4d", "#eff6ff", "#1c2a3f", "#60a5fa"),
	},
	// green — forest green
	"green": {
		navVars:     "--anb:#166534;--anf:#f0fdf4;--anl:#86efac",
		publicLight: publicVars("#16a34a", "#fff", "#14532d", "#f0fdf4", "#dcfce7", "#166534"),
		publicDark:  publicVars("#4ade80", "#052e16", "#052e16", "#f0fdf4", "#052e16", "#4ade80"),
	},
	// purple — deep plum
	"purple": {
		navVars:     "--anb:#581c87;--anf:#faf5ff;--anl:#d8b4fe",
		publicLight: publicVars("#7c3aed", "#fff", "#4c1d95", "#faf5ff", "#ede9fe", "#5b21b6"),
		publicDark:  publicVars("#a78bfa", "#1e0050", "#1e0050", "#faf5ff", "#2e1065", "#a78bfa"),
	},
	// red — crimson
	"red": {
		navVars:     "--anb:#991b1b;--anf:#fff1f2;--anl:#fca5a5",
		publicLight: publicVars("#dc2626", "#fff", "#7f1d1d", "#fff1f2", "#fee2e2", "#991b1b"),
		publicDark:  publicVars("#f87171", "#1c0202", "#1c0202", "#fff1f2", "#2d0a0a", "#f87171"),
	},
	// teal — deep teal
	"teal": {
		navVars:     "--anb:#0f766e;--anf:#f0fdfa;--anl:#5eead4",
		publicLight: publicVars("#0d9488", "#fff", "#134e4a", "#f0fdfa", "#ccfbf1", "#0f766e"),
		publicDark:  publicVars("#2dd4bf", "#001a18", "#001a18", "#f0fdfa", "#001a18", "#2dd4bf"),
	},
}

// adminThemes is an alias kept for the settings handler validation map.
var adminThemeNames = []string{"default", "grey", "white", "blue", "green", "purple", "red", "teal"}

func getSiteName() string {
	siteNameMu.RLock()
	defer siteNameMu.RUnlock()
	return siteNameVal
}

func setSiteName(s string) {
	siteNameMu.Lock()
	siteNameVal = s
	siteNameMu.Unlock()
}

func getAdminTheme() string {
	adminThemeMu.RLock()
	defer adminThemeMu.RUnlock()
	return adminThemeVal
}

func setAdminTheme(t string) {
	adminThemeMu.Lock()
	adminThemeVal = t
	adminThemeMu.Unlock()
}

func getAdminCompact() bool {
	adminCompactMu.RLock()
	defer adminCompactMu.RUnlock()
	return adminCompactOn
}

func setAdminCompact(on bool) {
	adminCompactMu.Lock()
	adminCompactOn = on
	adminCompactMu.Unlock()
}

func getFaviconPath() string {
	faviconMu.RLock()
	defer faviconMu.RUnlock()
	return faviconPathVal
}

func getFaviconVer() string {
	faviconMu.RLock()
	defer faviconMu.RUnlock()
	return faviconVerVal
}

func setFavicon(path, ver string) {
	faviconMu.Lock()
	faviconPathVal = path
	faviconVerVal = ver
	faviconMu.Unlock()
}

type App struct {
	db              *sql.DB
	adminPrefix     string
	adminCookiePath string
	uploadsDir      string
	resendKey       string
	baseURL         string
	secureCookies   bool

	publicTemplates map[string]*template.Template
	adminTemplates  map[string]*template.Template
	indexCacheMu    sync.RWMutex
	indexCacheHTML  []byte
	liveTrains      liveTrainsCache
}

type nonceKey struct{}

type adminPage struct {
	Title       string
	Flash       string
	CSRFToken   string
	AdminPrefix string
	Data        interface{}
	Authed      bool // an admin is logged in (nav section links shown)
	IsCentral   bool // logged-in admin is the central admin
	PermLevel   int  // cumulative permission level of the logged-in admin (0 = central/all)
	Nonce       string
}

var siteBaseURL string

var funcMap = template.FuncMap{
	"staticVer": func() string { return staticVersion },
	"delayText": func(min int) string {
		switch {
		case min > 1:
			return fmt.Sprintf("%d min late", min)
		case min < -1:
			return fmt.Sprintf("%d min early", -min)
		default:
			return "On time"
		}
	},
	"siteName":  func() string { return getSiteName() },
	"adminNavStyle": func() template.HTML {
		t, ok := siteThemes[getAdminTheme()]
		if !ok {
			t = siteThemes["default"]
		}
		css := "<style>:root{" + t.navVars + "}</style>"
		return template.HTML(css)
	},
	"siteAccentStyle": func() template.HTML {
		t, ok := siteThemes[getAdminTheme()]
		if !ok {
			t = siteThemes["default"]
		}
		css := ":root{" + t.publicLight + "}" +
			"@media(prefers-color-scheme:dark){:root:not([data-theme=light]){" + t.publicDark + "}}" +
			"[data-theme=dark]{" + t.publicDark + "}"
		if getAdminCompact() {
			// Compact mode: shrink the page hero header (keep the h1 at its
			// normal size) and remove the 2rem top margin on section headings
			// across the page.
			css += `.page-hero{padding:1rem 2rem}` +
				`.page-hero nav{margin-bottom:.25rem!important}` +
				`.page-hero .corridor-region{margin-top:.25rem}` +
				`.section-heading{margin-top:0!important}`
		}
		return template.HTML("<style>" + css + "</style>")
	},
	"faviconURL": func() string {
		p := getFaviconPath()
		if p == "" {
			return ""
		}
		return "/uploads/" + p + "?v=" + getFaviconVer()
	},
	"isYoutube": func(domain string) bool {
		return domain == "youtube.com" || domain == "youtu.be"
	},
	"isVimeo": func(domain string) bool {
		return domain == "vimeo.com"
	},
	"youtubeID": func(u string) string {
		if idx := strings.Index(u, "v="); idx >= 0 {
			id := u[idx+2:]
			if end := strings.Index(id, "&"); end >= 0 {
				id = id[:end]
			}
			return id
		}
		return ""
	},
	"truncate": func(s string, n int) string {
		if len(s) <= n {
			return s
		}
		return s[:n] + "…"
	},
	"year": func() string {
		return time.Now().Format("2006")
	},
	"mediaLabel": func(mt string) string {
		switch mt {
		case "image":
			return "Image"
		case "video":
			return "Video"
		case "website":
			return "Website"
		}
		return mt
	},
	"fmtFloat": func(f float64) string {
		return fmt.Sprintf("%.6f", f)
	},
	"fmtFloatN": func(f sql.NullFloat64) string {
		if !f.Valid {
			return ""
		}
		s := strconv.FormatFloat(f.Float64, 'f', 4, 64)
		s = strings.TrimRight(s, "0")
		s = strings.TrimSuffix(s, ".")
		return s
	},
	"mediaURL": func(m Media) string {
		return m.PublicURL()
	},
	"hasGeo": func(m Media) bool {
		return m.HasGeo()
	},
	"trainHasGeo": func(t Train) bool {
		return t.HasHeroGeo()
	},
	"heroLat": func(t Train) string {
		if !t.HeroLat.Valid {
			return ""
		}
		return fmt.Sprintf("%.6f", t.HeroLat.Float64)
	},
	"heroLon": func(t Train) string {
		if !t.HeroLon.Valid {
			return ""
		}
		return fmt.Sprintf("%.6f", t.HeroLon.Float64)
	},
	"nullInt64": func(n sql.NullInt64) int64 {
		return n.Int64
	},
	"nullFloat64": func(n sql.NullFloat64) float64 {
		return n.Float64
	},
	"safeURL": func(u string) template.URL {
		return template.URL(u)
	},
	"videoTags": func(tags string) []string {
		if tags == "" {
			return nil
		}
		labels := map[string]string{
			"long_consist": "Long consist", "doubleheader": "Doubleheader",
			"sandwich_set": "Sandwich set", "reverse_set": "Reverse set",
			"blow_overs": "Blow overs", "horn_show": "Horn show", "doppler": "Doppler",
		}
		var out []string
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(t)
			if l, ok := labels[t]; ok {
				out = append(out, l)
			}
		}
		return out
	},
	"hasTags": func(tags, tag string) bool {
		for _, t := range strings.Split(tags, ",") {
			if strings.TrimSpace(t) == tag {
				return true
			}
		}
		return false
	},
	"transitdocsURL": func(trainNumber string) string {
		now := time.Now()
		return fmt.Sprintf("https://asm.transitdocs.com/train/%d/%d/%d/A/%s",
			now.Year(), int(now.Month()), now.Day(), trainNumber)
	},
	"rarityBadges": func(tags string) []rarityBadge {
		return rarityBadgesForTags(tags)
	},
	"inc": func(i int) int { return i + 1 },
	"baseURL": func() string {
		return siteBaseURL
	},
	"absURL": func(u string) string {
		if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
			return u
		}
		return siteBaseURL + u
	},
	"metaDescription": func(s string, n int) string {
		s = strings.TrimSpace(s)
		if s == "" {
			return ""
		}
		if len(s) <= n {
			return s
		}
		return s[:n] + "…"
	},
}

// rarityBadge pairs a rarity tag with its display emoji and label, used to
// annotate videos on the train page and build the rarity-summary anchor line.
type rarityBadge struct {
	Tag   string
	Emoji string
	Label string
}

var rarityEmojis = map[string]string{
	"long_consist": "🚃", "doubleheader": "🚂", "sandwich_set": "🥪", "reverse_set": "🔄",
	"scenic": "🌄", "environment": "🌿", "historic": "🏛️", "special_event": "🎉",
}

var rarityOrder = []string{"long_consist", "doubleheader", "sandwich_set", "reverse_set", "scenic", "environment", "historic", "special_event"}

var rarityLabels = map[string]string{
	"long_consist": "Long consist", "doubleheader": "Doubleheader",
	"sandwich_set": "Sandwich set", "reverse_set": "Reverse set",
	"scenic": "Scenic", "environment": "Environment", "historic": "Historic", "special_event": "Special event",
}

func rarityBadgesForTags(tags string) []rarityBadge {
	if tags == "" {
		return nil
	}
	present := map[string]bool{}
	for _, t := range strings.Split(tags, ",") {
		present[strings.TrimSpace(t)] = true
	}
	var out []rarityBadge
	for _, tag := range rarityOrder {
		if present[tag] {
			out = append(out, rarityBadge{Tag: tag, Emoji: rarityEmojis[tag], Label: rarityLabels[tag]})
		}
	}
	return out
}

func computeStaticVersion() string {
	h := sha256.New()
	for _, path := range []string{"static/style.css", "static/theme.js"} {
		if data, err := embedFS.ReadFile(path); err == nil {
			h.Write(data)
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:8]
}

func buildTemplates() (pub, adm map[string]*template.Template, err error) {
	pub = make(map[string]*template.Template)
	adm = make(map[string]*template.Template)

	entries, err := embedFS.ReadDir("templates")
	if err != nil {
		return nil, nil, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".html") || e.Name() == "base.html" {
			continue
		}
		t, err := template.New("base").Funcs(funcMap).ParseFS(embedFS,
			"templates/base.html", "templates/"+e.Name())
		if err != nil {
			return nil, nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		pub[e.Name()] = t
	}

	entries, err = embedFS.ReadDir("templates/admin")
	if err != nil {
		return nil, nil, err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".html") {
			continue
		}
		t, err := template.New("admin_base").Funcs(funcMap).ParseFS(embedFS,
			"templates/admin_base.html", "templates/admin/"+e.Name())
		if err != nil {
			return nil, nil, fmt.Errorf("parse admin/%s: %w", e.Name(), err)
		}
		adm[e.Name()] = t
	}
	return pub, adm, nil
}

func (app *App) invalidateIndexCache() {
	app.indexCacheMu.Lock()
	app.indexCacheHTML = nil
	app.indexCacheMu.Unlock()
}

// withCurrentUser populates the logged-in user / CSRF token on a publicPage so
// the shared base template can render the account nav. Non-publicPage data
// (and the anonymous case) pass through unchanged.
func (app *App) withCurrentUser(r *http.Request, data interface{}) interface{} {
	pp, ok := data.(publicPage)
	if !ok || r == nil {
		return data
	}
	if u, csrf := app.getUserSession(r); u != nil {
		pp.CurrentUser = u
		pp.UserCSRF = csrf
	}
	if nonce, ok := r.Context().Value(nonceKey{}).(string); ok {
		pp.Nonce = nonce
	}
	return pp
}

func (app *App) renderPublicToBuffer(r *http.Request, page string, data interface{}) ([]byte, error) {
	tmpl, ok := app.publicTemplates[page]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", page)
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", app.withCurrentUser(r, data)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (app *App) renderPublic(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	b, err := app.renderPublicToBuffer(r, page, data)
	if err != nil {
		log.Printf("template error (%s): %v", page, err)
		http.Error(w, "Template error", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(b)
}

func (app *App) renderAdmin(w http.ResponseWriter, r *http.Request, page string, data adminPage) {
	data.AdminPrefix = app.adminPrefix
	if au := adminFromCtx(r); au != nil {
		data.Authed = true
		data.PermLevel = au.PermissionLevel
		data.IsCentral = au.IsCentral()
	}
	if nonce, ok := r.Context().Value(nonceKey{}).(string); ok {
		data.Nonce = nonce
	}
	tmpl, ok := app.adminTemplates[page]
	if !ok {
		log.Printf("admin template not found: %s", page)
		http.Error(w, "Template error", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "admin_base", data); err != nil {
		log.Printf("admin template exec error (%s): %v", page, err)
	}
}

func newToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func (app *App) limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 32<<20) // 32 MB
		next.ServeHTTP(w, r)
	})
}

func newNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (app *App) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := newNonce()
		ctx := context.WithValue(r.Context(), nonceKey{}, nonce)

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		if app.secureCookies {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		n := "'nonce-" + nonce + "'"
		if strings.HasPrefix(r.URL.Path, app.adminCookiePath+"/") || r.URL.Path == app.adminCookiePath {
			w.Header().Set("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self' "+n+"; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' https: data: blob:; "+
					"frame-src https://www.youtube.com https://player.vimeo.com; "+
					"object-src 'none';")
		} else {
			w.Header().Set("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self' "+n+" https://cdn.jsdelivr.net; "+
					"style-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'; "+
					"img-src 'self' https://tile.openstreetmap.org https://*.tile.openstreetmap.org "+
					"https://*.basemaps.cartocdn.com "+
					"https://upload.wikimedia.org "+
					"https://i.imgur.com https://live.staticflickr.com https://i.ytimg.com "+
					"https://img.youtube.com https://vumbnail.com data: blob:; "+
					"frame-src https://www.youtube.com https://player.vimeo.com; "+
					"object-src 'none';")
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// isTrue reports whether an env value means "on" (1/true/yes/on, case-insensitive).
func isTrue(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Bind to loopback only by default: the app sits behind nginx, which is the
	// only thing that should reach it. Override with BIND_ADDR if ever needed.
	// In dev mode (DEV=true) we bind to all interfaces so the server is reachable
	// from other machines on the LAN (e.g. http://192.168.x.x:PORT/).
	dev := isTrue(os.Getenv("DEV"))
	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		if dev {
			bindAddr = "0.0.0.0"
		} else {
			bindAddr = "127.0.0.1"
		}
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "amazingtrak.db"
	}
	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "./uploads"
	}
	adminSecret := os.Getenv("ADMIN_SECRET")
	if adminSecret == "" {
		adminSecret = "admin/secret"
	}
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "secret"
	}
	// Email is entirely optional: with no RESEND_API_KEY the site sends nothing
	// and no email-dependent feature blocks a user.
	resendKey := os.Getenv("RESEND_API_KEY")
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}
	secureCookies := strings.HasPrefix(baseURL, "https://")
	siteBaseURL = strings.TrimSuffix(baseURL, "/")

	firstSegment := adminSecret
	if i := strings.Index(adminSecret, "/"); i >= 0 {
		firstSegment = adminSecret[:i]
	}

	if err := os.MkdirAll(uploadsDir+"/images", 0755); err != nil {
		log.Fatal("create uploads dir:", err)
	}

	staticVersion = computeStaticVersion()

	pubTmpl, admTmpl, err := buildTemplates()
	if err != nil {
		log.Fatal("build templates:", err)
	}

	db, err := openDB(dbPath)
	if err != nil {
		log.Fatal("openDB:", err)
	}
	defer db.Close()

	if err := seedDB(db, adminUsername, adminPassword); err != nil {
		log.Fatal("seedDB:", err)
	}

	if prefs, err := getSitePrefs(db); err == nil {
		if prefs.SiteName != "" {
			setSiteName(prefs.SiteName)
		}
		if prefs.AdminTheme != "" {
			setAdminTheme(prefs.AdminTheme)
		}
		setAdminCompact(prefs.AdminCompact)
		if prefs.FaviconPath != "" {
			ver := "1"
			if fi, err := os.Stat(filepath.Join(uploadsDir, prefs.FaviconPath)); err == nil {
				ver = fmt.Sprintf("%d", fi.ModTime().Unix())
			}
			setFavicon(prefs.FaviconPath, ver)
		}
	}

	app := &App{
		db:              db,
		adminPrefix:     "/" + adminSecret,
		adminCookiePath: "/" + firstSegment,
		uploadsDir:      uploadsDir,
		resendKey:       resendKey,
		baseURL:         baseURL,
		secureCookies:   secureCookies,
		publicTemplates: pubTmpl,
		adminTemplates:  admTmpl,
	}

	go func() {
		for range time.Tick(time.Hour) {
			app.purgeExpiredSessions()
			app.purgeOldRateLimitLogs()
		}
	}()

	// Live train positions; a no-op while the feature is disabled in Settings.
	go app.pollLiveTrains()

	mux := http.NewServeMux()

	// Static embedded files — long cache when ?v= is present (versioned URL)
	staticFS := http.FileServer(http.FS(embedFS))
	mux.Handle("GET /static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		staticFS.ServeHTTP(w, r)
	}))

	// Uploaded files from disk — directory listing disabled.
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/",
		noListFileServer(http.Dir(uploadsDir))))

	// Health check
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Favicon (admin-uploaded; served from a fixed path regardless of /static versioning)
	mux.HandleFunc("GET /favicon.ico", app.handleFavicon)

	// Theme toggle (sets cookie, redirects back)
	mux.HandleFunc("POST /theme", app.handleThemeToggle)

	// Public pages
	mux.HandleFunc("GET /{$}", app.handleIndex)
	mux.HandleFunc("GET /overview", app.handleOverview)
	mux.HandleFunc("GET /trains-list", app.handleTrainsList)
	mux.HandleFunc("GET /map", app.handleMap)
	mux.HandleFunc("GET /api/amtrak-routes", app.handleAmtrakRoutes)
	mux.HandleFunc("GET /api/live-trains", app.handleLiveTrains)
	mux.HandleFunc("GET /api/vantage-spots", app.handleVantageSpotsAPI)
	mux.HandleFunc("GET /vantage-spots/suggest", app.handleVantageSpotForm)
	mux.HandleFunc("POST /vantage-spots/suggest", app.handleVantageSpotSubmit)
	mux.HandleFunc("GET /routes", app.handleCorridors)
	mux.HandleFunc("GET /routes/{slug}", app.handleCorridor)
	// Permanent redirects from the former /corridors URLs (renamed to /routes)
	// so existing bookmarks and search-engine links keep working.
	mux.HandleFunc("GET /corridors", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/routes", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /corridors/{slug}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/routes/"+r.PathValue("slug"), http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /trains/{slug}", app.handleTrain)
	mux.HandleFunc("GET /trains/{slug}/suggest", app.handleSuggestForm)
	mux.HandleFunc("POST /trains/{slug}/suggest", app.handleSuggestSubmit)
	mux.HandleFunc("POST /trains/{slug}/comment", app.requireUser(app.handleCommentSubmit))
	mux.HandleFunc("POST /routes/{slug}/comment", app.requireUser(app.handleCorridorCommentSubmit))
	mux.HandleFunc("GET /stations/{slug}", app.handleStation)

	// Conductor self-service (managed through public pages; never the admin URL).
	mux.HandleFunc("POST /routes/{slug}/conductor/request", app.requireUser(app.handleConductorRequest))
	mux.HandleFunc("GET /routes/{slug}/trains/new", app.requireUser(app.handleConductorTrainNewForm))
	mux.HandleFunc("POST /routes/{slug}/trains/new", app.requireUser(app.handleConductorTrainCreate))
	mux.HandleFunc("GET /trains/{slug}/edit", app.requireUser(app.handleConductorTrainEditForm))
	mux.HandleFunc("POST /trains/{slug}/edit", app.requireUser(app.handleConductorTrainEdit))
	mux.HandleFunc("POST /trains/{slug}/deactivate", app.requireUser(app.handleConductorTrainToggle))
	mux.HandleFunc("GET /trains/{slug}/edit/stops", app.requireUser(app.handleConductorStopsForm))
	mux.HandleFunc("POST /trains/{slug}/edit/stops", app.requireUser(app.handleConductorStopsUpdate))

	// Public user accounts (registration / login)
	mux.HandleFunc("GET /register", app.handleRegisterForm)
	mux.HandleFunc("POST /register", app.handleRegisterPost)
	mux.HandleFunc("GET /login", app.handleUserLoginForm)
	mux.HandleFunc("POST /login", app.handleUserLoginPost)
	mux.HandleFunc("POST /logout", app.handleUserLogout)
	mux.HandleFunc("GET /forgot-password", app.handleForgotPasswordForm)
	mux.HandleFunc("POST /forgot-password", app.handleForgotPasswordPost)
	mux.HandleFunc("GET /reset-password", app.handleResetPasswordForm)
	mux.HandleFunc("POST /reset-password", app.handleResetPasswordPost)
	mux.HandleFunc("GET /confirm-email", app.handleConfirmEmail)
	mux.HandleFunc("POST /resend-verification", app.requireUser(app.handleResendVerification))
	mux.HandleFunc("POST /users/change-password", app.requireUser(app.handleUserChangePassword))
	mux.HandleFunc("GET /users/{username}", app.requireUser(app.handleUserProfile))

	p := app.adminPrefix

	// Admin auth
	mux.HandleFunc("GET "+p+"/login", app.handleAdminLogin)
	mux.HandleFunc("POST "+p+"/login", app.handleAdminLoginPost)
	mux.HandleFunc("POST "+p+"/logout", app.requireAdmin(app.handleAdminLogout))

	// Admin dashboard
	mux.HandleFunc("GET "+p, app.requireAdmin(app.handleAdminDashboard))
	mux.HandleFunc("GET "+p+"/", app.requireAdmin(app.handleAdminDashboard))

	// Admin corridors (permission level 4)
	rc := func(h http.HandlerFunc) http.HandlerFunc { return app.requirePermission(4, h) }
	mux.HandleFunc("GET "+p+"/routes", rc(app.handleAdminCorridors))
	mux.HandleFunc("POST "+p+"/routes/new", rc(app.handleAdminCorridorCreate))
	mux.HandleFunc("POST "+p+"/routes/delete-inactive", rc(app.handleAdminCorridorDeleteInactive))
	mux.HandleFunc("POST "+p+"/routes/sync-schedule-urls", rc(app.handleAdminSyncScheduleURLs))
	mux.HandleFunc("GET "+p+"/routes/{id}", rc(app.handleAdminCorridorEdit))
	mux.HandleFunc("POST "+p+"/routes/{id}", rc(app.handleAdminCorridorUpdate))
	mux.HandleFunc("POST "+p+"/routes/{id}/toggle", rc(app.handleAdminCorridorToggle))
	mux.HandleFunc("GET "+p+"/routes/{id}/media", rc(app.handleAdminCorridorMedia))
	mux.HandleFunc("POST "+p+"/routes/{id}/media/add", rc(app.handleAdminCorridorMediaAdd))
	mux.HandleFunc("POST "+p+"/routes/{id}/media/{mid}/delete", rc(app.handleAdminCorridorMediaDelete))
	mux.HandleFunc("POST "+p+"/routes/{id}/media/{mid}/hero", rc(app.handleAdminCorridorMediaHero))
	mux.HandleFunc("POST "+p+"/routes/{id}/media/{mid}/thumbnail", rc(app.handleAdminCorridorMediaThumbnail))
	mux.HandleFunc("POST "+p+"/routes/{id}/media/{mid}/geo", rc(app.handleAdminCorridorMediaGeo))
	mux.HandleFunc("POST "+p+"/routes/{id}/media/{mid}/caption", rc(app.handleAdminCorridorMediaCaption))
	mux.HandleFunc("POST "+p+"/routes/{id}/media/{mid}/title", rc(app.handleAdminCorridorMediaTitle))

	// Conductors (corridor maintainers) — level 4
	mux.HandleFunc("GET "+p+"/conductors", rc(app.handleAdminConductors))
	mux.HandleFunc("POST "+p+"/conductors/{id}/approve", rc(app.handleAdminConductorApprove))
	mux.HandleFunc("POST "+p+"/conductors/{id}/reject", rc(app.handleAdminConductorReject))
	mux.HandleFunc("POST "+p+"/routes/{id}/conductor/set", rc(app.handleAdminCorridorConductorSet))
	mux.HandleFunc("POST "+p+"/routes/{id}/conductor/remove", rc(app.handleAdminCorridorConductorRemove))

	// Admin trains (permission level 3)
	rt := func(h http.HandlerFunc) http.HandlerFunc { return app.requirePermission(3, h) }
	mux.HandleFunc("GET "+p+"/trains", rt(app.handleAdminTrains))
	mux.HandleFunc("POST "+p+"/trains/new", rt(app.handleAdminTrainCreate))
	mux.HandleFunc("POST "+p+"/trains/delete-inactive", rt(app.handleAdminTrainDeleteInactive))
	mux.HandleFunc("GET "+p+"/trains/{id}", rt(app.handleAdminTrainDetail))
	mux.HandleFunc("POST "+p+"/trains/{id}", rt(app.handleAdminTrainUpdate))
	mux.HandleFunc("POST "+p+"/trains/{id}/toggle", rt(app.handleAdminTrainToggle))
	mux.HandleFunc("GET "+p+"/trains/{id}/media", rt(app.handleAdminTrainMedia))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/add", rt(app.handleAdminTrainMediaAdd))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/delete", rt(app.handleAdminTrainMediaDelete))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/hero", rt(app.handleAdminTrainMediaHero))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/thumbnail", rt(app.handleAdminTrainMediaThumbnail))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/map", rt(app.handleAdminTrainMediaMap))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/geo", rt(app.handleAdminTrainMediaGeo))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/caption", rt(app.handleAdminTrainMediaCaption))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/title", rt(app.handleAdminTrainMediaTitle))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/tags", rt(app.handleAdminTrainMediaTags))
	mux.HandleFunc("POST "+p+"/trains/{id}/media/{mid}/best", rt(app.handleAdminTrainMediaBest))
	mux.HandleFunc("GET "+p+"/trains/{id}/stops", rt(app.handleAdminTrainStops))
	mux.HandleFunc("POST "+p+"/trains/{id}/stops", rt(app.handleAdminTrainStopsUpdate))

	// Admin suggestions (permission level 1)
	rs := func(h http.HandlerFunc) http.HandlerFunc { return app.requirePermission(1, h) }
	mux.HandleFunc("GET "+p+"/suggestions", rs(app.handleAdminAllSuggestions))
	mux.HandleFunc("POST "+p+"/suggestions/approve-all-pending", rs(app.handleAdminSuggestionApproveAll))
	mux.HandleFunc("POST "+p+"/suggestions/reject-all-pending", rs(app.handleAdminSuggestionRejectAll))
	mux.HandleFunc("POST "+p+"/suggestions/delete-all-rejected", rs(app.handleAdminSuggestionDeleteRejected))
	mux.HandleFunc("POST "+p+"/suggestions/{id}/approve", rs(app.handleAdminSuggestionApprove))
	mux.HandleFunc("POST "+p+"/suggestions/{id}/reject", rs(app.handleAdminSuggestionReject))
	mux.HandleFunc("POST "+p+"/suggestions/{id}/unapprove", rs(app.handleAdminSuggestionUnapprove))
	mux.HandleFunc("POST "+p+"/suggestions/{id}/spam", rs(app.handleAdminSuggestionMarkSpam))
	mux.HandleFunc("POST "+p+"/suggestions/{id}/edit", rs(app.handleAdminSuggestionEdit))

	// Admin vantage spots (permission level 7)
	rv := func(h http.HandlerFunc) http.HandlerFunc { return app.requirePermission(7, h) }
	mux.HandleFunc("GET "+p+"/vantage-spots", rv(app.handleAdminVantageSpots))
	mux.HandleFunc("POST "+p+"/vantage-spots/{id}/approve", rv(app.handleAdminVantageSpotApprove))
	mux.HandleFunc("POST "+p+"/vantage-spots/{id}/reject", rv(app.handleAdminVantageSpotReject))
	mux.HandleFunc("POST "+p+"/vantage-spots/{id}/delete", rv(app.handleAdminVantageSpotDelete))

	// Admin comments (permission level 2)
	rco := func(h http.HandlerFunc) http.HandlerFunc { return app.requirePermission(2, h) }
	mux.HandleFunc("GET "+p+"/comments", rco(app.handleAdminComments))
	mux.HandleFunc("POST "+p+"/comments/approve-all-pending", rco(app.handleAdminCommentApproveAll))
	mux.HandleFunc("POST "+p+"/comments/reject-all-pending", rco(app.handleAdminCommentRejectAll))
	mux.HandleFunc("POST "+p+"/comments/delete-all-rejected", rco(app.handleAdminCommentDeleteRejected))
	mux.HandleFunc("POST "+p+"/comments/{id}/approve", rco(app.handleAdminCommentApprove))
	mux.HandleFunc("POST "+p+"/comments/{id}/reject", rco(app.handleAdminCommentReject))
	mux.HandleFunc("POST "+p+"/comments/{id}/unapprove", rco(app.handleAdminCommentUnapprove))

	// Admin users: list/approve/unapprove/delete at level 6; add-admin/delete-admin central only
	ru := func(h http.HandlerFunc) http.HandlerFunc { return app.requirePermission(6, h) }
	rca := app.requireCentralAdmin
	mux.HandleFunc("GET "+p+"/users", ru(app.handleAdminUsers))
	mux.HandleFunc("POST "+p+"/users/add-user", ru(app.handleAdminAddUser))
	mux.HandleFunc("POST "+p+"/users/add-admin", rca(app.handleAdminAddAdmin))
	mux.HandleFunc("POST "+p+"/users/delete-unapproved", ru(app.handleAdminUsersDeleteUnapproved))
	mux.HandleFunc("GET "+p+"/users/{id}/submissions", ru(app.handleAdminUserSubmissions))
	mux.HandleFunc("POST "+p+"/users/{id}/approve", ru(app.handleAdminUserApprove))
	mux.HandleFunc("POST "+p+"/users/{id}/unapprove", ru(app.handleAdminUserUnapprove))
	mux.HandleFunc("POST "+p+"/users/{id}/delete", ru(app.handleAdminUserDelete))
	mux.HandleFunc("POST "+p+"/users/{id}/anonymize", ru(app.handleAdminUserAnonymize))
	mux.HandleFunc("POST "+p+"/users/{id}/reset-password", ru(app.handleAdminUserResetPassword))
	mux.HandleFunc("POST "+p+"/users/{id}/unlock", ru(app.handleAdminUserUnlock))
	mux.HandleFunc("POST "+p+"/admins/{id}/delete", rca(app.handleAdminAdminDelete))

	// Admin settings (permission level 5)
	mux.HandleFunc("GET "+p+"/settings", app.requirePermission(5, app.handleAdminSettings))
	mux.HandleFunc("POST "+p+"/settings", app.requirePermission(5, app.handleAdminSettingsPost))
	mux.HandleFunc("GET "+p+"/email-errors", app.requirePermission(5, app.handleAdminEmailErrors))
	mux.HandleFunc("POST "+p+"/email-errors/clear", app.requirePermission(5, app.handleAdminEmailErrorsClear))

	handler := app.limitBody(app.securityHeaders(mux))
	srv := &http.Server{
		Addr:              bindAddr + ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	log.Printf("listening on %s:%s  admin at %s", bindAddr, port, app.adminPrefix)
	log.Fatal(srv.ListenAndServe())
}

func (app *App) handleFavicon(w http.ResponseWriter, r *http.Request) {
	path := getFaviconPath()
	if path == "" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, filepath.Join(app.uploadsDir, path))
}

// noListFileServer wraps an http.FileSystem to return 404 for directory requests,
// preventing directory listing while still serving individual files.
func noListFileServer(fs http.FileSystem) http.Handler {
	fserver := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(r.URL.Path)
		if err == nil {
			defer f.Close()
			if stat, err := f.Stat(); err == nil && stat.IsDir() {
				http.NotFound(w, r)
				return
			}
		}
		fserver.ServeHTTP(w, r)
	})
}

func (app *App) handleThemeToggle(w http.ResponseWriter, r *http.Request) {
	theme := r.FormValue("theme")
	if theme != "light" && theme != "dark" && theme != "auto" {
		theme = "auto"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "theme",
		Value:    theme,
		Path:     "/",
		MaxAge:   365 * 24 * 3600,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
	// Redirect back to the referring page, but only if it's the same host.
	// This prevents the Referer header from being used as an open redirect.
	ref := "/"
	if raw := r.Header.Get("Referer"); raw != "" {
		if u, err := url.Parse(raw); err == nil && (u.Host == "" || u.Host == r.Host) {
			ref = u.RequestURI()
		}
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}
