import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useFetcher } from "react-router-dom";
import { DefaultPaths } from "@/path";
import { toast } from "sonner";

export function NewEntityForm() {
  const createFetcher = useFetcher();
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (createFetcher.data) {
      if (createFetcher.data.ok) {
        toast.success(createFetcher.data.message || "Entity created");
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setOpen(false);
      } else {
        toast.error(createFetcher.data.message || "Failed to create entity");
      }
    }
  }, [createFetcher.data]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="bg-sidebar-primary hover:bg-sidebar-primary/90 self-start md:self-auto">
          New Mapping
        </Button>
      </DialogTrigger>
      <DialogContent className="data-[state=open]:!zoom-in-100 data-[state=open]:slide-in-from-bottom-20 data-[state=open]:duration-600 sm:max-w-[425px]">
        <createFetcher.Form
          method="post"
          action={`${DefaultPaths.ENTITY_LIST.path}/create-entity`}
          className="grid gap-4"
        >
          <DialogHeader>
            <DialogTitle>Create New Entity Mapping</DialogTitle>
            <DialogDescription>
              Provide an entity name to scaffold a fresh mapping configuration.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4">
            <div className="grid gap-3">
              <Label htmlFor="entity-name">Name</Label>
              <Input
                required
                id="entity-name"
                name="entity-name"
                placeholder="customer-to-user_customers"
              />
            </div>
          </div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
            <Button
              type="submit"
              disabled={createFetcher.state === "submitting"}
            >
              {createFetcher.state === "submitting" ? "Creating..." : "Create"}
            </Button>
          </DialogFooter>
        </createFetcher.Form>
      </DialogContent>
    </Dialog>
  );
}
