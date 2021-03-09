# huManUnited

This is an initiative towards building a reward-based (Good Samaritan Points) online platform for unitedly confronting common and personal(individual) problems through collaboration.

The application is completely serverless, using AWS serverless stack - AWS SAM, AWS Lambda, AWS APIGateway, DynamoDB, and s3. It's at the prototype stage, and can be accessed using https://master.d21skfmdtap7ea.amplifyapp.com/

Here are a couple of sample screenshots of the application 

<img width="1440" alt="image" src="https://user-images.githubusercontent.com/20017119/110473646-c43d1600-80ac-11eb-9080-80143c39b17c.png">

<img width="1439" alt="image" src="https://user-images.githubusercontent.com/20017119/110473790-eb93e300-80ac-11eb-87ac-256aeacc4307.png">


### Structure
```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── hello-world                 <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
└── template.yaml
```
