import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface KeyIdentityState {
  primaryKeyField: string;
  identityStrategy: string;
}

interface KeyIdentitySectionProps {
  state: KeyIdentityState;
  onUpdate: (state: KeyIdentityState) => void;
}

export function KeyIdentitySection({
  state,
  onUpdate,
}: KeyIdentitySectionProps) {
  const handleChange = (field: keyof KeyIdentityState, value: string) => {
    onUpdate({ ...state, [field]: value });
  };

  return (
    <div className="space-y-6">
      <Card className="p-6 border border-border bg-card">
        <h2 className="text-lg font-semibold text-foreground mb-6">
          Key & Identity Configuration
        </h2>

        <div className="space-y-6">
          {/* Primary Key Field */}
          <div>
            <Label
              htmlFor="primary-key"
              className="text-sm font-medium text-foreground mb-2 block"
            >
              Primary Key Field <span className="text-red-500">*</span>
            </Label>
            <Input
              id="primary-key"
              placeholder="e.g., id, user_id, entity_key"
              value={state.primaryKeyField}
              onChange={(e) => handleChange("primaryKeyField", e.target.value)}
              className="bg-input border border-border text-foreground placeholder:text-muted-foreground"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Column that uniquely identifies each row in this entity
            </p>
          </div>

          {/* Identity Strategy */}
          <div>
            <Label
              htmlFor="identity-strategy"
              className="text-sm font-medium text-foreground mb-2 block"
            >
              Identity Strategy <span className="text-red-500">*</span>
            </Label>
            <Select
              value={state.identityStrategy}
              onValueChange={(value) => handleChange("identityStrategy", value)}
            >
              <SelectTrigger className="bg-input border border-border text-foreground">
                <SelectValue placeholder="Select strategy" />
              </SelectTrigger>
              <SelectContent className="bg-card border border-border">
                <SelectItem value="auto_increment">Auto Increment</SelectItem>
                <SelectItem value="uuid">UUID</SelectItem>
                <SelectItem value="composite">Composite Key</SelectItem>
                <SelectItem value="natural">Natural Key</SelectItem>
                <SelectItem value="sequence">Database Sequence</SelectItem>
                <SelectItem value="none">None</SelectItem>
              </SelectContent>
            </Select>
            <p className="text-xs text-muted-foreground mt-1">
              How new identities are generated for this entity
            </p>
          </div>

          {/* Info Box */}
          <div className="p-4 bg-blue-500/10 border border-blue-500/20 rounded-lg">
            <p className="text-sm text-blue-200">
              <strong>Note:</strong> The primary key and identity strategy
              determine how rows are uniquely identified and created in the
              system.
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}
