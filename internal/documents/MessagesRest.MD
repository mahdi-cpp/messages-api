
## Messages REST API

## Create
curl -X POST http://localhost:50151/api/messages -H "Content-Type: application/json" \
-d '{
"userID": "01991bc4-faad-7b70-aedc-f20ea4146898",
"chatId": "018f3a8b-1b32-7295-a2c7-87654b4d4567",
"content": "Mahdi abdolmaleki 55"
}'

### Read
curl -X GET http://localhost:50151/api/messages?chatId=01992530-617f-7546-bf68-431fef3cd3ff


## Update
curl -X PATCH http://localhost:50151/api/messages -H "Content-Type: application/json" -d '{
"userID": "01991bc4-faad-7b70-aedc-f20ea4146898",
"chatId": "018f3a8b-1b32-7295-a2c7-87654b4d4567",
"messageId": "01991bc4-faad-7b70-aedc-f20ea4146898",
"content": "Mahdi abdolmaleki 55"
}'







