import { Fragment, useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ChevronRight, Loader2, RefreshCw, AlertTriangle } from "lucide-react";
import { useLoaderData, useRevalidator } from "react-router-dom";
import type { QueueLogLoaderData } from "./queue-loader";

const statusPalette: Record<string, string> = {
  success: "border-emerald-400/30 bg-emerald-500/10 text-emerald-500",
  completed: "border-emerald-400/30 bg-emerald-500/10 text-emerald-500",
  running: "border-sky-400/30 bg-sky-500/10 text-sky-500",
  processing: "border-sky-400/30 bg-sky-500/10 text-sky-500",
  pending: "border-amber-400/30 bg-amber-500/10 text-amber-500",
  queued: "border-amber-400/30 bg-amber-500/10 text-amber-500",
  error: "border-red-400/30 bg-red-500/10 text-red-500",
  failed: "border-red-400/30 bg-red-500/10 text-red-500",
};

const getStatusBadgeClasses = (status: string) => {
  if (!status) return "border border-border/60 bg-muted text-muted-foreground";
  const paletteKey = status.toLowerCase();
  return (
    "border px-2 py-1 rounded-full text-[11px] font-semibold tracking-wide uppercase " +
    (statusPalette[paletteKey] ||
      "border-border/60 bg-muted text-muted-foreground")
  );
};

const prettyStatusLabel = (status: string) => {
  if (!status) return "Unknown";
  return status
    .replace(/_/g, " ")
    .toLowerCase()
    .replace(/^.|\s./g, (match) => match.toUpperCase());
};

const normalizeNullableField = (value: unknown): string | null => {
  if (value === null || value === undefined) return null;
  if (typeof value === "string") return value.length > 0 ? value : null;
  if (typeof value === "object") {
    const maybe = value as { String?: unknown; Valid?: unknown };
    if (typeof maybe.String === "string") {
      if (typeof maybe.Valid === "boolean" && !maybe.Valid) return null;
      return maybe.String.length > 0 ? maybe.String : null;
    }
    try {
      return JSON.stringify(value);
    } catch {
      return String(value);
    }
  }
  return String(value);
};

const trimContent = (value: unknown, max = 140) => {
  const normalized = normalizeNullableField(value);
  if (!normalized) return "-";
  return normalized.length > max ? `${normalized.slice(0, max)}…` : normalized;
};

