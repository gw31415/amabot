import { liteAdaptor } from "@mathjax/src/js/adaptors/liteAdaptor.js";
import { RegisterHTMLHandler } from "@mathjax/src/js/handlers/html.js";
import { TeX } from "@mathjax/src/js/input/tex.js";
import { mathjax } from "@mathjax/src/js/mathjax.js";
import { SVG } from "@mathjax/src/js/output/svg.js";
import type { OptionList } from "@mathjax/src/js/util/Options.js";

const adaptor = liteAdaptor();
RegisterHTMLHandler(adaptor);

const mathDocument = mathjax.document("", {
  InputJax: new TeX({ packages: ["base", "ams", "newcommand"] }),
  OutputJax: new SVG({ fontCache: "none" }),
});

export async function renderMathSvg(
  latex: string,
  options?: OptionList,
): Promise<string> {
  try {
    const svgNode = await mathDocument.convert(latex, options);
    const svgString = adaptor.outerHTML(svgNode);
    const svgTag = svgString.match(/<svg[^>]*>[\s\S]*<\/svg>/g)![0];

    if (svgString.includes("data-mjx-error")) {
      const titleMatch = svgString.match(/title="([^"]+)"/);
      const title = titleMatch ? titleMatch[1] : "MathJax error";
      throw new Error(title);
    }

    return svgTag;
  } catch (error) {
    console.error("MathJax rendering error:", error);
    throw new Error(
      `Failed to render Math: ${error instanceof Error ? error.message : String(error)}`,
    );
  }
}
