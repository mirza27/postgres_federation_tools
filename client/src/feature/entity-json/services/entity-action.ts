import type { ActionFunctionArgs } from "react-router-dom";

const apiUrl = import.meta.env.VITE_API_URL;

const newBaseEntity = {
  entity: "<entity-name>",
  sources: [
    {
      alias: "<alias-table>",
      from: "<source-table-name>",
    },
  ],
  target_table: "<target-table-name>",
  key: {
    strategy: "<natural|shared_key>",
    source: "<alias-table.source-column-name>",
    resolver: {
      table: "<required if strategy is shared_key --> any name>",
      type: "<required if strategy is shared_key --> any name>",
    },
  },
  columns: {
    id: {
      from: "$key",
      expr: "<optional --> now|uuid|random_int>",
    },
    column_name1: {
      from: "<alias-table.source-column-name>",
      cast: "<optional --> string|number|boolean|int>",
    },
    column_name2: {
      from: "<alias-table.source-column-name>",
    },
  },
  split_table: [
    {
      table_name: "<split-table-name>",
      columns: {
        column_name1: {
          from: "<alias-table.source-column-name>",
          cast: "<optional --> string|number|boolean|int>",
        },
      },
    },
  ],
  routing: {
    on_snapshot: {
      mode: "insert",
    },
    on_create: {
      mode: "insert",
    },
    on_update: {
      mode: "update",
      matchKey: ["<oneofcolumns-names>"],
    },
  },
};

export async function CreateNewEntity({ request }: ActionFunctionArgs) {
  console.log("Creating new entity mapping...");
  const data = await request.formData();

  const entityName = data.get("entity-name") as string;

  const createdEntity = { ...newBaseEntity, entity: entityName };

  try {
    const res = await fetch(`${apiUrl}/mapping`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(createdEntity),
    });

    const response = await res.json();
    if (response.status !== "success") {
      return {
        message: response.message,
      };
    }

    return {
      ok: true,
      message: response.message,
    };
  } catch (error) {
    return {
      ok: false,
      message: `Internal server error, ${error}`,
    };
  }
}
