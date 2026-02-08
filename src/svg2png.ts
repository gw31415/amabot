import { initWasm, Resvg } from "@resvg/resvg-wasm";
import resvg from "@resvg/resvg-wasm/index_bg.wasm";

let resvgReady: Promise<void> | null = null;
function init() {
  resvgReady ??= initWasm(resvg);
  return resvgReady;
}

export async function svgToPng(
  ...opts: ConstructorParameters<typeof Resvg>
): Promise<Uint8Array> {
  await init();
  const resvg = new Resvg(...opts);
  return resvg.render().asPng();
}
