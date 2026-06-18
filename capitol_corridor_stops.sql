-- Stop IDs: ARN=131 RLN=132 RSV=133 SAC=134 DAV=135 FFV=136 SUI=137 MTZ=138
--           RIC=139 BKY=140 EMY=141 OKJ=142 OAC=143 HAY=144 FMT=145 GAC=146 SCC=147 SJC=148
-- Train IDs: 520=233 521=234 522=235 523=236 524=237 525=238 526=239 527=240 528=241
--            529=242 530=243 531=244 532=245 534=246 535=247 536=248 537=249 538=250
--            539=251 540=252 541=253 542=254 543=255 544=256 545=257 547=259 548=260
--            549=261 550=262 551=263 720=264 723=265 724=266 727=267 728=268 729=269
--            732=270 733=271 734=272 736=273 737=274 738=275 741=276 742=277 743=278
--            744=279 745=280 746=281 747=282 748=283 749=284 751=286

-- ── EASTBOUND short-turns: OKJ → SAC ────────────────────────────────────────
-- Train 520 (id=233)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(233,142,10,NULL,'5:28a'),(233,141,20,NULL,'5:38a'),(233,140,30,NULL,'5:42a'),
(233,139,40,NULL,'5:50a'),(233,138,50,NULL,'6:17a'),(233,137,60,NULL,'6:36a'),
(233,136,70,NULL,'6:42a'),(233,135,80,NULL,'7:06a'),(233,134,90,'7:32a',NULL);

-- Train 522 (id=235)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(235,142,10,NULL,'6:28a'),(235,141,20,NULL,'6:38a'),(235,140,30,NULL,'6:42a'),
(235,139,40,NULL,'6:50a'),(235,138,50,NULL,'7:17a'),(235,137,60,NULL,'7:36a'),
(235,136,70,NULL,'7:42a'),(235,135,80,NULL,'8:06a'),(235,134,90,'8:30a',NULL);

-- Train 526 (id=239)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(239,142,10,NULL,'8:58a'),(239,141,20,NULL,'9:08a'),(239,140,30,NULL,'9:12a'),
(239,139,40,NULL,'9:20a'),(239,138,50,NULL,'9:47a'),(239,137,60,NULL,'10:06a'),
(239,136,70,NULL,'10:12a'),(239,135,80,NULL,'10:36a'),(239,134,90,'11:00a',NULL);

-- Train 532 (id=245)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(245,142,10,NULL,'12:58p'),(245,141,20,NULL,'1:08p'),(245,140,30,NULL,'1:12p'),
(245,139,40,NULL,'1:20p'),(245,138,50,NULL,'1:47p'),(245,137,60,NULL,'2:06p'),
(245,136,70,NULL,'2:12p'),(245,135,80,NULL,'2:36p'),(245,134,90,'3:00p',NULL);

-- Train 534 (id=246)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(246,142,10,NULL,'1:58p'),(246,141,20,NULL,'2:08p'),(246,140,30,NULL,'2:12p'),
(246,139,40,NULL,'2:20p'),(246,138,50,NULL,'2:47p'),(246,137,60,NULL,'3:10p'),
(246,136,70,NULL,'3:16p'),(246,135,80,NULL,'3:40p'),(246,134,90,'4:05p',NULL);

-- Train 536 (id=248)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(248,142,10,NULL,'2:58p'),(248,141,20,NULL,'3:08p'),(248,140,30,NULL,'3:12p'),
(248,139,40,NULL,'3:20p'),(248,138,50,NULL,'3:47p'),(248,137,60,NULL,'4:06p'),
(248,136,70,NULL,'4:12p'),(248,135,80,NULL,'4:36p'),(248,134,90,'5:00p',NULL);

-- Train 540 (id=252)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(252,142,10,NULL,'4:22p'),(252,141,20,NULL,'4:32p'),(252,140,30,NULL,'4:36p'),
(252,139,40,NULL,'4:44p'),(252,138,50,NULL,'5:11p'),(252,137,60,NULL,'5:30p'),
(252,136,70,NULL,'5:40p'),(252,135,80,NULL,'6:02p'),(252,134,90,'6:25p',NULL);

