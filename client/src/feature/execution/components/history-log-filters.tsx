import { useEffect, useState } from "react";
import { useFetcher, type SetURLSearchParams } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { EntityListFilter } from "../services/entity";

type HistoryLogFiltersProps = {
  searchParams: URLSearchParams;
  setSearchParams: SetURLSearchParams;
};

type EntityFilterFetcherData = {
  ok: boolean;
  data: EntityListFilter[] | null;
  message?: string;
};

export function HistoryLogFilters({
  searchParams,
  setSearchParams,
}: HistoryLogFiltersProps) {
  const entityFetcher = useFetcher<EntityFilterFetcherData>();
  const [filterStatus, setFilterStatus] = useState(
    searchParams.get("status") ?? "",
  );
  const [filterEntity, setFilterEntity] = useState(
    searchParams.get("entity") ?? "",
  );
  const [filterSqlText, setFilterSqlText] = useState(
    searchParams.get("sql_text") ?? "",
  );
  const [filterSqlArgs, setFilterSqlArgs] = useState(
    searchParams.get("sql_args") ?? "",
  );
  const [paginationPage, setPaginationPage] = useState(
    parseInt(searchParams.get("page") ?? "1", 10),
  );
  const [paginationLimit, setPaginationLimit] = useState(
    parseInt(searchParams.get("limit") ?? "20", 10),
  );

  useEffect(() => {
    setFilterEntity(searchParams.get("entity") ?? "");
    setFilterStatus(searchParams.get("status") ?? "");
    setFilterSqlArgs(searchParams.get("sql_args") ?? "");
    setFilterSqlText(searchParams.get("sql_text") ?? "");
    setPaginationPage(parseInt(searchParams.get("page") ?? "1", 10));
    setPaginationLimit(parseInt(searchParams.get("limit") ?? "20", 10));
  }, [searchParams]);

  // load entity filters on component mount
  useEffect(() => {
    if (entityFetcher.state === "idle" && !entityFetcher.data) {
      entityFetcher.load("entity-filters");
    }
  }, [entityFetcher]);

  const applyFilters = () => {
    const next = new URLSearchParams(searchParams);

    const setOrDelete = (key: string, value: string) => {
      const trimmedValue = value.trim();
      if (!trimmedValue) {
        next.delete(key);
        return;
      }
      next.set(key, trimmedValue);
    };

    setOrDelete("entity", filterEntity);
    setOrDelete("status", filterStatus);
    setOrDelete("sql_text", filterSqlText);
    setOrDelete("sql_args", filterSqlArgs);

    const safePage =
      Number.isFinite(paginationPage) && paginationPage > 0
        ? paginationPage
        : 1;
    const safeLimit =
      Number.isFinite(paginationLimit) && paginationLimit > 0
        ? paginationLimit
        : 20;

    next.set("page", String(safePage));
    next.set("limit", String(safeLimit));

    setSearchParams(next, { replace: true });
  };

  return (
    <div className="mt-4">
      <div className="grid grid-cols-1 gap-2.5 md:grid-cols-2 xl:grid-cols-12">
        <div className="xl:col-span-3">
          <label className="mb-1 block text-xs font-medium text-muted-foreground">
            Entity
          </label>
          <select
            value={filterEntity}
            onChange={(e) => setFilterEntity(e.target.value)}
            className="h-8 w-full rounded-md border border-input bg-background px-2.5 text-xs text-foreground"
          >
            <option value="">All</option>
            {(entityFetcher.data?.data ?? []).map((item) => (
              <option key={item.entity} value={item.entity}>
                {item.entity}
              </option>
            ))}
          </select>
        </div>
        <div className="xl:col-span-2">
          <label className="mb-1 block text-xs font-medium text-muted-foreground">
            Status
          </label>
          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="h-8 w-full rounded-md border border-input bg-background px-2.5 text-xs text-foreground"
          >
            <option value="">All</option>
            <option value="ready">Ready</option>
            <option value="pending">Pending</option>
            <option value="waiting">Waiting</option>
            <option value="done">Done</option>
            <option value="error">Error</option>
          </select>
        </div>
        <div className="xl:col-span-3">
          <label className="mb-1 block text-xs font-medium text-muted-foreground">
            SQL Text
          </label>
          <Input
            type="text"
            placeholder="Filter by SQL text..."
            value={filterSqlText}
            onChange={(e) => setFilterSqlText(e.target.value)}
            className="h-8 bg-background text-sm placeholder:text-muted-foreground"
          />
        </div>
        <div className="xl:col-span-4">
          <label className="mb-1 block text-xs font-medium text-muted-foreground">
            SQL Args
          </label>
          <Input
            type="text"
            placeholder="Filter by SQL args..."
            value={filterSqlArgs}
            onChange={(e) => setFilterSqlArgs(e.target.value)}
            className="h-8 bg-background text-sm placeholder:text-muted-foreground"
          />
        </div>
      </div>

      <div className="mt-2.5 flex flex-wrap items-end justify-between gap-2.5">
        <div className="flex items-end gap-2.5">
          <div>
            <label className="mb-1 block text-xs font-medium text-muted-foreground">
              Page
            </label>
            <Input
              type="number"
              min={1}
              value={paginationPage}
              onChange={(e) => setPaginationPage(Number(e.target.value) || 1)}
              className="h-7 w-14 bg-background px-2 text-xs"
            />
          </div>
          <div>
            <label className="mb-1 block text-xs font-medium text-muted-foreground">
              Limit
            </label>
            <select
              value={paginationLimit}
              onChange={(e) => setPaginationLimit(Number(e.target.value))}
              className="h-7 w-16 rounded-md border border-input bg-background px-2 text-xs text-foreground"
            >
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
              <option value={100}>100</option>
            </select>
          </div>
        </div>

        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={applyFilters}
          className="h-7 px-3 text-xs"
        >
          Filter
        </Button>
      </div>
    </div>
  );
}
