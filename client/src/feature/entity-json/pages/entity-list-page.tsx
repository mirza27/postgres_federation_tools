import { ChevronLeft } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useEffect } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useFetcher, useRouteLoaderData } from "react-router-dom";
import { DefaultPaths } from "@/path";
import {
  getSourceTables,
  getTargetTables,
  type EntityListDataResponse,
} from "../services/entity-loader";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { NewEntityForm } from "../components/_new-entity-form";

export function EntityListPage() {
  const loaderData = useRouteLoaderData("entity-list-route") as {
    ok: boolean;
    data: EntityListDataResponse | null;
    message?: string;
  };

  const navigateMappingPage = (entityName: string) => {
    // eslint-disable-next-line react-hooks/immutability
    window.location.href = `${
      DefaultPaths.ENTITY_LIST.path
    }/${encodeURIComponent(entityName)}`;
  };

  const deleteFetcher = useFetcher();

  // message when delete entity
  useEffect(() => {
    if (deleteFetcher.data) {
      if (deleteFetcher.data.ok) {
        toast.success(
          deleteFetcher.data.message || "Entity deleted successfully",
        );
      } else {
        toast.error(deleteFetcher.data.message || "Failed to delete entity");
      }
    }
  }, [deleteFetcher.data]);

  return (
    <div className="h-screen w-full flex flex-col bg-background">
      {/* header */}
      <div className="border-b border-border px-6 py-4">
        <div className="flex items-center gap-3">
          <button className="p-1 hover:bg-muted rounded-md transition-colors">
            <ChevronLeft className="w-5 h-5 text-foreground" />
          </button>
          <div>
            <h1 className="text-2xl font-semibold">Mapping Entity List</h1>
            <p className="text-sm text-muted-foreground">
              Configure your mapping entities here
            </p>
          </div>
        </div>
      </div>

      {/* content */}
      <div className="flex-1 overflow-auto p-6">
        <Card className="border border-border bg-card">
          <CardContent className="py-2 px-6 flex flex-col gap-6">
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <CardHeader className="px-0 pb-0 md:flex-1 md:pr-8">
                <CardTitle>Entity List</CardTitle>
                <CardDescription>
                  You could create and manage your mapping entities here. Each
                  entity represents a mapping between source and target data
                  structures. 1 entity could consist more than 1 table in source
                  database
                </CardDescription>
              </CardHeader>

              <NewEntityForm />
            </div>

            <Table>
              <TableHeader>
                <TableRow className="hover:bg-transparent">
                  <TableHead>Entity Name</TableHead>
                  <TableHead>Source Tables</TableHead>
                  <TableHead>Target Tables</TableHead>
                  <TableHead className="text-right">Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loaderData.data?.entities?.map((entity) => (
                  <TableRow key={entity.entity}>
                    <TableCell>{entity.entity}</TableCell>

                    <TableCell>{getSourceTables(entity.sources)}</TableCell>

                    <TableCell>{getTargetTables(entity)}</TableCell>

                    <TableCell className="text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          type="button"
                          size="sm"
                          variant="outline"
                          className="px-3 text-xs"
                          onClick={() => navigateMappingPage(entity.entity)}
                        >
                          Edit
                        </Button>

                        {/* delete dialog */}
                        <AlertDialog>
                          <AlertDialogTrigger asChild>
                            <Button
                              type="button"
                              size="sm"
                              variant="destructive"
                              className="px-3 text-xs"
                            >
                              Delete
                            </Button>
                          </AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>
                                Delete mapping "{entity.entity}"?
                              </AlertDialogTitle>
                              <AlertDialogDescription>
                                This action cannot be undone. The mapping
                                configuration will be permanently removed.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Cancel</AlertDialogCancel>

                              {/* delete fetcher mechanism */}
                              <deleteFetcher.Form
                                method="delete"
                                action="delete-entity"
                              >
                                <Input
                                  type="hidden"
                                  name="entity_name"
                                  value={entity.entity}
                                />
                                <AlertDialogAction
                                  type="submit"
                                  className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                >
                                  Confirm delete
                                </AlertDialogAction>
                              </deleteFetcher.Form>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
