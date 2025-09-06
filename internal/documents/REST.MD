


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
curl -X GET http://localhost:50151/api/chats/018f3a8b-1b32-7293-c1d4-8765f4d1e2f3

curl -X GET http://localhost:50151/api/chats

### Update

### Delete

curl -X DELETE http://localhost:50151/api/chats/018f3a8b-1b32-729a-f7e5-5467c1b2d3e4

curl -X POST http://localhost:50151/api/chats/bulk-delete \
-H "Content-Type: application/json" \
-d '{"ids": ["chat_id_1", "chat_id_2", "chat_id_3"]}'















