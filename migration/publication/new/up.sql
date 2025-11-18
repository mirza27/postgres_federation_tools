CREATE TABLE public.author (
	code int8 NOT NULL,
	name varchar NULL,
	lastname1 varchar NULL,
	lastname2 varchar NULL,
	CONSTRAINT author_pk PRIMARY KEY ("code")
);

CREATE TABLE public.publicationtype (
	name varchar NOT NULL,
	CONSTRAINT publicationtype_pk PRIMARY KEY (name)
);

CREATE TABLE public.publisher (
	name varchar NOT null,
	CONSTRAINT publisher_pk PRIMARY KEY (name)
);

CREATE TABLE public.publication (
	code int8 NOT NULL,
	title varchar NULL,
	place varchar NULL,
	country varchar NULL,
	year int8 NULL,
	publicationtype_name varchar NULL,
	publisher_name varchar NULL,
	CONSTRAINT publication_pk PRIMARY KEY (code),
	CONSTRAINT publication_fk FOREIGN KEY (publicationtype_name) REFERENCES public.publicationtype(name),
	CONSTRAINT publication_fk_1 FOREIGN KEY (publisher_name) REFERENCES public.publisher(name)
);

CREATE TABLE public.paper (
	title varchar NULL,
	doi varchar NOT NULL,
	fulltext bytea NULL,
	startpage int8 NULL,
	endpage int8 NULL,
	publication_code int8 NULL,
	CONSTRAINT doi PRIMARY KEY (doi),
	CONSTRAINT paper_fk FOREIGN KEY (publication_code) REFERENCES public.publication(code)
);

CREATE TABLE public.authorship (
	code int8 NOT NULL,
	orderofauthor int8 NULL,
	author_code int8 NULL,
	paper_doi varchar NULL,
	CONSTRAINT authorship_pk PRIMARY KEY (code),
	CONSTRAINT authorship_fk FOREIGN KEY (author_code) REFERENCES public.author(code),
	CONSTRAINT authorship__fk_2 FOREIGN KEY (paper_doi) REFERENCES public.paper(doi)
);