export function LatestQueueLogPage() {
  const loaderData = useLoaderData() as QueueLogLoaderData;
  const queueList = loaderData?.data ?? [];
  const { state: revalidatorState, revalidate } = useRevalidator();
  const isRefreshing = revalidatorState !== "idle";
  const [lastUpdated, setLastUpdated] = useState(() => new Date());

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setLastUpdated(new Date());
  }, [loaderData?.data]);

  useEffect(() => {
    const interval = setInterval(() => {
      revalidate();
    }, 10000);

    return () => clearInterval(interval);
  }, [revalidate]);

  const lastUpdatedLabel = lastUpdated.toLocaleTimeString();

  return (
    <div className="h-screen w-full flex flex-col bg-background">
      {/* Header */}
      <div className="border-b border-border px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button
              onClick={() => {}}
              className="p-1 hover:bg-muted rounded-md transition-colors text-foreground"
            >
              <ChevronRight className="w-5 h-5 transform rotate-180" />
            </button>
            <div>
              <h1 className="text-2xl font-semibold text-foreground">
                Latest Log Queue
              </h1>
              <p className="text-sm text-muted-foreground">
                Inspect last runned execution logs and their statuses.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="flex-1 overflow-auto px-6 py-4">
        <Card className="border border-border bg-card">
          <CardHeader>
            <div className="flex flex-wrap items-center justify-between gap-4">
              <div>
                <CardTitle className="text-2xl">Execution Queue</CardTitle>
                <CardDescription className="text-base">
                  Live feed of migration jobs. Data auto-refreshes every 10
                  seconds.
                </CardDescription>
                <p className="text-xs text-muted-foreground mt-2">
                  Menampilkan {queueList.length} antrean aktif. Terakhir update:
                  <span className="font-semibold text-foreground ml-1">
                    {lastUpdatedLabel}
                  </span>
                </p>
              </div>
              <div className="flex items-center gap-3">
                {!loaderData?.ok && (
                  <div className="flex items-center gap-1 rounded-md border border-destructive/40 bg-destructive/10 px-3 py-2 text-destructive text-xs">
                    <AlertTriangle className="h-3.5 w-3.5" />
                    {loaderData?.message || "Gagal memuat data"}
                  </div>
                )}
                <Button
                  type="button"
                  variant="outline"
                  disabled={isRefreshing}
                  onClick={() => revalidate()}
                  className="min-w-[120px]"
                >
                  {isRefreshing ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <RefreshCw className="mr-2 h-4 w-4" />
                  )}
                  Refresh
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent className="p-0">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-border bg-muted/30 text-left text-xs font-semibold text-foreground/70">
                    <th className="px-4 py-3">Queue ID</th>
                    <th className="px-4 py-3">Entity</th>
                    <th className="px-4 py-3">Status</th>
                    <th className="px-4 py-3">SQL Statement</th>
                    <th className="px-4 py-3">Arguments</th>
                    <th className="px-4 py-3">Info</th>
                  </tr>
                </thead>
                <tbody>
                  {queueList.length === 0 ? (
                    <tr>
                      <td
                        colSpan={6}
                        className="py-10 text-center text-sm text-muted-foreground"
                      >
                        Belum ada antrean eksekusi yang tercatat.
                      </td>
                    </tr>
                  ) : (
                    queueList.map((queue) => {
                      const sqlTextFull = normalizeNullableField(queue.SQLText);
                      const sqlArgsFull = normalizeNullableField(queue.SQLArgs);
                      const lastError = normalizeNullableField(queue.LastError);

                      return (
                        <Fragment key={queue.QueueID}>
                          <tr className="border-b border-border/60 hover:bg-muted/20">
                            <td className="px-4 py-3 text-sm font-mono text-foreground">
                              {queue.QueueID}
                            </td>
                            <td className="px-4 py-3">
                              <p className="text-sm font-semibold text-foreground">
                                {queue.Entity}
                              </p>
                              <p className="text-xs text-muted-foreground">
                                {queue.SQLText ? "Primary statement" : "-"}
                              </p>
                            </td>
                            <td className="px-4 py-3">
                              <span
                                className={getStatusBadgeClasses(queue.Status)}
                              >
                                {prettyStatusLabel(queue.Status)}
                              </span>
                            </td>
                            <td className="px-4 py-3">
                              <p
                                className="font-mono text-xs text-foreground/80 break-words"
                                title={sqlTextFull || undefined}
                              >
                                {trimContent(sqlTextFull)}
                              </p>
                            </td>
                            <td className="px-4 py-3">
                              <p
                                className="font-mono text-xs text-muted-foreground break-words"
                                title={sqlArgsFull || undefined}
                              >
                                {trimContent(sqlArgsFull)}
                              </p>
                            </td>
                            <td className="px-4 py-3 text-sm">
                              {lastError ? (
                                <p className="text-destructive font-medium">
                                  Error: {trimContent(lastError, 80)}
                                </p>
                              ) : (
                                <p className="text-emerald-500 font-medium">
                                  No errors
                                </p>
                              )}
                              <p className="text-xs text-muted-foreground mt-1">
                                Split statements: {queue.ExecSplit?.length || 0}
                              </p>
                            </td>
                          </tr>
                          {queue.ExecSplit && queue.ExecSplit.length > 0 && (
                            <tr className="border-b border-border/50 bg-muted/20">
                              <td colSpan={6} className="px-6 py-4">
                                <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                                  Detail Split Statements
                                </p>
                                <div className="mt-3 grid gap-2 md:grid-cols-2">
                                  {queue.ExecSplit.map((split, index) => {
                                    const splitTextFull =
                                      normalizeNullableField(split.SQLText);
                                    const splitArgsFull =
                                      normalizeNullableField(split.SQLArgs);

                                    return (
                                      <div
                                        key={`${queue.QueueID}-${index}`}
                                        className="rounded-md border border-border bg-background/80 p-3"
                                      >
                                        <div className="flex items-center justify-between text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
                                          <span>Step {index + 1}</span>
                                          <span
                                            className={getStatusBadgeClasses(
                                              split.Status
                                            )}
                                          >
                                            {prettyStatusLabel(split.Status)}
                                          </span>
                                        </div>
                                        <p
                                          className="mt-2 font-mono text-xs text-foreground break-words"
                                          title={splitTextFull || undefined}
                                        >
                                          {trimContent(splitTextFull)}
                                        </p>
                                        {splitArgsFull && (
                                          <p
                                            className="mt-1 font-mono text-[11px] text-muted-foreground break-words"
                                            title={splitArgsFull}
                                          >
                                            Args:{" "}
                                            {trimContent(splitArgsFull, 80)}
                                          </p>
                                        )}
                                      </div>
                                    );
                                  })}
                                </div>
                              </td>
                            </tr>
                          )}
                        </Fragment>
                      );
                    })
                  )}
                </tbody>
              </table>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
