import { DOMParser } from "linkedom";

export function applyPaddingAndBackgroundToSvg(
  svg: string,
  { padding, background }: { padding: number; background: string },
): string {
  const parser = new DOMParser();
  const doc = parser.parseFromString(svg, "image/svg+xml");
  const svgEl = doc.documentElement;

  const parseLenPx = (v: string | null): number | null => {
    if (!v) return null;
    const m = String(v)
      .trim()
      .match(/^([0-9]+(?:\.[0-9]+)?)\s*(px)?$/i);
    return m ? Number(m[1]) : null;
  };

  // viewBox を優先して読む。無ければ width/height (px or unitless) から作る
  const viewBox: string = svgEl.getAttribute("viewBox");
  let minX = 0;
  let minY = 0;
  let vbW: number | null = null;
  let vbH: number | null = null;

  if (viewBox) {
    const parts = viewBox
      .trim()
      .split(/[,\s]+/)
      .map(Number);
    if (parts.length === 4 && parts.every((n) => Number.isFinite(n))) {
      [minX, minY, vbW, vbH] = parts;
    }
  }
  if (vbW == null || vbH == null) {
    const w = parseLenPx(svgEl.getAttribute("width"));
    const h = parseLenPx(svgEl.getAttribute("height"));
    if (w != null && h != null) {
      vbW = w;
      vbH = h;
      minX = 0;
      minY = 0;
      svgEl.setAttribute("viewBox", `0 0 ${vbW} ${vbH}`);
    }
  }

  // 背景 rect と translate 用 g を作る
  const rect = doc.createElementNS("http://www.w3.org/2000/svg", "rect", {});
  rect.setAttribute("fill", background);

  const g = doc.createElementNS("http://www.w3.org/2000/svg", "g", {});
  g.setAttribute("transform", `translate(${padding},${padding})`);

  if (vbW != null && vbH != null) {
    const newW = vbW + padding * 2;
    const newH = vbH + padding * 2;

    svgEl.setAttribute("viewBox", `${minX} ${minY} ${newW} ${newH}`);
    if (svgEl.hasAttribute("width")) svgEl.setAttribute("width", String(newW));
    if (svgEl.hasAttribute("height"))
      svgEl.setAttribute("height", String(newH));

    rect.setAttribute("x", String(minX));
    rect.setAttribute("y", String(minY));
    rect.setAttribute("width", String(newW));
    rect.setAttribute("height", String(newH));
  } else {
    // 寸法が取れないSVGでも「背景+translate」は入れる（サイズ拡張は諦める）
    rect.setAttribute("x", "0");
    rect.setAttribute("y", "0");
    rect.setAttribute("width", "100%");
    rect.setAttribute("height", "100%");
  }

  // rect を最背面に追加
  svgEl.insertBefore(rect, svgEl.firstChild);

  // rect以外を全て g に移動（defs等もまとめて移す：最小実装）
  let node = rect.nextSibling;
  while (node) {
    const next = node.nextSibling;
    g.appendChild(node);
    node = next;
  }
  svgEl.appendChild(g);

  return doc.documentElement.outerHTML;
}
