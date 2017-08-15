# Distributed logs system PostgresSQL to MongoDB

[![Go Report Card](https://goreportcard.com/badge/github.com/luk4z7/logres)](https://goreportcard.com/report/github.com/luk4z7/logres)

### Installation

```bash
go get github.com/luk4z7/logres
```

### Configuration
configure your postgres for work with logs type csv
```bash
vim /etc/postgresql/9.2/main/postgresql.conf
```

like this:
```bash
#------------------------------------------------------------------------------
# ERROR REPORTING AND LOGGING
#------------------------------------------------------------------------------

# - Where to Log -

log_destination = 'csvlog'              # Valid values are combinations of
# stderr, csvlog, syslog, and eventlog,
# depending on platform.  csvlog
# requires logging_collector to be on.

# This is used when logging to stderr:
logging_collector = on                  # Enable capturing of stderr and csvlog
# into log files. Required to be on for
# csvlogs.
# (change requires restart)
```

if do you like prefer change the directive `log_statement` for log all type of query, errors, warning, etc ...
```bash
log_statement = 'all'                 # none, ddl, mod, all
```

Restart your postgresql
```bash
/etc/init.d/postgresql restart
```

Configure your environment

```bash
logres --config
```

and then execute the `--run`
```bash
logres --run
```
