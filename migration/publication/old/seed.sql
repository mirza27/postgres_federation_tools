-- =================================================================
-- SEED DATA UNTUK SCHEMA LAMA (LEGACY)
-- Konsep: Shared Primary Key (Inheritance)
-- =================================================================

BEGIN;

-- Hapus data lama jika diperlukan untuk menghindari duplicate key error saat testing ulang
TRUNCATE TABLE public.publication_writer CASCADE;
TRUNCATE TABLE public.book_chapter CASCADE;
TRUNCATE TABLE public.book CASCADE;
TRUNCATE TABLE public.conference CASCADE;
TRUNCATE TABLE public.journal CASCADE;
TRUNCATE TABLE public.publication CASCADE;
TRUNCATE TABLE public.researcher CASCADE;
TRUNCATE TABLE public.country CASCADE;

-- -----------------------------------------------------------------
-- 1. COUNTRY
-- -----------------------------------------------------------------
INSERT INTO public.country (cntry_id, cntry_name_es, cntry_name_en, cntry_name_gl, cntry_aneca_id) VALUES
('ES', 'España',         'Spain',          'España',         34),
('US', 'Estados Unidos', 'United States',  'Estados Unidos', 1),
('UK', 'Reino Unido',    'United Kingdom', 'Reino Unido',    44),
('FR', 'Francia',        'France',         'Francia',        33),
('DE', 'Alemania',       'Germany',        'Alemaña',        49);

-- -----------------------------------------------------------------
-- 2. RESEARCHER (rese_id 1..12)
-- -----------------------------------------------------------------
INSERT INTO public.researcher (rese_id, rese_nif, rese_name, rese_first_surname, rese_sec_surname, rese_birth_day, rese_mail, rese_web, rese_phone, rese_ext, rese_fax, rese_univ, rese_sign, rese_coord, rese_photo, rese_gend, rese_pos, rese_title, rese_fig, rese_short_curr_es, rese_short_curr_en, rese_short_curr_gl) VALUES
(1,  '12345678A', 'Luis',   'García',    'Pérez',  '1978-03-12', 'lgarcia@example.es',    'http://lgarcia.es', '+34-600111111', '101', NULL, 'Univ. Vigo',     'LG', false, NULL, 1, 1, NULL, NULL, 'Profesor en informática',             'Professor in computer science',   NULL),
(2,  '23456789B', 'María',  'López',     'Santos', '1982-07-25', 'mlopez@example.es',     'http://mlopez.es',  '+34-600222222', '102', NULL, 'Univ. Vigo',     'ML', true,  NULL, 2, 1, NULL, NULL, 'Investigadora en IA',                 'Researcher in AI',                NULL),
(3,  '34567890C', 'André',  'Dubois',    NULL,     '1975-11-05', 'adubois@example.fr',    NULL,                '+33-610333333', NULL,  NULL, 'Univ. Lyon',     'AD', false, NULL, 1, 2, NULL, NULL, 'Experto en redes',                    'Networks expert',                 NULL),
(4,  '45678901D', 'Clara',  'Meier',     NULL,     '1985-02-18', 'cmeier@example.de',     NULL,                '+49-1704444444',NULL,  NULL, 'TU Berlin',      'CM', false, NULL, 2, 2, NULL, NULL, 'Especialista en bases de datos',      'Database specialist',             NULL),
(5,  '56789012E', 'John',   'Smith',     NULL,     '1970-09-09', 'jsmith@example.com',    NULL,                '+1-4155550101', NULL,  NULL, 'UC Berkeley',    'JS', true,  NULL, 1, 1, NULL, NULL, 'Profesor de sistemas distribuidos',   'Distributed systems professor',   NULL),
(6,  '67890123F', 'Ana',    'González',  'Rivas',  '1990-01-30', 'agonzalez@example.es',  NULL,                '+34-600333333', NULL,  NULL, 'Univ. Coruña',   'AG', false, NULL, 2, 3, NULL, NULL, 'Doctora en ingeniería del software',  'Software engineering PhD',        NULL),
(7,  '78901234G', 'Pedro',  'Alonso',    'Lago',   '1988-06-14', 'palonso@example.es',    NULL,                '+34-600444444', NULL,  NULL, 'Univ. Vigo',     'PA', false, NULL, 1, 3, NULL, NULL, 'Investigador postdoctoral',           'Postdoc researcher',              NULL),
(8,  '89012345H', 'Sofia',  'Klein',     NULL,     '1992-04-21', 'sklein@example.de',     NULL,                '+49-1705555555',NULL,  NULL, 'TU Munich',      'SK', false, NULL, 2, 3, NULL, NULL, 'Investigadora en aprendizaje automático','ML researcher',                 NULL),
(9,  '90123456I', 'Daniel', 'Ruiz',      'Martín', '1983-12-02', 'druiz@example.es',      NULL,                '+34-600555555', NULL,  NULL, 'Univ. Madrid',   'DR', false, NULL, 1, 2, NULL, NULL, 'Profesor ayudante',                   'Assistant professor',             NULL),
(10, '01234567J', 'Laura',  'Fernández', 'Gómez',  '1987-05-08', 'lfernandez@example.es', NULL,                '+34-600666666', NULL,  NULL, 'Univ. Sevilla',  'LF', false, NULL, 2, 2, NULL, NULL, 'Investigadora en minería de datos',   'Data mining researcher',          NULL),
(11, '11223344K', 'Miguel', 'Torres',    NULL,     '1979-01-19', 'mtorres@example.es',    NULL,                '+34-600777777', NULL,  NULL, 'Univ. Zaragoza', 'MT', false, NULL, 1, 2, NULL, NULL, 'Experto en sistemas embebidos',       'Embedded systems expert',         NULL),
(12, '22334455L', 'Elena',  'Costa',     NULL,     '1991-08-03', 'ecosta@example.pt',     NULL,                '+351-910000000', NULL,  NULL, 'Univ. Porto',    'EC', false, NULL, 2, 3, NULL, NULL, 'Investigadora en visualización',      'Visualization researcher',        NULL);

