import { generateAlias } from "../components/helper";

export interface SupportingTableItem {
  id: string;
  sourceTable: string;
  alias: string;
  factColumn: string;
  dimColumn: string;
}

export interface TargetMappingItem {
  id: string;
  targetColumn: string;
  mappingType: "key" | "string" | "kolom_sumber";
  sourceColumn?: string;
  transformation?: string;
}

export interface KeyMap {
  strategy: "shared_key" | "natural";
  targetKeyColumn: string;
}

export interface EditEntityFormState {
  primarySourceTable: string;
  supportingTables: SupportingTableItem[];
  primaryTargetTable: string;
  keyMap: KeyMap | null;
  targetMappings: TargetMappingItem[];
}

interface MappingSource {
  alias: string;
  from: string;
  join?: {
    fact_column: string;
    dim_column: string;
  };
}

interface MappingColumn {
  from: string;
  cast?: "string" | "number" | "boolean" | "int";
  expr?: "now" | "uuid" | "random_int";
}

interface MappingEntityPayload {
  entity: string;
  sources: MappingSource[];
  target_table: string;
  key: {
    strategy: "shared_key" | "natural";
    source: string;
    resolver?: {
      type: string;
      table: string;
    };
  };
  columns: Record<string, MappingColumn>;
  split_table?: Array<unknown>;
  routing?: unknown;
}

const randomId = () => Math.random().toString(36).slice(2, 11);

const buildKeymapTableName = (entityName: string) =>
  `_keymap_${entityName}`
    .replace(/\s+/g, "-")
    .replace(/[^a-zA-Z0-9_-]/g, "")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "");

const mapTransformationToJson = (
  transformation?: string,
): Partial<MappingColumn> => {
  if (!transformation) return {};

  if (transformation === "int to string") {
    return { cast: "string" };
  }

  if (transformation === "string to int") {
    return { cast: "int" };
  }

  if (transformation === "random key") {
    return { expr: "random_int" };
  }

  return {};
};

const mapJsonToTransformation = (column: MappingColumn): string | undefined => {
  if (column.expr === "random_int") return "random key";
  if (column.cast === "string") return "int to string";
  if (column.cast === "int") return "string to int";
  return undefined;
};

// json to state / interface
export const mappingJsonToFormState = (
  entity?: Partial<MappingEntityPayload> | null,
): EditEntityFormState => {
  if (!entity) {
    return {
      primarySourceTable: "",
      supportingTables: [],
      primaryTargetTable: "",
      keyMap: null,
      targetMappings: [],
    };
  }

  const sources = Array.isArray(entity.sources) ? entity.sources : [];
  const [primary, ...supporting] = sources;

  const targetMappings = Object.entries(entity.columns ?? {}).map(
    ([targetColumn, config]) => ({
      id: randomId(),
      targetColumn,
      mappingType:
        config.from === "$key" ? ("key" as const) : ("kolom_sumber" as const),
      sourceColumn: config.from === "$key" ? undefined : config.from,
      transformation: mapJsonToTransformation(config),
    }),
  );

  return {
    primarySourceTable: primary?.from ?? "",
    supportingTables: supporting.map((item) => ({
      id: randomId(),
      sourceTable: item.from,
      alias: item.alias || generateAlias(item.from),
      factColumn: item.join?.fact_column ?? "",
      dimColumn: item.join?.dim_column ?? "",
    })),
    primaryTargetTable: entity.target_table ?? "",
    keyMap: entity.key
      ? {
          strategy: entity.key.strategy,
          targetKeyColumn: entity.key.source,
        }
      : null,
    targetMappings,
  };
};

// state to interface
export const formStateToMappingJson = (params: {
  entityName: string;
  formState: EditEntityFormState;
  existingEntity?: Partial<MappingEntityPayload> | null;
}): MappingEntityPayload => {
  const { entityName, formState, existingEntity } = params;

  const sources: MappingSource[] = [];

  if (formState.primarySourceTable) {
    sources.push({
      alias: generateAlias(formState.primarySourceTable),
      from: formState.primarySourceTable,
    });
  }

  formState.supportingTables
    .filter((table) => table.sourceTable)
    .forEach((table) => {
      sources.push({
        alias: table.alias || generateAlias(table.sourceTable),
        from: table.sourceTable,
        join: {
          fact_column: table.factColumn,
          dim_column: table.dimColumn,
        },
      });
    });

  const columns = formState.targetMappings.reduce<
    Record<string, MappingColumn>
  >((acc, item) => {
    if (!item.targetColumn) return acc;

    if (item.mappingType === "key") {
      acc[item.targetColumn] = { from: "$key" };
      return acc;
    }

    if (!item.sourceColumn) return acc;

    acc[item.targetColumn] = {
      from: item.sourceColumn,
      ...mapTransformationToJson(item.transformation),
    };

    return acc;
  }, {});

  return {
    entity: entityName,
    sources,
    target_table: formState.primaryTargetTable,
    key: {
      strategy: formState.keyMap?.strategy ?? "natural",
      source: formState.keyMap?.targetKeyColumn ?? "",
      ...(formState.keyMap?.strategy === "shared_key"
        ? {
            resolver: {
              type: "mapping_table",
              table: buildKeymapTableName(entityName),
            },
          }
        : {}),
    },
    columns,
    split_table: existingEntity?.split_table,
    routing: existingEntity?.routing,
  };
};
