BINARY_NAME := bootstrap
FUNCTION_NAME := resizeImage
ZIP_NAME := resizeImageLambdaFunction.zip
AWS_REGION := ap-northeast-1

all: build

build:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY_NAME) cmd/lambda-handler/main.go

package: build
	zip $(ZIP_NAME) $(BINARY_NAME)

deploy: package
	aws lambda update-function-code --function-name $(FUNCTION_NAME) --zip-file fileb://$(ZIP_NAME) --region $(AWS_REGION)

clean:
	rm -f $(BINARY_NAME) $(ZIP_NAME)