-- Reset sequence researcher
SELECT pg_catalog.setval('public.researcher_rese_id_seq', 12, true);

-- -----------------------------------------------------------------
-- 3. PUBLICATION (PARENT TABLE) - IDs 1..30
-- Semua jenis publikasi masuk ke sini dulu untuk mendapatkan ID
-- -----------------------------------------------------------------
INSERT INTO public.publication (pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt, pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry) VALUES
-- 1-10: JOURNAL PAPERS (Detail ada di tabel Journal)
(1, NULL, 'Scalable Data Migration in Heterogeneous Systems', 2015, 'Vigo', 'Springer', 1, NULL, NULL, NULL, 'DB Journal', NULL, NULL, NULL, true, 'ES'),
(2, NULL, 'An Empirical Study on Change Data Capture', 2016, 'Madrid', 'Elsevier', 1, NULL, NULL, NULL, 'Journal of Information Systems', NULL, NULL, NULL, true, 'ES'),
(3, NULL, 'Improving CDC Pipelines with Kafka', 2017, 'Berlin', 'ACM', 2, NULL, NULL, NULL, 'Distributed Systems Journal', NULL, NULL, NULL, true, 'DE'),
(4, NULL, 'Performance Evaluation of Debezium Connectors', 2018, 'Paris', 'IEEE', 2, NULL, NULL, NULL, 'IEEE Data Engineering', NULL, NULL, NULL, true, 'FR'),
(5, NULL, 'A Survey on Database Modernization', 2019, 'London', 'Springer', 3, NULL, NULL, NULL, 'Software Engineering Journal', NULL, NULL, NULL, true, 'UK'),
(6, NULL, 'Streaming ETL for Legacy Systems', 2020, 'San Francisco', 'O''Reilly', 3, NULL, NULL, NULL, 'US Data Science Journal', NULL, NULL, NULL, true, 'US'),
(7, NULL, 'Near Real-Time Data Integration', 2021, 'Munich', 'IEEE', 4, NULL, NULL, NULL, 'Transactions on Big Data', NULL, NULL, NULL, true, 'DE'),
(8, NULL, 'Designing Pivot Databases for Migration', 2022, 'Barcelona', 'Elsevier', 4, NULL, NULL, NULL, 'International Journal of Databases', NULL, NULL, NULL, true, 'ES'),
(9, NULL, 'A DSL for Data Migration: DAMI Framework', 2023, 'Lyon', 'Springer', 5, NULL, NULL, NULL, 'Software & Systems Modeling', NULL, NULL, NULL, true, 'FR'),
(10,NULL, 'Benchmarking CDC Approaches in Microservices', 2024, 'Seattle', 'ACM', 5, NULL, NULL, NULL, 'Journal of Cloud Computing', NULL, NULL, NULL, true, 'US'),

