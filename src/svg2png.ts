import { initWasm, Resvg } from "@resvg/resvg-wasm";
import resvg from "@resvg/resvg-wasm/index_bg.wasm";

let resvgReady: Promise<void> | null = null;
function init() {
  resvgReady ??= initWasm(resvg);
  return resvgReady;
}

export async function svgToPngBlob(
  ...opts: ConstructorParameters<typeof Resvg>
): Promise<Blob> {
  await init();
  const resvg = new Resvg(...opts);

  const pngData = resvg.render();
  const pngUint8 = pngData.asPng();
  return new Blob([pngUint8], { type: "image/png" });
}
