-- 590 RailRat external links across 46 corridors
-- Source: https://railrat.net/trains/
-- Seed these after trains table is populated.

INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/816/', 'Amtrak 816 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '816' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/817/', 'Amtrak 817 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '817' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/880/', 'Amtrak 880 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '880' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2102/', 'Amtrak 2102 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2102' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2103/', 'Amtrak 2103 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2103' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2104/', 'Amtrak 2104 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2104' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2108/', 'Amtrak 2108 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2108' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2109/', 'Amtrak 2109 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2109' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2110/', 'Amtrak 2110 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2110' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2113/', 'Amtrak 2113 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2113' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2115/', 'Amtrak 2115 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2115' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2121/', 'Amtrak 2121 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2121' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2122/', 'Amtrak 2122 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2122' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2123/', 'Amtrak 2123 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2123' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2124/', 'Amtrak 2124 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2124' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2126/', 'Amtrak 2126 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2126' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2130/', 'Amtrak 2130 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2130' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2150/', 'Amtrak 2150 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2150' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2151/', 'Amtrak 2151 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2151' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2152/', 'Amtrak 2152 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2152' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2153/', 'Amtrak 2153 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2153' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2154/', 'Amtrak 2154 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2154' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2155/', 'Amtrak 2155 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2155' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2159/', 'Amtrak 2159 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2159' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2162/', 'Amtrak 2162 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2162' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2163/', 'Amtrak 2163 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2163' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2166/', 'Amtrak 2166 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2166' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2167/', 'Amtrak 2167 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2167' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2168/', 'Amtrak 2168 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2168' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2169/', 'Amtrak 2169 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2169' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2170/', 'Amtrak 2170 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2170' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2171/', 'Amtrak 2171 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2171' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2172/', 'Amtrak 2172 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2172' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2173/', 'Amtrak 2173 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2173' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2174/', 'Amtrak 2174 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2174' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2190/', 'Amtrak 2190 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2190' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2192/', 'Amtrak 2192 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2192' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2193/', 'Amtrak 2193 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2193' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2201/', 'Amtrak 2201 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2201' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2203/', 'Amtrak 2203 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2203' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2205/', 'Amtrak 2205 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2205' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2206/', 'Amtrak 2206 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2206' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2207/', 'Amtrak 2207 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2207' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2214/', 'Amtrak 2214 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2214' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2215/', 'Amtrak 2215 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2215' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2216/', 'Amtrak 2216 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2216' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2218/', 'Amtrak 2218 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2218' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2220/', 'Amtrak 2220 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2220' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2222/', 'Amtrak 2222 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2222' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2223/', 'Amtrak 2223 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2223' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2224/', 'Amtrak 2224 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2224' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2226/', 'Amtrak 2226 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2226' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2228/', 'Amtrak 2228 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2228' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2233/', 'Amtrak 2233 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2233' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2247/', 'Amtrak 2247 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2247' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2248/', 'Amtrak 2248 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2248' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2249/', 'Amtrak 2249 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2249' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2250/', 'Amtrak 2250 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2250' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2251/', 'Amtrak 2251 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2251' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2252/', 'Amtrak 2252 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2252' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2253/', 'Amtrak 2253 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2253' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2254/', 'Amtrak 2254 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2254' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2255/', 'Amtrak 2255 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2255' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2256/', 'Amtrak 2256 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2256' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2257/', 'Amtrak 2257 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2257' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2258/', 'Amtrak 2258 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2258' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2259/', 'Amtrak 2259 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2259' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2262/', 'Amtrak 2262 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2262' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2263/', 'Amtrak 2263 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2263' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2265/', 'Amtrak 2265 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2265' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2271/', 'Amtrak 2271 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2271' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2274/', 'Amtrak 2274 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2274' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2275/', 'Amtrak 2275 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2275' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2290/', 'Amtrak 2290 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2290' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2292/', 'Amtrak 2292 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2292' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2295/', 'Amtrak 2295 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2295' AND c.name = 'Amtrak Acela';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/68/', 'Amtrak 68 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '68' AND c.name = 'Adirondack';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/69/', 'Amtrak 69 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '69' AND c.name = 'Adirondack';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/500/', 'Amtrak 500 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '500' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/502/', 'Amtrak 502 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '502' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/503/', 'Amtrak 503 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '503' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/504/', 'Amtrak 504 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '504' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/505/', 'Amtrak 505 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '505' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/506/', 'Amtrak 506 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '506' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/507/', 'Amtrak 507 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '507' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/508/', 'Amtrak 508 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '508' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/509/', 'Amtrak 509 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '509' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/511/', 'Amtrak 511 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '511' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/516/', 'Amtrak 516 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '516' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/517/', 'Amtrak 517 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '517' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/518/', 'Amtrak 518 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '518' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/519/', 'Amtrak 519 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '519' AND c.name = 'Amtrak Cascades';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/52/', 'Amtrak 52 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '52' AND c.name = 'Amtrak Auto Train';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/53/', 'Amtrak 53 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '53' AND c.name = 'Amtrak Auto Train';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1233/', 'Amtrak 1233 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1233' AND c.name = 'Berkshire Flyer';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1234/', 'Amtrak 1234 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1234' AND c.name = 'Berkshire Flyer';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1246/', 'Amtrak 1246 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1246' AND c.name = 'Berkshire Flyer';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/364/', 'Amtrak 364 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '364' AND c.name = 'Blue Water';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/365/', 'Amtrak 365 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '365' AND c.name = 'Blue Water';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1364/', 'Amtrak 1364 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1364' AND c.name = 'Blue Water';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1333/', 'Amtrak 1333 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1333' AND c.name = 'Borealis';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1340/', 'Amtrak 1340 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1340' AND c.name = 'Borealis';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/5/', 'Amtrak 5 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '5' AND c.name = 'Amtrak California Zephyr';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/6/', 'Amtrak 6 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '6' AND c.name = 'Amtrak California Zephyr';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1005/', 'Amtrak 1005 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1005' AND c.name = 'Amtrak California Zephyr';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1006/', 'Amtrak 1006 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1006' AND c.name = 'Amtrak California Zephyr';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/520/', 'Amtrak 520 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '520' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/521/', 'Amtrak 521 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '521' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/522/', 'Amtrak 522 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '522' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/523/', 'Amtrak 523 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '523' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/524/', 'Amtrak 524 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '524' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/525/', 'Amtrak 525 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '525' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/526/', 'Amtrak 526 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '526' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/527/', 'Amtrak 527 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '527' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/528/', 'Amtrak 528 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '528' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/529/', 'Amtrak 529 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '529' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/530/', 'Amtrak 530 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '530' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/531/', 'Amtrak 531 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '531' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/532/', 'Amtrak 532 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '532' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/534/', 'Amtrak 534 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '534' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/535/', 'Amtrak 535 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '535' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/536/', 'Amtrak 536 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '536' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/537/', 'Amtrak 537 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '537' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/538/', 'Amtrak 538 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '538' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/539/', 'Amtrak 539 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '539' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/540/', 'Amtrak 540 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '540' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/541/', 'Amtrak 541 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '541' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/542/', 'Amtrak 542 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '542' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/543/', 'Amtrak 543 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '543' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/544/', 'Amtrak 544 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '544' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/545/', 'Amtrak 545 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '545' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/546/', 'Amtrak 546 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '546' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/547/', 'Amtrak 547 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '547' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/548/', 'Amtrak 548 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '548' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/549/', 'Amtrak 549 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '549' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/550/', 'Amtrak 550 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '550' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/551/', 'Amtrak 551 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '551' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/720/', 'Amtrak 720 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '720' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/723/', 'Amtrak 723 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '723' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/724/', 'Amtrak 724 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '724' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/727/', 'Amtrak 727 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '727' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/728/', 'Amtrak 728 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '728' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/729/', 'Amtrak 729 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '729' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/732/', 'Amtrak 732 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '732' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/733/', 'Amtrak 733 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '733' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/734/', 'Amtrak 734 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '734' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/736/', 'Amtrak 736 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '736' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/737/', 'Amtrak 737 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '737' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/738/', 'Amtrak 738 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '738' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/741/', 'Amtrak 741 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '741' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/742/', 'Amtrak 742 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '742' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/743/', 'Amtrak 743 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '743' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/744/', 'Amtrak 744 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '744' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/745/', 'Amtrak 745 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '745' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/746/', 'Amtrak 746 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '746' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/747/', 'Amtrak 747 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '747' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/748/', 'Amtrak 748 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '748' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/749/', 'Amtrak 749 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '749' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/750/', 'Amtrak 750 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '750' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/751/', 'Amtrak 751 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '751' AND c.name = 'Capitol Corridor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/50/', 'Amtrak 50 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '50' AND c.name = 'Amtrak Cardinal';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/51/', 'Amtrak 51 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '51' AND c.name = 'Amtrak Cardinal';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/71/', 'Amtrak 71 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '71' AND c.name = 'Piedmont';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/72/', 'Amtrak 72 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '72' AND c.name = 'Piedmont';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/73/', 'Amtrak 73 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '73' AND c.name = 'Piedmont';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/74/', 'Amtrak 74 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '74' AND c.name = 'Piedmont';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/75/', 'Amtrak 75 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '75' AND c.name = 'Piedmont';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/76/', 'Amtrak 76 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '76' AND c.name = 'Piedmont';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/77/', 'Amtrak 77 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '77' AND c.name = 'Piedmont';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/78/', 'Amtrak 78 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '78' AND c.name = 'Piedmont';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/79/', 'Amtrak 79 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '79' AND c.name = 'Carolinian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/80/', 'Amtrak 80 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '80' AND c.name = 'Carolinian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/105/', 'Amtrak 105 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '105' AND c.name = 'Carolinian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1072/', 'Amtrak 1072 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1072' AND c.name = 'Carolinian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1075/', 'Amtrak 1075 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1075' AND c.name = 'Carolinian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1123/', 'Amtrak 1123 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1123' AND c.name = 'Carolinian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1171/', 'Amtrak 1171 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1171' AND c.name = 'Carolinian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1172/', 'Amtrak 1172 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1172' AND c.name = 'Carolinian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/58/', 'Amtrak 58 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '58' AND c.name = 'Amtrak City of New Orleans';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/59/', 'Amtrak 59 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '59' AND c.name = 'Amtrak City of New Orleans';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1058/', 'Amtrak 1058 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1058' AND c.name = 'Amtrak City of New Orleans';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1059/', 'Amtrak 1059 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1059' AND c.name = 'Amtrak City of New Orleans';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/11/', 'Amtrak 11 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '11' AND c.name = 'Amtrak Coast Starlight';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/14/', 'Amtrak 14 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '14' AND c.name = 'Amtrak Coast Starlight';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1011/', 'Amtrak 1011 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1011' AND c.name = 'Amtrak Coast Starlight';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1014/', 'Amtrak 1014 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1014' AND c.name = 'Amtrak Coast Starlight';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/19/', 'Amtrak 19 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '19' AND c.name = 'Amtrak Crescent';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/20/', 'Amtrak 20 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '20' AND c.name = 'Amtrak Crescent';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1019/', 'Amtrak 1019 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1019' AND c.name = 'Amtrak Crescent';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1020/', 'Amtrak 1020 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1020' AND c.name = 'Amtrak Crescent';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/680/', 'Amtrak 680 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '680' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/681/', 'Amtrak 681 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '681' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/682/', 'Amtrak 682 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '682' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/683/', 'Amtrak 683 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '683' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/684/', 'Amtrak 684 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '684' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/685/', 'Amtrak 685 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '685' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/686/', 'Amtrak 686 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '686' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/687/', 'Amtrak 687 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '687' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/688/', 'Amtrak 688 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '688' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/689/', 'Amtrak 689 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '689' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/690/', 'Amtrak 690 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '690' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/691/', 'Amtrak 691 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '691' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/692/', 'Amtrak 692 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '692' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/693/', 'Amtrak 693 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '693' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/694/', 'Amtrak 694 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '694' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/695/', 'Amtrak 695 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '695' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/696/', 'Amtrak 696 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '696' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/697/', 'Amtrak 697 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '697' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/698/', 'Amtrak 698 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '698' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/699/', 'Amtrak 699 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '699' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1689/', 'Amtrak 1689 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1689' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1697/', 'Amtrak 1697 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1697' AND c.name = 'Downeaster';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/7/', 'Amtrak 7 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '7' AND c.name = 'Amtrak Empire Builder';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/8/', 'Amtrak 8 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '8' AND c.name = 'Amtrak Empire Builder';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/27/', 'Amtrak 27 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '27' AND c.name = 'Amtrak Empire Builder';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/28/', 'Amtrak 28 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '28' AND c.name = 'Amtrak Empire Builder';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1007/', 'Amtrak 1007 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1007' AND c.name = 'Amtrak Empire Builder';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1008/', 'Amtrak 1008 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1008' AND c.name = 'Amtrak Empire Builder';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1027/', 'Amtrak 1027 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1027' AND c.name = 'Amtrak Empire Builder';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1028/', 'Amtrak 1028 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1028' AND c.name = 'Amtrak Empire Builder';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/230/', 'Amtrak 230 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '230' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/232/', 'Amtrak 232 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '232' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/233/', 'Amtrak 233 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '233' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/234/', 'Amtrak 234 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '234' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/235/', 'Amtrak 235 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '235' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/236/', 'Amtrak 236 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '236' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/237/', 'Amtrak 237 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '237' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/238/', 'Amtrak 238 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '238' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/239/', 'Amtrak 239 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '239' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/240/', 'Amtrak 240 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '240' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/241/', 'Amtrak 241 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '241' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/243/', 'Amtrak 243 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '243' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/244/', 'Amtrak 244 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '244' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/245/', 'Amtrak 245 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '245' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/280/', 'Amtrak 280 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '280' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/281/', 'Amtrak 281 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '281' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/283/', 'Amtrak 283 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '283' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/284/', 'Amtrak 284 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '284' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1237/', 'Amtrak 1237 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1237' AND c.name = 'Empire Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/290/', 'Amtrak 290 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '290' AND c.name = 'Ethan Allen Express';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/291/', 'Amtrak 291 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '291' AND c.name = 'Ethan Allen Express';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/40/', 'Amtrak 40 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '40' AND c.name = 'Amtrak Floridian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/41/', 'Amtrak 41 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '41' AND c.name = 'Amtrak Floridian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1040/', 'Amtrak 1040 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1040' AND c.name = 'Amtrak Floridian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1041/', 'Amtrak 1041 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1041' AND c.name = 'Amtrak Floridian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/701/', 'Amtrak 701 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '701' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/702/', 'Amtrak 702 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '702' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/703/', 'Amtrak 703 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '703' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/704/', 'Amtrak 704 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '704' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/710/', 'Amtrak 710 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '710' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/711/', 'Amtrak 711 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '711' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/712/', 'Amtrak 712 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '712' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/713/', 'Amtrak 713 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '713' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/714/', 'Amtrak 714 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '714' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/715/', 'Amtrak 715 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '715' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/716/', 'Amtrak 716 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '716' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/717/', 'Amtrak 717 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '717' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/718/', 'Amtrak 718 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '718' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/719/', 'Amtrak 719 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '719' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1701/', 'Amtrak 1701 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1701' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1702/', 'Amtrak 1702 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1702' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1703/', 'Amtrak 1703 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1703' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1704/', 'Amtrak 1704 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1704' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1710/', 'Amtrak 1710 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1710' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1711/', 'Amtrak 1711 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1711' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1712/', 'Amtrak 1712 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1712' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1715/', 'Amtrak 1715 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1715' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1716/', 'Amtrak 1716 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1716' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1717/', 'Amtrak 1717 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1717' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1718/', 'Amtrak 1718 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1718' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1719/', 'Amtrak 1719 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1719' AND c.name = 'Gold Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/821/', 'Amtrak 821 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '821' AND c.name = 'Heartland Flyer';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/822/', 'Amtrak 822 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '822' AND c.name = 'Heartland Flyer';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/329/', 'Amtrak 329 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '329' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/330/', 'Amtrak 330 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '330' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/331/', 'Amtrak 331 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '331' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/332/', 'Amtrak 332 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '332' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/334/', 'Amtrak 334 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '334' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/335/', 'Amtrak 335 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '335' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/336/', 'Amtrak 336 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '336' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/337/', 'Amtrak 337 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '337' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/338/', 'Amtrak 338 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '338' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/339/', 'Amtrak 339 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '339' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/341/', 'Amtrak 341 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '341' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/342/', 'Amtrak 342 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '342' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/343/', 'Amtrak 343 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '343' AND c.name = 'Hiawatha';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/380/', 'Amtrak 380 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '380' AND c.name = 'Carl Sandburg / Illinois Zephyr';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/381/', 'Amtrak 381 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '381' AND c.name = 'Carl Sandburg / Illinois Zephyr';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/382/', 'Amtrak 382 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '382' AND c.name = 'Carl Sandburg / Illinois Zephyr';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/383/', 'Amtrak 383 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '383' AND c.name = 'Carl Sandburg / Illinois Zephyr';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/110/', 'Amtrak 110 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '110' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/123/', 'Amtrak 123 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '123' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/600/', 'Amtrak 600 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '600' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/601/', 'Amtrak 601 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '601' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/605/', 'Amtrak 605 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '605' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/607/', 'Amtrak 607 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '607' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/609/', 'Amtrak 609 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '609' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/610/', 'Amtrak 610 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '610' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/611/', 'Amtrak 611 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '611' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/612/', 'Amtrak 612 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '612' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/615/', 'Amtrak 615 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '615' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/620/', 'Amtrak 620 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '620' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/622/', 'Amtrak 622 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '622' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/623/', 'Amtrak 623 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '623' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/624/', 'Amtrak 624 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '624' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/626/', 'Amtrak 626 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '626' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/637/', 'Amtrak 637 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '637' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/639/', 'Amtrak 639 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '639' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/640/', 'Amtrak 640 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '640' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/641/', 'Amtrak 641 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '641' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/642/', 'Amtrak 642 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '642' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/643/', 'Amtrak 643 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '643' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/644/', 'Amtrak 644 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '644' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/645/', 'Amtrak 645 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '645' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/646/', 'Amtrak 646 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '646' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/647/', 'Amtrak 647 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '647' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/648/', 'Amtrak 648 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '648' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/649/', 'Amtrak 649 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '649' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/650/', 'Amtrak 650 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '650' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/651/', 'Amtrak 651 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '651' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/652/', 'Amtrak 652 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '652' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/653/', 'Amtrak 653 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '653' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/654/', 'Amtrak 654 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '654' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/655/', 'Amtrak 655 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '655' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/656/', 'Amtrak 656 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '656' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/657/', 'Amtrak 657 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '657' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/658/', 'Amtrak 658 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '658' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/660/', 'Amtrak 660 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '660' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/661/', 'Amtrak 661 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '661' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/662/', 'Amtrak 662 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '662' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/663/', 'Amtrak 663 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '663' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/664/', 'Amtrak 664 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '664' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/665/', 'Amtrak 665 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '665' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/666/', 'Amtrak 666 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '666' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/667/', 'Amtrak 667 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '667' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/669/', 'Amtrak 669 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '669' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/670/', 'Amtrak 670 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '670' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/671/', 'Amtrak 671 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '671' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/672/', 'Amtrak 672 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '672' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/674/', 'Amtrak 674 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '674' AND c.name = 'Keystone';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/48/', 'Amtrak 48 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '48' AND c.name = 'Amtrak Lake Shore Limited';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/49/', 'Amtrak 49 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '49' AND c.name = 'Amtrak Lake Shore Limited';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/448/', 'Amtrak 448 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '448' AND c.name = 'Amtrak Lake Shore Limited';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/449/', 'Amtrak 449 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '449' AND c.name = 'Amtrak Lake Shore Limited';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/318/', 'Amtrak 318 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '318' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/319/', 'Amtrak 319 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '319' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/301/', 'Amtrak 301 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '301' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/302/', 'Amtrak 302 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '302' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/305/', 'Amtrak 305 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '305' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/307/', 'Amtrak 307 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '307' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/300/', 'Amtrak 300 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '300' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/306/', 'Amtrak 306 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '306' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/63/', 'Amtrak 63 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '63' AND c.name = 'Maple Leaf';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/64/', 'Amtrak 64 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '64' AND c.name = 'Maple Leaf';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/23/', 'Amtrak 23 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '23' AND c.name = 'Mardi Gras Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/24/', 'Amtrak 24 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '24' AND c.name = 'Mardi Gras Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/25/', 'Amtrak 25 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '25' AND c.name = 'Mardi Gras Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/26/', 'Amtrak 26 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '26' AND c.name = 'Mardi Gras Service';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/311/', 'Amtrak 311 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '311' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/316/', 'Amtrak 316 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '316' AND c.name = 'Lincoln / Missouri River Runner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/65/', 'Amtrak 65 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '65' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/66/', 'Amtrak 66 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '66' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/67/', 'Amtrak 67 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '67' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/82/', 'Amtrak 82 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '82' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/84/', 'Amtrak 84 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '84' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/85/', 'Amtrak 85 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '85' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/86/', 'Amtrak 86 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '86' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/87/', 'Amtrak 87 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '87' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/88/', 'Amtrak 88 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '88' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/93/', 'Amtrak 93 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '93' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/94/', 'Amtrak 94 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '94' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/95/', 'Amtrak 95 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '95' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/96/', 'Amtrak 96 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '96' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/99/', 'Amtrak 99 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '99' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/101/', 'Amtrak 101 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '101' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/103/', 'Amtrak 103 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '103' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/104/', 'Amtrak 104 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '104' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/106/', 'Amtrak 106 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '106' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/108/', 'Amtrak 108 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '108' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/109/', 'Amtrak 109 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '109' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/111/', 'Amtrak 111 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '111' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/112/', 'Amtrak 112 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '112' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/113/', 'Amtrak 113 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '113' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/114/', 'Amtrak 114 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '114' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/116/', 'Amtrak 116 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '116' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/118/', 'Amtrak 118 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '118' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/119/', 'Amtrak 119 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '119' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/120/', 'Amtrak 120 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '120' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/121/', 'Amtrak 121 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '121' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/122/', 'Amtrak 122 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '122' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/124/', 'Amtrak 124 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '124' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/125/', 'Amtrak 125 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '125' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/127/', 'Amtrak 127 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '127' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/128/', 'Amtrak 128 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '128' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/129/', 'Amtrak 129 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '129' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/130/', 'Amtrak 130 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '130' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/131/', 'Amtrak 131 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '131' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/132/', 'Amtrak 132 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '132' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/133/', 'Amtrak 133 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '133' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/134/', 'Amtrak 134 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '134' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/135/', 'Amtrak 135 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '135' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/136/', 'Amtrak 136 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '136' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/137/', 'Amtrak 137 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '137' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/138/', 'Amtrak 138 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '138' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/139/', 'Amtrak 139 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '139' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/140/', 'Amtrak 140 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '140' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/141/', 'Amtrak 141 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '141' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/142/', 'Amtrak 142 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '142' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/143/', 'Amtrak 143 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '143' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/144/', 'Amtrak 144 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '144' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/145/', 'Amtrak 145 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '145' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/146/', 'Amtrak 146 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '146' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/147/', 'Amtrak 147 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '147' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/148/', 'Amtrak 148 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '148' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/149/', 'Amtrak 149 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '149' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/150/', 'Amtrak 150 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '150' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/151/', 'Amtrak 151 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '151' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/152/', 'Amtrak 152 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '152' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/153/', 'Amtrak 153 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '153' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/154/', 'Amtrak 154 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '154' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/155/', 'Amtrak 155 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '155' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/156/', 'Amtrak 156 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '156' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/157/', 'Amtrak 157 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '157' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/158/', 'Amtrak 158 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '158' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/159/', 'Amtrak 159 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '159' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/160/', 'Amtrak 160 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '160' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/161/', 'Amtrak 161 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '161' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/162/', 'Amtrak 162 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '162' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/163/', 'Amtrak 163 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '163' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/164/', 'Amtrak 164 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '164' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/165/', 'Amtrak 165 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '165' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/166/', 'Amtrak 166 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '166' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/167/', 'Amtrak 167 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '167' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/168/', 'Amtrak 168 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '168' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/169/', 'Amtrak 169 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '169' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/170/', 'Amtrak 170 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '170' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/171/', 'Amtrak 171 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '171' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/172/', 'Amtrak 172 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '172' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/173/', 'Amtrak 173 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '173' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/174/', 'Amtrak 174 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '174' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/175/', 'Amtrak 175 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '175' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/176/', 'Amtrak 176 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '176' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/177/', 'Amtrak 177 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '177' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/178/', 'Amtrak 178 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '178' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/179/', 'Amtrak 179 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '179' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/181/', 'Amtrak 181 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '181' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/182/', 'Amtrak 182 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '182' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/183/', 'Amtrak 183 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '183' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/184/', 'Amtrak 184 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '184' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/185/', 'Amtrak 185 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '185' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/186/', 'Amtrak 186 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '186' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/189/', 'Amtrak 189 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '189' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/190/', 'Amtrak 190 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '190' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/192/', 'Amtrak 192 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '192' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/193/', 'Amtrak 193 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '193' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/194/', 'Amtrak 194 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '194' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/195/', 'Amtrak 195 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '195' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/197/', 'Amtrak 197 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '197' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/198/', 'Amtrak 198 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '198' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/400/', 'Amtrak 400 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '400' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/405/', 'Amtrak 405 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '405' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/409/', 'Amtrak 409 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '409' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/416/', 'Amtrak 416 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '416' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/417/', 'Amtrak 417 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '417' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/425/', 'Amtrak 425 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '425' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/426/', 'Amtrak 426 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '426' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/450/', 'Amtrak 450 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '450' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/460/', 'Amtrak 460 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '460' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/461/', 'Amtrak 461 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '461' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/463/', 'Amtrak 463 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '463' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/464/', 'Amtrak 464 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '464' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/465/', 'Amtrak 465 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '465' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/467/', 'Amtrak 467 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '467' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/470/', 'Amtrak 470 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '470' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/471/', 'Amtrak 471 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '471' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/473/', 'Amtrak 473 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '473' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/474/', 'Amtrak 474 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '474' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/475/', 'Amtrak 475 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '475' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/478/', 'Amtrak 478 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '478' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/479/', 'Amtrak 479 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '479' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/486/', 'Amtrak 486 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '486' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/488/', 'Amtrak 488 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '488' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/490/', 'Amtrak 490 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '490' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/494/', 'Amtrak 494 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '494' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/495/', 'Amtrak 495 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '495' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/497/', 'Amtrak 497 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '497' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/499/', 'Amtrak 499 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '499' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/630/', 'Amtrak 630 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '630' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/632/', 'Amtrak 632 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '632' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/636/', 'Amtrak 636 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '636' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1108/', 'Amtrak 1108 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1108' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1161/', 'Amtrak 1161 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1161' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1175/', 'Amtrak 1175 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1175' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1194/', 'Amtrak 1194 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1194' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2107/', 'Amtrak 2107 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2107' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2117/', 'Amtrak 2117 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2117' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/100/', 'Amtrak 100 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '100' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/102/', 'Amtrak 102 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '102' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/117/', 'Amtrak 117 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '117' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/126/', 'Amtrak 126 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '126' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/196/', 'Amtrak 196 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '196' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/199/', 'Amtrak 199 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '199' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/627/', 'Amtrak 627 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '627' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/631/', 'Amtrak 631 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '631' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/806/', 'Amtrak 806 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '806' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/807/', 'Amtrak 807 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '807' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/887/', 'Amtrak 887 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '887' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/888/', 'Amtrak 888 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '888' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1195/', 'Amtrak 1195 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1195' AND c.name = 'Amtrak Northeast Regional';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/562/', 'Amtrak 562 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '562' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/564/', 'Amtrak 564 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '564' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/566/', 'Amtrak 566 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '566' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/567/', 'Amtrak 567 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '567' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/572/', 'Amtrak 572 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '572' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/573/', 'Amtrak 573 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '573' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/577/', 'Amtrak 577 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '577' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/579/', 'Amtrak 579 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '579' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/580/', 'Amtrak 580 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '580' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/581/', 'Amtrak 581 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '581' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/582/', 'Amtrak 582 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '582' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/584/', 'Amtrak 584 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '584' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/586/', 'Amtrak 586 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '586' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/587/', 'Amtrak 587 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '587' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/588/', 'Amtrak 588 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '588' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/591/', 'Amtrak 591 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '591' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/593/', 'Amtrak 593 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '593' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/595/', 'Amtrak 595 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '595' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/757/', 'Amtrak 757 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '757' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/761/', 'Amtrak 761 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '761' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/765/', 'Amtrak 765 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '765' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/769/', 'Amtrak 769 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '769' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/770/', 'Amtrak 770 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '770' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/774/', 'Amtrak 774 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '774' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/777/', 'Amtrak 777 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '777' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/779/', 'Amtrak 779 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '779' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/782/', 'Amtrak 782 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '782' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/784/', 'Amtrak 784 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '784' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/785/', 'Amtrak 785 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '785' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/786/', 'Amtrak 786 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '786' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/790/', 'Amtrak 790 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '790' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/791/', 'Amtrak 791 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '791' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/794/', 'Amtrak 794 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '794' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1562/', 'Amtrak 1562 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1562' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1591/', 'Amtrak 1591 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1591' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1595/', 'Amtrak 1595 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1595' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1765/', 'Amtrak 1765 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1765' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1769/', 'Amtrak 1769 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1769' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1770/', 'Amtrak 1770 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1770' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1774/', 'Amtrak 1774 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1774' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1777/', 'Amtrak 1777 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1777' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1784/', 'Amtrak 1784 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1784' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1785/', 'Amtrak 1785 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1785' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1790/', 'Amtrak 1790 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1790' AND c.name = 'Pacific Surfliner';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/90/', 'Amtrak 90 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '90' AND c.name = 'Amtrak Palmetto';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/42/', 'Amtrak 42 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '42' AND c.name = 'Pennsylvanian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/43/', 'Amtrak 43 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '43' AND c.name = 'Pennsylvanian';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/370/', 'Amtrak 370 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '370' AND c.name = 'Pere Marquette';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/371/', 'Amtrak 371 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '371' AND c.name = 'Pere Marquette';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/390/', 'Amtrak 390 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '390' AND c.name = 'Illini / Saluki';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/391/', 'Amtrak 391 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '391' AND c.name = 'Illini / Saluki';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/392/', 'Amtrak 392 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '392' AND c.name = 'Illini / Saluki';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/393/', 'Amtrak 393 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '393' AND c.name = 'Illini / Saluki';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/97/', 'Amtrak 97 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '97' AND c.name = 'Amtrak Silver Meteor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/98/', 'Amtrak 98 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '98' AND c.name = 'Amtrak Silver Meteor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1098/', 'Amtrak 1098 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1098' AND c.name = 'Amtrak Silver Meteor';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/89/', 'Amtrak 89 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '89' AND c.name = 'Amtrak Palmetto';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/3/', 'Amtrak 3 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '3' AND c.name = 'Amtrak Southwest Chief';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/4/', 'Amtrak 4 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '4' AND c.name = 'Amtrak Southwest Chief';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1003/', 'Amtrak 1003 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1003' AND c.name = 'Amtrak Southwest Chief';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1004/', 'Amtrak 1004 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1004' AND c.name = 'Amtrak Southwest Chief';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1/', 'Amtrak 1 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1' AND c.name = 'Amtrak Sunset Limited';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/2/', 'Amtrak 2 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '2' AND c.name = 'Amtrak Sunset Limited';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/21/', 'Amtrak 21 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '21' AND c.name = 'Amtrak Texas Eagle';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/22/', 'Amtrak 22 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '22' AND c.name = 'Amtrak Texas Eagle';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1021/', 'Amtrak 1021 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1021' AND c.name = 'Amtrak Texas Eagle';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1022/', 'Amtrak 1022 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1022' AND c.name = 'Amtrak Texas Eagle';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/54/', 'Amtrak 54 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '54' AND c.name = 'Vermonter';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/55/', 'Amtrak 55 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '55' AND c.name = 'Vermonter';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/56/', 'Amtrak 56 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '56' AND c.name = 'Vermonter';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/57/', 'Amtrak 57 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '57' AND c.name = 'Vermonter';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/107/', 'Amtrak 107 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '107' AND c.name = 'Vermonter';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1105/', 'Amtrak 1105 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1105' AND c.name = 'Winter Park Express';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1106/', 'Amtrak 1106 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1106' AND c.name = 'Winter Park Express';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/350/', 'Amtrak 350 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '350' AND c.name = 'Wolverine';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/351/', 'Amtrak 351 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '351' AND c.name = 'Wolverine';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/352/', 'Amtrak 352 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '352' AND c.name = 'Wolverine';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/353/', 'Amtrak 353 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '353' AND c.name = 'Wolverine';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/354/', 'Amtrak 354 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '354' AND c.name = 'Wolverine';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/355/', 'Amtrak 355 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '355' AND c.name = 'Wolverine';
INSERT INTO media (train_id, url, title, media_type, source_type, source_domain, is_published, added_by)
  SELECT t.id, 'https://railrat.net/trains/1354/', 'Amtrak 1354 live status on RailRat', 'website', 'seed', 'railrat.net', 1, 'seed'
  FROM trains t JOIN corridors c ON c.id = t.corridor_id
  WHERE t.train_number = '1354' AND c.name = 'Wolverine';
