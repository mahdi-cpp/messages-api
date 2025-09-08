


````
curl -X POST "http://localhost:50151/api/chats" \
-H "Content-Type: application/json" \
-d '{
"userId": "bbbbbb",
"memebers": ["bbbbbb", "ccccccc"]
}'
````
````
curl -X GET "http://localhost:50151/api/chats?offset=2&limit=20"
````

http://localhost:50151/api/chats?title=group&isVerified=true&offset=10&limit=50&sortBy=title&sortOrder=asc


### Updates

curl -X PATCH http://localhost:50151/api/chats/bulk-update \
-H "Content-Type: application/json" \
-d '{"ids": ["chat_id_1002", "chat_id_1004", "chat_id_1006"]}'

curl -X POST http://localhost:50151/api/chats/018f3a8b-1b32-729a-f7e5-5467c1b2d3e4 \
-H "Content-Type: application/json" \
-d '{
"userID": "new_user_id_1",
"role": "member",
"customTitle": "Team Member",
"isActive": true,
"lastActive": "2025-09-06T12:00:00Z",
"joinedAt": "2025-09-01T10:00:00Z"
}'

### Read
curl -X GET http://localhost:50151/api/chats/018f3a8b-1b32-7299-e6f4-12398a7b6c5d
curl -X GET http://localhost:50151/api/chats//018f3a8b-1b32-7299-e6f4-12398a7b6c5d/messages/019928c6-9e55-7479-884e-92cc2e244eb2
curl -X GET http://localhost:50151/api/chats
curl -X GET "http://localhost:50151/api/chats?chatId=01992530-617f-7546-bf68-431fef3cd3ff"

### Update

### Delete

curl -X DELETE http://localhost:50151/api/chats/018f3a8b-1b32-729a-f7e5-5467c1b2d3e4

curl -X POST http://localhost:50151/api/chats/bulk-delete \
-H "Content-Type: application/json" \
-d '{"ids": ["chat_id_1", "chat_id_2", "chat_id_3"]}'




### ---------------------------------------------------------------
## Messages

### Read

curl -X GET http://localhost:50151/api/messages/019928c3-895d-7286-9cd3-919cd4ddf9ea

curl -X GET http://localhost:50151/api/messages