-- Train 548 (id=260)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(260,142,10,NULL,'7:28p'),(260,141,20,NULL,'7:38p'),(260,140,30,NULL,'7:43p'),
(260,139,40,NULL,'7:51p'),(260,138,50,NULL,'8:17p'),(260,137,60,NULL,'8:36p'),
(260,136,70,NULL,'8:42p'),(260,135,80,NULL,'9:06p'),(260,134,90,'9:30p',NULL);

-- Train 720 (id=264, weekend)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(264,142,10,NULL,'7:48a'),(264,141,20,NULL,'7:58a'),(264,140,30,NULL,'8:02a'),
(264,139,40,NULL,'8:10a'),(264,138,50,NULL,'8:37a'),(264,137,60,NULL,'8:56a'),
(264,136,70,NULL,'9:02a'),(264,135,80,NULL,'9:26a'),(264,134,90,'9:50a',NULL);

-- Train 732 (id=270, weekend)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(270,142,10,NULL,'12:58p'),(270,141,20,NULL,'1:08p'),(270,140,30,NULL,'1:12p'),
(270,139,40,NULL,'1:20p'),(270,138,50,NULL,'1:47p'),(270,137,60,NULL,'2:06p'),
(270,136,70,NULL,'2:12p'),(270,135,80,NULL,'2:36p'),(270,134,90,'3:01p',NULL);

-- Train 736 (id=273, weekend)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(273,142,10,NULL,'2:58p'),(273,141,20,NULL,'3:08p'),(273,140,30,NULL,'3:12p'),
(273,139,40,NULL,'3:20p'),(273,138,50,NULL,'3:47p'),(273,137,60,NULL,'4:06p'),
(273,136,70,NULL,'4:12p'),(273,135,80,NULL,'4:36p'),(273,134,90,'5:02p',NULL);

-- Train 742 (id=277, weekend)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(277,142,10,NULL,'4:58p'),(277,141,20,NULL,'5:08p'),(277,140,30,NULL,'5:12p'),
(277,139,40,NULL,'5:20p'),(277,138,50,NULL,'5:47p'),(277,137,60,NULL,'6:06p'),
(277,136,70,NULL,'6:12p'),(277,135,80,NULL,'6:36p'),(277,134,90,'7:03p',NULL);

