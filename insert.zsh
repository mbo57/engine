#!/bin/zsh

curl -X POST -H "Content-Type: application/json" http://localhost:1323/test
while read line
do
    echo $line
    curl -X POST -H "Content-Type: application/json" -d '{"text":"'"$line"'"}' http://localhost:1323/test/_doc
done < ./sample3.txt
