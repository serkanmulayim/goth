UI

#Tested with 
#nvm v12.6.0
#npm 6.9.0

```
#build ui
cd src/fe/client
npm install
npm run build

cd ../admin
npm install
npm run build
cd ../..
```

BACKEND
```
export GOPATH=`pwd`
go get ./...
go build src/goth/main/main.go
```


#for atom ES6 and JSX highlights
apm install language-babel
npm install -g babel-cli
