


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