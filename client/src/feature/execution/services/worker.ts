const apiUrl = import.meta.env.VITE_API_URL;

export async function checkWorkerStatus() {
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


export async function runWorker() {
  try {
    const res = await fetch(`${apiUrl}/worker/start`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
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

export async function stopWorker() {
  try {
    const res = await fetch(`${apiUrl}/worker/stop`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
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
