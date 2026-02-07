import { register } from "discord-hono";
import { factory, handlers } from ".";

const DISCORD_TEST_GUILD_ID = process.env.DISCORD_TEST_GUILD_ID;

register(
  factory.getCommands(handlers),
  process.env.DISCORD_APPLICATION_ID,
  process.env.DISCORD_TOKEN,
  DISCORD_TEST_GUILD_ID.length > 0 ? DISCORD_TEST_GUILD_ID : undefined,
);
