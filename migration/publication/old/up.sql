-- DROP SEQUENCE public.publication_pub_id_seq;

CREATE SEQUENCE public.publication_pub_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;
-- DROP SEQUENCE public.researcher_rese_id_seq;

CREATE SEQUENCE public.researcher_rese_id_seq
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	START 1
	CACHE 1
	NO CYCLE;-- public.country definition

-- Drop table

-- DROP TABLE public.country;

CREATE TABLE public.country (
	cntry_id bpchar(2) NOT NULL,
	cntry_name_es varchar(150) NULL,
	cntry_name_en varchar(150) NULL,
	cntry_name_gl varchar(150) NULL,
	cntry_aneca_id int4 NULL,
	CONSTRAINT pk_cntry PRIMARY KEY (cntry_id)
);


-- public.researcher definition

-- Drop table

-- DROP TABLE public.researcher;

CREATE TABLE public.researcher (
	rese_id serial4 NOT NULL,
	rese_nif varchar(10) NULL,
	rese_name varchar(30) NULL,
	rese_first_surname varchar(50) NULL,
	rese_sec_surname varchar(50) NULL,
	rese_birth_day date NULL,
	rese_mail varchar(50) NULL,
	rese_web varchar(90) NULL,
	rese_phone varchar(20) NULL,
	rese_ext varchar(10) NULL,
	rese_fax varchar(20) NULL,
	rese_univ varchar(50) NULL,
	rese_sign varchar(40) NULL,
	rese_coord bool NOT NULL DEFAULT false,
	rese_photo varchar(100) NULL,
	rese_gend int4 NOT NULL,
	rese_pos int4 NOT NULL,
	rese_title int4 NULL,
	rese_fig int4 NULL,
	rese_short_curr_es text NULL,
	rese_short_curr_en text NULL,
	rese_short_curr_gl text NULL,
	CONSTRAINT pk_rese PRIMARY KEY (rese_id)
);


-- public."publication" definition

-- Drop table

-- DROP TABLE public."publication";

CREATE TABLE public."publication" (
	pub_id serial4 NOT NULL,
	pub_isbn varchar(30) NULL,
	pub_title varchar(250) NULL,
	pub_year numeric(4) NULL,
	pub_loc varchar(50) NULL,
	pub_edit varchar(75) NULL,
	pub_subj int4 NULL,
	pub_draft varchar(150) NULL,
	pub_prsnt varchar(150) NULL,
	pub_minute varchar(150) NULL,
	pub_source varchar(150) NULL,
	pub_book varchar(150) NULL,
	pub_link varchar(250) NULL,
	pub_repo varchar(250) NULL,
	pub_visib bool NOT NULL DEFAULT false,
	pub_cntry bpchar(2) NULL,
	CONSTRAINT pk_pub PRIMARY KEY (pub_id),
	CONSTRAINT publication_country_fk FOREIGN KEY (pub_cntry) REFERENCES public.country(cntry_id)
);


-- public.publication_writer definition

-- Drop table

-- DROP TABLE public.publication_writer;

CREATE TABLE public.publication_writer (
	rese_id int4 NOT NULL,
	pub_id int4 NOT NULL,
	pwrit_is_edit bool NULL,
	pwrit_ord int4 NULL,
	CONSTRAINT pk_pwrit PRIMARY KEY (rese_id, pub_id),
	CONSTRAINT fk_pwrit_pub FOREIGN KEY (pub_id) REFERENCES public."publication"(pub_id) ON DELETE CASCADE,
	CONSTRAINT fk_pwrit_rese FOREIGN KEY (rese_id) REFERENCES public.researcher(rese_id),
	CONSTRAINT publication_writer_publication_fk FOREIGN KEY (pub_id) REFERENCES public."publication"(pub_id),
	CONSTRAINT publication_writer_researcher_fk FOREIGN KEY (rese_id) REFERENCES public.researcher(rese_id)
);


-- public.book definition

-- Drop table

-- DROP TABLE public.book;

CREATE TABLE public.book (
	book_id int4 NOT NULL,
	book_isbn varchar(30) NULL,
	book_edit varchar(75) NULL,
	book_vol varchar(75) NULL,
	book_editors varchar(500) NULL,
	CONSTRAINT pk_book PRIMARY KEY (book_id),
	CONSTRAINT fk_book_pub FOREIGN KEY (book_id) REFERENCES public."publication"(pub_id)
);


-- public.book_chapter definition

-- Drop table

-- DROP TABLE public.book_chapter;

CREATE TABLE public.book_chapter (
	chapt_id int4 NOT NULL,
	chapt_book_tit varchar(100) NULL,
	chapt_vol varchar(75) NULL,
	chapt_start_page int2 NULL,
	chapt_end_page int2 NULL,
	chapt_edit varchar(75) NULL,
	chapt_editors varchar(500) NULL,
	CONSTRAINT pk_bchap PRIMARY KEY (chapt_id),
	CONSTRAINT fk_bchap_pub FOREIGN KEY (chapt_id) REFERENCES public."publication"(pub_id)
);


-- public.conference definition

-- Drop table

-- DROP TABLE public.conference;

CREATE TABLE public.conference (
	conf_id int4 NOT NULL,
	conf_name varchar(650) NULL,
	conf_ref varchar(100) NULL,
	conf_publ varchar(600) NULL,
	conf_edit varchar(600) NULL,
	conf_start_date date NULL,
	conf_end_date date NULL,
	conf_start_page int4 NULL,
	conf_end_page int4 NULL,
	conf_type int4 NULL,
	conf_serie varchar(150) NULL,
	conf_organi varchar(150) NULL,
	conf_particip int4 NOT NULL,
	CONSTRAINT pk_confe PRIMARY KEY (conf_id),
	CONSTRAINT conference_publication_fk FOREIGN KEY (conf_id) REFERENCES public."publication"(pub_id)
);


-- public.journal definition

-- Drop table

-- DROP TABLE public.journal;

CREATE TABLE public.journal (
	jrnl_id int4 NOT NULL,
	jrnl_journ_title varchar(100) NULL,
	jrnl_vol varchar(50) NULL,
	jrnl_numb varchar(50) NULL,
	jrnl_start_page int2 NULL,
	jrnl_end_page int2 NULL,
	jrnl_month int4 NULL,
	jrnl_scope int4 NULL,
	jrnl_indexed bool NOT NULL DEFAULT false,
	jrnl_pos_numb int4 NULL,
	jrnl_tot_pos int4 NULL,
	CONSTRAINT pk_jrnl PRIMARY KEY (jrnl_id),
	CONSTRAINT fk_jrnl_pub FOREIGN KEY (jrnl_id) REFERENCES public."publication"(pub_id)
);



CREATE OR REPLACE AGGREGATE public.textcat_all(pg_catalog.text) (
	SFUNC = textcat,
	STYPE = text
);