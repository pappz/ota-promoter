# OTA-promoter

Over the air MicroPython application update tool for ESP8266 type microcomputer.

With this tool you can track and you can publish the changed files from 
your developer PC to your microchips via HTTP.

## Motivation
I designed this solution for decrease the testing-development period times during
the MicroPython application development.
In my personal cases I am not using or I can not using USB wire for upload the 
Python files to my devices during the development circles. 

## How does it work
The service registers some information about the promoted files.

- file names and full path of it
- generated unique hashes of the files
- version hash code of the current state of the promoted files

This service is listening on HTTP and with inotify are watching the modifications in a 
specified folder. Every time has occur a changes it will update the version information from the files
and it will generate a new version hash code.
With an HTTP client can get the current version code and the list of the available 
files. Based on the hash code the HTTP client can download the modified files.

## API

### Get current version code
```
$ curl -v http://192.168.0.10:9090/files/version
{"version":"8a841114726b9f327a6b94d6d129ec8588b5bdc7"}
```

### List of available files
```
$ curl  http://192.168.0.10:9090/files | python3 -m json.tool
{
    "files": [
        {
            "path": "README.txt",
            "checksum": "96f264583956281570cc591158c9371f8bba3736"
        }
    ],
    "version": "8a841114726b9f327a6b94d6d129ec8588b5bdc7"
}
```

### Download file by hash
```
curl  http://192.168.0.10:9090/files/96f264583956281570cc591158c9371f8bba3736
```
 
