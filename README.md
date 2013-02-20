json2csv
========

Converts a stream of newline separated json data to csv format. This version of json2csv supports nested fields.

Installation
============

If you have a working golang install, you can use `go get`.

```bash
go get github.com/jeroenjanssens/json2csv
```

Usage
=====

```
usage: json2csv
    -k fields,and,nested.fields,to,output
    -i /path/to/input.json (optional; default is stdin)
    -o /path/to/output.csv (optional; default is stdout)
    -v verbose output (to stderr)
    --version
    -h this help
```

To convert:

```json
{"user": {"name":"jehiah", "password": "root"}, "remote_ip": "127.0.0.1", "dt" : "[20/Aug/2010:01:12:44 -0400]"}
{"user": {"name":"jeroenjanssens", "password": "123"}, "remote_ip": "192.168.0.1", "dt" : "[20/Aug/2010:01:12:44 -0400]"}
{"user": {"name":"unknown", "password": ""}, "remote_ip": "76.216.210.0", "dt" : "[20/Aug/2010:01:12:45 -0400]"}
```

to:

```
"jehiah","127.0.0.1"
"jeroenjanssens","192.168.0.1"
"unknown","76.216.210.0"
```
    
you would either

```bash
json2csv -k user.name,remote_ip -i input.json -o output.csv
```

or

```bash
cat input.json | json2csv -k user.name,remote_ip > output.csv
```

