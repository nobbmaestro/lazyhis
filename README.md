# LazyHis

[![GitHub Release](https://img.shields.io/github/v/release/nobbmaestro/lazyhis)](github-release)
[![GitHub last commit](https://img.shields.io/github/last-commit/nobbmaestro/lazyhis/development)](github-last-commit)
[![GitHub commits since](https://img.shields.io/github/commits-since/nobbmaestro/lazyhis/0.1.0-beta/development)](githut-commits-since)
![License](https://img.shields.io/github/license/nobbmaestro/lazyhis)

LazyHis: `FIXME`

`TODO: Add short demo video`

## Table of contents

- [Installation](#installation)
  - [Brew](#brew)
  - [Manual](#manual)
  - [Configure zsh](#configure-zsh)
- [Feature Roadmap](#feature-roadmap)
- [Customization](#customization)
- [Alternatives](#alternatives)

## Installation

### Brew (Recommended)

`TODO: Finalize this`

### Manual

```sh
git clone git@github.com:nobbmaestro/lazyhis.git
cd lazyhis
make
```

`TODO: Finalize this`

## Configure zsh

```sh
echo 'eval "$(lazyhis init zsh)"' >> ~/.zshrc
```

## Feature Roadmap

- [ ] Customizable GUI theme
- [ ] Support for inline GUI mode
- [ ] Edit history entries via GUI
- [ ] Delete history entries via GUI
- [ ] Delete selected history entries via GUI
- [ ] Copy to clipboard via GUI
- [ ] Fuzzy-finder search strategy in GUI
- [ ] Filter history by context via GUI
- [ ] Add doctor CLI command for verifying shell configuration
- [ ] Add prune CLI command for removing history based on ignore pattern
- [ ] Add generate shell-completions CLI command
- [ ] Add export CLI command for exporting to HISTFILE
- [ ] Add support for command execution duration
- [ ] Customizable keybindings

## Customization

Check out the [configuration docs](docs/config.md).

## Alternatives

If you find that `lazyhis` does not quite satisfy your needs, following may be a better fit:

- [Atuin](https://github.com/atuinsh/atuin)
