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

export async function GetLatestQueueLogsLoader() {
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

export async function GetHistoryQueueLogsLoader() {
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