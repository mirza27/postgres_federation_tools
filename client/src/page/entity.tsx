import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { ChevronLeft, Save, X } from "lucide-react";
import { GeneralSection } from "../components/entity-sections/general-section.tsx";
import { SourcesJoinsSection } from "../components/entity-sections/sources-join-section.tsx";
import { KeyIdentitySection } from "../components/entity-sections/key-identity-section.tsx";
import { ColumnsMappingSection } from "../components/entity-sections/columns-mapping-section.tsx";
import { RoutingSection } from "../components/entity-sections/routing-section.tsx";

interface EntityState {
  general: {
    name: string;
    description: string;
    entityType: string;
  };
  sourcesJoins: {
    sources: Array<{ id: string; name: string; table: string }>;
    joins: Array<{ id: string; sourceId: string; condition: string }>;
  };
  keyIdentity: {
    primaryKeyField: string;
    identityStrategy: string;
  };
  columnsMapping: {
    mappings: Array<{
      id: string;
      sourceField: string;
      targetField: string;
      dataType: string;
    }>;
  };
  routing: {
    onSnapshot: {
      action: string;
      matchKeyColumns: string[];
      keyFormat: string;
    };
    onCreate: {
      action: string;
      matchKeyColumns: string[];
      keyFormat: string;
    };
    onUpdate: {
      action: string;
      matchKeyColumns: string[];
      keyFormat: string;
    };
  };
}

export function EntityEditor() {
  const [activeTab, setActiveTab] = useState("general");
  const [entityState, setEntityState] = useState<EntityState>({
    general: {
      name: "",
      description: "",
      entityType: "",
    },
    sourcesJoins: {
      sources: [],
      joins: [],
    },
    keyIdentity: {
      primaryKeyField: "",
      identityStrategy: "",
    },
    columnsMapping: {
      mappings: [],
    },
    routing: {
      onSnapshot: {
        action: "",
        matchKeyColumns: [],
        keyFormat: "",
      },
      onCreate: {
        action: "",
        matchKeyColumns: [],
        keyFormat: "",
      },
      onUpdate: {
        action: "",
        matchKeyColumns: [],
        keyFormat: "",
      },
    },
  });

  const updateState = (section: keyof EntityState, data: any) => {
    setEntityState((prev) => ({
      ...prev,
      [section]: data,
    }));
  };

  return (
    <div className="h-screen w-full flex flex-col bg-background">
      {/* Header */}
      <div className="border-b border-border px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button className="p-1 hover:bg-muted rounded-md transition-colors">
              <ChevronLeft className="w-5 h-5 text-foreground" />
            </button>
            <div>
              <h1 className="text-2xl font-semibold text-foreground">
                {entityState.general.name || "Entity Editor"}
              </h1>
              <p className="text-sm text-muted-foreground">
                Configure database entity mapping
              </p>
            </div>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              className="gap-2 bg-transparent"
            >
              <X className="w-4 h-4" />
              Cancel
            </Button>
            <Button
              size="sm"
              className="gap-2 bg-sidebar-primary hover:bg-sidebar-primary/90"
            >
              <Save className="w-4 h-4" />
              Save Entity
            </Button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto">
        <div className="w-full mx-auto p-6">
          <Tabs
            value={activeTab}
            onValueChange={setActiveTab}
            className="w-full"
          >
            <TabsList className="grid w-full grid-cols-5 bg-muted p-1 rounded-lg">
              <TabsTrigger value="general" className="text-sm">
                General
              </TabsTrigger>
              <TabsTrigger value="sources" className="text-sm">
                Sources & Joins
              </TabsTrigger>
              <TabsTrigger value="key" className="text-sm">
                Key/Identity
              </TabsTrigger>
              <TabsTrigger value="columns" className="text-sm">
                Columns Mapping
              </TabsTrigger>
              <TabsTrigger value="routing" className="text-sm">
                Routing
              </TabsTrigger>
            </TabsList>

            {/* Tab Content */}
            <div className="mt-6 w-full">
              <TabsContent value="general" className="w-full">
                <GeneralSection
                  state={entityState.general}
                  onUpdate={(data) => updateState("general", data)}
                />
              </TabsContent>

              <TabsContent value="sources" className="w-full">
                <SourcesJoinsSection
                  state={entityState.sourcesJoins}
                  onUpdate={(data) => updateState("sourcesJoins", data)}
                />
              </TabsContent>

              <TabsContent value="key" className="w-full">
                <KeyIdentitySection
                  state={entityState.keyIdentity}
                  onUpdate={(data) => updateState("keyIdentity", data)}
                />
              </TabsContent>

              <TabsContent value="columns" className="w-full">
                <ColumnsMappingSection
                  state={entityState.columnsMapping}
                  onUpdate={(data) => updateState("columnsMapping", data)}
                />
              </TabsContent>

              <TabsContent value="routing" className="w-full">
                <RoutingSection
                  state={entityState.routing}
                  onUpdate={(data) => updateState("routing", data)}
                />
              </TabsContent>
            </div>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
