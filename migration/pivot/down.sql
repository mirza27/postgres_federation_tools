TRUNCATE TABLE 
    "_keymap_generic",
    "_exec_queue",
    "_need_join",
    "_join_map",
    "_join_map_topic"
RESTART IDENTITY CASCADE;

DROP TABLE IF EXISTS _keymap_generic CASCADE;
DROP TABLE IF EXISTS _exec_queue CASCADE;
DROP TABLE IF EXISTS _need_join CASCADE;
DROP TABLE IF EXISTS _join_map CASCADE;
DROP TABLE IF EXISTS _join_map_topic CASCADE;