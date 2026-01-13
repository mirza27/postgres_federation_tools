export interface DatabaseCredentials {
  source_database: DBCredential;
  target_database: DBCredential;
}

export interface DBCredential {
  type: string;
  host: string;
  port: number;
  user: string;
  password: string;
  db_name: string;
}

const apiUrl = import.meta.env.VITE_API_URL;

export async function databaseCredentialsLoader() {
  try {
    const res = await fetch(`${apiUrl}/database`);
    if (!res.ok) {
      return {
        ok: false,
        message: "Failed to fetch database credentials",
      };
    }

    const response = await res.json();
    if (response.status !== "success") {
      return {
        ok: false,
        message: response.message || "Failed to fetch database credentials",
      };
    }

    const data = response.data as DatabaseCredentials;

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

export async function checkDebeziumConnectorStatus() {
  try {
    const res = await fetch(`${apiUrl}/connector`);
    if (!res.ok) {
      return {
        ok: false,
        message: "Failed to fetch create debezium connector",
      };
    }

    const response = await res.json();
    if (response.status !== "success") {
      console.log("kok sini", response);
      return {
        ok: false,
        message: response.message,
      };
    }

    console.log("DEBEZIUM CONNECTOR STATUS RESPONSE", response);

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
