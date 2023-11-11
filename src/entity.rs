use poise::serenity_prelude::{self as serenity, Color};
use resvg::usvg;

use crate::util::code_block;

/// Error (to be resolved during execution)
#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error(transparent)]
    LaTeX(#[from] mathjax_svg::Error),
    #[error(transparent)]
    Png(#[from] png::EncodingError),
    #[error(transparent)]
    Svg(#[from] usvg::Error),
    #[error(transparent)]
    Serenity(#[from] serenity::Error),
    #[error(transparent)]
    Other(#[from] anyhow::Error),
}

impl Error {
    pub fn get_embed(
        &self,
    ) -> impl FnOnce(&mut serenity::CreateEmbed) -> &mut serenity::CreateEmbed {
        let msg = match self {
            LaTeX(err) => format!("LaTeX Error:{}", code_block(err.to_string())),
            _ => format!(
                "Unknown error:{}Please contact the owner.",
                code_block(self.to_string())
            ),
        };
        use Error::*;
        |embed| embed.title("Error").color(Color::RED).description(msg)
    }
    pub fn log_required(&self) -> bool {
        !matches!(self, Error::LaTeX(_))
    }
}
