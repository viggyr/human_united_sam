# huManUnited

This is an initiative towards building a reward-based (Good Samaritan Points) online platform for unitedly confronting common and personal(individual) problems through collaboration.

The application is completely serverless, developed using the AWS serverless framework - AWS SAM, AWS Lambda, AWS APIGateway, DynamoDB, and s3. The lambda functions are developed using GoLang while Vue.js (https://github.com/ck090/AWS-Serverless-App-Hackathon) is used for the frontend. It's currently at the prototype stage, and can be accessed using https://master.d21skfmdtap7ea.amplifyapp.com/

Here are a couple of sample screenshots of the application.

<img width="1440" alt="image" src="https://user-images.githubusercontent.com/20017119/110473646-c43d1600-80ac-11eb-9080-80143c39b17c.png">

<img width="1439" alt="image" src="https://user-images.githubusercontent.com/20017119/110473790-eb93e300-80ac-11eb-87ac-256aeacc4307.png">


### Structure
```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── issues                      <-- Source code for a lambda function concerning issue management functionality
├── userlogin                   <-- Source code for a lambda function concerning user login/logout functionality
├── users                       <-- Source code for a lambda function concerning user management functionality
├── json                        <-- This has the static data for the prototype purpose
└── template.yaml               <-- Config file for defining the infrastructure (similar to AWS Cloudformation)
└── samconfig.toml              <-- Config file for deployment.
```
Different resources/functionalities (login, user management, etc.,) can be developed using different languages, but for time being only Go is being used. 
