// define all routing
export const DefaultPaths = {
  CONNECTION_PAGE: {
    path: "/connection",
    pathname: "connection",
    childPaths: {
      DATABASE: {
        path: "database",
        pathname: "database",
      },
      DEBEZIUM: {
        path: "debezium",
        pathname: "debezium",
      },
    },
  },
  ENTITY_LIST: {
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
