import { Command, createFactory, Embed, Option } from "discord-hono";
import { renderMathSvg } from "./mathjax";
import { svgToPngBlob } from "./svg2png";

export const factory = createFactory<{ Bindings: Env }>();

export const handlers = [
  factory.command(
    new Command("tex", "Render math using Mathjax").options(
      new Option("math", "Mathjax expression").required(),
    ),
    async (c) => {
      c.var.math;
      const svg = await renderMathSvg(c.var.math, { display: true });
      const pngBlob = await svgToPngBlob(svg);
      const filename = "math.png";

      const msg = new Embed()
        .title("`tex` result:")
        .color(0x006400)
        .image({ url: `attachment://${filename}` });

      return c.res({ embeds: [msg] }, { blob: pngBlob, name: filename });
    },
  ),
];

export default factory.discord().loader(handlers);
