## Table Top Generator
Table Top Generator, or ttgen for short, is a CLI utility to facilitate in generating table top roleplaying game worlds. I am primarily focused on making it easier to build the base framework of a location, be it a farm or city, and then using a little personal creativity to customize the town in my session notes. There will be some weird sections of code in ttgen due to the fact that I'm also using this project to familiarize myself with Go.

## Setting up ttgen
For now ttgen uses [Cobra](https://github.com/spf13/cobra) to provide a uniform CLI experience regardless of which component you are interacting with. ttgen can be built with a simple `go build ttgen.go` and then interacted with via `./ttgen --help`. The available commands should have enough documentation to get you going.

## Future for ttgen
I'm sure this section will quickly become outdated, but for now the goals of ttgen are to provide a simple CLI that interacts with an API that can be used in future projects. Frankly even in early development I am finding dumping large towns/cities in text format isn't easily consumable meanwhile single character generation is perfect as is. I'd like to eventually link this to another project that enables all of this content to be consumed via a web browser, which I expect to be a slightly more palatable way to use this.

## Bugs and Suggestions
Please report bugs and suggestions to the [issues section](https://github.com/gomeler/LocationGenerator/issues) on github. I'll try to do my best at logging and resolving issues around new development.
