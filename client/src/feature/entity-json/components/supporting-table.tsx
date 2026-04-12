import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";
import { Input } from "@/components/ui/input";
import { generateAlias } from "./helper";

interface SupportingTableProps {
  item: {
    id: string;
    sourceTable: string;
    alias: string;
    factColumn: string;
    dimColumn: string;
  };
  onUpdate: (id: string, updates: any) => void;
  onRemove: (id: string) => void;
  allTables: string[];
  sourceColumns: string[];
  allColumnsMap: Record<string, string[]>;
}

export function SupportingTable({
  item,
  onUpdate,
  onRemove,
  allTables,
  sourceColumns,
  allColumnsMap,
}: SupportingTableProps) {
  const supportingTableColumns = item.sourceTable
    ? allColumnsMap[item.sourceTable as keyof typeof allColumnsMap] || []
    : [];

  return (
    <div className="bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg p-4">
      <div className="grid grid-cols-1 md:grid-cols-5 gap-3 items-end">
        <div className="flex flex-col gap-1">
          <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
            Tabel Pendukung
          </label>
          <Select
            value={item.sourceTable}
            onValueChange={(value) =>
              onUpdate(item.id, {
                sourceTable: value,
                alias: generateAlias(value),
                dimColumn: "",
              })
            }
            required
          >
            <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
              <SelectValue placeholder="Pilih..." />
            </SelectTrigger>
            <SelectContent>
              {allTables.map((table) => (
                <SelectItem key={table} value={table}>
                  {table}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="flex flex-col gap-1">
          <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
            Alias
          </label>
          <Input value={item.alias} disabled readOnly required />
        </div>

        <div className="flex flex-col gap-1">
          <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
            Kolom Fact (Sumber)
          </label>
          <Select
            value={item.factColumn}
            onValueChange={(value) => onUpdate(item.id, { factColumn: value })}
            required
          >
            <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
              <SelectValue placeholder="Pilih..." />
            </SelectTrigger>
            <SelectContent>
              {sourceColumns.map((col) => (
                <SelectItem key={col} value={col}>
                  {col}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="flex flex-col gap-1">
          <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
            Kolom Dim (Pendukung)
          </label>
          <Select
            value={item.dimColumn}
            onValueChange={(value) => onUpdate(item.id, { dimColumn: value })}
            required
          >
            <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
              <SelectValue placeholder="Pilih..." />
            </SelectTrigger>
            <SelectContent>
              {supportingTableColumns.map((col) => (
                <SelectItem key={col} value={col}>
                  {col}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="flex flex-col gap-1">
          <label className="text-xs font-semibold text-slate-600 dark:text-slate-400 invisible">
            Action
          </label>
          <Button
            onClick={() => onRemove(item.id)}
            variant="ghost"
            size="sm"
            className="text-red-600 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-950"
          >
            <Trash2 className="w-4 h-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
