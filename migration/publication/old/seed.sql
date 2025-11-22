-- ============================================
-- SEED DATA UNTUK SCHEMA public.-- ============================================

BEGIN;

-- --------------------------------------------
-- 1. COUNTRY
-- --------------------------------------------
INSERT INTO public.country (cntry_id, cntry_name_es, cntry_name_en, cntry_name_gl, cntry_aneca_id)
VALUES
  ('ES', 'España',        'Spain',        'España',        34),
  ('US', 'Estados Unidos','United States','Estados Unidos', 1),
  ('UK', 'Reino Unido',   'United Kingdom','Reino Unido', 44),
  ('FR', 'Francia',       'France',       'Francia',       33),
  ('DE', 'Alemania',      'Germany',      'Alemaña',       49);

-- --------------------------------------------
-- 2. RESEARCHER (rese_id 1..12)
-- --------------------------------------------
INSERT INTO public.researcher
(rese_id, rese_nif, rese_name, rese_first_surname, rese_sec_surname,
 rese_birth_day, rese_mail, rese_web, rese_phone, rese_ext,
 rese_fax, rese_univ, rese_sign, rese_coord, rese_photo,
 rese_gend, rese_pos, rese_title, rese_fig,
 rese_short_curr_es, rese_short_curr_en, rese_short_curr_gl)
VALUES
  (1, '12345678A', 'Luis',   'García',    'Pérez',   '1978-03-12',
   'lgarcia@example.es', 'http://lgarcia.es', '+34-600111111', '101', NULL,
   'Univ. Vigo', 'LG', false, NULL, 1, 1, NULL, NULL,
   'Profesor en informática', 'Professor in computer science', NULL),

  (2, '23456789B', 'María',  'López',     'Santos',  '1982-07-25',
   'mlopez@example.es', 'http://mlopez.es', '+34-600222222', '102', NULL,
   'Univ. Vigo', 'ML', true, NULL, 2, 1, NULL, NULL,
   'Investigadora en IA', 'Researcher in AI', NULL),

  (3, '34567890C', 'André',  'Dubois',    NULL,      '1975-11-05',
   'adubois@example.fr', NULL, '+33-610333333', NULL, NULL,
   'Univ. Lyon', 'AD', false, NULL, 1, 2, NULL, NULL,
   'Experto en redes', 'Networks expert', NULL),

  (4, '45678901D', 'Clara',  'Meier',     NULL,      '1985-02-18',
   'cmeier@example.de', NULL, '+49-1704444444', NULL, NULL,
   'TU Berlin', 'CM', false, NULL, 2, 2, NULL, NULL,
   'Especialista en bases de datos', 'Database specialist', NULL),

  (5, '56789012E', 'John',   'Smith',     NULL,      '1970-09-09',
   'jsmith@example.com', NULL, '+1-4155550101', NULL, NULL,
   'UC Berkeley', 'JS', true, NULL, 1, 1, NULL, NULL,
   'Profesor de sistemas distribuidos', 'Distributed systems professor', NULL),

  (6, '67890123F', 'Ana',    'González',  'Rivas',   '1990-01-30',
   'agonzalez@example.es', NULL, '+34-600333333', NULL, NULL,
   'Univ. Coruña', 'AG', false, NULL, 2, 3, NULL, NULL,
   'Doctora en ingeniería del software', 'Software engineering PhD', NULL),

  (7, '78901234G', 'Pedro',  'Alonso',    'Lago',    '1988-06-14',
   'palonso@example.es', NULL, '+34-600444444', NULL, NULL,
   'Univ. Vigo', 'PA', false, NULL, 1, 3, NULL, NULL,
   'Investigador postdoctoral', 'Postdoc researcher', NULL),

  (8, '89012345H', 'Sofia',  'Klein',     NULL,      '1992-04-21',
   'sklein@example.de', NULL, '+49-1705555555', NULL, NULL,
   'TU Munich', 'SK', false, NULL, 2, 3, NULL, NULL,
   'Investigadora en aprendizaje automático', 'ML researcher', NULL),

  (9, '90123456I', 'Daniel', 'Ruiz',      'Martín',  '1983-12-02',
   'druiz@example.es', NULL, '+34-600555555', NULL, NULL,
   'Univ. Madrid', 'DR', false, NULL, 1, 2, NULL, NULL,
   'Profesor ayudante', 'Assistant professor', NULL),

  (10,'01234567J', 'Laura',  'Fernández', 'Gómez',   '1987-05-08',
   'lfernandez@example.es', NULL, '+34-600666666', NULL, NULL,
   'Univ. Sevilla', 'LF', false, NULL, 2, 2, NULL, NULL,
   'Investigadora en minería de datos', 'Data mining researcher', NULL),

  (11,'11223344K', 'Miguel', 'Torres',    NULL,      '1979-01-19',
   'mtorres@example.es', NULL, '+34-600777777', NULL, NULL,
   'Univ. Zaragoza', 'MT', false, NULL, 1, 2, NULL, NULL,
   'Experto en sistemas embebidos', 'Embedded systems expert', NULL),

  (12,'22334455L', 'Elena',  'Costa',     NULL,      '1991-08-03',
   'ecosta@example.pt', NULL, '+351-910000000', NULL, NULL,
   'Univ. Porto', 'EC', false, NULL, 2, 3, NULL, NULL,
   'Investigadora en visualización', 'Visualization researcher', NULL);

