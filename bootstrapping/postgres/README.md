# Postgres refresher

## psql commands

\d __database__

\c __database__

\d __table__

## Schema

    CREATE TABLE $SHORTY_DB_TABLE(
    id             SERIAL PRIMARY KEY,
    shorty         VARCHAR(50) NOT NULL,
    url            VARCHAR(255) NOT NULL,
    UNIQUE (shorty)
    );

## SQL
INSERT INTO shorties VALUES(DEFAULT,'ggl.c', 'https://google.com');

INSERT INTO shorties(shorty,url) VALUES('ggl.u', 'https://google.co.uk');

## Postgres running in a Docker container

### Access the postgres interactive terminal
docker exec -it url-shorty-postgres psql -U shorty_user url_shortener

### Run SQL command locally
docker exec -it url-shorty-postgres psql -U shorty_user url_shortener --command "SELECT * FROM shorties"

    - url-shortener-postgres    : container name
    - url_shortener             : database
    - psql ... --command "...." : running an SQL command outside the database env


#### Without gorm AutoMigrate feature
docker exec -it url-shorty-postgres psql -U postgres url_shortener --command "INSERT INTO shorties VALUES(DEFAULT,'ggl.c', 'https://google.com')"

docker exec -it url-shorty-postgres psql -U shorty_user --password url_shortener --command "INSERT INTO shorties VALUES(DEFAULT,'amzn.u', 'https://amazon.co.uk')"

    asks interactively for a password

#### With gorm AutoMigrate feature

gorm AutoMigrade feature takes the provided struct/slice and creates a corresponding table. The table contains additional columns:

    - created_at | timestamp with time zone |
    - updated_at | timestamp with time zone |
    - deleted_at | timestamp with time zone |

CAVEAT regarding <id> field

If the data structure corresponding to the table contains gorm.Model, eg:

    type Shorty struct {
        gorm.Model
        Shorty, Url string
    }

the primary key id is shorty.Model.ID. Adding an id as the struct field will result in an error message about a duplicate id and the table will not be created.

#### Without gorm AutoMigrate feature

If not using the gorm AutoMigrate feature, the database table must be create (see bootstrapping/postgres/init/ script copmments)


### Run SQL commands against a container
docker exec -it url-shorty-postgres psql -U postgres  url_shortener --command "INSERT INTO shorties VALUES(DEFAULT,DEFAULT,DEFAULT,DEFAULT,'ggl.c', 'https://google.com')"

docker exec -it url-shorty-postgres psql -U shorty_user --password url_shortener --command "INSERT INTO shorties VALUES(DEFAULT,DEFAULT,DEFAULT,DEFAULT,'amzn.u', 'https://amazon.co.uk')"

    asks interactively for a password

## Useful links

https://gist.github.com/Kartones/dd3ff5ec5ea238d4c546

https://www.w3resource.com/PostgreSQL/unique.php

http://doc.gorm.io/crud.html#delete