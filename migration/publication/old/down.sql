TRUNCATE TABLE
    public.publication_writer,
    public.book,
    public.book_chapter,
    public.conference,
    public.journal,
    public.publication,
    public.researcher,
    public.country
CASCADE;


DROP TABLE IF EXISTS public.publication_writer CASCADE;
DROP TABLE IF EXISTS public.book CASCADE;
DROP TABLE IF EXISTS public.book_chapter CASCADE;
DROP TABLE IF EXISTS public.conference CASCADE;
DROP TABLE IF EXISTS public.journal CASCADE;
DROP TABLE IF EXISTS public.publication CASCADE;
DROP TABLE IF EXISTS public.researcher CASCADE;
DROP TABLE IF EXISTS public.country CASCADE;