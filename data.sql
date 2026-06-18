-- ============================================================
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
-- Note: above stops are Floridian extension -- Crescent proper ends below
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

