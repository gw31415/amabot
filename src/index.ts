import { Command, createFactory, Option } from "discord-hono";
import { renderMathSvg } from "./mathjax";

export const factory = createFactory<{ Bindings: Env }>();

export const handlers = [
  factory.command(
    new Command("tex", "Render math using Mathjax").options(
      new Option("math", "Mathjax expression").required(),
    ),
    async (c) => {
      c.var.math;
      const svg = await renderMathSvg(c.var.math, { display: true });
      return c.res(svg);
    },
  ),
];

export default factory.discord().loader(handlers);
