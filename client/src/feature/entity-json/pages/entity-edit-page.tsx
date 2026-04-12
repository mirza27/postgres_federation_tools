import { useMemo, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ChevronLeft, Plus } from "lucide-react";
import { TargetColumnMapping } from "../components/target-column-mapping";
import { SupportingTable } from "../components/supporting-table";
import { generateAlias } from "../components/helper";
import { useFetcher, useLoaderData, useNavigate } from "react-router-dom";
import { DefaultPaths } from "@/path";

interface SupportingTableItem {
  id: string;
  sourceTable: string;
  alias: string;
  factColumn: string;
  dimColumn: string;
}

interface TargetMappingItem {
  id: string;
  targetColumn: string;
  mappingType: "key" | "string" | "kolom_sumber";
  sourceColumn?: string;
  transformation?: string;
}

interface KeyMap {
  strategy: "shared_key" | "natural";
  targetKeyColumn: string;
}

interface SchemaColumn {
  column_name: string;
  data_type: string;
}

interface SchemaTable {
  table_name: string;
  columns: SchemaColumn[];
}

const toQualifiedColumnName = (alias: string, column: string) =>
  alias ? `${alias}.${column}` : column;

export function EditEntityPage() {
  const navigate = useNavigate();
  const loaderData = useLoaderData() as any;

  const entity = loaderData.data?.entity;
  const fileName = loaderData.data?.filename ?? "Unavailable";
  const entityName = entity?.entity ?? "Unknown Entity";
  const sourceSchema: SchemaTable[] = useMemo(
    () => loaderData.data?.schemas.source_schema ?? [],
    [loaderData],
  );
  const targetSchema: SchemaTable[] = useMemo(
    () => loaderData.data?.schemas.target_schema ?? [],
    [loaderData],
  );

  const [primarySourceTable, setPrimarySourceTable] = useState("");
  const [supportingTables, setSupportingTables] = useState<
    SupportingTableItem[]
  >([]);
  const [primaryTargetTable, setPrimaryTargetTable] = useState("");
  const [keyMap, setKeyMap] = useState<KeyMap | null>(null);

  const [targetMappings, setTargetMappings] = useState<TargetMappingItem[]>([]);

  // update
  const updateFetcher = useFetcher();

  const sourceColumnsMap = useMemo(
    () =>
      sourceSchema.reduce<Record<string, string[]>>((acc, table) => {
        acc[table.table_name] = table.columns.map((col) => col.column_name);
        return acc;
      }, {}),
    [sourceSchema],
  );

  const targetColumnsMap = useMemo(
    () =>
      targetSchema.reduce<Record<string, string[]>>((acc, table) => {
        acc[table.table_name] = table.columns.map((col) => col.column_name);
        return acc;
      }, {}),
    [targetSchema],
  );

  // list semua tabel dari source schema kecuali yang sudah dipilih sebagai primary source table
  const availableSupportingTables = useMemo(
    () =>
      sourceSchema
        .map((table) => table.table_name)
        .filter((tableName) => tableName !== primarySourceTable),
    [sourceSchema, primarySourceTable],
  );

  // list column dari primary source table
  const primarySourceColumns = useMemo(
    () => sourceColumnsMap[primarySourceTable] ?? [],
    [primarySourceTable, sourceColumnsMap],
  );

  const primaryTargetColumns = useMemo(
    () => targetColumnsMap[primaryTargetTable] ?? [],
    [primaryTargetTable, targetColumnsMap],
  );

  const primarySourceAlias = useMemo(
    () => generateAlias(primarySourceTable),
    [primarySourceTable],
  );

  const primaryTargetAlias = useMemo(
    () => generateAlias(primaryTargetTable),
    [primaryTargetTable],
  );

  // tambah join table
  const addSupportingTable = () => {
    const newItem: SupportingTableItem = {
      id: Math.random().toString(36).substr(2, 9),
      sourceTable: "",
      alias: "",
      factColumn: "",
      dimColumn: "",
    };
    setSupportingTables([...supportingTables, newItem]);
  };

  // hapus join table
  const removeSupportingTable = (id: string) => {
    setSupportingTables(supportingTables.filter((item) => item.id !== id));
  };

  // update join table
  const updateSupportingTable = (
    id: string,
    updates: Partial<SupportingTableItem>,
  ) => {
    setSupportingTables(
      supportingTables.map((item) =>
        item.id === id ? { ...item, ...updates } : item,
      ),
    );
  };

  // list kolom for source column mapping (primary source columns + all supporting table columns)
  const allSourceColumnsForMapping = useMemo(() => {
    const columns = primarySourceColumns.map((column) =>
      toQualifiedColumnName(primarySourceAlias, column),
    );

    supportingTables.forEach((table) => {
      if (table.sourceTable) {
        const tableColumns = sourceColumnsMap[table.sourceTable] || [];
        const tableAlias = table.alias || generateAlias(table.sourceTable);

        columns.push(
          ...tableColumns.map((column) =>
            toQualifiedColumnName(tableAlias, column),
          ),
        );
      }
    });

    return Array.from(new Set(columns));
  }, [
    primarySourceAlias,
    primarySourceColumns,
    supportingTables,
    sourceColumnsMap,
  ]);

  const qualifiedPrimaryTargetColumns = useMemo(
    () =>
      primaryTargetColumns.map((column) =>
        toQualifiedColumnName(primaryTargetAlias, column),
      ),
    [primaryTargetAlias, primaryTargetColumns],
  );

  const addTargetMapping = () => {
    const newMapping: TargetMappingItem = {
      id: Math.random().toString(36).substr(2, 9),
      targetColumn: "",
      mappingType: "string",
    };
    setTargetMappings([...targetMappings, newMapping]);
  };

  const removeTargetMapping = (id: string) => {
    setTargetMappings(targetMappings.filter((item) => item.id !== id));
  };

  const updateTargetMapping = (
    id: string,
    updates: Partial<TargetMappingItem>,
  ) => {
    setTargetMappings(
      targetMappings.map((item) =>
        item.id === id ? { ...item, ...updates } : item,
      ),
    );
  };

  const getTargetColumns = (tableName: string) => {
    for (const table of targetSchema) {
      if (table.table_name === tableName) {
        return table.columns.map((col) => col.column_name);
      }
    }

    return [];
  };

  // ketika ganti primary source table
  const handlePrimarySourceTableChange = (tableName: string) => {
    setPrimarySourceTable(tableName);

    setSupportingTables([]); // delete semua support / join table
  };

  return (
    <div className="h-screen w-full flex flex-col bg-background">
      {/* Header */}
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
          {/* Action Buttons */}
          <div className="ml-auto flex gap-4">
            <Button
              type="submit"
              variant="outline"
              className="border-slate-200 dark:border-slate-700"
            >
              Reset
            </Button>
            <Button
              className="text-white gap-2 bg-sidebar-primary hover:bg-sidebar-primary/90"
              type="submit"
              disabled={updateFetcher.state === "submitting"}
            >
              Simpan Konfigurasi
            </Button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* sisi kiri =========================*/}
          <div className="lg:col-span-2 space-y-6">
            {/* Pilih Tabel */}
            <Card className="shadow-sm border-slate-200 dark:border-slate-700">
              <CardHeader className="border-b border-slate-200 dark:border-slate-700p">
                <CardTitle className="text-xl">Database Sumber</CardTitle>
                <CardDescription>
                  Tambahkan tabel sumber utama dan tambahkan satu atau lebih
                  tabel pendukung untuk operasi join / foreign key
                </CardDescription>
              </CardHeader>

              <CardContent className="">
                {/* tabel sumber utama */}
                <div className="grid grid-cols-1 md:grid-cols-3 items-center gap-4 mb-4">
                  <label className="text-sm font-semibold text-slate-700 dark:text-slate-300 md:col-span-1">
                    Pilih Tabel Sumber Utama
                  </label>
                  <div className="md:col-span-2">
                    <Select
                      value={primarySourceTable}
                      onValueChange={handlePrimarySourceTableChange}
                    >
                      <SelectTrigger className="w-full bg-white dark:bg-slate-800">
                        <SelectValue placeholder="Pilih tabel..." />
                      </SelectTrigger>
                      <SelectContent>
                        {sourceSchema.map((table: any) => (
                          <SelectItem
                            key={table.table_name}
                            value={table.table_name}
                          >
                            {table.table_name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                {/* tabel sumber pendukung */}
                <div className="space-y-4 border-t m-x-4 pt-4">
                  {supportingTables.length === 0 ? (
                    <p className="text-sm text-slate-500 dark:text-slate-400">
                      Belum ada tabel pendukung
                    </p>
                  ) : (
                    supportingTables.map((item) => (
                      <SupportingTable
                        key={item.id}
                        item={item}
                        onUpdate={updateSupportingTable}
                        onRemove={removeSupportingTable}
                        allTables={availableSupportingTables}
                        sourceColumns={primarySourceColumns}
                        allColumnsMap={sourceColumnsMap}
                      />
                    ))
                  )}
                  <Button
                    onClick={addSupportingTable}
                    variant="outline"
                    className="w-full border-slate-200 dark:border-slate-700"
                  >
                    <Plus className="w-4 h-4 mr-2" />
                    Tambah Tabel Pendukung
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* 2. Target Database */}
            <Card className="shadow-sm border-slate-200 dark:border-slate-700">
              <CardHeader className="border-b border-slate-200 dark:border-slate-700">
                <CardTitle className="text-xl">Database Target</CardTitle>
                <CardDescription>
                  Pilih tabel target untuk pemetaan data
                </CardDescription>
              </CardHeader>
              <CardContent className="">
                <div className="grid grid-cols-1 md:grid-cols-3 items-center mb-4">
                  <label className="text-sm font-semibold text-slate-700 dark:text-slate-300 md:col-span-1">
                    Pilih Tabel Target
                  </label>
                  <div className="md:col-span-2">
                    <Select
                      value={primaryTargetTable}
                      onValueChange={setPrimaryTargetTable}
                    >
                      <SelectTrigger className="w-full bg-white dark:bg-slate-800">
                        <SelectValue placeholder="Pilih tabel..." />
                      </SelectTrigger>
                      <SelectContent>
                        {targetSchema.map((table: any) => (
                          <SelectItem
                            key={table.table_name}
                            value={table.table_name}
                          >
                            {table.table_name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <div className="border-t py-4">
                  <CardDescription className="pb-4">
                    Pilih aturan key ke tabel database target untuk memastikan
                    relasi yang benar antara sumber dan target
                  </CardDescription>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="flex flex-col gap-2">
                      <label className="text-sm font-semibold text-slate-700 dark:text-slate-300">
                        Tipe Key
                      </label>
                      <Select
                        value={keyMap?.strategy || ""}
                        onValueChange={(value: any) =>
                          setKeyMap({
                            ...keyMap,
                            strategy: value ?? "",
                          } as KeyMap)
                        }
                      >
                        <SelectTrigger className="w-full bg-white dark:bg-slate-800">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="shared_key">Shared Key</SelectItem>
                          <SelectItem value="natural">Natural</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex flex-col gap-2">
                      <label className="text-sm font-semibold text-slate-700 dark:text-slate-300">
                        Kolom Key dari Tabel Target
                      </label>
                      <Select
                        value={keyMap?.targetKeyColumn || ""}
                        onValueChange={(value: any) =>
                          setKeyMap({
                            ...keyMap,
                            targetKeyColumn: value ?? "",
                          } as KeyMap)
                        }
                      >
                        <SelectTrigger className="w-full bg-white dark:bg-slate-800">
                          <SelectValue placeholder="Pilih kolom..." />
                        </SelectTrigger>
                        <SelectContent>
                          {qualifiedPrimaryTargetColumns.map((col) => (
                            <SelectItem key={col} value={col}>
                              {col}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* 3. Target Column Mapping */}
            <Card className="shadow-sm border-slate-200 dark:border-slate-700">
              <CardHeader className="border-b border-slate-200 dark:border-slate-700">
                <CardTitle className="text-xl">Pemetaan Kolom</CardTitle>
                <CardDescription>
                  Tentukan bagaimana setiap kolom target dipetakan dari sumber
                </CardDescription>
              </CardHeader>
              <CardContent className="">
                <div className="space-y-4">
                  {targetMappings.length === 0 ? (
                    <p className="text-sm text-slate-500 dark:text-slate-400">
                      Belum ada pemetaan kolom
                    </p>
                  ) : (
                    targetMappings.map((item) => (
                      <TargetColumnMapping
                        key={item.id}
                        item={item}
                        onUpdate={updateTargetMapping}
                        onRemove={removeTargetMapping}
                        targetColumns={getTargetColumns(primaryTargetTable)}
                        sourceColumns={allSourceColumnsForMapping}
                      />
                    ))
                  )}
                  <Button
                    onClick={addTargetMapping}
                    variant="outline"
                    className="w-full border-slate-200 dark:border-slate-700"
                  >
                    <Plus className="w-4 h-4 mr-2" />
                    Tambah Mapping
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* SISI KANAN ====================================== */}
          <div className="lg:col-span-1">
            <Card className="sticky top-6">
              <CardHeader>
                <CardTitle>Informasi</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4 text-sm text-slate-600">
                <p>Pilih tabel sumber utama sebagai basis data utama.</p>

                <p>Tabel pendukung digunakan untuk enrich data (join).</p>

                <p>Pastikan kolom yang dipilih memiliki relasi yang sesuai</p>
                <p>File ini akan dikonversi menjadi file {fileName}.</p>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
