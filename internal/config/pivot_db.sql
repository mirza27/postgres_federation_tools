create extension if not exists "pgcrypto";

create table if not exists _keymap_generic (
  keymap_id bigserial primary key,
  map_name text not null,
  src_key text not null,
  tgt_key text,
  src_column text,
  tgt_column text,
  src_table text,
  tgt_table text,
  queue_id uuid,
  tgt_key_status text not null default 'requested',
  requested_at timestamptz not null default now(),
  fulfilled_at timestamptz,
  attempts integer not null default 0,
  last_error text,
  unique (map_name, src_key, src_table, tgt_table)
);

create table if not exists _exec_queue (
  id bigserial primary key,
  queue_id uuid not null default gen_random_uuid(),
  entity text not null,
  op text not null,
  sql_text text not null,
  sql_args jsonb,
  returning_cols text[],
  keymap_payload jsonb,
  is_split boolean not null default false,
  split_name text,
  status text not null default 'pending',
  need_keymap boolean not null default false,
  keymap_id bigint references _keymap_generic(keymap_id),
  retry_count integer not null default 0,
  last_error text,
  locked_at timestamptz,
  locked_by text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (queue_id)
);

create table if not exists _exec_split (
  exec_split_id bigserial primary key,
  queue_id uuid not null references _exec_queue(queue_id) on delete cascade,
  sql_text text not null,
  sql_args jsonb,
  returning_cols text[],
  keymap_payload jsonb,
  status text not null default 'pending',
  last_error text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);


create table if not exists _need_join (
  id bigserial primary key,
  queue_id uuid not null default gen_random_uuid(),
  entity text not null,
  op text not null,
  join_key text not null,
  join_topic text not null,
  join_payload jsonb,
  join_fields jsonb,
  status text not null default 'pending',
  last_error text,
  locked_at timestamptz,
  locked_by text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (queue_id)
);

-- join map untuk menyimpan potongan payload per topic sampai lengkap
create table if not exists _join_map (
  join_map_id bigserial primary key,
  entity text not null,
  join_key text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (entity, join_key)
);

create table if not exists _join_map_topic (
  join_map_topic_id bigserial primary key,
  join_map_id bigint not null references _join_map(join_map_id) on delete cascade,
  topic text not null,
  source_key text not null default '',
  join_fields jsonb,
  payload jsonb,
  updated_at timestamptz not null default now(),
  unique (join_map_id, topic, source_key)
);
