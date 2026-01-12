import { useState } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Plus, Trash2, GripVertical } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface Source {
  id: string;
  name: string;
  table: string;
}

interface Join {
  id: string;
  sourceId: string;
  condition: string;
}

interface SourcesJoinsState {
  sources: Source[];
  joins: Join[];
}

interface SourcesJoinsSectionProps {
  state: SourcesJoinsState;
  onUpdate: (state: SourcesJoinsState) => void;
}

export function SourcesJoinsSection({
  state,
  onUpdate,
}: SourcesJoinsSectionProps) {
  const [newSourceName, setNewSourceName] = useState("");
  const [newSourceTable, setNewSourceTable] = useState("");
  const [newJoinSourceId, setNewJoinSourceId] = useState("");
  const [newJoinCondition, setNewJoinCondition] = useState("");

  const addSource = () => {
    if (newSourceName && newSourceTable) {
      const newSource: Source = {
        id: `source-${Date.now()}`,
        name: newSourceName,
        table: newSourceTable,
      };
      onUpdate({ ...state, sources: [...state.sources, newSource] });
      setNewSourceName("");
      setNewSourceTable("");
    }
  };

  const removeSource = (id: string) => {
    onUpdate({
      ...state,
      sources: state.sources.filter((s) => s.id !== id),
      joins: state.joins.filter((j) => j.sourceId !== id),
    });
  };

  const addJoin = () => {
    if (newJoinSourceId && newJoinCondition) {
      const newJoin: Join = {
        id: `join-${Date.now()}`,
        sourceId: newJoinSourceId,
        condition: newJoinCondition,
      };
      onUpdate({ ...state, joins: [...state.joins, newJoin] });
      setNewJoinSourceId("");
      setNewJoinCondition("");
    }
  };

  const removeJoin = (id: string) => {
    onUpdate({ ...state, joins: state.joins.filter((j) => j.id !== id) });
  };

  const getSourceName = (id: string) => {
    return state.sources.find((s) => s.id === id)?.name || "Unknown";
  };

  return (
    <div className="space-y-6">
      {/* Sources */}
      <Card className="p-6 border border-border bg-card">
        <h2 className="text-lg font-semibold text-foreground mb-6">
          Data Sources
        </h2>

        <div className="space-y-4">
          {state.sources.map((source) => (
            <div
              key={source.id}
              className="flex items-center gap-3 p-3 bg-muted rounded-lg border border-border"
            >
              <GripVertical className="w-4 h-4 text-muted-foreground cursor-grab" />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-foreground truncate">
                  {source.name}
                </p>
                <p className="text-xs text-muted-foreground truncate">
                  {source.table}
                </p>
              </div>
              <button
                onClick={() => removeSource(source.id)}
                className="p-2 hover:bg-input rounded transition-colors"
              >
                <Trash2 className="w-4 h-4 text-destructive" />
              </button>
            </div>
          ))}

          <div className="space-y-3 pt-4 border-t border-border">
            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label className="text-xs font-medium text-foreground mb-1 block">
                  Source Name
                </Label>
                <Input
                  placeholder="e.g., Primary Source"
                  value={newSourceName}
                  onChange={(e) => setNewSourceName(e.target.value)}
                  className="bg-input border border-border text-foreground placeholder:text-muted-foreground"
                />
              </div>
              <div>
                <Label className="text-xs font-medium text-foreground mb-1 block">
                  Table Name
                </Label>
                <Input
                  placeholder="e.g., users"
                  value={newSourceTable}
                  onChange={(e) => setNewSourceTable(e.target.value)}
                  className="bg-input border border-border text-foreground placeholder:text-muted-foreground"
                />
              </div>
            </div>
            <Button
              onClick={addSource}
              variant="outline"
              size="sm"
              className="w-full gap-2 bg-transparent"
            >
              <Plus className="w-4 h-4" />
              Add Source
            </Button>
          </div>
        </div>
      </Card>

      {/* Joins */}
      <Card className="p-6 border border-border bg-card">
        <h2 className="text-lg font-semibold text-foreground mb-6">
          Join Conditions
        </h2>

        <div className="space-y-4">
          {state.joins.map((join) => (
            <div
              key={join.id}
              className="flex items-center gap-3 p-3 bg-muted rounded-lg border border-border"
            >
              <GripVertical className="w-4 h-4 text-muted-foreground cursor-grab" />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-foreground">
                  {getSourceName(join.sourceId)}
                </p>
                <p className="text-xs text-muted-foreground font-mono truncate">
                  {join.condition}
                </p>
              </div>
              <button
                onClick={() => removeJoin(join.id)}
                className="p-2 hover:bg-input rounded transition-colors"
              >
                <Trash2 className="w-4 h-4 text-destructive" />
              </button>
            </div>
          ))}

          <div className="space-y-3 pt-4 border-t border-border">
            <div>
              <Label className="text-xs font-medium text-foreground mb-1 block">
                Source
              </Label>
              <Select
                value={newJoinSourceId}
                onValueChange={setNewJoinSourceId}
              >
                <SelectTrigger className="bg-input border border-border text-foreground">
                  <SelectValue placeholder="Select source" />
                </SelectTrigger>
                <SelectContent className="bg-card border border-border">
                  {state.sources.map((source) => (
                    <SelectItem key={source.id} value={source.id}>
                      {source.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div>
              <Label className="text-xs font-medium text-foreground mb-1 block">
                Join Condition
              </Label>
              <Input
                placeholder="e.g., primary.id = source.user_id"
                value={newJoinCondition}
                onChange={(e) => setNewJoinCondition(e.target.value)}
                className="bg-input border border-border text-foreground placeholder:text-muted-foreground font-mono text-xs"
              />
            </div>
            <Button
              onClick={addJoin}
              variant="outline"
              size="sm"
              className="w-full gap-2 bg-transparent"
            >
              <Plus className="w-4 h-4" />
              Add Join
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}
