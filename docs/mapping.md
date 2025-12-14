# Mapping Config Guide

Dokumen ini merangkum struktur dan variasi konfigurasi JSON yang dipakai untuk mendefinisikan mapping entitas. File mapping ditempatkan di `mappingPath` (lihat `config.DefaultConfigPath` untuk default.json yang memuat bagian global seperti `sources`, `target`, `engine`).

## Struktur Root (default.json)
- `version` : string versi mapping (opsional).
- `sources[]` : daftar sumber Debezium (name/type/default_schema).
- `target` : info DB target (`name`, `type`, `schema`).
- `engine` : setelan runtime (dialect, batch.maxRows/maxIntervalMs, logging, pivotDb, expr.dialect).
- `entities[]` : entitas bisa langsung ditaruh di root default.json, namun biasanya diletakkan di file terpisah per entitas pada `mappingPath`.

## Entity
Contoh kerangka (tanpa field `kind`, karena tidak dipakai di pipeline):
```json
{
  "entity": "order",
  "sources": [...],
  "target_table": "public.orders",
  "key": {...},
  "columns": {...},
  "fact_condition": {"column": "status", "op": "equal", "value": "PAID"},
  "split_table": [...],
  "routing": {...}
}
```

### Sources dan Join
- `alias` : nama pendek dipakai di `key.source`, `columns.from`, dan join config.
- `from` : nama tabel sumber (dipakai juga untuk derivasi topic jika `topic` kosong).
- `topic` : override nama topic Debezium.
- `join` (opsional, untuk dimensi/fragment):
  - `fact_column` : kolom di fact untuk join.
  - `dim_column`  : kolom di dimensi/fragment untuk join.
  - catatan: alias sumber pertama dalam `sources` selalu dianggap fact; join selalu mengacu ke fact tersebut.

Jika entitas memiliki lebih dari satu source, worker parser akan menyimpan payload fact ke `_need_join` dan fragment ke `_join_map_topic`, lalu joiner menunggu semua topic lengkap sebelum menulis eksekusi.

### Key
- `strategy` : `natural` | `shared_key` (alias sebelumnya `surrogate` tidak lagi didukung).
- `source[]` : daftar `<alias>.<kolom>` yang dipakai untuk menghasilkan `$key` (natural atau lookup untuk keymap).
- `join_key` : opsional, dipakai jika kunci join berbeda dari `source`.
- `resolver` (wajib untuk `shared_key`):
  - `type` : `mapping_table` / `mapping_table_lookup`.
  - `table` : nama map.
  - `source_key_col`, `target_key_col` : kolom pada map.
  - `value_from` : sumber nilai (opsional).
  - `gen` : generator (`type`: mis. `uuid_v7`).

Behavior:
- `natural` : langsung pakai nilai dari `key.source`.
- `shared_key` : cek `_keymap_generic`. Jika belum ada, baris `_exec_queue` akan ditandai `need_keymap` dan checker membuat permintaan map; executor mengambil nilai `RETURNING` untuk mengisi map.

### Columns
Setiap entry `columns.<target_col>` mendukung:
- `from` : `"<alias>.<kolom>"` untuk ambil field terakhir; `"$key"` untuk memakai key hasil resolusi.
- `expr` : ekspresi sederhana via `expr.Evaluate`.
- `cast` : nama cast yang dikenali `cast.Value`.
- `default` : nilai default jika `from`/`expr` tidak terpakai.
- `resolver` : lookup tambahan (mis. foreign key) menggunakan keymap lain.

### Routing
```json
"routing": {
  "on_create": {"mode": "insert"},
  "on_update": {"mode": "update", "matchKey": ["id"]},
  "on_snapshot": {"mode": "insert"}
}
```
- `mode` : `insert` / `update` (default fallback ke insert).
- `matchKey` : daftar kolom target untuk klausa `where` saat update (fallback ke kolom `$key`).

### Fact Condition
Filter event fact sebelum diproses:
- `column` : nama kolom pada payload fact.
- `op` : `equal` | `notEquals`.
- `value` : dibandingkan sebagai string.

### Split Table
Menulis tambahan insert paralel selain target utama:
```json
"split_table": [
  {
    "table_name": "public.order_items",
    "columns": {
      "order_id": {"from": "$key"},
      "sku": {"from": "items_sku"},
      "qty": {"from": "items_qty"}
    }
  }
]
```
Setiap split disusun ulang kolomnya dengan payload yang sama dan di-enqueue sebagai pekerjaan `is_split=true`.

## Contoh Ringkas
### Fact tunggal (tanpa join)
```json
{
  "entity": "customer",
  "sources": [{"alias": "c", "from": "public.customer"}],
  "target_table": "dw.customer",
  "key": {"strategy": "natural", "source": ["c.id"]},
  "columns": {
    "customer_id": {"from": "$key"},
    "name": {"from": "c.name"},
    "created_at": {"from": "c.created_at"}
  },
  "routing": {"on_create": {"mode": "insert"}, "on_update": {"mode": "update"}}
}
```

### Fact + Dimension join
```json
{
  "entity": "order",
  "sources": [
    {"alias": "o", "from": "public.orders"},
    {"alias": "u", "from": "public.users", "join": {"fact_column": "user_id", "dim_column": "id"}}
  ],
  "target_table": "dw.orders",
  "key": {"strategy": "shared_key", "source": ["o.id"], "resolver": {"table": "orders_map", "source_key_col": "src_id", "target_key_col": "tgt_id"}},
  "columns": {
    "order_id": {"from": "$key"},
    "user_name": {"from": "u.name"},
    "amount": {"from": "o.amount"}
  },
  "routing": {"on_create": {"mode": "insert"}, "on_update": {"mode": "update", "matchKey": ["order_id"]}}
}
```
Parser akan menaruh fact ke `_need_join`, fragmen user ke `_join_map_topic`; joiner menunggu keduanya, lalu processor menyusun SQL insert/update dengan `sql_args` sesuai urutan kolom.
