import { useState } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Plus, Trash2 } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

interface Mapping {
  id: string;
  sourceField: string;
  targetField: string;
  dataType: string;
}

interface ColumnsMappingState {
  mappings: Mapping[];
}

interface ColumnsMappingSectionProps {
  state: ColumnsMappingState;
  onUpdate: (state: ColumnsMappingState) => void;
}

export function ColumnsMappingSection({
  state,
  onUpdate,
}: ColumnsMappingSectionProps) {
  const [newSourceField, setNewSourceField] = useState("");
  const [newTargetField, setNewTargetField] = useState("");
  const [newDataType, setNewDataType] = useState("");

  const addMapping = () => {
    if (newSourceField && newTargetField && newDataType) {
      const newMapping: Mapping = {
        id: `mapping-${Date.now()}`,
        sourceField: newSourceField,
        targetField: newTargetField,
        dataType: newDataType,
      };
      onUpdate({ ...state, mappings: [...state.mappings, newMapping] });
      setNewSourceField("");
      setNewTargetField("");
      setNewDataType("");
    }
  };

  const removeMapping = (id: string) => {
    onUpdate({ ...state, mappings: state.mappings.filter((m) => m.id !== id) });
  };

  return (
    <div className="space-y-6">
      <Card className="p-6 border border-border bg-card">
        <h2 className="text-lg font-semibold text-foreground mb-6">
          Columns Mapping
        </h2>

        {/* Mappings Table */}
        {state.mappings.length > 0 && (
          <div className="mb-6 overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-b border-border hover:bg-transparent">
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Source Field
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Target Field
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Data Type
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground w-8"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {state.mappings.map((mapping) => (
                  <TableRow
                    key={mapping.id}
                    className="border-b border-border hover:bg-muted/50 transition-colors"
                  >
                    <TableCell className="text-sm text-foreground font-mono">
                      {mapping.sourceField}
                    </TableCell>
                    <TableCell className="text-sm text-foreground font-mono">
                      {mapping.targetField}
                    </TableCell>
                    <TableCell className="text-sm text-foreground">
                      <span className="px-2 py-1 bg-muted rounded text-xs">
                        {mapping.dataType}
                      </span>
                    </TableCell>
                    <TableCell>
                      <button
                        onClick={() => removeMapping(mapping.id)}
                        className="p-1 hover:bg-input rounded transition-colors"
                      >
                        <Trash2 className="w-3 h-3 text-destructive" />
                      </button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}

        {/* Add Mapping Form */}
        <div className="space-y-4 p-4 bg-muted/50 rounded-lg border border-border">
          <h3 className="text-sm font-medium text-foreground">
            Add New Mapping
          </h3>
          <div className="grid grid-cols-3 gap-3">
            <div>
              <Label className="text-xs font-medium text-foreground mb-1 block">
                Source Field
              </Label>
              <Input
                placeholder="e.g., user_email"
                value={newSourceField}
                onChange={(e) => setNewSourceField(e.target.value)}
                className="bg-input border border-border text-foreground placeholder:text-muted-foreground text-sm"
              />
            </div>
            <div>
              <Label className="text-xs font-medium text-foreground mb-1 block">
                Target Field
              </Label>
              <Input
                placeholder="e.g., email"
                value={newTargetField}
                onChange={(e) => setNewTargetField(e.target.value)}
                className="bg-input border border-border text-foreground placeholder:text-muted-foreground text-sm"
              />
            </div>
            <div>
              <Label className="text-xs font-medium text-foreground mb-1 block">
                Data Type
              </Label>
              <Select value={newDataType} onValueChange={setNewDataType}>
                <SelectTrigger className="bg-input border border-border text-foreground text-sm">
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent className="bg-card border border-border">
                  <SelectItem value="varchar">Varchar</SelectItem>
                  <SelectItem value="int">Integer</SelectItem>
                  <SelectItem value="bigint">BigInt</SelectItem>
                  <SelectItem value="decimal">Decimal</SelectItem>
                  <SelectItem value="boolean">Boolean</SelectItem>
                  <SelectItem value="datetime">DateTime</SelectItem>
                  <SelectItem value="json">JSON</SelectItem>
                  <SelectItem value="uuid">UUID</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <Button
            onClick={addMapping}
            variant="outline"
            size="sm"
            className="w-full gap-2 bg-transparent"
          >
            <Plus className="w-4 h-4" />
            Add Mapping
          </Button>
        </div>
      </Card>
    </div>
  );
}
