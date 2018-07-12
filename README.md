# feedall
A feed reader writen in Golang .

## About
feedall is a pratice project writen in Golang.
I choose Macaron as backend framwork and mongodb to storage data just because they can help me learn more :)

## Requirement
* go1.10+
* [mongodb3.6+](https://docs.mongodb.com/v3.6/installation/)

## Install

### Get feedall
```bash
go get github.com/looyun/feedall
```

### Create mongodb user
```
mongo
>>use feedall
>>db.createUser({user:"username",pwd:"password",roles:[{role:"dbOwner",db:"feedall"}]})
```

### Start feedall
```
cd $GOPATH/src/github.com/looyun/feedall
go run main.go
```
