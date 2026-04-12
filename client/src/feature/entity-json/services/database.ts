import type { LoaderFunctionArgs } from "react-router-dom";
import { getEntityByName } from "./entity";

const apiUrl = import.meta.env.VITE_API_URL;

interface getDatabaseSchemaResponse {
  source_schema: Record<string, string>;
  target_schema: Record<string, string>;
}

export async function getDatabaseSchema() {
  try {
    const res = await fetch(`${apiUrl}/database/schema`);
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

    const data = response.data as getDatabaseSchemaResponse;

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

export async function getEntityAndSchema({ params }: LoaderFunctionArgs) {
  // get entity
  const entityName = params?.name;
  if (!entityName) {
    console.error("Entity name parameter is missing in the route.");
    return {
      data: null,
      ok: false,
      message: "Entity name is required",
    };
  }

  const entityData = await getEntityByName(entityName);
  const entity = entityData.data?.entity;
  const filename = entityData.data?.file;

  const schemaData = await getDatabaseSchema();
  const schemas = schemaData.data;

  return {
    data: {
      entity: entity,
      schemas: schemas,
      filename: filename,
    },
    ok: true,
    message: "Entity and schema fetched successfully",
  };
}
