curl -X POST http://localhost:1323/test
curl -X POST http://localhost:1323/test/_doc -d '{"text": "She read books"}'
curl -X POST http://localhost:1323/test/_doc -d '{"text": "She read books in library"}'
curl -X POST http://localhost:1323/test/_doc -d '{"text": "I read books and she read books"}'
curl -X POST http://localhost:1323/test/_doc -d '{"text": "I play games"}'
curl -X POST http://localhost:1323/test/_doc -d '{"text": "I play games and read books"}'
curl -X POST http://localhost:1323/test/_doc -d '{"text": "I like to cook"}'
curl -X GET http://localhost:1323/test/_search -d '{"text": "read books"}'
