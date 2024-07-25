Test for applied System

This is a gin based API server to save graphs. This will start server on port 8080

It has 3 API's
1. POST `/graphs` -> Post a new graph and returns with a id
   ```bash
   curl --location 'localhost:8080/graphs' 
    --header 'Content-Type: application/json' \
    --data '[
        {
        "v1":"A",
        "v2":"B"
        },
        {
        "v1":"B",
        "v2":"C"
        },
        {
        "v1":"C",
        "v2":"D"
        },
        {
        "v1":"D",
        "v2":"E"
        },
        {
        "v1":"E",
        "v2":"F"
        },
        {
        "v1":"B",
        "v2":"F"
        }
    ]'
   ```

2. GET `/graphs/:id/shortest-path?start=<start_vertice>&end=<end_vertice>`.  
`id`:- id of the graph to find shortest path.   
`query-params`
   start and end vertices 
```bash
curl --location --request GET 'localhost:8080/graphs/cf4dcc78-f1cd-46e8-af96-171510306b8d/shortest-path?start=A&end=F' \
--header 'Content-Type: application/json' 
```

3. DELETE `/graphs/:id` 
   Delete a graph with id

# Steps to run 
1. Clone this repo.
2. start the server with `go run main.go`
3. Run test cases with `go test ./...`
4. Run benchmark test case with `cd controllers go test -bench .`