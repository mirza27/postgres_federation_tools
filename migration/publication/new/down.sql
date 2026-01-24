TRUNCATE TABLE
    public.author,
    public.publicationtype,
    public.publisher,
    public.publication,
    public.paper,
    public.authorship
CASCADE;


DROP TABLE IF EXISTS public.author CASCADE;
DROP TABLE IF EXISTS public.publicationtype CASCADE;
DROP TABLE IF EXISTS public.publisher CASCADE;
DROP TABLE IF EXISTS public.publication CASCADE;
DROP TABLE IF EXISTS public.paper CASCADE;
DROP TABLE IF EXISTS public.researcher CASCADE;
DROP TABLE IF EXISTS public.authorship CASCADE;