-- Sesuaikan sequence researcher setelah insert manual
SELECT pg_catalog.setval('public.researcher_rese_id_seq', 12, true);

-- --------------------------------------------
-- 3. PUBLICATION (pub_id 1..30)
-- 1-10  : journal
-- 11-20 : conference
-- 21-25 : book
-- 26-30 : book chapter
-- --------------------------------------------

INSERT INTO public."publication"
(pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit,
 pub_subj, pub_draft, pub_prsnt, pub_minute, pub_source,
 pub_book, pub_link, pub_repo, pub_visib, pub_cntry)
VALUES
-- Journals (1..10)
(1, NULL, 'Scalable Data Migration in Heterogeneous Systems', 2015, 'Vigo',
 'Springer', 1, NULL, NULL, NULL, 'DB Journal', NULL, NULL, NULL, true, 'ES'),
(2, NULL, 'An Empirical Study on Change Data Capture', 2016, 'Madrid',
 'Elsevier', 1, NULL, NULL, NULL, 'Journal of Information Systems', NULL, NULL, NULL, true, 'ES'),
(3, NULL, 'Improving CDC Pipelines with Kafka', 2017, 'Berlin',
 'ACM', 2, NULL, NULL, NULL, 'Distributed Systems Journal', NULL, NULL, NULL, true, 'DE'),
(4, NULL, 'Performance Evaluation of Debezium Connectors', 2018, 'Paris',
 'IEEE', 2, NULL, NULL, NULL, 'IEEE Data Engineering', NULL, NULL, NULL, true, 'FR'),
(5, NULL, 'A Survey on Database Modernization', 2019, 'London',
 'Springer', 3, NULL, NULL, NULL, 'Software Engineering Journal', NULL, NULL, NULL, true, 'UK'),
(6, NULL, 'Streaming ETL for public.Systems', 2020, 'San Francisco',
 'O''Reilly', 3, NULL, NULL, NULL, 'US Data Science Journal', NULL, NULL, NULL, true, 'US'),
(7, NULL, 'Near Real-Time Data Integration', 2021, 'Munich',
 'IEEE', 4, NULL, NULL, NULL, 'Transactions on Big Data', NULL, NULL, NULL, true, 'DE'),
(8, NULL, 'Designing Pivot Databases for Migration', 2022, 'Barcelona',
 'Elsevier', 4, NULL, NULL, NULL, 'International Journal of Databases', NULL, NULL, NULL, true, 'ES'),
(9, NULL, 'A DSL for Data Migration: DAMI Framework', 2023, 'Lyon',
 'Springer', 5, NULL, NULL, NULL, 'Software & Systems Modeling', NULL, NULL, NULL, true, 'FR'),
(10,NULL, 'Benchmarking CDC Approaches in Microservices', 2024, 'Seattle',
 'ACM', 5, NULL, NULL, NULL, 'Journal of Cloud Computing', NULL, NULL, NULL, true, 'US'),

