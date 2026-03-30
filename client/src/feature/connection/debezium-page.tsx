/* eslint-disable react-hooks/set-state-in-effect */
import { useEffect, useState } from "react";
import { useLoaderData, useFetcher, useLocation } from "react-router-dom";
import { ChevronLeft, RefreshCw } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export function DebeziumPage() {
  const loaderData = useLoaderData() as {
    ok: boolean;
    message?: string;
  };
  const [connectorStatus, setConnectorStatus] = useState(
    loaderData?.ok ?? false
  );
  const [message, setMessage] = useState(loaderData?.message ?? "");

  const fetcher = useFetcher();
  const location = useLocation();
  const isChecking = fetcher.state !== "idle";

  useEffect(() => {
    setConnectorStatus(loaderData?.ok ?? false);
    setMessage(loaderData?.message ?? "");
  }, [loaderData]);

  useEffect(() => {
    if (!fetcher.data) return;
    const d = fetcher.data as {
      ok?: boolean;
      message?: string;
    };

    if (typeof d.ok === "boolean") {
      setConnectorStatus(d.ok);
      toast.success("debezium is connected successfully");
    } else {
      toast.error("failed to connect debezium");
    }

    if (d.message) setMessage(d.message);
  }, [fetcher.data]);

  const handleCheckStatus = () => {
    // same as refresh page but without routing
    fetcher.load(location.pathname);
  };

  return (
    <>
      <div className="h-screen w-full flex flex-col bg-background">
        {/* header */}
        <div className="border-b border-border px-6 py-4">
          <div className="flex items-center gap-3">
            <button className="p-1 hover:bg-muted rounded-md transition-colors">
              <ChevronLeft className="w-5 h-5 text-foreground" />
            </button>
            <div>
              <h1 className="text-2xl font-semibold">Database Connection</h1>
              <p className="text-sm text-muted-foreground">
                Configure Debezium Connector (CDC Tool)
              </p>
            </div>
          </div>
        </div>

        {/* content */}
        <div className="flex-1 overflow-auto p-6">
          <Card className="border border-border bg-card">
            <div className="grid grid-cols-1 md:grid-cols-2">
              {/* Tutorial column */}
              <CardContent className="space-y-6 border-b border-border md:border-b-0 md:border-r p-6">
                <div>
                  <CardHeader className="px-0 pb-4">
                    <CardTitle>
                      PostgreSQL to Debezium Configuration Guide
                    </CardTitle>
                    <CardDescription>
                      Essential steps to configure PostgreSQL for Debezium CDC
                      (Change Data Capture)
                    </CardDescription>
                  </CardHeader>
                  <p className="font-semibold mb-3">
                    🚀 Step-by-Step Configuration:
                  </p>
                  <ol className="list-decimal flex flex-col gap-3 pl-6 text-sm">
                    <li>
                      <div className="font-medium">
                        Access PostgreSQL Configuration File
                      </div>
                      <pre className="text-xs bg-muted text-foreground/90 p-3 rounded mt-1 overflow-x-auto">
                        {`# Navigate to PostgreSQL data directory
cd /var/lib/postgresql/data/

# Edit the main configuration file
sudo nano postgresql.conf

# Alternative locations if above doesn't exist:
# Ubuntu/Debian: /etc/postgresql/{version}/main/postgresql.conf
# RHEL/CentOS: /var/lib/pgsql/data/postgresql.conf`}
                      </pre>
                    </li>
                    <li>
                      <div className="font-medium">
                        Update Critical Parameters
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">
                        Locate and modify these lines:
                      </div>
                      <pre className="text-xs bg-muted text-foreground/90 p-3 rounded mt-1 overflow-x-auto">
                        {`# Change from 'replica' or 'minimal' to 'logical'
wal_level = logical

# Ensure at least 1 replication slot (increase if needed)
max_replication_slots = 1

# Optional but recommended parameters:
# max_wal_senders = 1
# wal_keep_size = 64MB`}
                      </pre>
                    </li>
                    <li>
                      <div className="font-medium">
                        Restart PostgreSQL Service
                      </div>
                      <pre className="text-xs bg-muted text-foreground/90 p-3 rounded mt-1 overflow-x-auto">
                        {`# For systemd systems (Ubuntu 16.04+, CentOS 7+):
sudo systemctl restart postgresql
# or specific version
sudo systemctl restart postgresql-15

# For older init.d systems:
sudo service postgresql restart

# Verify service status
sudo systemctl status postgresql`}
                      </pre>
                    </li>
                    <li>
                      <div className="font-medium">
                        Connect to PostgreSQL via psql
                      </div>
                      <pre className="text-xs bg-muted text-foreground/90 p-3 rounded mt-1 overflow-x-auto">
                        {`# Connect as postgres user (adjust as needed)
sudo -u postgres psql

# Or with explicit credentials
psql -h localhost -U postgres -d postgres

# For Docker containers:
docker exec -it postgres_container psql -U postgres`}
                      </pre>
                    </li>
                    <li>
                      <div className="font-medium">
                        Verify Configuration Applied
                      </div>
                      <pre className="text-xs bg-muted text-foreground/90 p-3 rounded mt-1 overflow-x-auto">
                        {`-- Check WAL level (should return 'logical')
SHOW wal_level;

-- Verify replication slots parameter
SHOW max_replication_slots;

-- Additional useful checks
SELECT name, setting FROM pg_settings 
WHERE name IN ('wal_level', 'max_replication_slots', 'max_wal_senders');`}
                      </pre>
                    </li>
                  </ol>
                </div>

                <div className="p-4 bg-primary/5 rounded-lg border border-border text-sm">
                  <p className="font-semibold mb-2 text-foreground/80">
                    📝 Important Notes
                  </p>
                  <ul className="list-disc pl-5 space-y-2 text-muted-foreground">
                    <li>
                      <strong>Privilege Requirements:</strong> Debezium user
                      needs
                      <span className="bg-muted px-1 rounded ml-1">
                        REPLICATION
                      </span>
                      and
                      <span className="bg-muted px-1 rounded ml-1">SELECT</span>
                      privileges
                    </li>
                    <li>
                      <strong>Security:</strong> Update
                      <span className="bg-muted px-1 rounded mx-1">
                        pg_hba.conf
                      </span>
                      when enabling remote access
                    </li>
                    <li>
                      <strong>Monitoring:</strong> Inspect replication slots via
                      <span className="bg-muted px-1 rounded ml-1">
                        SELECT * FROM pg_replication_slots;
                      </span>
                    </li>
                    <li>
                      <strong>Troubleshooting:</strong> Ensure edits are in the
                      correct config file and restart PostgreSQL when changes do
                      not apply
                    </li>
                  </ul>
                </div>
              </CardContent>

              {/* Status column */}
              <CardContent className="p-6 flex flex-col gap-6">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="text-lg font-semibold">
                        Connector Status
                      </h3>
                      <p className="text-sm text-muted-foreground">
                        Current Debezium connector health
                      </p>
                    </div>
                    <span
                      className={`inline-flex items-center rounded-full border px-3 py-1 text-sm font-medium ${
                        connectorStatus
                          ? "border-emerald-200 bg-emerald-50 text-emerald-700"
                          : "border-red-200 bg-red-50 text-red-700"
                      }`}
                    >
                      {connectorStatus ? "Connected" : "Disconnected"}
                    </span>
                  </div>

                  <div className="rounded-md border border-dashed border-border p-4 text-sm text-muted-foreground bg-muted/30">
                    {message || "No status message yet"}
                  </div>
                </div>

                <div>
                  <Button
                    type="button"
                    onClick={handleCheckStatus}
                    disabled={isChecking}
                    className="w-full gap-2"
                  >
                    <RefreshCw
                      className={`h-4 w-4 ${isChecking ? "animate-spin" : ""}`}
                    />
                    {isChecking ? "Checking connector..." : "Check Connector"}
                  </Button>
                </div>
              </CardContent>
            </div>
          </Card>
        </div>
      </div>
    </>
  );
}
