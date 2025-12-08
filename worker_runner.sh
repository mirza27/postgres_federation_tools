make parser  > logs/parser.log  2>&1 & echo $! > /tmp/parser.pid
make joiner > logs/joiner.log 2>&1 & echo $! > /tmp/joiner.pid
make checker > logs/checker.log 2>&1 & echo $! > /tmp/checker.pid
make executor > logs/executor.log 2>&1 & echo $! > /tmp/executor.pid