-- Conferences (11..20)
(11,NULL, 'Proceedings of the CDC & Streaming Workshop 2018', 2018, 'Vigo',
 'ACM', 6, NULL, NULL, NULL, 'CDC Workshop', NULL, NULL, NULL, true, 'ES'),
(12,NULL, 'Proceedings of the Data Migration Summit 2019', 2019, 'Berlin',
 'IEEE', 6, NULL, NULL, NULL, 'Data Migration Summit', NULL, NULL, NULL, true, 'DE'),
(13,NULL, 'Proceedings of the ETL and Integration Conf 2020', 2020, 'Paris',
 'Springer', 7, NULL, NULL, NULL, 'ETL & Integration Conf', NULL, NULL, NULL, true, 'FR'),
(14,NULL, 'Proceedings of the Debezium User Conference 2021', 2021, 'Madrid',
 'Red Hat', 7, NULL, NULL, NULL, 'Debezium Conf', NULL, NULL, NULL, true, 'ES'),
(15,NULL, 'Proceedings of Kafka and Friends 2022', 2022, 'London',
 'Confluent', 7, NULL, NULL, NULL, 'Kafka & Friends', NULL, NULL, NULL, true, 'UK'),
(16,NULL, 'Proceedings of Distributed Systems Conf 2020', 2020, 'San Jose',
 'IEEE', 8, NULL, NULL, NULL, 'Distributed Systems Conf', NULL, NULL, NULL, true, 'US'),
(17,NULL, 'Proceedings of Big Data Europe 2021', 2021, 'Berlin',
 'Springer', 8, NULL, NULL, NULL, 'Big Data Europe', NULL, NULL, NULL, true, 'DE'),
(18,NULL, 'Proceedings of Cloud Native Databases 2022', 2022, 'Amsterdam',
 'ACM', 9, NULL, NULL, NULL, 'Cloud Native DB Conf', NULL, NULL, NULL, true, 'FR'),
(19,NULL, 'Proceedings of Real-Time Analytics 2023', 2023, 'New York',
 'O''Reilly', 9, NULL, NULL, NULL, 'Real-Time Analytics', NULL, NULL, NULL, true, 'US'),
(20,NULL, 'Proceedings of Modernization and Migration 2024', 2024, 'Lisbon',
 'IEEE', 9, NULL, NULL, NULL, 'Modernization & Migration', NULL, NULL, NULL, true, 'ES'),

-- Books (21..25)
(21,'978-3-16-148410-0', 'public.Systems Modernization Handbook', 2014, 'Madrid',
 'Springer', 10, NULL, NULL, NULL, 'Handbook Series', NULL, NULL, NULL, true, 'ES'),
(22,'978-1-491-99999-1', 'Streaming Architectures with Kafka', 2017, 'San Francisco',
 'O''Reilly', 10, NULL, NULL, NULL, 'Professional Series', NULL, NULL, NULL, true, 'US'),
(23,'978-0-12-345678-9', 'Data Integration Patterns', 2018, 'London',
 'Elsevier', 11, NULL, NULL, NULL, 'Patterns Series', NULL, NULL, NULL, true, 'UK'),
(24,'978-1-234-56789-7', 'Microservices and Databases', 2019, 'Berlin',
 'ACM', 11, NULL, NULL, NULL, 'Software Architecture Series', NULL, NULL, NULL, true, 'DE'),
(25,'978-9-876-54321-0', 'Hands-On CDC with Debezium', 2021, 'Paris',
 'Red Hat', 11, NULL, NULL, NULL, 'Practical Guides', NULL, NULL, NULL, true, 'FR'),

-- Book Chapters (26..30)
(26,NULL, 'Chapter: Designing Pivot Tables', 2014, 'Madrid',
 'Springer', 12, NULL, NULL, NULL, 'public.Systems Modernization Handbook', NULL, NULL, NULL, true, 'ES'),
(27,NULL, 'Chapter: Kafka Connect in Practice', 2017, 'San Francisco',
 'O''Reilly', 12, NULL, NULL, NULL, 'Streaming Architectures with Kafka', NULL, NULL, NULL, true, 'US'),
(28,NULL, 'Chapter: CDC Design Patterns', 2018, 'London',
 'Elsevier', 12, NULL, NULL, NULL, 'Data Integration Patterns', NULL, NULL, NULL, true, 'UK'),