-- 11-20: CONFERENCE PAPERS (Detail ada di tabel Conference)
-- Catatan: Di schema lama, Paper Konferensi dan Event Konferensi sering tercampur/duplikat
(11,NULL, 'Proceedings of the CDC & Streaming Workshop 2018', 2018, 'Vigo', 'ACM', 6, NULL, NULL, NULL, 'CDC Workshop', NULL, NULL, NULL, true, 'ES'),
(12,NULL, 'Proceedings of the Data Migration Summit 2019', 2019, 'Berlin', 'IEEE', 6, NULL, NULL, NULL, 'Data Migration Summit', NULL, NULL, NULL, true, 'DE'),
(13,NULL, 'Proceedings of the ETL and Integration Conf 2020', 2020, 'Paris', 'Springer', 7, NULL, NULL, NULL, 'ETL & Integration Conf', NULL, NULL, NULL, true, 'FR'),
(14,NULL, 'Proceedings of the Debezium User Conference 2021', 2021, 'Madrid', 'Red Hat', 7, NULL, NULL, NULL, 'Debezium Conf', NULL, NULL, NULL, true, 'ES'),
(15,NULL, 'Proceedings of Kafka and Friends 2022', 2022, 'London', 'Confluent', 7, NULL, NULL, NULL, 'Kafka & Friends', NULL, NULL, NULL, true, 'UK'),
(16,NULL, 'Proceedings of Distributed Systems Conf 2020', 2020, 'San Jose', 'IEEE', 8, NULL, NULL, NULL, 'Distributed Systems Conf', NULL, NULL, NULL, true, 'US'),
(17,NULL, 'Proceedings of Big Data Europe 2021', 2021, 'Berlin', 'Springer', 8, NULL, NULL, NULL, 'Big Data Europe', NULL, NULL, NULL, true, 'DE'),
(18,NULL, 'Proceedings of Cloud Native Databases 2022', 2022, 'Amsterdam', 'ACM', 9, NULL, NULL, NULL, 'Cloud Native DB Conf', NULL, NULL, NULL, true, 'FR'),
(19,NULL, 'Proceedings of Real-Time Analytics 2023', 2023, 'New York', 'O''Reilly', 9, NULL, NULL, NULL, 'Real-Time Analytics', NULL, NULL, NULL, true, 'US'),
(20,NULL, 'Proceedings of Modernization and Migration 2024', 2024, 'Lisbon', 'IEEE', 9, NULL, NULL, NULL, 'Modernization & Migration', NULL, NULL, NULL, true, 'ES'),

-- 21-25: BOOKS (Detail ada di tabel Book)
(21,'978-3-16-148410-0', 'Legacy Systems Modernization Handbook', 2014, 'Madrid', 'Springer', 10, NULL, NULL, NULL, 'Handbook Series', NULL, NULL, NULL, true, 'ES'),
(22,'978-1-491-99999-1', 'Streaming Architectures with Kafka', 2017, 'San Francisco', 'O''Reilly', 10, NULL, NULL, NULL, 'Professional Series', NULL, NULL, NULL, true, 'US'),
(23,'978-0-12-345678-9', 'Data Integration Patterns', 2018, 'London', 'Elsevier', 11, NULL, NULL, NULL, 'Patterns Series', NULL, NULL, NULL, true, 'UK'),
(24,'978-1-234-56789-7', 'Microservices and Databases', 2019, 'Berlin', 'ACM', 11, NULL, NULL, NULL, 'Software Architecture Series', NULL, NULL, NULL, true, 'DE'),
(25,'978-9-876-54321-0', 'Hands-On CDC with Debezium', 2021, 'Paris', 'Red Hat', 11, NULL, NULL, NULL, 'Practical Guides', NULL, NULL, NULL, true, 'FR'),

-- 26-30: BOOK CHAPTERS (Detail ada di tabel Book Chapter)
(26,NULL, 'Chapter: Designing Pivot Tables', 2014, 'Madrid', 'Springer', 12, NULL, NULL, NULL, 'Legacy Systems Modernization Handbook', NULL, NULL, NULL, true, 'ES'),
(27,NULL, 'Chapter: Kafka Connect in Practice', 2017, 'San Francisco', 'O''Reilly', 12, NULL, NULL, NULL, 'Streaming Architectures with Kafka', NULL, NULL, NULL, true, 'US'),
(28,NULL, 'Chapter: CDC Design Patterns', 2018, 'London', 'Elsevier', 12, NULL, NULL, NULL, 'Data Integration Patterns', NULL, NULL, NULL, true, 'UK'),
(29,NULL, 'Chapter: Microservices Data Ownership', 2019, 'Berlin', 'ACM', 12, NULL, NULL, NULL, 'Microservices and Databases', NULL, NULL, NULL, true, 'DE'),
(30,NULL, 'Chapter: Debezium in Production', 2021, 'Paris', 'Red Hat', 12, NULL, NULL, NULL, 'Hands-On CDC with Debezium', NULL, NULL, NULL, true, 'FR');

