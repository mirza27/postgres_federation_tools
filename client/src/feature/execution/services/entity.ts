const apiUrl = import.meta.env.VITE_API_URL;

export interface EntityListFilter {
  entity: string;
}

export async function getEntityFilters() {
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

    // cast to get entity name only
    const data = response.data;
    const entities: EntityListFilter[] = data.entities.map((item: any) => ({
      entity: item.entity,
    }));

    return {
      ok: true,
      data: entities,
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
