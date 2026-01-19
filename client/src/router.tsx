import { ExecutionLogPage } from "@/page/execution/execution-log-page";
import { createBrowserRouter, Navigate } from "react-router-dom";
import { DefaultPaths } from "./path";
import { ConnectionLayout } from "./layout/connection-layout";
import MainLayout from "./layout/main-layout";
import { EntityLayout } from "./layout/entity-layout";
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
import { EntityListPage } from "./page/entity-json/entity-list-page";
import {
  GetEntityLoader,
  ListAllEntitiesLoader,
} from "./page/entity-json/entity-loader";
import { EntityPage as EntityDetailPage } from "./page/entity-json/entity-page";
import {
  CreateNewEntity,
  DeleteEntity,
  UpdateEntity,
} from "./page/entity-json/entity-action";
import { ExecutionLayout } from "./layout/execution-layout";
import { RunnerPage } from "./page/execution/runner-page";
import { CheckWorkerStatus } from "./page/execution/worker-loader";
import {
  RunWorkerAction,
  StopWorkerAction,
} from "./page/execution/woker-action";
import { GetExecutionQueueListLoader } from "./page/execution/queue-loader";

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
        path: DefaultPaths.ENTITY_LIST.path,
        id: "entity-list-route",
        element: <EntityLayout />,
        loader: ListAllEntitiesLoader,
        children: [
          {
            index: true,
            element: <EntityListPage />,
          },
          {
            path: ":name",
            element: <EntityDetailPage />,
            loader: GetEntityLoader,
            children: [
              {
                path: "update-entity",
                action: UpdateEntity,
              },
            ],
          },
          {
            path: "delete-entity",
            action: DeleteEntity,
          },
          {
            path: "create-entity",
            action: CreateNewEntity,
          },
        ],
      },
      {
        path: DefaultPaths.EXECUTION_LOG.path,
        element: <ExecutionLayout />,
        children: [
          {
            index: true,
            element: (
              <Navigate
                to={
                  DefaultPaths.EXECUTION_LOG.path +
                  "/" +
                  DefaultPaths.EXECUTION_LOG.childPaths.RUNNER.path
                }
                replace
              />
            ),
          },
          {
            path: DefaultPaths.EXECUTION_LOG.childPaths.RUNNER.path,
            element: <RunnerPage />,
            loader: CheckWorkerStatus,
            children: [
              {
                path: "run-worker",
                action: RunWorkerAction,
              },
              {
                path: "stop-worker",
                action: StopWorkerAction,
              },
            ],
          },
          {
            path: DefaultPaths.EXECUTION_LOG.childPaths.LOGS.path,
            element: <ExecutionLogPage />,
            loader: GetExecutionQueueListLoader,
          },
        ],
      },
    ],
  },
]);
