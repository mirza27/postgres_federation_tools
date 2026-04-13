import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Trash2 } from "lucide-react";

interface TargetColumnMappingProps {
  item: {
    id: string;
    targetColumn: string;
    mappingType: "key" | "string" | "kolom_sumber";
    sourceColumn?: string;
    cast?: "int to string" | "string to int";
    expr?: "now" | "uuid" | "random_string";
    defaultValue?: string;
  };
  onUpdate: (id: string, updates: any) => void;
  onRemove: (id: string) => void;
  targetColumns: string[];
  sourceColumns: string[];
}

const CAST_TYPES = ["int to string", "string to int"] as const;
const EXPRESSION_TYPES = ["now", "uuid", "random_string"] as const;
const NO_CAST_VALUE = "__no_cast__";
const NO_EXPR_VALUE = "__no_expr__";

export function TargetColumnMapping({
  item,
  onUpdate,
  onRemove,
  targetColumns,
  sourceColumns,
}: TargetColumnMappingProps) {
  const showSourceColumnField = item.mappingType === "kolom_sumber";
  const showDefaultField = item.mappingType === "string";
  const showTransformationFields =
    item.mappingType === "kolom_sumber" || item.mappingType === "key";

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
                const shouldKeepTransformation =
                  value === "kolom_sumber" || value === "key";
                const shouldKeepDefault = value === "string";

                onUpdate(item.id, {
                  mappingType: value,
                  sourceColumn:
                    value === "kolom_sumber" ? item.sourceColumn : undefined,
                  cast: shouldKeepTransformation ? item.cast : undefined,
                  expr: shouldKeepTransformation ? item.expr : undefined,
                  defaultValue: shouldKeepDefault ? item.defaultValue : undefined,
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

        {/* Conditional Fields */}
        {(showSourceColumnField || showTransformationFields || showDefaultField) && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3 pt-2 border-t border-slate-200 dark:border-slate-700 mt-3">
            {showSourceColumnField && (
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
            )}

            {showDefaultField && (
              <div className="flex flex-col gap-1 md:col-span-2">
                <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
                  Nilai Default
                </label>
                <Input
                  value={item.defaultValue || ""}
                  onChange={(event) =>
                    onUpdate(item.id, { defaultValue: event.target.value })
                  }
                  placeholder="Masukkan nilai string default..."
                  className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600"
                />
              </div>
            )}

            {showTransformationFields && (
              <div className="grid grid-cols-1 gap-3 md:grid-cols-2 md:col-span-1">
                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
                    Casting
                  </label>
                  <Select
                    value={item.cast || NO_CAST_VALUE}
                    onValueChange={(value) =>
                      onUpdate(item.id, {
                        cast: value === NO_CAST_VALUE ? undefined : value,
                      })
                    }
                  >
                    <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
                      <SelectValue placeholder="Pilih casting..." />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value={NO_CAST_VALUE}>Tidak Ada</SelectItem>
                      {CAST_TYPES.map((cast) => (
                        <SelectItem key={cast} value={cast}>
                          {cast}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="flex flex-col gap-1">
                  <label className="text-xs font-semibold text-slate-600 dark:text-slate-400">
                    Expression
                  </label>
                  <Select
                    value={item.expr || NO_EXPR_VALUE}
                    onValueChange={(value) =>
                      onUpdate(item.id, {
                        expr: value === NO_EXPR_VALUE ? undefined : value,
                      })
                    }
                  >
                    <SelectTrigger className="bg-white dark:bg-slate-700 border-slate-200 dark:border-slate-600">
                      <SelectValue placeholder="Pilih expression..." />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value={NO_EXPR_VALUE}>Tidak Ada</SelectItem>
                      {EXPRESSION_TYPES.map((expr) => (
                        <SelectItem key={expr} value={expr}>
                          {expr}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