-- Reset sequence publication (Parent)
SELECT pg_catalog.setval('public.publication_pub_id_seq', 30, true);

-- -----------------------------------------------------------------
-- 4. JOURNAL (CHILD) - IDs 1..10
-- PK (jrnl_id) SAMA DENGAN Publication(pub_id)
-- -----------------------------------------------------------------
INSERT INTO public.journal (jrnl_id, jrnl_journ_title, jrnl_vol, jrnl_numb, jrnl_start_page, jrnl_end_page, jrnl_month, jrnl_scope, jrnl_indexed, jrnl_pos_numb, jrnl_tot_pos) VALUES
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

-- -----------------------------------------------------------------
-- 5. CONFERENCE (CHILD) - IDs 11..20
-- PK (conf_id) SAMA DENGAN Publication(pub_id)
-- Masalah di Schema A: Satu entry conference ini juga mewakili paper,
-- sehingga detail konferensi (lokasi/tgl) terduplikasi untuk tiap paper.
-- -----------------------------------------------------------------
INSERT INTO public.conference (conf_id, conf_name, conf_ref, conf_publ, conf_edit, conf_start_date, conf_end_date, conf_start_page, conf_end_page, conf_type, conf_serie, conf_organi, conf_particip) VALUES
(11, 'CDC & Streaming Workshop 2018', 'CDC2018', 'Proceedings', 'ACM',       '2018-06-10', '2018-06-12', 1, 200, 1, 'Workshop Series', 'Univ. Vigo', 150),
(12, 'Data Migration Summit 2019',    'DMS2019', 'Proceedings', 'IEEE',      '2019-09-20', '2019-09-22', 1, 250, 1, 'Summit Series',   'TU Berlin',  200),
(13, 'ETL & Integration Conf 2020',   'ETL2020', 'Proceedings', 'Springer',  '2020-03-05', '2020-03-07', 1, 180, 2, 'Conference Series', 'Univ. Paris', 180),
(14, 'Debezium User Conference 2021', 'DBZ2021', 'Proceedings', 'Red Hat',   '2021-04-12', '2021-04-14', 1, 160, 2, 'User Conf',       'Red Hat',    220),
(15, 'Kafka and Friends 2022',        'KAF2022', 'Proceedings', 'Confluent', '2022-05-18', '2022-05-20', 1, 210, 3, 'Industry Conf',   'Confluent',  300),
(16, 'Distributed Systems Conf 2020', 'DSC2020', 'Proceedings', 'IEEE',      '2020-10-01', '2020-10-03', 1, 230, 2, 'Conference Series', 'IEEE',       250),
(17, 'Big Data Europe 2021',          'BDE2021', 'Proceedings', 'Springer',  '2021-11-10', '2021-11-12', 1, 260, 3, 'Conference Series', 'Springer',   280),
(18, 'Cloud Native Databases 2022',   'CND2022', 'Proceedings', 'ACM',       '2022-07-01', '2022-07-03', 1, 190, 1, 'Workshop',        'ACM',        120),
(19, 'Real-Time Analytics 2023',      'RTA2023', 'Proceedings', 'O''Reilly', '2023-09-15', '2023-09-17', 1, 220, 1, 'Industry Conf',   'O''Reilly',  260),
(20, 'Modernization and Migration 2024','MM2024', 'Proceedings', 'IEEE',      '2024-04-20', '2024-04-22', 1, 240, 2, 'Conference Series', 'IEEE',       300);

-- -----------------------------------------------------------------
-- 6. BOOK (CHILD) - IDs 21..25
-- PK (book_id) SAMA DENGAN Publication(pub_id)
-- -----------------------------------------------------------------
INSERT INTO public.book (book_id, book_isbn, book_edit, book_vol, book_editors) VALUES
(21, '978-3-16-148410-0', 'Springer', '1', 'Luis García; María López'),
(22, '978-1-491-99999-1', 'O''Reilly', '1', 'John Smith'),
(23, '978-0-12-345678-9', 'Elsevier', '1', 'Clara Meier; André Dubois'),
(24, '978-1-234-56789-7', 'ACM',      '1', 'Pedro Alonso'),
(25, '978-9-876-54321-0', 'Red Hat',  '1', 'Ana González');

