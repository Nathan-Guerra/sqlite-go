[![progress-banner](https://backend.codecrafters.io/progress/sqlite/3ee4e69a-746c-4e9f-8774-8110424498c1)](https://app.codecrafters.io/users/Nathan-Guerra?r=2qF)

This is my attempt on implementing the sqlite cli from scratch based on codecrafters's
["Build Your Own SQLite" Challenge](https://codecrafters.io/challenges/sqlite).

Right now it supports the following dot commands:
- **.dbinfo**
   It correctly outputs the size of the database pages and the number of rows in
   sample.db. However, it 
- **.tables**
   Display all tables registered at `sqlite_schema`.

## Known limitations of this implementation

- [ ] At `.dbinfo` the size of tables was achieved only using the first page's
row counting. This solution will output the incorrect solution if there are 
triggers or indexes in the `sqlite_schema` table. To fix this, I believe we have
to walk through the sqlite_schema rows counting every row with type *table*.
- [ ] At `.tables` the payload implementation do not take into account the
possibility of the sqlite_schema having overflow pages, so the solution is only
90% complete.
