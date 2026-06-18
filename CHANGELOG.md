# Changelog

All notable changes to AmazingTrak are documented here.

===3.9.0=====

### Added
- **Spammer comment cleanup**: marking a user as a spammer now deletes all their existing comments immediately.
- **Spammer silent drop**: comments from spammer accounts are silently discarded (fake success shown) rather than saved.
- **Per-user daily comment limit**: max 10 comments per user per 24 hours.
- **Site-wide weekly comment cap**: max 100 comments across all users per 7 days.
- **Auto-approve first 3 comments** for approved/auto-approved users; their comment appears immediately without admin review. After 3, subsequent comments revert to pending.
- **Pending comment preview**: logged-in users now see their own pending comments on the train page, greyed out with an "awaiting review" label.
- **Comments on user profile**: approved comments appear on `/users/{username}` grouped under a Comments section below Submissions.
- **Latest Comments panel** on the `/overview` page showing the 5 most recent approved comments site-wide.
- **Two new admin permission levels** (levels now 1–6): L2=Comments, L3=Trains, L4=Corridors, L5=Settings, L6=Users. Comments split from Suggestions; Users management accessible without being central admin (add-admin/delete-admin remain central-admin only).
- Train page breadcrumb last segment (train name) is now a hyperlink.
- Unapprove buttons on Suggestions and Comments admin pages are now grey (`btn-secondary`) instead of red.
- Delete button in Admin Users only shown for unapproved users.

===3.8.0=====

### Added
- **Comments on train pages.** Registered, logged-in users can post comments on a train; anonymous visitors see a prompt to log in or register (no anonymous comments).
  - Comments are plain text only — HTML is never rendered (bodies are output-escaped by `html/template`, so any markup shows as literal text).
  - Comments enter as `pending` and are moderated by an admin (approve / reject / unapprove / bulk approve-all / reject-all / delete-all-rejected), mirroring the Suggestions workflow. Only approved comments appear publicly.
  - New **Comments** admin page (permission level 1, same as Suggestions) and a "Pending comments" count + "Review Comments" button on the dashboard.
  - Rate limited per network like submissions: honeypot + minimum fill-time bot defenses and a per-IP hourly/daily cap, configurable under Settings → Comment Rate Limits (`comment_rate_per_hour` default 10, `comment_rate_per_day` default 50). Approved/auto-approved contributors are exempt.
- New `comments` table (`train_id`, `user_id`, `body`, `status`, moderation metadata) with supporting indexes.

===3.7.6=====

### Added
- `/overview` page (linked in nav) with four panels: top 5 videos by rarity count, top 5 contributors by rarities submitted, latest 5 videos, and top corridors with the most trains missing video.
- OG/Twitter meta description on train pages now includes rarity highlights (e.g. "🥪 Sandwich set (2) · 🚂 Doubleheader (1)") when the train has rarity-tagged videos.

===3.7.5=====

### Added
- **Spam marking** on the Suggestions page: admins can flag any submission as spam via a "🚫 Spam" button.
  - Anonymous submission (no linked user): only the one submission is set to *pending-spam*.
  - Linked regular user with no approved submissions: user marked as spammer, all their submissions set to *rejected-spam*.
  - Linked regular user with at least one approved submission: user marked as spammer, all non-rarity submissions set to *pending-spam* (rarity-tagged submissions preserved).
- Spam status is stored as an `is_spam` flag on `suggestions` alongside the existing `status` column (avoids SQLite CHECK-constraint rebuild). The Submissions page treats pending-spam identically to pending, and rejected-spam identically to rejected for filtering purposes.
- Users table gains `is_spammer` flag; spammer accounts display a "🚫 spammer" badge in the Admin Users page.

===3.7.4=====

