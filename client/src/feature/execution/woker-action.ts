const apiUrl = import.meta.env.VITE_API_URL;

export async function RunWorkerAction() {
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

export async function StopWorkerAction() {
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
