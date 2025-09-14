http://localhost:50151/api/chats?title=group&isVerified=true&offset=10&limit=50&sortBy=title&sortOrder=asc

### Read

curl -X GET http://localhost:50151/api/chats/018f3a8b-1b32-7299-e6f4-12398a7b6c5d
curl -X GET http://localhost:50151/api/chats
curl -X GET "http://localhost:50151/api/chats?chatId=01992530-617f-7546-bf68-431fef3cd3ff"

### Update

curl -X PATCH http://localhost:50151/api/chats/0199357b-8351-7a36-932d-18576ac281e0 \
-H "Content-Type: application/json" \
-d '{
"chatIds": ["chat_id_1002", "chat_id_1004", "chat_id_1006"]
"chatId": "0199357b-8351-7a36-932d-18576ac281e0",
"title": "sara",
}'

curl -X PATCH http://localhost:50151/api/chats/bulk-update \
-H "Content-Type: application/json" \
-d '{"ids": ["chat_id_1002", "chat_id_1004", "chat_id_1006"]}'


### Delete

curl -X DELETE http://localhost:50151/api/chats/018f3a8b-1b32-7292-b2d9-1237a7b8c8d2

curl -X POST http://localhost:50151/api/chats/bulk-delete \
-H "Content-Type: application/json" \
-d '{"ids": ["chat_id_1", "chat_id_2", "chat_id_3"]}'
