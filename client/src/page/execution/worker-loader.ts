const apiUrl = import.meta.env.VITE_API_URL;

export async function CheckWorkerStatus() {
  try {
    const res = await fetch(`${apiUrl}/worker/status`);
    if (!res.ok) {
      return {
        ok: false,
        message: "Failed to fetch worker status",
      };
    }

    const response = await res.json();
    if (response.status !== "success") {
      return {
        ok: false,
        message: response.message || "Failed to fetch worker status",
      };
    }

    return {
      ok: true,
      message: response.message,
    };
  } catch (error) {
    return {
      ok: false,
      message: `Network error: Unable to reach the server, ${error}`,
    };
  }
}
