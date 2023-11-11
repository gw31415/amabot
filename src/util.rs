use std::sync::OnceLock;

use anyhow::Context as _;
use mathjax_svg::convert_to_svg;
use resvg::usvg::{
    self,
    fontdb::{self, Database},
    Tree, TreeParsing, TreeTextToPath,
};
use tiny_skia::{Pixmap, PixmapPaint, Transform};

/// Get code block from string
pub fn code_block(string: impl AsRef<str>) -> String {
    let mut out = "```".to_string();
    out.push_str(&string.as_ref().replace("```", "`\u{200B}``"));
    if out.ends_with('`') {
        out.push('\u{200B}');
    }
    out.push_str("```");
    out
}

/// Convert math to PNG
pub fn math_to_png(math: impl AsRef<str>) -> Result<Vec<u8>, crate::Error> {
    /// The height of the PNG
    const HEIGHT: u32 = 100;
    /// Padding size
    const PADDING: u32 = 20;
    /// Default font-family for <text> tag
    #[cfg(target_os = "macos")]
    const FONT_FAMILY: &str = "Hiragino Mincho ProN";
    #[cfg(target_os = "windows")]
    const FONT_FAMILY: &str = "Yu Mincho";
    #[cfg(not(any(target_os = "macos", target_os = "windows")))]
    const FONT_FAMILY: &str = "Noto Serif CJK JP";
    /// Font database: only needs to be initialized once
    static FONTDB: OnceLock<Database> = OnceLock::new();

    let svg = convert_to_svg(math)?;
    let png = {
        let image = {
            // Convert to Pixmap
            let svg_data = svg.into_bytes();
            let rtree = {
                let opt = usvg::Options::default();

                let mut tree = Tree::from_data(&svg_data, &opt)?;
                tree.convert_text(FONTDB.get_or_init(|| {
                    let mut fdb = fontdb::Database::new();
                    fdb.load_system_fonts();
                    // Set default serif font
                    fdb.set_serif_family(FONT_FAMILY);
                    fdb
                }));
                resvg::Tree::from_usvg(&tree)
            };

            // Vertical length is scaled to be HEIGHT
            let (mut math_pix, scale_x, scale_y) = {
                let original_size = rtree.size;
                let target_size = original_size
                    .to_int_size()
                    .scale_to_height(HEIGHT)
                    .context("scaling Pixmap")?;
                (
                    tiny_skia::Pixmap::new(target_size.width(), target_size.height())
                        .context("creating new Pixmap to draw svg in")?,
                    target_size.width() as f32 / original_size.width(),
                    target_size.height() as f32 / original_size.height(),
                )
            };
            rtree.render(
                tiny_skia::Transform::from_scale(scale_x, scale_y),
                &mut math_pix.as_mut(),
            );
            math_pix
        };

        let image = {
            // Add padding and white background
            let mut background =
                Pixmap::new(PADDING * 2 + image.width(), PADDING * 2 + image.height())
                    .context("creating new Pixmap for padding")?;
            background.fill(tiny_skia::Color::WHITE);
            background.draw_pixmap(
                PADDING as i32,
                PADDING as i32,
                image.as_ref(),
                &PixmapPaint::default(),
                Transform::default(),
                None,
            );
            background
        };

        image.encode_png()?
    };
    Ok(png)
}
