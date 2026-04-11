import { ChevronDown, Database, Settings, FileInput } from "lucide-react";
import { useNavigate, useLocation, useRouteLoaderData } from "react-router-dom";
import { DefaultPaths } from "../path";
import clsx from "clsx";
import type { EntityListDataResponse } from "@/feature/entity-json/services/entity-loader";

export function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();
  const pathname = location.pathname;
  const entityLoaderData = useRouteLoaderData("entity-list-route") as
    | {
        ok: boolean;
        data: EntityListDataResponse | null;
        message?: string;
      }
    | undefined;

  const entities = entityLoaderData?.data?.entities ?? [];
  const activeEntity = pathname.startsWith(`${DefaultPaths.ENTITY_LIST.path}/`)
    ? decodeURIComponent(
        pathname.slice(DefaultPaths.ENTITY_LIST.path.length + 1),
      )
    : null;

  const isPathActive = (current: string, target: string) => {
    if (target === "/") return current === "/";
    return current === target || current.startsWith(`${target}/`);
  };

  const expandedMenu = (() => {
    if (isPathActive(pathname, DefaultPaths.CONNECTION_PAGE.path))
      return "connection";
    if (isPathActive(pathname, DefaultPaths.ENTITY_LIST.path)) return "entity";
    if (isPathActive(pathname, DefaultPaths.EXECUTION_LOG.path))
      return "execution";
    return null;
  })();

  const baseButton =
    "w-full flex items-center justify-between px-3 py-2 rounded-md text-sm transition-colors";
  const childButton =
    "w-full text-left px-3 py-1.5 rounded text-xs transition-colors";

  return (
    <div className="w-64 bg-sidebar border-r border-sidebar-border h-screen flex flex-col">
      {/* Header */}
      <div className="p-4 border-b border-sidebar-border">
        <h1 className="flex items-center gap-2 text-lg font-semibold">
          <Database className="w-5 h-5 text-sidebar-primary" />
          DB Migrator
        </h1>
      </div>

      {/* Menu */}
      <div className="flex-1 p-4 space-y-2 overflow-y-auto">
        {/* Connection */}
        <div>
          <button
            onClick={() => navigate(DefaultPaths.CONNECTION_PAGE.path)}
            className={clsx(
              baseButton,
              isPathActive(pathname, DefaultPaths.CONNECTION_PAGE.path) &&
                "bg-sidebar-accent/30",
            )}
          >
            <div className="flex items-center gap-2">
              <Database className="w-4 h-4" />
              <span>{DefaultPaths.CONNECTION_PAGE.pathname}</span>
            </div>
            <ChevronDown
              className={clsx(
                "w-4 h-4 transition-transform",
                expandedMenu === "connection" && "rotate-180",
              )}
            />
          </button>

          {expandedMenu === "connection" && (
            <div className="ml-4 mt-2 space-y-1">
              {Object.values(DefaultPaths.CONNECTION_PAGE.childPaths).map(
                (child) => {
                  const fullPath = `${DefaultPaths.CONNECTION_PAGE.path}/${child.path}`;
                  const isActive = pathname === fullPath;

                  return (
                    <button
                      key={child.path}
                      onClick={() => navigate(fullPath)}
                      className={clsx(
                        childButton,
                        isActive
                          ? "bg-sidebar-primary/30 text-sidebar-primary font-medium"
                          : "text-sidebar-foreground/80 hover:bg-sidebar-primary/20",
                      )}
                    >
                      {child.pathname}
                    </button>
                  );
                },
              )}
            </div>
          )}
        </div>

        {/* Entity */}
        <div>
          <button
            onClick={() => navigate(DefaultPaths.ENTITY_LIST.path)}
            className={clsx(
              baseButton,
              isPathActive(pathname, DefaultPaths.ENTITY_LIST.path) &&
                "bg-sidebar-accent/30",
            )}
          >
            <div className="flex items-center gap-2">
              <Settings className="w-4 h-4" />
              <span>{DefaultPaths.ENTITY_LIST.pathname}</span>
            </div>
            <ChevronDown
              className={clsx(
                "w-4 h-4 transition-transform",
                expandedMenu === "entity" && "rotate-180",
              )}
            />
          </button>

          {expandedMenu === "entity" && entities.length > 0 && (
            <div className="ml-4 mt-2 space-y-1">
              {entities.map((entity) => {
                const entityName = entity.entity;
                const isActive = activeEntity === entityName;

                return (
                  <button
                    key={entityName}
                    onClick={() =>
                      navigate(
                        `${DefaultPaths.ENTITY_LIST.path}/${encodeURIComponent(
                          entityName,
                        )}`,
                      )
                    }
                    className={clsx(
                      childButton,
                      isActive
                        ? "bg-sidebar-primary/30 text-sidebar-primary font-medium"
                        : "text-sidebar-foreground/80 hover:bg-sidebar-primary/20",
                    )}
                  >
                    {entityName}
                  </button>
                );
              })}
            </div>
          )}
          {expandedMenu === "entity" && entities.length === 0 && (
            <div className="ml-4 mt-2 text-xs text-sidebar-foreground/70">
              No entities available yet
            </div>
          )}
        </div>

        {/* Execution Log */}
        <div>
          <button
            onClick={() => navigate(DefaultPaths.EXECUTION_LOG.path)}
            className={clsx(
              baseButton,
              isPathActive(pathname, DefaultPaths.EXECUTION_LOG.path) &&
                "bg-sidebar-accent/30",
            )}
          >
            <div className="flex items-center gap-2">
              <FileInput className="w-4 h-4" />
              <span>{DefaultPaths.EXECUTION_LOG.pathname}</span>
            </div>
            <ChevronDown
              className={clsx(
                "w-4 h-4 transition-transform",
                expandedMenu === "execution" && "rotate-180",
              )}
            />
          </button>

          {expandedMenu === "execution" && (
            <div className="ml-4 mt-2 space-y-1">
              {Object.values(DefaultPaths.EXECUTION_LOG.childPaths).map(
                (child) => {
                  const fullPath = `${DefaultPaths.EXECUTION_LOG.path}/${child.path}`;
                  const isActive = pathname === fullPath;

                  return (
                    <button
                      key={child.path}
                      onClick={() => navigate(fullPath)}
                      className={clsx(
                        childButton,
                        isActive
                          ? "bg-sidebar-primary/30 text-sidebar-primary font-medium"
                          : "text-sidebar-foreground/80 hover:bg-sidebar-primary/20",
                      )}
                    >
                      {child.pathname}
                    </button>
                  );
                },
              )}
            </div>
          )}
        </div>
      </div>

      {/* Footer */}
      <div className="p-4 border-t border-sidebar-border text-xs opacity-60">
        v1.0.0
      </div>
    </div>
  );
}
