mod entity;
mod util;

use std::borrow::Cow;

use entity::Error;
use poise::{
    serenity_prelude::{self as serenity, AttachmentType, Color},
    Event,
};

use crate::util::math_to_png;

type Data = ();
type Context<'a> = poise::Context<'a, Data, Error>;

/// Render math using Mathjax
#[poise::command(slash_command)]
async fn tex(
    ctx: Context<'_>,
    #[description = "Mathjax expression"] math: String,
) -> Result<(), Error> {
    let data: std::borrow::Cow<'_, [u8]> = Cow::Owned(math_to_png(math)?);
    let filename = "tex-result.png".to_string();
    let attachment = AttachmentType::Bytes {
        filename: filename.clone(),
        data,
    };
    ctx.send(|b| {
        b.embed(|embed| {
            embed
                .title("`tex` result:")
                .color(Color::DARK_GREEN)
                .attachment(filename)
        })
        .attachment(attachment)
    })
    .await?;
    Ok(())
}

#[tokio::main]
async fn main() {
    #[cfg(feature = "dotenv")]
    dotenvy::dotenv().unwrap();

    #[cfg(feature = "logger")]
    env_logger::init();

    let framework = poise::Framework::builder()
        .options(poise::FrameworkOptions {
            commands: vec![tex()],
            event_handler: |ctx, event, framework, data| {
                Box::pin(event_handler(ctx, event, framework, data))
            },
            on_error: |err| Box::pin(on_error(err)),
            ..Default::default()
        })
        .intents(serenity::GatewayIntents::non_privileged())
        .token(std::env::var("DISCORD_TOKEN").expect("missing DISCORD_TOKEN"))
        .setup(|_ctx, _ready, _framework| Box::pin(async move { Ok(()) }));
    framework.run().await.unwrap();
}

/// On Error handler
async fn on_error(err: poise::FrameworkError<'_, Data, Error>) {
    if let Some(ctx) = err.ctx() {
        use poise::FrameworkError::*;
        let _ = ctx
            .send(|b| match &err {
                Command { error, .. }
                | Setup { error, .. }
                | EventHandler { error, .. }
                | DynamicPrefix { error, .. } => {
                    #[cfg(feature = "logger")]
                    if error.log_required() {
                        log::error!("{:?}", err);
                    }
                    b.embed(error.get_embed())
                }
                _ => {
                    #[cfg(feature = "logger")]
                    {
                        log::error!("{:?}", err);
                    }
                    b.embed(|embed| {
                        embed
                            .color(Color::RED)
                            .title("Error")
                            .description("Unhandled error occured.")
                    })
                }.ephemeral(true)
            })
            .await;
    }
}

/// Automatic configuration of guild command
async fn event_handler(
    ctx: &serenity::Context,
    event: &Event<'_>,
    framework: poise::FrameworkContext<'_, Data, Error>,
    _data: &Data,
) -> Result<(), Error> {
    match event {
        Event::GuildMemberAddition { new_member } => {
            if framework.bot_id == new_member.user.id {
                poise::builtins::register_in_guild(
                    ctx,
                    &framework.options().commands,
                    new_member.guild_id,
                )
                .await?;
            }
        }
        Event::Ready { data_about_bot } => {
            for guild in &data_about_bot.guilds {
                poise::builtins::register_in_guild(ctx, &framework.options().commands, guild.id)
                    .await?;
            }
        }
        _ => (),
    }
    Ok(())
}
