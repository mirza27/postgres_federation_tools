package pivot

// DefaultSchemaSQL returns the SQL that ensures the required pivot tables exist.
func DefaultSchemaSQL() string {
	return `
create table if not exists _keymap_generic (
  map_name text not null,
  src_key text not null,
  tgt_key uuid not null,
  primary key (map_name, src_key),
  unique (map_name, tgt_key)
);
create table if not exists _exec_queue (
  id bigserial primary key,
  entity text not null,
  op text not null,
  sql_text text not null,
  sql_args jsonb,
  status text not null default 'pending',
  error text,
  created_at timestamptz default now()
);
create table if not exists _batch_log (
  id bigserial primary key,
  entity text not null,
  op text not null,
  key_values jsonb,
  payload jsonb,
  status text not null,
  error text,
  processed_at timestamptz default now(),
  batch_id text
);`
}
