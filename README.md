# Guildmaster

Guildmaster is a proof of concept Discord bot written entirely in Go. The purpose of the bot is to track the activity of all users in a server for leveling purposes, including messages sent and time spent within a voice channel. The bot will also notify users upon leveling up, and provide users with a leaderboard to track activity across the server.

## Development

### Requirements

- Docker Engine and Compose
- Golang version 1.22 or higher
- Discord Developer Account with a bot token

### Starting the development server

Prior to starting the application, certain environment variables need to be set. Create a `.env` file in the root of the project and set the `DISCORD_TOKEN` environment variable.

Running `sudo docker compose up` will start the development server. Any changes made to Go files will automatically recompile the bot.
