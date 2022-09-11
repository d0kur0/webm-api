# WEBM API 

## HTTP Endpoints:                                                                                

### Get current server grabber schema: [request type: GET]
http://localhost:3000/schema
                                                                                          
### Get all grabbed files: [request type: GET]
http://localhost:3000/files
                                                                                          
### Get grabbed files by specific vendor and boards: [request type: POST]                        
http://localhost:3000/filesWithCondition

This request need body with condition struct:
```json
{ "<vendor name from grabber schema>": ["<board name1>", "<board name 2>"] }
```
for example:

```json
{ "2ch": ["b", "media", "vg"], "4chan": ["b"] }
```