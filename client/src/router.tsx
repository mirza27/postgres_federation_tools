import { EntityEditor } from "@/page/entity/entity-page";
import { ExecutionLog } from "@/page/execution/execution-log-page";
import { createBrowserRouter, Navigate } from "react-router-dom";
import { DefaultPaths } from "./path";
import { ConnectionLayout } from "./layout/connection-layout";
import MainLayout from "./layout/main-layout";
import { DebeziumPage } from "./page/connection/debezium-page";
import { DatabasePage } from "./page/connection/database-page";
import {
  checkDebeziumConnectorStatus,
  databaseCredentialsLoader,
} from "./page/connection/loader";
import {
  saveDatabaseSourceCredentials,
  saveDatabaseTargetCredentials,
} from "./page/connection/action";

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
              DefaultPaths.CONNECTION_PAGE.childPaths.DATABASE.path
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
                to={DefaultPaths.CONNECTION_PAGE.childPaths.DATABASE.path}
                replace
              />
            ),
          },

          {
            path: DefaultPaths.CONNECTION_PAGE.childPaths.DATABASE.path,
            element: <DatabasePage />,
            loader: databaseCredentialsLoader,
            children: [
              {
                path: "save-source",
                action: saveDatabaseSourceCredentials,
              },
              {
                path: "save-target",
                action: saveDatabaseTargetCredentials,
              },
            ],
          },
          {
            path: DefaultPaths.CONNECTION_PAGE.childPaths.DEBEZIUM.path,
            element: <DebeziumPage />,
            loader: checkDebeziumConnectorStatus,
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