(29,NULL, 'Chapter: Microservices Data Ownership', 2019, 'Berlin',
 'ACM', 12, NULL, NULL, NULL, 'Microservices and Databases', NULL, NULL, NULL, true, 'DE'),
(30,NULL, 'Chapter: Debezium in Production', 2021, 'Paris',
 'Red Hat', 12, NULL, NULL, NULL, 'Hands-On CDC with Debezium', NULL, NULL, NULL, true, 'FR');

-- Sesuaikan sequence publication setelah insert manual
SELECT pg_catalog.setval('public.publication_pub_id_seq', 30, true);

-- --------------------------------------------
-- 4. JOURNAL (pub_id 1..10)
-- --------------------------------------------
INSERT INTO public.journal
(jrnl_id, jrnl_journ_title, jrnl_vol, jrnl_numb,
 jrnl_start_page, jrnl_end_page, jrnl_month,
 jrnl_scope, jrnl_indexed, jrnl_pos_numb, jrnl_tot_pos)
VALUES
(1, 'Journal of Data Migration',          '21', '3',  1, 20, 3,  1, true,  5,  50),
(2, 'Journal of Information Systems',     '10', '2', 21, 45, 4,  1, true, 10, 100),
(3, 'Distributed Systems Journal',        '8',  '1',  5, 25, 2,  2, true,  8,  80),
(4, 'IEEE Data Engineering',              '35', '4', 30, 60, 6,  1, true, 15, 200),
(5, 'Software Engineering Journal',       '18', '1', 10, 35, 1,  3, true, 20, 250),
(6, 'US Data Science Journal',            '5',  '3', 40, 65, 7,  2, false, 50, 300),
(7, 'IEEE Transactions on Big Data',      '7',  '2', 15, 40, 5,  1, true, 12, 150),
(8, 'International Journal of Databases', '12', '1',  5, 30, 2,  2, true, 18, 180),
(9, 'Software & Systems Modeling',        '22', '3', 35, 70, 9,  2, true,  7,  90),
(10,'Journal of Cloud Computing',         '9',  '4',  1, 25, 11, 3, true, 25, 320);

-- --------------------------------------------
-- 5. CONFERENCE (pub_id 11..20)
-- --------------------------------------------
INSERT INTO public.conference
(conf_id, conf_name, conf_ref, conf_publ, conf_edit,
 conf_start_date, conf_end_date, conf_start_page, conf_end_page,
 conf_type, conf_serie, conf_organi, conf_particip)
VALUES
(11, 'CDC & Streaming Workshop 2018', 'CDC2018', 'Proceedings', 'ACM',
 '2018-06-10', '2018-06-12', 1, 200, 1, 'Workshop Series', 'Univ. Vigo', 150),
(12, 'Data Migration Summit 2019', 'DMS2019', 'Proceedings', 'IEEE',
 '2019-09-20', '2019-09-22', 1, 250, 1, 'Summit Series', 'TU Berlin', 200),
(13, 'ETL & Integration Conf 2020', 'ETL2020', 'Proceedings', 'Springer',
 '2020-03-05', '2020-03-07', 1, 180, 2, 'Conference Series', 'Univ. Paris', 180),
(14, 'Debezium User Conference 2021', 'DBZ2021', 'Proceedings', 'Red Hat',
 '2021-04-12', '2021-04-14', 1, 160, 2, 'User Conf', 'Red Hat', 220),
(15, 'Kafka and Friends 2022', 'KAF2022', 'Proceedings', 'Confluent',
 '2022-05-18', '2022-05-20', 1, 210, 3, 'Industry Conf', 'Confluent', 300),
(16, 'Distributed Systems Conf 2020', 'DSC2020', 'Proceedings', 'IEEE',
 '2020-10-01', '2020-10-03', 1, 230, 2, 'Conference Series', 'IEEE', 250),
(17, 'Big Data Europe 2021', 'BDE2021', 'Proceedings', 'Springer',
 '2021-11-10', '2021-11-12', 1, 260, 3, 'Conference Series', 'Springer', 280),
(18, 'Cloud Native Databases 2022', 'CND2022', 'Proceedings', 'ACM',
 '2022-07-01', '2022-07-03', 1, 190, 1, 'Workshop', 'ACM', 120),
