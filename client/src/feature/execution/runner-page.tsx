import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import type { DefaultLoaderResponse } from "@/helper/loader";
import { ChevronLeft, Loader2, PauseCircle, Play, Square } from "lucide-react";
import { useFetcher, useLoaderData } from "react-router-dom";

export function RunnerPage() {
  const loaderData = useLoaderData() as DefaultLoaderResponse<any>;

  const isWorkerRunning = loaderData.ok ?? false;

  const startWorkerFetcher = useFetcher();
  const stopWorkerFetcher = useFetcher();
  const isStarting = startWorkerFetcher.state !== "idle";
  const isStopping = stopWorkerFetcher.state !== "idle";
  const statusClassName = isWorkerRunning
    ? "border-emerald-200 bg-emerald-50 text-emerald-700"
    : "border-muted bg-muted text-muted-foreground";
  const statusLabel = isWorkerRunning ? "Runner Active" : "Runner Idle";
  const migrationMessage = isWorkerRunning
    ? "Migration in progress."
    : "Migration has not started yet. Turn on the worker to begin the data migration process.";

  return (
    <div className="h-screen w-full flex flex-col bg-background">
      {/* header */}
      <div className="border-b border-border px-6 py-4">
        <div className="flex items-center gap-3">
          <button className="p-1 hover:bg-muted rounded-md transition-colors">
            <ChevronLeft className="w-5 h-5 text-foreground" />
          </button>
          <div>
            <h1 className="text-2xl font-semibold">Start Migration</h1>
            <p className="text-sm text-muted-foreground">
              Check and start runner for data migration
            </p>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6">
        <Card className="max-w border border-border bg-card">
          <CardHeader>
            <div className="flex items-center justify-between gap-4">
              <div>
                <CardTitle className="text-2xl">Migration Runner</CardTitle>
                <CardDescription className="text-base">
                  Monitor and control the worker responsible for executing the
                  migration jobs.
                </CardDescription>
              </div>
              <span
                className={
                  "inline-flex items-center gap-2 rounded-full border px-3 py-1 text-sm font-medium " +
                  statusClassName
                }
              >
                {isWorkerRunning ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <PauseCircle className="h-4 w-4" />
                )}
                {statusLabel}
              </span>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="rounded-lg border border-dashed border-border/80 bg-muted/40 p-4">
              <p className="text-sm text-muted-foreground">
                {migrationMessage}
              </p>
            </div>
          </CardContent>
          <CardFooter className="flex flex-wrap items-center gap-3">
            <startWorkerFetcher.Form method="post" action="run-worker">
              <Button
                type="submit"
                disabled={isWorkerRunning || isStarting}
                className="min-w-[150px]"
              >
                {isStarting ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Play className="mr-2 h-4 w-4" />
                )}
                Start Worker
              </Button>
            </startWorkerFetcher.Form>

            <stopWorkerFetcher.Form method="post" action="stop-worker">
              <Button
                type="submit"
                variant="destructive"
                disabled={!isWorkerRunning || isStopping}
                className="min-w-[150px]"
              >
                {isStopping ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Square className="mr-2 h-4 w-4" />
                )}
                Stop Worker
              </Button>
            </stopWorkerFetcher.Form>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
