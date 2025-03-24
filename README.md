# simpleSRE

A simple project aimed at learning the basics of Site Reliability Engineering

## Go Commands

- Initialize module:
  - `go mod init <folder_name>`
- Run code:
  - `go run main.go`
  - `go run .`
  - To create a binary:
    - `go build .` then run the `<filename>.exe`
- Run SQL code:
  - `source C:\Users\Marwan\Desktop\projects\simpleSRE\backend\create-tables.sql`

### Notes

- [Creating counter table in MySQL](https://dba.stackexchange.com/questions/51736/counter-table-in-mysql)
  - According to this, it makes more sense to keep more than one row to represent the counter. You can then choose a random row and update it, then the counter is the sum of all the rows. Instead of having 1 row, 1 variable, you can run into 'global Mutex' issue for any transaction that tries to update the counter
  - [Another link that supports having multiple rows to represent a counter, instead of just 1](https://planetscale.com/blog/the-slotted-counter-pattern)
- [INT can overflow at 2.1 billion](https://stackoverflow.com/questions/47567287/bigint-signed-vs-unsigned) so I'll probably use this instead of BIGINT
- Unsigned means that the column can only store non-negative numbers
- `create-tables.sql` notes:

```sql
--Typical increment query could look like:
INSERT INTO simple_sre_db(slot, count)
VALUES (RAND() * 100, 1)
ON DUPLICATE KEY UPDATE count = count + 1;

--Getting the count
--Allows you to execute counter increments in parallel without causing contention and affecting concurrency
SELECT SUM(count) as count FROM simple_sre_db

```

- Following [this tutorial](https://go.dev/doc/tutorial/database-access) for connecting Go to MySQL
