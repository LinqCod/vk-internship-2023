build:
	sudo docker build --tag numbers-telegram-bot .
run:
	sudo docker run numbers-telegram-bot

.PHONY: build, run
