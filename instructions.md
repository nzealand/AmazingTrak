Spec 1
ok i want a simple train tracker website, it will allow the secret unpublished administrator admin/secret to manage the current and inactive trains by default with all the following trains configured in the amtrak sqlite table by route, it will allow the admin to authenticate and access an admin train page to add edit existing Corridors or create new Corridors, and also inactivate and add additional trains and indicate if they are active or inactive and put them under the existing or new Corridor as relevant. the admin will also be able to drill into each train from the admin train page and add links, the admin can add links to either videos or images or websites in general. at least one image should be marked by the admin as a hero image. there will be a public set of pages to list all the trains and drill into the train details to see the train details and the hero image will be shown by default, and a list of published videos and published images and published links. any visitor will be able to contribute but not publish photos or videos but not contribute other links, and so only links to videos and photos should be allowed from the public either by limiting the domain to youtube etc or by validating the link is a video/image. the admin should see a list of suggested urls in each train admin page and be able to approve or reject individual suggestions or all suggestions from a single click
follow the go-web-simple including to use sqlite and go, the tables should be designed to track corridors, trains, links and suggestions.
Currently Operating Amtrak Train Numbers
Exact numbers shift slightly with schedules, but the following ranges are stable as of 2025–2026.
Long Distance Trains
Sunset Limited: 1, 2
Southwest Chief: 3, 4
California Zephyr: 5, 6
Empire Builder: 7, 8, 27, 28
Coast Starlight: 11, 14
Crescent: 19, 20
Texas Eagle: 21, 22, 421, 422
Floridian: 40, 41
Pennsylvanian: 42, 43
Lake Shore Limited: 48, 49, 448, 449
Cardinal: 50, 51
Auto Train: 52, 53
Vermonter: 54–57
City of New Orleans: 58, 59
Adirondack: 68, 69
Piedmont: 71–78
Carolinian: 79, 80
Palmetto: 89, 90
Silver Meteor: 97, 98
Northeast Regional (long-distance extensions): many in 65–67, 82–88, 93–96, 99, 111, 121–198, etc.
Winter Park Express (seasonal): 1105, 1106
NEC / Empire / Corridor Trains
Acela: 2100–2290 (varies by trip)
Northeast Regional (core NEC): roughly 100–199 block (many 101–198 variants)
Keystone Service: 600–674 range
Empire Service: 230–288 range
Hartford Line / Valley Flyer: 400–499 range (e.g., 400, 405–432, 450–499)
Capitol Corridor: 520–553, 720–751
Pacific Surfliner: 562–595, 761–794
Amtrak Cascades: 500–519
Hiawatha: 329–343
Lincoln Service: 300–302, 304–307, 318–319
Illinois Zephyr / Carl Sandburg: 380–383
Illini / Saluki: 390–393
Blue Water: 364, 365
Pere Marquette: 370, 371
Wolverine: 350–355
Missouri River Runner: 311, 316, 318–319 (some overlap by schedule)
Heartland Flyer: 821, 822
Downeaster: 680–699, 1689
Borealis: 1333, 1340
Berkshire Flyer (seasonal): 1235, 1244
Ethan Allen Express: 290, 291
Gold Runner: 701–719
# Train Tracker App Design
## Verdict
Use:
```text
Go + SQLite + server-rendered HTML + Caddy + systemd
```
Core entities:
```text
corridors
trains
links
suggestions
admin_users
```
Avoid user accounts for visitors. Public submissions should be anonymous but heavily rate-limited and moderated.
---
## Recommended Tables
```sql
CREATE TABLE corridors (
id INTEGER PRIMARY KEY,
name TEXT NOT NULL UNIQUE,
slug TEXT NOT NULL UNIQUE,
description TEXT,
region TEXT,
is_active INTEGER NOT NULL DEFAULT 1,
sort_order INTEGER NOT NULL DEFAULT 0,
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE trains (
id INTEGER PRIMARY KEY,
corridor_id INTEGER NOT NULL REFERENCES corridors(id),
train_number TEXT NOT NULL,
name TEXT,
slug TEXT NOT NULL UNIQUE,
direction TEXT,
notes TEXT,
is_active INTEGER NOT NULL DEFAULT 1,
sort_order INTEGER NOT NULL DEFAULT 0,
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
UNIQUE(corridor_id, train_number)
);
CREATE TABLE links (
id INTEGER PRIMARY KEY,
train_id INTEGER NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
url TEXT NOT NULL,
title TEXT,
link_type TEXT NOT NULL CHECK (link_type IN ('image', 'video', 'website')),
source_domain TEXT,
is_hero INTEGER NOT NULL DEFAULT 0,
is_published INTEGER NOT NULL DEFAULT 1,
added_by TEXT NOT NULL DEFAULT 'admin',
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE suggestions (
id INTEGER PRIMARY KEY,
train_id INTEGER NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
url TEXT NOT NULL,
title TEXT,
link_type TEXT NOT NULL CHECK (link_type IN ('image', 'video')),
source_domain TEXT,
status TEXT NOT NULL DEFAULT 'pending'
CHECK (status IN ('pending', 'approved', 'rejected')),
submitter_ip_hash TEXT,
submitter_user_agent TEXT,
rejection_reason TEXT,
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
reviewed_at TEXT,
reviewed_by INTEGER REFERENCES admin_users(id)
);
CREATE TABLE admin_users (
id INTEGER PRIMARY KEY,
username TEXT NOT NULL UNIQUE,
password_hash TEXT NOT NULL,
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
last_login_at TEXT
);
CREATE INDEX idx_trains_corridor ON trains(corridor_id);
CREATE INDEX idx_links_train ON links(train_id);
CREATE INDEX idx_links_hero ON links(train_id, is_hero);
CREATE INDEX idx_suggestions_train_status ON suggestions(train_id, status);
CREATE INDEX idx_suggestions_ip_time ON suggestions(submitter_ip_hash, created_at);
```
---
## Hero Image Rule
Enforce in app logic:
```text
Each train may have zero or one hero image.
Only image links can be marked as hero.
When setting a new hero image, unset the old one.
```
Optional SQLite partial unique index:
```sql
CREATE UNIQUE INDEX one_hero_per_train
ON links(train_id)
WHERE is_hero = 1;
```
---
## Public Submission Rules
Visitors can submit only:
```text
image links
video links
```
Do not allow arbitrary websites from public users.
Allowed public domains:
```text
youtube.com
youtu.be
vimeo.com
flickr.com
imgur.com
instagram.com
railpictures.net
rrpicturearchives.net
```
For stronger control, start with only:
```text
youtube.com
youtu.be
vimeo.com
flickr.com
imgur.com
```
---
## Anti-Spam Protections
Use all of these:
1. Rate limit by IP hash.
2. Max 3 submissions per IP per hour.
3. Max 10 submissions per IP per day.
4. Require hidden honeypot field.
5. Reject forms submitted too quickly, under 3 seconds.
6. Normalize and deduplicate URLs.
7. Block duplicate pending suggestions.
8. Limit URL length.
9. Validate domain allowlist.
10. Require admin approval before publishing.
Suggested constraints:
```text
Max URL length: 500 chars
Max title length: 120 chars
Public submissions: pending only
No HTML in titles
No file uploads
No comments
No user accounts
```
---
## Admin Features
Admin can:
```text
Create/edit corridors
Activate/inactivate corridors
Create/edit trains
Activate/inactivate trains
Add/edit/delete links
Mark one image as hero
View pending suggestions per train
Approve one suggestion
Reject one suggestion
Approve all pending suggestions for a train
Reject all pending suggestions for a train
```
Approval flow:
```text
suggestions.status = 'approved'
copy suggestion into links
links.added_by = 'public_approved'
```
Do not directly expose suggestions on public pages.
---
## Security Hardening
### Authentication
Do:
```text
Use bcrypt or argon2id password hashing.
Use secure random session tokens.
Store sessions server-side.
Set HttpOnly, Secure, SameSite=Lax cookies.
Add CSRF tokens to all admin forms.
```
### Admin Path
Fine to use:
```text
/admin/secret
```
But treat it only as obscurity.
Real protection must be:
```text
password + session + CSRF + rate limiting
```
### Input Safety
Do:
```text
Use prepared SQL statements.
Escape all HTML output.
Validate slugs.
Normalize URLs.
Reject javascript:, data:, file:, ftp: URLs.
Only allow http/https.
```
### Headers
Set:
```text
Content-Security-Policy: default-src 'self'; img-src 'self' https:; frame-src https://www.youtube.com https://player.vimeo.com; object-src 'none';
X-Content-Type-Options: nosniff
Referrer-Policy: strict-origin-when-cross-origin
```
---
## Suggested Pages
Public:
```text
/
/corridors
/corridors/{slug}
/trains/{slug}
/trains/{slug}/suggest
```
Admin:
```text
/admin/secret/login
/admin/secret/trains
/admin/secret/corridors
/admin/secret/trains/{id}
/admin/secret/trains/{id}/suggestions
```
---
## Recommendation
Build this as a moderation-first directory, not a social site.
Best MVP:
1. Seed corridors and trains.
2. Public train list and train detail pages.
3. Admin login.
4. Admin CRUD for trains and links.
5. Public suggestion form.
6. Admin approve/reject queue.
7. Rate limiting and domain allowlist.
Avoid uploads, comments, accounts, ratings, and open-ended links until the site has real traffic.
Missing updated_at trigger. SQLite doesn't auto-update it. You need:
sqlCREATE TRIGGER links_updated_at
AFTER UPDATE ON links
BEGIN
UPDATE links SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
Repeat for all tables with updated_at.
Missing sessions table — you mentioned server-side sessions but didn't define the table:
sqlCREATE TABLE sessions (
id TEXT PRIMARY KEY, -- crypto/rand hex, 32+ bytes
admin_user_id INTEGER NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
expires_at TEXT NOT NULL,
ip_hash TEXT,
user_agent TEXT
);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
Missing audit_log — for an admin tool managing public content, you want a record of who approved/rejected what:
sqlCREATE TABLE audit_log (
id INTEGER PRIMARY KEY,
admin_user_id INTEGER REFERENCES admin_users(id),
action TEXT NOT NULL, -- 'approve_suggestion', 'reject_suggestion', 'edit_train', etc.
entity_type TEXT NOT NULL, -- 'suggestion', 'train', 'link'
entity_id INTEGER NOT NULL,
detail TEXT, -- JSON blob of before/after or extra context
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
suggestions.submitter_user_agent TEXT — store raw UA only if you need it for abuse analysis. It's somewhat PII-adjacent. Fine to keep, but trim it at 200 chars.
Anti-Spam Hardening
The list is good. These are the gaps:
Rate limiting needs a backing table, not just in-memory state. If the process restarts, limits reset. Add:
sqlCREATE TABLE rate_limit_log (
id INTEGER PRIMARY KEY,
ip_hash TEXT NOT NULL,
action TEXT NOT NULL DEFAULT 'suggest',
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_rate_limit_ip_time ON rate_limit_log(ip_hash, created_at);
Query count WHERE created_at > datetime('now', '-1 hour'). Purge rows older than 24h on a background tick.
Honeypot field name matters. Don't call it honeypot. Call it something plausible:
html<input type="text" name="website" autocomplete="off" tabindex="-1" 
style="position:absolute;left:-9999px" aria-hidden="true">
Timing check should be server-side. Client-side timing is trivially bypassed. Store a hidden form token with a timestamp in the session or a signed cookie, then validate server-side that at least 4 seconds elapsed between token issuance and form submission.
URL normalization steps to enforce:
Lowercase scheme and host
Strip tracking params (?utm_*, ?fbclid, etc.)
Normalize YouTube: youtu.be/xyz → youtube.com/watch?v=xyz before dedup check
Strip trailing slashes
Reject URLs with auth components (user:pass@host)
Deduplication scope: Check for duplicate pending and already-approved links, not just pending:
sqlSELECT 1 FROM suggestions WHERE train_id = ? AND url = ? AND status = 'pending'
UNION
SELECT 1 FROM links WHERE train_id = ? AND url = ?
Security Hardening
Authentication gaps:
Add login attempt throttling with lockout — not covered in the schema:
sqlCREATE TABLE login_attempts (
id INTEGER PRIMARY KEY,
ip_hash TEXT NOT NULL,
username TEXT,
succeeded INTEGER NOT NULL DEFAULT 0,
created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
Block IPs with 5+ failed attempts in 15 minutes.
Session token generation in Go:
gofunc newSessionToken() string {
b := make([]byte, 32)
_, err := rand.Read(b)
if err != nil {
panic(err)
}
return hex.EncodeToString(b)
}
64-char hex token, HttpOnly; Secure; SameSite=Lax; Path=/admin; Max-Age=86400.
CSRF implementation note: Use the double-submit cookie pattern or a synchronizer token stored in the session. Don't generate CSRF tokens from static secrets — use crypto/rand per session.
Missing: session expiry enforcement on read. When loading a session from the DB, always check expires_at > CURRENT_TIMESTAMP. Add a background goroutine to purge expired sessions every hour.
CSP adjustment: Your current CSP allows img-src 'self' https: which is broad. Tighten for the public pages:
img-src 'self' https://i.imgur.com https://live.staticflickr.com https://i.ytimg.com https://vumbnail.com;
For admin pages you can be more permissive.
Missing: Permissions-Policy header:
Permissions-Policy: camera=(), microphone=(), geolocation=()
SQL injection note: With Go's database/sql and ? placeholders, you're protected as long as you never use fmt.Sprintf to build queries. Add a linter rule or code review checklist item for this.
URL Validation Logic
The allowlist approach is correct. Implement as a two-stage check in Go:
govar allowedDomains = map[string]string{
"youtube.com": "video",
"youtu.be": "video",
"vimeo.com": "video",
"flickr.com": "image",
"imgur.com": "image",
"railpictures.net": "image",
"rrpicturearchives.net": "image",
}
func classifyPublicURL(raw string) (domain, linkType string, ok bool) {
u, err := url.Parse(strings.TrimSpace(raw))
if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
return "", "", false
}
host := strings.ToLower(u.Hostname())
host = strings.TrimPrefix(host, "www.")
t, found := allowedDomains[host]
if !found {
return "", "", false
}
return host, t, true
}
Auto-populate link_type from this — don't trust the user's submitted link_type.
Admin UX Recommendations
Bulk approve/reject should use a transaction and the audit log atomically:
sqlBEGIN;
INSERT INTO audit_log ...;
UPDATE suggestions SET status = 'approved', reviewed_at = ..., reviewed_by = ? 
WHERE train_id = ? AND status = 'pending';
INSERT INTO links (train_id, url, title, link_type, source_domain, added_by)
SELECT train_id, url, title, link_type, source_domain, 'public_approved'
FROM suggestions WHERE train_id = ? AND status = 'pending';
COMMIT;
Hero image enforcement in admin UI: When the admin marks a new hero, the UI should show the current hero with a warning before replacing it, not silently swap.
also consider having a suggestions page for all suggestions for easy review/approval, and seed the trains into the scema and seed the admin user admin is the username secret is the password, but dont show that anywhere in url or pages or anything it is secret.
Define a single env file with admin/secret as the admin path and reference that once
Include a page on the admin ui that allows the admin to change their username and password if they remember the current username and password
Update
Major Update: Delete and recreate the database. You have a California corridor and a capitol corridor. This is a cool way to group corridors, but California corridor does not really exist. You can group the capitol corridor and related corridors under the California region but there is no need for a page for California region as no one searches for such a thing. The regional grouping can just be an editable text box on the corridor page, the admin can ensure everything is consistent. There should be a corridor page, which includes capitol corridor and there should be a train page which includes Amtrak 742.  Each photo should be assigned to a train e.g. Amtrak 742, and each train should have its own page, where the admin can pick a hero image for the train. The admin should be able to pick a hero train for the corridor as well. The admin should also be able to pick a thumbnail. The admin might struggle with providing image urls so provide either the option to specify a url or the option to upload an image if you can to a local images folder or if you can also allow the admin to paste in an image, and name local images with meaningful names e.g. Amtrak-742-1. When the admin submits an image url or uploads the image provide an optional lat & long. If you can, when processing an uploaded image, figure out the lat and long geo encoding. On the train page, if you can, include a small map showing the location where the hero photo was taken. 
Provide a dark mode and a light mode, and either if you can detect the browser mode and auto set the mode, or allow users to quickly change and remember the change in a cookie that persists
Revised Train Tracker Schema
Core Change
Do not create fake regional corridor pages.
Use:
•	corridors for real public corridor pages, e.g. Capitol Corridor
•	region as editable text metadata, e.g. California
•	trains for individual train pages, e.g. Amtrak 742
•	media for photos, videos, websites, uploads
•	hero and thumbnail selection at train and corridor level
 
Schema
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;

CREATE TABLE corridors (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    region TEXT,
    description TEXT,
    on_time_percent REAL,
    service_quality_summary TEXT,
    hero_train_id INTEGER,
    hero_media_id INTEGER,
    thumbnail_media_id INTEGER,
    is_active INTEGER NOT NULL DEFAULT 1,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE trains (
    id INTEGER PRIMARY KEY,
    corridor_id INTEGER NOT NULL REFERENCES corridors(id) ON DELETE CASCADE,
    train_number TEXT NOT NULL,
    display_name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    direction TEXT,
    notes TEXT,
    hero_media_id INTEGER,
    thumbnail_media_id INTEGER,
    is_active INTEGER NOT NULL DEFAULT 1,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(corridor_id, train_number)
);

CREATE TABLE stops (
    id INTEGER PRIMARY KEY,
    corridor_id INTEGER NOT NULL REFERENCES corridors(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    station_code TEXT,
    latitude REAL,
    longitude REAL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE train_stops (
    id INTEGER PRIMARY KEY,
    train_id INTEGER NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    stop_id INTEGER NOT NULL REFERENCES stops(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    scheduled_arrival TEXT,
    scheduled_departure TEXT,
    UNIQUE(train_id, stop_id)
);

CREATE TABLE media (
    id INTEGER PRIMARY KEY,
    train_id INTEGER REFERENCES trains(id) ON DELETE CASCADE,
    corridor_id INTEGER REFERENCES corridors(id) ON DELETE CASCADE,

    media_type TEXT NOT NULL CHECK (media_type IN ('image', 'video', 'website')),
    source_type TEXT NOT NULL CHECK (source_type IN ('url', 'upload', 'paste', 'seed')),
    
    url TEXT,
    local_path TEXT,
    original_filename TEXT,
    stored_filename TEXT,

    title TEXT,
    caption TEXT,
    source_domain TEXT,

    latitude REAL,
    longitude REAL,
    location_name TEXT,
    location_source TEXT CHECK (location_source IN ('admin', 'exif', 'inferred', 'unknown')) DEFAULT 'unknown',

    is_published INTEGER NOT NULL DEFAULT 1,
    added_by TEXT NOT NULL DEFAULT 'admin',

    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (
        train_id IS NOT NULL OR corridor_id IS NOT NULL
    )
);

CREATE TABLE suggestions (
    id INTEGER PRIMARY KEY,
    train_id INTEGER NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    title TEXT,
    media_type TEXT NOT NULL CHECK (media_type IN ('image', 'video')),
    source_domain TEXT,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'approved', 'rejected')),
    submitter_ip_hash TEXT,
    submitter_user_agent TEXT,
    rejection_reason TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TEXT,
    reviewed_by INTEGER REFERENCES admin_users(id)
);

CREATE TABLE admin_users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    must_change_password INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TEXT
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    admin_user_id INTEGER NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TEXT NOT NULL,
    ip_hash TEXT,
    user_agent TEXT
);

CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY,
    admin_user_id INTEGER REFERENCES admin_users(id),
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INTEGER NOT NULL,
    detail TEXT,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE rate_limit_log (
    id INTEGER PRIMARY KEY,
    ip_hash TEXT NOT NULL,
    action TEXT NOT NULL DEFAULT 'suggest',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE login_attempts (
    id INTEGER PRIMARY KEY,
    ip_hash TEXT NOT NULL,
    username TEXT,
    succeeded INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE site_preferences (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    default_theme TEXT NOT NULL DEFAULT 'auto'
        CHECK (default_theme IN ('light', 'dark', 'auto')),
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
Indexes
CREATE INDEX idx_trains_corridor ON trains(corridor_id);
CREATE INDEX idx_media_train ON media(train_id);
CREATE INDEX idx_media_corridor ON media(corridor_id);
CREATE INDEX idx_suggestions_train_status ON suggestions(train_id, status);
CREATE INDEX idx_rate_limit_ip_time ON rate_limit_log(ip_hash, created_at);
CREATE INDEX idx_login_attempts_ip_time ON login_attempts(ip_hash, created_at);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
Seed Example
INSERT INTO corridors (name, slug, region, description)
VALUES (
    'Capitol Corridor',
    'capitol-corridor',
    'California',
    'Amtrak corridor service between the Sacramento area and the San Francisco Bay Area.'
);

INSERT INTO trains (corridor_id, train_number, display_name, slug, direction)
VALUES (
    1,
    '742',
    'Amtrak 742',
    'amtrak-742',
    NULL
);
Other
Image Handling
Admin should be able to add images three ways:
1.	URL
2.	Upload file
3.	Paste image from clipboard
Store uploaded or pasted files in:
/static/uploads/images/
Filename pattern:
Amtrak-742-1.jpg
Amtrak-742-2.jpg
Capitol-Corridor-1.jpg
Use database stored_filename and local_path, not raw filesystem guessing.
Geo Handling
When an image is uploaded:
1.	Try to read EXIF GPS.
2.	If GPS exists, populate latitude and longitude.
3.	If not, allow admin to enter optional latitude and longitude.
4.	Store location_source as exif, admin, or unknown.
Train page should show a small map only when hero image has latitude and longitude.
Theme
Use CSS media query by default:
@media (prefers-color-scheme: dark) {
  body { background: #111; color: #eee; }
}
Also allow manual override:
theme=light
theme=dark
theme=auto
Store user preference in a cookie.
Admin Pages
/admin/login
/admin/corridors
/admin/corridors/{id}
.admin/corridors/{id}/media
/admin/trains
/admin/trains/{id}
/admin/trains/{id}/media
/admin/suggestions
/admin/trains/{id}/suggestions
Public Pages
/
 /corridors
 /corridors/capitol-corridor
 /trains/amtrak-742
 /trains/amtrak-742/suggest
Security Notes
•	Do not put usernames or passwords in URLs.
•	Store only password hashes.
•	Require CSRF tokens on admin forms.
•	Use secure, HttpOnly, SameSite cookies.
•	Rate limit login attempts.
•	Rate limit public suggestions.
•	Allow public submissions only from approved image/video domains.
•	Do not allow public file uploads.
•	Admin uploads only.
•	Validate MIME type and file extension.
•	Re-encode uploaded images before serving.
•	Strip EXIF from public copies unless you intentionally display location.
Additional Recommendations
2. Image processing — specify the libraries explicitly
Add this to your implementation constraints:
For EXIF extraction use: github.com/rwcarlsen/goexif/exif
For image re-encoding (strip EXIF from public copies) use: standard library image/jpeg
For file type validation use: net/http DetectContentType on first 512 bytes
 
3. Map embed — specify Leaflet explicitly
Add:
For the hero photo location map on train pages, use Leaflet.js with OpenStreetMap 
tiles (no API key required). Load from CDN. Show a single marker at the hero 
image lat/long. Only render the map div if hero_media.latitude IS NOT NULL.

Map and route info: 

Data Prep
can you read and extract the relevant text from this? https://railroads.dot.gov/sites/fra.dot.gov/files/2026-03/FY26%20Q1%20Performance%20and%20Service%20Quality%20Report_PDFa.pdf

also extract the routes
extract the following links from [https://railrat.net/trains/](https://railrat.net/trains/542/)
e.g. create a link for the amtrak 542 to https://railrat.net/trains/542/


Data
Acela Express
816, 817, 880, 2102, 2103, 2104, 2108, 2109, 2110, 2113, 2115, 2121, 2122, 2123, 2124, 2126, 2130, 2150, 2151, 2152, 2153, 2154, 2155, 2159, 2162, 2163, 2166, 2167, 2168, 2169, 2170, 2171, 2172, 2173, 2174, 2190, 2192, 2193, 2201, 2203, 2205, 2206, 2207, 2214, 2215, 2216, 2218, 2220, 2222, 2223, 2224, 2226, 2228, 2233, 2247, 2248, 2249, 2250, 2251, 2252, 2253, 2254, 2255, 2256, 2257, 2258, 2259, 2262, 2263, 2265, 2271, 2274, 2275, 2290, 2292, 2295

Adirondack
68, 69

Amtrak Cascades
500, 502, 503, 504, 505, 506, 507, 508, 509, 511, 516, 517, 518, 519

Auto Train
52, 53

Berkshire Flyer
1233, 1234, 1246

Blue Water/Michigan Service
364, 365, 1364

Borealis
1333, 1340

California Zephyr
5, 6(2), 1005, 1006

Capitol Corridor
520, 521, 522, 523, 524, 525, 526, 527, 528, 529, 530, 531, 532, 534, 535, 536, 537, 538, 539, 540, 541, 542, 543, 544, 545, 546, 547, 548, 549, 550, 551, 720, 723, 724, 727, 728, 729, 732, 733, 734, 736, 737, 738, 741, 742, 743, 744, 745, 746, 747, 748, 749, 750, 751

Cardinal
50, 51

Carolinian / Piedmont
71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 105, 1072, 1075, 1123, 1171, 1172

City of New Orleans
58, 59, 1058, 1059

Coast Starlight
11(2), 14(2), 1011, 1014

Crescent
19(2), 20(2), 1019, 1020

Downeaster
680, 681, 682, 683, 684, 685, 686, 687, 688, 689, 690, 691, 692, 693, 694, 695, 696, 697, 698, 699, 1689, 1697

Empire Builder
7(2), 8(2), 27, 28, 1007, 1008, 1027, 1028

Empire Service
230, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 243, 244, 245, 280, 281, 283, 284, 1237

Ethan Allen Express
290, 291

Floridian
40(2), 41(2), 1040, 1041

Gold Runner
701, 702, 703, 704, 710, 711, 712, 713, 714, 715, 716, 717, 718, 719, 1701, 1702, 1703, 1704, 1710, 1711, 1712, 1715, 1716, 1717, 1718, 1719

Heartland Flyer
821, 822

Hiawatha
329, 330, 331, 332, 334, 335, 336, 337, 338, 339, 341, 342, 343

Illinois Zephyr/Carl Sandburg
380, 381, 382, 383

Keystone
110, 123, 600, 601, 605, 607, 609, 610, 611, 612, 615, 620, 622, 623, 624, 626, 637, 639, 640, 641, 642, 643, 644, 645, 646, 647, 648, 649, 650, 651, 652, 653, 654, 655, 656, 657, 658, 660, 661, 662, 663, 664, 665, 666, 667, 669, 670, 671, 672, 674

Lake Shore Limited
48, 49, 448, 449

Lincoln Service Missouri River Runner
318, 319

Lincoln Service/Illinois Service
301, 302, 305, 307

Lincoln Service/Michigan Service
300, 306

Maple Leaf
63, 64

Mardi Gras Service
23, 24, 25, 26

Missouri River Runner
311, 316

Northeast Regional
65, 66, 67, 82, 84, 85, 86, 87, 88, 93, 94, 95, 96, 99, 101, 103, 104, 106, 108, 109, 111, 112, 113, 114, 116, 118, 119, 120, 121, 122, 124, 125, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 181, 182, 183, 184, 185, 186, 189, 190, 192, 193, 194, 195, 197, 198, 400, 405, 409, 416, 417, 425, 426, 450, 460, 461, 463, 464, 465, 467, 470, 471, 473, 474, 475, 478, 479, 486, 488, 490, 494, 495, 497, 499, 630, 632, 636, 1108, 1161, 1175, 1194, 2107, 2117

Northest Regional
100, 102, 117, 126, 196, 199, 627, 631, 806, 807, 887, 888, 1195

Pacific Surfliner
562, 564, 566, 567, 572, 573, 577, 579, 580, 581, 582, 584, 586, 587, 588, 591, 593, 595, 757, 761, 765, 769, 770, 774, 777, 779, 782, 784, 785, 786, 790, 791, 794, 1562, 1591, 1595, 1765, 1769, 1770, 1774, 1777, 1784, 1785, 1790

Palmetto
90

Pennsylvanian
42, 43

Pere Marquette/Michigan Service
370, 371

Saluki/Illinois Service
390, 391, 392, 393

Silver Meteor
97(2), 98, 1098

Silver Service / Palmetto
89

Southwest Chief
3, 4(2), 1003, 1004

Sunset Limited
1, 2

Texas Eagle
21, 22(2), 1021, 1022

Vermonter
54, 55, 56, 57, 107

Winter Park Express
1105, 1106

Wolverine/Michigan Service
350, 351, 352, 353, 354, 355, 1354


FY2026 Q1 OTP Data — paste this block into your prompt:
Seed the following into corridors.on_time_percent and corridors.service_quality_summary.
Data source: FRA Quarterly Report FY2026 Q1 (Oct–Dec 2025). OTP = % customers arriving
within 15 minutes of scheduled arrival.

NEC / Northeast Corridor service line overall: 78%
Long Distance service line overall: 66%
State Supported service line overall: 80%

Route-level OTP (FY2026 Q1):
Acela: 76%
Northeast Regional (On Spine): 80%
Richmond/Newport News/Norfolk: 76%
Roanoke: 82%
Springfield Shuttles: 88%
Adirondack: 63%
Blue Water: 69%
Borealis: 74%
Capitol Corridor: 89%
Carl Sandburg / Illinois Zephyr: 73%
Carolinian: 71%
Cascades (Amtrak Cascades): 76%
Downeaster: 86%
Ethan Allen Express: 70%
Gold Runner (formerly San Joaquins): 63%
Heartland Flyer: 57%
Hiawatha: 87%
Illini / Saluki: 72%
Keystone: 90%
Lincoln / Missouri (combined): 51%
Lincoln Service (Illinois portion): 68%
Missouri (Missouri portion): 77%
Maple Leaf: 59%
Mardi Gras Service: 88%
New York - Albany: 86%
New York - Niagara Falls: 83%
Pacific Surfliner: 85%
Pennsylvanian: 89%
Pere Marquette: 82%
Piedmont: 76%
Vermonter: 72%
Wolverine: 66%
Auto Train: 83%
California Zephyr: 71%
Cardinal: 56%
City of New Orleans: 78%
Coast Starlight: 68%
Crescent: 76%
Empire Builder: 61%
Floridian: 44%
Lake Shore Limited: 78%
Palmetto: 83%
Silver Meteor: 67%
Southwest Chief: 40%
Sunset Limited: 75%
Texas Eagle: 62%
SQL
============================================================
-- AMTRAK ROUTE STOPS SEED DATA
-- Source: FRA FY2026 Q1 Performance & Service Quality Report
-- Oct 1 – Dec 31, 2025
-- ============================================================
-- Format: INSERT INTO stops (corridor_id, name, station_code, sort_order)
-- corridor_id references must match your corridors seed.
-- sort_order reflects geographic sequence along route.
-- ============================================================

-- ============================================================
-- NORTHEAST CORRIDOR
-- ============================================================

-- ACELA (corridor_id = 1)
-- Boston → Washington DC
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(1, 'Boston (South Station), MA', 'BOS', 10),
(1, 'Boston (Back Bay Station), MA', 'BBY', 20),
(1, 'Route 128 (Westwood), MA', 'RTE', 30),
(1, 'Providence, RI', 'PVD', 40),
(1, 'New Haven (Union Station), CT', 'NHV', 50),
(1, 'Stamford, CT', 'STM', 60),
(1, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 70),
(1, 'Newark (Penn Station), NJ', 'NWK', 80),
(1, 'Metropark (Iselin), NJ', 'MET', 90),
(1, 'Philadelphia (30th St Station), PA', 'PHL', 100),
(1, 'Wilmington, DE', 'WIL', 110),
(1, 'Baltimore (Penn Station), MD', 'BAL', 120),
(1, 'BWI Thurgood Marshall Airport Station, MD', 'BWI', 130),
(1, 'Washington, DC', 'WAS', 140);

-- NORTHEAST REGIONAL - ON SPINE (corridor_id = 2)
-- Boston → Washington DC (with additional stops vs Acela)
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(2, 'Boston (South Station), MA', 'BOS', 10),
(2, 'Boston (Back Bay Station), MA', 'BBY', 20),
(2, 'Route 128, MA', 'RTE', 30),
(2, 'Providence, RI', 'PVD', 40),
(2, 'Kingston, RI', 'KIN', 50),
(2, 'Westerly, RI', 'WLY', 60),
(2, 'Mystic, CT', 'MYS', 70),
(2, 'New London, CT', 'NLC', 80),
(2, 'Old Saybrook, CT', 'OSB', 90),
(2, 'New Haven (State St Station), CT', 'STS', 95),
(2, 'New Haven (Union Station), CT', 'NHV', 100),
(2, 'Bridgeport, CT', 'BRP', 110),
(2, 'Stamford, CT', 'STM', 120),
(2, 'New Rochelle, NY', 'NRO', 130),
(2, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 140),
(2, 'Newark (Penn Station), NJ', 'NWK', 150),
(2, 'Newark Liberty International Airport, NJ', 'EWR', 160),
(2, 'Metropark (Iselin), NJ', 'MET', 170),
(2, 'New Brunswick, NJ', 'NBK', 180),
(2, 'Princeton Junction, NJ', 'PJC', 190),
(2, 'Trenton, NJ', 'TRE', 200),
(2, 'Philadelphia (30th St Station), PA', 'PHL', 210),
(2, 'Wilmington, DE', 'WIL', 220),
(2, 'Newark, DE', 'NRK', 230),
(2, 'Aberdeen, MD', 'ABE', 240),
(2, 'Baltimore (Penn Station), MD', 'BAL', 250),
(2, 'BWI Thurgood Marshall Airport Station, MD', 'BWI', 260),
(2, 'New Carrollton, MD', 'NCR', 270),
(2, 'Washington, DC', 'WAS', 280),
(2, 'Alexandria, VA', 'ALX', 290),
(2, 'Woodbridge, VA', 'WDB', 300),
(2, 'Quantico, VA', 'QAN', 310),
(2, 'Fredericksburg, VA', 'FBG', 320),
(2, 'Ashland, VA', 'ASD', 330),
(2, 'Richmond (Staples Mill Rd), VA', 'RVR', 340),
(2, 'Newport News, VA', 'NPN', 350),
(2, 'Williamsburg, VA', 'WBG', 360),
(2, 'Richmond, VA', 'RVM', 370),
(2, 'Petersburg, VA', 'PTB', 380),
(2, 'Norfolk, VA', 'NFK', 390),
(2, 'Springfield, MA', 'SPG', 400),
(2, 'Windsor Locks, CT', 'WNL', 410),
(2, 'Windsor, CT', 'WND', 420),
(2, 'Hartford, CT', 'HFD', 430),
(2, 'Berlin, CT', 'BER', 440),
(2, 'Meriden, CT', 'MDN', 450),
(2, 'Wallingford, CT', 'WFD', 460);

-- RICHMOND / NEWPORT NEWS / NORFOLK (corridor_id = 3)
-- Same spine as NEC Regional plus Virginia extensions - shares stops with NEC Regional above
-- Roanoke branch (corridor_id = 4)
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(4, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(4, 'Newark (Penn Station), NJ', 'NWK', 20),
(4, 'Trenton, NJ', 'TRE', 30),
(4, 'Philadelphia (30th St Station), PA', 'PHL', 40),
(4, 'Wilmington, DE', 'WIL', 50),
(4, 'Baltimore (Penn Station), MD', 'BAL', 60),
(4, 'Washington, DC', 'WAS', 70),
(4, 'Alexandria, VA', 'ALX', 80),
(4, 'Burke Centre, VA', 'BCV', 90),
(4, 'Manassas, VA', 'MSS', 100),
(4, 'Culpeper, VA', 'CLP', 110),
(4, 'Charlottesville, VA', 'CVS', 120),
(4, 'Lynchburg, VA', 'LYH', 130),
(4, 'Roanoke, VA', 'RNK', 140);

-- SPRINGFIELD SHUTTLES (corridor_id = 5)
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(5, 'New Haven (Union Station), CT', 'NHV', 10),
(5, 'New Haven (State Street Station), CT', 'STS', 20),
(5, 'Wallingford, CT', 'WFD', 30),
(5, 'Meriden, CT', 'MDN', 40),
(5, 'Berlin, CT', 'BER', 50),
(5, 'Hartford, CT', 'HFD', 60),
(5, 'Windsor, CT', 'WND', 70),
(5, 'Windsor Locks, CT', 'WNL', 80),
(5, 'Springfield, MA', 'SPG', 90),
(5, 'Holyoke, MA', 'HLK', 100),
(5, 'Northampton, MA', 'NHT', 110),
(5, 'Greenfield, MA', 'GFD', 120);

-- ============================================================
-- STATE SUPPORTED
-- ============================================================

-- ADIRONDACK (corridor_id = 6)
-- New York → Montreal
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(6, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(6, 'Yonkers, NY', 'YNY', 20),
(6, 'Croton-Harmon, NY', 'CRT', 30),
(6, 'Poughkeepsie, NY', 'POU', 40),
(6, 'Rhinecliff, NY', 'RHI', 50),
(6, 'Hudson, NY', 'HUD', 60),
(6, 'Albany-Rensselaer, NY', 'ALB', 70),
(6, 'Schenectady, NY', 'SDY', 80),
(6, 'Saratoga Springs, NY', 'SAR', 90),
(6, 'Fort Edward-Glens Falls, NY', 'FED', 100),
(6, 'Whitehall, NY', 'WHL', 110),
(6, 'Ticonderoga, NY', 'FTC', 120),
(6, 'Port Henry, NY', 'POH', 130),
(6, 'Westport, NY', 'WSP', 140),
(6, 'Port Kent, NY', 'PRK', 150),
(6, 'Plattsburgh, NY', 'PLB', 160),
(6, 'Rouses Point, NY', 'RSP', 170),
(6, 'Saint-Lambert, Quebec, Canada', 'SLQ', 180),
(6, 'Montreal, Quebec, Canada', 'MTR', 190);

-- BLUE WATER (corridor_id = 7)
-- Port Huron, MI → Chicago
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(7, 'Port Huron, MI', 'PTH', 10),
(7, 'Lapeer, MI', 'LPE', 20),
(7, 'Flint, MI', 'FLN', 30),
(7, 'Durand, MI', 'DRD', 40),
(7, 'East Lansing, MI', 'LNS', 50),
(7, 'Battle Creek, MI', 'BTL', 60),
(7, 'Kalamazoo, MI', 'KAL', 70),
(7, 'Dowagiac, MI', 'DOA', 80),
(7, 'Niles, MI', 'NLS', 90),
(7, 'New Buffalo, MI', 'NBU', 100),
(7, 'Chicago (Union Station), IL', 'CHI', 110);

-- BOREALIS (corridor_id = 8)
-- Chicago → St. Paul-Minneapolis
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(8, 'Chicago (Union Station), IL', 'CHI', 10),
(8, 'Glenview, IL', 'GLN', 20),
(8, 'Sturtevant, WI', 'SVT', 30),
(8, 'Milwaukee Airport, WI', 'MKA', 40),
(8, 'Milwaukee, WI', 'MKE', 50),
(8, 'Columbus, WI', 'CBS', 60),
(8, 'Portage, WI', 'POG', 70),
(8, 'Wisconsin Dells, WI', 'WDL', 80),
(8, 'Tomah, WI', 'TOH', 90),
(8, 'La Crosse, WI', 'LSE', 100),
(8, 'Winona, MN', 'WIN', 110),
(8, 'Red Wing, MN', 'RDW', 120),
(8, 'St. Paul-Minneapolis, MN', 'MSP', 130);

-- CAPITOL CORRIDOR (corridor_id = 9)
-- Auburn, CA → San Jose, CA
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(9, 'Auburn, CA', 'ARN', 10),
(9, 'Rocklin, CA', 'RLN', 20),
(9, 'Roseville, CA', 'RSV', 30),
(9, 'Sacramento, CA', 'SAC', 40),
(9, 'Davis, CA', 'DAV', 50),
(9, 'Fairfield-Vacaville, CA', 'FFV', 60),
(9, 'Suisun-Fairfield, CA', 'SUI', 70),
(9, 'Martinez, CA', 'MTZ', 80),
(9, 'Richmond, CA', 'RIC', 90),
(9, 'Berkeley, CA', 'BKY', 100),
(9, 'Emeryville, CA', 'EMY', 110),
(9, 'Oakland (Jack London Square), CA', 'OKJ', 120),
(9, 'Oakland (Coliseum/Airport), CA', 'OAC', 130),
(9, 'Hayward, CA', 'HAY', 140),
(9, 'Fremont (Capitol Trains), CA', 'FMT', 150),
(9, 'Santa Clara (Great America), CA', 'GAC', 160),
(9, 'Santa Clara (Transit Center), CA', 'SCC', 170),
(9, 'San Jose, CA', 'SJC', 180);

-- CARL SANDBURG / ILLINOIS ZEPHYR (corridor_id = 10)
-- Chicago → Quincy, IL
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(10, 'Chicago (Union Station), IL', 'CHI', 10),
(10, 'La Grange, IL', 'LAG', 20),
(10, 'Naperville, IL', 'NPV', 30),
(10, 'Plano, IL', 'PLO', 40),
(10, 'Mendota, IL', 'MDT', 50),
(10, 'Princeton, IL', 'PCT', 60),
(10, 'Kewanee, IL', 'KEE', 70),
(10, 'Galesburg, IL', 'GBB', 80),
(10, 'Macomb, IL', 'MAC', 90),
(10, 'Quincy, IL', 'QCY', 100);

-- CAROLINIAN (corridor_id = 11)
-- New York → Charlotte, NC
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(11, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(11, 'Newark (Penn Station), NJ', 'NWK', 20),
(11, 'Trenton, NJ', 'TRE', 30),
(11, 'Philadelphia (30th St Station), PA', 'PHL', 40),
(11, 'Wilmington, DE', 'WIL', 50),
(11, 'Baltimore (Penn Station), MD', 'BAL', 60),
(11, 'Washington, DC', 'WAS', 70),
(11, 'Alexandria, VA', 'ALX', 80),
(11, 'Quantico, VA', 'QAN', 90),
(11, 'Fredericksburg, VA', 'FBG', 100),
(11, 'Richmond (Staples Mill Rd), VA', 'RVR', 110),
(11, 'Petersburg, VA', 'PTB', 120),
(11, 'Rocky Mount, NC', 'RMT', 130),
(11, 'Wilson, NC', 'WLN', 140),
(11, 'Selma, NC', 'SSM', 150),
(11, 'Raleigh, NC', 'RGH', 160),
(11, 'North Carolina State Fair, NC (Seasonal)', 'NSF', 165),
(11, 'Cary, NC', 'CYN', 170),
(11, 'Durham, NC', 'DNC', 180),
(11, 'Burlington, NC', 'BNC', 190),
(11, 'Greensboro, NC', 'GRO', 200),
(11, 'High Point, NC', 'HPT', 210),
(11, 'Salisbury, NC', 'SAL', 220),
(11, 'Kannapolis, NC', 'KAN', 230),
(11, 'Charlotte, NC', 'CLT', 240);

-- AMTRAK CASCADES (corridor_id = 12)
-- Vancouver, BC → Eugene, OR
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(12, 'Vancouver, British Columbia, Canada', 'VAC', 10),
(12, 'Bellingham, WA', 'BEL', 20),
(12, 'Mount Vernon, WA', 'MVW', 30),
(12, 'Stanwood, WA', 'STW', 40),
(12, 'Everett, WA', 'EVR', 50),
(12, 'Edmonds, WA', 'EDM', 60),
(12, 'Seattle (King Street Station), WA', 'SEA', 70),
(12, 'Tukwila, WA', 'TUK', 80),
(12, 'Tacoma, WA', 'TAC', 90),
(12, 'Olympia-Lacey, WA', 'OLW', 100),
(12, 'Centralia, WA', 'CTL', 110),
(12, 'Kelso-Longview, WA', 'KEL', 120),
(12, 'Vancouver, WA', 'VAN', 130),
(12, 'Portland (Union Station), OR', 'PDX', 140),
(12, 'Oregon City, OR', 'ORC', 150),
(12, 'Salem, OR', 'SLM', 160),
(12, 'Albany, OR', 'ALY', 170),
(12, 'Eugene, OR', 'EUG', 180);

-- DOWNEASTER (corridor_id = 13)
-- Boston → Brunswick, ME
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(13, 'Boston (North Station), MA', 'BON', 10),
(13, 'Woburn, MA', 'WOB', 20),
(13, 'Haverhill, MA', 'HHL', 30),
(13, 'Exeter, NH', 'EXR', 40),
(13, 'Durham, NH', 'DHM', 50),
(13, 'Dover, NH', 'DOV', 60),
(13, 'Wells, ME', 'WEM', 70),
(13, 'Saco, ME', 'SAO', 80),
(13, 'Old Orchard Beach, ME (Seasonal)', 'ORB', 85),
(13, 'Portland, ME', 'POR', 90),
(13, 'Freeport, ME', 'FRE', 100),
(13, 'Brunswick, ME', 'BRK', 110);

-- ETHAN ALLEN EXPRESS (corridor_id = 14)
-- New York → Burlington, VT
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(14, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(14, 'Yonkers, NY', 'YNY', 20),
(14, 'Croton-Harmon, NY', 'CRT', 30),
(14, 'Poughkeepsie, NY', 'POU', 40),
(14, 'Rhinecliff, NY', 'RHI', 50),
(14, 'Hudson, NY', 'HUD', 60),
(14, 'Albany-Rensselaer, NY', 'ALB', 70),
(14, 'Fort Edward-Glens Falls, NY', 'FED', 80),
(14, 'Saratoga Springs, NY', 'SAR', 90),
(14, 'Schenectady, NY', 'SDY', 100),
(14, 'Castleton, VT', 'CNV', 110),
(14, 'Rutland, VT', 'RUD', 120),
(14, 'Middlebury, VT', 'MBY', 130),
(14, 'Ferrisburgh-Vergennes, VT', 'VRN', 140),
(14, 'Burlington (Union Station), VT', 'BTN', 150);

-- GOLD RUNNER (formerly San Joaquins) (corridor_id = 15)
-- Oakland/Emeryville → Bakersfield, CA
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(15, 'Oakland (Jack London Square), CA', 'OKJ', 10),
(15, 'Oakland (Coliseum/Airport), CA', 'OAC', 20),
(15, 'Emeryville, CA', 'EMY', 30),
(15, 'Richmond, CA', 'RIC', 40),
(15, 'Martinez, CA', 'MTZ', 50),
(15, 'Antioch-Pittsburg, CA', 'ACA', 60),
(15, 'Sacramento, CA', 'SAC', 70),
(15, 'Lodi, CA', 'LOD', 80),
(15, 'Stockton (San Joaquin Street), CA', 'SKN', 90),
(15, 'Stockton (Channel Street), CA', 'SKT', 100),
(15, 'Turlock-Denair, CA', 'TRK', 110),
(15, 'Modesto, CA', 'MOD', 120),
(15, 'Merced, CA', 'MCD', 130),
(15, 'Madera, CA', 'MDR', 140),
(15, 'Fresno, CA', 'FNO', 150),
(15, 'Hanford, CA', 'HNF', 160),
(15, 'Colonel Allensworth State Park, CA (Seasonal)', 'CNL', 165),
(15, 'Corcoran, CA', 'COC', 170),
(15, 'Wasco, CA', 'WAC', 180),
(15, 'Bakersfield, CA', 'BFD', 190);

-- HEARTLAND FLYER (corridor_id = 16)
-- Fort Worth, TX → Oklahoma City, OK
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(16, 'Fort Worth, TX', 'FTW', 10),
(16, 'Gainesville, TX', 'GLE', 20),
(16, 'Ardmore, OK', 'ADM', 30),
(16, 'Pauls Valley, OK', 'PVL', 40),
(16, 'Purcell, OK', 'PUR', 50),
(16, 'Norman, OK', 'NOR', 60),
(16, 'Oklahoma City, OK', 'OKC', 70);

-- HIAWATHA (corridor_id = 17)
-- Chicago → Milwaukee
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(17, 'Chicago (Union Station), IL', 'CHI', 10),
(17, 'Glenview, IL', 'GLN', 20),
(17, 'Sturtevant, WI', 'SVT', 30),
(17, 'Milwaukee Airport, WI', 'MKA', 40),
(17, 'Milwaukee (Downtown), WI', 'MKE', 50);

-- ILLINI / SALUKI (corridor_id = 18)
-- Chicago → Carbondale, IL
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(18, 'Chicago (Union Station), IL', 'CHI', 10),
(18, 'Homewood, IL', 'HMW', 20),
(18, 'Kankakee, IL', 'KKI', 30),
(18, 'Gilman, IL', 'GLM', 40),
(18, 'Rantoul, IL', 'RTL', 50),
(18, 'Champaign-Urbana, IL', 'CHM', 60),
(18, 'Mattoon, IL', 'MAT', 70),
(18, 'Effingham, IL', 'EFG', 80),
(18, 'Centralia, IL', 'CEN', 90),
(18, 'Du Quoin, IL', 'DQN', 100),
(18, 'Carbondale, IL', 'CDL', 110);

-- KEYSTONE (corridor_id = 19)
-- New York → Harrisburg, PA
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(19, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(19, 'Newark (Penn Station), NJ', 'NWK', 20),
(19, 'Newark Liberty International Airport, NJ', 'EWR', 30),
(19, 'Metropark, NJ', 'MET', 40),
(19, 'New Brunswick, NJ', 'NBK', 50),
(19, 'Princeton Junction, NJ', 'PJC', 60),
(19, 'Trenton, NJ', 'TRE', 70),
(19, 'Cornwells Heights, PA', 'CWH', 80),
(19, 'North Philadelphia, PA', 'PHN', 90),
(19, 'Philadelphia (30th St Station), PA', 'PHL', 100),
(19, 'Ardmore, PA', 'ARD', 110),
(19, 'Paoli, PA', 'PAO', 120),
(19, 'Exton, PA', 'EXT', 130),
(19, 'Downingtown, PA', 'DOW', 140),
(19, 'Coatesville, PA', 'COT', 150),
(19, 'Parkesburg, PA', 'PAR', 160),
(19, 'Lancaster, PA', 'LNC', 170),
(19, 'Mount Joy, PA', 'MJY', 180),
(19, 'Elizabethtown, PA', 'ELT', 190),
(19, 'Middletown, PA', 'MID', 200),
(19, 'Harrisburg, PA', 'HAR', 210);

-- LINCOLN / MISSOURI (combined corridor, corridor_id = 20)
-- Chicago → Kansas City, MO
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(20, 'Chicago (Union Station), IL', 'CHI', 10),
(20, 'Summit, IL', 'SMT', 20),
(20, 'Joliet, IL', 'JOL', 30),
(20, 'Dwight, IL', 'DWT', 40),
(20, 'Pontiac, IL', 'PON', 50),
(20, 'Bloomington-Normal, IL', 'BNL', 60),
(20, 'Lincoln, IL', 'LCN', 70),
(20, 'Springfield, IL', 'SPI', 80),
(20, 'Carlinville, IL', 'CRV', 90),
(20, 'Alton, IL', 'ALN', 100),
(20, 'St. Louis, MO', 'STL', 110),
(20, 'Kirkwood, MO', 'KWD', 120),
(20, 'Washington, MO', 'WAH', 130),
(20, 'Hermann, MO', 'HEM', 140),
(20, 'Jefferson City, MO', 'JEF', 150),
(20, 'Sedalia, MO', 'SED', 160),
(20, 'Warrensburg, MO', 'WAR', 170),
(20, 'Lee''s Summit, MO', 'LEE', 180),
(20, 'Independence, MO', 'IDP', 190),
(20, 'Kansas City (Union Station), MO', 'KCY', 200);

-- MAPLE LEAF (corridor_id = 21)
-- New York → Toronto
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(21, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(21, 'Yonkers, NY', 'YNY', 20),
(21, 'Croton-Harmon, NY', 'CRT', 30),
(21, 'Poughkeepsie, NY', 'POU', 40),
(21, 'Rhinecliff, NY', 'RHI', 50),
(21, 'Hudson, NY', 'HUD', 60),
(21, 'Albany-Rensselaer, NY', 'ALB', 70),
(21, 'Schenectady, NY', 'SDY', 80),
(21, 'Amsterdam, NY', 'AMS', 90),
(21, 'Utica, NY', 'UCA', 100),
(21, 'Rome, NY', 'ROM', 110),
(21, 'Syracuse, NY', 'SYR', 120),
(21, 'New York State Fair, NY (Seasonal)', 'NYF', 125),
(21, 'Rochester, NY', 'ROC', 130),
(21, 'Buffalo-Depew, NY', 'BUF', 140),
(21, 'Buffalo, NY', 'BFX', 150),
(21, 'Niagara Falls, NY', 'NFL', 160),
(21, 'Canadian Border NY', 'CBN', 165),
(21, 'Niagara Falls, Ontario, Canada', 'NFS', 170),
(21, 'St. Catharines, Ontario, Canada', 'SCA', 180),
(21, 'Grimsby, Ontario, Canada', 'GMS', 190),
(21, 'Aldershot, Ontario, Canada', 'AST', 200),
(21, 'Oakville, Ontario, Canada', 'OKL', 210),
(21, 'Toronto Union, Ontario, Canada', 'TWO', 220);

-- MARDI GRAS SERVICE (corridor_id = 22)
-- New Orleans → Mobile, AL
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(22, 'New Orleans, LA', 'NOL', 10),
(22, 'Bay Saint Louis, MS', 'BAS', 20),
(22, 'Gulfport, MS', 'GUF', 30),
(22, 'Biloxi, MS', 'BIX', 40),
(22, 'Pascagoula, MS', 'PAG', 50),
(22, 'Mobile, AL', 'MOE', 60);

-- NEW YORK - ALBANY (corridor_id = 23)
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(23, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(23, 'Yonkers, NY', 'YNY', 20),
(23, 'Croton-Harmon, NY', 'CRT', 30),
(23, 'Poughkeepsie, NY', 'POU', 40),
(23, 'Rhinecliff, NY', 'RHI', 50),
(23, 'Hudson, NY', 'HUD', 60),
(23, 'Albany-Rensselaer, NY', 'ALB', 70);

-- NEW YORK - NIAGARA FALLS (corridor_id = 24)
-- (Empire Service)
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(24, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(24, 'Yonkers, NY', 'YNY', 20),
(24, 'Croton-Harmon, NY', 'CRT', 30),
(24, 'Poughkeepsie, NY', 'POU', 40),
(24, 'Rhinecliff, NY', 'RHI', 50),
(24, 'Hudson, NY', 'HUD', 60),
(24, 'Albany-Rensselaer, NY', 'ALB', 70),
(24, 'Schenectady, NY', 'SDY', 80),
(24, 'Amsterdam, NY', 'AMS', 90),
(24, 'Utica, NY', 'UCA', 100),
(24, 'Rome, NY', 'ROM', 110),
(24, 'Syracuse, NY', 'SYR', 120),
(24, 'New York State Fair, NY (Seasonal)', 'NYF', 125),
(24, 'Rochester, NY', 'ROC', 130),
(24, 'Buffalo-Depew, NY', 'BUF', 140),
(24, 'Buffalo, NY', 'BFX', 150),
(24, 'Niagara Falls, NY', 'NFL', 160);

-- PACIFIC SURFLINER (corridor_id = 25)
-- San Luis Obispo → San Diego
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(25, 'San Luis Obispo, CA', 'SLO', 10),
(25, 'Grover Beach, CA', 'GVB', 20),
(25, 'Guadalupe-Santa Maria, CA', 'GUA', 30),
(25, 'Lompoc-Surf, CA', 'LPS', 40),
(25, 'Goleta, CA', 'GTA', 50),
(25, 'Santa Barbara, CA', 'SBA', 60),
(25, 'Carpinteria, CA', 'CPN', 70),
(25, 'Ventura, CA', 'VEC', 80),
(25, 'Oxnard, CA', 'OXN', 90),
(25, 'Camarillo, CA', 'CML', 100),
(25, 'Moorpark, CA', 'MPK', 110),
(25, 'Simi Valley, CA', 'SIM', 120),
(25, 'Chatsworth, CA', 'CWT', 130),
(25, 'Northridge Station', 'NRG', 140),
(25, 'Van Nuys, CA', 'VNC', 150),
(25, 'Burbank (Airport), CA', 'BUR', 160),
(25, 'Burbank, CA', 'BBK', 170),
(25, 'Glendale, CA', 'GDL', 180),
(25, 'Los Angeles (Union Station), CA', 'LAX', 190),
(25, 'Fullerton, CA', 'FUL', 200),
(25, 'Anaheim, CA', 'ANA', 210),
(25, 'Santa Ana, CA', 'SNA', 220),
(25, 'Irvine, CA', 'IRV', 230),
(25, 'San Juan Capistrano, CA', 'SNC', 240),
(25, 'San Clemente Pier, CA', 'SNP', 250),
(25, 'Oceanside, CA', 'OSD', 260),
(25, 'Solana Beach, CA', 'SOL', 270),
(25, 'San Diego (Old Town), CA', 'OLT', 280),
(25, 'San Diego (Downtown), CA', 'SAN', 290);

-- PENNSYLVANIAN (corridor_id = 26)
-- New York → Pittsburgh
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(26, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(26, 'Newark (Penn Station), NJ', 'NWK', 20),
(26, 'Trenton, NJ', 'TRE', 30),
(26, 'Philadelphia (30th St Station), PA', 'PHL', 40),
(26, 'Paoli, PA', 'PAO', 50),
(26, 'Exton, PA', 'EXT', 60),
(26, 'Elizabethtown, PA', 'ELT', 70),
(26, 'Lancaster, PA', 'LNC', 80),
(26, 'Harrisburg, PA', 'HAR', 90),
(26, 'Lewistown, PA', 'LEW', 100),
(26, 'Huntingdon, PA', 'HGD', 110),
(26, 'Tyrone, PA', 'TYR', 120),
(26, 'Altoona, PA', 'ALT', 130),
(26, 'Johnstown, PA', 'JST', 140),
(26, 'Latrobe, PA', 'LAB', 150),
(26, 'Greensburg, PA', 'GNB', 160),
(26, 'Pittsburgh (Union Station), PA', 'PGH', 170);

-- PERE MARQUETTE (corridor_id = 27)
-- Chicago → Grand Rapids, MI
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(27, 'Chicago (Union Station), IL', 'CHI', 10),
(27, 'St. Joseph, MI', 'SJM', 20),
(27, 'Bangor, MI', 'BAM', 30),
(27, 'Holland, MI', 'HOM', 40),
(27, 'Grand Rapids, MI', 'GRR', 50);

-- PIEDMONT (corridor_id = 28)
-- Raleigh → Charlotte, NC
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(28, 'Raleigh, NC', 'RGH', 10),
(28, 'Cary, NC', 'CYN', 20),
(28, 'Durham, NC', 'DNC', 30),
(28, 'Burlington, NC', 'BNC', 40),
(28, 'Greensboro, NC', 'GRO', 50),
(28, 'High Point, NC', 'HPT', 60),
(28, 'Lexington, NC', 'LEX', 70),
(28, 'Salisbury, NC', 'SAL', 80),
(28, 'Kannapolis, NC', 'KAN', 90),
(28, 'Charlotte, NC', 'CLT', 100),
(28, 'North Carolina State Fair, NC (Seasonal)', 'NSF', 105);

-- VERMONTER (corridor_id = 29)
-- St. Albans, VT → Washington, DC
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(29, 'St. Albans, VT', 'SAB', 10),
(29, 'Essex Junction, VT', 'ESX', 20),
(29, 'Waterbury, VT', 'WAB', 30),
(29, 'Montpelier-Berlin, VT', 'MPR', 40),
(29, 'Randolph, VT', 'RPH', 50),
(29, 'White River Junction, VT', 'WRJ', 60),
(29, 'Windsor, VT', 'WNM', 70),
(29, 'Claremont, NH', 'CLA', 80),
(29, 'Bellows Falls, VT', 'BLF', 90),
(29, 'Brattleboro, VT', 'BRA', 100),
(29, 'Greenfield, MA', 'GFD', 110),
(29, 'Northampton, MA', 'NHT', 120),
(29, 'Holyoke, MA', 'HLK', 130),
(29, 'Springfield, MA', 'SPG', 140),
(29, 'Windsor Locks, CT', 'WNL', 150),
(29, 'Hartford, CT', 'HFD', 160),
(29, 'Meriden, CT', 'MDN', 170),
(29, 'New Haven (Union Station), CT', 'NHV', 180),
(29, 'Bridgeport, CT', 'BRP', 190),
(29, 'Stamford, CT', 'STM', 200),
(29, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 210),
(29, 'Newark (Penn Station), NJ', 'NWK', 220),
(29, 'Metropark (Iselin), NJ', 'MET', 230),
(29, 'Trenton, NJ', 'TRE', 240),
(29, 'Philadelphia (30th St Station), PA', 'PHL', 250),
(29, 'Wilmington, DE', 'WIL', 260),
(29, 'Baltimore (Penn Station), MD', 'BAL', 270),
(29, 'BWI Thurgood Marshall Airport Station, MD', 'BWI', 280),
(29, 'New Carrollton, MD', 'NCR', 290),
(29, 'Washington, DC', 'WAS', 300);

-- WOLVERINE (corridor_id = 30)
-- Pontiac, MI → Chicago
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(30, 'Pontiac, MI', 'PNT', 10),
(30, 'Troy, MI', 'TRM', 20),
(30, 'Royal Oak, MI', 'ROY', 30),
(30, 'Detroit, MI', 'DET', 40),
(30, 'Dearborn, MI', 'DER', 50),
(30, 'Ann Arbor, MI', 'ARB', 60),
(30, 'Jackson, MI', 'JXN', 70),
(30, 'Albion, MI', 'ALI', 80),
(30, 'Battle Creek, MI', 'BTL', 90),
(30, 'Kalamazoo, MI', 'KAL', 100),
(30, 'Dowagiac, MI', 'DOA', 110),
(30, 'Niles, MI', 'NLS', 120),
(30, 'New Buffalo, MI', 'NBU', 130),
(30, 'Hammond-Whiting, IN', 'HMI', 140),
(30, 'Chicago (Union Station), IL', 'CHI', 150);

-- ============================================================
-- LONG DISTANCE
-- ============================================================

-- AUTO TRAIN (corridor_id = 31)
-- Lorton, VA → Sanford, FL
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(31, 'Lorton (Auto Train), VA', 'LOR', 10),
(31, 'Sanford (Auto Train), FL', 'SFA', 20);

-- CALIFORNIA ZEPHYR (corridor_id = 32)
-- Chicago → Emeryville, CA
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(32, 'Chicago (Union Station), IL', 'CHI', 10),
(32, 'Naperville, IL', 'NPV', 20),
(32, 'Princeton, IL', 'PCT', 30),
(32, 'Galesburg, IL', 'GBB', 40),
(32, 'Burlington, IA', 'BRL', 50),
(32, 'Mount Pleasant, IA', 'MTP', 60),
(32, 'Ottumwa, IA', 'OTM', 70),
(32, 'Osceola, IA', 'OSC', 80),
(32, 'Creston, IA', 'CRN', 90),
(32, 'Omaha, NE', 'OMA', 100),
(32, 'Lincoln, NE', 'LNK', 110),
(32, 'Hastings, NE', 'HAS', 120),
(32, 'Holdrege, NE', 'HLD', 130),
(32, 'McCook, NE', 'MCK', 140),
(32, 'Fort Morgan, CO', 'FMG', 150),
(32, 'Denver (Union Station), CO', 'DEN', 160),
(32, 'Winter Park/Fraser, CO', 'WIP', 170),
(32, 'Granby, CO', 'GRA', 180),
(32, 'Glenwood Springs, CO', 'GSC', 190),
(32, 'Grand Junction, CO', 'GJT', 200),
(32, 'Green River, UT', 'GRI', 210),
(32, 'Helper, UT', 'HER', 220),
(32, 'Provo, UT', 'PRO', 230),
(32, 'Salt Lake City, UT', 'SLC', 240),
(32, 'Elko, NV', 'ELK', 250),
(32, 'Winnemucca, NV', 'WNN', 260),
(32, 'Reno, NV', 'RNO', 270),
(32, 'Truckee, CA', 'TRU', 280),
(32, 'Colfax, CA', 'COX', 290),
(32, 'Roseville, CA', 'RSV', 300),
(32, 'Sacramento, CA', 'SAC', 310),
(32, 'Davis, CA', 'DAV', 320),
(32, 'Martinez, CA', 'MTZ', 330),
(32, 'Richmond, CA', 'RIC', 340),
(32, 'Emeryville, CA', 'EMY', 350);

-- CARDINAL (corridor_id = 33)
-- New York → Chicago (3x weekly)
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(33, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(33, 'Newark (Penn Station), NJ', 'NWK', 20),
(33, 'Trenton, NJ', 'TRE', 30),
(33, 'Philadelphia (30th St Station), PA', 'PHL', 40),
(33, 'Wilmington, DE', 'WIL', 50),
(33, 'Baltimore (Penn Station), MD', 'BAL', 60),
(33, 'Washington, DC', 'WAS', 70),
(33, 'Alexandria, VA', 'ALX', 80),
(33, 'Manassas, VA', 'MSS', 90),
(33, 'Culpeper, VA', 'CLP', 100),
(33, 'Charlottesville, VA', 'CVS', 110),
(33, 'Staunton, VA', 'STA', 120),
(33, 'Clifton Forge, VA', 'CLF', 130),
(33, 'White Sulphur Springs, WV', 'WSS', 140),
(33, 'Alderson, WV', 'ALD', 150),
(33, 'Hinton, WV', 'HIN', 160),
(33, 'Prince, WV', 'PRC', 170),
(33, 'Thurmond, WV', 'THN', 180),
(33, 'Montgomery, WV', 'MNG', 190),
(33, 'Charleston, WV', 'CHW', 200),
(33, 'Huntington, WV', 'HUN', 210),
(33, 'Ashland, KY', 'AKY', 220),
(33, 'South Shore, KY - Portsmouth, OH', 'SPM', 230),
(33, 'Maysville, KY', 'MAY', 240),
(33, 'Cincinnati (Union Terminal), OH', 'CIN', 250),
(33, 'Connersville, IN', 'COI', 260),
(33, 'Indianapolis, IN', 'IND', 270),
(33, 'Crawfordsville, IN', 'CRF', 280),
(33, 'Lafayette, IN', 'LAF', 290),
(33, 'Rensselaer, IN', 'REN', 300),
(33, 'Dyer, IN', 'DYE', 310),
(33, 'Chicago (Union Station), IL', 'CHI', 320);

-- CITY OF NEW ORLEANS (corridor_id = 34)
-- Chicago → New Orleans
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(34, 'Chicago (Union Station), IL', 'CHI', 10),
(34, 'Homewood, IL', 'HMW', 20),
(34, 'Kankakee, IL', 'KKI', 30),
(34, 'Champaign-Urbana, IL', 'CHM', 40),
(34, 'Mattoon, IL', 'MAT', 50),
(34, 'Effingham, IL', 'EFG', 60),
(34, 'Centralia, IL', 'CEN', 70),
(34, 'Carbondale, IL', 'CDL', 80),
(34, 'Fulton, KY', 'FTN', 90),
(34, 'Newbern-Dyersburg, TN', 'NBN', 100),
(34, 'Memphis, TN', 'MEM', 110),
(34, 'Marks, MS', 'MKS', 120),
(34, 'Greenwood, MS', 'GWD', 130),
(34, 'Yazoo City, MS', 'YAZ', 140),
(34, 'Jackson, MS', 'JAN', 150),
(34, 'Hazlehurst, MS', 'HAZ', 160),
(34, 'Brookhaven, MS', 'BRH', 170),
(34, 'McComb, MS', 'MCB', 180),
(34, 'Hammond, LA', 'HMD', 190),
(34, 'New Orleans, LA', 'NOL', 200);

-- COAST STARLIGHT (corridor_id = 35)
-- Seattle → Los Angeles
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(35, 'Seattle (King Street Station), WA', 'SEA', 10),
(35, 'Tacoma, WA', 'TAC', 20),
(35, 'Olympia-Lacey, WA', 'OLW', 30),
(35, 'Centralia, WA', 'CTL', 40),
(35, 'Kelso-Longview, WA', 'KEL', 50),
(35, 'Vancouver, WA', 'VAN', 60),
(35, 'Portland (Union Station), OR', 'PDX', 70),
(35, 'Salem, OR', 'SLM', 80),
(35, 'Albany, OR', 'ALY', 90),
(35, 'Eugene, OR', 'EUG', 100),
(35, 'Chemult, OR', 'CMO', 110),
(35, 'Klamath Falls, OR', 'KFS', 120),
(35, 'Dunsmuir, CA', 'DUN', 130),
(35, 'Redding, CA', 'RDD', 140),
(35, 'Chico, CA', 'CIC', 150),
(35, 'Sacramento, CA', 'SAC', 160),
(35, 'Davis, CA', 'DAV', 170),
(35, 'Martinez, CA', 'MTZ', 180),
(35, 'Emeryville, CA', 'EMY', 190),
(35, 'Oakland (Jack London Square), CA', 'OKJ', 200),
(35, 'San Jose, CA', 'SJC', 210),
(35, 'Salinas, CA', 'SNS', 220),
(35, 'Paso Robles, CA', 'PRB', 230),
(35, 'San Luis Obispo, CA', 'SLO', 240),
(35, 'Santa Barbara, CA', 'SBA', 250),
(35, 'Oxnard, CA', 'OXN', 260),
(35, 'Simi Valley, CA', 'SIM', 270),
(35, 'Van Nuys, CA', 'VNC', 280),
(35, 'Burbank (Airport), CA', 'BUR', 290),
(35, 'Los Angeles (Union Station), CA', 'LAX', 300);

-- CRESCENT (corridor_id = 36)
-- New York → New Orleans
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(36, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(36, 'Newark (Penn Station), NJ', 'NWK', 20),
(36, 'Trenton, NJ', 'TRE', 30),
(36, 'Philadelphia (30th St Station), PA', 'PHL', 40),
(36, 'Wilmington, DE', 'WIL', 50),
(36, 'Baltimore (Penn Station), MD', 'BAL', 60),
(36, 'Washington, DC', 'WAS', 70),
(36, 'Alexandria, VA', 'ALX', 80),
(36, 'Fredericksburg, VA', 'FBG', 90),
(36, 'Richmond (Staples Mill Rd), VA', 'RVR', 100),
(36, 'Petersburg, VA', 'PTB', 110),
(36, 'Rocky Mount, NC', 'RMT', 120),
(36, 'Raleigh, NC', 'RGH', 130),
(36, 'Cary, NC', 'CYN', 140),
(36, 'Southern Pines, NC', 'SOP', 150),
(36, 'Hamlet, NC', 'HAM', 160),
(36, 'Camden, SC', 'CAM', 170),
(36, 'Columbia, SC', 'CLB', 180),
(36, 'Denmark, SC', 'DNK', 190),
(36, 'Savannah, GA', 'SAV', 200),
(36, 'Jacksonville, FL', 'JAX', 210),
(36, 'Palatka, FL', 'PAK', 220),
(36, 'DeLand, FL', 'DLD', 230),
(36, 'Winter Park, FL', 'WPK', 240),
(36, 'Orlando, FL', 'ORL', 250),
(36, 'Kissimmee, FL', 'KIS', 260),
-- Note: above stops are Floridian extension; Crescent proper ends below
(36, 'Gainesville, GA', 'GNS', 270),
(36, 'Toccoa, GA', 'TCA', 280),
(36, 'Clemson, SC', 'CSN', 290),
(36, 'Greenville, SC', 'GRV', 300),
(36, 'Spartanburg, SC', 'SPB', 310),
(36, 'Gastonia, NC', 'GAS', 320),
(36, 'Charlotte, NC', 'CLT', 330),
(36, 'Salisbury, NC', 'SAL', 340),
(36, 'High Point, NC', 'HPT', 350),
(36, 'Greensboro, NC', 'GRO', 360),
(36, 'Danville, VA', 'DAN', 370),
(36, 'Lynchburg, VA', 'LYH', 380),
(36, 'Charlottesville, VA', 'CVS', 390),
(36, 'Culpeper, VA', 'CLP', 400),
(36, 'Manassas, VA', 'MSS', 410),
(36, 'Atlanta, GA', 'ATL', 420),
(36, 'Anniston, AL', 'ATN', 430),
(36, 'Birmingham, AL', 'BHM', 440),
(36, 'Tuscaloosa, AL', 'TCL', 450),
(36, 'Meridian, MS', 'MEI', 460),
(36, 'Laurel, MS', 'LAU', 470),
(36, 'Hattiesburg, MS', 'HBG', 480),
(36, 'Picayune, MS', 'PIC', 490),
(36, 'Slidell, LA', 'SDL', 500),
(36, 'New Orleans, LA', 'NOL', 510);

-- EMPIRE BUILDER (corridor_id = 37)
-- Chicago → Seattle/Portland
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(37, 'Chicago (Union Station), IL', 'CHI', 10),
(37, 'Glenview, IL', 'GLN', 20),
(37, 'Milwaukee, WI', 'MKE', 30),
(37, 'Columbus, WI', 'CBS', 40),
(37, 'Portage, WI', 'POG', 50),
(37, 'Wisconsin Dells, WI', 'WDL', 60),
(37, 'Tomah, WI', 'TOH', 70),
(37, 'La Crosse, WI', 'LSE', 80),
(37, 'Winona, MN', 'WIN', 90),
(37, 'Red Wing, MN', 'RDW', 100),
(37, 'St. Paul-Minneapolis, MN', 'MSP', 110),
(37, 'St. Cloud, MN', 'SCD', 120),
(37, 'Staples, MN', 'SPL', 130),
(37, 'Detroit Lakes, MN', 'DLK', 140),
(37, 'Fargo, ND', 'FAR', 150),
(37, 'Grand Forks, ND', 'GFK', 160),
(37, 'Devils Lake, ND', 'DVL', 170),
(37, 'Rugby, ND', 'RUG', 180),
(37, 'Minot, ND', 'MOT', 190),
(37, 'Stanley, ND', 'STN', 200),
(37, 'Williston, ND', 'WTN', 210),
(37, 'Wolf Point, MT', 'WPT', 220),
(37, 'Glasgow, MT', 'GGW', 230),
(37, 'Malta, MT', 'MAL', 240),
(37, 'Havre, MT', 'HAV', 250),
(37, 'Shelby, MT', 'SBY', 260),
(37, 'Cut Bank, MT', 'CUT', 270),
(37, 'Browning, MT', 'BRO', 280),
(37, 'East Glacier Park, MT', 'GPK', 290),
(37, 'Essex, MT', 'ESM', 300),
(37, 'West Glacier, MT', 'WGL', 310),
(37, 'Whitefish, MT', 'WFH', 320),
(37, 'Libby, MT', 'LIB', 330),
(37, 'Sandpoint, ID', 'SPT', 340),
(37, 'Spokane, WA', 'SPK', 350),
-- Seattle branch
(37, 'Ephrata, WA', 'EPH', 360),
(37, 'Wenatchee, WA', 'WEN', 370),
(37, 'Leavenworth, WA', 'LWA', 380),
(37, 'Everett, WA', 'EVR', 390),
(37, 'Edmonds, WA', 'EDM', 400),
(37, 'Seattle (King Street Station), WA', 'SEA', 410),
-- Portland branch
(37, 'Pasco, WA', 'PSC', 420),
(37, 'Wishram, WA', 'WIH', 430),
(37, 'Bingen-White Salmon, WA', 'BNG', 440),
(37, 'Vancouver, WA', 'VAN', 450),
(37, 'Portland, OR', 'PDX', 460);

-- FLORIDIAN (corridor_id = 38) - Temporary route, began Nov 2024
-- Chicago → Miami
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(38, 'Chicago (Union Station), IL', 'CHI', 10),
(38, 'Harpers Ferry, WV', 'HFY', 20),
(38, 'Rockville, MD', 'RKV', 30),
(38, 'Washington, DC', 'WAS', 40),
(38, 'Alexandria, VA', 'ALX', 50),
(38, 'Richmond (Staples Mill Rd), VA', 'RVR', 60),
(38, 'Petersburg, VA', 'PTB', 70),
(38, 'Rocky Mount, NC', 'RMT', 80),
(38, 'Raleigh, NC', 'RGH', 90),
(38, 'Cary, NC', 'CYN', 100),
(38, 'Southern Pines, NC', 'SOP', 110),
(38, 'Hamlet, NC', 'HAM', 120),
(38, 'Camden, SC', 'CAM', 130),
(38, 'Columbia, SC', 'CLB', 140),
(38, 'Denmark, SC', 'DNK', 150),
(38, 'Savannah, GA', 'SAV', 160),
(38, 'Jacksonville, FL', 'JAX', 170),
(38, 'Palatka, FL', 'PAK', 180),
(38, 'DeLand, FL', 'DLD', 190),
(38, 'Winter Park, FL', 'WPK', 200),
(38, 'Orlando, FL', 'ORL', 210),
(38, 'Kissimmee, FL', 'KIS', 220),
(38, 'Lakeland, FL', 'LAK', 230),
(38, 'Tampa, FL', 'TPA', 240),
(38, 'Lakeland, FL', 'LKL', 250),
(38, 'Winter Haven, FL', 'WTH', 260),
(38, 'Sebring, FL', 'SBG', 270),
(38, 'Okeechobee, FL', 'OKE', 280),
(38, 'West Palm Beach, FL', 'WPB', 290),
(38, 'Delray Beach, FL', 'DLB', 300),
(38, 'Deerfield Beach, FL', 'DFB', 310),
(38, 'Fort Lauderdale, FL', 'FTL', 320),
(38, 'Hollywood, FL', 'HOL', 330),
(38, 'Miami, FL', 'MIA', 340),
-- Also includes Pittsburgh via Cumberland route
(38, 'Cumberland, MD', 'CUM', 350),
(38, 'Connellsville, PA', 'COV', 360),
(38, 'Pittsburgh (Union Station), PA', 'PGH', 370),
(38, 'Alliance, OH', 'ALC', 380),
(38, 'Cleveland, OH', 'CLE', 390),
(38, 'Elyria, OH', 'ELY', 400),
(38, 'Sandusky, OH', 'SKY', 410),
(38, 'Toledo, OH', 'TOL', 420),
(38, 'Waterloo, IN', 'WTI', 430),
(38, 'Elkhart, IN', 'EKH', 440),
(38, 'South Bend, IN', 'SOB', 450);

-- LAKE SHORE LIMITED (corridor_id = 39)
-- Chicago → Boston/New York
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(39, 'Chicago (Union Station), IL', 'CHI', 10),
(39, 'South Bend, IN', 'SOB', 20),
(39, 'Elkhart, IN', 'EKH', 30),
(39, 'Waterloo, IN', 'WTI', 40),
(39, 'Bryan, OH', 'BYN', 50),
(39, 'Toledo, OH', 'TOL', 60),
(39, 'Sandusky, OH', 'SKY', 70),
(39, 'Elyria, OH', 'ELY', 80),
(39, 'Cleveland, OH', 'CLE', 90),
(39, 'Erie, PA', 'ERI', 100),
(39, 'Buffalo-Depew, NY', 'BUF', 110),
(39, 'Poughkeepsie, NY', 'POU', 120),
(39, 'Rochester, NY', 'ROC', 130),
(39, 'Syracuse, NY', 'SYR', 140),
(39, 'Utica, NY', 'UCA', 150),
(39, 'Schenectady, NY', 'SDY', 160),
(39, 'Albany-Rensselaer, NY', 'ALB', 170),
(39, 'Rhinecliff, NY', 'RHI', 180),
(39, 'Croton-Harmon, NY', 'CRT', 190),
(39, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 200),
-- Boston branch
(39, 'Springfield, MA', 'SPG', 210),
(39, 'Pittsfield, MA', 'PIT', 220),
(39, 'Worcester, MA', 'WOR', 230),
(39, 'Framingham, MA', 'FRA', 240),
(39, 'Boston (Back Bay Station), MA', 'BBY', 250),
(39, 'Boston (South Station), MA', 'BOS', 260);

-- PALMETTO (corridor_id = 40)
-- New York → Savannah, GA
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(40, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(40, 'Newark (Penn Station), NJ', 'NWK', 20),
(40, 'Trenton, NJ', 'TRE', 30),
(40, 'Philadelphia (30th St Station), PA', 'PHL', 40),
(40, 'Wilmington, DE', 'WIL', 50),
(40, 'Baltimore (Penn Station), MD', 'BAL', 60),
(40, 'Washington, DC', 'WAS', 70),
(40, 'Alexandria, VA', 'ALX', 80),
(40, 'Fredericksburg, VA', 'FBG', 90),
(40, 'Richmond (Staples Mill Rd), VA', 'RVR', 100),
(40, 'Petersburg, VA', 'PTB', 110),
(40, 'Rocky Mount, NC', 'RMT', 120),
(40, 'Fayetteville, NC', 'FAY', 130),
(40, 'Florence, SC', 'FLO', 140),
(40, 'Kingstree, SC', 'KTR', 150),
(40, 'Charleston, SC', 'CHS', 160),
(40, 'Yemassee, SC', 'YEM', 170),
(40, 'Savannah, GA', 'SAV', 180);

-- SILVER METEOR (corridor_id = 41)
-- New York → Miami
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(41, 'NY Moynihan Train Hall at Penn Station, NY', 'NYP', 10),
(41, 'Newark (Penn Station), NJ', 'NWK', 20),
(41, 'Trenton, NJ', 'TRE', 30),
(41, 'Philadelphia (30th St Station), PA', 'PHL', 40),
(41, 'Wilmington, DE', 'WIL', 50),
(41, 'Baltimore (Penn Station), MD', 'BAL', 60),
(41, 'Washington, DC', 'WAS', 70),
(41, 'Alexandria, VA', 'ALX', 80),
(41, 'Fredericksburg, VA', 'FBG', 90),
(41, 'Richmond (Staples Mill Rd), VA', 'RVR', 100),
(41, 'Petersburg, VA', 'PTB', 110),
(41, 'Rocky Mount, NC', 'RMT', 120),
(41, 'Fayetteville, NC', 'FAY', 130),
(41, 'Florence, SC', 'FLO', 140),
(41, 'Kingstree, SC', 'KTR', 150),
(41, 'Charleston, SC', 'CHS', 160),
(41, 'Yemassee, SC', 'YEM', 170),
(41, 'Savannah, GA', 'SAV', 180),
(41, 'Jesup, GA', 'JSP', 190),
(41, 'Jacksonville, FL', 'JAX', 200),
(41, 'Palatka, FL', 'PAK', 210),
(41, 'DeLand, FL', 'DLD', 220),
(41, 'Winter Park, FL', 'WPK', 230),
(41, 'Orlando, FL', 'ORL', 240),
(41, 'Kissimmee, FL', 'KIS', 250),
(41, 'Winter Haven, FL', 'WTH', 260),
(41, 'Sebring, FL', 'SBG', 270),
(41, 'West Palm Beach, FL', 'WPB', 280),
(41, 'Delray Beach, FL', 'DLB', 290),
(41, 'Deerfield Beach, FL', 'DFB', 300),
(41, 'Fort Lauderdale, FL', 'FTL', 310),
(41, 'Hollywood, FL', 'HOL', 320),
(41, 'Miami, FL', 'MIA', 330);

-- SOUTHWEST CHIEF (corridor_id = 42)
-- Chicago → Los Angeles
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(42, 'Chicago (Union Station), IL', 'CHI', 10),
(42, 'Naperville, IL', 'NPV', 20),
(42, 'Mendota, IL', 'MDT', 30),
(42, 'Princeton, IL', 'PCT', 40),
(42, 'Galesburg, IL', 'GBB', 50),
(42, 'Fort Madison, IA', 'FMD', 60),
(42, 'La Plata, MO', 'LAP', 70),
(42, 'Kansas City (Union Station), MO', 'KCY', 80),
(42, 'Lawrence, KS', 'LRC', 90),
(42, 'Topeka, KS', 'TOP', 100),
(42, 'Newton, KS', 'NEW', 110),
(42, 'Hutchinson, KS', 'HUT', 120),
(42, 'Dodge City, KS', 'DDG', 130),
(42, 'Garden City, KS', 'GCK', 140),
(42, 'Lamar, CO', 'LMR', 150),
(42, 'La Junta, CO', 'LAJ', 160),
(42, 'Trinidad, CO', 'TRI', 170),
(42, 'Raton, NM', 'RAT', 180),
(42, 'Las Vegas, NM', 'LSV', 190),
(42, 'Lamy, NM', 'LMY', 200),
(42, 'Albuquerque, NM', 'ABQ', 210),
(42, 'Gallup, NM', 'GLP', 220),
(42, 'Winslow, AZ', 'WLO', 230),
(42, 'Flagstaff, AZ', 'FLG', 240),
(42, 'Kingman, AZ', 'KNG', 250),
(42, 'Needles, CA', 'NDL', 260),
(42, 'Barstow, CA', 'BAR', 270),
(42, 'Victorville, CA', 'VRV', 280),
(42, 'San Bernardino, CA', 'SNB', 290),
(42, 'Riverside (Downtown), CA', 'RIV', 300),
(42, 'Fullerton, CA', 'FUL', 310),
(42, 'Los Angeles (Union Station), CA', 'LAX', 320);

-- SUNSET LIMITED (corridor_id = 43) - 3x weekly
-- New Orleans → Los Angeles
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(43, 'New Orleans, LA', 'NOL', 10),
(43, 'Schriever, LA', 'SCH', 20),
(43, 'New Iberia, LA', 'NIB', 30),
(43, 'Lafayette, LA', 'LFT', 40),
(43, 'Lake Charles, LA', 'LCH', 50),
(43, 'Beaumont, TX', 'BMT', 60),
(43, 'Houston, TX', 'HOS', 70),
(43, 'San Antonio, TX', 'SAS', 80),
(43, 'Del Rio, TX', 'DRT', 90),
(43, 'Sanderson, TX', 'SND', 100),
(43, 'Alpine, TX', 'ALP', 110),
(43, 'El Paso, TX', 'ELP', 120),
(43, 'Deming, NM', 'DEM', 130),
(43, 'Lordsburg, NM', 'LDB', 140),
(43, 'Benson, AZ', 'BEN', 150),
(43, 'Tucson, AZ', 'TUS', 160),
(43, 'Maricopa, AZ', 'MRC', 170),
(43, 'Yuma, AZ', 'YUM', 180),
(43, 'Palm Springs, CA', 'PSN', 190),
(43, 'Ontario, CA', 'ONA', 200),
(43, 'Pomona, CA', 'POS', 210),
(43, 'Los Angeles (Union Station), CA', 'LAX', 220);

-- TEXAS EAGLE (corridor_id = 44)
-- Chicago → San Antonio, TX
INSERT INTO stops (corridor_id, name, station_code, sort_order) VALUES
(44, 'Chicago (Union Station), IL', 'CHI', 10),
(44, 'Joliet, IL', 'JOL', 20),
(44, 'Pontiac, IL', 'PON', 30),
(44, 'Bloomington-Normal, IL', 'BNL', 40),
(44, 'Lincoln, IL', 'LCN', 50),
(44, 'Springfield, IL', 'SPI', 60),
(44, 'Carlinville, IL', 'CRV', 70),
(44, 'Alton, IL', 'ALN', 80),
(44, 'St. Louis, MO', 'STL', 90),
(44, 'Arcadia, MO', 'ACD', 100),
(44, 'Poplar Bluff, MO', 'PBF', 110),
(44, 'Walnut Ridge, AR', 'WNR', 120),
(44, 'Little Rock, AR', 'LRK', 130),
(44, 'Malvern, AR', 'MVN', 140),
(44, 'Arkadelphia, AR', 'ARK', 150),
(44, 'Hope, AR', 'HOP', 160),
(44, 'Texarkana, AR', 'TXA', 170),
(44, 'Marshall, TX', 'MHL', 180),
(44, 'Longview, TX', 'LVW', 190),
(44, 'Mineola, TX', 'MIN', 200),
(44, 'Dallas, TX', 'DAL', 210),
(44, 'Fort Worth, TX', 'FTW', 220),
(44, 'Cleburne, TX', 'CBR', 230),
(44, 'McGregor, TX', 'MCG', 240),
(44, 'Temple, TX', 'TPL', 250),
(44, 'Taylor, TX', 'TAY', 260),
(44, 'Austin, TX', 'AUS', 270),
(44, 'San Marcos, TX', 'SMC', 280),
(44, 'San Antonio, TX', 'SAS', 290);

More
Allow an image to be marked either as a hero photo or as a map, display both the photo (longscape) and the map image. Put the map image on the left and the routes on the right, as map images are tall and the routes are tall. 
Seed all routes and trains in data.txt, you can probably just execute data.sql for the routes, but data.txt lists all the trains
Append Amtrak against name of corridors that are pure Amtrak 
Create a readme.md formatted for github
Create a changelog.md with all changes tracked

Create a dynamic link to https://asm.transitdocs.com/train/2026/6/13/A/703 where 2026/6/13 is todays date, A is for Amtrak and V is for Via, and 703 is the train number
Also run links.sql 
On the admin page allow the admin to update the username which defaults to admin and the password which defaults to secret by clicking on some sort of settings link, also allow the admin to enter an email address where they will get emailed whenever anyone submits a photo or video, and actually send an email whenever anyone submits a photo or video to the admin with the corridor and train info, a link to approve the submission and a link to the train page and a link to review all submissions
When the admin is adding media, instead of by url, upload file, paste image , offer video, website, image url, upload image, paste image. Next to the train display name e.g. Amtrak 542, Give the admin a search link that will take the train display name and open a new browser tab to search google for videos using that train name Amtrak 542. When the admin is on the trains admin page, the hyperlinked train name should take the admin to the media page, as that is the most common action, the media, edit & deactivate buttons are great 
When viewing routes on the corridor you have a nice visual /corridors/capitol-corridor also add that nice visual when viewing the route of a specific train trains/amtrak-742, hyperlink the station names in both places and create a station page where you can see a link to the corridor and a link to all trains that traverse the corridor and a time at which each train arrives and departs

When the admin is adding media, instead of by  video, website, image url, upload image, paste image, offer video, website, image url, upload image, paste image, paste banner, paste thumbnail. Next to the train display name e.g. Amtrak 542, Give the admin a search link that will take the train display name and open a new browser tab to search google for videos using that train name Amtrak 542, and add a search button to search google images that takes the same train display name. When users are viewing a train, do not display photos that are thumbnails or hero images, they are already displayed.  

Allow users who are Suggest a Photo or Video to add a comment. This comment goes into the existing caption. Ensure captions are displayed to the administrator when managing media, and when emailing the suggestions to the administrator for review. When the admin is reviewing trains, e.g. admin/secret/trains you currently show the following columns (Train	Number	Corridor	Media	Pending	Status	Actions) , add additional columns to indicate if there is no thumbnail and no video and put a red x next to the train. 

Populate all corridors and trains
populate all amtrak corridors and trains Currently Operating Amtrak Train Numbers Exact numbers shift slightly with schedules, but the following ranges are stable as of 2025–2026. Long Distance Trains Sunset Limited: 1, 2 Southwest Chief: 3, 4 California Zephyr: 5, 6 Empire Builder: 7, 8, 27, 28 Coast Starlight: 11, 14 Crescent: 19, 20 Texas Eagle: 21, 22, 421, 422 Floridian: 40, 41 Pennsylvanian: 42, 43 Lake Shore Limited: 48, 49, 448, 449 Cardinal: 50, 51 Auto Train: 52, 53 Vermonter: 54–57 City of New Orleans: 58, 59 Adirondack: 68, 69 Piedmont: 71–78 Carolinian: 79, 80 Palmetto: 89, 90 Silver Meteor: 97, 98 Northeast Regional (long-distance extensions): many in 65–67, 82–88, 93–96, 99, 111, 121–198, etc. Winter Park Express (seasonal): 1105, 1106 NEC / Empire / Corridor Trains Acela: 2100–2290 (varies by trip) Northeast Regional (core NEC): roughly 100–199 block (many 101–198 variants) Keystone Service: 600–674 range Empire Service: 230–288 range Hartford Line / Valley Flyer: 400–499 range (e.g., 400, 405–432, 450–499) Capitol Corridor: 520–553, 720–751 Pacific Surfliner: 562–595, 761–794 Amtrak Cascades: 500–519 Hiawatha: 329–343 Lincoln Service: 300–302, 304–307, 318–319 Illinois Zephyr / Carl Sandburg: 380–383 Illini / Saluki: 390–393 Blue Water: 364, 365 Pere Marquette: 370, 371 Wolverine: 350–355 Missouri River Runner: 311, 316, 318–319 (some overlap by schedule) Heartland Flyer: 821, 822 Downeaster: 680–699, 1689 Borealis: 1333, 1340 Berkshire Flyer (seasonal): 1235, 1244 Ethan Allen Express: 290, 291 Gold Runner: 701–719

Use the following list of all current valid trains and include only amtrak trains above but include also only trains below, only use the intersection of both to ensure you are only listing valid current trains that are only Amtrak trains
Acela Express
816, 817, 880, 2102, 2103, 2104, 2108, 2109, 2110, 2113, 2115, 2121, 2122, 2123, 2124, 2126, 2130, 2150, 2151, 2152, 2153, 2154, 2155, 2159, 2162, 2163, 2166, 2167, 2168, 2169, 2170, 2171, 2172, 2173, 2174, 2190, 2192, 2193, 2201, 2203, 2205, 2206, 2207, 2214, 2215, 2216, 2218, 2220, 2222, 2223, 2224, 2226, 2228, 2233, 2247, 2248, 2249, 2250, 2251, 2252, 2253, 2254, 2255, 2256, 2257, 2258, 2259, 2262, 2263, 2265, 2271, 2274, 2275, 2290, 2292, 2295

Adirondack
68, 69

Amtrak Cascades
500, 502, 503, 504, 505, 506, 507, 508, 509, 511, 516, 517, 518, 519

Auto Train
52, 53

Berkshire Flyer
1233, 1234, 1246

Blue Water/Michigan Service
364, 365, 1364

Borealis
1333, 1340

California Zephyr
5(2), 6(3), 1005, 1006

Capitol Corridor
520, 521, 522, 523, 524, 525, 526, 527, 528, 529, 530, 531, 532, 534, 535, 536, 537, 538, 539, 540, 541, 542, 543, 544, 545, 546, 547, 548, 549, 550, 551, 720, 723, 724, 727, 728, 729, 732, 733, 734, 736, 737, 738, 741, 742, 743, 744, 745, 746, 747, 748, 749, 750, 751

Cardinal
50, 51

Carolinian / Piedmont
71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 105, 1072, 1075, 1123, 1171, 1172

City of New Orleans
58(2), 59, 1058, 1059

Coast Starlight
11(2), 14(2), 1011, 1014

Crescent
19(2), 20(2), 1019, 1020

Downeaster
680, 681, 682, 683, 684, 685, 686, 687, 688, 689, 690, 691, 692, 693, 694, 695, 696, 697, 698, 699, 1689, 1697

Empire Builder
7(2), 8(2), 27, 28, 1007, 1008, 1027, 1028

Empire Service
230, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 243, 244, 245, 280, 281, 283, 284, 1237

Ethan Allen Express
290, 291

Floridian
40(2), 41(2), 1040, 1041

Gold Runner
701, 702, 703, 704, 710, 711, 712, 713, 714, 715, 716, 717, 718, 719, 1701, 1702, 1703, 1704, 1710, 1711, 1712, 1715, 1716, 1717, 1718, 1719

Heartland Flyer
821, 822

Hiawatha
329, 330, 331, 332, 334, 335, 336, 337, 338, 339, 341, 342, 343

Illinois Zephyr/Carl Sandburg
380, 381, 382, 383

Keystone
110, 123, 600, 601, 605, 607, 609, 610, 611, 612, 615, 620, 622, 623, 624, 626, 637, 639, 640, 641, 642, 643, 644, 645, 646, 647, 648, 649, 650, 651, 652, 653, 654, 655, 656, 657, 658, 660, 661, 662, 663, 664, 665, 666, 667, 669, 670, 671, 672, 674

Lake Shore Limited
48, 49, 448, 449

Lincoln Service Missouri River Runner
318, 319

Lincoln Service/Illinois Service
301, 302, 305, 307

Lincoln Service/Michigan Service
300, 306

Maple Leaf
63, 64

Mardi Gras Service
23, 24, 25, 26

Missouri River Runner
311, 316

Northeast Regional
65, 66, 67, 82, 84, 85, 86, 87, 88, 93, 94, 95, 96, 99, 101, 103, 104, 106, 108, 109, 111, 112, 113, 114, 116, 118, 119, 120, 121, 122, 124, 125, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 181, 182, 183, 184, 185, 186, 189, 190, 192, 193, 194, 195, 197, 198, 400, 405, 409, 416, 417, 425, 426, 450, 460, 461, 463, 464, 465, 467, 470, 471, 473, 474, 475, 478, 479, 486, 488, 490, 494, 495, 497, 499, 630, 632, 636, 1108, 1161, 1175, 1194, 2107, 2117

Northest Regional
100, 102, 117, 126, 196, 199, 627, 631, 806, 807, 887, 888, 1195

Pacific Surfliner
562, 564, 566, 567, 572, 573, 577, 579, 580, 581, 582, 584, 586, 587, 588, 591, 593, 595, 757, 761, 765, 769, 770, 774, 777, 779, 782, 784, 785, 786, 790, 791, 794, 1562, 1591, 1595, 1765, 1769, 1770, 1774, 1777, 1784, 1785, 1790

Palmetto
90

Pennsylvanian
42, 43

Pere Marquette/Michigan Service
370, 371

Saluki/Illinois Service
390, 391, 392, 393

Silver Meteor
97(2), 98, 1098

Silver Service / Palmetto
89

Southwest Chief
3(3), 4(2), 1003, 1004

Sunset Limited
1, 2

Texas Eagle
21, 22(2), 1021, 1022

Vermonter
54, 55, 56, 57, 107

Winter Park Express
1105, 1106

Wolverine/Michigan Service
350, 351, 352, 353, 354, 355, 1354
Even more
When a suggestion is submitted, do not allow suggestions already submitted by the administrator to be submitted again, strip unnecessary stuff from the URL to avoid duplicates that lead to the same basic url e.g. https://www.youtube.com/watch?v=WQVUrv_AsVI&test should be https://www.youtube.com/watch?v=WQVUrv_AsVI as anything after the & is irrelevant. On the admin submissions page give the admin an option to approve all pending, and reject all pending, and delete all rejected. On the Admin trains page, give the admin the option to delete all deactivated. . On the Admin trains page, give the admin the option to sort by media, thumb, video, pending, status, train, number, corridor with a default by maybe clicking on the column title and first sorting no media and no thumb and no video and most pending to the top. On the Admin corridors page, give the admin the option to delete all deactivated.

When admins or visitors submit a video, give the following checkboxes: 
Rarities include long consist, doubleheader, sandwich set, reverse set   
Highlights include blow overs, horn  shows, Doppler 

On the Admin corridors page, give the admin the option to sort by trains, region, status, media, thumb, hero, with a default by maybe clicking on the column title and first sorting no media and no thumb and no video. On the Admin corridors page, give the admin the option to delete all deactivated. When the admin is adding media to a corridor, instead of by  video, website, image url, upload image, paste image, offer schedule URL, website, image url, upload image, paste image, paste banner, paste thumbnail.
Allow links to commons.wikimedia.org
Display captains for train videos under the video in a manner similar to photos 
When managing media show the captions for media and allow edits. A small textbox or multi line textbox should suffice with some sort of update button or action. 

Allow administrators to view and update the train videos  checkboxes: 
Rarities include long consist, doubleheader, sandwich set, reverse set   
Highlights include blow overs, horn  shows, Doppler.

Give the administrator the ability to also check a box for best video. 

Sort train videos based on checkboxes, the best video is first and then sort by the number of rarities next. Optimize the database for performance. 

Optimize the homepage for viewers, have a tab for overview, corridors and trains. For the overview create a map with corridor routes with the routes clickable. For the corridors show the corridors with the thumbnail image. For the trains show the trains in table form only. Explain to the viewers why one train video is best. For the overview page show   all amtrak corridor routes on a map, get the data from https://data-usdot.opendata.arcgis.com/datasets/usdot::amtrak-routes and put it on a simple map Leaflet.js preferred or Mapbox GL JS or Google Maps JS API if that doent work.

Move the maps to the last tab (<div class="tabs" style="border-bottom:2px solid var(--border);margin-bottom:0">)  move the trains and maps tabs to the top alongside the corridor link/tab (class="nav-links")

The tabs don’t work on the home page, create a new page for trains and a new page for maps 
Make all the columns clickable on the new trains tab

The maps tab just says Loading route map…
Test cases
Harden my code
Create a Shell script to configure VPS dime with dependencies etc.
Populate the capitol corridor train stops  for capitol corridor from here https://juckins.net/amtrak_timetables/archive/timetables_Capitol_Corridor_California_latest.pdf and here https://content.amtrak.com/content/timetable/Capitol%20Corridor.pdf

Augment the trainstops to indicate if the train stop is for weekday or weekend, add editable fields for the time and checkboxes for the weekeday and weekend to the admin interface for trains, and indicate on the end user interface if the train stop is weekday only or, weekend only 

Is all the data ready to be loaded into the vps server including corridors and trains and corridor routes and train stops except perhaps the train media and corridor media?

For end user submissions of links add a couple of hidden fields the humans never see a and b, if a is filled state the submission was successful but do not save, if b is not equals to “ok” state the submission was successful but do not save, write a couple of javascripts one to set the field b to the letter k and then another to prepend the letter o at the beginning of field b, store the page load time and if the form is submitted in under 2 seconds, reject it. Implement rate limits
 •  Max 5 end user submissions per hour 
•  Max 1 end user submission every minute
do not rate limit the admin on admin pages

make the following configurable from admin settings
 •  Max 5 end user submissions per hour 
•  Max 1 end user submission every minute

Make per day configurable and bump it up to 20 a day

On the map you have a link to the Amtrak schedule e.g. 
 Capitol Corridor
Amtrak schedule →

add a link to the capitol corridor page e.g. http://192.168.8.108:8000/corridors/capitol-corridor and on the capitol corridor page add the Amtrak schedule link in a manner similar to the train page link to transitdocs. E.g. Live Status
Track on TransitDocs →

What can be done to optimize performance of the home page, and which page is fastest the corridor train or maps

Ok go ahead with your recommendation to cache the templates, invalidate the cache on updates, and tell the browser to cache js and css

Submission Rate Limit Controls how many public suggestions all users can submit, not for a single IP address can submit, because single ip addresses can be faked

This doesn’t work for a lot of the corridors as a lot of them have Amtrak prepended to theh corridor name, remove all Amtrak from the beginning of the corridor name so that Amtrak Acela becomes Acela, do that for all corridors that begin with Amtrak and then redo this On the map you have a link to the Amtrak schedule e.g. 
 Capitol Corridor
Amtrak schedule →

add a link to the capitol corridor page e.g. http://192.168.8.108:8000/corridors/capitol-corridor . 

Also on the capitol corridor page add the Amtrak schedule link in a manner similar to the train page link to transitdocs. E.g. on http://192.168.8.108:8080/corridors/amtrak-acela add a link to Amtrak schedule →

When I click on most routes on the map all I get is Texas Eagle
Amtrak schedule → no link to the actual corridor page http://192.168.8.108:8080/corridors/amtrak-texas-eagle 
is should look more like
Texas Eagle
 Texas Eagle →
Amtrak schedule →

also when I drill into the corridor page there is no Amtrak schedule →



Round the Corridor On-time performance: 76.0000% to 76%

Provide an admin option to change the name of the site AmazingTrak and to upload a new thumbnail which is used as the ico of the site 

When either a user or admin submits a youtube url can you retrieve the youtube title description and populate the title and comment based on that?

On the trains page display videos above photos and include rarities tags etc on the videos after the caption

Approve all pending just hangs the entire server.

Allow the administrator to edit the title of all youtube media

Evaluate these requirements and recommend suggestions based on three biggest things to motivate users to visit the site, contribute to the site and share this site on other sites


On the page http://192.168.8.108:8000/admin/secret/trains/264/media the youtube link is not a hyperlink, make it drillable 

If a user submits a train video, and the system retrieves the youtube title, and the title contains the train number in the title, and if the train number is three digits, then auto approve the video. On the suggestions page… Allow the administrator to see what was auto approved and allow the administrator to unapproved the approved suggestions and allow the administrator to view the title of all suggestions and allow the administrator to edit the title and the comments on the suggestion for pending suggestions and take the train hyperlink for suggestions and open up the media for that train, not editing that train name. 

Create a user registration, admin registration flows. Up until now the user referred to as the admin is now referred to by me as the central admin, he doesn’t change, the settings page still allows users to quickly modify the central admins username and password. Provide a users page to view all other users. Allow the central admin to approve user registrations/ admin registrations manually, and allow the central admin to delete users/admins. Once a user has registered with a username, password and email address, if email is configured, then email a confirmation link, if the user confirms the email address then the user is confirmed but not approved. Do not assume the email will work, so allow the user to proceed without email confirmation. If the administrator adds a rarity to a user submitted link they should be approved, if the administrator approves a user submitted video they should be approved. Provide a link from user admin to see the users link submissions. If the user is deleted then delete all the user submissions without rarities. Deleting users should preserve rare submissions. Give the administrator the option to delete all unapproved users. Allow the administrator to delete a user who is confirmed. Allow the administrator to unapproved an approved user. When the administrator adds another admin account, that admin account give that user one of three permission levels, first level is manage submissions, the second level is manage trains, the third level is manage corridors, the fourth level is manage settings. Only the central admin can manage users. Email confirmation should not block approval as email is optional. Auto-approving users based on approved submissions should be explicitly listed as auto-approved status. Central admin username/password editable in settings” is risky. Require re-auth before changing it.

Allow the admin to register regular users as well as admin users, make them approved by default.   



Add open graph/ twitter meta tags on the train page and the corridor page including hero image, train name, best video link and contributor name. 

Show only the best of video embedded in the train page, under that display a list of all youtube videos with the title and the list of rarities with relevant emojis to highlight why certain videos are listed above all others sort best first and then based on most rarities after that then reverse chronological as third sort order. Above the best of video, provide in a single line some anchor links with the relevant emoji’s showing that there are other videos with rarities, anchoring to the table, only show those rarities that exist in the table e.g. long consist (3), doubleheader (1) etc Add in a contributor name to the tables for admins and users.

Don’t auto approve videos if they have rarities. Allow the admin reviewing suggestions to edit rarities. 

Create a user page where any user logged in can see the submissions from the user and when the user registered and sort submission by the number of rarities and give the user badges for submitting first video, submitting first rarity, submitting ten videos, submitting ten rarities, submitting a hundred videos, submitting a hundred rarities

On the http://192.168.8.108:8080/admin/secret/suggestions page remove the tags column and make the actions column double width

Whereever the username is displayed make it clickable e.g. http://192.168.8.108:8080/trains/amtrak-816#video-table and if the user isn’t signed in when they click on the username prompt them to sign in 

Implement the same honeypots and rate limiters for registering a new username as for submissions - For registering a user add a couple of hidden fields the humans never see a and b, if a is filled state the submission was successful but do not save, if b is not equals to “ok” state the submission was successful but do not save, write a couple of javascripts one to set the field b to the letter k and then another to prepend the letter o at the beginning of field b, store the page load time and if the form is submitted in under 2 seconds, reject it.  make the following configurable from admin settings
 •  Max 5 registrations per hour 
Max 20 registrations a day


On every train page:
Share:
[Reddit]
[X]
[Facebook]
[Copy Link]

Optimize train pages open graph/ twitter meta tags  for reddit, x, facebook

Allow admins to mark submissions as spam. When a submission is marked as spam, if the user is just a regular user with no approved links then mark the user as spammer and mark all their links as  rejected -spam.  When a submission is marked as spam, if the user is just a regular user with approved links then mark the user as spammer and mark all their links with no rarities as pending - spam.  if an admin is marking an admins submission as spam, just mark the submission as pending - spam. Update the submissions page to treat pending - spam the same as pending, and rejected -spam the same as rejected. 

Approved submissions should only be unapproved, not from approved to spam

The reject button should have a dropdown clicking reject ust rejects clicking the dropdown allows you to mark as spam, so there should be no need for a separate spam button, 

the train page takes two lines, 
Share: Reddit X Facebook 
Corridor Acela Train Number 816 Live Status Track on TransitDocs →

put that on one line with the share function right aligned and theh other left aligned
Corridor Acela Train Number 816 Live Status Track on TransitDocs → Share: Reddit X Facebook 

Change the message “Please wait a moment before submitting again.” To indicate that if you register once you are approved you can submit unlimited number of links.

Highlight the rarities on the open graph/ twitter meta tags

Create an overview tab for users to see the top five train videos by rarities, top five users by rarities submitted, latest five train videos, and top corridors with the largest number of trains with no videos

On submission of a youtube video with no thumbnail, automatically set it, don’t store the thumbnail locally, just reference it

Add in scenic rarity, environment rarity, historic rarity, special event rarity

Change top contributers on the overview page from Top Contributors by Rarities to Top Contributors by submissions

On the http://192.168.8.108:8080/trains/amtrak-816/suggest provide a search youtube function to search for youtube videos with the train name – remove most of what you added “Search YouTube, Opens YouTube in a new tab — paste the video URL below once you find it.” Just have the search youtube button and move it to the right of the  url


On the http://192.168.8.108:8080/admin/secret/trains/248/media page and hte http://192.168.8.108:8080/admin/secret/trains/248  page provide a quick way to see the public train page, 

Allow Registered users to submit comments, No anonymous comments, No HTML
Comments are pending by default and approved or rejected by admin in manner similar to suggestions
Rate limited, in a manner similar to submissions

http://192.168.8.108:8000/admin/secret/suggestions color the unapprove button as grey

silently drop comments from users marked as spammers but tell the user the comment is pending review, if a user is market as a spammer drop all their existing comments. Add in a max number of comments per user per day of 10 and limit the number of comments to 100 a week across all users. Once a user is approved, auto approve their first three comments. Show users their own pending comments. On the users page show the users comments for each train video. 


Add two more levels to administration. Suggestions, Comments, Trains, Corridors, Settings, Users. 

Add latest comments to the users overview page. 

Only allow unapproved users to be deleted. 

On this page make the unapproved button grey http://192.168.8.108:8000/admin/secret/comments

On this page http://192.168.8.108:8000/trains/amtrak-816 the train breadcrumb is not hyperlinked Corridors › Acela › Amtrak 816 – make Amtrak 816 hyperlinked

If you unapproved a user then unapproved their comments. If you delete a user then delete their comments.  On the http://192.168.8.108:8000/admin/secret/users page give a link to the users page so you can see their comments before deleting them. 

Give the administrator some basic color theme choices, the default and then like seven other good combos including grey, white etc.

Dude the theme should also update the user colors. The user site can be light or dark mode so the themes should invert white becomes black if in dark mode.

