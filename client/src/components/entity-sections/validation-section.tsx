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

interface ValidationRule {
  id: string;
  field: string;
  rule: string;
  message: string;
}

interface ValidationState {
  rules: ValidationRule[];
}

interface ValidationSectionProps {
  state: ValidationState;
  onUpdate: (state: ValidationState) => void;
}

export function ValidationSection({ state, onUpdate }: ValidationSectionProps) {
  const [newField, setNewField] = useState("");
  const [newRule, setNewRule] = useState("");
  const [newMessage, setNewMessage] = useState("");

  const addRule = () => {
    if (newField && newRule) {
      const newValidation: ValidationRule = {
        id: `rule-${Date.now()}`,
        field: newField,
        rule: newRule,
        message: newMessage,
      };
      onUpdate({ ...state, rules: [...state.rules, newValidation] });
      setNewField("");
      setNewRule("");
      setNewMessage("");
    }
  };

  const removeRule = (id: string) => {
    onUpdate({ ...state, rules: state.rules.filter((r) => r.id !== id) });
  };

  return (
    <div className="space-y-6">
      <Card className="p-6 border border-border bg-card">
        <h2 className="text-lg font-semibold text-foreground mb-6">
          Validation Rules
        </h2>

        {/* Rules Table */}
        {state.rules.length > 0 && (
          <div className="mb-6 overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-b border-border hover:bg-transparent">
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Field
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Rule
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Error Message
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground w-8"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {state.rules.map((rule) => (
                  <TableRow
                    key={rule.id}
                    className="border-b border-border hover:bg-muted/50 transition-colors"
                  >
                    <TableCell className="text-sm text-foreground font-medium">
                      {rule.field}
                    </TableCell>
                    <TableCell className="text-sm text-foreground">
                      <span className="px-2 py-1 bg-muted rounded text-xs">
                        {rule.rule}
                      </span>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground truncate max-w-xs">
                      {rule.message}
                    </TableCell>
                    <TableCell>
                      <button
                        onClick={() => removeRule(rule.id)}
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

        {/* Add Rule Form */}
        <div className="space-y-4 p-4 bg-muted/50 rounded-lg border border-border">
          <h3 className="text-sm font-medium text-foreground">Add New Rule</h3>
          <div className="space-y-3">
            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label className="text-xs font-medium text-foreground mb-1 block">
                  Field Name
                </Label>
                <Input
                  placeholder="e.g., email, age"
                  value={newField}
                  onChange={(e) => setNewField(e.target.value)}
                  className="bg-input border border-border text-foreground placeholder:text-muted-foreground text-sm"
                />
              </div>
              <div>
                <Label className="text-xs font-medium text-foreground mb-1 block">
                  Validation Rule
                </Label>
                <Select value={newRule} onValueChange={setNewRule}>
                  <SelectTrigger className="bg-input border border-border text-foreground text-sm">
                    <SelectValue placeholder="Select rule type" />
                  </SelectTrigger>
                  <SelectContent className="bg-card border border-border">
                    <SelectItem value="required">Required</SelectItem>
                    <SelectItem value="email">Email Format</SelectItem>
                    <SelectItem value="min_length">Min Length</SelectItem>
                    <SelectItem value="max_length">Max Length</SelectItem>
                    <SelectItem value="numeric">Numeric</SelectItem>
                    <SelectItem value="unique">Unique</SelectItem>
                    <SelectItem value="pattern">Pattern Match</SelectItem>
                    <SelectItem value="range">Range</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div>
              <Label className="text-xs font-medium text-foreground mb-1 block">
                Error Message
              </Label>
              <Input
                placeholder="e.g., Please provide a valid email address"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                className="bg-input border border-border text-foreground placeholder:text-muted-foreground text-sm"
              />
            </div>
          </div>
          <Button
            onClick={addRule}
            variant="outline"
            size="sm"
            className="w-full gap-2 bg-transparent"
          >
            <Plus className="w-4 h-4" />
            Add Rule
          </Button>
        </div>
      </Card>
    </div>
  );
}
