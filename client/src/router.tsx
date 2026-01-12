import { EntityEditor } from "@/page/entity";
import { ExecutionLog } from "@/page/execution-log";
import { createBrowserRouter, Navigate } from "react-router-dom";
import { DefaultPaths } from "./path";
import { ConnectionLayout } from "./layout/connection-layout";
import MainLayout from "./layout/main-layout";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: (
          <Navigate
            to={
              DefaultPaths.CONNECTION_PAGE.path +
              "/" +
              DefaultPaths.CONNECTION_PAGE.childPaths.SOURCE_DATABASE.path
            }
            replace
          />
        ),
      },
      {
        path: DefaultPaths.CONNECTION_PAGE.path,
        element: <ConnectionLayout />,
        children: [
          {
            index: true,
            element: (
              <Navigate
                to={
                  DefaultPaths.CONNECTION_PAGE.childPaths.SOURCE_DATABASE.path
                }
                replace
              />
            ),
          },
          {
            element: <div>Connection Source Database</div>,
            path: DefaultPaths.CONNECTION_PAGE.childPaths.SOURCE_DATABASE.path,
          },
          {
            path: DefaultPaths.CONNECTION_PAGE.childPaths.TARGET_DATABASE.path,
            element: <div>Connection Target Database</div>,
          },
          {
            path: DefaultPaths.CONNECTION_PAGE.childPaths.DEBEZIUM.path,
            element: <div>Connection Debezium</div>,
          },
        ],
      },
      {
        path: DefaultPaths.ENTITY_EDITOR.path,
        element: <EntityEditor />,
      },
      {
        path: DefaultPaths.EXECUTION_LOG.path,
        element: <ExecutionLog onBack={() => null} />,
      },
    ],
  },
]);
