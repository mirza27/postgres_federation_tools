TRUNCATE TABLE
    public.album,
    public.artist,
    public.customer,
    public.employee,
    public.genre,
    public.invoice,
    public.invoice_line,
    public.media_type,
    public.playlist,
    public.playlist_track,
    public.track
CASCADE;

DROP TABLE IF EXISTS album;
DROP TABLE IF EXISTS artist;
DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS genre;
DROP TABLE IF EXISTS invoice;
DROP TABLE IF EXISTS invoice_line;
DROP TABLE IF EXISTS media_type;
DROP TABLE IF EXISTS playlist;
DROP TABLE IF EXISTS playlist_track;
DROP TABLE IF EXISTS track;