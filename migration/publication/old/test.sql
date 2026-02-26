-- scenario 1
INSERT INTO public.researcher 
(rese_id, rese_nif, rese_name, rese_first_surname, rese_sec_surname, rese_birth_day, rese_mail, rese_web, rese_phone, rese_ext, rese_fax, rese_univ, rese_sign, rese_coord, rese_photo, rese_gend, rese_pos, rese_title, rese_fig, rese_short_curr_es, rese_short_curr_en, rese_short_curr_gl) 
VALUES
(99, 'X8273645J', 'Marco', 'Davies', 'Sullivan', '1985-11-24', 'm.davies@research-net.org', 'http://mdavies-lab.io', '+44-7700900123', '552', '913-221-00', 'Tech Institute of Berlin', 'MD', true, NULL, 2, 3, 'PhD', 'Senior Researcher', 'Investigador en sistemas distribuidos', 'Researcher in distributed systems', 'Investigador en sistemas distribuídos');


-- scenario 3
BEGIN TRANSACTION;

INSERT INTO public.country (cntry_id, cntry_name_es, cntry_name_en, cntry_name_gl, cntry_aneca_id) VALUES
('DY', 'Reino Unido', 'United Kingdom', 'Reino Unido', '88');

INSERT INTO public.publication (pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt, pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry) VALUES
(100001, '978-3-16-148410-0', 'Neural Networks: A Modern Approach to Intelligence', 2019, 'London', 'Oxford Press', 42, 'DRAFT_V2', 'Oral', 'MIN-992', 'AI Quarterly Series', 'Volume 4', 'https://archive.org/nn-2019', 'repo://oxford/nn', false, 'DY');

INSERT INTO public.book (book_id, book_isbn, book_edit, book_vol, book_editors) VALUES
(100001, '978-3-16-148410-0', 'Oxford Press', 'V.9', 'Helena M. Thompson');

END TRANSACTION;


-- scenario 4
BEGIN TRANSACTION;

INSERT INTO public.country (cntry_id, cntry_name_es, cntry_name_en, cntry_name_gl, cntry_aneca_id) VALUES
('DY2', 'Canadá', 'Canada', 'Canadá', '442');

INSERT INTO public.publication (pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt, pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry) VALUES
(100002, '978-1-4028-9462-6', 'Scalable Microservices Architecture', 2022, 'Toronto', 'O-Reilly Tech', 15, NULL, 'Poster', NULL, 'Cloud Computing Digest', NULL, 'http://tech-library.ca/ms2022', 'git://internal/ms-arch', true, 'DY2');

INSERT INTO public.conference (conf_id, conf_name, conf_ref, conf_publ, conf_edit, conf_start_date, conf_end_date, conf_start_page, conf_end_page, conf_type, conf_serie, conf_organi, conf_particip) VALUES
(100002, 'Global Summit on Cloud Native 2024', 'GSCN24', 'Tech Proceedings', 'ACM', '2024-05-15', '2024-05-18', 215, 240, 2, 'Cloud Series', 'CNCF', 1200);

END TRANSACTION;


-- scenario 5
BEGIN TRANSACTION;
INSERT INTO public.country (cntry_id, cntry_name_es, cntry_name_en, cntry_name_gl, cntry_aneca_id) VALUES
('DY3', 'Alemania', 'Germany', 'Alemaña', '901');

INSERT INTO public.publication (pub_id, pub_isbn, pub_title, pub_year, pub_loc, pub_edit, pub_subj, pub_draft, pub_prsnt, pub_minute, pub_source, pub_book, pub_link, pub_repo, pub_visib, pub_cntry) VALUES
(100003, '978-0-596-52068-7', 'Advanced Database Internals and Storage', 2025, 'Berlin', 'Springer Verlag', 88, 'FINAL', NULL, 'REF-12', 'Modern Data Engineering', 'B-102', NULL, NULL, true, 'DY3');

INSERT INTO public.journal (jrnl_id, jrnl_journ_title, jrnl_vol, jrnl_numb, jrnl_start_page, jrnl_end_page, jrnl_month, jrnl_scope, jrnl_indexed, jrnl_pos_numb, jrnl_tot_pos) VALUES
(100003, 'Transactions on Knowledge Engineering', '104', '12', 45, 68, 11, 2, true, 3, 120);

END TRANSACTION;
-- scenario 2
INSERT INTO public.publication_writer (rese_id, pub_id, pwrit_is_edit, pwrit_ord) VALUES
(99,  100003, false, 1)



-- remove test data
BEGIN TRANSACTION;

-- 1. Hapus dari tabel anak (child tables) terlebih dahulu
DELETE FROM public.book WHERE book_id = 100001;
DELETE FROM public.conference WHERE conf_id = 100002;
DELETE FROM public.journal WHERE jrnl_id = 100003;

-- 2. Hapus dari tabel publication (setelah referensi di tabel anak hilang)
DELETE FROM public.publication WHERE pub_id IN (100001, 100002, 100003);

-- 3. Hapus dari tabel country (setelah referensi di publication hilang)
DELETE FROM public.country WHERE cntry_id IN ('DY', 'DY2', 'DY3');

-- 4. Hapus dari tabel researcher (tabel mandiri dalam skenario ini)
DELETE FROM public.researcher WHERE rese_id = 99;

COMMIT;


