FROM amazon/dynamodb-local
EXPOSE 8000/tcp
WORKDIR /home/dynamodblocal
VOLUME ["dynamodata:/home/dynamodblocal"]
CMD ["-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "."]