### Added
- Share buttons (Reddit, X, Facebook, Copy Link) on every train page.
- Train page Open Graph/Twitter tags now include `og:site_name`, `og:locale`, `og:image:alt`/`twitter:image:alt`, `og:image:secure_url`, and always include a description (falls back to a generic line when there's no best video).

===3.7.3=====

### Added
- Contributor usernames are now clickable links to their profile everywhere they appear (public train pages, admin train/corridor media pages); anonymous visitors are prompted to log in first.
- Anti-spam defenses on `/register` matching the existing suggestion form: hidden honeypot fields, a minimum 2-second fill-time check, and per-IP rate limiting (default 5/hour, 20/day), all configurable from admin Settings.

### Changed
- Suggestions admin table: removed the Tags column, widened Actions column.

===3.7.2=====

### Added
- Public user profile page (`/users/{username}`, requires being logged in to view) showing join date, account status, earned badges, and submissions sorted by rarity count.
- Badges for submission milestones: first video, first rarity, 10 videos, 10 rarities, 100 videos, 100 rarities.
- Admins reviewing pending suggestions can now edit a submission's rarity tags (long consist, doubleheader, sandwich set, reverse set) directly from the suggestion review table, in addition to title/comment.

### Changed
- A video submission is no longer auto-approved (even if its YouTube title matches the train number) when it's flagged with a rarity tag — rarity claims always require manual admin review.

===3.7.1=====

### Added
- Central admin can directly create regular user accounts (username/email/password), approved by default, from the Users page — alongside the existing "add admin" flow.
- Open Graph and Twitter Card meta tags on train and corridor pages: hero image, name, best video link, and contributor name.
- Train page now embeds only the single best video (highest rarity count, then newest), with a one-line anchor summary of rarity emoji counts (e.g. 🚃 Long consist (3)) linking to a full video table below listing every video's title, rarity emoji badges, and contributor name.
- Media items now track and display a "Contributor" (registered username, or a generic label for admin/anonymous submissions) in the admin train/corridor media views and the public train video table.

### Changed
- Video sort order on the train page: best-flagged first, then by rarity-tag count, then reverse chronological.

===3.7.0=====

### Added
- User accounts: visitors can self-register (username, password, optional email) and log in at `/register` and `/login`. Logged-in users have their photo/video submissions attributed to their account; anonymous submission still works unchanged.
- Optional email confirmation: when SMTP is configured a confirmation link is emailed, but confirmation is never required to use the site or be approved. Confirmed accounts show a distinct "confirmed" status.
- Central admin Users page: approve / unapprove / delete registered users, "delete all unapproved users", view a user's submissions, and create additional admin accounts.
- Sub-admin accounts with cumulative permission levels (L1 Submissions → L2 Trains → L3 Corridors → L4 Settings). Only the central admin can manage users. Admin nav links and routes are gated by level.
- Auto-approval of users: adding a rarity to a user's media, or approving a user's video submission, promotes that user to a distinct "auto_approved" status.
- User deletion preserves rare submissions: non-rare submissions (and their uploaded files) are removed; submissions carrying a rarity tag are kept with their user link cleared.

### Changed
- The original single admin is now the immutable "central admin" (never deleted, full access). Changing its username/password in Settings now requires re-entering the current password.

### Security
- Re-authentication required before admin credential changes; sub-admin routes enforce least-privilege by permission level.



### Added
- Route stop data for all 44 Amtrak corridors (Northeast Corridor, State Supported, Long Distance) sourced from FRA FY2026 Q1 Performance Report
- Amtrak schedule URL field per corridor (admin-editable); shown on the corridor page in the same style as the TransitDocs link on train pages
- Map popups now include a link to the corridor page when the route name matches a known corridor
- Admin settings: submission rate limits (per minute, per hour, per day) are now configurable without redeployment
- Per-day submission limit raised from 10 to 20 (default)
- Compiled template cache — all public and admin templates compiled once at startup instead of per request
- Homepage HTML cache — rendered index page stored in memory and served directly on repeat loads; invalidated automatically on any corridor write (create, update, toggle, delete, hero/thumbnail change)
- Static asset cache-busting — `style.css` and `theme.js` URLs include a content-hash query string; browsers cache them for one year and re-fetch only when the file changes

### Changed
- Public submission rate limit raised from 3/hour to 5/hour
- Added per-minute rate limit (1/minute default) for public submissions
- Timing check for public submissions tightened from 4 seconds to 2 seconds
- Admin version bumped to v3.6.5

### Security
- Added two additional bot traps on the public suggestion form: hidden field `a` (honeypot) and hidden field `b` (must equal `"ok"`, set by two separate JS scripts); submissions that fail either check receive a fake success response without being saved
- Client-side JS blocks form submission if less than 2 seconds have elapsed since page load

### Fixed
- Paste upload: files were being saved to disk with original extension (`.png`) while the database recorded `.jpg`, causing a path mismatch and broken images
- Corridor media handler was calling `r.FormFile` twice on the same field, which could produce an empty file header on the second read

---

## [0.1.0] — 2026-06-13

Initial working release.

### Added
- Public pages: corridor list, corridor detail, train detail with hero image, photos, videos, links, and schedule
- Admin panel at secret URL: corridor management, train management, media management, suggestion review
- Image uploads re-encoded to JPEG (strips EXIF metadata); GPS coordinates extracted from JPEG EXIF before stripping
- Paste-to-upload in admin — clipboard images submitted via DataTransfer
- Public media suggestion form with rate limiting (3/hr, 10/day per IP), server-side timing check, honeypot field, and duplicate detection
- URL allowlist for public submissions: YouTube, Vimeo, Flickr, Imgur, railpictures.net, rrpicturearchives.net, Instagram
- Hero image, thumbnail, and route map image per train
- Leaflet map showing photo location when EXIF GPS is present
- Live status link to TransitDocs per train
- Dark/light/auto theme toggle (persisted via cookie)
- Admin audit log
- Session-based auth with bcrypt passwords, CSRF tokens, login throttling (5 failures/15 min), 24h session expiry
- `/healthz` endpoint
- `X-Content-Type-Options`, `Referrer-Policy`, `Permissions-Policy`, and CSP security headers
