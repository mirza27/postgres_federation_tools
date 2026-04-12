import { LatestQueueLogPage } from "@/feature/execution/page/latest-log-page";
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
import { EntityListPage } from "./feature/entity-json/pages/entity-list-page";
import {
  createEmptyEntity,
  deleteEntity,
  GetEntityLoader,
  ListAllEntitiesLoader,
  updateEntity,
} from "./feature/entity-json/services/entity";
import { EntityPage as EntityDetailPage } from "./feature/entity-json/pages/entity-page";
import { ExecutionLayout } from "./layout/execution-layout";
import { RunnerPage } from "./feature/execution/page/runner-page";
import {
  checkWorkerStatus,
  runWorker,
  stopWorker,
} from "./feature/execution/services/worker";
import { HistoryLogPage } from "./feature/execution/page/history-log-page";
import {
  getHistoryQueueLogs,
  getLatestQueueLogs,
} from "./feature/execution/services/queue";
import { getEntityFilters } from "./feature/execution/services/entity";
import { EditEntityPage } from "./feature/entity-json/pages/entity-edit-page";
import { getEntityAndSchema } from "./feature/entity-json/services/database";

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
            path: "edit/:name",
            element: <EntityDetailPage />,
            loader: GetEntityLoader,
            children: [
              {
                path: "update-entity",
                action: updateEntity,
              },
            ],
          },
          {
            path: ":name",
            element: <EditEntityPage />,
            loader: getEntityAndSchema,
            children: [
              {
                path: "update-entity",
                action: updateEntity,
              },
            ],
          },

          {
            path: "delete-entity",
            action: deleteEntity,
          },
          {
            path: "create-entity",
            action: createEmptyEntity,
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
            loader: checkWorkerStatus,
            children: [
              {
                path: "run-worker",
                action: runWorker,
              },
              {
                path: "stop-worker",
                action: stopWorker,
              },
            ],
          },
          {
            path: DefaultPaths.EXECUTION_LOG.childPaths.LATEST_LOGS.path,
            element: <LatestQueueLogPage />,
            loader: getLatestQueueLogs,
          },
          {
            path: DefaultPaths.EXECUTION_LOG.childPaths.HISTORY_LOGS.path,
            element: <HistoryLogPage />,
            loader: getHistoryQueueLogs,
            children: [
              {
                path: "entity-filters",
                loader: getEntityFilters,
              },
            ],
          },
        ],
      },
    ],
  },
]);
