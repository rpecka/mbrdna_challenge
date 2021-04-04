# MBRDNA Chatbot

A chatbot server and client for remotely accessing a Mercedes-Benz connected vehicle.

# Running
I've hardcoded my own Houndify credentials into the project so that you can run with my rules and intent configurations.

I also hardcoded my Mercedes-Benz Connected Vehicle (MBCV) credentials in case you don't want to use your own, but you will be prompted to provide your own on startup if you wish.

```
git clone https://github.com/rpecka/mbrdna_challenge.git
cd mbrdna_challenge
```
Run the server:
```
go run ./cmd/server/...
```
Run the client:
```
go run ./cmd/client/...
```
