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

## Decisions

- No error handling on the frontend
- Counter updates are initiated via POST requests and broadcasted via WebSockets
  - Revisit Server Sent Events if need be

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
- Details on [how to deploy a Go/HTMX application](https://community.aws/content/2hYjbCwWyM3KAuR77j9DqE1P4p7/deploying-a-go-application-with-htmx-to-aws-elastic-beanstalk-a-step-by-step-guide?lang=en)
- Read more about preflight, OPTIONS requests, and setting proper headers to your APIs
- Tests are following [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests) structure. It currently does not include subtests properly or good takedown, but it's working
- Good reference on [whether or not I need websockets](https://hntrl.io/posts/you-dont-need-websockets/), although a lot of comments are saying this is bad advice
- Have to use `go test -race` at some point
- Based off of [this article](https://netflixtechblog.com/netflixs-distributed-counter-abstraction-8d0c45eb66b2) you should maybe consider restructing the API's request/response to look more like this to have `generation time`

```json
{
  "namespace": "my_dataset",
  "counter_name": "counter123",
  "delta": 2,
  "idempotency_token": { 
    "token": "some_event_id",
    "generation_time": "2024-10-05T14:48:00Z"
  }
}
```

- Use `ifconfig | grep "inet " | grep -v 127.0.0.1` to get your ip address on Mac
  - The first ip address should be usable along with the port number on other devices in the same network
  - You might need to disable firewall
  - At this point, other devices are getting `htmx:afterRequest`, `htmx:sendError`, `Failed to load resource: net::ERR_CONNECTION_REFUSED` errors
- Some good links on websockets:
  - [HTMX extension](https://htmx.org/extensions/ws/)
  - [Simple guide on websockets in Go](https://medium.com/wisemonks/implementing-websockets-in-golang-d3e8e219733b)
  - [Guide on gorilla + gin in Go](https://medium.com/@abhishekranjandev/building-a-production-grade-websocket-for-notifications-with-golang-and-gin-a-detailed-guide-5b676dcfbd5a)
- Look up MySQl security best practices
- Look up Gin security best practices
- Look up HTMX security best practices
- Look up API security best practices
- Look up WebSockets security best practices
- Is there anything wrong with spamming API POST request?
- [The Twelve Factor App](https://12factor.net/) - an amazing resource on deploying an app. LITERALLY A GODSEND
- Definitely need to think about [most common security vulnerabilites](https://owasp.org/www-project-top-ten/) soon
  - [Good resource](https://senowijayanto.medium.com/securing-your-go-backend-encryption-vulnerability-prevention-and-more-3fc980f45a8f) on securing Go backend via encryption - which is #2
- `ngrok http http://localhost:8080` amazing to access apps on other devices [ngrok](https://dashboard.ngrok.com/get-started/setup/macos) amazing service
- [Amazing resource](https://threedots.tech/post/live-website-updates-go-sse-htmx/) on live website updates with Go, SSE and htmx
