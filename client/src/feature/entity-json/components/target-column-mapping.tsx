import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";

interface TargetColumnMappingProps {
  item: {
    id: string;
    targetColumn: string;
    mappingType: "key" | "string" | "kolom_sumber";
    sourceColumn?: string;
    transformation?: string;
  };
  onUpdate: (id: string, updates: any) => void;
  onRemove: (id: string) => void;
  targetColumns: string[];
  sourceColumns: string[];
}

const TRANSFORMATION_TYPES = ["int to string", "string to int", "random key"];
const NO_TRANSFORMATION_VALUE = "__none__";

export function TargetColumnMapping({
  item,
  onUpdate,
  onRemove,
  targetColumns,
  sourceColumns,
}: TargetColumnMappingProps) {
  const showSourceFields = item.mappingType === "kolom_sumber";

  return (
    <div className="bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg p-4">
      <div className="space-y-3">
        {/* First Row */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-3 items-end">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
              Kolom Target
            </label>
            <Select
              value={item.targetColumn}
              onValueChange={(value) =>
                onUpdate(item.id, { targetColumn: value })
              }
            >
              <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
                <SelectValue placeholder="Pilih kolom..." />
              </SelectTrigger>
              <SelectContent>
                {targetColumns.map((col) => (
                  <SelectItem key={col} value={col}>
                    {col}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
              Tipe Mapping
            </label>
            <Select
              value={item.mappingType}
              onValueChange={(value: any) => {
                onUpdate(item.id, {
                  mappingType: value,
                  sourceColumn:
                    value === "kolom_sumber" ? item.sourceColumn : undefined,
                  transformation:
                    value === "kolom_sumber" ? item.transformation : undefined,
                });
              }}
            >
              <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="key">Key</SelectItem>
                <SelectItem value="string">String (Nilai Langsung)</SelectItem>
                <SelectItem value="kolom_sumber">Kolom Sumber</SelectItem>
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
              <Trash2 className="w-4 h-4 mr-2" />
              Hapus
            </Button>
          </div>
        </div>

        {/* Conditional Fields - Show when mappingType is 'kolom_sumber' */}
        {showSourceFields && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3 pt-2 border-t border-slate-200 dark:border-slate-700 mt-3">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
                Kolom Sumber
              </label>
              <Select
                value={item.sourceColumn || ""}
                onValueChange={(value) =>
                  onUpdate(item.id, { sourceColumn: value })
                }
              >
                <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
                  <SelectValue placeholder="Pilih kolom..." />
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
                Transformasi
              </label>
              <Select
                value={item.transformation || NO_TRANSFORMATION_VALUE}
                onValueChange={(value) =>
                  onUpdate(item.id, {
                    transformation:
                      value === NO_TRANSFORMATION_VALUE ? undefined : value,
                  })
                }
              >
                <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
                  <SelectValue placeholder="Pilih transformasi..." />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value={NO_TRANSFORMATION_VALUE}>
                    Tidak Ada
                  </SelectItem>
                  {TRANSFORMATION_TYPES.map((t) => (
                    <SelectItem key={t} value={t}>
                      {t}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
