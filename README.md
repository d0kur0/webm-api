# WEBM API

```
WEBM-API: Web server realisation for grabbing and share media from image-boards

Usage:
  webm-pwa [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       Start http server
  version     Print the version number of Hugo

Flags:
      --config string   config file (default is $CWD/.webm-pwa.json | $HOME/.webm-pwa.json)
  -h, --help            help for webm-pwa

Use "webm-pwa [command] --help" for more information about a command.
```

## Installation

// TODO

## HTTP Endpoints:

#### Get current server grabber schema: [request type: GET]
http://localhost:3000/schema
                                                                                          
#### Get all grabbed files: [request type: GET]
http://localhost:3000/files
                                                                                          
#### Get grabbed files by specific vendor and boards: [request type: POST]                        
http://localhost:3000/filesWithCondition

This request need body with condition struct:
```json
{ "<vendor name from grabber schema>": ["<board name1>", "<board name 2>"] }
```
for example:

```json
{ "2ch": ["b", "media", "vg"], "4chan": ["b"] }
```