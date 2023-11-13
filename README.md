# Amabot

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

A high performance, simple Discord bot for science majors, written in Rust.

- [INVITE LINK](https://discord.com/api/oauth2/authorize?client_id=603145258188013568&permissions=0&scope=bot)

## List of commands

- `/tex`: Render math expression using Mathjax.

## Deploying

### Prerequisites

- Font: Install a font file to the system according to your environment.

  | Platform         | Font                   |
  | ---------------- | ---------------------- |
  | Macos            | `Hiragino Mincho ProN` |
  | Windows          | `Yu Mincho`            |
  | Linux and others | `Noto Serif CJK JP`    |

### Configuration

- `DISCORD_TOKEN` (Required): Environment variable that holds the discord token.

### Using Docker ([Fly.io](https://fly.io) and etc.)

Assuming deployment with Docker, a Feature `docker` for Docker is provided: for
example, the search for font files is made in the current directory to simplify
the procedure of installing fonts in a Distroless image.

Please read the [Dockerfile](./Dockerfile) for more information.

## License

This Program is licensed under [AGPL-3.0](./LICENSE).

## Acknowledgments

Thanks to [gaato](https://github.com/gaato): wrote a JavaScript to get math SVG
code using Mathjax. For more information, please see
[this project](https://github.com/gw31415/mathjax_svg).

## Author

gw31415 <git@amas.dev>
