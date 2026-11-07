# Skill Badges

[![Go][go-shield]][go-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

High-performance, lightweight HTTP service written in Go for dynamically generating customized SVG skill badge grids for GitHub profiles or personal websites.

<p align="center">
  <img align="center" src="./.github/example-dark.svg#gh-dark-mode-only"/>
  <img align="center" src="./.github/example-light.svg#gh-light-mode-only"/>
</p>

> [!NOTE]
> This project is a custom Go implementation based on the original idea by [tandpfun/skill-icons](https://github.com/tandpfun/skill-icons) and inspired by [LelouchFR/skill-icons](https://github.com/LelouchFR/skill-icons).

---

## Table Of Contents

1. [Examples](#examples)
   - [Show Your Skills](#show-your-skills)
   - [Themes](#themes)
   - [Badges Per Line](#badges-per-line)
   - [Size Presets](#size-presets)
   - [Centering Badges](#centering-badges)
4. [Built With](#built-with)
6. [Acknowledgements](#acknowledgements)
7. [License](#license)

---

## Examples

### Show Your Skills

Copy this line into your README. Change the list after `b=` to your own skills, separated by commas.

```md
[![My Skills](https://skill-badges.fly.dev/api/v1/badges?b=java,spring,go,js)](https://github.com/RomanGulevatiy/skill-badges)
```

[![My Skills](https://skill-badges.fly.dev/api/v1/badges?b=java,spring,go,js)](https://github.com/RomanGulevatiy/skill-badges)

### Themes

Some badges have a dark and a light version. The default theme is `dark`. Add `&theme=light` to switch.

```md
[![My Skills](https://skill-badges.fly.dev/api/v1/badges?b=java,spring,arch,neovim&theme=light)](https://github.com/RomanGulevatiy/skill-badges)
```

[![My Skills](https://skill-badges.fly.dev/api/v1/badges?b=java,spring,arch,neovim&theme=light)](https://github.com/RomanGulevatiy/skill-badges)

> [!NOTE]
> Badges without theme variants (docker, git, postman, golang, javascript, typescript) look the same in both themes.

### Badges Per Line

By default, you get 15 badges in a row. Add `&perline=` with a number from 1 to 32 to change it.
 
```md
[![My Skills](https://skill-badges.fly.dev/api/v1/badges?b=docker,git,github,gitlab,idea,goland&perline=3)](https://github.com/RomanGulevatiy/skill-badges)
```
 
[![My Skills](https://skill-badges.fly.dev/api/v1/badges?b=docker,git,github,gitlab,idea,goland&perline=3)](https://github.com/RomanGulevatiy/skill-badges)

### Size Presets

Choose a size with `&size=`: `sm`, `md` (default), `lg` or `xl`.

```md
[![My Skills](https://skill-badges.fly.dev/api/v1/badges?b=linux,debian,mongodb,postgres,redis&size=lg)](https://github.com/RomanGulevatiy/skill-badges)
```
 
[![My Skills](https://skill-badges.fly.dev/api/v1/badges?b=linux,debian,mongodb,postgres,redis&size=lg)](https://github.com/RomanGulevatiy/skill-badges)

### Centering Badges
 
The badges are a normal image, so you can center them like any other image.
 
```html
<p align="center">
  <a href="https://github.com/RomanGulevatiy/skill-badges">
    <img src="https://skill-badges.fly.dev/api/v1/badges?b=js,ts,vite,vue,pinia,tailwind"/>
  </a>
</p>
```
 
<p align="center">
  <a href="https://github.com/RomanGulevatiy/skill-badges">
    <img src="https://skill-badges.fly.dev/api/v1/badges?b=js,ts,vite,vue,pinia,tailwind"/>
  </a>
</p>

---

## Built With

- Go 1.27.2
- Standard `net/http`
- [caarlos0/env](https://github.com/caarlos0/env) - configuration from environment variables
- [joho/godotenv](https://github.com/joho/godotenv) - `.env` file support
- Docker & Docker Compose (multi-stage build)

---

## Acknowledgements

- [tandpfun/skill-icons](https://github.com/tandpfun/skill-icons) - original idea and visual layout concept.
- [LelouchFR/skill-icons](https://github.com/LelouchFR/skill-icons) - reference implementation in Go.

---

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

[go-shield]: https://img.shields.io/github/go-mod/go-version/RomanGulevatiy/skill-badges
[go-url]: https://github.com/RomanGulevatiy/skill-badges/blob/main/go.mod
[issues-shield]: https://img.shields.io/github/issues/RomanGulevatiy/skill-badges?label=Issues
[issues-url]: https://github.com/RomanGulevatiy/skill-badges/issues
[license-shield]: https://img.shields.io/github/license/RomanGulevatiy/skill-badges?label=License
[license-url]: https://github.com/RomanGulevatiy/skill-badges/blob/main/LICENSE
