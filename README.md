### Setup infrastructure

- Start postgres container:

    ```bash
    make postgres
    ```

- Create simple_bank database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```


- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run test:

    ```bash
    make test
    ```

- Run program:

    ```bash
    go run main.go -table rims -version V1 -file data/rims_good.dat -max-fail-rate 5
    ```

    ```bash
    go run main.go -table rims -file data/rims.dat
    ```

    ```bash
    go run main.go -table timespans  -file data/timespans.dat
    ```
    

## Program structure:
- in db/migration are the up and down migrate schema
- in db/query are the used database queries for rims, timespans and logs
- make sqlc generates golang types and database connectors, using a docker image and the db/sqlc.yaml config file
- the store.go i define the database transaction and commit or rollback logic
- the store.go file also contains nested param structs and a mechanism to handle different table types and versions




## Anmerkungen
- Data Conversion can be optimized. It would be nice to generate the Converters with via some kind of Table Header File, that contain all relevant info.
- Data Validation can be optimized. It would be nice to add min/max values and other forms of validation to each DataType or Column.
- I did not find a good way to store the 00.00.0000 Date Format yet, so I safed a Null Value as equivalent.
- I can imagine automated database migration gerated by new table header files
- I did not have the time to do all testing. Some unit tests for database transfers with random test data are shown.
- I would probably reorganize the project, after a couple of more tables, and tableversions.
- I would probably rename some of the functions and param container soon.
- Also I would do some performance optimizations, after the refactoring is complete. There are lots of unnescessairy local variables to have a better readability.
- It would also do some performance tracking on the conversion and query times. It doesn't make sense to do concurrent queries, since it is always one connection on the same transaction, but it might make sens to spread the conversion and validation on different threads, if they take much longer, then the query.