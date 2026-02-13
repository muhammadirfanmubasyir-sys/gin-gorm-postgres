This is README!
===============
go get -u github.com/gin-gonic/gin

go get -u gorm.io/gorm

go get -u gorm.io/driver/postgres
==================================================

psql -h localhost -p 5432 -d irfan-db -U postgres

irfan-db=# \c postgres
You are now connected to database "postgres" as user "postgres".

irfan-db=# \list
                                                                    List of databases
   Name    |  Owner   | Encoding | Locale Provider |          Collate           |           Ctype            | Locale | ICU Rules |   Access privileges
-----------+----------+----------+-----------------+----------------------------+----------------------------+--------+-----------+-----------------------
 irfan-db  | postgres | UTF8     | libc            | English_United States.1252 | English_United States.1252 |        |           |
 postgres  | postgres | UTF8     | libc            | English_United States.1252 | English_United States.1252 |        |           |
(2 rows)

irfan-db=# \d+ users
                                                                  Table "public.users"
   Column   |           Type           | Collation | Nullable |              Default              | Storage  | Compression | Stats target | Descr
iption
------------+--------------------------+-----------+----------+-----------------------------------+----------+-------------+--------------+------
-------
 id         | bigint                   |           | not null | nextval('users_id_seq'::regclass) | plain    |             |              |
 created_at | timestamp with time zone |           |          |                                   | plain    |             |              |
 updated_at | timestamp with time zone |           |          |                                   | plain    |             |              |
 deleted_at | timestamp with time zone |           |          |                                   | plain    |             |              |
 name       | text                     |           |          |                                   | extended |             |              |
 email      | text                     |           |          |                                   | extended |             |              |
 password   | text                     |           |          |                                   | extended |             |              |
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "idx_users_deleted_at" btree (deleted_at)
Not-null constraints:
    "users_id_not_null" NOT NULL "id"