# MBRDNA Chatbot

A chatbot server and client for remotely accessing a Mercedes-Benz connected vehicle.

![Screen Shot 2021-04-04 at 5 09 40 PM](https://user-images.githubusercontent.com/13580126/113521612-91315980-9568-11eb-8023-9b90663ed01f.png)

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