-- ── EASTBOUND full runs: SJC → SAC ──────────────────────────────────────────
-- Train 524 (id=237, Mo-Fr, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(237,148,10,NULL,'6:18a'),(237,147,20,NULL,'6:25a'),(237,146,30,NULL,'6:33a'),
(237,145,40,NULL,'6:50a'),(237,144,50,NULL,'7:10a'),(237,143,60,NULL,'7:20a'),
(237,142,70,NULL,'7:35a'),(237,141,80,NULL,'7:45a'),(237,140,90,NULL,'7:52a'),
(237,139,100,NULL,'8:00a'),(237,138,110,NULL,'8:29a'),(237,137,120,NULL,'8:48a'),
(237,136,130,NULL,'8:54a'),(237,135,140,NULL,'9:18a'),(237,134,150,'9:40a',NULL);

-- Train 528 (id=241, Mo-Fr, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(241,148,10,NULL,'8:38a'),(241,147,20,NULL,'8:46a'),(241,146,30,NULL,'8:53a'),
(241,145,40,NULL,'9:13a'),(241,144,50,NULL,'9:32a'),(241,143,60,NULL,'9:42a'),
(241,142,70,NULL,'9:59a'),(241,141,80,NULL,'10:09a'),(241,140,90,NULL,'10:13a'),
(241,139,100,NULL,'10:21a'),(241,138,110,NULL,'10:51a'),(241,137,120,NULL,'11:10a'),
(241,136,130,NULL,'11:16a'),(241,135,140,NULL,'11:40a'),(241,134,150,'12:04p',NULL);

-- Train 530 (id=243, Mo-Fr, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(243,148,10,NULL,'10:35a'),(243,147,20,NULL,'10:43a'),(243,146,30,NULL,'10:50a'),
(243,145,40,NULL,'11:08a'),(243,144,50,NULL,'11:23a'),(243,143,60,NULL,'11:33a'),
(243,142,70,NULL,'11:43a'),(243,141,80,NULL,'11:53a'),(243,140,90,NULL,'11:57a'),
(243,139,100,NULL,'12:05p'),(243,138,110,NULL,'12:32p'),(243,137,120,NULL,'12:51p'),
(243,136,130,NULL,'12:57p'),(243,135,140,NULL,'1:21p'),(243,134,150,'1:45p',NULL);

-- Train 542 (id=254, Mo-Fr, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(254,148,10,NULL,'4:17p'),(254,147,20,NULL,'4:25p'),(254,146,30,NULL,'4:32p'),
(254,145,40,NULL,'4:49p'),(254,144,50,NULL,'5:04p'),(254,143,60,NULL,'5:14p'),
(254,142,70,NULL,'5:23p'),(254,141,80,NULL,'5:33p'),(254,140,90,NULL,'5:38p'),
(254,139,100,NULL,'5:47p'),(254,138,110,NULL,'6:13p'),(254,137,120,NULL,'6:32p'),
(254,136,130,NULL,'6:38p'),(254,135,140,NULL,'7:02p'),(254,134,150,'7:26p',NULL);

-- Train 544 (id=256, Mo-Fr, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(256,148,10,NULL,'5:25p'),(256,147,20,NULL,'5:33p'),(256,146,30,NULL,'5:40p'),
(256,145,40,NULL,'5:58p'),(256,144,50,NULL,'6:13p'),(256,143,60,NULL,'6:23p'),
(256,142,70,NULL,'6:35p'),(256,141,80,NULL,'6:48p'),(256,140,90,NULL,'6:53p'),
(256,139,100,NULL,'7:01p'),(256,138,110,NULL,'7:28p'),(256,137,120,NULL,'7:47p'),
(256,136,130,NULL,'7:54p'),(256,135,140,NULL,'8:18p'),(256,134,150,'8:41p',NULL);

-- Train 550 (id=262, Mo-Fr, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(262,148,10,NULL,'6:55p'),(262,147,20,NULL,'7:03p'),(262,146,30,NULL,'7:10p'),
(262,145,40,NULL,'7:31p'),(262,144,50,NULL,'7:46p'),(262,143,60,NULL,'7:56p'),
(262,142,70,NULL,'8:05p'),(262,141,80,NULL,'8:15p'),(262,140,90,NULL,'8:19p'),
(262,139,100,NULL,'8:27p'),(262,138,110,NULL,'8:53p'),(262,137,120,NULL,'9:12p'),
(262,136,130,NULL,'9:18p'),(262,135,140,NULL,'9:42p'),(262,134,150,'10:05p',NULL);

-- Train 724 (id=266, weekend, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(266,148,10,NULL,'8:50a'),(266,147,20,NULL,'8:57a'),(266,146,30,NULL,'9:05a'),
(266,145,40,NULL,'9:23a'),(266,144,50,NULL,'9:38a'),(266,143,60,NULL,'9:48a'),
(266,142,70,NULL,'9:58a'),(266,141,80,NULL,'10:08a'),(266,140,90,NULL,'10:12a'),
(266,139,100,NULL,'10:20a'),(266,138,110,NULL,'10:47a'),(266,137,120,NULL,'11:06a'),
(266,136,130,NULL,'11:12a'),(266,135,140,NULL,'11:36a'),(266,134,150,'12:07p',NULL);

-- Train 728 (id=268, weekend, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(268,148,10,NULL,'10:50a'),(268,147,20,NULL,'10:57a'),(268,146,30,NULL,'11:05a'),
(268,145,40,NULL,'11:23a'),(268,144,50,NULL,'11:38a'),(268,143,60,NULL,'11:48a'),
(268,142,70,NULL,'11:58a'),(268,141,80,NULL,'12:08p'),(268,140,90,NULL,'12:12p'),
(268,139,100,NULL,'12:20p'),(268,138,110,NULL,'12:47p'),(268,137,120,NULL,'1:06p'),
(268,136,130,NULL,'1:12p'),(268,135,140,NULL,'1:36p'),(268,134,150,'2:05p',NULL);

-- Train 734 (id=272, weekend, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(272,148,10,NULL,'12:50p'),(272,147,20,NULL,'12:57p'),(272,146,30,NULL,'1:05p'),
(272,145,40,NULL,'1:23p'),(272,144,50,NULL,'1:38p'),(272,143,60,NULL,'1:48p'),
(272,142,70,NULL,'1:58p'),(272,141,80,NULL,'2:08p'),(272,140,90,NULL,'2:12p'),
(272,139,100,NULL,'2:20p'),(272,138,110,NULL,'2:47p'),(272,137,120,NULL,'3:06p'),
(272,136,130,NULL,'3:12p'),(272,135,140,NULL,'3:36p'),(272,134,150,'4:05p',NULL);

-- Train 738 (id=275, weekend, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(275,148,10,NULL,'2:50p'),(275,147,20,NULL,'2:57p'),(275,146,30,NULL,'3:05p'),
(275,145,40,NULL,'3:23p'),(275,144,50,NULL,'3:38p'),(275,143,60,NULL,'3:48p'),
(275,142,70,NULL,'3:58p'),(275,141,80,NULL,'4:08p'),(275,140,90,NULL,'4:12p'),
(275,139,100,NULL,'4:20p'),(275,138,110,NULL,'4:47p'),(275,137,120,NULL,'5:06p'),
(275,136,130,NULL,'5:12p'),(275,135,140,NULL,'5:36p'),(275,134,150,'6:05p',NULL);

-- Train 746 (id=281, weekend, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(281,148,10,NULL,'5:50p'),(281,147,20,NULL,'5:57p'),(281,146,30,NULL,'6:05p'),
(281,145,40,NULL,'6:23p'),(281,144,50,NULL,'6:38p'),(281,143,60,NULL,'6:49p'),
(281,142,70,NULL,'7:01p'),(281,141,80,NULL,'7:11p'),(281,140,90,NULL,'7:15p'),
(281,139,100,NULL,'7:23p'),(281,138,110,NULL,'7:49p'),(281,137,120,NULL,'8:08p'),
(281,136,130,NULL,'8:14p'),(281,135,140,NULL,'8:38p'),(281,134,150,'9:05p',NULL);

-- Train 748 (id=283, weekend, SJC→SAC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(283,148,10,NULL,'8:05p'),(283,147,20,NULL,'8:10p'),(283,146,30,NULL,'8:18p'),
(283,145,40,NULL,'8:35p'),(283,144,50,NULL,'8:50p'),(283,143,60,NULL,'9:00p'),
(283,142,70,NULL,'9:11p'),(283,141,80,NULL,'9:21p'),(283,140,90,NULL,'9:25p'),
(283,139,100,NULL,'9:33p'),(283,138,110,NULL,'9:59p'),(283,137,120,NULL,'10:18p'),
(283,136,130,NULL,'10:24p'),(283,135,140,NULL,'10:48p'),(283,134,150,'11:15p',NULL);

-- ── EASTBOUND long runs: SJC → ARN ──────────────────────────────────────────
-- Train 538 (id=250, Mo-Fr, SJC→ARN)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(250,148,10,NULL,'2:50p'),(250,147,20,NULL,'2:58p'),(250,146,30,NULL,'3:07p'),
(250,145,40,NULL,'3:26p'),(250,144,50,NULL,'3:41p'),(250,143,60,NULL,'3:51p'),
(250,142,70,NULL,'4:02p'),(250,141,80,NULL,'4:13p'),(250,140,90,NULL,'4:17p'),
(250,139,100,NULL,'4:25p'),(250,138,110,NULL,'4:52p'),(250,137,120,NULL,'5:11p'),
(250,136,130,NULL,'5:17p'),(250,135,140,NULL,'5:41p'),
(250,134,150,'6:03p','6:03p'),
(250,133,160,NULL,'6:28p'),(250,132,170,NULL,'6:38p'),(250,131,180,'7:10p',NULL);

-- Train 744 (id=279, weekend, SJC→ARN)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(279,148,10,NULL,'4:50p'),(279,147,20,NULL,'4:57p'),(279,146,30,NULL,'5:05p'),
(279,145,40,NULL,'5:23p'),(279,144,50,NULL,'5:38p'),(279,143,60,NULL,'5:48p'),
(279,142,70,NULL,'5:58p'),(279,141,80,NULL,'6:08p'),(279,140,90,NULL,'6:12p'),
(279,139,100,NULL,'6:20p'),(279,138,110,NULL,'6:47p'),(279,137,120,NULL,'7:06p'),
(279,136,130,NULL,'7:12p'),(279,135,140,NULL,'7:36p'),
(279,134,150,'8:00p','8:00p'),
(279,133,160,NULL,'8:23p'),(279,132,170,NULL,'8:32p'),(279,131,180,'9:03p',NULL);

-- ── WESTBOUND: SAC → OKJ (short-turns) ──────────────────────────────────────
-- Train 525 (id=238, Mo-Fr, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(238,134,10,NULL,'6:03a'),(238,135,20,NULL,'6:19a'),(238,136,30,NULL,'6:39a'),
(238,137,40,NULL,'6:45a'),(238,138,50,NULL,'7:05a'),(238,139,60,NULL,'7:31a'),
(238,140,70,NULL,'7:39a'),(238,141,80,NULL,'7:45a'),(238,142,90,'8:00a',NULL);

-- Train 531 (id=244, Mo-Fr, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(244,134,10,NULL,'8:53a'),(244,135,20,NULL,'9:09a'),(244,136,30,NULL,'9:29a'),
(244,137,40,NULL,'9:35a'),(244,138,50,NULL,'9:54a'),(244,139,60,NULL,'10:21a'),
(244,140,70,NULL,'10:29a'),(244,141,80,NULL,'10:36a'),(244,142,90,'10:50a',NULL);

-- Train 535 (id=247, Mo-Fr, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(247,134,10,NULL,'9:53a'),(247,135,20,NULL,'10:09a'),(247,136,30,NULL,'10:29a'),
(247,137,40,NULL,'10:35a'),(247,138,50,NULL,'10:54a'),(247,139,60,NULL,'11:21a'),
(247,140,70,NULL,'11:29a'),(247,141,80,NULL,'11:36a'),(247,142,90,'11:50a',NULL);

-- Train 537 (id=249, Mo-Fr, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(249,134,10,NULL,'10:53a'),(249,135,20,NULL,'11:09a'),(249,136,30,NULL,'11:29a'),
(249,137,40,NULL,'11:35a'),(249,138,50,NULL,'11:54a'),(249,139,60,NULL,'12:21p'),
(249,140,70,NULL,'12:29p'),(249,141,80,NULL,'12:36p'),(249,142,90,'12:50p',NULL);

-- Train 541 (id=253, Mo-Fr, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(253,134,10,NULL,'12:53p'),(253,135,20,NULL,'1:09p'),(253,136,30,NULL,'1:29p'),
(253,137,40,NULL,'1:35p'),(253,138,50,NULL,'1:56p'),(253,139,60,NULL,'2:25p'),
(253,140,70,NULL,'2:33p'),(253,141,80,NULL,'2:39p'),(253,142,90,'2:50p',NULL);

-- Train 549 (id=261, Mo-Fr, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(261,134,10,NULL,'7:03p'),(261,135,20,NULL,'7:19p'),(261,136,30,NULL,'7:39p'),
(261,137,40,NULL,'7:45p'),(261,138,50,NULL,'8:04p'),(261,139,60,NULL,'8:31p'),
(261,140,70,NULL,'8:39p'),(261,141,80,NULL,'8:46p'),(261,142,90,'9:00p',NULL);

-- Train 551 (id=263, Mo-Fr, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(263,134,10,NULL,'9:03p'),(263,135,20,NULL,'9:19p'),(263,136,30,NULL,'9:39p'),
(263,137,40,NULL,'9:47p'),(263,138,50,NULL,'10:06p'),(263,139,60,NULL,'10:33p'),
(263,140,70,NULL,'10:41p'),(263,141,80,NULL,'10:48p'),(263,142,90,'11:03p',NULL);

-- Train 737 (id=274, weekend, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(274,134,10,NULL,'11:53a'),(274,135,20,NULL,'12:09p'),(274,136,30,NULL,'12:29p'),
(274,137,40,NULL,'12:35p'),(274,138,50,NULL,'12:54p'),(274,139,60,NULL,'1:21p'),
(274,140,70,NULL,'1:29p'),(274,141,80,NULL,'1:35p'),(274,142,90,'1:50p',NULL);

-- Train 743 (id=278, weekend, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(278,134,10,NULL,'1:53p'),(278,135,20,NULL,'2:09p'),(278,136,30,NULL,'2:29p'),
(278,137,40,NULL,'2:35p'),(278,138,50,NULL,'2:54p'),(278,139,60,NULL,'3:21p'),
(278,140,70,NULL,'3:29p'),(278,141,80,NULL,'3:36p'),(278,142,90,'3:54p',NULL);

-- Train 749 (id=284, weekend, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(284,134,10,NULL,'6:53p'),(284,135,20,NULL,'7:09p'),(284,136,30,NULL,'7:29p'),
(284,137,40,NULL,'7:35p'),(284,138,50,NULL,'7:54p'),(284,139,60,NULL,'8:21p'),
(284,140,70,NULL,'8:29p'),(284,141,80,NULL,'8:36p'),(284,142,90,'8:53p',NULL);

-- Train 751 (id=286, weekend, SAC→OKJ)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(286,134,10,NULL,'8:53p'),(286,135,20,NULL,'9:09p'),(286,136,30,NULL,'9:29p'),
(286,137,40,NULL,'9:35p'),(286,138,50,NULL,'9:54p'),(286,139,60,NULL,'10:21p'),
(286,140,70,NULL,'10:29p'),(286,141,80,NULL,'10:36p'),(286,142,90,'10:49p',NULL);

-- ── WESTBOUND: SAC → SJC (full runs, skipping OAC/HAY/FMT) ─────────────────
-- Train 527 (id=240, Mo-Fr, SAC→SJC via express: skips OAC HAY FMT)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(240,134,10,NULL,'6:43a'),(240,135,20,NULL,'6:58a'),(240,136,30,NULL,'7:18a'),
(240,137,40,NULL,'7:24a'),(240,138,50,NULL,'7:43a'),(240,139,60,NULL,'8:09a'),
(240,140,70,NULL,'8:17a'),(240,141,80,NULL,'8:23a'),(240,142,90,NULL,'8:39a'),
(240,146,100,NULL,'9:31a'),(240,147,110,NULL,'9:40a'),(240,148,120,'9:54a',NULL);

-- ── WESTBOUND: SAC → SJC (full runs, all stops) ─────────────────────────────
-- Train 521 (id=234, Mo-Fr, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(234,134,10,NULL,'4:20a'),(234,135,20,NULL,'4:36a'),(234,136,30,NULL,'4:56a'),
(234,137,40,NULL,'5:02a'),(234,138,50,NULL,'5:21a'),(234,139,60,NULL,'5:48a'),
(234,140,70,NULL,'5:56a'),(234,141,80,NULL,'6:03a'),(234,142,90,NULL,'6:13a'),
(234,143,100,NULL,'6:22a'),(234,144,110,NULL,'6:32a'),(234,145,120,NULL,'6:49a'),
(234,146,130,NULL,'7:05a'),(234,147,140,NULL,'7:14a'),(234,148,150,'7:30a',NULL);

-- Train 523 (id=236, Mo-Fr, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(236,134,10,NULL,'5:23a'),(236,135,20,NULL,'5:39a'),(236,136,30,NULL,'5:59a'),
(236,137,40,NULL,'6:05a'),(236,138,50,NULL,'6:24a'),(236,139,60,NULL,'6:51a'),
(236,140,70,NULL,'6:59a'),(236,141,80,NULL,'7:06a'),(236,142,90,NULL,'7:16a'),
(236,143,100,NULL,'7:25a'),(236,144,110,NULL,'7:35a'),(236,145,120,NULL,'7:52a'),
(236,146,130,NULL,'8:08a'),(236,147,140,NULL,'8:17a'),(236,148,150,'8:32a',NULL);

-- Train 539 (id=251, Mo-Fr, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(251,134,10,NULL,'11:53a'),(251,135,20,NULL,'12:09p'),(251,136,30,NULL,'12:29p'),
(251,137,40,NULL,'12:35p'),(251,138,50,NULL,'12:54p'),(251,139,60,NULL,'1:21p'),
(251,140,70,NULL,'1:29p'),(251,141,80,NULL,'1:36p'),(251,142,90,NULL,'1:46p'),
(251,143,100,NULL,'1:55p'),(251,144,110,NULL,'2:05p'),(251,145,120,NULL,'2:22p'),
(251,146,130,NULL,'2:38p'),(251,147,140,NULL,'2:47p'),(251,148,150,'3:04p',NULL);

-- Train 543 (id=255, Mo-Fr, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(255,134,10,NULL,'2:23p'),(255,135,20,NULL,'2:39p'),(255,136,30,NULL,'2:59p'),
(255,137,40,NULL,'3:05p'),(255,138,50,NULL,'3:23p'),(255,139,60,NULL,'3:49p'),
(255,140,70,NULL,'3:57p'),(255,141,80,NULL,'4:02p'),(255,142,90,NULL,'4:15p'),
(255,143,100,NULL,'4:24p'),(255,144,110,NULL,'4:33p'),(255,145,120,NULL,'4:50p'),
(255,146,130,NULL,'5:11p'),(255,147,140,NULL,'5:20p'),(255,148,150,'5:35p',NULL);

-- Train 545 (id=257, Mo-Fr, SAC→SJC — using Tue/Fri schedule which runs full route)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(257,134,10,NULL,'3:20p'),(257,135,20,NULL,'3:36p'),(257,136,30,NULL,'3:55p'),
(257,137,40,NULL,'4:01p'),(257,138,50,NULL,'4:21p'),(257,139,60,NULL,'4:48p'),
(257,140,70,NULL,'4:57p'),(257,141,80,NULL,'5:04p'),(257,142,90,NULL,'5:20p'),
(257,143,100,NULL,'5:29p'),(257,144,110,NULL,'5:39p'),(257,145,120,NULL,'5:57p'),
(257,146,130,NULL,'6:14p'),(257,147,140,NULL,'6:23p'),(257,148,150,'6:41p',NULL);

-- Train 547 (id=259, Mo-Fr, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(259,134,10,NULL,'5:01p'),(259,135,20,NULL,'5:17p'),(259,136,30,NULL,'5:37p'),
(259,137,40,NULL,'5:43p'),(259,138,50,NULL,'6:02p'),(259,139,60,NULL,'6:29p'),
(259,140,70,NULL,'6:37p'),(259,141,80,NULL,'6:43p'),(259,142,90,NULL,'6:53p'),
(259,143,100,NULL,'7:02p'),(259,144,110,NULL,'7:12p'),(259,145,120,NULL,'7:30p'),
(259,146,130,NULL,'7:46p'),(259,147,140,NULL,'7:55p'),(259,148,150,'8:14p',NULL);

-- Train 723 (id=265, weekend, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(265,134,10,NULL,'6:03a'),(265,135,20,NULL,'6:19a'),(265,136,30,NULL,'6:39a'),
(265,137,40,NULL,'6:45a'),(265,138,50,NULL,'7:04a'),(265,139,60,NULL,'7:31a'),
(265,140,70,NULL,'7:39a'),(265,141,80,NULL,'7:45a'),(265,142,90,NULL,'7:57a'),
(265,143,100,NULL,'8:06a'),(265,144,110,NULL,'8:16a'),(265,145,120,NULL,'8:32a'),
(265,146,130,NULL,'8:49a'),(265,147,140,NULL,'8:57a'),(265,148,150,'9:14a',NULL);

-- Train 727 (id=267, weekend, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(267,134,10,NULL,'7:26a'),(267,135,20,NULL,'7:42a'),(267,136,30,NULL,'8:02a'),
(267,137,40,NULL,'8:08a'),(267,138,50,NULL,'8:27a'),(267,139,60,NULL,'8:54a'),
(267,140,70,NULL,'9:02a'),(267,141,80,NULL,'9:08a'),(267,142,90,NULL,'9:20a'),
(267,143,100,NULL,'9:29a'),(267,144,110,NULL,'9:39a'),(267,145,120,NULL,'9:58a'),
(267,146,130,NULL,'10:14a'),(267,147,140,NULL,'10:22a'),(267,148,150,'10:41a',NULL);

-- Train 733 (id=271, weekend, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(271,134,10,NULL,'9:53a'),(271,135,20,NULL,'10:09a'),(271,136,30,NULL,'10:29a'),
(271,137,40,NULL,'10:35a'),(271,138,50,NULL,'10:54a'),(271,139,60,NULL,'11:21a'),
(271,140,70,NULL,'11:29a'),(271,141,80,NULL,'11:35a'),(271,142,90,NULL,'11:47a'),
(271,143,100,NULL,'11:56a'),(271,144,110,NULL,'12:06p'),(271,145,120,NULL,'12:22p'),
(271,146,130,NULL,'12:39p'),(271,147,140,NULL,'12:47p'),(271,148,150,'1:04p',NULL);

-- Train 741 (id=276, weekend, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(276,134,10,NULL,'12:53p'),(276,135,20,NULL,'1:09p'),(276,136,30,NULL,'1:29p'),
(276,137,40,NULL,'1:35p'),(276,138,50,NULL,'1:54p'),(276,139,60,NULL,'2:21p'),
(276,140,70,NULL,'2:29p'),(276,141,80,NULL,'2:36p'),(276,142,90,NULL,'2:46p'),
(276,143,100,NULL,'2:55p'),(276,144,110,NULL,'3:05p'),(276,145,120,NULL,'3:22p'),
(276,146,130,NULL,'3:38p'),(276,147,140,NULL,'3:47p'),(276,148,150,'4:02p',NULL);

-- Train 745 (id=280, weekend, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(280,134,10,NULL,'3:53p'),(280,135,20,NULL,'4:09p'),(280,136,30,NULL,'4:29p'),
(280,137,40,NULL,'4:35p'),(280,138,50,NULL,'4:54p'),(280,139,60,NULL,'5:21p'),
(280,140,70,NULL,'5:29p'),(280,141,80,NULL,'5:35p'),(280,142,90,NULL,'5:46p'),
(280,143,100,NULL,'5:55p'),(280,144,110,NULL,'6:05p'),(280,145,120,NULL,'6:22p'),
(280,146,130,NULL,'6:39p'),(280,147,140,NULL,'6:47p'),(280,148,150,'7:05p',NULL);

-- Train 747 (id=282, weekend, SAC→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(282,134,10,NULL,'5:03p'),(282,135,20,NULL,'5:19p'),(282,136,30,NULL,'5:39p'),
(282,137,40,NULL,'5:45p'),(282,138,50,NULL,'6:04p'),(282,139,60,NULL,'6:31p'),
(282,140,70,NULL,'6:39p'),(282,141,80,NULL,'6:46p'),(282,142,90,NULL,'6:56p'),
(282,143,100,NULL,'7:05p'),(282,144,110,NULL,'7:15p'),(282,145,120,NULL,'7:32p'),
(282,146,130,NULL,'7:48p'),(282,147,140,NULL,'7:57p'),(282,148,150,'8:14p',NULL);

-- ── WESTBOUND long runs: ARN → SJC ──────────────────────────────────────────
-- Train 529 (id=242, Mo-Fr, ARN→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(242,131,10,NULL,'6:36a'),(242,132,20,NULL,'7:00a'),(242,133,30,NULL,'7:09a'),
(242,134,40,NULL,'7:41a'),(242,135,50,NULL,'7:57a'),(242,136,60,NULL,'8:16a'),
(242,137,70,NULL,'8:23a'),(242,138,80,NULL,'8:43a'),(242,139,90,NULL,'9:10a'),
(242,140,100,NULL,'9:18a'),(242,141,110,NULL,'9:25a'),(242,142,120,NULL,'9:36a'),
(242,143,130,NULL,'9:45a'),(242,144,140,NULL,'9:54a'),(242,145,150,NULL,'10:11a'),
(242,146,160,NULL,'10:28a'),(242,147,170,NULL,'10:39a'),(242,148,180,'11:02a',NULL);

-- Train 729 (id=269, weekend, ARN→SJC)
INSERT INTO train_stops (train_id,stop_id,sort_order,scheduled_arrival,scheduled_departure) VALUES
(269,131,10,NULL,'7:51a'),(269,132,20,NULL,'8:14a'),(269,133,30,NULL,'8:23a'),
(269,134,40,NULL,'8:53a'),(269,135,50,NULL,'9:09a'),(269,136,60,NULL,'9:29a'),
(269,137,70,NULL,'9:35a'),(269,138,80,NULL,'9:54a'),(269,139,90,NULL,'10:21a'),
(269,140,100,NULL,'10:29a'),(269,141,110,NULL,'10:35a'),(269,142,120,NULL,'10:47a'),
(269,143,130,NULL,'10:56a'),(269,144,140,NULL,'11:07a'),(269,145,150,NULL,'11:23a'),
(269,146,160,NULL,'11:39a'),(269,147,170,NULL,'11:47a'),(269,148,180,'12:05p',NULL);

-- Set weekday/weekend flags based on train number series
UPDATE train_stops SET runs_weekday=1, runs_weekend=0
WHERE train_id IN (
  SELECT t.id FROM trains t
  JOIN corridors c ON t.corridor_id=c.id
  WHERE c.slug='capitol-corridor'
  AND CAST(t.train_number AS INTEGER) BETWEEN 500 AND 599
);

UPDATE train_stops SET runs_weekday=0, runs_weekend=1
WHERE train_id IN (
  SELECT t.id FROM trains t
  JOIN corridors c ON t.corridor_id=c.id
  WHERE c.slug='capitol-corridor'
  AND CAST(t.train_number AS INTEGER) BETWEEN 700 AND 799
);
