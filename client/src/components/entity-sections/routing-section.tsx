import { Card } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";

interface RoutingState {
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
}

interface RoutingSectionProps {
  state: RoutingState;
  onUpdate: (state: RoutingState) => void;
}

export function RoutingSection({ state, onUpdate }: RoutingSectionProps) {
  const handleActionChange = (event: string, action: string) => {
    const updatedState = { ...state };
    updatedState[event as keyof RoutingState].action = action;
    updatedState[event as keyof RoutingState].matchKeyColumns =
      action === "Ignore"
        ? []
        : updatedState[event as keyof RoutingState].matchKeyColumns;
    onUpdate(updatedState);
  };

  const handleKeyFormatChange = (event: string, value: string) => {
    const updatedState = { ...state };
    updatedState[event as keyof RoutingState].keyFormat = value;
    onUpdate(updatedState);
  };

  const targetColumns = [
    "user_id",
    "order_id",
    "product_id",
    "timestamp",
    "status",
  ];

  return (
    <div className="space-y-6 w-full">
      <Card className="p-6 border border-border bg-card">
        <h3 className="text-base font-semibold text-foreground mb-4">
          On Snapshot
        </h3>
        <div className="space-y-4">
          <div className="flex items-center gap-6">
            <Label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="on-snapshot"
                value="Ignore"
                checked={state.onSnapshot.action === "Ignore"}
                onChange={(e) =>
                  handleActionChange("onSnapshot", e.target.value)
                }
                className="w-4 h-4"
              />
              <span className="text-sm text-foreground">Ignore</span>
            </Label>
            <Label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="on-snapshot"
                value="Insert"
                checked={state.onSnapshot.action === "Insert"}
                onChange={(e) =>
                  handleActionChange("onSnapshot", e.target.value)
                }
                className="w-4 h-4"
              />
              <span className="text-sm text-foreground">Insert</span>
            </Label>
          </div>
          <div className="pt-2">
            <Label
              htmlFor="snapshot-key-format"
              className="text-sm font-medium text-foreground mb-2 block"
            >
              Key Format String
            </Label>
            <Input
              id="snapshot-key-format"
              type="text"
              placeholder="e.g., user_{id}, order-{timestamp}"
              value={state.onSnapshot.keyFormat}
              onChange={(e) =>
                handleKeyFormatChange("onSnapshot", e.target.value)
              }
              className="bg-input border border-border text-foreground placeholder:text-muted-foreground"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Define the naming format for the key column
            </p>
          </div>
        </div>
      </Card>

      <Card className="p-6 border border-border bg-card">
        <h3 className="text-base font-semibold text-foreground mb-4">
          On Create
        </h3>
        <div className="space-y-4">
          <div className="flex items-center gap-6">
            <Label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="on-create"
                value="Ignore"
                checked={state.onCreate.action === "Ignore"}
                onChange={(e) => handleActionChange("onCreate", e.target.value)}
                className="w-4 h-4"
              />
              <span className="text-sm text-foreground">Ignore</span>
            </Label>
            <Label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="on-create"
                value="Insert"
                checked={state.onCreate.action === "Insert"}
                onChange={(e) => handleActionChange("onCreate", e.target.value)}
                className="w-4 h-4"
              />
              <span className="text-sm text-foreground">Insert</span>
            </Label>
          </div>
          <div className="pt-2">
            <Label
              htmlFor="create-key-format"
              className="text-sm font-medium text-foreground mb-2 block"
            >
              Key Format String
            </Label>
            <Input
              id="create-key-format"
              type="text"
              placeholder="e.g., user_{id}, order-{timestamp}"
              value={state.onCreate.keyFormat}
              onChange={(e) =>
                handleKeyFormatChange("onCreate", e.target.value)
              }
              className="bg-input border border-border text-foreground placeholder:text-muted-foreground"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Define the naming format for the key column
            </p>
          </div>
        </div>
      </Card>

      {/* On Update */}
      <Card className="p-6 border border-border bg-card">
        <h3 className="text-base font-semibold text-foreground mb-4">
          On Update
        </h3>
        <div className="space-y-4">
          <div className="flex items-center gap-6">
            <Label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="on-update"
                value="Ignore"
                checked={state.onUpdate.action === "Ignore"}
                onChange={(e) => handleActionChange("onUpdate", e.target.value)}
                className="w-4 h-4"
              />
              <span className="text-sm text-foreground">Ignore</span>
            </Label>
            <Label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="on-update"
                value="Update"
                checked={state.onUpdate.action === "Update"}
                onChange={(e) => handleActionChange("onUpdate", e.target.value)}
                className="w-4 h-4"
              />
              <span className="text-sm text-foreground">Update</span>
            </Label>
            <Label className="flex items-center gap-3 cursor-pointer">
              <input
                type="radio"
                name="on-update"
                value="Upsert"
                checked={state.onUpdate.action === "Upsert"}
                onChange={(e) => handleActionChange("onUpdate", e.target.value)}
                className="w-4 h-4"
              />
              <span className="text-sm text-foreground">Upsert</span>
            </Label>
          </div>

          {/* Conditional Match Key Columns */}
          {(state.onUpdate.action === "Update" ||
            state.onUpdate.action === "Upsert") && (
            <div className="mt-4 p-4 bg-muted/50 rounded-lg border border-border">
              <Label
                htmlFor="match-columns"
                className="text-sm font-medium text-foreground mb-2 block"
              >
                Match Key Columns <span className="text-accent">*</span>
              </Label>
              <div className="space-y-2">
                {targetColumns.map((col) => (
                  <Label
                    key={col}
                    className="flex items-center gap-3 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      checked={state.onUpdate.matchKeyColumns.includes(col)}
                      onChange={(e) => {
                        const updated = { ...state };
                        if (e.target.checked) {
                          updated.onUpdate.matchKeyColumns = [
                            ...updated.onUpdate.matchKeyColumns,
                            col,
                          ];
                        } else {
                          updated.onUpdate.matchKeyColumns =
                            updated.onUpdate.matchKeyColumns.filter(
                              (c) => c !== col
                            );
                        }
                        onUpdate(updated);
                      }}
                      className="w-4 h-4"
                    />
                    <span className="text-sm text-foreground/80">{col}</span>
                  </Label>
                ))}
              </div>
              <p className="text-xs text-muted-foreground mt-2">
                Typically use the column marked as Special Type = Key
              </p>
            </div>
          )}

          {/* Key Format Input for On Update */}
          <div className="pt-2">
            <Label
              htmlFor="update-key-format"
              className="text-sm font-medium text-foreground mb-2 block"
            >
              Key Format String
            </Label>
            <Input
              id="update-key-format"
              type="text"
              placeholder="e.g., user_{id}, order-{timestamp}"
              value={state.onUpdate.keyFormat}
              onChange={(e) =>
                handleKeyFormatChange("onUpdate", e.target.value)
              }
              className="bg-input border border-border text-foreground placeholder:text-muted-foreground"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Define the naming format for the key column
            </p>
          </div>
        </div>
      </Card>

      {/* Info Box */}
      <div className="p-4 bg-sidebar-primary/10 border border-sidebar-primary/20 rounded-lg">
        <p className="text-sm text-sidebar-primary/90">
          <strong>Event Handlers:</strong> Define how to handle data on
          different events. Select the appropriate action, configure key format
          strings, and set matching columns for updates.
        </p>
      </div>
    </div>
  );
}
