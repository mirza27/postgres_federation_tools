import type { LoaderFunctionArgs } from "react-router-dom";

const apiUrl = import.meta.env.VITE_API_URL;

interface ExecutionQueueList {
  QueueID: string;
  Entity: string;
  SQLText: string;
  SQLArgs: string;
  Status: string;
  LastError: string | null;
  ExecSplit: Array<{
    SQLText: string;
    SQLArgs: string;
    Status: string;
  }>;
}

export type QueueLogLoaderData = {
  ok: boolean;
  data: ExecutionQueueList[] | null;
  message?: string;
};

export async function getLatestQueueLogs() {
  try {
    const res = await fetch(`${apiUrl}/progress?limit=20`);
    if (!res.ok) {
      return {
        data: null,
        ok: false,
        message: "Failed to fetch execution queue list",
      };
    }

    const response = await res.json();
    if (response.status !== "success") {
      return {
        data: null,
        ok: false,
        message: response.message || "Failed to fetch execution queue list",
      };
    }

    const data = response.data as ExecutionQueueList[];

    return {
      ok: true,
      data: data,
      message: response.message,
    };
  } catch (error) {
    return {
      data: null,
      ok: false,
      message: `Network error: Unable to reach the server, ${error}`,
    };
  }
}

export async function getHistoryQueueLogs({ request }: LoaderFunctionArgs) {
  try {
    const url = new URL(request.url);
    const queryParams = new URLSearchParams();

    const entity = url.searchParams.get("entity");
    const status = url.searchParams.get("status");
    const sqlText = url.searchParams.get("sql_text");
    const sqlArgs = url.searchParams.get("sql_args");
    const page = url.searchParams.get("page");
    const limit = url.searchParams.get("limit");

    if (entity) queryParams.set("filter.entity", entity);
    if (status) queryParams.set("filter.status", status);
    if (sqlText) queryParams.set("search.sql_text", sqlText);
    if (sqlArgs) queryParams.set("search.sql_args", sqlArgs);

    if (page) {
      queryParams.set("page", page);
    } else {
      queryParams.set("page", "1");
    }
    if (limit) {
      queryParams.set("limit", limit);
    } else {
      queryParams.set("limit", "20");
    }

    const queryString = queryParams.toString();
    const endpoint = queryString
      ? `${apiUrl}/progress/queue/list?${queryString}`
      : `${apiUrl}/progress/queue/list`;

    const res = await fetch(endpoint);
    if (!res.ok) {
      return {
        data: null,
        ok: false,
        message: "Failed to fetch execution queue list",
      };
    }

    const response = await res.json();
    if (response.status !== "success") {
      return {
        data: null,
        ok: false,
        message: response.message || "Failed to fetch execution queue list",
      };
    }

    const data = response.data as ExecutionQueueList[];

    return {
      ok: true,
      data: data,
      message: response.message,
    };
  } catch (error) {
    return {
      data: null,
      ok: false,
      message: `Network error: Unable to reach the server, ${error}`,
    };
  }
}