(19, 'Real-Time Analytics 2023', 'RTA2023', 'Proceedings', 'O''Reilly',
 '2023-09-15', '2023-09-17', 1, 220, 1, 'Industry Conf', 'O''Reilly', 260),
(20, 'Modernization and Migration 2024', 'MM2024', 'Proceedings', 'IEEE',
 '2024-04-20', '2024-04-22', 1, 240, 2, 'Conference Series', 'IEEE', 300);

-- --------------------------------------------
-- 6. BOOK (pub_id 21..25)
-- --------------------------------------------
INSERT INTO public.book
(book_id, book_isbn, book_edit, book_vol, book_editors)
VALUES
(21, '978-3-16-148410-0', 'Springer', '1', 'Luis García; María López'),
(22, '978-1-491-99999-1', 'O''Reilly', '1', 'John Smith'),
(23, '978-0-12-345678-9', 'Elsevier', '1', 'Clara Meier; André Dubois'),
(24, '978-1-234-56789-7', 'ACM', '1', 'Pedro Alonso'),
(25, '978-9-876-54321-0', 'Red Hat', '1', 'Ana González');

-- --------------------------------------------
-- 7. BOOK CHAPTER (pub_id 26..30)
-- --------------------------------------------
INSERT INTO public.book_chapter
(chapt_id, chapt_book_tit, chapt_vol,
 chapt_start_page, chapt_end_page, chapt_edit, chapt_editors)
VALUES
(26, 'public.Systems Modernization Handbook', '1',  1,  25, 'Springer', 'Luis García'),
(27, 'Streaming Architectures with Kafka',     '1', 26,  50, 'O''Reilly', 'John Smith'),
(28, 'Data Integration Patterns',              '1', 51,  80, 'Elsevier', 'Clara Meier'),
(29, 'Microservices and Databases',            '1', 81, 110, 'ACM', 'Pedro Alonso'),
(30, 'Hands-On CDC with Debezium',             '1',111, 140, 'Red Hat', 'Ana González');

-- --------------------------------------------
-- 8. PUBLICATION_WRITER (rese_id ↔ pub_id)
--    Beberapa publikasi multi-penulis, dengan pwrit_ord
-- --------------------------------------------

INSERT INTO public.publication_writer (rese_id, pub_id, pwrit_is_edit, pwrit_ord)
VALUES
-- Journals
(1,  1, false, 1),
(2,  1, false, 2),

(2,  2, false, 1),
(6,  2, false, 2),

(3,  3, false, 1),
(4,  3, false, 2),

(4,  4, false, 1),
(9,  4, false, 2),

(5,  5, false, 1),
(10, 5, false, 2),

(6,  6, false, 1),
(8,  6, false, 2),

(7,  7, false, 1),
(1,  7, false, 2),

(8,  8, false, 1),
(2,  8, false, 2),

(9,  9, false, 1),
(3,  9, false, 2),
(6,  9, false, 3),

(10, 10, false, 1),
(5,  10, false, 2),

-- Conferences
(1,  11, false, 1),
(7,  11, false, 2),

(2,  12, false, 1),
(6,  12, false, 2),
(11, 12, false, 3),

(3,  13, false, 1),
(4,  13, false, 2),

(4,  14, false, 1),
(1,  14, false, 2),

(5,  15, false, 1),
(8,  15, false, 2),

(6,  16, false, 1),
(9,  16, false, 2),

(7,  17, false, 1),
(10, 17, false, 2),

(8,  18, false, 1),
(12, 18, false, 2),

(9,  19, false, 1),
(11, 19, false, 2),

(10, 20, false, 1),
(12, 20, false, 2),

-- Books
(1,  21, true,  1),
(2,  21, true,  2),

(5,  22, true,  1),

(3,  23, true,  1),
(4,  23, true,  2),

(7,  24, true,  1),

(6,  25, true,  1),

-- Book Chapters
(1,  26, false, 1),
(6,  26, false, 2),

(5,  27, false, 1),

(3,  28, false, 1),
(4,  28, false, 2),

(7,  29, false, 1),

(6,  30, false, 1),
(8,  30, false, 2);

COMMIT;
