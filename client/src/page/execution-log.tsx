import { useState } from "react";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { ChevronRight, MoreHorizontal } from "lucide-react";

interface ExecutionLogProps {
  onBack: () => void;
}

export function ExecutionLog({ onBack }: ExecutionLogProps) {
  const [filters, setFilters] = useState({
    entity: "",
    status: "All",
    timeRange: "",
  });

  const mockExecutions = [
    {
      id: "1",
      timestamp: "2024-01-15 14:32:45",
      entity: "users",
      operation: "INSERT",
      items: 1250,
      duration: "2.3s",
      status: "success",
    },
    {
      id: "2",
      timestamp: "2024-01-15 14:15:22",
      entity: "orders",
      operation: "UPSERT",
      items: 856,
      duration: "1.8s",
      status: "success",
    },
    {
      id: "3",
      timestamp: "2024-01-15 13:48:10",
      entity: "products",
      operation: "UPDATE",
      items: 542,
      duration: "3.1s",
      status: "error",
    },
  ];

  return (
    <div className="flex-1 flex flex-col h-screen bg-background">
      {/* Header */}
      <div className="border-b border-border px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button
              onClick={onBack}
              className="p-1 hover:bg-muted rounded-md transition-colors text-foreground"
            >
              <ChevronRight className="w-5 h-5 transform rotate-180" />
            </button>
            <div>
              <h1 className="text-2xl font-semibold text-foreground">
                Execution Log
              </h1>
              <p className="text-sm text-muted-foreground">
                Inspect dry-runs or published executions
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="border-b border-border px-6 py-4">
        <div className="flex gap-4">
          <div className="flex-1">
            <label className="text-xs font-medium text-foreground block mb-1">
              Entity
            </label>
            <Input
              placeholder="Filter by entity..."
              value={filters.entity}
              onChange={(e) =>
                setFilters({ ...filters, entity: e.target.value })
              }
              className="bg-input border border-border text-foreground text-sm"
            />
          </div>
          <div className="flex-1">
            <label className="text-xs font-medium text-foreground block mb-1">
              Status
            </label>
            <select
              value={filters.status}
              onChange={(e) =>
                setFilters({ ...filters, status: e.target.value })
              }
              className="w-full px-3 py-1.5 rounded-md bg-input border border-border text-foreground text-sm"
            >
              <option>All</option>
              <option>Success</option>
              <option>Error</option>
            </select>
          </div>
          <div className="flex-1">
            <label className="text-xs font-medium text-foreground block mb-1">
              Time Range
            </label>
            <Input
              type="text"
              placeholder="Filter by time..."
              value={filters.timeRange}
              onChange={(e) =>
                setFilters({ ...filters, timeRange: e.target.value })
              }
              className="bg-input border border-border text-foreground text-sm"
            />
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="flex-1 overflow-auto px-6 py-4">
        <Card className="border border-border bg-card">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-border">
                  <th className="text-left px-4 py-3 text-xs font-semibold text-foreground/70 bg-muted/30">
                    Timestamp
                  </th>
                  <th className="text-left px-4 py-3 text-xs font-semibold text-foreground/70 bg-muted/30">
                    Entity
                  </th>
                  <th className="text-left px-4 py-3 text-xs font-semibold text-foreground/70 bg-muted/30">
                    Operation
                  </th>
                  <th className="text-right px-4 py-3 text-xs font-semibold text-foreground/70 bg-muted/30">
                    Items
                  </th>
                  <th className="text-left px-4 py-3 text-xs font-semibold text-foreground/70 bg-muted/30">
                    Duration
                  </th>
                  <th className="text-left px-4 py-3 text-xs font-semibold text-foreground/70 bg-muted/30">
                    Status
                  </th>
                  <th className="text-center px-4 py-3 text-xs font-semibold text-foreground/70 bg-muted/30">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody>
                {mockExecutions.length === 0 ? (
                  <tr>
                    <td
                      colSpan={7}
                      className="text-center py-8 text-muted-foreground text-sm"
                    >
                      No executions found
                    </td>
                  </tr>
                ) : (
                  mockExecutions.map((exec) => (
                    <tr
                      key={exec.id}
                      className="border-b border-border hover:bg-muted/30 transition-colors"
                    >
                      <td className="px-4 py-3 text-sm text-foreground font-mono">
                        {exec.timestamp}
                      </td>
                      <td className="px-4 py-3 text-sm text-foreground">
                        {exec.entity}
                      </td>
                      <td className="px-4 py-3 text-sm text-foreground font-mono">
                        {exec.operation}
                      </td>
                      <td className="px-4 py-3 text-sm text-foreground text-right">
                        {exec.items}
                      </td>
                      <td className="px-4 py-3 text-sm text-foreground">
                        {exec.duration}
                      </td>
                      <td className="px-4 py-3 text-sm">
                        <span
                          className={`px-2 py-1 rounded text-xs font-medium ${
                            exec.status === "success"
                              ? "bg-green-500/20 text-green-300"
                              : "bg-red-500/20 text-red-300"
                          }`}
                        >
                          {exec.status === "success" ? "Success" : "Error"}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-center">
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-7 w-7 p-0"
                        >
                          <MoreHorizontal className="w-4 h-4" />
                        </Button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </Card>
      </div>
    </div>
  );
}