-- -----------------------------------------------------------------
-- 7. BOOK CHAPTER (CHILD) - IDs 26..30
-- PK (chapt_id) SAMA DENGAN Publication(pub_id)
-- -----------------------------------------------------------------
INSERT INTO public.book_chapter (chapt_id, chapt_book_tit, chapt_vol, chapt_start_page, chapt_end_page, chapt_edit, chapt_editors) VALUES
(26, 'Legacy Systems Modernization Handbook', '1',  1,  25, 'Springer', 'Luis García'),
(27, 'Streaming Architectures with Kafka',     '1', 26,  50, 'O''Reilly', 'John Smith'),
(28, 'Data Integration Patterns',              '1', 51,  80, 'Elsevier', 'Clara Meier'),
(29, 'Microservices and Databases',            '1', 81, 110, 'ACM',      'Pedro Alonso'),
(30, 'Hands-On CDC with Debezium',             '1',111, 140, 'Red Hat',  'Ana González');

-- -----------------------------------------------------------------
-- 8. PUBLICATION_WRITER (rese_id <-> pub_id)
-- Menghubungkan Peneliti dengan Publikasi (Parent ID)
-- -----------------------------------------------------------------
INSERT INTO public.publication_writer (rese_id, pub_id, pwrit_is_edit, pwrit_ord) VALUES
-- Journals (IDs 1-10)
(1,  1, false, 1), (2,  1, false, 2),
(2,  2, false, 1), (6,  2, false, 2),
(3,  3, false, 1), (4,  3, false, 2),
(4,  4, false, 1), (9,  4, false, 2),
(5,  5, false, 1), (10, 5, false, 2),
(6,  6, false, 1), (8,  6, false, 2),
(7,  7, false, 1), (1,  7, false, 2),
(8,  8, false, 1), (2,  8, false, 2),
(9,  9, false, 1), (3,  9, false, 2), (6, 9, false, 3),
(10, 10, false, 1), (5,  10, false, 2),

-- Conferences (IDs 11-20)
(1,  11, false, 1), (7,  11, false, 2),
(2,  12, false, 1), (6,  12, false, 2), (11, 12, false, 3),
(3,  13, false, 1), (4,  13, false, 2),
(4,  14, false, 1), (1,  14, false, 2),
(5,  15, false, 1), (8,  15, false, 2),
(6,  16, false, 1), (9,  16, false, 2),
(7,  17, false, 1), (10, 17, false, 2),
(8,  18, false, 1), (12, 18, false, 2),
(9,  19, false, 1), (11, 19, false, 2),
(10, 20, false, 1), (12, 20, false, 2),

-- Books (IDs 21-25) - usually editors (pwrit_is_edit = true)
(1,  21, true,  1), (2,  21, true,  2),
(5,  22, true,  1),
(3,  23, true,  1), (4,  23, true,  2),
(7,  24, true,  1),
(6,  25, true,  1),

-- Book Chapters (IDs 26-30)
(1,  26, false, 1), (6,  26, false, 2),
(5,  27, false, 1),
(3,  28, false, 1), (4,  28, false, 2),
(7,  29, false, 1),
(6,  30, false, 1), (8,  30, false, 2);

COMMIT;


-- =================================================================
-- TAMBAHAN DATA BARU (BATCH 2)
-- Melanjutkan sequence dari ID 30
-- =================================================================

BEGIN;

-- -----------------------------------------------------------------
-- 1. INSERT KE PARENT TABLE (PUBLICATION)
-- ID 31-60
-- -----------------------------------------------------------------
INSERT INTO public.publication (pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt, pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry) VALUES

-- --- JURNAL TAMBAHAN (ID 31-40) ---
(31, NULL, 'Advances in Quantum Cryptography', 2023, 'Zurich', 'Springer', 20, NULL, NULL, NULL, 'Journal of Cryptology', NULL, NULL, NULL, true, 'DE'),
(32, NULL, 'Deep Learning for Medical Imaging', 2024, 'London', 'Elsevier', 20, NULL, NULL, NULL, 'Artificial Intelligence in Medicine', NULL, NULL, NULL, true, 'UK'),
(33, NULL, 'Blockchain Scalability Solutions', 2022, 'New York', 'ACM', 21, NULL, NULL, NULL, 'ACM Transactions on Blockchain', NULL, NULL, NULL, true, 'US'),
(34, NULL, 'Edge Computing in 5G Networks', 2023, 'Paris', 'IEEE', 21, NULL, NULL, NULL, 'IEEE Internet of Things Journal', NULL, NULL, NULL, true, 'FR'),
(35, NULL, 'Ethical Implications of AI', 2024, 'Oxford', 'Oxford Univ Press', 22, NULL, NULL, NULL, 'Journal of AI Ethics', NULL, NULL, NULL, true, 'UK'),
(36, NULL, 'Serverless Architecture Patterns', 2023, 'San Francisco', 'O''Reilly', 22, NULL, NULL, NULL, 'Cloud Computing Reports', NULL, NULL, NULL, true, 'US'),
(37, NULL, 'Cybersecurity in Smart Grids', 2022, 'Berlin', 'Springer', 23, NULL, NULL, NULL, 'Energy Systems Security', NULL, NULL, NULL, true, 'DE'),
(38, NULL, 'Natural Language Processing for Low-Resource Languages', 2024, 'Barcelona', 'Elsevier', 23, NULL, NULL, NULL, 'Computational Linguistics', NULL, NULL, NULL, true, 'ES'),
(39, NULL, 'DevOps Metrics and KPIs', 2023, 'Madrid', 'IEEE', 24, NULL, NULL, NULL, 'IEEE Software', NULL, NULL, NULL, true, 'ES'),
(40, NULL, 'Green Computing Initiatives', 2022, 'Munich', 'Springer', 24, NULL, NULL, NULL, 'Sustainable Computing', NULL, NULL, NULL, true, 'DE'),

