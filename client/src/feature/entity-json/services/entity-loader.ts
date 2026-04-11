import type { LoaderFunctionArgs } from "react-router-dom";

const apiUrl = import.meta.env.VITE_API_URL;

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
