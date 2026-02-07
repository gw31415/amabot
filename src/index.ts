import { Command, createFactory, Embed, Option } from "discord-hono";
import { renderMathSvg } from "./mathjax";
import { svgToPngBlob } from "./svg2png";
import { applyPaddingAndBackgroundToSvg } from "./transformSvg";

export const factory = createFactory<{ Bindings: Env }>();

export const handlers = [
  factory.command(
    new Command("tex", "Render math using Mathjax").options(
      new Option("math", "Mathjax expression").required(),
    ),
    async (c) => {
      let svg = await renderMathSvg(c.var.math, { display: true });
      svg = applyPaddingAndBackgroundToSvg(svg, {
        padding: 20,
        background: "white",
      });
      const pngBlob = await svgToPngBlob(svg, {
        fitTo: { mode: "height", value: 100 },
      });

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
