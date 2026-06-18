-- Restoration script: re-inserts user-added media after fresh DB reseed.
-- Run AFTER the server has started and reseeded (all 46 corridors + trains present).
-- RailRat seed links (590 entries) are already inserted by seed.go via links.sql.

BEGIN TRANSACTION;

-- Capitol Corridor level: Wikipedia and Amtrak.com seed links
INSERT INTO media (corridor_id, media_type, source_type, url, title, is_published, added_by)
VALUES (
  (SELECT id FROM corridors WHERE slug='capitol-corridor'),
  'website', 'seed', 'https://en.wikipedia.org/wiki/Capitol_Corridor', 'Capitol Corridor — Wikipedia', 1, 'seed'
);

INSERT INTO media (corridor_id, media_type, source_type, url, title, is_published, added_by)
VALUES (
  (SELECT id FROM corridors WHERE slug='capitol-corridor'),
  'website', 'seed', 'https://www.amtrak.com/capitol-corridor-train', 'Capitol Corridor — Amtrak Official Route Page', 1, 'seed'
);

-- System map (upload)
INSERT INTO media (corridor_id, media_type, source_type, local_path, stored_filename, title, caption, is_published, added_by)
VALUES (
  (SELECT id FROM corridors WHERE slug='capitol-corridor'),
  'image', 'upload',
  'images/FRA-Amtrak-System-Map-FY2026-Q1.png', 'FRA-Amtrak-System-Map-FY2026-Q1.png',
  'Amtrak System Map — FY2026 Q1',
  'Full Amtrak system map showing all routes as of FY2026 Q1. The Capitol Corridor (State Supported) runs between San Jose and Auburn, CA. Source: FRA FY2026 Q1 Performance Report, Figure 2.',
  1, 'admin'
);

-- Corridor hero (paste)
INSERT INTO media (corridor_id, media_type, source_type, local_path, stored_filename, is_published, added_by)
VALUES (
  (SELECT id FROM corridors WHERE slug='capitol-corridor'),
  'image', 'paste', 'images/Capitol-Corridor-1.jpg', 'Capitol-Corridor-1.jpg', 1, 'admin'
);

-- Train 742: CCJPA seed website
INSERT INTO media (train_id, media_type, source_type, url, title, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='742' AND c.slug='capitol-corridor'),
  'website', 'seed', 'https://www.capitolcorridor.org/', 'Capitol Corridor Joint Powers Authority', 1, 'seed'
);

-- Train 742: hero paste image
INSERT INTO media (train_id, media_type, source_type, local_path, stored_filename, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='742' AND c.slug='capitol-corridor'),
  'image', 'paste', 'images/Amtrak-742-1.jpg', 'Amtrak-742-1.jpg', 1, 'admin'
);

-- Train 742: thumbnail paste image
INSERT INTO media (train_id, media_type, source_type, local_path, stored_filename, title, caption, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='742' AND c.slug='capitol-corridor'),
  'image', 'paste', 'images/Amtrak-742-2.jpg', 'Amtrak-742-2.jpg', 'test', 'test', 1, 'admin'
);

-- Train 742: video
INSERT INTO media (train_id, media_type, source_type, url, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='742' AND c.slug='capitol-corridor'),
  'video', 'url', 'https://www.youtube.com/watch?v=dJLlIPW1lH4', 1, 'admin'
);

-- Train 742: set hero and thumbnail
UPDATE trains
SET hero_media_id=(SELECT id FROM media WHERE local_path='images/Amtrak-742-1.jpg')
WHERE train_number='742'
  AND corridor_id=(SELECT id FROM corridors WHERE slug='capitol-corridor');

UPDATE trains
SET thumbnail_media_id=(SELECT id FROM media WHERE local_path='images/Amtrak-742-2.jpg')
WHERE train_number='742'
  AND corridor_id=(SELECT id FROM corridors WHERE slug='capitol-corridor');

-- Train 743: hero paste image
INSERT INTO media (train_id, media_type, source_type, local_path, stored_filename, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='743' AND c.slug='capitol-corridor'),
  'image', 'paste', 'images/Amtrak-743-1.jpg', 'Amtrak-743-1.jpg', 1, 'admin'
);

-- Train 743: video (was an approved public suggestion)
INSERT INTO media (train_id, media_type, source_type, url, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='743' AND c.slug='capitol-corridor'),
  'video', 'url', 'https://www.youtube.com/watch?v=XdIvEU2r9oo', 1, 'approved_suggestion'
);

-- Train 743: set hero
UPDATE trains
SET hero_media_id=(SELECT id FROM media WHERE local_path='images/Amtrak-743-1.jpg')
WHERE train_number='743'
  AND corridor_id=(SELECT id FROM corridors WHERE slug='capitol-corridor');

-- Train 526: hero paste image
INSERT INTO media (train_id, media_type, source_type, local_path, stored_filename, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='526' AND c.slug='capitol-corridor'),
  'image', 'paste', 'images/Amtrak-526-1.jpg', 'Amtrak-526-1.jpg', 1, 'admin'
);

-- Train 526: set hero
UPDATE trains
SET hero_media_id=(SELECT id FROM media WHERE local_path='images/Amtrak-526-1.jpg')
WHERE train_number='526'
  AND corridor_id=(SELECT id FROM corridors WHERE slug='capitol-corridor');

-- Train 529: video
INSERT INTO media (train_id, media_type, source_type, url, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='529' AND c.slug='capitol-corridor'),
  'video', 'url', 'https://www.youtube.com/shorts/Ey8_i9QhjzA', 1, 'admin'
);

-- Train 535: video
INSERT INTO media (train_id, media_type, source_type, url, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='535' AND c.slug='capitol-corridor'),
  'video', 'url', 'https://www.youtube.com/watch?v=iWYwNDk9jPs', 1, 'admin'
);

-- Train 535: thumbnail paste image
INSERT INTO media (train_id, media_type, source_type, local_path, stored_filename, is_published, added_by)
VALUES (
  (SELECT t.id FROM trains t JOIN corridors c ON c.id=t.corridor_id
   WHERE t.train_number='535' AND c.slug='capitol-corridor'),
  'image', 'paste', 'images/Amtrak-535-1.jpg', 'Amtrak-535-1.jpg', 1, 'admin'
);

-- Train 535: set thumbnail
UPDATE trains
SET thumbnail_media_id=(SELECT id FROM media WHERE local_path='images/Amtrak-535-1.jpg')
WHERE train_number='535'
  AND corridor_id=(SELECT id FROM corridors WHERE slug='capitol-corridor');

COMMIT;
