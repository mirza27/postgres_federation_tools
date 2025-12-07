TRUNCATE TABLE 
    "_keymap_generic",
    "_exec_queue",
    "_batch_log",
    "_polling_store",
    "_need_join",
    "_join_map",
    "_join_map_topic"
RESTART IDENTITY CASCADE;

DROP TABLE IF EXISTS _keymap_generic CASCADE;
DROP TABLE IF EXISTS _exec_queue CASCADE;
DROP TABLE IF EXISTS _batch_log CASCADE;
DROP TABLE IF EXISTS _polling_store CASCADE; 
DROP TABLE IF EXISTS _need_join CASCADE;
DROP TABLE IF EXISTS _join_map CASCADE;
DROP TABLE IF EXISTS _join_map_topic CASCADE;