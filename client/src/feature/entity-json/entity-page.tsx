import { useCallback, useEffect, useState } from "react";
import { AlertTriangle, ChevronLeft } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { DefaultPaths } from "@/path";
import { useFetcher, useLoaderData, useNavigate } from "react-router-dom";
import type { EntityDetailLoaderResult } from "./entity-loader";
import { toast } from "sonner";

export function EntityPage() {
  const navigate = useNavigate();
  const loaderData = useLoaderData() as EntityDetailLoaderResult;

  const entity = loaderData.data?.entity;
  const fileName = loaderData.data?.file ?? "Unavailable";
  const entityName = entity?.entity ?? "Unknown Entity";
  const [rawJson, setRawJson] = useState(() =>
    JSON.stringify(entity ?? {}, null, 2)
  );

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setRawJson(JSON.stringify(entity ?? {}, null, 2));
  }, [entity]);

  const formatJson = useCallback(() => {
    try {
      const parsed = JSON.parse(rawJson);
      setRawJson(JSON.stringify(parsed, null, 2));
    } catch (error) {
      console.error("Failed to format JSON:", error);
      toast.error(`Invalid JSON format. Please correct it.`);
    }
  }, [rawJson]);

  const updateFetcher = useFetcher();

  // message when update entity
  useEffect(() => {
    if (updateFetcher.data) {
      if (updateFetcher.data.ok) {
        toast.success(
          updateFetcher.data.message || "Entity updated successfully"
        );
      } else {
        toast.error(updateFetcher.data.message || "Failed to update entity");
      }
    }
  }, [updateFetcher.data]);

  return (
    <>
      <div className="h-screen w-full flex flex-col bg-background">
        <div className="border-b border-border px-6 py-4">
          <div className="flex items-center gap-3">
            <button
              className="p-1 hover:bg-muted rounded-md transition-colors"
              onClick={() => navigate(DefaultPaths.ENTITY_LIST.path)}
              aria-label="Back to entity list"
            >
              <ChevronLeft className="w-5 h-5 text-foreground" />
            </button>
            <div>
              <h1 className="text-2xl font-semibold">{entityName}</h1>
              <p className="text-sm text-muted-foreground">
                Edit your mapping entity here
              </p>
            </div>
          </div>
        </div>

        <div className="flex-1 overflow-auto p-6">
          <updateFetcher.Form
            method="put"
            action="update-entity"
            className="space-y-4"
          >
            <Card className="border border-border bg-card">
              <CardHeader className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                <div className="space-y-3">
                  <div>
                    <CardTitle>Raw entity JSON</CardTitle>
                    <CardDescription>
                      Edit the entity definition directly before saving it.
                      (Preview: {fileName})
                    </CardDescription>
                  </div>
                  <div className="flex gap-2 rounded-md border border-amber-300 bg-amber-50 px-3 py-2 text-sm text-amber-900">
                    <AlertTriangle className="h-4 w-4 shrink-0 text-amber-600" />
                    <p className="leading-relaxed">
                      Make sure the key 'entity' is not changed to avoid issues.
                      If you need to change the entity name, please create a new
                      entity mapping instead.
                    </p>
                  </div>
                </div>
                <Button
                  type="submit"
                  disabled={updateFetcher.state === "submitting"}
                  className="bg-sidebar-primary hover:bg-sidebar-primary/90"
                >
                  Save
                </Button>
              </CardHeader>
              <CardContent>
                <Textarea
                  value={rawJson}
                  onChange={(event) => setRawJson(event.target.value)}
                  onBlur={formatJson}
                  rows={24}
                  spellCheck={false}
                  name="raw_entity"
                  className="font-mono text-base leading-relaxed"
                />
              </CardContent>
            </Card>
          </updateFetcher.Form>
        </div>
      </div>
    </>
  );
}
