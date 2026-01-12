import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface GeneralState {
  name: string;
  description: string;
  entityType: string;
}

interface GeneralSectionProps {
  state: GeneralState;
  onUpdate: (state: GeneralState) => void;
}

export function GeneralSection({ state, onUpdate }: GeneralSectionProps) {
  const handleChange = (field: keyof GeneralState, value: string) => {
    onUpdate({ ...state, [field]: value });
  };

  return (
    <div className="space-y-6">
      <Card className="p-6 border border-border bg-card">
        <h2 className="text-lg font-semibold text-foreground mb-6">
          Entity Information
        </h2>

        <div className="space-y-6">
          {/* Entity Name */}
          <div>
            <Label
              htmlFor="entity-name"
              className="text-sm font-medium text-foreground mb-2 block"
            >
              Entity Name <span className="text-red-500">*</span>
            </Label>
            <Input
              id="entity-name"
              placeholder="e.g., Customer, Order, Product"
              value={state.name}
              onChange={(e) => handleChange("name", e.target.value)}
              className="bg-input border border-border text-foreground placeholder:text-muted-foreground"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Unique identifier for this entity in the system
            </p>
          </div>

          {/* Entity Type */}
          <div>
            <Label
              htmlFor="entity-type"
              className="text-sm font-medium text-foreground mb-2 block"
            >
              Entity Type <span className="text-red-500">*</span>
            </Label>
            <Select
              value={state.entityType}
              onValueChange={(value) => handleChange("entityType", value)}
            >
              <SelectTrigger className="bg-input border border-border text-foreground">
                <SelectValue placeholder="Select entity type" />
              </SelectTrigger>
              <SelectContent className="bg-card border border-border">
                <SelectItem value="table">Table</SelectItem>
                <SelectItem value="view">View</SelectItem>
                <SelectItem value="materialized_view">
                  Materialized View
                </SelectItem>
                <SelectItem value="virtual">Virtual Entity</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* Description */}
          <div>
            <Label
              htmlFor="entity-description"
              className="text-sm font-medium text-foreground mb-2 block"
            >
              Description
            </Label>
            <Textarea
              id="entity-description"
              placeholder="Describe the purpose and usage of this entity"
              value={state.description}
              onChange={(e) => handleChange("description", e.target.value)}
              rows={4}
              className="bg-input border border-border text-foreground placeholder:text-muted-foreground resize-none"
            />
          </div>
        </div>
      </Card>
    </div>
  );
}
