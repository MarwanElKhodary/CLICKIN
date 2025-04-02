# Backend Documentation

## Go Commands

- Initialize module:
  - `go mod init <folder_name>`
- Run code:
  - `go run main.go`
  - `go run .`
  - To create a binary:
    - `go build .` then run the `<filename>.exe`
- Run SQL code:
  - `use simple_sre_db;`
  - `source C:/Users/Marwan/Desktop/projects/simpleSRE/backend/create-tables.sql;`

### Notes

- [Creating counter table in MySQL](https://dba.stackexchange.com/questions/51736/counter-table-in-mysql)
  - According to this, it makes more sense to keep more than one row to represent the counter. You can then choose a random row and update it, then the counter is the sum of all the rows. Instead of having 1 row, 1 variable, you can run into 'global Mutex' issue for any transaction that tries to update the counter
  - [Another link that supports having multiple rows to represent a counter, instead of just 1](https://planetscale.com/blog/the-slotted-counter-pattern)
- [INT can overflow at 2.1 billion](https://stackoverflow.com/questions/47567287/bigint-signed-vs-unsigned) so I'll probably use this instead of BIGINT
- Unsigned means that the column can only store non-negative numbers
- `create-tables.sql` notes:

```sql
--To use db in MySQL CLI:
use simple_sre_db;

--Typical increment query could look like:
INSERT INTO simple_sre_db(slot, count)
VALUES (RAND() * 100, 1)
ON DUPLICATE KEY UPDATE count = count + 1;

--Getting the count
--Allows you to execute counter increments in parallel without causing contention and affecting concurrency
SELECT SUM(count) as count FROM count_table;

```

- Following [this tutorial](https://go.dev/doc/tutorial/database-access) for connecting Go to MySQL
- How you install packages is as follows:
  - `go install <the_package_you_want>`
  - `go get .` to track the module as a dependency
- [How to Write Go Code](https://go.dev/doc/code)
- [Effective Go](https://go.dev/doc/effective_go)
- Some relevant links on project architecture:
  - [Repository pattern](https://threedots.tech/post/repository-pattern-in-go/)
  - [Go Backend Rest API services](https://medium.com/@janishar.ali/how-to-architecture-good-go-backend-rest-api-services-14cc4730c05b)
  - [Simple Go Project Folder Structures](https://medium.com/@smart_byte_labs/organize-like-a-pro-a-simple-guide-to-go-project-folder-structures-e85e9c1769c2)
    - For now, basing the project structure on this
- Go through [this](https://go.dev/tour/methods/1) to understand more about Functions, Methods and Receivers like `func (r *Repository) IncrementCount(slot int, count int) (int64, error) {}`
- Better way to organize services [here](https://medium.com/@ott.kristian/how-i-structure-services-in-go-19147ad0e6bd)
  - But like come back to this when you have more services, or close to the end
- Interesting conceptual articles:
  - [About stack vs heap in Go](https://medium.com/eureka-engineering/understanding-allocations-in-go-stack-heap-memory-9a2631b5035d)
  - [When to use pointers](https://medium.com/@meeusdylan/when-to-use-pointers-in-go-44c15fe04eac)
  - Cool command "escape analysis": `go build -gcflags="-m"`
- What happens if a user is on different browsers on the same pc?
- Should we consider the option of using a websocket when the user is spamming the button? Is that better than making HTTP calls on every click?