-- --- KONFERENSI TAMBAHAN (ID 41-50) ---
(41, NULL, 'Proc. of Int. Conf. on Computer Vision 2023', 2023, 'Paris', 'IEEE', 25, NULL, NULL, NULL, 'ICCV 2023', NULL, NULL, NULL, true, 'FR'),
(42, NULL, 'Proc. of Neural Information Processing Systems 2024', 2024, 'New Orleans', 'NeurIPS', 25, NULL, NULL, NULL, 'NeurIPS 2024', NULL, NULL, NULL, true, 'US'),
(43, NULL, 'Proc. of Int. Conf. on Software Engineering 2023', 2023, 'Melbourne', 'IEEE', 26, NULL, NULL, NULL, 'ICSE 2023', NULL, NULL, NULL, true, 'UK'),
(44, NULL, 'Proc. of Kubernetes Community Days 2024', 2024, 'Amsterdam', 'CNCF', 26, NULL, NULL, NULL, 'KCD 2024', NULL, NULL, NULL, true, 'DE'),
(45, NULL, 'Proc. of Black Hat Europe 2023', 2023, 'London', 'Black Hat', 27, NULL, NULL, NULL, 'Black Hat Eu', NULL, NULL, NULL, true, 'UK'),
(46, NULL, 'Proc. of PyData Global 2024', 2024, 'Online', 'NumFOCUS', 27, NULL, NULL, NULL, 'PyData 2024', NULL, NULL, NULL, true, 'US'),
(47, NULL, 'Proc. of European Conf. on Info Systems 2023', 2023, 'Kristiansand', 'AIS', 28, NULL, NULL, NULL, 'ECIS 2023', NULL, NULL, NULL, true, 'DE'),
(48, NULL, 'Proc. of DefCon 32', 2024, 'Las Vegas', 'DefCon', 28, NULL, NULL, NULL, 'DefCon 2024', NULL, NULL, NULL, true, 'US'),
(49, NULL, 'Proc. of Web Summit 2023', 2023, 'Lisbon', 'Web Summit', 29, NULL, NULL, NULL, 'Web Summit 23', NULL, NULL, NULL, true, 'ES'),
(50, NULL, 'Proc. of Agile Alliance 2024', 2024, 'Washington', 'Agile Alliance', 29, NULL, NULL, NULL, 'Agile 2024', NULL, NULL, NULL, true, 'US'),

