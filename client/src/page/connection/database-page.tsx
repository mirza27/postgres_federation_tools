import { useLoaderData, useFetcher } from "react-router-dom";
import { useEffect, useState } from "react";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import type { DatabaseCredentials, DBCredential } from "./loader";
import { ChevronLeft, Save } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { toast } from "sonner";

const emptyCred: DBCredential = {
  type: "",
  host: "",
  port: 0,
  user: "",
  password: "",
  db_name: "",
};

export function DatabasePage() {
  const loaderData = useLoaderData() as {
    ok: boolean;
    data: DatabaseCredentials | undefined;
    message?: string;
  };

  const [source, setSource] = useState<DBCredential>(
    loaderData.data?.source_database ?? emptyCred
  );
  const [target, setTarget] = useState<DBCredential>(
    loaderData.data?.target_database ?? emptyCred
  );

  const fetcherSource = useFetcher();
  const fetcherTarget = useFetcher();

  useEffect(() => {
    if (loaderData.ok === false) {
      toast.error(loaderData.message);
    }
  }, [loaderData]);

  useEffect(() => {
    if (!fetcherSource.data) return;
    const data = fetcherSource.data as {
      ok: boolean;
      message: string;
    };
    if (data.ok) toast.success(data.message);
    else toast.error(data.message);
  }, [fetcherSource.data]);

  useEffect(() => {
    if (!fetcherTarget.data) return;
    const data = fetcherTarget.data as {
      ok: boolean;
      message: string;
    };
    if (data.ok) {
      toast.success(data.message);
    } else toast.error(data.message);
  }, [fetcherTarget.data]);

  return (
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
              Configure source and target database
            </p>
          </div>
        </div>
      </div>

      {/* content */}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* SOURCE */}
        <fetcherSource.Form
          method="post"
          action="save-source"
          className="flex-1 overflow-auto p-6"
        >
          <Card className="p-6 border border-border bg-card">
            <h2 className="text-lg font-semibold mb-6">
              Source Database Credentials
            </h2>

            <div className="space-y-4">
              {/* database type */}
              <div className="grid grid-cols-3 items-center gap-4">
                <Label className="text-md font-medium text-foreground text-right">
                  Database Type <span className="text-red-500">*</span>
                </Label>
                <Select
                  value={source.type}
                  name="type"
                  onValueChange={(value) =>
                    setSource({ ...source, type: value })
                  }
                  defaultValue="postgres"
                >
                  <SelectTrigger className="bg-input border border-border text-foreground">
                    <SelectValue placeholder="Select database type" />
                  </SelectTrigger>
                  <SelectContent className="bg-card border border-border">
                    <SelectItem value="postgres">PostgreSQL</SelectItem>
                    <SelectItem value="mysql">MySQL</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* hostname */}
              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-host"
                  className="text-md font-medium text-foreground text-right"
                >
                  Host
                  <span className="text-red-500">*</span>
                </Label>

                <Input
                  id="source-host"
                  placeholder="example.db or 128.X.X.X"
                  name="host"
                  value={source.host}
                  onChange={(e) =>
                    setSource({ ...source, host: e.target.value })
                  }
                  className="col-span-2"
                />
              </div>

              {/* port */}
              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-port"
                  className="text-md font-medium text-foreground text-right"
                >
                  Port
                  <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="source-port"
                  type="number"
                  name="port"
                  placeholder="5432"
                  value={source.port}
                  onChange={(e) =>
                    setSource({ ...source, port: Number(e.target.value) })
                  }
                  className="col-span-2"
                />
              </div>

              {/* user */}
              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-username"
                  className="text-md font-medium text-foreground text-right"
                >
                  Username
                  <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="source-username"
                  value={source.user}
                  name="user"
                  placeholder="postgres"
                  onChange={(e) =>
                    setSource({ ...source, user: e.target.value })
                  }
                  className="col-span-2"
                />
              </div>

              {/* password */}
              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-password"
                  className="text-md font-medium text-foreground text-right"
                >
                  Password
                  <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="source-password"
                  type="password"
                  name="password"
                  placeholder="adminXXX"
                  value={source.password}
                  onChange={(e) =>
                    setSource({ ...source, password: e.target.value })
                  }
                  className="col-span-2"
                />
              </div>
              {/* db name */}
              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-database"
                  className="text-md font-medium text-foreground text-right"
                >
                  Database Name
                  <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="source-database"
                  placeholder="e-commerce-db"
                  name="db_name"
                  value={source.db_name}
                  onChange={(e) =>
                    setSource({ ...source, db_name: e.target.value })
                  }
                  className="col-span-2"
                />
              </div>
            </div>

            {/* submit button */}
            <div className="flex justify-end gap-2">
              <Button
                type="submit"
                size="sm"
                className="gap-2 bg-sidebar-primary hover:bg-sidebar-primary/90"
              >
                <Save className="w-4 h-4" />
                Save and Check Connection
              </Button>
            </div>
          </Card>
        </fetcherSource.Form>

        {/* TARGET */}
        <fetcherTarget.Form
          method="post"
          action="save-target"
          className="flex-1 overflow-auto p-6"
        >
          <Card className="p-6 border border-border bg-card">
            <h2 className="text-lg font-semibold mb-6">
              Target Database Credentials
            </h2>

            <div className="space-y-4">
              {/* database type */}
              <div className="grid grid-cols-3 items-center gap-4">
                <Label className="text-md font-medium text-foreground text-right">
                  Database Type <span className="text-red-500">*</span>
                </Label>
                <Select
                  value={target.type}
                  name="type"
                  onValueChange={(value) =>
                    setTarget({ ...target, type: value })
                  }
                  defaultValue="postgres"
                >
                  <SelectTrigger className="bg-input border border-border text-foreground">
                    <SelectValue placeholder="Select database type" />
                  </SelectTrigger>
                  <SelectContent className="bg-card border border-border">
                    <SelectItem value="postgres">PostgreSQL</SelectItem>
                    <SelectItem value="mysql">MySQL</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* hostname */}
              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-host"
                  className="text-md font-medium text-foreground text-right"
                >
                  Host
                  <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="target-host"
                  placeholder="example.db or 128.X.X.X"
                  value={target.host}
                  name="host"
                  onChange={(e) =>
                    setTarget({ ...target, host: e.target.value })
                  }
                  className="col-span-2"
                />
              </div>

              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-port"
                  className="text-md font-medium text-foreground text-right"
                >
                  Port
                  <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="target-port"
                  type="number"
                  name="port"
                  placeholder="5432"
                  value={target.port}
                  onChange={(e) =>
                    setTarget({ ...target, port: Number(e.target.value) })
                  }
                  className="col-span-2"
                />
              </div>

              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-username"
                  className="text-md font-medium text-foreground text-right"
                >
                  Username
                  <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="target-username"
                  value={target.user}
                  name="user"
                  placeholder="postgres"
                  onChange={(e) =>
                    setTarget({ ...target, user: e.target.value })
                  }
                  className="col-span-2"
                />
              </div>

              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="target-password"
                  className="text-md font-medium text-foreground text-right"
                >
                  Password
                </Label>
                <Input
                  id="target-password"
                  type="password"
                  name="password"
                  placeholder="adminXXX"
                  value={target.password}
                  onChange={(e) =>
                    setTarget({ ...target, password: e.target.value })
                  }
                  className="col-span-2"
                />
              </div>

              <div className="grid grid-cols-3 items-center gap-4">
                <Label
                  htmlFor="source-database"
                  className="text-md font-medium text-foreground text-right"
                >
                  Database Name
                  <span className="text-red-500">*</span>
                </Label>
                <Input
                  id="target-database"
                  value={target.db_name}
                  name="db_name"
                  placeholder="e-commerce-db"
                  onChange={(e) =>
                    setTarget({ ...target, db_name: e.target.value })
                  }
                  className="col-span-2"
                />
              </div>
            </div>

            {/* submit button */}
            <div className="flex justify-end gap-2">
              <Button
                type="submit"
                size="sm"
                className="gap-2 bg-sidebar-primary hover:bg-sidebar-primary/90"
              >
                <Save className="w-4 h-4" />
                Save and Check Connection
              </Button>
            </div>
          </Card>
        </fetcherTarget.Form>
      </div>
    </div>
  );
}
