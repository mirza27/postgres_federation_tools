import type { ActionFunctionArgs, LoaderFunctionArgs } from "react-router-dom";

const apiUrl = import.meta.env.VITE_API_URL;

type KeyStrategy = "natural" | "shared_key";
type CastType = "string" | "number" | "boolean" | "int";
type ExprType = "now" | "uuid" | "random_int";

interface SourceConfig {
  alias: string;
  from: string;
}

interface KeyConfig {
  strategy: KeyStrategy;
  source: string;
}

interface ColumnConfig {
  from: string;
  cast?: CastType;
  expr?: ExprType;
}

type ColumnsConfig = Record<string, ColumnConfig>;

interface SplitTableConfig {
  table_name: string;
  columns: ColumnsConfig;
}

interface RoutingConfig {
  on_snapshot: { mode: "insert" | "update" | "delete" };
  on_create: { mode: "insert" | "update" | "delete" };
  on_update: { mode: "insert" | "update" | "delete"; matchKey: string[] };
}

export interface BaseEntityConfig {
  entity: string;
  sources: SourceConfig[];
  target_table: string;
  key: KeyConfig;
  columns: ColumnsConfig;
  split_table?: SplitTableConfig[];
  routing?: RoutingConfig;
}

const newEmptyEntity: BaseEntityConfig = {
  entity: "",
  sources: [
    {
      alias: "",
      from: "",
    },
  ],
  target_table: "",
  key: {
    strategy: "natural",
    source: "",
  },
  columns: {},
  split_table: [],
  routing: {
    on_snapshot: { mode: "insert" },
    on_create: { mode: "insert" },
    on_update: { mode: "update", matchKey: ["code"] },
  },
};

// create empty entities
export async function createEmptyEntity({ request }: ActionFunctionArgs) {
  const data = await request.formData();

  const entityName = data.get("entity-name") as string;

  const createdEntity = { ...newEmptyEntity, entity: entityName };

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

// list entities
export interface EntityListDataResponse {
  entities: Array<any>;
  files: Array<string>;
}

export async function ListAllEntitiesLoader() {
  try {
    const res = await fetch(`${apiUrl}/mapping/list`);
    if (!res.ok) {
      return {
        data: null,
        ok: false,
        message: "Failed to fetch entities",
      };
    }

    const response = await res.json();
    if (response.status !== "success") {
      return {
        data: null,
        ok: false,
        message: response.message || "Failed to fetch entities",
      };
    }

    const data = response.data as EntityListDataResponse;

    return {
      ok: true,
      data: data,
      message: response.message,
    };
  } catch (error) {
    return {
      data: null,
      ok: false,
      message: `Internal server error, ${error}`,
    };
  }
}

// get entity-detail
export interface EntityDataDetailResponse {
  entity: any;
  file: string;
}

export async function getEntityByName(name: string) {
  try {
    const queryParams = new URLSearchParams();
    queryParams.append("name", name);

    const res = await fetch(`${apiUrl}/mapping?${queryParams.toString()}`);
    if (!res.ok) {
      return {
        data: null,
        ok: false,
        message: "Failed to fetch entity details",
      };
    }

    const response = await res.json();
    if (response.status !== "success") {
      return {
        data: null,
        ok: false,
        message: response.message,
      };
    }

    const data = response.data as EntityDataDetailResponse;

    return {
      ok: true,
      data: data,
      message: response.message,
    };
  } catch (error) {
    return {
      data: null,
      ok: false,
      message: `Internal server error, ${error}`,
    };
  }
}

export async function GetEntityLoader({ params }: LoaderFunctionArgs) {
  const entityName = params?.name;
  if (!entityName) {
    console.error("Entity name parameter is missing in the route.");
    return {
      data: null,
      ok: false,
      message: "Entity name is required",
    };
  }

  return getEntityByName(entityName);
}

export type EntityDetailLoaderResult = Awaited<
  ReturnType<typeof getEntityByName>
>;

const getSourceTables = (sources: any[]) =>
  sources.map((s) => s.from).join(", ");

const getTargetTables = (entity: any) => {
  const main = entity.target_table;
  const splits = entity.split_table?.map((t: any) => t.table_name) ?? [];
  return [main, ...splits].join(", ");
};

export { getSourceTables, getTargetTables };

// update entity
export async function updateEntity({ request }: ActionFunctionArgs) {
  const data = await request.formData();

  const rawJsonEntity = data.get("raw_entity") as string;
  const parsedEntity = JSON.parse(rawJsonEntity);

  const entityName = parsedEntity.entity as string;
  const queryParams = new URLSearchParams();
  queryParams.append("name", entityName);

  try {
    const res = await fetch(`${apiUrl}/mapping?${queryParams.toString()}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(parsedEntity),
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

// delete entity
export async function deleteEntity({ request }: ActionFunctionArgs) {
  const data = await request.formData();

  const entityName = data.get("entity_name") as string;
  const queryParams = new URLSearchParams();
  queryParams.append("name", entityName);

  try {
    const res = await fetch(`${apiUrl}/mapping?${queryParams.toString()}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    const response = await res.json();
    if (response.status !== "success") {
      return {
        data: null,
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
