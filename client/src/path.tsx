// define all routing
export const DefaultPaths = {
  CONNECTION_PAGE: {
    path: "/connection",
    pathname: "Connection",
    childPaths: {
      DATABASE: {
        path: "database",
        pathname: "Database",
      },
      DEBEZIUM: {
        path: "debezium",
        pathname: "Debezium",
      },
    },
  },
  ENTITY_LIST: {
    path: "/entity",
    pathname: "Entities",
    childPaths: {},
  },
  EXECUTION_LOG: {
    path: "/execution",
    pathname: "Execution",
    childPaths: {
      RUNNER: {
        path: "runner",
        pathname: "Runner",
      },
      LATEST_LOGS: {
        path: "logs",
        pathname: "Latest Logs",
      },
      HISTORY_LOGS: {
        path: "history",
        pathname: "History Logs",
      }
    },
  },
};
