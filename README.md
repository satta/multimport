# multimport

This is a tool that evenly distributes VAST import input across a given number of import processes, restarting the import processes if they return for any reason.

## Usage

```
$ multimport
Helps with importing lots of events in parallel.

Usage:
  multimport [command]

Available Commands:
  help        Help about any command
  suricata    import Suricata data

Flags:
  -p, --extra-params strings   extra parameters to pass to 'vast import'
  -h, --help                   help for multimport
  -j, --jobs uint              amount of parallel VAST import processes (default 4)
      --vast-path string       VAST executable (default "vast")
  -v, --verbose                be verbose and show debug level output

Use "multimport [command] --help" for more information about a command.
```

## Example

Here's an example invocation for 3000 simulated input events per second (via [speeve](https://github.com/satta/speeve)) and 3 importers, with one of them killed during the import:


```
$ speeve spew -p ~/profile.yaml -s 3000 | multimport -j3 suricata
INFO[0000] 2021-01-08T14:17:31.995 loaded configuration file: "/etc/vast/vast.yaml"  importer=suri_0
INFO[0000] 2021-01-08T14:17:31.996 loaded configuration file: "/etc/vast/vast.yaml"  importer=suri_1
INFO[0000] 2021-01-08T14:17:31.996 loaded configuration file: "/etc/vast/vast.yaml"  importer=suri_2
INFO[0000] 2021-01-08T14:17:32.003 connecting to VAST node localhost:42000  importer=suri_2
INFO[0000] 2021-01-08T14:17:32.005 suricata-reader reads data from STDIN  importer=suri_2
INFO[0000] 2021-01-08T14:17:32.005 connecting to VAST node localhost:42000  importer=suri_0
INFO[0000] 2021-01-08T14:17:32.006 connecting to VAST node localhost:42000  importer=suri_1
INFO[0000] 2021-01-08T14:17:32.007 suricata-reader reads data from STDIN  importer=suri_0
INFO[0000] 2021-01-08T14:17:32.008 suricata-reader reads data from STDIN  importer=suri_1
INFO[0010] 2021-01-08T14:17:42.013 suricata-reader produced 10155 events at a rate of 1015 events/sec in 10.0s  importer=suri_1
INFO[0010] 2021-01-08T14:17:42.013 suricata-reader produced 10027 events at a rate of 1002 events/sec in 10.0s  importer=suri_0
INFO[0010] 2021-01-08T14:17:42.013 suricata-reader produced 9911 events at a rate of 991 events/sec in 10.01s  importer=suri_2
ERRO[0015] importer finished with error: signal: killed  importer=suri_1
INFO[0015] 2021-01-08T14:17:47.518 loaded configuration file: "/etc/vast/vast.yaml"  importer=suri_1
INFO[0015] 2021-01-08T14:17:47.528 connecting to VAST node localhost:42000  importer=suri_1
INFO[0015] 2021-01-08T14:17:47.530 suricata-reader reads data from STDIN  importer=suri_1
INFO[0025] 2021-01-08T14:17:57.534 suricata-reader produced 10025 events at a rate of 1003 events/sec in 9.1s  importer=suri_1
INFO[0030] 2021-01-08T14:18:02.014 suricata-reader produced 19953 events at a rate of 998 events/sec in 20.0s  importer=suri_0
INFO[0030] 2021-01-08T14:18:02.014 suricata-reader produced 20113 events at a rate of 1006 events/sec in 20.0s  importer=suri_2
...
```
The dead process (`suri_1`) is restarted instantly.