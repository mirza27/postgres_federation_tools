// define all routing
export const DefaultPaths = {
  CONNECTION_PAGE: {
    path: "/connection",
    pathname: "connection",
    childPaths: {
      SOURCE_DATABASE: {
        path: "source-database",
        pathname: "source-database",
      },
      TARGET_DATABASE: {
        path: "target-database",
        pathname: "target-database",
      },
      DEBEZIUM: {
        path: "debezium",
        pathname: "debezium",
      },
    },
  },
  ENTITY_EDITOR: {
    path: "/entity",
    pathname: "entity",
    childPaths: [],
  },
  EXECUTION_LOG: {
    path: "/execution-log",
    pathname: "execution-log",
    childPaths: [],
  },
};
