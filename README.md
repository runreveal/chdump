chdump
======

This tool dumps the schema from the database given in the connection string or
the default database.  It does not dump or export any data contained in those
tables.

This is useful for cloning database schemas between instances or deployments of
clickhouse.

There are existing tools like `clickhouse-backup` which read the metadata
directories to achieve the same thing.  However, that tool (and other similar
tools) require filesystem access to the deployed database which isn't the case
if you're running clickhouse cloud or interacting with other hosted clickhouse
databases.

## Usage

To dump the schema, run the command like so:

```
chdump 'clickhouse://username:password@hostname[:port]/[database][?setting=value][&setting2=value]'
```

It spits the schema to standard output, with table definitions separated by a
semicolon and SQL comment line to make the visual identification of tables
clear.

You can then apply this schema to a new database like so:

```
clickhouse client --host <hostname> --user <username> --database [database] -n < schema.sql
```



