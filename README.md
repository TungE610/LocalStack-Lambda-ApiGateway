# LocalStack-Lambda-ApiGateway

GOOS=linux GOARCH=amd64 go build -o main main.go

zip function.zip main

aws lambda create-function \
  --function-name ConvertCurrency \
  --runtime go1.x \
  --handler main \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --endpoint-url=http://localhost:4566

  aws lambda list-functions --endpoint-url=http://localhost:4566


//create api gateway

aws apigateway create-rest-api --name "CurrencyConversionAPI" --endpoint-url http://localhost:4566

// create resource and method
aws apigateway create-resource --rest-api-id <api-id> --parent-id <parent-id> --path-part "convert" --endpoint-url http://localhost:4566
aws apigateway put-method --rest-api-id <api-id> --resource-id <resource-id> --http-method POST --authorization-type NONE --endpoint-url http://localhost:4566

//Link the resource method to your Lambda function:
aws apigateway put-integration --rest-api-id <api-id> --resource-id <resource-id> --http-method POST --integration-http-method POST --type AWS_PROXY --uri arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:ConvertCurrency/invocations --endpoint-url http://localhost:4566

//deploy the API
aws apigateway create-deployment --rest-api-id <api-id> --stage-name prod --endpoint-url http://localhost:4566

aws lambda invoke \
  --function-name ConvertCurrency \
  --cli-binary-format raw-in-base64-out \
  --payload '{"amount": 100, "from": "USD", "to": "EUR"}' \ 
  --endpoint-url http://localhost:4566
