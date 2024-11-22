#!/bin/zsh
while read line
do
    echo $line
    curl -X POST -H "Content-Type: application/json" -d '{"text":"'"$line"'"}' http://localhost:1323/test/_doc
done < ./sample.txt
