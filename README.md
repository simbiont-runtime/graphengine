# What is GraphEngine

graphEngine is a [PGQL](https://pgql-lang.org/)-compatible embeddable graph database for large-scale vertices. 
You can find the specification [here](https://pgql-lang.org/spec/1.5/).

## Installation

```bash
go get -u github.com/simbiont-runtime/graphengine
```

## Quick Start

### Playground

```bash
> make build
> ./bin/graphengine play
```

### Build Application
graphEngine implements driver for [database/sql](https://golang.org/pkg/database/sql/). To use graphEngien, you can simply import GraphEngine package and use `graphEngine` as the driver name in `sql.Open`.

Here is an example of how to use create a graph and query it. The graph is an example in [PGQL specification](https://pgql-lang.org/spec/1.5/#edge-patterns).

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/simbiont-runtime/graphengine"
)

func main() {
	db, err := sql.Open("graphEngine", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create a graph.
	mustExec(ctx, conn, "CREATE GRAPH student_network")

	// Change the current graph.
	mustExec(ctx, conn, "USE student_network")

	// Create labels.
	mustExec(ctx, conn, "CREATE LABEL Person")
	mustExec(ctx, conn, "CREATE LABEL University")
	mustExec(ctx, conn, "CREATE LABEL knows")
	mustExec(ctx, conn, "CREATE LABEL studentOf")

	// Create vertices.
	mustExec(ctx, conn, `INSERT VERTEX x LABELS (Person) PROPERTIES (x.name = 'Kathrine', x.dob = DATE '1994-01-15')`)
	mustExec(ctx, conn, `INSERT VERTEX x LABELS (Person) PROPERTIES (x.name = 'Riya', x.dob = DATE '1995-03-20')`)
	mustExec(ctx, conn, `INSERT VERTEX x LABELS (Person) PROPERTIES (x.name = 'Lee', x.dob = DATE '1996-01-20')`)
	mustExec(ctx, conn, `INSERT VERTEX x LABELS (University) PROPERTIES (x.name = 'UC Berkeley')`)

	// Create edges.
	mustExec(ctx, conn, `INSERT EDGE e BETWEEN x AND y LABELS ( knows ) FROM MATCH (x), MATCH (y) WHERE x.name = 'Kathrine' AND y.name = 'Lee'`)
	mustExec(ctx, conn, `INSERT EDGE e BETWEEN x AND y LABELS ( knows ) FROM MATCH (x), MATCH (y) WHERE x.name = 'Kathrine' AND y.name = 'Riya'`)
	mustExec(ctx, conn, `INSERT EDGE e BETWEEN x AND y LABELS ( knows ) FROM MATCH (x), MATCH (y) WHERE x.name = 'Lee' AND y.name = 'Kathrine'`)
	mustExec(ctx, conn, `INSERT EDGE e BETWEEN x AND y LABELS ( studentOf ) FROM MATCH (x), MATCH (y) WHERE x.name = 'Kathrine' AND y.name = 'UC Berkeley'`)
	mustExec(ctx, conn, `INSERT EDGE e BETWEEN x AND y LABELS ( studentOf ) FROM MATCH (x), MATCH (y) WHERE x.name = 'Lee' AND y.name = 'UC Berkeley'`)
	mustExec(ctx, conn, `INSERT EDGE e BETWEEN x AND y LABELS ( studentOf ) FROM MATCH (x), MATCH (y) WHERE x.name = 'Riya' AND y.name = 'UC Berkeley'`)

	// Query the graph.
	rows, err := conn.QueryContext(ctx, "SELECT a.name AS a, b.name AS b FROM MATCH (a:Person) -[e:knows]-> (b:Person)")
	if err != nil {
		log.Fatal(err)
	}
	var a, b string
	for rows.Next() {
		if err := rows.Scan(&a, &b); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("'%s' knows '%s'\n", a, b)
	}
}

func mustExec(ctx context.Context, conn *sql.Conn, query string) {
	_, err := conn.ExecContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Contributing

We welcome contributions from everyone. GraphEngine is in its early stages, if you have any ideas or suggestions, please feel free to open an issue or pull request.

## License

graphEngine is licensed under the Apache 2.0 license. See [LICENSE](LICENSE) for the full license text.