-- --- BUKU TAMBAHAN (ID 51-60) ---
(51, '978-0-13-235088-4', 'Clean Code: A Handbook of Agile Software Craftsmanship', 2008, 'Boston', 'Prentice Hall', 30, NULL, NULL, NULL, 'Robert C. Martin Series', NULL, NULL, NULL, true, 'US'),
(52, '978-0-321-12521-7', 'Domain-Driven Design: Tackling Complexity in the Heart of Software', 2003, 'Boston', 'Addison-Wesley', 30, NULL, NULL, NULL, 'Evans Series', NULL, NULL, NULL, true, 'US'),
(53, '978-1-492-04034-7', 'Designing Data-Intensive Applications', 2017, 'Sebastopol', 'O''Reilly', 31, NULL, NULL, NULL, 'Big Data Series', NULL, NULL, NULL, true, 'US'),
(54, '978-0-201-63361-0', 'Design Patterns: Elements of Reusable Object-Oriented Software', 1994, 'Boston', 'Addison-Wesley', 31, NULL, NULL, NULL, 'GoF Series', NULL, NULL, NULL, true, 'US'),
(55, '978-1-617-29454-9', 'Microservices Patterns', 2018, 'Shelter Island', 'Manning', 32, NULL, NULL, NULL, 'Manning Publications', NULL, NULL, NULL, true, 'US'),
(56, '978-0-13-449416-6', 'Refactoring: Improving the Design of Existing Code', 2018, 'Boston', 'Addison-Wesley', 32, NULL, NULL, NULL, 'Fowler Series', NULL, NULL, NULL, true, 'US'),
(57, '978-0-321-35668-0', 'Effective Java', 2017, 'Boston', 'Addison-Wesley', 33, NULL, NULL, NULL, 'Java Series', NULL, NULL, NULL, true, 'US'),
(58, '978-1-118-00818-8', 'The Lean Startup', 2011, 'New York', 'Crown Business', 33, NULL, NULL, NULL, 'Business Series', NULL, NULL, NULL, true, 'US'),
(59, '978-1-593-27950-9', 'Automate the Boring Stuff with Python', 2019, 'San Francisco', 'No Starch Press', 34, NULL, NULL, NULL, 'Python Series', NULL, NULL, NULL, true, 'US'),
(60, '978-0-13-409266-9', 'The Pragmatic Programmer: Your Journey to Mastery', 2019, 'Boston', 'Addison-Wesley', 34, NULL, NULL, NULL, 'Pragmatic Series', NULL, NULL, NULL, true, 'US');

-- Update Sequence
SELECT pg_catalog.setval('public.publication_pub_id_seq', 60, true);


-- -----------------------------------------------------------------
-- 2. INSERT KE CHILD TABLE: JOURNAL
-- ID 31-40 (Harus match dengan pub_id)
-- -----------------------------------------------------------------
INSERT INTO public.journal (jrnl_id, jrnl_journ_title, jrnl_vol, jrnl_numb, jrnl_start_page, jrnl_end_page, jrnl_month, jrnl_scope, jrnl_indexed, jrnl_pos_numb, jrnl_tot_pos) VALUES
(31, 'Journal of Cryptology',                '36', '2',  45,  80, 2,  1, true,  10, 100),
(32, 'Artificial Intelligence in Medicine',  '55', '4', 112, 145, 5,  1, true,  12, 120),
(33, 'ACM Transactions on Blockchain',       '5',  '1',  10,  35, 1,  2, true,   8,  60),
(34, 'IEEE Internet of Things Journal',      '10', '6', 200, 230, 6,  1, true,  20, 300),
(35, 'Journal of AI Ethics',                 '3',  '2',  15,  40, 4,  3, true,   5,  40),
(36, 'Cloud Computing Reports',              '12', '3',  50,  75, 7,  2, false, 15, 150),
(37, 'Energy Systems Security',              '8',  '1',   1,  25, 3,  2, true,   6,  80),
(38, 'Computational Linguistics',            '49', '4', 300, 340, 12, 1, true,  18, 180),
(39, 'IEEE Software',                        '41', '5',  20,  55, 9,  1, true,  10, 100),
(40, 'Sustainable Computing',                '15', '2',  60,  90, 8,  2, true,  14, 140);


-- -----------------------------------------------------------------
-- 3. INSERT KE CHILD TABLE: CONFERENCE
-- ID 41-50 (Harus match dengan pub_id)
-- -----------------------------------------------------------------
INSERT INTO public.conference (conf_id, conf_name, conf_ref, conf_publ, conf_edit, conf_start_date, conf_end_date, conf_start_page, conf_end_page, conf_type, conf_serie, conf_organi, conf_particip) VALUES
(41, 'Int. Conf. on Computer Vision 2023',         'ICCV23',  'Proceedings', 'IEEE',           '2023-10-02', '2023-10-06', 1, 300, 1, 'ICCV Series', 'IEEE', 500),
(42, 'Neural Info Processing Systems 2024',         'NIPS24',  'Proceedings', 'NeurIPS',        '2024-12-09', '2024-12-14', 1, 400, 1, 'NeurIPS Series', 'NeurIPS Foundation', 800),
(43, 'Int. Conf. on Software Engineering 2023',     'ICSE23',  'Proceedings', 'IEEE',           '2023-05-14', '2023-05-20', 1, 250, 1, 'ICSE Series', 'IEEE CS', 600),
(44, 'Kubernetes Community Days Amsterdam 2024',    'KCD24',   'Proceedings', 'CNCF',           '2024-02-23', '2024-02-24', 1, 100, 2, 'KCD Series', 'CNCF', 350),
(45, 'Black Hat Europe 2023',                       'BHEU23',  'Proceedings', 'Black Hat',      '2023-12-04', '2023-12-07', 1, 150, 3, 'Security Series', 'Informa', 400),
(46, 'PyData Global 2024',                          'PYD24',   'Proceedings', 'NumFOCUS',       '2024-12-01', '2024-12-03', 1, 120, 2, 'PyData Series', 'NumFOCUS', 1000),
(47, 'European Conf. on Info Systems 2023',         'ECIS23',  'Proceedings', 'AIS',            '2023-06-11', '2023-06-16', 1, 220, 1, 'ECIS Series', 'AIS', 300),
(48, 'DefCon 32',                                   'DC32',    'Proceedings', 'DefCon',         '2024-08-08', '2024-08-11', 1, 200, 3, 'Hacker Series', 'DefCon', 2000),
(49, 'Web Summit Lisbon 2023',                      'WS23',    'Proceedings', 'Web Summit',     '2023-11-13', '2023-11-16', 1, 180, 3, 'Tech Series', 'Web Summit', 5000),
(50, 'Agile 2024',                                  'AG24',    'Proceedings', 'Agile Alliance', '2024-07-22', '2024-07-26', 1, 160, 2, 'Agile Series', 'Agile Alliance', 450);


