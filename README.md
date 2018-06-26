# FeedAll
A feed reader writen in Golang .

## About
FeedAll is a pratice project writen in Golang.
I choose Macaron as backend framwork and mongodb to storage data just because they can help me learn more :)

## Requirement
* go1.10+
* mongodb3.6+

## Install

### Get feedall
```bash
git clone https://github.com/looyun/feedall.git
```

### Create mongodb user
```
mongo
>>db.createUser({user:"username",pwd:"password",roles:[{role:"dbOwner",db:"feedall"}]})
```

### Start feedall
```
cd /path-of-feedall/
go run main.go
```
