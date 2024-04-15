all:
	@echo "no target."

run:
	go run .

docker:
	docker build --tag manenfu/medichat-be:latest .

docker-push:
	docker push manenfu/medichat-be:latest

mock:
	mockery --dir=./domain --outpkg=domainmocks --output=./mocks/domainmocks --name=DataRepository
	mockery --dir=./domain --outpkg=domainmocks --output=./mocks/domainmocks --name=AccountRepository
	mockery --dir=./domain --outpkg=domainmocks --output=./mocks/domainmocks --name=ResetPasswordTokenRepository
	mockery --dir=./domain --outpkg=domainmocks --output=./mocks/domainmocks --name=VerifyEmailTokenRepository
	
	mockery --dir=./domain --outpkg=domainmocks --output=./mocks/domainmocks --name=AccountService
	mockery --dir=./domain --outpkg=domainmocks --output=./mocks/domainmocks --name=GoogleService
	mockery --dir=./domain --outpkg=domainmocks --output=./mocks/domainmocks --name=OAuth2Service

	mockery --dir=./cryptoutil --outpkg=cryptomocks --output=./mocks/cryptomocks --name=JWTProvider 
	mockery --dir=./cryptoutil --outpkg=cryptomocks --output=./mocks/cryptomocks --name=OAuth2Provider 
	mockery --dir=./cryptoutil --outpkg=cryptomocks --output=./mocks/cryptomocks --name=PasswordHasher
	mockery --dir=./cryptoutil --outpkg=cryptomocks --output=./mocks/cryptomocks --name=RandomTokenProvider 