-- -----------------------------------------------------------------
-- 4. INSERT KE CHILD TABLE: BOOK
-- ID 51-60 (Harus match dengan pub_id)
-- -----------------------------------------------------------------
INSERT INTO public.book (book_id, book_isbn, book_edit, book_vol, book_editors) VALUES
(51, '978-0-13-235088-4', 'Prentice Hall',  '1', 'Robert C. Martin'),
(52, '978-0-321-12521-7', 'Addison-Wesley', '1', 'Eric Evans'),
(53, '978-1-492-04034-7', 'O''Reilly',      '1', 'Martin Kleppmann'),
(54, '978-0-201-63361-0', 'Addison-Wesley', '1', 'Erich Gamma; Richard Helm'),
(55, '978-1-617-29454-9', 'Manning',        '1', 'Chris Richardson'),
(56, '978-0-13-449416-6', 'Addison-Wesley', '2', 'Martin Fowler'),
(57, '978-0-321-35668-0', 'Addison-Wesley', '3', 'Joshua Bloch'),
(58, '978-1-118-00818-8', 'Crown Business', '1', 'Eric Ries'),
(59, '978-1-593-27950-9', 'No Starch Press','1', 'Al Sweigart'),
(60, '978-0-13-409266-9', 'Addison-Wesley', '2', 'Andrew Hunt; David Thomas');


-- -----------------------------------------------------------------
-- 5. INSERT KE PUBLICATION_WRITER (Link Penulis ke Publikasi Baru)
-- Menggunakan researcher ID 1-12 secara acak
-- -----------------------------------------------------------------
INSERT INTO public.publication_writer (rese_id, pub_id, pwrit_is_edit, pwrit_ord) VALUES
-- Jurnal Baru (31-40)
(3, 31, false, 1), (5, 31, false, 2),
(6, 32, false, 1), (2, 32, false, 2),
(1, 33, false, 1),
(11, 34, false, 1), (12, 34, false, 2),
(8, 35, false, 1),
(5, 36, false, 1), (4, 36, false, 2),
(3, 37, false, 1),
(10, 38, false, 1),
(7, 39, false, 1), (9, 39, false, 2),
(4, 40, false, 1), (2, 40, false, 2),

-- Konferensi Baru (41-50)
(6, 41, false, 1), (8, 41, false, 2), (2, 41, false, 3),
(2, 42, false, 1), (8, 42, false, 2),
(1, 43, false, 1), (7, 43, false, 2),
(5, 44, false, 1),
(3, 45, false, 1), (11, 45, false, 2),
(10, 46, false, 1), (12, 46, false, 2),
(4, 47, false, 1), (9, 47, false, 2),
(11, 48, false, 1), (3, 48, false, 2),
(5, 49, false, 1),
(1, 50, false, 1), (7, 50, false, 2),

-- Buku Baru (51-60)
-- Penulis buku biasanya editor atau penulis utama, set is_edit=true
(1, 51, true, 1),
(4, 52, true, 1),
(5, 53, true, 1),
(7, 54, true, 1), (9, 54, true, 2),
(11, 55, true, 1),
(7, 56, true, 1),
(5, 57, true, 1),
(2, 58, true, 1),
(12, 59, true, 1),
(1, 60, true, 1), (6, 60, true, 2);

COMMIT;