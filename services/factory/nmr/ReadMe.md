# NMR

## Setup database

### Start database

```
 docker run --name postgresql-container -p 5432:5432 -e POSTGRES_PASSWORD=Pass2020! -d postgres
 docker run --net=host -e 'PGADMIN_DEFAULT_EMAIL=user@domain.com' -e 'PGADMIN_DEFAULT_PASSWORD=SuperSecret' dpage/pgadmin4
```
