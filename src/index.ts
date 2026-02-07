import { Command, createFactory } from "discord-hono";

export const factory = createFactory<{ Bindings: Env }>();

export const handlers = [
  factory.command(new Command("hello", "hello-world message"), (c) =>
    c.res("world!"),
  ),
];

export default factory.discord().loader(handlers);
