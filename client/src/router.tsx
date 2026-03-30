import { LatestQueueLogPage } from "@/feature/execution/latest-log-page";
import { createBrowserRouter, Navigate } from "react-router-dom";
import { DefaultPaths } from "./path";
import { ConnectionLayout } from "./layout/connection-layout";
import MainLayout from "./layout/main-layout";
import { EntityLayout } from "./layout/entity-layout";
import { DebeziumPage } from "./feature/connection/debezium-page";
import { DatabasePage } from "./feature/connection/database-page";
import {
  checkDebeziumConnectorStatus,
  databaseCredentialsLoader,
} from "./feature/connection/loader";
import {
  saveDatabaseSourceCredentials,
  saveDatabaseTargetCredentials,
} from "./feature/connection/action";
import { EntityListPage } from "./feature/entity-json/entity-list-page";
import {
  GetEntityLoader,
  ListAllEntitiesLoader,
} from "./feature/entity-json/entity-loader";
import { EntityPage as EntityDetailPage } from "./feature/entity-json/entity-page";
import {
  CreateNewEntity,
  DeleteEntity,
  UpdateEntity,
} from "./feature/entity-json/entity-action";
import { ExecutionLayout } from "./layout/execution-layout";
import { RunnerPage } from "./feature/execution/runner-page";
import { CheckWorkerStatus } from "./feature/execution/worker-loader";
import {
  RunWorkerAction,
  StopWorkerAction,
} from "./feature/execution/woker-action";
import { GetHistoryQueueLogsLoader, GetLatestQueueLogsLoader } from "./feature/execution/queue-loader";
import { HistoryLogPage } from "./feature/execution/history-log-page";

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
            path: DefaultPaths.EXECUTION_LOG.childPaths.LATEST_LOGS.path,
            element: <LatestQueueLogPage />,
            loader: GetLatestQueueLogsLoader,
          },
          {
            path: DefaultPaths.EXECUTION_LOG.childPaths.HISTORY_LOGS.path,
            element: <HistoryLogPage />,
            loader: GetHistoryQueueLogsLoader,
          },
        ],
      },
    ],
  },
]);
