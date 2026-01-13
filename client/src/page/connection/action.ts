import type { ActionFunctionArgs } from "react-router-dom";
import type { DBCredential } from "./loader";

const apiUrl = import.meta.env.VITE_API_URL;

export async function saveDatabaseSourceCredentials({
  request,
}: ActionFunctionArgs) {
  const data = await request.formData();

  const requestedData: DBCredential = {
    type: data.get("type") as string,
    host: data.get("host") as string,
    port: Number(data.get("port")),
    user: data.get("user") as string,
    password: data.get("password") as string,
    db_name: data.get("db_name") as string,
  };

  try {
    const res = await fetch(`${apiUrl}/database/source`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(requestedData),
    });

    const response = await res.json();
    if (!res.ok || response.status !== "success") {
      return {
        ok: false,
        message:
          response.message || "Failed to save source database credentials",
      };
    }
    return {
      ok: true,
      message: response.message,
    }; // refetch
  } catch (error) {
    return {
      ok: false,
      message: `Network error: Unable to reach the server, ${error}`,
    };
  }
}

export async function saveDatabaseTargetCredentials({
  request,
}: ActionFunctionArgs) {
  const data = await request.formData();

  const requestedData: DBCredential = {
    type: data.get("type") as string,
    host: data.get("host") as string,
    port: Number(data.get("port")),
    user: data.get("user") as string,
    password: data.get("password") as string,
    db_name: data.get("db_name") as string,
  };

  try {
    const res = await fetch(`${apiUrl}/database/target`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(requestedData),
    });

    const response = await res.json();
    if (!res.ok || response.status !== "success") {
      return {
        ok: false,
        message:
          response.message || "Failed to save target database credentials",
      };
    }
    return {
      ok: true,
      message: response.message,
    }; // refetch
  } catch (error) {
    return {
      ok: false,
      message: `Network error: Unable to reach the server, ${error}`,
    };
  }
}
